package workerpool

import (
	"sync"
	"sync/atomic"
)

type WorkerPool struct {
	atomicSize int32
	wg         sync.WaitGroup
	fn         func(interface{}) interface{}
	data       <-chan interface{}
	result     chan interface{}
}

func New(fn func(interface{}) interface{}, data <-chan interface{}) *WorkerPool {
	return &WorkerPool{
		fn:     fn,
		data:   data,
		result: make(chan interface{}),
	}
}

func NewBuffered(buffer int, fn func(interface{}) interface{}, data <-chan interface{}) *WorkerPool {
	return &WorkerPool{
		fn:     fn,
		data:   data,
		result: make(chan interface{}, buffer),
	}
}

func (w *WorkerPool) Start() {
	w.Inc(1)
}

func (w *WorkerPool) Wait() {
	w.wg.Wait()
	close(w.result)
}

func (w *WorkerPool) Results() <-chan interface{} {
	return w.result
}

func (w *WorkerPool) Inc(delta int) {
	if delta <= 0 {
		return
	}
	for i := 1; i < delta; i++ {
		w.invokeNew()
	}
}

func (w *WorkerPool) Dec(delta int) {
	if delta <= 0 {
		return
	}
	w.wg.Add(-delta)
	atomic.AddInt32(&w.atomicSize, -int32(delta))
}

func (w *WorkerPool) invokeNew() {
	id := atomic.LoadInt32(&w.atomicSize)
	atomic.AddInt32(&w.atomicSize, 1)
	w.wg.Add(1)

	go func(id int32) {
		for d := range w.data {
			w.result <- w.fn(d)

			if atomic.LoadInt32(&w.atomicSize) < id {
				break
			}
		}
		w.wg.Done()
	}(id + 1)
}
