package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"testing"
)

type tuple struct {
	a, b, c, d int
	q          string
}

func TestPile(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{4, 3, 3, -3, "(1 [10,9,4] [9,4,2])"}:                                          "1",
		tuple{-1, -5, 5, -2, "(1 [4,7,8] [2,9,0]);(2 [0,7,1] [5,9,8])"}:                     "1,2",
		tuple{-4, -5, -5, -3, "(1 [4,8,6] [0,9,2]);(2 [8,-1,3] [0,5,4])"}:                   "-",
		tuple{5, 5, -5, -5, "(11 [0,0,0] [1,1,1]);(3 [0,0,0] [1,1,1]);(2 [0,0,0] [1,1,1])"}: "2,3,11"} {
		if r := pile(k.a, k.b, k.c, k.d, k.q); r != v {
			t.Errorf("failed: pile %d %d %d %d %s is %s, got %s",
				k.a, k.b, k.c, k.d, k.q)
		}
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func diagonal(c, h []int) bool {
	a, b, p, q := float64(h[0]), float64(h[1]), float64(c[0]), float64(c[1])
	return (p > a) && (b >= (2*p*q*a+(p*p-q*q)*math.Sqrt(p*p+q*q-a*a))/(p*p+q*q))
}

func pile(a, b, c, d int, q string) string {
	var e, f, j int
	g, h := make([]int, 2), make([]int, 3)
	g[0], g[1] = abs(a-c), abs(b-d)
	sort.Ints(g)
	var r []int
	for _, i := range strings.Split(q, ";") {
		fmt.Sscanf(i, "(%d [%d,%d,%d] [%d,%d,%d])", &j, &a, &b, &c, &d, &e, &f)
		h[0], h[1], h[2] = abs(a-d), abs(b-e), abs(c-f)
		sort.Ints(h)
		if h[0] <= g[0] && h[1] <= g[1] { // || diagonal(h, g)
			r = append(r, j)
		}
	}
	if len(r) == 0 {
		return "-"
	}
	sort.Ints(r)
	var r2 []string
	for _, i := range r {
		r2 = append(r2, fmt.Sprint(i))
	}
	return strings.Join(r2, ",")
}
