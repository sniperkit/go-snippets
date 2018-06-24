package main

import "testing"

type tuple struct {
	n, p, q uint
}

func TestXz(t *testing.T) {
	for k, v := range map[tuple]bool{
		tuple{86, 2, 3}:  true,
		tuple{125, 1, 2}: false,
		tuple{1, 1, 1}:   true,
		tuple{0, 4, 16}:  true} {
		if r := bitEqual(k.n, k.p, k.q); r != v {
			t.Errorf("failed: bitEqual %d %d %d is %t, got %t",
				k.n, k.p, k.q, v, r)
		}
	}
}

func BenchmarkBitEqual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bitEqual(uint(i/256), uint((i/16)%16), uint(i%16))
	}
}

func bitEqual(n, p, q uint) bool {
	return ((n & (1 << (p - 1))) == 0) == ((n & (1 << (q - 1))) == 0)
}
