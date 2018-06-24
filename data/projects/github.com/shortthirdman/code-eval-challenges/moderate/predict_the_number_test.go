package main

import "testing"

func TestKernighan(t *testing.T) {
	for k, v := range map[uint]uint{
		0: 0, 1: 1, 2: 1, 13: 3, 42: 3, 127: 7, 8192: 1,
		8200: 2, 4194305: 2, 4294967222: 29} {
		if r := kernighan(k); r != v {
			t.Errorf("failed: kernighan %d is %d, got %d",
				k, v, r)
		}
	}
}

func BenchmarkKernighan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		kernighan(uint(i))
	}
}

func kernighan(a uint) (r uint) {
	for ; a > 0; a &= a - 1 {
		r++
	}
	return r
}
