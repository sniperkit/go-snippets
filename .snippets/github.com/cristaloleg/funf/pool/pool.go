package pool

import "time"

// Pool interface
type Pool interface {
	Get() interface{}
	GetTimed(time.Duration) interface{}
	Put(interface{})
	Size() int
	Free() int
}

// NewStatic returns a static pool
func NewStatic(size int) Pool {
	p := &staticPool{}
	return p
}

type staticPool struct {
	size int
	free int
	pool chan interface{}
}

func (p *staticPool) Get() interface{} {
	return <-p.pool
}

func (p *staticPool) GetTimed(delay time.Duration) interface{} {
	select {
	case obj := <-p.pool:
		return obj
	case <-time.After(delay):
		return nil
	}
}

func (p *staticPool) Put(obj interface{}) {
	p.pool <- obj
}

func (p *staticPool) Size() int {
	return p.size
}
func (p *staticPool) Free() int {
	return p.free
}
