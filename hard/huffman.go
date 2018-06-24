package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"strings"
)

type treenode struct {
	leftchild  *symbol
	rightchild *symbol
	symbols    string
}

type symbol struct {
	priority int
	tree     treenode
}

type symbolHeap []symbol

func (h symbolHeap) Len() int { return len(h) }
func (h symbolHeap) Less(i, j int) bool {
	if h[i].priority == h[j].priority {
		if len(h[i].tree.symbols) == len(h[j].tree.symbols) {
			return h[i].tree.symbols < h[j].tree.symbols
		}
		return len(h[i].tree.symbols) > len(h[j].tree.symbols)
	}
	return h[i].priority < h[j].priority
}
func (h symbolHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *symbolHeap) Push(x interface{}) {
	*h = append(*h, x.(symbol))
}
func (h *symbolHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

var m map[byte]string

func traverse(prefix string, t symbol) {
	if t.tree.leftchild == nil && t.tree.rightchild == nil {
		m[t.tree.symbols[0]] = prefix
		return
	}
	if t.tree.leftchild != nil {
		traverse(prefix+"0", *(t.tree.leftchild))
	}
	if t.tree.rightchild != nil {
		traverse(prefix+"1", *(t.tree.rightchild))
	}
	return
}

func main() {
	var (
		a rune
		i int
	)
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		c := make([]int, 26)
		for _, a = range scanner.Text() {
			c[a-'a']++
		}
		d := &symbolHeap{}
		heap.Init(d)
		for i = 0; i < 26; i++ {
			if c[i] > 0 {
				heap.Push(d, symbol{c[i], treenode{nil, nil, string(rune('a' + i))}})
			}
		}
		for d.Len() > 1 {
			var t, u symbol
			t = heap.Pop(d).(symbol)
			u = heap.Pop(d).(symbol)
			heap.Push(d, symbol{t.priority + u.priority, treenode{&t, &u, t.tree.symbols + u.tree.symbols}})
		}
		m = make(map[byte]string, 26)
		v := heap.Pop(d).(symbol)
		traverse("", v)

		var ret []string
		for i = 0; i < 26; i++ {
			if r, f := m[byte('a'+i)]; f {
				ret = append(ret, fmt.Sprintf("%c: %s;", rune('a'+i), r))
			}
		}
		fmt.Println(strings.Join(ret, " "))
	}
}
