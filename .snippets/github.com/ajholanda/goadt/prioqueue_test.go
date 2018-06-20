package adt_test

import (
	adt "github.com/ajholanda/goadt"
	
	"fmt"
	"testing"
)

var pq *adt.PriorityQueue

func initPQ() {
	items := map[string]int{"2^3" : 8, "2^1" : 2, "2^2" : 4}
	for name, prio := range items {
		pq.Push(prio, name, nil)
	}
}


func TestMinPriorityQueue(t *testing.T) {
	expPriorities := []int{1, 2, 4, 8}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq = adt.NewPriorityQueue(adt.MIN)

	initPQ()
	
	// Insert a new item and then modify its distance.
	pq.Push(10, "2^0", nil)
	pq.Update(1, "2^0", nil)

	// Take the items out; they arrive in decreasing priority order.
	i := 0
	for pq.IsEmpty() == false {
		item := pq.Pop()

		if item.Prio() != expPriorities[i] {
			t.Errorf("heap.Pop(pq)=%d, want %d", item.Prio(), expPriorities[i])
		}
		i++
	}
}

func TestMaxPriorityQueue(t *testing.T) {
	expPriorities := []int{8, 4, 2, 1}

	fmt.Println()
	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq = adt.NewPriorityQueue(adt.MAX)

	initPQ()
	
	// Insert a new item and then modify its distance.
	pq.Push(10, "2^0", nil)
	pq.Update(1, "2^0", nil)

	// Take the items out; they arrive in decreasing priority order.
	i := 0
	for pq.IsEmpty() == false {
		hn := pq.Pop() // heap node

		if hn.Index() != 1 {
			t.Errorf("heap.Pop(pq).index=%d, want %d", hn.Index(), 1)
		}
		
		if hn.Prio() != expPriorities[i] {
			t.Errorf("heap.Pop(pq).prio=%d, want %d", hn.Prio(), expPriorities[i])
		}
		i++
	}
}
