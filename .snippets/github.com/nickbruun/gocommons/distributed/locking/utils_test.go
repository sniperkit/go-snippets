package locking

import (
	"testing"
	"time"
)

func AssertLock(t *testing.T, l Lock) <-chan struct{} {
	failCh, err := l.Lock()
	if err != nil {
		t.Fatalf("Unexpected error acquring lock: %v", err)
	}
	return failCh
}

func AssertLockTimedOut(t *testing.T, l Lock, timeout time.Duration) {
	_, err := l.LockTimeout(timeout)
	if err == nil {
		t.Fatal("Expected timeout while acquiring lock, but no error occured")
	} else if err != ErrTimeout {
		t.Fatalf("Unexpected error acquring lock expected to time out: %v", err)
	}
}

func AssertLockBeforeTimeout(t *testing.T, l Lock, timeout time.Duration) <-chan struct{} {
	failCh, err := l.LockTimeout(timeout)
	if err != nil {
		t.Fatalf("Unexpected error acquring lock: %v", err)
	}
	return failCh
}

func AssertUnlock(t *testing.T, l Lock) {
	if err := l.Unlock(); err != nil {
		t.Fatalf("Unexpected error releasing lock: %v", err)
	}
}
