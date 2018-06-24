package main

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

func TestDistinctTriangles(t *testing.T) {
	for k, v := range map[string]uint{
		"0 2,0 1,1 2,1 3,2 3": 2,
		"1 3,1 8,3 8":         1,
		"5 6,5 7,6 7":         1} {
		if r := distinctTriangles(k); r != v {
			t.Errorf("failed: distinctTriangles %s is %d, got %d",
				k, v, r)
		}
	}
}

func contains(a []int, b int) bool {
	ix := sort.SearchInts(a, b)
	return ix < len(a) && a[ix] == b
}

func distinctTriangles(q string) (r uint) {
	g := make(map[int][]int)
	t := strings.Split(q, ",")
	p := make([]int, 2)
	for _, i := range t {
		fmt.Sscanf(i, "%d %d", &p[0], &p[1])
		sort.Ints(p)
		g[p[0]] = append(g[p[0]], p[1])
	}
	for k := range g {
		sort.Ints(g[k])
	}
	for _, v := range g {
		for i := 0; i < len(v)-1; i++ {
			for j := i + 1; j < len(v); j++ {
				if contains(g[v[i]], v[j]) {
					r++
				}
			}
		}
	}
	return r
}
