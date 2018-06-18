package locking

import (
	log "github.com/nickbruun/gocommons/logging"
	"sync"
	"time"
)

// Internal mock lock provider lock representation.
type mockLockProviderLock struct {
	locker  *mockLock
	waiters []chan struct{}
}

func (l *mockLockProviderLock) OfferToNext() {
	if len(l.waiters) == 0 {
		return
	}

	waiter := l.waiters[0]
	l.waiters = l.waiters[1:]

	waiter <- struct{}{}
}

// Mock lock provider.
//
// Provides mock locks. All mock locks must use the same provider for mutual
// exclusivity to be in effect.
type MockLockProvider struct {
	failing bool
	locks   map[string]*mockLockProviderLock
	lock    sync.Mutex
}

// Fail the lock provider.
//
// Fails any currently held locks and does not allow new locks to be acquried
// until the lock provider is recovered.
func (p *MockLockProvider) Fail() {
	p.lock.Lock()
	defer p.lock.Unlock()

	log.Debug("Failing provider")

	p.failing = true

	for _, il := range p.locks {
		if il.locker != nil {
			il.locker.fail()
			il.locker = nil
		}
	}

	log.Debug("Provider failed")
}

// Recover the lock provider.
func (p *MockLockProvider) Recover() {
	p.lock.Lock()
	defer p.lock.Unlock()

	log.Debug("Recovering provider")

	p.failing = false

	for path, il := range p.locks {
		log.Debugf("Offering lock to waiters for %s: %v", path, il.waiters)
		il.OfferToNext()
	}

	log.Debug("Provider recovered")
}

func (p *MockLockProvider) GetLock(path string) Lock {
	return &mockLock{
		p:    p,
		path: path,
	}
}

// Acquire a lock.
func (p *MockLockProvider) acquire(lock *mockLock) {
	// To avoid duplication, we just acquire a lock with a really long timeout.
	if !p.acquireTimeout(lock, 24*time.Hour) {
		panic("Lock with long timeout to emulate no timeout should always lock")
	}
}

// Acquire a lock with timeout.
func (p *MockLockProvider) acquireTimeout(lock *mockLock, timeout time.Duration) bool {
	p.lock.Lock()

	// Get the internal lock representation.
	il, exists := p.locks[lock.path]
	if !exists {
		il = &mockLockProviderLock{
			locker:  nil,
			waiters: make([]chan struct{}, 0),
		}
		p.locks[lock.path] = il
	}

	// If the lock is already held, or the provider is failing, enqueue the
	// lock.
	if il.locker != nil || p.failing {
		if il.locker != nil {
			log.Debugf("Lock is currently held by %v, waiting for lock to become available or time out after %s", lock, timeout)
		} else {
			log.Debugf("Lock provider is currently failing, waiting for recovery or time out after %s", timeout)
		}

		var wg sync.WaitGroup
		wg.Add(1)
		wc := make(chan struct{}, 1)

		il.waiters = append(il.waiters, wc)
		p.lock.Unlock()

		timedOut := false

		go func() {
			select {
			case <-wc:
				log.Debugf("Lock for %s was released, accepting offer", lock.path)
			case <-time.After(timeout):
				log.Debugf("Timed out waiting for lock for %s after %s", lock.path, timeout)
				timedOut = true
			}

			wg.Done()
		}()

		wg.Wait()

		// If a timeout occurs, make sure the lock gets offered to the next in
		// line if relevant.
		if timedOut {
			go func() {
				<-wc

				log.Debugf("Lock for %s was released, but request has timed out, passing on the offer", lock.path)

				p.lock.Lock()
				defer p.lock.Unlock()

				il.OfferToNext()
			}()

			return false
		}

		p.lock.Lock()
	}

	// Lock.
	il.locker = lock

	p.lock.Unlock()

	return true
}

// Release a lock.
func (p *MockLockProvider) release(lock *mockLock) {
	p.lock.Lock()
	defer p.lock.Unlock()

	il, exists := p.locks[lock.path]
	if !exists || il.locker != lock {
		panic("Releasing lock not held by the provided locker")
	}

	il.locker = nil
	il.OfferToNext()
}

// New mock lock provider.
func NewMockLockProvider() *MockLockProvider {
	return &MockLockProvider{
		locks: make(map[string]*mockLockProviderLock),
	}
}
