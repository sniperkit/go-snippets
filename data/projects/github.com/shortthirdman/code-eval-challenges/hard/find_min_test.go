package main

import (
	"container/heap"
	"sort"
	"testing"
)

type tuple struct {
	n, k, a, b, c, r int
}

func TestFindMin(t *testing.T) {
	for k, v := range map[tuple]int{
		tuple{97, 39, 34, 37, 656, 97}:   8,
		tuple{186, 75, 68, 16, 539, 186}: 38,
		tuple{137, 49, 48, 17, 461, 137}: 41,
		tuple{98, 59, 6, 30, 524, 98}:    40,
		tuple{46, 18, 7, 11, 9, 46}:      12} {
		if r := findMin(k.n, k.k, k.a, k.b, k.c, k.r); r != v {
			t.Errorf("failed: findMin %d,%d,%d,%d,%d,%d is %d, got %d",
				k.n, k.k, k.a, k.b, k.c, k.r, v, r)
		}
	}
}

type intHeap []int

func (h intHeap) Len() int           { return len(h) }
func (h intHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h intHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *intHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}
func (h *intHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h intHeap) notyet(x int) bool {
	for _, i := range h {
		if i == x {
			return false
		}
	}
	return true
}

func notagain(t []int, x int) bool {
	for _, i := range t {
		if i == x {
			return false
		}
	}
	return true
}

func findMin(n, k, a, b, c, r int) int {
	m := []int{a}
	for i := 1; i < k; i++ {
		m = append(m, (b*m[i-1]+c)%r)
	}
	o := make([]int, k)
	copy(o, m)
	sort.Ints(o)
	h := &intHeap{}
	heap.Init(h)
	var x, y int
	for i := 0; i <= k; {
		if y >= k || x < o[y] {
			heap.Push(h, x)
			x++
			i++
		} else {
			if x == o[y] {
				x++
			}
			y++
		}
	}
	for len(m)+1 < n {
		p := heap.Pop(h).(int)
		if h.notyet(m[len(m)-k]) && notagain(m[len(m)-k+1:len(m)], m[len(m)-k]) {
			heap.Push(h, m[len(m)-k])
		}
		m = append(m, p)
	}
	return heap.Pop(h).(int)
}
