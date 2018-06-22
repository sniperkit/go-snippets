package roundrobin

import (
	"sync"
)

// RoundRobin ...
type RoundRobin struct {
	sync.Mutex

	current int
	pool    []int
}

// NewRoundRobin ...
func NewRoundRobin(resources []int) *RoundRobin {
	return &RoundRobin{
		current: 0,
		pool:    resources,
	}
}

// Get ...
func (r *RoundRobin) Get() int {
	r.Lock()
	defer r.Unlock()

	if r.current >= len(r.pool) {
		r.current = r.current % len(r.pool)
	}

	result := r.pool[r.current]
	r.current++
	return result
}

// RobinHood ...
type RobinHood struct {
	current int
	pool    []int

	requestQ chan chan int
}

// NewRobinHood ...
func NewRobinHood(resources []int) *RobinHood {
	r := &RobinHood{
		current:  0,
		pool:     resources,
		requestQ: make(chan chan int),
	}
	go r.balancer()
	return r
}

// Get ...
func (r *RobinHood) Get() int {
	output := make(chan int, 1)
	r.requestQ <- output
	return <-output
}

// balancer ...
func (r *RobinHood) balancer() {
	for {
		select {
		case output := <-r.requestQ:
			if r.current >= len(r.pool) {
				r.current = 0
			}
			output <- r.pool[r.current]
			r.current++
		}
	}
}
