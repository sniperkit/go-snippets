package locking

import (
	log "github.com/nickbruun/gocommons/logging"
	"github.com/nickbruun/gocommons/unittest"
	"sync"
	"time"
)

// Lock test suite.
type LockTestSuite struct {
	unittest.TestSuite

	LockProvider LockProvider
}

// Assert that a lock is acquired by a call to Lock.
func (s *LockTestSuite) AssertLock(l Lock) (failCh <-chan struct{}) {
	var err error
	if failCh, err = l.Lock(); err != nil {
		s.Fatalf("Unexpected error acquring lock: %v", err)
	}
	return
}

// Assert that the attempt to acquire a lock with timeout times out.
func (s *LockTestSuite) AssertLockTimedOut(l Lock, timeout time.Duration) {
	_, err := l.LockTimeout(timeout)
	if err == nil {
		s.Fatal("Expected timeout while acquiring lock, but no error occured")
	} else if err != ErrTimeout {
		s.Fatalf("Unexpected error acquring lock expected to time out: %v", err)
	}
}

// Assert that a lock is acquired before a timeout.
func (s *LockTestSuite) AssertLockBeforeTimeout(l Lock, timeout time.Duration) <-chan struct{} {
	failCh, err := l.LockTimeout(timeout)
	if err != nil {
		s.Fatalf("Unexpected error acquring lock: %v", err)
	}
	return failCh
}

// Assert that a lock is unlocked successfully.
func (s *LockTestSuite) AssertUnlock(l Lock) {
	if err := l.Unlock(); err != nil {
		s.Fatalf("Unexpected error releasing lock: %v", err)
	}
}

// Assert that unlocking a lock returns that the lock is not held.
func (s *LockTestSuite) AssertUnlockNotHeld(l Lock) {
	err := l.Unlock()

	if err == nil {
		s.Fatal("Expected ErrNotLocked, but not error occured")
	} else if err != ErrNotLocked {
		s.Fatalf("Expected ErrNotLock, but an unexpected error occured: %v", err)
	}
}

func (s *LockTestSuite) TestLock() {
	// Test simple locking and unlocking.
	{
		l := s.LockProvider.GetLock("/my/lock")

		for i := 0; i < 3; i++ {
			s.AssertLock(l)
			log.Debug("Acquired lock")
			s.AssertUnlock(l)
			log.Debug("Released lock")
			s.AssertLockBeforeTimeout(l, time.Second)
			log.Debug("Acquired lock before 1 s timeout")
			s.AssertUnlock(l)
			log.Debug("Released lock")
		}
	}

	// Test that locks that have timed out are skipped when acquring locks.
	{
		l1 := s.LockProvider.GetLock("/lock/path")
		l2 := s.LockProvider.GetLock("/lock/path")
		l3 := s.LockProvider.GetLock("/lock/path")
		l4 := s.LockProvider.GetLock("/lock/path")

		for i := 0; i < 3; i++ {
			var l2TimedOut, l3Acquired, l4Acquired sync.WaitGroup
			l2TimedOut.Add(1)
			l3Acquired.Add(1)
			l4Acquired.Add(1)

			s.AssertLock(l1)
			log.Debug("Acquired lock 1, waiting for lock 2 to time out")

			go func() {
				s.AssertLockTimedOut(l2, 10*time.Millisecond)
				log.Debug("Lock 2 timed out")
				l2TimedOut.Done()
			}()

			go func() {
				s.AssertLock(l3)
				log.Debug("Lock 3 acquired")
				l3Acquired.Done()
			}()

			go func() {
				s.AssertLock(l4)
				log.Debug("Lock 4 acquired")
				l4Acquired.Done()
			}()

			l2TimedOut.Wait()

			s.AssertUnlock(l1)
			log.Debug("Released lock 1, waiting for lock 3 to be acquired")

			l3Acquired.Wait()

			s.AssertUnlock(l3)
			log.Debug("Released lock 3, waiting for lock 4 to be acquired")

			l4Acquired.Wait()

			s.AssertUnlock(l4)
			log.Debug("Released lock 4")
		}
	}
}
