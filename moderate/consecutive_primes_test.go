package main

import (
	"sort"
	"testing"
)

func TestConsecutivePrimes(t *testing.T) {
	for k, v := range map[int]uint{
		3: 0, 5: 0, 7: 0, 9: 0, 11: 0, 13: 0, 15: 0, 17: 0,
		2: 1, 4: 2, 6: 2, 8: 4, 10: 96, 12: 1024, 14: 2880,
		16: 81024, 18: 770144} {
		if r := consecutivePrimes(k); r != v {
			t.Errorf("failed: consecutivePrimes %d is %d, got %d",
				k, v, r)
		}
	}
}

func BenchmarkConsecutivePrimes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		consecutivePrimes(i%17 + 2)
	}
}

var (
	primes = []int{2, 3, 5, 7, 11, 13}
	maxP   = primes[len(primes)-1]
)

func isPrime(a int) bool {
	if a <= primes[len(primes)-1] {
		return contains(primes, a)
	}
	for b := 0; a%primes[b] > 0; b++ {
		if primes[b]*primes[b] > a {
			return true
		}
	}
	return false
}

func contains(a []int, b int) bool {
	ix := sort.SearchInts(a, b)
	return ix < len(a) && a[ix] == b
}

func consecutive(g [][]uint, d int, c uint) (r uint) {
	if d == 0 && g[c][0] == 1 {
		return 1
	}
	for _, i := range g[c] {
		if d&(1<<i) > 0 {
			r += consecutive(g, d-(1<<i), i)
		}
	}
	return r
}

func consecutivePrimes(a int) uint {
	if a&1 > 0 {
		return 0
	}
	for maxP < 2*a-1 {
		maxP += 2
		if isPrime(maxP) {
			primes = append(primes, maxP)
		}
	}
	g := make([][]uint, a+1)
	for i := 1; i <= a; i++ {
		for j := 1; j <= a; j++ {
			if i != j && isPrime(i+j) {
				g[i] = append(g[i], uint(j))
			}
		}
	}
	d := (1 << uint(a+1)) - 4
	return consecutive(g, d, 1)
}
