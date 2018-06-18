package scheduler

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	tomb "gopkg.in/tomb.v1"
)

var (
	ErrDynamicSchedulerOverMaxPendings error = errors.New("dynamic scheduler: pendding jobs over max limit")
	ErrDynamicSchedulerCurrence        error = errors.New("dynamic scheduler: currence number must greater than 0")
	ErrDynamicSchedulerThresold        error = errors.New("dynamic scheduler: thresold number less than 0")
	ErrDynamicSchedulerDispatch        error = errors.New("dynamic scheduler: dispatch to worker failed")
	ErrDynamicSchedulerShutdown        error = errors.New("dynamic scheduler: already shutdown")
)

type dynamicReactor struct {
	tomb tomb.Tomb

	id          int64
	thresold    int
	max         int
	maxPendings int

	waiters        []JobWrapper
	idleWorkers    []*worker
	runningWorkers map[int64]*worker
	jobCh          chan jobEntry
	controlCh      chan controlMsger

	workersWg sync.WaitGroup
}

func NewDynamicReactor(max, thresold, maxPendings int) *dynamicReactor {
	if max <= 0 {
		panic(ErrDynamicSchedulerCurrence)
	}
	if thresold < 0 {
		panic(ErrDynamicSchedulerThresold)
	}
	d := &dynamicReactor{
		thresold:    thresold,
		max:         max,
		maxPendings: maxPendings,

		waiters:        make([]JobWrapper, 0, 200),
		idleWorkers:    make([]*worker, 0, max),
		runningWorkers: make(map[int64]*worker),
		jobCh:          make(chan jobEntry),
		controlCh:      make(chan controlMsger),
	}

	// start main loop
	go func() {
		defer d.tomb.Done()
		d.tomb.Kill(d.runLoop())
	}()

	return d
}

func (d *dynamicReactor) Kill() error {
	d.tomb.Kill(nil)
	return d.tomb.Wait()
}

func (d *dynamicReactor) Idle() bool {
	idleCh := make(chan bool, 1)
	select {
	case d.controlCh <- idle{idleCh}:
		return <-idleCh
	case <-d.tomb.Dead():
		return true
	}
}

func (d *dynamicReactor) Schedule(j JobWrapper) error {
	ackC := make(chan error, 1)
	select {
	case d.jobCh <- jobEntry{j, ackC}:
	case <-d.tomb.Dead():
		return ErrDynamicSchedulerShutdown
	}

	log.Println("[scheduler] put in")
	// scheduler get the request must return a value
	return <-ackC
}

func (d *dynamicReactor) pushWaiter(j JobWrapper) {
	d.waiters = append(d.waiters, j)
}

func (d *dynamicReactor) popWaiter() JobWrapper {
	if len(d.waiters) == 0 {
		return nil
	}
	j := d.waiters[0]
	d.waiters = d.waiters[1:]
	return j
}

func (d *dynamicReactor) pushIdleWorker(w *worker) {
	d.idleWorkers = append(d.idleWorkers, w)
}

func (d *dynamicReactor) popIdleWorker() *worker {
	if len(d.idleWorkers) == 0 {
		return nil
	}
	w := d.idleWorkers[0]
	d.idleWorkers = d.idleWorkers[1:]
	return w
}

func (d *dynamicReactor) insertRunningWorker(w *worker) {
	d.runningWorkers[w.id] = w
}

func (d *dynamicReactor) removeRunningWorker(id int64) {
	// FIXME: Is need to check there is one exist ?
	delete(d.runningWorkers, id)
}

func (d *dynamicReactor) spawnWorker(ctx context.Context, workersCh chan<- *worker) *worker {
	id := d.id
	d.id++
	w := newWorker(id)
	d.workersWg.Add(1)
	go func() {
		defer d.workersWg.Done()
		w.run(ctx, workersCh)
	}()
	return w
}

func (d *dynamicReactor) totalWorkers() int {
	return len(d.idleWorkers) + len(d.runningWorkers)
}

func (d *dynamicReactor) isIdle(idle chan bool) {
	idle <- (len(d.runningWorkers) == 0 && len(d.waiters) == 0)
}

func (d *dynamicReactor) clean(cancel func()) {
	// clean all waiters
	for _, waiter := range d.waiters {
		waiter.Interrupt(ErrDynamicSchedulerShutdown)
	}
	d.waiters = []JobWrapper{}

	// ctx cancel
	cancel()

	// wait all workers done
	d.workersWg.Wait()

	// clean all references to these workers
	d.idleWorkers = []*worker{}
	d.runningWorkers = make(map[int64]*worker)
}

// for debug
func (d *dynamicReactor) report() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			log.Printf("[report] idle workers: %d running workers: %d waiters: %d idle: %v\n", len(d.idleWorkers), len(d.runningWorkers), len(d.waiters), d.Idle())
		case <-d.tomb.Dying():
			return
		}
	}
}

func (d *dynamicReactor) runLoop() error {
	var workersCh = make(chan *worker, d.thresold)
	ctx, cancel := context.WithCancel(context.Background())

	go d.report()

	for {
		select {
		case je := <-d.jobCh: //job channel
			log.Printf("[scheduler] got a job")
			var err error
			worker := d.popIdleWorker()
			if worker != nil {
				log.Printf("[scheduler] got a idle woker id: %d", worker.id)
				d.insertRunningWorker(worker)
				worker.jobCh <- je.job
			} else if d.totalWorkers() < d.max {
				w := d.spawnWorker(ctx, workersCh)
				log.Printf("[scheduler] spawn a new worker id: %d", w.id)
				d.insertRunningWorker(w)
				w.jobCh <- je.job
			} else {
				if len(d.waiters) > d.maxPendings {
					err = ErrDynamicSchedulerOverMaxPendings
				} else {
					d.pushWaiter(je.job)
				}
			}
			je.ackC <- err
		case w := <-workersCh: // worker channel
			log.Printf("[scheduler] worker: %d done. waiters: %d", w.id, len(d.waiters))
			if len(d.waiters) > 0 {
				j := d.popWaiter()
				w.jobCh <- j
			} else {
				d.removeRunningWorker(w.id)
				if len(d.idleWorkers) < d.thresold {
					d.pushIdleWorker(w)
				} else {
					// stop this worker
					close(w.jobCh)
				}
			}
		case m := <-d.controlCh: // control channel
			switch msg := m.(type) {
			case idle:
				d.isIdle(msg.idleCh)
			default:
				panic("unknow control message")
			}
		case <-d.tomb.Dying():
			d.clean(cancel)
			log.Printf("[scheduler] I'm done")
			return nil
		}
	}
}

type jobEntry struct {
	job  JobWrapper
	ackC chan error
}

type controlMsger interface {
	Control()
}

type idle struct {
	idleCh chan bool
}

func (i idle) Control() {}

// worker
type worker struct {
	id    int64
	jobCh chan JobWrapper
}

func newWorker(id int64) *worker {
	log.Printf("[worker] new woker id: %d\n", id)
	w := &worker{
		id:    id,
		jobCh: make(chan JobWrapper),
	}
	return w
}

func (w *worker) run(ctx context.Context, workersCh chan<- *worker) {
	for {
		// waiting for a job or dead message
		select {
		case <-ctx.Done():
			log.Printf("[worker] id: %d Parent done, I'm done", w.id)
			return
		case job, ok := <-w.jobCh:
			if !ok {
				log.Printf("[worker] Parent let me done id: %d I'm done", w.id)
				return
			}
			log.Printf("[worker] id: %d got a job", w.id)
			job.Run(ctx)
		}
		// job done
		select {
		case workersCh <- w:
		case <-ctx.Done():
			log.Printf("[worker] id: %d Parent done, I'm done", w.id)
			return
		}
	}
}
