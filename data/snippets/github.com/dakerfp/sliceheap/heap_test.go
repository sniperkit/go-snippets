package sliceheap

import (
	"container/heap"
	"testing"
)

func TestIntHeap(t *testing.T) {
	h := New(make([]int, 0, 256), func(a, b interface{}) bool {
		return a.(int) < b.(int)
	})

	heap.Init(h)
	heap.Push(h, 8)
	heap.Push(h, -1)
	heap.Push(h, 2)
	heap.Push(h, 5)
	heap.Push(h, 8)

	last := heap.Pop(h).(int)
	for h.Len() > 0 {
		x := heap.Pop(h).(int)
		if last > x {
			t.Fatal(x)
		}
		last = x
	}
}
