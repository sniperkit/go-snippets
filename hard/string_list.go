package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func sortu(s string) (r string) {
	for _, i := range s {
		var done bool
		for jx, j := range r {
			if j == i {
				done = true
				break
			} else if j > i {
				if jx > 0 {
					r = r[0:jx] + string(i) + r[jx:len(r)]
				} else {
					r = string(i) + r
				}
				done = true
				break
			}
		}
		if !done {
			r = r + string(i)
		}
	}
	return r
}

func slist(s, t string, n int) string {
	if n == 0 {
		return t
	}
	var r []string
	for _, i := range s {
		r = append(r, slist(s, t+string(i), n-1))
	}
	return strings.Join(r, ",")
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
		s := strings.Split(scanner.Text(), ",")
		fmt.Sscan(s[0], &n)
		fmt.Println(slist(sortu(s[1]), "", n))
	}
}
