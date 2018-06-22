package lock

import (
	"sync"
	"testing"
)

type myPanel struct {
	sync.Mutex
	data int
}

// This Test will block forever
func TestLock(t *testing.T) {
	var l myPanel
	l.Lock()
	l.Lock()
	l.Unlock()
	l.Unlock()
}
