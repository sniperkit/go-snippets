package main

import "testing"

func TestEven(t *testing.T) {
	for k, v := range map[int]int{
		701:  0,
		4123: 0,
		2936: 1,
		1:    0,
		4999: 0} {
		if r := even(k); r != v {
			t.Errorf("failed: even %d is %d, got %d",
				k, v, r)
		}
	}
}

func BenchmarkEven(b *testing.B) {
	for i := 0; i < b.N; i++ {
		even(i)
	}
}

func even(n int) int {
	return 1 - n&1
}
