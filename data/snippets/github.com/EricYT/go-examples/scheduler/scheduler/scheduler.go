package scheduler

import (
	"context"
	"errors"
	"log"
	"sync"
	"sync/atomic"

	tomb "gopkg.in/tomb.v1"
)

var (
	ErrSchedulerOverMaxPendings error = errors.New("scheduler: pendding jobs over max limit")
	ErrSchedulerReset           error = errors.New("scheduler: reset")
	ErrSchedulerCurrence        error = errors.New("scheduler: currence number must greater than 0")
	ErrSchedulerDispatch        error = errors.New("scheduler: dispatch to worker failed")
	ErrSchedulerShutdown        error = errors.New("scheduler: already shutdown")
)

type Scheduler interface {
	// .Schedule schedule a job wrapper into reactor
	Schedule(job JobWrapper) error
	// .Kill kill the reactor and wait it done and return a error message
	Kill(reason error) error
	// .Idle returns true if there are no running jobs and waiters.
	Idle() bool
}

type reactor struct {
	tomb  *tomb.Tomb
	mutex sync.Mutex

	currence    int
	maxPendings int

	runningCount int64
	pendingCount int64

	cancel   func()
	jobsC    chan chan<- JobWrapper
	pendings []JobWrapper
	resume   chan struct{}
}

func NewReactor(currence int, maxPendings int) *reactor {
	if currence <= 0 {
		panic(ErrSchedulerCurrence)
	}
	r := &reactor{
		tomb:        new(tomb.Tomb),
		currence:    currence,
		maxPendings: maxPendings,
		jobsC:       make(chan chan<- JobWrapper, currence),
		pendings:    make([]JobWrapper, 0, 200),
		resume:      make(chan struct{}, 1),
	}

	// cancel context
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	// start workers
	for i := 0; i < currence; i++ {
		go r.workerLoop(ctx, r.jobsC)
	}

	// start main loop
	go func() {
		defer r.tomb.Done()
		r.tomb.Kill(r.runLoop())
	}()
	return r
}

func (r *reactor) Pendings() int64 {
	return atomic.LoadInt64(&r.pendingCount)
}

func (r *reactor) Running() int64 {
	return atomic.LoadInt64(&r.runningCount)
}

func (r *reactor) Schedule(j JobWrapper) error {
	log.Printf("reactor: add job: %#v", j)
	select {
	case <-r.tomb.Dead():
		return ErrSchedulerShutdown
	default:
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()
	if len(r.pendings) >= r.maxPendings {
		return ErrSchedulerOverMaxPendings
	}

	atomic.AddInt64(&r.pendingCount, 1)
	r.pendings = append(r.pendings, j)
	// maybe reset resume
	if len(r.pendings) == 1 {
		r.resume <- struct{}{}
	}
	return nil
}

func (r *reactor) Kill() error {
	r.tomb.Kill(nil)
	r.cancel()
	return r.tomb.Wait()
}

func (r *reactor) runLoop() error {
	log.Println("reactor(loop): main loop run")

	// select
	closed := make(chan struct{})
	close(closed)

	var next chan struct{}

	for {
		select {
		case <-r.resume:
			log.Println("reactor(loop): receive a resume signal")
			next = closed
		case <-next:
			// pop one job
			job, empty := r.popOne()
			if empty {
				next = nil
			} else {
				next = closed
			}

			// job run
			if job != nil {
				// this will block main loop if there is not a idle worker
				workerC := r.pickOneWorker()
				select {
				case workerC <- job:
					atomic.AddInt64(&r.pendingCount, -1)
				default:
					// worker already dead
					log.Printf("reactor(loop): dispatch job: %#v error", job)
					job.Interrupt(ErrSchedulerDispatch)
				}
			}
		case <-r.tomb.Dying():
			log.Println("reactor(loop): ready to die")
			r.reset()
			return nil
		}
	}
	return nil
}

func (r *reactor) reset() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, job := range r.pendings {
		log.Printf("reactor(rest): kill job: %#v", job)
		job.Interrupt(ErrSchedulerReset)
	}
	r.pendings = []JobWrapper{}
	atomic.StoreInt64(&r.pendingCount, 0)
}

func (r *reactor) popOne() (JobWrapper, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if len(r.pendings) == 0 {
		return nil, true
	}
	job := r.pendings[0]
	r.pendings = r.pendings[1:]
	return job, len(r.pendings) == 0
}

func (r *reactor) pickOneWorker() chan<- JobWrapper {
	select {
	case <-r.tomb.Dying():
		return nil
	case jobC := <-r.jobsC:
		return jobC
	}
}

func (r *reactor) workerLoop(ctx context.Context, jobsC chan chan<- JobWrapper) {
	log.Println("reactor(worker): loop run")
	jobC := make(chan JobWrapper)
	for {
		select {
		case jobsC <- jobC:
		case <-r.tomb.Dying():
			return
		}

		select {
		case <-r.tomb.Dying():
			log.Println("reactor(worker): worker done")
			return
		case job := <-jobC:
			atomic.AddInt64(&r.runningCount, 1)
			job.Run(ctx)
			job = nil
			atomic.AddInt64(&r.runningCount, -1)
		}
	}
}
