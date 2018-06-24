package main

import (
	"fmt"
	"strings"
	"testing"
)

type tuple struct {
	n, m int
}

func TestFlavius(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{10, 3}: "2 5 8 1 6 0 7 4 9 3",
		tuple{5, 2}:  "1 3 0 4 2"} {
		if r := flavius(k.n, k.m); r != v {
			t.Errorf("failed: flavius %d %d is %s, got %s",
				k.n, k.m, v, r)
		}
	}
}

type node struct {
	value int
	next  *node
}

func flavius(n, m int) string {
	t := make([]string, n)
	tail := node{value: n - 1}
	list := &tail
	for i := n - 2; i >= 0; i-- {
		list = &node{i, list}
	}
	tail.next, list = list, &tail

	for i := 0; i < n; i++ {
		for j := 0; j < m-1; j++ {
			list = list.next
		}
		t[i] = fmt.Sprint(list.next.value)
		list.next = list.next.next
	}
	return strings.Join(t, " ")
}
