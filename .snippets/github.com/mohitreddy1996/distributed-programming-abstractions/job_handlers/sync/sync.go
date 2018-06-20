package sync

import (
	"distributed-programming-abstractions/job_handlers"
)

type JobHandler struct {
	process func(j job.Job)
	confirm func(j job.Job)

	mutex_chan chan func()
}

func New() *JobHandler {
	jh := &JobHandler {
		process: func(j job.Job) {}, // Do nothing when process is called.
		confirm: func(j job.Job) {}, // Do nothing when confirm is called.

		mutex_chan: make(chan func()),
	}

	go jh.internal()
	return jh
}

func (jh *JobHandler) Process(f func(job.Job)) {
	jh.process = f
}

func (jh *JobHandler) Confirm(f func(job.Job)) {
	jh.confirm = f
}

func (jh *JobHandler) internal() {
	for f := range jh.mutex_chan {
		f()
	}
}

func (jh *JobHandler) Submit(j job.Job) {
	jh.mutex_chan <- func () {
		jh.process(j)

		go jh.confirm(j)
	}
}