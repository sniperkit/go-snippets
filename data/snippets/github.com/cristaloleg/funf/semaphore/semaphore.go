package semaphore

import "time"

// Semaphore is a counting semaphore.
type Semaphore struct {
	locks chan struct{}
}

// New return new semaphore with `count` locks.
func New(count int) *Semaphore {
	s := &Semaphore{
		locks: make(chan struct{}, count),
	}
	for i := 0; i < count; i++ {
		s.locks <- struct{}{}
	}
	return s
}

// Lock decreases available locks by 1. It's a blocking operation.
func (s *Semaphore) Lock() {
	<-s.locks
}

// Unlock increases available locks by 1.
func (s *Semaphore) Unlock() {
	s.locks <- struct{}{}
}

// TryLock tries to lock semaphore, returns true if succeeded.
func (s *Semaphore) TryLock() bool {
	select {
	case <-s.locks:
		return true
	default:
		return false
	}
}

// TryDelayedLock tries to lock semaphore until delay, returns true if succeeded.
func (s *Semaphore) TryDelayedLock(delay time.Duration) bool {
	select {
	case <-s.locks:
		return true
	case <-time.After(delay):
		return false
	}
}

// Count returns number of locks for the semaphore.
func (s *Semaphore) Count() int {
	return len(s.locks)
}
