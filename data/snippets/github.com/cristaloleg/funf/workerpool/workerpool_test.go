package workerpool

import "testing"

func TestNewWorkersPool(t *testing.T) {
	fn := func(interface{}) interface{} { return nil }
	data := make(chan interface{})
	size := 4
	w := NewBuffered(size, fn, data)

	if w == nil {
		t.Error("Should be instantiated")
	}
}
