package main

import "testing"

func TestMinCoins(t *testing.T) {
	for k, v := range map[uint]uint{
		11: 3, 20: 4, 12: 4, 13: 3} {
		if r := minCoins(k); r != v {
			t.Errorf("failed: minCoins %d is %d, got %d",
				k, v, r)
		}
	}
}

func BenchmarkMinCoins(b *testing.B) {
	for i := 0; i < b.N; i++ {
		minCoins(uint(i))
	}
}

var b []uint

func init() {
	b = []uint{0, 1, 2, 1, 2}
}

func minCoins(n uint) uint {
	return n/5 + b[n%5]
}
