package scheduler

import (
	"time"
	"fmt"
	"github.com/maxpowel/dislet"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/jinzhu/gorm"
	"math/rand"
	"strconv"
)


var kernel *dislet.Kernel

func consumptionSignature(lineId uint) *tasks.Signature {
	return 	&tasks.Signature{
		Name: "consumptionPeriodic",
		Args: []tasks.Arg{
			{
				Type:  "uint",
				Value: lineId,
			},
		},
	}
}

type MysqlJob struct {
	gorm.Model
	LastResult string
	Period uint
	LastRun time.Time
	NextRun time.Time
	ContinueRunning bool
	MaxPeriod uint
	MinPeriod uint
	IncrementRatio float32   `sql:"type:decimal(4,2);"`
	DecrementRatio float32   `sql:"type:decimal(4,2);"`
}

func loadJob(lineId uint) (*Job, error) {
	db := kernel.Container.MustGet("database").(*gorm.DB)
	mj := MysqlJob{}
	db.First(&mj, lineId)
	if mj.ID == 0 {
		fmt.Sprintf("job %d does not exists", lineId)
		return nil, fmt.Errorf("job %d does not exists", lineId)
	}
	j := Job{
		LineId: lineId,
		LastResult: mj.LastResult,
		Period: time.Duration(mj.Period),
		LastRun: mj.LastRun,
		NextRun: mj.NextRun,
		ContinueRunning: mj.ContinueRunning,
		MaxPeriod: time.Duration(mj.MaxPeriod),
		MinPeriod: time.Duration(mj.MinPeriod),
		IncrementRatio: mj.IncrementRatio,
		DecrementRatio: mj.DecrementRatio,
	}

	return &j, nil
}

func updateJob(job *Job) (error){
	db := kernel.Container.MustGet("database").(*gorm.DB)
	mj := MysqlJob{}
	db.First(&mj, job.LineId)
	mj.LastRun = job.LastRun
	mj.NextRun = job.NextRun

	mj.ContinueRunning = job.ContinueRunning
	mj.Period = uint(job.Period)
	mj.MaxPeriod = uint(job.MaxPeriod)
	mj.MinPeriod = uint(job.MinPeriod)
	mj.IncrementRatio = job.IncrementRatio
	mj.DecrementRatio = job.DecrementRatio
	mj.LastResult = job.LastResult

	job.LastRun = mj.LastRun
	db.Save(&mj)

	return nil
}

func Consumption(lineId uint) (error) {
	fmt.Println("ENTRANDO")
	j, err := loadJob(lineId)
	if err != nil {
		return err
	}

	if j.ContinueRunning {
		// HACER EL TRABAJO
		res := strconv.Itoa(rand.Intn(3))

		fmt.Println("RESULTADO ES ", res)
		j.scheduleNextRun(res)
	} else {
		fmt.Println("Planificacion cancelada para ", j.LineId)
	}

	err = updateJob(j)
	return err
}

type Job struct {
	LineId uint
	LastResult string
	Period time.Duration
	LastRun time.Time
	NextRun time.Time
	ContinueRunning bool
	MaxPeriod time.Duration
	MinPeriod time.Duration
	IncrementRatio float32
	DecrementRatio float32
}

func (job *Job) scheduleNextRun(result string) {
	fmt.Println("LAST", job.LastResult, "NEW", result)
	if job.LastResult == result {
		// El resultado es el mismo, entonces incrementar tiempo
		fmt.Println("Aumentar tiempo", result, job.LastResult, job.Period)
		//Hasta un maximo
		if job.Period < job.MaxPeriod {
			job.Period = time.Duration(float32(job.Period.Nanoseconds()) * job.IncrementRatio)
			if job.Period > job.MaxPeriod {
				job.Period = job.MaxPeriod
			}
		} else if job.Period > job.MaxPeriod {
			job.Period = job.MaxPeriod
		}
		fmt.Println("Aumentado", job.Period)
	} else {
		// Ha habido cambios reducir tiempo
		fmt.Println("Reducir tiempo", result, job.LastResult, job.Period)
		if job.Period > job.MinPeriod {
			job.Period = time.Duration(float32(job.Period.Nanoseconds()) / job.DecrementRatio)
			if job.Period < job.MinPeriod {
				job.Period = job.MinPeriod
			}
		} else if job.Period < job.MinPeriod {
			job.Period = job.MinPeriod
		}
		fmt.Println("Reducido", job.Period)

	}

	job.LastResult = result
	job.LastRun = time.Now()
	job.NextRun = job.LastRun.Add(job.Period)
	m := kernel.Container.MustGet("machinery").(*machinery.Server)

	//eta := time.Now().UTC().Add(time.Second * 5)
	if job.ContinueRunning {
		fmt.Println("Planificado a ", job.Period)
		eta := time.Now().UTC().Add(job.Period)
		signature := consumptionSignature(job.LineId)
		signature.ETA = &eta
		_, err := m.SendTask(signature)
		if err != nil {
			fmt.Println("ERROR", err)
		} else {
			fmt.Println("RSEULTADO")
			//fmt.Println(asyncResult.Get(time.Second * 10))
		}
	} else {
		fmt.Println("Planificacion cancelada para ", job.LineId)
	}
}

func startSchedule(lineId uint) (error){
	db := kernel.Container.MustGet("database").(*gorm.DB)
	m := kernel.Container.MustGet("machinery").(*machinery.Server)

	j, err := loadJob(lineId)
	//No existe, entonces lo creamos (con los valores por defecto)
	if err != nil {
		mj := MysqlJob{
			DecrementRatio: 1.5,
			IncrementRatio: 1.5,
			Period: uint(5 * time.Second),
			ContinueRunning: true,
			MinPeriod: uint(2 * time.Second),
			MaxPeriod: uint(20 * time.Second),
		}
		mj.ID = lineId
		db.Save(&mj)
		// Check that everything is ok
		j, err = loadJob(lineId)
		if err != nil {
			return err
		}
	} else {
		//Simplemente lo activamos
		if j.ContinueRunning == false {
			j.ContinueRunning = true
			updateJob(j)
		} else {
			return fmt.Errorf("Already running")
		}
	}

	_, err = m.SendTask(consumptionSignature(j.LineId))
	if err != nil {
		return err
	} else {
		return nil
	}
	/*if err != nil {
		fmt.Println("ERROR", err)
	} else {
		fmt.Println("RSEULTADO")
		//fmt.Println(asyncResult.Get(time.Second * 10))
	}*/
}

func stopSchedule(lineId uint) (error){

	j, err := loadJob(lineId)
	if err != nil {
		return err
	}

	j.ContinueRunning = false
	updateJob(j)
	return nil
}


func Bootstrap(k *dislet.Kernel) {
	//mapping := k.Config.Mapping

	kernel = k
	var baz dislet.OnKernelReady = func(k *dislet.Kernel){
		//r := k.Container.MustGet("redis").(*redis.Client)
		db := kernel.Container.MustGet("database").(*gorm.DB)
		db.AutoMigrate(&MysqlJob{})

		m := k.Container.MustGet("machinery").(*machinery.Server)

		m.RegisterTasks(map[string]interface{}{
			//"add":      Add,
			"consumptionPeriodic": Consumption,
		})
		startSchedule(8)

	}
	k.Subscribe(baz)
}

// Ver trabajos que deberían estar funcionando pero no lo están
// select continue_running, last_run, period, next_run from mysql_jobs WHERE next_run < DATE_SUB(now(), INTERVAL 2 HOUR) and continue_running = 1;