package scheduler

import "context"

type JobWrapper interface {
	// .Run execute this job correctly
	Run(ctx context.Context)
	// .Interrupt break this execute abnormally
	Interrupt(err error)
}

type jobWrapper struct {
	run       func(ctx context.Context)
	interrupt func(err error)
}

func NewJobWrapper(run func(ctx context.Context), interrupt func(error)) *jobWrapper {
	j := &jobWrapper{
		run:       run,
		interrupt: interrupt,
	}
	return j
}

func (j *jobWrapper) Run(ctx context.Context) {
	j.run(ctx)
}

func (j *jobWrapper) Interrupt(err error) {
	j.interrupt(err)
}
