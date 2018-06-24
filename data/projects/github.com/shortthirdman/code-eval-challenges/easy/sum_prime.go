package main

import "fmt"

var primes = []int{2, 3, 5, 7, 11, 13}

func isPrime(a int) bool {
	for b := 0; a%primes[b] > 0; b++ {
		if primes[b]*primes[b] > a {
			primes = append(primes, a)
			return true
		}
	}
	return false
}

func primeSeq() func() int {
	p0 := 0
	var i int
	return func() int {
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

func main() {
	var primesum int
	for i, nextPrime := 0, primeSeq(); i < 1000; i++ {
		primesum += nextPrime()
	}
	fmt.Println(primesum)
}
