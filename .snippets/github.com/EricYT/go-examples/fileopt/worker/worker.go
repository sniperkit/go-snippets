package worker

import (
	"errors"
	"log"

	tomb "gopkg.in/tomb.v1"
)

var (
	ErrReceiveWrongJob error = errors.New("Worker: receive a wrong job")
)

type Worker struct {
	tomb *tomb.Tomb
	g    Generator
}

func NewWorker(g Generator) *Worker {
	w := &Worker{
		tomb: new(tomb.Tomb),
		g:    g,
	}

	// start this worker
	go func() {
		defer w.tomb.Done()
		w.tomb.Kill(w.runLoop())
	}()

	return w
}

func (w *Worker) runLoop() error {
	for {
		select {
		case job, ok := <-w.g.Generate():
			if !ok {
				return ErrReceiveWrongJob
			}
			// operate this job
			log.Printf("worker(%p) run\n", w)
			job.Done(job.Execute())
		case <-w.tomb.Dying():
			return nil
		}
	}
}

type Jobber interface {
	// run this job and returns a error message
	Execute() error
	// let others know this job already done and returns a result
	Done(error)
	// wait this job done and get a result
	Wait() error
}
