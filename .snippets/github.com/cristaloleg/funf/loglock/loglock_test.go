package loglock

import "testing"

func TestNewLogLock(t *testing.T) {
	locked, unlocked := 0, 0
	fLock := func() {
		locked++
	}
	fUnlock := func() {
		unlocked++
	}

	var ll LogLock

	ll.OnLock(fLock)
	ll.OnUnlock(fUnlock)

	ll.OnRLock(fLock)
	ll.OnRUnlock(fUnlock)

	ll.Lock()
	if locked != 1 {
		t.Error("must be locked")
	}
	ll.Unlock()
	if unlocked != 1 {
		t.Error("must be unlocked")
	}

	ll.RLock()
	if locked != 2 {
		t.Error("must be locked")
	}
	ll.RUnlock()
	if unlocked != 2 {
		t.Error("must be unlocked")
	}
}
