package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestHighestScore(t *testing.T) {
	for k, v := range map[string]string{
		"72 64 150 | 100 18 33 | 13 250 -6":                "100 250 150",
		"10 25 -30 44 | 5 16 70 8 | 13 1 31 12":            "13 25 70 44",
		"100 6 300 20 10 | 5 200 6 9 500 | 1 10 3 400 143": "100 200 300 400 500"} {
		if r := highestScore(k); r != v {
			t.Errorf("failed: highestScore %s is %s, got %s",
				k, v, r)
		}
	}
}

func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}

func highestScore(s string) string {
	var (
		a int
		r []int
		u []string
	)
	for ix, i := range strings.Split(s, " | ") {
		for jx, j := range strings.Fields(i) {
			fmt.Sscan(j, &a)
			if ix == 0 {
				r = append(r, a)
			} else {
				r[jx] = max(r[jx], a)
			}
		}
	}
	for _, i := range r {
		u = append(u, fmt.Sprint(i))
	}
	return strings.Join(u, " ")
}
