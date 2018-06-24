package main

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

func TestMinDist(t *testing.T) {
	for k, v := range map[string]int{
		"4 3 3 5 7":  6,
		"3 20 30 40": 20} {
		if r := minDist(k); r != v {
			t.Errorf("failed: minDist %s is %d, got %d",
				k, v, r)
		}
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func minDist(q string) (r int) {
	var n int
	s := strings.Fields(q)
	fmt.Sscan(s[0], &n)
	t := make([]int, n)
	for i := 1; i <= n; i++ {
		fmt.Sscan(s[i], &t[i-1])
	}
	sort.Ints(t)
	for i, u := 0, t[n/2]; i < n; i++ {
		r += abs(t[i] - u)
	}
	return r
}
