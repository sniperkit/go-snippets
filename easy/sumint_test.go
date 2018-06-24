package main

import "testing"

type tuple struct {
	n, s int
}

func TestSumint(t *testing.T) {
	for k, v := range map[tuple]int{
		tuple{3, 0}:    3,
		tuple{23, 3}:   26,
		tuple{321, 26}: 347,
		tuple{21, 0}:   21,
		tuple{14, 21}:  35,
		tuple{35, 7}:   42} {
		if r := sumint(k.n, k.s); r != v {
			t.Errorf("failed: sumint %d %d is %d, got %d",
				k.n, k.s, v, r)
		}
	}
}

func sumint(n, s int) int {
	return n + s
}
