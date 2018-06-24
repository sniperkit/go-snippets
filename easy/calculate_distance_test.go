package main

import (
	"math"
	"testing"
)

func TestSq(t *testing.T) {
	for i := 0; i <= 80000; i++ {
		if r := int(math.Sqrt(float64(i * i))); r != i {
			t.Errorf("failed: math.Sqrt %d = %d, got %d", i*i, i, r)
		}
		if r := sqInt(i * i); r != i {
			t.Errorf("failed: sqInt %d = %d, got %d", i*i, i, r)
		}
	}
}

func BenchmarkSq(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = int(math.Sqrt(float64(i % 80000)))
	}
}

func BenchmarkSqInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sqInt(i % 80000)
	}
}

func sqInt(a int) (r int) {
	for r = 0; r*r < a; r++ {
	}
	return r
}
