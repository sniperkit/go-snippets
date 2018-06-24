package main

import "testing"

type tuple struct {
	n, a int
}

func TestXz(t *testing.T) {
	for k, v := range map[tuple]bool{
		tuple{1, 5}:     true,
		tuple{1, 7}:     false,
		tuple{7, 10922}: true} {
		if r := xz(k.n, k.a); r != v {
			t.Errorf("failed: xz %d %d is %t, got %t",
				k.n, k.a, v, r)
		}
	}
}

func BenchmarkXz(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xz(i%11, (i/11)%995+5)
	}
}

func xz(n, a int) bool {
	for a > 0 && n >= 0 {
		if a&1 == 0 {
			n--
		}
		a >>= 1
	}
	return n == 0
}
