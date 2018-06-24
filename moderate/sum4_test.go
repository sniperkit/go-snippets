package main

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

func TestNumzero(t *testing.T) {
	if r := numzero([]int{-4, -1, 0, 1, 2, 3}, 4, 0); r != 2 {
		t.Errorf("failed: numzero test 1 is 2, got %d",
			r)
	}
	if r := numzero([]int{-2, -1, 0, 3}, 4, 0); r != 1 {
		t.Errorf("failed: numzero test 2 is 1, got %d",
			r)
	}
	if r := numzero([]int{0, 0, 0, 0, 0, 0, 0, 0}, 4, 0); r != 70 {
		t.Errorf("failed: numzero test 3 is 70, got %d",
			r)
	}
}

func TestSum4(t *testing.T) {
	for k, v := range map[string]int{
		"2,3,1,0,-4,-1":   2,
		"0,-1,3,-2":       1,
		"0,0,0,0,0,0,0,0": 70} {
		if r := sum4(k); r != v {
			t.Errorf("failed: sum4 %s is %d, got %d",
				k, v, r)
		}
	}
}

func numzero(t []int, c, z int) int {
	switch {
	case c == 0:
		if z == 0 {
			return 1
		}
		return 0
	case len(t) < c || z+c*t[0] > 0 || z+c*t[len(t)-1] < 0:
		return 0
	default:
		return numzero(t[1:len(t)], c-1, z+t[0]) + numzero(t[1:len(t)], c, z)
	}
}

func sum4(q string) int {
	s := strings.Split(q, ",")
	t := make([]int, len(s))
	for ix, i := range s {
		fmt.Sscan(i, &t[ix])
	}
	sort.Ints(t)
	return numzero(t, 4, 0)
}
