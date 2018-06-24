package main

import (
	"fmt"
	"strings"
	"testing"
)

type tuple struct {
	p, q string
}

func TestPanacea(t *testing.T) {
	for k, v := range map[tuple]bool{
		tuple{"64 6e 78", "100101100 11110"}:          true,
		tuple{"5e 7d 59", "1101100 10010101 1100111"}: true,
		tuple{"93 75", "1000111 1011010 1100010"}:     false} {
		if r := panacea(k.p, k.q); r != v {
			t.Errorf("failed: panacea %s | %s is %t, got %t",
				k.p, k.q, v, r)
		}
	}
}

func panacea(p, q string) bool {
	var v, m int
	l, r := strings.Fields(p), strings.Fields(q)
	for _, i := range l {
		fmt.Sscanf(i, "%x", &v)
		m += v
	}
	for _, i := range r {
		fmt.Sscanf(i, "%b", &v)
		m -= v
	}
	return m <= 0
}
