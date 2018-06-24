package main

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

func TestBubble(t *testing.T) {
	for k, v := range map[string]string{
		"36 47 78 28 20 79 87 16 8 45 72 69 81 66 60 8 3 86 90 90 | 1":       "36 47 28 20 78 79 16 8 45 72 69 81 66 60 8 3 86 87 90 90",
		"40 69 52 42 24 16 66 | 2":                                           "40 42 24 16 52 66 69",
		"54 46 1 34 15 48 47 53 25 18 50 5 21 76 62 48 74 2 43 74 78 29 | 6": "1 15 25 18 34 5 21 46 47 48 48 2 43 50 53 29 54 62 74 74 76 78",
		"48 51 5 61 18 | 2":                                                  "5 48 18 51 61",
		"59 68 55 31 73 4 2 25 26 19 60 1 | 2":                               "55 31 59 4 2 25 26 19 60 1 68 73",
		"59 68 55 31 73 4 2 25 26 19 60 1 | 4321123456":                      "59 68 55 31 73 4 2 25 26 19 60 1",
		"42 | 1":        "42",
		"1 2 3 4 5 | 0": "1 2 3 4 5"} {
		if r := bubble(k); r != v {
			t.Errorf("failed: bubble %s is %s, got %s",
				k, v, r)
		}
	}
}

func bubble(q string) string {
	var n, k int
	t := strings.Split(q, " | ")
	fmt.Sscan(t[1], &n)
	u := strings.Fields(t[0])
	v := make([]int, len(u))
	for ix, i := range u {
		fmt.Sscan(i, &k)
		v[ix] = k
	}
	for i := 0; i < len(v)-n; i++ {
		sort.Ints(v[i : i+n+1])
	}
	for ix, i := range v {
		u[ix] = fmt.Sprint(i)
	}
	return strings.Join(u, " ")
}
