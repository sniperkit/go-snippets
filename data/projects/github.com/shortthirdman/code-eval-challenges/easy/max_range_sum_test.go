package main

import (
	"fmt"
	"strings"
	"testing"
)

type tuple struct {
	n int
	q string
}

func TestMaxRangeSum(t *testing.T) {
	for k, v := range map[tuple]int{
		tuple{5, "7 -3 -10 4 2 8 -2 4 -5 -2"}: 16,
		tuple{6, "-4 3 -10 5 3 -7 -3 7 -6 3"}: 0,
		tuple{3, "-7 0 -45 34 -24 7"}:         17} {
		if r := maxRangeSum(k.n, k.q); r != v {
			t.Errorf("failed: maxRangeSum %d;%s is %d, got %d",
				k.n, k.q, v, r)
		}
	}
}

func maxRangeSum(n int, q string) (r int) {
	t := strings.Fields(q)
	var c int
	u := make([]int, len(t))
	for ix, i := range t {
		fmt.Sscan(i, &u[ix])
	}
	for i := 0; i < n; i++ {
		c += u[i]
	}
	if c > r {
		r = c
	}
	for len(u) > n {
		c = c - u[0] + u[n]
		if c > r {
			r = c
		}
		u = u[1:]
	}
	return r
}
