package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func pow3(x int) int {
	y, ret := 3, 1
	for x > 0 {
		if x&1 > 0 {
			ret *= y
		}
		y *= y
		x >>= 1
	}
	return ret
}

func ugly(j int, n []int) bool {
	var (
		k  int
		s  int64
		c  = int64(n[0])
		cj = j
		p  = true
	)
	for k < len(n)-1 {
		ops := cj % 3
		cj /= 3
		if ops == 0 {
			c *= 10
		} else {
			if p {
				s += c
				if ops == 1 {
					p = false
				}
			} else {
				s -= c
				if ops == 2 {
					p = true
				}
			}
			c = 0
		}
		k++
		c += int64(n[k])
	}
	if p {
		s += c
	} else {
		s -= c
	}
	return s%2 == 0 || s%3 == 0 || s%5 == 0 || s%7 == 0
}

func uglyNumbers(q string) (r int) {
	t := make([]int, len(q))
	for ix, i := range q {
		t[ix] = int(i - '0')
	}
	p := pow3(len(t) - 1)
	for j := 0; j < p; j++ {
		if ugly(j, t) {
			r++
		}
	}
	return r
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(uglyNumbers(scanner.Text()))
	}
}
