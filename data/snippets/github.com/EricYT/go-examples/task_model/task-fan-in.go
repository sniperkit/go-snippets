package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/oleiade/lane"
	"runtime"
	//	"runtime/debug"
	"sync"
	"time"
)

var msgNums = 10000
var taskChanNums = 1024
var workerNums = 1000

func init() {
	log.Errorln("Cpu numbers:", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
	//	debug.SetGCPercent(-1)
}

type TaskQueue struct {
	workers  chan *Worker
	taskChan chan *Task

	queue *lane.PQueue
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		queue:    lane.NewPQueue(lane.MINPQ),
		workers:  make(chan *Worker, workerNums),
		taskChan: make(chan *Task, taskChanNums),
	}
}

func (tq *TaskQueue) Serve(wg *sync.WaitGroup) {
	wg.Done()
	var count int
	var startTime *time.Time
	for {
		select {
		case task := <-tq.taskChan:
			count++
			if count == msgNums {
				log.Errorln("Total usedtime:", time.Now().Sub(*startTime))
			}
			if startTime == nil {
				tmp := time.Now()
				startTime = &tmp
			}
			go func(t *Task) {
				worker := tq.getWorker()
				worker.WakeUp(task)
			}(task)
		}
	}
}

func (tq *TaskQueue) getWorker() *Worker {
	return <-tq.workers
}

func (tq *TaskQueue) AddTask(t *Task) {
	tq.taskChan <- t
}

func (tq *TaskQueue) PopTask() interface{} {
	//	start := time.Now()
	value, _ := tq.queue.Pop()
	//	log.Errorln("Task pop timeused:", time.Now().Sub(start))
	return value
}

type Worker struct {
	id int
	tq *TaskQueue

	task chan *Task
}

func NewWorker(id int, tq *TaskQueue) *Worker {
	return &Worker{
		id:   id,
		tq:   tq,
		task: make(chan *Task),
	}
}

func (w *Worker) Serve(wg *sync.WaitGroup) {
	wg.Done()
	for {
		w.Register()
		select {
		case task := <-w.task:
			w.Run(task)
		}
	}
}

func (w *Worker) Run(t *Task) {
}

func (w *Worker) WakeUp(t *Task) {
	w.task <- t
}

func (w *Worker) Register() {
	w.tq.workers <- w
}

type Task struct {
	id   int
	time *time.Time
}

func (t Task) Void() {}

func NewTask(index int) *Task {
	//	now := time.Now()
	return &Task{
		id: index,
		//		time: &now,
	}
}

func main() {
	log.Errorln("Task queue start testing")

	tq := NewTaskQueue()
	var wg sync.WaitGroup

	wg.Add(1)
	go tq.Serve(&wg)
	wg.Wait()

	log.Errorln("start workers")
	wg.Add(workerNums)
	for i := 0; i < workerNums; i++ {
		go func(i int, tq *TaskQueue, wg *sync.WaitGroup) {
			worker := NewWorker(i, tq)
			worker.Serve(wg)
		}(i, tq, &wg)
	}
	wg.Wait()

	log.Errorln("send message")
	for i := 0; i < msgNums; i++ {
		go func() {
			task := NewTask(i)
			tq.AddTask(task)
		}()
	}

	time.Sleep(time.Second * 60)
}
