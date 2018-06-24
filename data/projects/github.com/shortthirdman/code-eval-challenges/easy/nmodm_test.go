package main

import "testing"

type tuple struct {
	n, m int
}

func TestMod(t *testing.T) {
	for k, v := range map[tuple]int{
		tuple{20, 6}: 2,
		tuple{2, 3}:  2,
		tuple{1, 1}:  0} {
		if r := mod(k.n, k.m); r != v {
			t.Errorf("failed: mod %d %d is %d, got %d",
				k.n, k.m, v, r)
		}
	}
}

func BenchmarkMultiples(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mod(i+1, 10-i%10)
	}
}

func mod(n, m int) int {
	return n - (n/m)*m
}
