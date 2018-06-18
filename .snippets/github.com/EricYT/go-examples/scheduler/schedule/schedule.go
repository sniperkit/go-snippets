package schedule

import (
	"context"
	"sync"
)

type Job func(ctx context.Context)

type Scheduler interface {
	// Schedule asks the scheduler to schedule a job defined by the given func.
	// Schedule to a stopped scheduler might panic.
	Schedule(j Job)

	// Pending returns the number of pending jobs
	Pending() int

	// Scheduled returns the number of scheduled jobs (excluding pending jobs)
	Scheduled() int

	// Finished returns the number of finished jobs
	Finished() int

	// WaitFinish waits until at least n job are finished and all pending jobs are finished
	WaitFinish(n int)

	// Stop stop the scheduler
	Stop()
}

type fifo struct {
	mu sync.Mutex

	currence  int
	scheduled int
	finished  int
	pendings  []Job
	resume    chan struct{}

	ctx    context.Context
	cancel context.CancelFunc

	finishCond *sync.Cond
	donec      chan struct{}
}

func NewFIFOScheduler(currence int) Scheduler {
	f := &fifo{
		currence: currence,
		resume:   make(chan struct{}, currence),
		donec:    make(chan struct{}, 1),
	}
	f.finishCond = sync.NewCond(&f.mu)
	f.ctx, f.cancel = context.WithCancel(context.Background())
	go f.run()
	return f
}

func (f *fifo) Schedule(j Job) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.cancel == nil {
		panic("schedule: schedule to stopped scheduler")
	}

	f.pendings = append(f.pendings, j)
	select {
	case f.resume <- struct{}{}:
	default:
	}
}

func (f *fifo) Pending() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.pendings)
}

func (f *fifo) Scheduled() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.scheduled
}

func (f *fifo) Finished() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.finished
}

func (f *fifo) WaitFinish(n int) {
	f.finishCond.L.Lock()
	for f.finished < n || len(f.pendings) != 0 {
		f.finishCond.Wait()
	}
	f.finishCond.L.Unlock()
}

func (f *fifo) Stop() {
	f.mu.Lock()
	f.cancel()
	f.cancel = nil
	f.mu.Unlock()
	<-f.donec
}

func (f *fifo) run() {
	defer func() {
		close(f.donec)
		close(f.resume)
	}()

	for {
		select {
		case <-f.resume:
			f.mu.Lock()
			if len(f.pendings) != 0 {
				todo := f.pendings[0]
				f.pendings = f.pendings[1:]
				f.scheduled++
				go func(j Job) {
					defer func() {
						f.finishCond.L.Lock()
						f.finished++
						f.finishCond.Broadcast()
						f.finishCond.L.Unlock()
						select {
						case f.resume <- struct{}{}:
						default:
						}
					}()
					j(f.ctx)
				}(todo)
			}
			f.mu.Unlock()
		case <-f.ctx.Done():
			f.mu.Lock()
			pendings := f.pendings
			f.pendings = nil
			f.mu.Unlock()
			for _, todo := range pendings {
				go func(j Job) {
					j(f.ctx)
				}(todo)
			}
		}
	}
}
