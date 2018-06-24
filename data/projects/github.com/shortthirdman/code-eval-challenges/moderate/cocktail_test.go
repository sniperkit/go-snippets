package main

import (
	"fmt"
	"strings"
	"testing"
)

type tuple struct {
	q string
	n int
}

func TestCocktail(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{"5 4 9 10 7 3 2 1 6", 1}: "1 4 5 9 7 3 2 6 10",
		tuple{"9 8 7 6 5 4 3 2 1", 3}:  "1 2 3 6 5 4 7 8 9"} {
		if r := cocktail(k.q, k.n); r != v {
			t.Errorf("failed: cocktail %s | %d id %s, got %s",
				k.q, k.n, v, r)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func cocktail(q string, n int) string {
	u := strings.Fields(q)
	n = min(n, len(u)/2)
	v := make([]int, len(u))
	for ix, i := range u {
		fmt.Sscan(i, &v[ix])
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < len(v)-i; j++ {
			if v[j-1] > v[j] {
				v[j-1], v[j] = v[j], v[j-1]
			}
		}
		for j := len(v) - i - 2; j >= i+1; j-- {
			if v[j-1] > v[j] {
				v[j-1], v[j] = v[j], v[j-1]
			}
		}
	}
	for ix, i := range v {
		u[ix] = fmt.Sprint(i)
	}
	return strings.Join(u, " ")
}
