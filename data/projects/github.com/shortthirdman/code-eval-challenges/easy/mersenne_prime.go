package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var primes = []uint{2, 3}

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

func mersenne(a uint) string {
	var r []string
	nextPrime := primeSeq()
	for i := nextPrime(); a > 1<<i-1; i = nextPrime() {
		r = append(r, fmt.Sprint(1<<i-1))
	}
	return strings.Join(r, ", ")
}

func main() {
	var a uint
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscan(scanner.Text(), &a)
		fmt.Println(mersenne(a))
	}
}
