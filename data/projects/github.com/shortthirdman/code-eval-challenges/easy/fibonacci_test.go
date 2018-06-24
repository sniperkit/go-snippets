package main

import "testing"

func TestFibo(t *testing.T) {
	for k, v := range map[uint64]uint64{
		0: 0, 1: 1, 2: 1, 3: 2, 4: 3, 5: 5, 6: 8, 7: 13,
		12: 144, 23: 28657, 64: 10610209857723} {
		if r := fibo(k); r != v {
			t.Errorf("failed: fibo %d = %d, got %d",
				k, v, r)
		}
	}
}

func BenchmarkFibo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fibo(uint64(i % 93))
	}
}

func fibo(n uint64) uint64 {
	if n < 2 {
		return n
	}
	var (
		r uint64 = 1
		p uint64 = 1
	)
	for ; n > 2; n-- {
		r, p = r+p, r
	}
	return r
}
