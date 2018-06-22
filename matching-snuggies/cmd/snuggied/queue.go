package main

import (
	"fmt"
	"log"
	"sync"
)

// Scheduler is the write-end of a job queue.  It takes a mesh file url, the
// name of a slicer and a preset for that slicer.  Scheduler is responsible for
// routing the job to a machine capable of servicing the request.
type Scheduler interface {
	ScheduleSliceJob(id, meshurl, slicer, preset string) error
	CancelSliceJob(id string)
}

// Consumer is the read-end of a job queue.  It reserves a from the queue and
// ensures any remote mesh file locations are downloaded to local paths.
type Consumer interface {
	NextSliceJob() (*Job, error)
}

type Job struct {
	// NodeID is the job's originating node.
	NodeID  string
	ID      string
	MeshURL string
	Slicer  string
	Preset  string

	// Cancel receives a value if the job has been cancelled by the scheduling
	// process.
	Cancel <-chan error

	// Done is called when the slicing process has terminated.  Done is passed
	// a path at which the output G-code can be retreived.  If the G-code could
	// not be generated due to failure a non-nil error must be passed to Done.
	Done func(path string, err error)
}

// MemQueue is an in memory database and job queue that implements the
// Scheduler and Consumer interfaces.  MemQueue is safe for many producers and
// consumers to be calling interface methods simultaneously.
type MemQueue struct {
	NodeID  string
	Started func(id string)
	Done    func(id, path string, err error)
	cond    sync.Cond
	jobs    []*memJob
	db      map[string]*memJob
}

var _ Scheduler = new(MemQueue)
var _ Consumer = new(MemQueue)

// MemoryQueue allocates and initializes a new MemQueue.  The function argument
// is called when consumers finish work on a job.
func MemoryQueue(done func(id, path string, err error)) *MemQueue {
	return &MemQueue{
		Done: done,
		cond: sync.Cond{L: new(sync.Mutex)},
		db:   make(map[string]*memJob),
	}
}

func (q *MemQueue) jobTerminated(id string) {
	q.cond.L.Lock()
	delete(q.db, id)
	qlen := len(q.jobs)
	dblen := len(q.db)
	q.cond.L.Unlock()
	log.Printf("jobs running:%d queued:%d", dblen-qlen, qlen)
}

// ScheduleSliceJob enqueues a job in q.
func (q *MemQueue) ScheduleSliceJob(id, meshurl, slicer, preset string) error {
	j := &memJob{
		ID:       id,
		NodeID:   q.NodeID,
		Location: meshurl,
		Slicer:   slicer,
		Preset:   preset,
		Cancel:   make(chan error, 1),
		Done:     make(chan struct{}),
		Fin: func(id, path string, err error) {
			q.jobTerminated(id)
			if q.Done != nil {
				q.Done(id, path, err)
			}
		},
	}

	// append the job to the queue and signal a waiting consumer goroutine to
	// wake up and process the job.
	q.cond.L.Lock()
	q.jobs = append(q.jobs, j)
	q.db[j.ID] = j
	qlen := len(q.jobs)
	dblen := len(q.db)
	q.cond.Signal()
	q.cond.L.Unlock()
	log.Printf("jobs running:%d queued:%d", dblen-qlen, qlen)

	return nil
}

// BUG:
// CancelSliceJob can temporarily skew the count of queued jobs because pending
// jobs are not removed from the queue immediately on cancellation.
func (q *MemQueue) CancelSliceJob(id string) {
	q.cond.L.Lock()
	if j := q.db[id]; j != nil {
		j.Cancel <- fmt.Errorf("the job was cancelled")
		delete(q.db, id)
	}
	q.cond.L.Unlock()
}

// NextSliceJob dequeues a job from q or blocks until one is available.
func (q *MemQueue) NextSliceJob() (*Job, error) {
	q.cond.L.Lock()
	for len(q.jobs) == 0 {
		q.cond.Wait()
	}
	j := q.jobs[0]
	q.jobs = q.jobs[1:]
	qlen := len(q.jobs)
	dblen := len(q.db)
	q.cond.L.Unlock()

	select {
	case <-j.Cancel:
		// the job was cancelled previously get another job
		return q.NextSliceJob()
	default:
	}

	go func() {
		if q.Started != nil {
			q.Started(j.ID)
		}
	}()
	log.Printf("jobs running:%d queued:%d", dblen-qlen, qlen)

	return j.Job(), nil
}

type memJob struct {
	ID       string
	NodeID   string
	Location string
	Slicer   string
	Preset   string
	Cancel   chan error
	Done     chan struct{}
	Fin      func(string, string, error)
}

func (m *memJob) Job() *Job {
	return &Job{
		ID:      m.ID,
		NodeID:  m.NodeID,
		MeshURL: m.Location,
		Slicer:  m.Slicer,
		Preset:  m.Preset,
		Cancel:  m.Cancel,
		Done: func(path string, err error) {
			close(m.Done)
			m.Fin(m.ID, path, err)
		},
	}
}
