package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestLowestUnique(t *testing.T) {
	for k, v := range map[string]uint{
		"3 3 9 1 6 5 8 1 5 3":   5,
		"9 2 9 9 1 8 8 8 2 1 1": 0} {
		if r := lowestUnique(k); r != v {
			t.Errorf("failed: lowestUnique %s is %d, got %d",
				k, v, r)
		}
	}
}

func BenchmarkLowestUnique(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var s []string
		for j := i + 1; j > 0; j /= 9 {
			s = append(s, fmt.Sprint(j%9+1))
		}
		lowestUnique(strings.Join(s, " "))
	}
}

func lowestUnique(s string) uint {
	var k uint
	m := make([]int, 9)
	for ix, i := range strings.Fields(s) {
		fmt.Sscan(i, &k)
		if m[k-1] == 0 {
			m[k-1] = ix + 1
		} else {
			m[k-1] = -1
		}
	}
	for i := uint(0); i < 9; i++ {
		if m[i] > 0 {
			return i
		}
	}
	return 0
}
