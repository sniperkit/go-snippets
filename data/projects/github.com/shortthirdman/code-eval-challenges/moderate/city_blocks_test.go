package main

import (
	"fmt"
	"strings"
	"testing"
)

type tuple struct {
	s, t string
}

func TestCityBlocks(t *testing.T) {
	for k, v := range map[tuple]int{
		tuple{"0,2,4,8,10,13,14,18,22,23,24,33,40,42,44,47,49,53,55,63,66,81,87,91", "0,147,220"}: 24,
		tuple{"0,1,2,4", "0,1,3,4,5"}:                                                             6,
		tuple{"0,1,3,4,6", "0,1,2,4"}:                                                             5,
		tuple{"1,3", "1,2"}:                                                                       1,
		tuple{"1", "1,2"}:                                                                         0,
		tuple{"1,3", "2"}:                                                                         0} {
		if r := cityBlocks(k.s, k.t); r != v {
			t.Errorf("failed: cityBlocks %s %s is %d, got %d",
				k.s, k.t, v, r)
		}
	}
}

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
