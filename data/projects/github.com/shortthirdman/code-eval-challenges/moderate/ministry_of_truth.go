package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func contains(a, b string) (string, bool) {
	r := strings.Repeat("_", len(a))
	i := strings.Index(a, b)
	if i == -1 {
		return r, false
	}
	return r[:i] + b + r[i+len(b):], true
}

func ministry(p, q string) string {
	var (
		c int
		r []string
	)
	t, u := strings.Fields(p), strings.Fields(q)
	for i := range t {
		var v string
		if c < len(u) {
			var f bool
			if v, f = contains(t[i], u[c]); f {
				c++
			}
		} else {
			v = strings.Repeat("_", len(t[i]))
		}
		if len(v) > 0 {
			r = append(r, v)
		}
	}
	if c < len(u) {
		return "I cannot fix history"
	}
	return strings.Join(r, " ")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ";")
		fmt.Println(ministry(s[0], s[1]))
	}
}
