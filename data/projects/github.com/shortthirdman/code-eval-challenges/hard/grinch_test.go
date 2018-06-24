package main

import (
	"container/heap"
	"fmt"
	"math"
	"strings"
	"testing"
)

type tuple struct {
	q    string
	a, b uint
}

func TestGrinch(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{"1 2 2, 1 3 3, 3 4 3, 2 4 6, 4 5 16, 3 5 7", 1, 5}: "10",
		tuple{"1 2 3, 2 8 10, 1 9 4, 8 9 2", 2, 8}:               "9",
		tuple{"1 2 3, 4 5 6", 1, 5}:                              "False"} {
		if r := grinch(k.q, k.a, k.b); r != v {
			t.Errorf("failed: grinch %s %d %d is %s, got %s",
				k.q, k.a, k.b, v, r)
		}
	}
}

const inf = math.MaxInt32

type vertex struct {
	id, dist uint
}

type edge struct {
	to, dist uint
}

type vertexHeap []vertex

func (h vertexHeap) Len() int           { return len(h) }
func (h vertexHeap) Less(i, j int) bool { return h[i].dist < h[j].dist }
func (h vertexHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *vertexHeap) Push(x interface{}) {
	*h = append(*h, x.(vertex))
}
func (h *vertexHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h vertexHeap) FindID(a uint) int {
	for ix, i := range h {
		if i.id == a {
			return ix
		}
	}
	return -1
}

func grinch(q string, a, b uint) string {
	var d uint = inf
	t := strings.Split(q, ", ")
	p := &vertexHeap{}
	heap.Init(p)
	neigh := make(map[uint][]edge)
	var v, x, y, z uint
	heap.Push(p, vertex{a, 0})
	v |= 1 << a
	if b != a {
		heap.Push(p, vertex{b, inf})
		v |= 1 << b
	}
	for _, i := range t {
		fmt.Sscanf(i, "%d %d %d", &x, &y, &z)
		neigh[x] = append(neigh[x], edge{y, z})
		neigh[y] = append(neigh[y], edge{x, z})
		if v&(1<<x) == 0 {
			heap.Push(p, vertex{x, inf})
			v |= 1 << x
		}
		if v&(1<<y) == 0 {
			heap.Push(p, vertex{y, inf})
			v |= 1 << y
		}
	}

	for p.Len() > 0 {
		u := heap.Pop(p).(vertex)
		if u.id == b {
			d = u.dist
			break
		}
		for _, i := range neigh[u.id] {
			dis := u.dist + i.dist
			ix := p.FindID(i.to)
			if ix >= 0 && dis < (*p)[ix].dist {
				(*p)[ix].dist = dis
				heap.Fix(p, ix)
			}
		}
	}

	if d == inf {
		return "False"
	}
	return fmt.Sprint(d)
}
