package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestCombSort(t *testing.T) {
	for k, v := range map[string]uint{
		"3 1 2":     2,
		"5 4 3 2 1": 3} {
		if r := combSort(k); r != v {
			t.Errorf("failed: combSort %s is %d, got %d",
				k, v, r)
		}
	}
}

const shrink = 1.25

func combSort(q string) uint {
	var (
		n, l uint
		k, g int
		c    bool
	)
	u := strings.Fields(q)
	g, c = len(u), false
	v := make([]int, len(u))
	for ix, i := range u {
		fmt.Sscan(i, &k)
		v[ix] = k
	}
	for g > 1 || c {
		n, c = n+1, false
		if g > 1 {
			g = int(float32(g) / shrink)
		}
		for j := 0; j < len(u)-g; j++ {
			if u[j] > u[j+g] {
				u[j], u[j+g], c, l = u[j+g], u[j], true, n
			}
		}
	}
	return l
}
