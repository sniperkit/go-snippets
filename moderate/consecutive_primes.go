package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

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

func main() {
	var n int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscan(scanner.Text(), &n)
		fmt.Println(consecutivePrimes(n))
	}
}
