package main

import "testing"

type tuple struct {
	n, m int
}

func TestLocks(t *testing.T) {
	for k, v := range map[tuple]int{
		tuple{3, 1}:     2,
		tuple{100, 100}: 50,
		tuple{10, 10}:   5,
		tuple{10, 7}:    7} {
		if r := locks(k.n, k.m); r != v {
			t.Errorf("failed: locks %d %d is %d, got %d",
				k.n, k.m, v, r)
		}
	}
}

func locks(n, m int) int {
	if n == 0 || m == 0 {
		return n
	} else if m == 1 {
		return n - 1
	}
	l := make([]bool, n)
	for i := 0; i <= m%2; i++ {
		for j := 1; j < n; j += 2 {
			l[j] = true
		}
		for j := 2; j < n; j += 3 {
			l[j] = !l[j]
		}
	}
	l[n-1] = !l[n-1]
	var c int
	for i := 0; i < n; i++ {
		if !l[i] {
			c++
		}
	}
	return c
}
