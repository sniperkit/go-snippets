package cfs

import (
	"container/heap"
	"fmt"
	"testing"
)

func TestMaxHeap(t *testing.T) {
	mh := NewMaxHeap()
	for i := 0; i < 10; i++ {
		heap.Push(mh, i)
	}
	top := heap.Pop(mh).(int)
	if top != 9 {
		t.Fatalf("This is not a max heap top is %d not 9\n", top)
	}
}

func TestIdGenerator(t *testing.T) {
	idg := NewIdGenerator(1000, 0)
	ids := make([]int, 100)
	for i := 0; i < 100; i++ {
		id, _ := idg.Get()
		ids[i] = id
	}
	fmt.Println("(1)id generator anchor ", idg.Anchor())
	var index int
	for index = 0; index < 98; index++ {
		idg.PutBack(ids[index])
	}
	fmt.Println("ids:", ids)
	fmt.Printf("id ids[%d + 1] = %d\n", index, ids[index+1])
	idg.PutBack(ids[index+1])
	fmt.Println("(2)id generator anchor ", idg.Anchor())
	if idg.Anchor() != ids[index+1] {
		t.Fatal("(1)id generator anchor error:", idg.Anchor())
	}
	idg.PutBack(ids[index])
	fmt.Println("(3)id generator anchor ", idg.Anchor())
	if idg.Anchor() != ids[0] {
		t.Fatal("(2)id generator anchor error:", idg.Anchor())
	}
}
