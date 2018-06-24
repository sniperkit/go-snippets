package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

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

func findPath(n, p, b int) (ret int) {
	var r int
	if p >= n && (b&(1<<uint(p-n)) == 0) {
		ret = 1 + findPath(n, p-n, b|u[ms[p-n]])
	}
	if p < n*(n-1) && (b&(1<<uint(p+n)) == 0) {
		r = 1 + findPath(n, p+n, b|u[ms[p+n]])
		if r > ret {
			ret = r
		}
	}
	if p%n > 0 && (b&(1<<uint(p-1)) == 0) {
		r = 1 + findPath(n, p-1, b|u[ms[p-1]])
		if r > ret {
			ret = r
		}
	}
	if p%n < n-1 && (b&(1<<uint(p+1)) == 0) {
		r = 1 + findPath(n, p+1, b|u[ms[p+1]])
		if r > ret {
			ret = r
		}
	}
	return ret
}

var (
	u  map[uint8]int
	ms string
)

func main() {
	var n int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		m := len(scanner.Text())
		ms, n = scanner.Text(), sq(m)
		u = make(map[uint8]int)
		for i := 0; i < m; i++ {
			u[ms[i]] |= 1 << uint(i)
		}
		upperBound := len(u)
		var longPath int
		if upperBound < 3 || n == 2 || upperBound == m {
			longPath = upperBound
		} else if upperBound == m-1 {
			longPath = upperBound - n%2
		} else {
			longPath = 2
		}
		for i := 0; i < m-2 && longPath < upperBound; i++ {
			x := 1 + findPath(n, i, u[ms[i]])
			if x > longPath {
				longPath = x
			}
		}
		fmt.Printf("%d\n", longPath)
	}
}
