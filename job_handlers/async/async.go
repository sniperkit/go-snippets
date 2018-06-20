package async

import (
	"distributed-programming-abstractions/job_handlers"
)

type JobHandler struct {
	confirm func(job.Job)
	process func(job.Job)
	buffer []job.Job

	mutex_chan chan func()
}

func New() *JobHandler {
	jh := &JobHandler {
		process: func(job.Job) {},
		confirm: func(job.Job) {},
		mutex_chan: make(chan func()),
	}

	jh.Init()
	go jh.internal()
	return jh
}

func (jh *JobHandler) Init() {
	jh.buffer = make([]job.Job, 0)
}

func (jh *JobHandler) internal() {
	for f:= range jh.mutex_chan {
		f()
	}
}

func (jh *JobHandler) Process(f func(job.Job)) {
	jh.process = f
}

func (jh *JobHandler) Confirm(f func(job.Job)) {
	jh.confirm = f
}

func (jh *JobHandler) Submit(j job.Job) {
	jh.mutex_chan <- func() {
		jh.buffer = append(jh.buffer, j)
		go jh.confirm(j)
		go jh.Handlebuffer()
	}
}

func (jh *JobHandler) Handlebuffer() {
	jh.mutex_chan <- func() {
		j:= jh.buffer[0]
		jh.process(j)
		jh.buffer = jh.buffer[1:]
	}
}