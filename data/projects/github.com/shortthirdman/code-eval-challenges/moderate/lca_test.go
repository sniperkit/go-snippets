package main

import "testing"

type tuple struct {
	a, b int
}

func TestLca(t *testing.T) {
	for k, v := range map[tuple]int{
		tuple{8, 52}:  30,
		tuple{3, 29}:  8,
		tuple{52, 29}: 30,
		tuple{29, 30}: 30} {
		if r := lca(k.a, k.b); r != v {
			t.Errorf("failed: lca %d %d is %d, got %d",
				k.a, k.b, v, r)
		}
	}
}

var parent map[int]int

func init() {
	parent = map[int]int{3: 8, 8: 30, 10: 20, 20: 8, 29: 20, 30: 0, 52: 30}
}

func lca(a, b int) int {
	for ; a != 0 && a != b; a = parent[a] {
		for i := b; i != 0; i = parent[i] {
			if a == i {
				return a
			}
		}
	}
	return a
}
