package main

import "testing"

type tuple struct {
	a, b uint
}

func TestCountPrimes(t *testing.T) {
	for k, v := range map[tuple]uint{
		tuple{2, 10}:  4,
		tuple{24, 25}: 0,
		tuple{20, 30}: 2,
		tuple{31, 31}: 1,
		tuple{31, 32}: 1,
		tuple{30, 31}: 1,
		tuple{29, 31}: 2,
		tuple{29, 29}: 1} {
		if r := countPrimes(k.a, k.b); r != v {
			t.Errorf("failed: countPrimes %d,%d is %d, got %d",
				k.a, k.b, r, v)
		}
	}
}

var primes = []uint{2, 3, 5, 7, 11, 13}

func isPrime(a uint) bool {
	for b := 0; a%primes[b] > 0; b++ {
		if primes[b]*primes[b] > a {
			primes = append(primes, a)
			return true
		}
	}
	return false
}

func primeSeq() func() uint {
	p0 := 0
	var i uint
	return func() uint {
		if p0 < len(primes) {
			i = primes[p0]
			p0++
			return i
		}
		for i += 2; !isPrime(i); i += 2 {
		}
		p0++
		return i
	}
}

func countPrimes(a, b uint) (r uint) {
	nextPrime := primeSeq()
	for i := nextPrime(); i <= b; i = nextPrime() {
		if i >= a {
			r++
		}
	}
	return r
}
