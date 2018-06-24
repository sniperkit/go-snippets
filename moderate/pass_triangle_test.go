package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestPass(t *testing.T) {
	for k, v := range map[string]int{
		"5\n9 6\n4 6 8\n0 7 1 5": 27} {
		u = []int{}
		for _, i := range strings.Split(k, "\n") {
			pass(i)
		}
		if r := maxis(u); r != v {
			t.Errorf("failed: pass %s is %d, got %d",
				k, v, r)
		}
	}
}

var u []int

func maxi(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxis(a []int) (r int) {
	r = a[0]
	for i := 1; i < len(a); i++ {
		if a[i] > r {
			r = a[i]
		}
	}
	return r
}

func pass(q string) {
	s := strings.Fields(q)
	t := make([]int, len(s))
	for ix, i := range s {
		fmt.Sscan(i, &t[ix])
	}
	if len(u) > 0 {
		t[0] += u[0]
		if len(u) > 1 {
			t[len(u)] += u[len(u)-1]
		}
	}
	for i := 1; i < len(t)-1; i++ {
		t[i] += maxi(u[i-1], u[i])
	}
	u = make([]int, len(t))
	copy(u, t)
	return
}
