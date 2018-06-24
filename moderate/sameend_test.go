package main

import (
	"strings"
	"testing"
)

type tuple struct {
	p, q string
}

func TestSameend(t *testing.T) {
	for k, v := range map[tuple]bool{
		tuple{"Hello World", "World"}:       true,
		tuple{"Hello CodeEval", "CodeEval"}: true,
		tuple{"San Francisco", "San Jose"}:  false,
		tuple{"spaceman", "man "}:           false} {
		if r := sameend(k.p, k.q); r != v {
			t.Errorf("failed: sameend %s, %s is %t, got %t",
				k.p, k.q, v, r)
		}
	}
}

func sameend(p, q string) bool {
	return strings.HasSuffix(p, q)
}
