package promise

import "sync"

// Promise is a handle to a value or error that will be
// available sometime in the future.
type Promise struct {
	v          interface{}
	err        error
	cond       *sync.Cond
	isComplete bool
}

// NewPromise returns an incomplete Promise.
func NewPromise() *Promise {
	return &Promise{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

// Get returns the promised value or error when it is ready.
// It blocks until the value or error are ready.
func (p *Promise) Get() (interface{}, error) {
	p.cond.L.Lock()
	for !p.isComplete {
		p.cond.Wait()
	}
	p.cond.L.Unlock()
	return p.v, p.err
}

// Complete completes the promise, storing the provided value
// and unblocking and returning the provided value to all
// calls to Get.
func (p *Promise) Complete(v interface{}) {
	p.complete(v, nil)
}

// CompleteWithError completes the promise with an error,
// storing the provided error and unblocking and returning
// the provided error to all calls to Get.
func (p *Promise) CompleteWithError(err error) {
	p.complete(nil, err)
}

func (p *Promise) complete(v interface{}, err error) {
	p.cond.L.Lock()
	if p.isComplete {
		p.cond.L.Unlock()
		return
	}
	p.v = v
	p.err = err
	p.isComplete = true
	p.cond.L.Unlock()
	p.cond.Broadcast()
}

// All returns all promised values or errors when all provided
// promises are completed. The order of the returned values
// and errors match the order of the provided promises.
func All(promises ...*Promise) ([]interface{}, []error) {
	if len(promises) == 0 {
		return nil, nil
	}
	vs := make([]interface{}, len(promises))
	errs := make([]error, len(promises))
	for i, p := range promises {
		vs[i], errs[i] = p.Get()
	}
	return vs, errs
}
