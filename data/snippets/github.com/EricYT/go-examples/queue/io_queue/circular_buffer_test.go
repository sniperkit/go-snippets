package io_queue

import "testing"

func Test1(t *testing.T) {
	cb := NewCircularBuffer(3)
	if !cb.Empty() {
		t.Fatalf("circular buffer should be empty")
	}
	cb.PushBack(1)
	cb.PushBack(2)
	cb.PushBack(3)
	cb.PushBack(4)
	cb.PushBack(5)

	item := cb.PopFront().(int)
	if item != 2 {
		t.Fatalf("circular buffer should return 3 but we got %d", item)
	}
}
