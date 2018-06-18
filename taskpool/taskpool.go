package taskpool

// Pool ...
type Pool struct {
	semaphor chan struct{}
	tasks    chan func()
}

// New ...
func New(size, spawn int) *Pool {
	p := &Pool{}
	for i := 0; i < spawn; i++ {
		p.semaphor <- struct{}{}
		go p.run(func() {})
	}
	return p
}

// Add ...
func (p *Pool) Add(task func()) {
	select {
	case p.tasks <- task:
		return
	case p.semaphor <- struct{}{}:
		go p.run(task)
	}
}

func (p *Pool) run(task func()) {
	task()
	for task := range p.tasks {
		task()
	}
	<-p.semaphor
}
