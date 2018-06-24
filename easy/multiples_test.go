package main

import "testing"

type tuple struct {
	x, n int
}

func TestMultiples(t *testing.T) {
	for k, v := range map[tuple]int{
		tuple{13, 8}:               16,
		tuple{17, 16}:              32,
		tuple{2000000000000000, 1}: 2000000000000000} {
		if r := multiples(k.x, k.n); r != v {
			t.Errorf("failed: multiples %d %d is %d, got %d",
				k.x, k.n, v, r)
		}
	}
}

func BenchmarkMultiples(b *testing.B) {
	for i := 0; i < b.N; i++ {
		multiples(i/16+1, i%16+1)
	}
}

func BenchmarkMultiplesNaive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		multiplesNaive(i/16+1, i%16+1)
	}
}

func multiples(x, n int) (r int) {
	var c uint
	r = n
	for (r << 1) <= x {
		r <<= 1
		c++
	}
	for c > 0 {
		c--
		for r+(n<<c) <= x {
			r += n << c
		}
	}
	for r < x {
		r += n
	}
	return r
}

func multiplesNaive(x, n int) (r int) {
	r = n
	for r < x {
		r += n
	}
	return r
}
