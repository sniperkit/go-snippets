package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func maxi(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxis(a []int) (r int) {
	r = a[0]
	for i := 1; i < len(a); i++ {
		if a[i] > r {
			r = a[i]
		}
	}
	return r
}

var u []int

func pass(q string) {
	s := strings.Fields(q)
	t := make([]int, len(s))
	for ix, i := range s {
		fmt.Sscan(i, &t[ix])
	}
	if len(u) > 0 {
		t[0] += u[0]
		if len(u) > 1 {
			t[len(u)] += u[len(u)-1]
		}
	}
	for i := 1; i < len(t)-1; i++ {
		t[i] += maxi(u[i-1], u[i])
	}
	u = make([]int, len(t))
	copy(u, t)
	return
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		pass(scanner.Text())
	}
	fmt.Println(maxis(u))
}
