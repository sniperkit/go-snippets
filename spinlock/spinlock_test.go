package spinlock_test

import (
	"testing"
	"time"

	"github.com/cristaloleg/funf/spinlock"
)

func TestSpinLock(t *testing.T) {
	s := spinlock.New()
	if s == nil {
		t.Error("Should be instantiated")
	}

	s.Lock()

	if !s.IsLocked() {
		t.Error("Should be locked by 1")
	}

	s.Unlock()

	if s.IsLocked() {
		t.Error("Should be locked by no one")
	}
}

func TestSpinLock_TryLock(t *testing.T) {
	s := spinlock.New()
	if s == nil {
		t.Error("Should be instantiated")
	}

	s.Lock()

	if !s.IsLocked() {
		t.Error("Must be free for 1")
	}

	if s.TryLock() {
		t.Error("Should be locked")
	}

	if s.TryLock() {
		t.Error("Cannot be locked")
	}
}

func TestSpinLock_TryDelayedLock(t *testing.T) {
	s := spinlock.New()
	if s == nil {
		t.Error("Should be instantiated")
	}

	s.Lock()

	if !s.IsLocked() {
		t.Error("Must be totally locked")
	}

	go func() {
		<-time.After(150 * time.Millisecond)
		s.Unlock()
	}()

	if s.TryLock() {
		t.Error("Should be locked")
	}

	// if s.TryDelayedLock(10 * time.Millisecond) {
	// 	t.Error("Should be locked")
	// }

	// if !s.TryDelayedLock(200 * time.Millisecond) {
	// 	t.Error("Should be locked")
	// }
}
