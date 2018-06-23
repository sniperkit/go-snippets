package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/oleiade/lane"
	"runtime"
	"sync"
	"time"
)

var buffNums = 7000
var workerNums = 5000
var msgNums = 6000

func init() {
	log.Errorln("Cpu numbers:", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type TaskQueue struct {
	signal chan *time.Time

	queue *lane.PQueue
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		queue:  lane.NewPQueue(lane.MINPQ),
		signal: make(chan *time.Time, buffNums),
	}
}

func (tq *TaskQueue) AddTask(index int) {
	tq.queue.Push(index, 1)
	t := time.Now()
	tq.signal <- &t
}

func (tq *TaskQueue) PopTask() interface{} {
	//	start := time.Now()
	value, _ := tq.queue.Pop()
	//	log.Errorln("Task pop timeused:", time.Now().Sub(start))
	return value
}

func (tq *TaskQueue) Wait() <-chan *time.Time {
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
		case t := <-w.tq.Wait():
			//w.tq.PopTask()
			timeused := time.Now().Sub(*t)
			go func() {
				time.Sleep(timeused)
				log.Errorln("Receive timeused:", timeused)
			}()
			//time.Sleep(time.Millisecond * 1)
		}
	}
}

func main() {
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
		go tq.AddTask(i)
	}

	time.Sleep(time.Second * 60)
}
