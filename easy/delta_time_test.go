package main

import (
	"fmt"
	"testing"
)

type tuple struct {
	h1, m1, s1, h2, m2, s2 int
}

func TestDelta(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{14, 1, 57, 12, 47, 11}:  "01:14:46",
		tuple{13, 9, 42, 22, 16, 15}:  "09:06:33",
		tuple{8, 8, 6, 8, 38, 28}:     "00:30:22",
		tuple{23, 35, 7, 2, 49, 59}:   "20:45:08",
		tuple{14, 31, 45, 14, 46, 56}: "00:15:11"} {
		if r := delta(k.h1, k.m1, k.s1, k.h2, k.m2, k.s2); r != v {
			t.Errorf("failed: delta %02d:%02d:%02d %02d:%02d:%02d is %s, got %s",
				k.h1, k.m1, k.s1, k.h2, k.m2, k.s2, v, r)
		}
	}
}

func BenchmarkMultiples(b *testing.B) {
	var h1, m1, s1, h2, m2, s2 int
	for i := 0; i < b.N; i++ {
		j := i
		s1, j = 10*j%10, j/10
		s2, j = 10*j%10, j/10
		m1, j = 10*j%10, j/10
		m2, j = 10*j%10, j/10
		h1, j = j%24, j/24
		h2 = j % 24
		delta(h1, m1, s1, h2, m2, s2)
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func delta(h1, m1, s1, h2, m2, s2 int) string {
	t := abs((h1-h2)*3600 + (m1-m2)*60 + s1 - s2)
	return fmt.Sprintf("%02d:%02d:%02d", t/3600, (t/60)%60, t%60)
}
