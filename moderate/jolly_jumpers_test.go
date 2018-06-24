package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestJolly(t *testing.T) {
	for k, v := range map[string]bool{
		"4 1 4 2 3":        true,
		"5 76 34 75 78 77": false,
		"5 76 74 75 78 74": true,
		"6 7 5 8 5 6 9":    false,
		"1 5":              true,
		"2 6 5":            true,
		"3 3 3 3":          false} {
		if r := jolly(k); r != v {
			t.Errorf("failed: jolly %s is %t, got %t",
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

func jolly(q string) bool {
	var v, vp int
	s := strings.Fields(q)[1:]
	n := len(s)
	if n == 1 {
		return true
	}
	u := make([]bool, n-1)
	fmt.Sscan(s[0], &vp)
	for i := 1; i < n; i++ {
		fmt.Sscan(s[i], &v)
		x := abs(v - vp)
		if x >= n || x == 0 || u[x-1] {
			return false
		}
		u[x-1], vp = true, v
	}
	return true
}
