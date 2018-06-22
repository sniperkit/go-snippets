package promise

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPromise(t *testing.T) {
	tests := []struct {
		description string
		v           interface{}
	}{
		{
			description: "Given an integer value, the promise should return the same value",
			v:           4,
		},
		{
			description: "Given a nil value, the promise should return nil",
			v:           nil,
		},
	}
	for _, test := range tests {
		test := test // Capture range variable.
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()
			p := NewPromise()
			go p.Complete(test.v)
			v, err := p.Get()
			assert.NoError(t, err, "Unexpected error")
			assert.Equal(t, test.v, v, "Got unexpected value from promise")
		})
	}
}

func TestPromise_timeout(t *testing.T) {
	p := NewPromise()
	result := make(chan interface{})
	defer close(result)
	go func() {
		v, err := p.Get()
		assert.NoError(t, err, "Unexpected error")
		result <- v
	}()
	select {
	case <-time.After(1 * time.Second):
	case got := <-result:
		t.Errorf("Expecting timeout, but got value %+v", got)
	}
}

func TestPromise_getFirst(t *testing.T) {
	p := NewPromise()
	go func() {
		time.Sleep(1 * time.Second)
		p.Complete("ok")
	}()
	v, err := p.Get()
	assert.NoError(t, err, "Unexpected error")
	assert.Equal(t, "ok", v, "Got unexpected value from promise")
}

func TestPromise_setFirst(t *testing.T) {
	p := NewPromise()
	result := make(chan interface{})
	defer close(result)
	go func() {
		time.Sleep(1 * time.Second)
		v, err := p.Get()
		assert.NoError(t, err, "Unexpected error")
		result <- v
	}()
	p.Complete("ok")
	assert.Equal(t, "ok", <-result, "Got unexpected value from promise")
}

func TestPromise_multipleComplete(t *testing.T) {
	p := NewPromise()
	go func() {
		p.Complete("ok")
		p.Complete("not ok")
	}()
	v, err := p.Get()
	assert.NoError(t, err, "Unexpected error")
	assert.Equal(t, "ok", v, "Got unexpected value from promise")
}

func TestPromise_error(t *testing.T) {
	p := NewPromise()
	go func() {
		time.Sleep(1 * time.Second)
		p.CompleteWithError(errors.New("Expected error"))
	}()
	v, err := p.Get()
	assert.Error(t, err, "Expected an error")
	assert.Equal(t, nil, v, "Got unexpected value from promise")
}

func TestAll(t *testing.T) {
	tests := []struct {
		description string
		putters     []func(*Promise)
		vs          []interface{}
	}{
		{
			description: "",
			putters: []func(*Promise){
				func(p *Promise) { p.Complete(1) },
				func(p *Promise) { p.Complete(2) },
				func(p *Promise) { p.Complete(3) },
				func(p *Promise) { p.Complete(4) },
			},
			vs: []interface{}{1, 2, 3, 4},
		},
	}
	for _, test := range tests {
		test := test // Capture range variable.
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()
			promises := make([]*Promise, len(test.putters))
			for i := 0; i < len(promises); i++ {
				promises[i] = NewPromise()
			}
			for i, putter := range test.putters {
				putter(promises[i])
			}
			vs, errs := All(promises...)
			for _, err := range errs {
				assert.NoError(t, err, "Unexpected error")
			}
			assert.Equal(t, test.vs, vs, "")
		})
	}
}

func BenchmarkPromise(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p := NewPromise()
			p.Complete("ok")
			p.Get()
		}
	})
}

type chanPromise struct {
	v          interface{}
	isComplete chan interface{}
	once       sync.Once
}

func newChanPromise() *chanPromise {
	return &chanPromise{isComplete: make(chan interface{})}
}

func (cp *chanPromise) Get() interface{} {
	<-cp.isComplete
	return cp.v
}

func (cp *chanPromise) Complete(v interface{}) {
	cp.once.Do(func() {
		cp.v = v
		close(cp.isComplete)
	})
}

func BenchmarkChanPromise(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p := newChanPromise()
			p.Complete("ok")
			p.Get()
		}
	})
}
