package main

import (
	"fmt"
	"strings"
	"testing"
)

type tuple struct {
	p, q string
}

func TestSwapElements(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{"1 2 3 4 5 6 7 8 9", "0-8"}:         "9 2 3 4 5 6 7 8 1",
		tuple{"1 2 3 4 5 6 7 8 9 10", "0-1, 1-3"}: "2 4 3 1 5 6 7 8 9 10"} {
		if r := swapElements(k.p, k.q); r != v {
			t.Errorf("failed: swapElements %s : %s is %s, got %s",
				k.p, k.q, v, r)
		}
	}
}

func swapElements(p, q string) string {
	var a, b int
	r, s := strings.Fields(p), strings.Fields(q)
	for _, i := range s {
		fmt.Sscanf(i, "%d-%d,", &a, &b)
		r[a], r[b] = r[b], r[a]
	}
	return strings.Join(r, " ")
}
