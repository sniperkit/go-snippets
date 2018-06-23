package transformation

import (
	"distributed-programming-abstractions/job_handlers"
)

type TransformationHandler struct {
	job_handler job.Handler // Could be sync or async
	buffer []job.Job
	N int 
	top int
	bottom int
	handling bool

	confirm func(job.Job)
	error func(job.Job)

	mutex_chan chan func()
}

func New(job_handler job.Handler, N int) *TransformationHandler {
	transformation_handler := &TransformationHandler {
		job_handler: job_handler, 
		N: N,
		confirm: func(job.Job) {},
		error: func(job.Job) {},
		mutex_chan: make(chan func()),
	}

	transformation_handler.Init()
	go transformation_handler.internal()

	return transformation_handler
}

func (transformation_handler *TransformationHandler) Init() {
	transformation_handler.top = 1
	transformation_handler.bottom = 1
	transformation_handler.handling = false
	transformation_handler.buffer = make([]job.Job, transformation_handler.N)
	transformation_handler.job_handler.Confirm(transformation_handler.JobHandlerConfirm)
}

func (transformation_handler *TransformationHandler) Submit(j job.Job) {
	transformation_handler.mutex_chan <- func() {
		if transformation_handler.bottom + transformation_handler.N == transformation_handler.top {
			go transformation_handler.error(j)
			return
		}
		transformation_handler.buffer[transformation_handler.top % transformation_handler.N] = j
		transformation_handler.top++
		go transformation_handler.confirm(j)
		go transformation_handler.handle_job()
	}
}

func (transformation_handler *TransformationHandler) handle_job() {
	transformation_handler.mutex_chan <- func() {
		if transformation_handler.bottom >= transformation_handler.top || transformation_handler.handling != false {
			return
		}
		j := transformation_handler.buffer[transformation_handler.bottom % transformation_handler.N]
		transformation_handler.bottom++
		transformation_handler.handling = true
		go transformation_handler.job_handler.Submit(j)
	}
}

func (transformation_handler *TransformationHandler) JobHandlerConfirm(j job.Job) {
	transformation_handler.mutex_chan <- func() {
		transformation_handler.handling = false
	}
}

func (transformation_handler *TransformationHandler) internal() {
	for f := range transformation_handler.mutex_chan {
		f()
	}
}

// Abstractions for setting the underlying functionalities.
func (transformation_handler *TransformationHandler) Confirm(f func(job.Job)) {
	transformation_handler.confirm = f
}

func (transformation_handler *TransformationHandler) Error(f func(job.Job)) {
	transformation_handler.error = f
}

func (transformation_handler *TransformationHandler) Process(f func(job.Job)) {
	transformation_handler.job_handler.Process(f)
}

