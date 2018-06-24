package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/giagiannis/data-profiler/core"
)

// TaskEngine is deployed once for the server's lifetime and keeps the tasks
// which are executed.
type TaskEngine struct {
	Tasks []*Task
	lock  sync.Mutex
}

// NewTaskEngine initializes a new TasEngine object. Called once per server deployment.
func NewTaskEngine() *TaskEngine {
	te := new(TaskEngine)
	te.lock = *new(sync.Mutex)
	return te
}

// Submit appends a new task to the task engine and initializes its execution.
func (e *TaskEngine) Submit(t *Task) {
	if t == nil {
		return
	}
	e.lock.Lock()
	go t.Run()
	e.Tasks = append(e.Tasks, t)
	e.lock.Unlock()
}

// Task is the primitive struct that represents a task of the server.
type Task struct {
	Status      string
	Started     time.Time
	Duration    float64
	Description string
	Dataset     *ModelDataset
	fnc         func() error
}

// Run is responsible to execute to task's method and update the task status
// accordingly.
func (t *Task) Run() {
	t.Started = time.Now()
	t.Status = "RUNNING"
	err := t.fnc()
	if err != nil {
		t.Status = "ERROR - " + err.Error()
	} else {
		t.Status = "DONE"
		t.Duration = time.Since(t.Started).Seconds()
	}
}

// NewSMComputationTask initializes a new Similarity Matrix computation task.
func NewSMComputationTask(datasetID string, conf map[string]string) *Task {
	task := new(Task)
	dts := modelDatasetGetInfo(datasetID)
	task.Dataset = dts
	task.Description = fmt.Sprintf("SM Computation for %s, type %s\n",
		dts.Name, conf["estimatorType"])
	task.fnc = func() error {
		datasets := core.DiscoverDatasets(dts.Path)
		estType := *core.NewDatasetSimilarityEstimatorType(conf["estimatorType"])
		est := core.NewDatasetSimilarityEstimator(estType, datasets)
		est.Configure(conf)
		if conf["popPolicy"] == "aprx" {
			pop := new(core.DatasetSimilarityPopulationPolicy)
			pop.PolicyType = core.PopulationPolicyAprx
			val, err := strconv.ParseFloat(conf["popParameterValue"], 64)
			if err != nil {
				log.Println(err)
			}
			if conf["popParameter"] == "popCount" {
				pop.Parameters = map[string]float64{"count": val}
			}
			if conf["popParameter"] == "popThreshold" {
				pop.Parameters = map[string]float64{"threshold": val}
			}
			est.SetPopulationPolicy(*pop)
		}
		err := est.Compute()
		if err != nil {
			return err
		}
		sm := est.SimilarityMatrix()
		//		var smID string
		if sm != nil {
			//smID =
			modelSimilarityMatrixInsert(datasetID, sm.Serialize(), est.Serialize(), conf)
		}
		//modelEstimatorInsert(datasetID, smID, est.Serialize(), conf)
		return nil
	}
	return task
}

// NewMDSComputationTask initializes a new Multidimensional Scaling execution task.
func NewMDSComputationTask(smID, datasetID string, conf map[string]string) *Task {
	smModel := modelSimilarityMatrixGet(smID)
	if smModel == nil {
		log.Println("SM not found")
		return nil
	}

	cnt, err := ioutil.ReadFile(smModel.Path)
	if err != nil {
		log.Println(err)
	}
	sm := new(core.DatasetSimilarityMatrix)
	sm.Deserialize(cnt)
	k, err := strconv.ParseInt(conf["k"], 10, 64)
	if err != nil {
		log.Println(err)
	}

	dat := modelDatasetGetInfo(datasetID)
	task := new(Task)
	task.Dataset = dat
	task.Description = fmt.Sprintf("MDS Execution for %s with k=%d\n",
		dat.Name, k)
	task.fnc = func() error {
		mds := core.NewMDScaling(sm, int(k), Conf.Scripts.MDS)
		err = mds.Compute()
		if err != nil {
			return err
		}
		gof := fmt.Sprintf("%.5f", mds.Gof())
		stress := fmt.Sprintf("%.5f", mds.Stress())
		modelCoordinatesInsert(mds.Coordinates(), dat.ID, conf["k"], gof, stress, smID)
		return nil
	}
	return task
}

// NewOperatorRunTask initializes a new operator execution task.
func NewOperatorRunTask(operatorID string) *Task {
	m := modelOperatorGet(operatorID)
	if m == nil {
		log.Println("Operator was not found")
		return nil
	}
	dat := modelDatasetGetInfo(m.DatasetID)
	for _, f := range modelDatasetGetFiles(dat.ID) {
		dat.Files = append(dat.Files, dat.Path+"/"+f)
	}
	task := new(Task)
	task.Description = fmt.Sprintf("%s evaluation", m.Name)
	task.Dataset = dat
	task.fnc = func() error {
		eval, err := core.NewDatasetEvaluator(core.OnlineEval,
			map[string]string{
				"script":  m.Path,
				"testset": "",
			})
		if err != nil {
			log.Println(err)
		}
		scores := core.NewDatasetScores()
		for _, f := range dat.Files {
			s, err := eval.Evaluate(f)
			if err != nil {
				log.Println(err)
			} else {
				scores.Scores[path.Base(f)] = s
			}
		}
		cnt, _ := scores.Serialize()
		modelOperatorScoresInsert(operatorID, cnt)
		return nil
	}
	return task
}

// NewModelTrainTask generates a new ML for a given (dataset,operator) combination,
// according to the user-specified parameters.
func NewModelTrainTask(datasetID, operatorID string, sr float64,
	modelType string,
	coordinatesID, mlScript string,
	matrixID, k string) *Task {
	m := modelDatasetGetInfo(datasetID)
	task := new(Task)
	task.Description = fmt.Sprintf("Model training (%s for %s)", path.Base(mlScript), m.Name)
	task.Dataset = m
	task.fnc = func() error {
		datasets := core.DiscoverDatasets(m.Path)
		o := modelOperatorGet(operatorID)
		var evaluator core.DatasetEvaluator
		var err error
		if o.ScoresFile != "" {
			evaluator, err = core.NewDatasetEvaluator(core.FileBasedEval, map[string]string{"scores": o.ScoresFile})
		} else {
			evaluator, err = core.NewDatasetEvaluator(core.OnlineEval, map[string]string{"script": o.Path, "testset": ""})
		}
		if err != nil {
			log.Println(err)
			return err
		}
		t := core.NewModelerType(modelType)
		modeler := core.NewModeler(t, datasets, sr, evaluator)
		var conf map[string]string
		if t == core.ScriptBasedModelerType {
			c := modelCoordinatesGet(coordinatesID)
			conf = map[string]string{"script": mlScript, "coordinates": c.Path}
		} else if t == core.KNNModelerType {
			m := modelSimilarityMatrixGet(matrixID)
			conf = map[string]string{"k": k, "smatrix": m.Path}
		}
		modeler.Configure(conf)
		err = modeler.Run()
		if err != nil {
			log.Println(err)
			return err
		}
		// serialze appxValues
		var cnt [][]float64
		cnt = append(cnt, modeler.AppxValues())
		appxBuffer := serializeCSVFile(cnt)
		samplesBuffer, _ := json.Marshal(modeler.Samples())
		errors := make(map[string]string)
		for k, v := range modeler.ErrorMetrics() {
			errors[k] = fmt.Sprintf("%.5f", v)
		}
		modelDatasetModelInsert(coordinatesID, operatorID, datasetID, samplesBuffer, appxBuffer, conf, errors, sr)
		return nil
	}
	return task
}
