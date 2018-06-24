package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"sort"
)

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

func main() {
	var n, k, a, b, c, r int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d,%d,%d,%d,%d,%d", &n, &k, &a, &b, &c, &r)
		fmt.Println(findMin(n, k, a, b, c, r))
	}
}
