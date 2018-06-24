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

func isPalindrome(a int) bool {
	if a%10 == 0 {
		return a == 0
	}
	var rev int
	for ; a > rev; a /= 10 {
		rev = 10*rev + a%10
		if a == rev {
			return true
		}
	}
	return rev == a
}

func main() {
	var maxprime int
	nextPrime := primeSeq()
	for i := nextPrime(); i < 1000; i = nextPrime() {
		if isPalindrome(i) {
			maxprime = i
		}
	}
	fmt.Println(maxprime)
}
