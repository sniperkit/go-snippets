package main

import (
	"fmt"
	"strings"
	"testing"
)

type tuple struct {
	a string
	n int
}

func TestNotSoClever(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{"4 3 2 1", 1}:   "3 4 2 1",
		tuple{"5 4 3 2 1", 2}: "4 3 5 2 1"} {
		if r := notSoClever(k.a, k.n); r != v {
			t.Errorf("failed: notSoClever %s | %d is %s, got %s",
				k.a, k.n, v, r)
		}
	}
}

func notSoClever(q string, n int) string {
	var j, k int
	u := strings.Fields(q)
	v := make([]int, len(u))
	for ix, i := range u {
		fmt.Sscan(i, &k)
		v[ix] = k
	}
	l := 1
	for i := 0; i < n; i++ {
		for j, l = l, 0; j < len(v); j++ {
			if v[j-1] <= v[j] {
				continue
			}
			v[j-1], v[j] = v[j], v[j-1]
			if j > 1 {
				l = j - 1
			} else {
				l = 2
			}
			break
		}
		if l == 0 {
			break
		}
	}
	for ix, i := range v {
		u[ix] = fmt.Sprint(i)
	}
	return strings.Join(u, " ")
}
