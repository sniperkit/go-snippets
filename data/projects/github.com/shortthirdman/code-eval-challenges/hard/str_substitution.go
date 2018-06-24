package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type segment struct {
	s string
	r bool
}

func strsubs(p, q string) (r string) {
	t := strings.Split(q, ",")
	trans := []segment{segment{p, false}}
	for len(t) > 0 {
		for i := len(trans) - 1; i >= 0; i-- {
			if trans[i].r == false && strings.Contains(trans[i].s, t[0]) {
				u := strings.Split(trans[i].s, t[0])
				var v []segment
				for jx, j := range u {
					if len(j) > 0 {
						v = append(v, segment{j, false})
					}
					if jx != len(u)-1 {
						v = append(v, segment{t[1], true})
					}
				}
				if i > 0 {
					if i < len(trans)-1 {
						trans = append(trans[0:i], append(v, trans[i+1:len(trans)]...)...)
					} else {
						trans = append(trans[0:i], v...)
					}
				} else {
					if i < len(trans)-1 {
						trans = append(v, trans[i+1:len(trans)]...)
					} else {
						trans = v
					}
				}
			}
		}
		t = t[2:len(t)]
	}
	for _, i := range trans {
		r += i.s
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
		s := strings.Split(scanner.Text(), ";")
		fmt.Println(strsubs(s[0], s[1]))
	}
}
