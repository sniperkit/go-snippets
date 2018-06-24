package main

import (
	"strings"
	"testing"
)

type tuple struct {
	p, q string
}

func TestSubstring(t *testing.T) {
	for k, v := range map[tuple]bool{
		tuple{"Hello", "ell"}:             true,
		tuple{"This is good", " is "}:     true,
		tuple{"CodeEval", "C*Eval"}:       true,
		tuple{"test*ing", "test\\*ing"}:   true,
		tuple{"test\\ing", "test\\\\ing"}: false,
		tuple{"test*ing", "t*s*n*g"}:      true,
		tuple{"test*ing", "t*t*s*g"}:      false,
		tuple{"Old", "Young"}:             false} {
		if r := substring(k.p, k.q); r != v {
			t.Errorf("failed: substring %s,%s is %t, got %t",
				k.p, k.q, v, r)
		}
	}
}

func subcheck(s0, s1 string, q, r int) bool {
	for q < len(s0) && r < len(s1) {
		if s1[r] == '*' {
			var m bool
			for i := q; i < len(s0); i++ {
				if subcheck(s0, s1, i, r+1) {
					m = true
					break
				}
			}
			return m
		} else if s1[r] == '\\' && r < len(s1)+1 && s1[r+1] == '*' {
			if s0[q] != '*' {
				return false
			}
			q, r = q+1, r+2
		}
		if s0[q] != s1[r] {
			return false
		}
		q, r = q+1, r+1
	}
	return r >= len(s1)
}

func substring(p, q string) bool {
	q = strings.TrimLeft(q, "*") // stars without item. ignore.
	for len(q) > 1 && q[len(q)-1] == '*' && q[len(q)-2] != '\\' {
		q = q[0 : len(q)-1]
	} // redundant terminal stars.
	for i := 0; i < len(p); i++ {
		if subcheck(p, q, i, 0) {
			return true
		}
	}
	return false
}
