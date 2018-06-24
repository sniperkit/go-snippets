package main

import (
	"strings"
	"testing"
)

type tuple struct {
	p, q string
}

func TestRightmost(t *testing.T) {
	for k, v := range map[tuple]int{
		tuple{"Hello World", "r"}:                       8,
		tuple{"Hello CodeEval", "E"}:                    10,
		tuple{"Yesofcourse", "X"}:                       -1,
		tuple{"in nova fert animus", "n"}:               14,
		tuple{"Willst Du einen Schneemann bauen?", "W"}: 0} {
		if r := rightmost(k.p, k.q); r != v {
			t.Errorf("failed: rightmost %s, %s is %d, got %d",
				k.p, k.q, v, r)
		}
	}
}

func rightmost(p, q string) int {
	return strings.LastIndex(p, q)
}
