package main

import "testing"

func TestZsum(t *testing.T) {
	for k, v := range map[int]int{
		0: 0, 7: 7, 1003: 4, 1234567890: 45} {
		if r := zsum(k); r != v {
			t.Errorf("failed: zsum %d = %d, got %d",
				k, v, r)
		}
	}
}

func BenchmarkZsum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		zsum(i)
	}
}

func zsum(a int) (r int) {
	for a > 0 {
		r, a = r+a%10, a/10
	}
	return r
}
