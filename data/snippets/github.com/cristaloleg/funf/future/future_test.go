package future

import (
	"testing"
	"time"
)

func TestNewFuture(t *testing.T) {
	fn := func() interface{} {
		time.Sleep(100 * time.Millisecond)
		return 42
	}
	f := New(fn)

	if f == nil {
		t.Error("Should be instantiated")
	}

	if f.IsClosed() {
		t.Error("Should be closed already")
	}
	if f.Value() != 42 {
		t.Error("Should be equal to 42")
	}
	if !f.IsClosed() {
		t.Error("Should be closed already")
	}
	if f.Value() != 42 {
		t.Error("Should be equal to 42")
	}
	if !f.IsClosed() {
		t.Error("Should be closed already")
	}
}
