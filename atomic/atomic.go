package atomic

import "sync/atomic"

// Atomic is an atomic value update by predicate
type Atomic struct {
	value atomic.Value
	fn    func(x, y interface{}) bool
}

// New return a pointer to the new Atomic
func New(f func(x, y interface{}) bool) *Atomic {
	a := &Atomic{
		fn: f,
	}
	return a
}

// Get returns a value
func (a *Atomic) Get() interface{} {
	return a.value.Load()
}

// Update return true if object satisfies fn, false otherwise
func (a *Atomic) Update(value interface{}) bool {
	v := a.value.Load()
	if v == nil || a.fn(v, value) {
		a.value.Store(value)
		return true
	}
	return false
}

// NewMin return a pointer to the new Atomic with minimum predicate
func NewMin() *Atomic {
	a := &Atomic{
		fn: minFn,
	}
	return a
}

// NewMax return a pointer to the new Atomic with maximum predicate
func NewMax() *Atomic {
	a := &Atomic{
		fn: maxFn,
	}
	return a
}

func minFn(x, y interface{}) bool {
	return x.(int) > y.(int)
}

func maxFn(x, y interface{}) bool {
	return x.(int) < y.(int)
}
