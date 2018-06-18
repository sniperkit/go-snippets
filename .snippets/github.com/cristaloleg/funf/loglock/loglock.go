package loglock

import (
	"sync"
	"time"
)

// LogLock is a logged lock
type LogLock struct {
	mutex sync.RWMutex

	locked    int64
	unlocked  int64
	rlocked   int64
	runlocked int64

	lastLock    time.Time
	lastUnlock  time.Time
	lastRLock   time.Time
	lastRUnlock time.Time

	onLock    func()
	onUnlock  func()
	onRLock   func()
	onRUnlock func()
}

// Lock locks for read/write
func (ll *LogLock) Lock() {
	ll.mutex.Lock()
	ll.locked++
	ll.lastLock = time.Now()
	if ll.onLock != nil {
		ll.onLock()
	}
}

// Unlock unlocks for read/write
func (ll *LogLock) Unlock() {
	ll.mutex.Unlock()
	ll.unlocked++
	ll.lastLock = time.Now()
	if ll.onUnlock != nil {
		ll.onUnlock()
	}
}

// RLock locks for reading
func (ll *LogLock) RLock() {
	ll.mutex.RLock()
	ll.rlocked++
	ll.lastRLock = time.Now()
	if ll.onRLock != nil {
		ll.onRLock()
	}
}

// RUnlock unlocks reading
func (ll *LogLock) RUnlock() {
	ll.mutex.RUnlock()
	ll.runlocked++
	ll.lastRUnlock = time.Now()
	if ll.onRUnlock != nil {
		ll.onRUnlock()
	}
}

// OnLock executes on locking
func (ll *LogLock) OnLock(f func()) {
	ll.onLock = f
}

// OnUnlock executes on unlocking
func (ll *LogLock) OnUnlock(f func()) {
	ll.onUnlock = f
}

// OnRLock executes on read locking
func (ll *LogLock) OnRLock(f func()) {
	ll.onRLock = f
}

// OnRUnlock executes on read unlocking
func (ll *LogLock) OnRUnlock(f func()) {
	ll.onRUnlock = f
}

// Counters returrns counters for lock toggles
func (ll *LogLock) Counters() (lock, unlock, rlock, runlock int64) {
	return ll.locked, ll.unlocked, ll.rlocked, ll.runlocked
}

// LastToggles returrns last lock toggles
func (ll *LogLock) LastToggles() (lock, unlock, rlock, runlock time.Time) {
	return ll.lastLock, ll.lastUnlock, ll.lastRLock, ll.lastRUnlock
}
