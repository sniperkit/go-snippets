package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func kernighan(a int) (r int) {
	for ; a > 0; a &= a - 1 {
		r++
	}
	return r
}

func numCar(m []int, n int) (r int) {
	var i uint
	for i = 0; i < uint(n); i++ {
		r += kernighan(m[i] & ((1 << uint(n)) - 1))
	}
	return r
}

func isNum(m []int, n, x, y, nc int) bool {
	var (
		r int
		i uint
	)
	for i = uint(x); i < uint(x+n); i++ {
		r += kernighan(m[i] & (((1 << uint(n)) - 1) << uint(y)))
		if r > nc {
			return false
		}
	}
	return r == nc
}

func blackOrWhite(q string) string {
	var rm, rn int
	s := strings.Split(q, "|")
	m := make([]int, len(s))
	for ix, i := range s {
		fmt.Sscanf(strings.TrimSpace(i), "%b", &m[ix])
	}

	for rm = 1; rm <= len(s); rm++ {
		f := false
		rn = numCar(m, rm)
		for i := 0; i < len(s)-rm+1; i++ {
			for j := 0; j < len(s)-rm+1; j++ {
				if !isNum(m, rm, i, j, rn) {
					f = true
					break
				}
			}
			if f {
				break
			}
		}
		if !f {
			break
		}
	}
	return fmt.Sprintf("%dx%d, %d", rm, rm, rn)
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(blackOrWhite(scanner.Text()))
	}
}
