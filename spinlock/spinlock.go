package spinlock

import (
	"runtime"
	"sync/atomic"
)

// SpinLock is a spinned lock
type SpinLock struct {
	lock int32
}

// New return pointer to new SpinLock
func New() *SpinLock {
	s := &SpinLock{
		lock: 0,
	}
	return s
}

// Lock will lock SpinLock
func (s *SpinLock) Lock() {
	for !s.TryLock() {
		runtime.Gosched()
	}
}

// Unlock will unlock SpinLock
func (s *SpinLock) Unlock() {
	atomic.StoreInt32(&s.lock, 0)
}

// TryLock will try to lock SpinLock
func (s *SpinLock) TryLock() bool {
	return atomic.CompareAndSwapInt32(&s.lock, 0, 1)
}

// IsLocked returns true if SpinLock is locked
func (s *SpinLock) IsLocked() bool {
	return atomic.LoadInt32(&s.lock) == 1
}

// TryDelayedLock tries to lock semaphore until delay, returns true if succeeded.
// func (s *SpinLock) TryDelayedLock(delay time.Duration) bool {
// 	for {
// 		select {
// 		case <-time.After(delay):
// 			return false

// 		default:
// 			if !s.TryLock() {
// 				runtime.Gosched()
// 				continue
// 			}
// 			return true
// 		}
// 	}
// 	return false
// }
