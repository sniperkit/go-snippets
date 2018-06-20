package job

import (
	"sync"
)

// manages a pool of goroutines.
// it uses a buffered pool of workers in a Job/Worker pattern
type WorkerPool struct {
	// A pool of worker channels that are registered with the pool
	WorkerPool chan chan Job
	Workers    []Worker
	maxWorkers int
	waitGroup  sync.WaitGroup
}

func NewWorkerPool(maxWorkers int) WorkerPool {
	pool := make(chan chan Job, maxWorkers)
	workers := make([]Worker, 0)
	return WorkerPool{
		WorkerPool: pool,
		Workers:    workers,
		maxWorkers: maxWorkers,
		waitGroup:  sync.WaitGroup{}}
}

// Starts the WorkerPool
func (p *WorkerPool) Run(queue chan Job) {
	w := p.waitGroup

	// starting n number of workers
	for i := 0; i < p.maxWorkers; i++ {
		worker := NewWorker(p.WorkerPool)
		p.Workers = append(p.Workers, worker)
		w.Add(1)
		worker.Start(&w)
	}

	go p.dispatch(queue)
}

// stops the Pool
func (p *WorkerPool) Stop() bool {
	// stops all workers
	for _, worker := range p.Workers {
		worker.Stop()
	}
	p.waitGroup.Wait() //Wait for the goroutines to shutdown

	close(p.WorkerPool)
	for x := range p.WorkerPool {
		_ = x // ignore channel var using blank identifier
	}
	ok := p.IsOpen()
	return ok
}

// dispatches a job to be handled by an idle Worker of the pool
func (p *WorkerPool) dispatch(jobQueue chan Job) {
	for {
		select {
		case job := <-jobQueue:
			// a model request has been received
			go func(job Job) {
				// try to obtain a worker model channel that is available.
				// this will block until a worker is idle
				jobChannel := <-p.WorkerPool

				// dispatch the model to the worker model channel
				jobChannel <- job
			}(job)
		}
	}
}

// checks if a Worker Pool is open or closed - If we can recieve on the channel then it is NOT closed
func (p *WorkerPool) IsOpen() bool {
	_, ok := <-p.WorkerPool
	return ok
}

// Utility function to merge multiple JobResult output channels into one
func Merge(cs ...<-chan JobResult) <-chan JobResult {
	var wg sync.WaitGroup
	out := make(chan JobResult)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan JobResult) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
