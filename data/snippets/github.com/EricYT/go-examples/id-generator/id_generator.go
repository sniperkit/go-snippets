package cfs

import (
	"errors"
	"sync"

	"container/heap"
)

// file descriptor id generator
//

type IdGenerator struct {
	sync.Mutex

	max     int
	anchor  int
	freeIds *maxHeap
}

func (ig IdGenerator) Anchor() int { return ig.anchor }
func (ig IdGenerator) Max() int    { return ig.max }

func (ig *IdGenerator) Get() (int, error) {
	ig.Lock()
	defer ig.Unlock()
	if ig.freeIds.Len() == 0 {
		// There is not a id in free id pool, use anchor increment
		if ig.anchor >= ig.max {
			return -1, errors.New("id over flow")
		} else {
			id := ig.anchor
			ig.anchor++
			return id, nil
		}
	} else {
		// pop a id from free id pool
		id := heap.Pop(ig.freeIds).(int)
		return id, nil
	}
}

func (ig *IdGenerator) PutBack(id int) {
	ig.Lock()
	defer ig.Unlock()
	heap.Push(ig.freeIds, id)
	// maybe need decrease anchor
	for ig.freeIds.Len() > 0 {
		currTop := heap.Pop(ig.freeIds).(int)
		if currTop+1 == ig.anchor {
			ig.anchor--
		} else {
			heap.Push(ig.freeIds, currTop)
			return
		}
	}
}

func NewIdGenerator(m, a int) *IdGenerator {
	return &IdGenerator{
		max:     m,
		anchor:  a,
		freeIds: NewMaxHeap(),
	}
}

// An maxHeap is a max-heap of ints.
type maxHeap []int

func (mh maxHeap) Len() int           { return len(mh) }
func (mh maxHeap) Less(i, j int) bool { return mh[i] > mh[j] }
func (mh maxHeap) Swap(i, j int)      { mh[i], mh[j] = mh[j], mh[i] }

func (mh *maxHeap) Push(x interface{}) {
	*mh = append(*mh, x.(int))
}

func (mh *maxHeap) Pop() interface{} {
	old := *mh
	n := len(old)
	x := old[n-1]
	*mh = old[0 : n-1]
	return x
}

func NewMaxHeap() *maxHeap {
	mh := &maxHeap{}
	heap.Init(mh)
	return mh
}
