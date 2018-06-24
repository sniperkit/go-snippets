package main

import "testing"

func TestSq(t *testing.T) {
	for i := 2; i <= 6; i++ {
		if r := sq(i * i); r != i {
			t.Errorf("failed: sq %d = %d, got %d", i*i, i, r)
		}
	}
}

func BenchmarkSq(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sq(i%5 + 2)
	}
}

func BenchmarkSqGeneric(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sqGeneric(i%5 + 2)
	}
}

func sq(a int) (ret int) {
	switch a {
	case 4:
		ret = 2
	case 9:
		ret = 3
	case 16:
		ret = 4
	case 25:
		ret = 5
	case 36:
		ret = 6
	}
	return ret
}

func sqGeneric(a int) (r int) {
	for r = 1; r*r < a; r++ {
	}
	return r
}
