package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func notSoClever(q string, n int) string {
	var j, k int
	u := strings.Fields(q)
	v := make([]int, len(u))
	for ix, i := range u {
		fmt.Sscan(i, &k)
		v[ix] = k
	}
	l := 1
	for i := 0; i < n; i++ {
		for j, l = l, 0; j < len(v); j++ {
			if v[j-1] <= v[j] {
				continue
			}
			v[j-1], v[j] = v[j], v[j-1]
			if j > 1 {
				l = j - 1
			} else {
				l = 2
			}
			break
		}
		if l == 0 {
			break
		}
	}
	for ix, i := range v {
		u[ix] = fmt.Sprint(i)
	}
	return strings.Join(u, " ")
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
		t := strings.Split(scanner.Text(), " | ")
		fmt.Sscan(t[1], &n)
		fmt.Println(notSoClever(t[0], n))
	}
}
