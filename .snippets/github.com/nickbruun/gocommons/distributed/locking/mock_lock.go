package locking

import (
	"sync"
	"time"
)

// Mock lock.
//
// Should only be used for testing.
type mockLock struct {
	p      *MockLockProvider
	path   string
	locked bool
	failCh chan struct{}
	lock   sync.Mutex
}

func (l *mockLock) Lock() (<-chan struct{}, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.locked {
		return nil, ErrLocked
	}

	l.failCh = make(chan struct{}, 1)
	l.p.acquire(l)
	l.locked = true

	return l.failCh, nil
}

func (l *mockLock) LockTimeout(timeout time.Duration) (<-chan struct{}, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.locked {
		return nil, ErrLocked
	}

	l.failCh = make(chan struct{}, 1)
	l.locked = l.p.acquireTimeout(l, timeout)

	if l.locked {
		return l.failCh, nil
	} else {
		return nil, ErrTimeout
	}
}

func (l *mockLock) Unlock() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if !l.locked {
		return ErrNotLocked
	}

	l.p.release(l)
	l.locked = false
	return nil
}

func (l *mockLock) fail() {
	l.lock.Lock()
	defer l.lock.Unlock()

	if !l.locked {
		return
	}

	l.locked = false
	l.failCh <- struct{}{}
}
