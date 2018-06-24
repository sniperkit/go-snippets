package main

import (
	"fmt"
	"testing"
)

func TestRev(t *testing.T) {
	for k, v := range map[int]int{
		0: 0, 1: 1, 10: 1, 15: 51, 29: 92, 110: 11,
		1001: 1001, 9990: 999, 9999: 9999} {
		if r := rev(k); r != v {
			t.Errorf("failed: rev %d is %d, got %d",
				k, v, r)
		}
	}
}

func TestRevAdd(t *testing.T) {
	for k, v := range map[int]string{
		195:  "4 9339",
		155:  "3 4444",
		2139: "2 26862",
		196:  "not found",
		0:    "0 0",
		1:    "0 1",
		2:    "0 2",
		3:    "0 3",
		4:    "0 4",
		5:    "0 5",
		9998: "6 8836388",
		9999: "0 9999"} {
		if r := revAdd(k); r != v {
			t.Errorf("failed: revAdd %d is %s, got %s",
				k, v, r)
		}
	}
}

func BenchmarkRev(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rev(i % 10000)
	}
}

func BenchmarkRevAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		revAdd(i % 10000)
	}
}

func rev(n int) (r int) {
	for ; n > 0; n /= 10 {
		r = 10*r + n%10
	}
	return r
}

func revAdd(n int) string {
	var c, r int
	for c < 100 {
		r = rev(n)
		if r == n {
			return fmt.Sprint(c, n)
		}
		n, c = n+r, c+1
	}
	return "not found"
}
