package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/oleiade/lane"
	"runtime"
	//	"runtime/debug"
	"sync"
	"time"
)

var msgNums = 600
var taskChanNums = msgNums
var workerNums = 1

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
	for tmp := range tq.taskChan {
		go func(t *time.Time) {
			log.Errorln("(3) get task used:", time.Now().Sub(*t))
			//now := time.Now()
			//worker := tq.getWorker()
			//log.Errorln("(2) get worker used:", time.Now().Sub(now))
			//worker.WakeUp(task)
		}(tmp.time)
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
	w.Register()
	for {
		select {
		case task := <-w.task:
			w.Run(task)
			w.Register()
		}
	}
}

func (w *Worker) Run(t *Task) {
	log.Errorln("(1) worker get job timeused:", time.Now().Sub(*t.time))
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
	now := time.Now()
	return &Task{
		id:   index,
		time: &now,
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
		go func() {
			worker := NewWorker(i, tq)
			worker.Serve(&wg)
		}()
	}
	wg.Wait()

	log.Errorln("send message")
	for i := 0; i < msgNums; i++ {
		func(index int, tq_ *TaskQueue) {
			task := NewTask(index)
			tq_.AddTask(task)
		}(i, tq)
	}

	time.Sleep(time.Second * 60)
}
