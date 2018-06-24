package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestPrimeLess(t *testing.T) {
	for k, v := range map[uint]string{
		10:  "2,3,5,7",
		20:  "2,3,5,7,11,13,17,19",
		100: "2,3,5,7,11,13,17,19,23,29,31,37,41,43,47,53,59,61,67,71,73,79,83,89,97",
		42:  "2,3,5,7,11,13,17,19,23,29,31,37,41"} {
		if r := primeLess(k); r != v {
			t.Errorf("failed: primeLess %d is %s, got %s",
				k, v, r)
		}
	}
}

func BenchmarkPrimeLess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		primeLess(uint(i))
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

func primeLess(a uint) string {
	var r []string
	nextPrime := primeSeq()
	for i := nextPrime(); i < a; i = nextPrime() {
		r = append(r, fmt.Sprint(i))
	}
	return strings.Join(r, ",")
}
