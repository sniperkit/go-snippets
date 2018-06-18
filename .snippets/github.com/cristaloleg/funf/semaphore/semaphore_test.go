package semaphore

import (
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	size := 10
	s := New(size)
	if s == nil {
		t.Error("Should be instantiated")
	}

	s.Lock()

	if s.Count() != size-1 {
		t.Error("Should be locked by 1")
	}

	s.Unlock()

	if s.Count() != size {
		t.Error("Should be locked by no one")
	}
}

func TestSemaphore_TryLock(t *testing.T) {
	size := 10
	s := New(size)
	if s == nil {
		t.Error("Should be instantiated")
	}

	for i := 0; i < size-1; i++ {
		s.Lock()
	}

	if s.Count() != 1 {
		t.Error("Must be free for 1")
	}

	if !s.TryLock() {
		t.Error("Should be locked")
	}

	if s.TryLock() {
		t.Error("Cannot be locked")
	}
}

func TestSemaphore_TryDelayedLock(t *testing.T) {
	size := 10
	s := New(size)
	if s == nil {
		t.Error("Should be instantiated")
	}

	for i := 0; i < size; i++ {
		s.Lock()
	}

	if s.Count() != 0 {
		t.Error("Must be totally locked")
	}

	go func() {
		<-time.After(150 * time.Millisecond)
		s.Unlock()
	}()

	if s.TryLock() {
		t.Error("Should be locked")
	}

	if s.TryDelayedLock(10 * time.Millisecond) {
		t.Error("Should be locked")
	}

	if !s.TryDelayedLock(200 * time.Millisecond) {
		t.Error("Should be locked")
	}
}
