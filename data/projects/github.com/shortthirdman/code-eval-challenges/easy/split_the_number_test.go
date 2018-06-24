package main

import (
	"fmt"
	"strings"
	"testing"
)

type tuple struct {
	s, t string
}

func TestMask(t *testing.T) {
	for k, v := range map[tuple]int{
		tuple{"3413289830", "a-bcdefghij"}: -413289827,
		tuple{"776", "a+bc"}:               83,
		tuple{"12345", "a+bcde"}:           2346,
		tuple{"1232", "ab+cd"}:             44,
		tuple{"90602", "a+bcde"}:           611} {
		if r := split(k.s, k.t); r != v {
			t.Errorf("failed: split %s %s is %d, got %d",
				k.s, k.t, v, r)
		}
	}
}

func split(s, t string) int {
	var n, m, p int
	var neg bool
	if strings.Contains(t, "+") {
		p = strings.Index(t, "+")
	} else {
		p = strings.Index(t, "-")
		neg = true
	}
	fmt.Sscan(s[0:p], &n)
	fmt.Sscanf(s[p:len(s)], "%d", &m) // may have leading 0
	if neg {
		return n - m
	}
	return n + m
}
