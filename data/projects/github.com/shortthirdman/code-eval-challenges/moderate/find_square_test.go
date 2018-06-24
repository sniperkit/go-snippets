package main

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

type tuple struct {
	x1, y1, x2, y2 int
}

func TestDist2(t *testing.T) {
	for k, v := range map[tuple]int{
		tuple{0, 0, 10, 10}: 200,
		tuple{3, 7, 1, 5}:   8,
		tuple{1, 9, 8, 3}:   85} {
		if r := dist2(point{k.x1, k.y1}, point{k.x2, k.y2}); r != v {
			t.Errorf("failed: dist2 (%d, %d) (%d, %d) is %d, got %d",
				k.x1, k.y1, k.x2, k.y2, v, r)
		}
	}
}

func TestFindSquare(t *testing.T) {
	for k, v := range map[string]bool{
		"(1,6), (6,7), (2,7), (9,1)": false,
		"(4,1), (3,4), (0,5), (1,2)": false,
		"(4,6), (5,5), (5,6), (4,5)": true} {
		if r := findSquare(k); r != v {
			t.Errorf("failed: findSquare %s is %t, got %t",
				k, v, r)
		}
	}
}

func BenchmarkDist2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dist2(point{i % 11, (i / 11) % 11}, point{(i / 121) % 11, (i / 1331) % 11})
	}
}

type point struct {
	x, y int
}

func dist2(a, b point) int {
	x, y := a.x-b.x, a.y-b.y
	return x*x + y*y
}

func findSquare(q string) bool {
	var (
		x, y int
		d    []int
	)
	t := strings.Split(q, ", ")
	u := make([]point, 4)
	for ix, i := range t {
		fmt.Sscanf(i, "(%d,%d)", &x, &y)
		u[ix] = point{x, y}
	}
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 4; j++ {
			d = append(d, dist2(u[i], u[j]))
		}
	}
	sort.Ints(d)
	return d[0] == d[1] && d[0] == d[2] && d[0] == d[3] &&
		d[4] == d[5] && 2*d[0] == d[4]
}
