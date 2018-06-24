package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func cityBlocks(s, t string) int {
	p, q := strings.Split(s, ","), strings.Split(t, ",")
	if len(p) < 2 || len(q) < 2 {
		return 0
	}
	var intersec int
	st, av := make([]float32, len(p)), make([]float32, len(q))
	for ix, i := range p {
		fmt.Sscan(i, &st[ix])
	}
	for ix, i := range q {
		fmt.Sscan(i, &av[ix])
	}
	if st[0] != 0 {
		st0 := st[0]
		for ix, i := range st {
			st[ix] = i - st0
		}
	}
	if av[0] != 0 {
		av0 := av[0]
		for ix, i := range av {
			av[ix] = i - av0
		}
	}
	stl, avl := st[len(st)-1], av[len(av)-1]
	for ix, i := range av {
		av[ix] = i * stl / avl
	}
	for i, j := 0, 0; i < len(st) && j < len(av); {
		if st[i] == av[j] {
			intersec++
			i++
			j++
		} else if st[i] < av[j] {
			i++
		} else if st[i] > av[j] {
			j++
		}
	}
	return len(st) + len(av) - intersec - 1
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ") (")
		fmt.Println(cityBlocks(s[0][1:len(s[0])], s[1][0:len(s[1])-1]))
	}
}
