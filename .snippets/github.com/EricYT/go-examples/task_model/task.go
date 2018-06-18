package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/oleiade/lane"
	"github.com/pkg/profile"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var buffNums = 1024
var workerNums = 20
var msgNums = 10000

// count
var count int32
var start *time.Time

func init() {
	log.Errorln("Cpu numbers:", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetLevel(log.ErrorLevel)
}

type TaskQueue struct {
	signal chan int

	queue *lane.PQueue
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		queue:  lane.NewPQueue(lane.MINPQ),
		signal: make(chan int, buffNums),
	}
}

func (tq *TaskQueue) AddTask(index int) {
	tq.queue.Push(index, 1)
	tq.signal <- index
}

func (tq *TaskQueue) PopTask() interface{} {
	value, _ := tq.queue.Pop()
	return value
}

func (tq *TaskQueue) Wait() <-chan int {
	return tq.signal
}

type Worker struct {
	id int
	tq *TaskQueue
}

func NewWorker(id int, tq *TaskQueue) *Worker {
	return &Worker{
		id: id,
		tq: tq,
	}
}

func (w *Worker) Run(wg *sync.WaitGroup) {
	wg.Done()
	for {
		select {
		case <-w.tq.Wait():
			value := atomic.AddInt32(&count, 1)
			log.Debugln("number:", value)
			if value == int32(msgNums) {
				log.Errorln("total timeused:", time.Now().Sub(*start))
				return
			} else if value == 1 {
				tmp := time.Now()
				start = &tmp
			}
			w.tq.PopTask()
		}
	}
}

func main() {
	//defer profile.Start(profile.BlockProfile).Stop()
	defer profile.Start(profile.CPUProfile).Stop()
	log.Errorln("Task queue start testing")

	tq := NewTaskQueue()
	var wg sync.WaitGroup

	wg.Add(workerNums)
	for i := 0; i < workerNums; i++ {
		worker := NewWorker(i, tq)
		go worker.Run(&wg)
	}
	wg.Wait()

	for i := 0; i < msgNums; i++ {
		go func(tq *TaskQueue) {
			tq.AddTask(i)
		}(tq)
	}

	time.Sleep(time.Second * 2)
}
