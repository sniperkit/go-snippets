package main

import "testing"

type tuple struct {
	a, b string
}

func TestSubs(t *testing.T) {
	for k, v := range map[tuple]uint{
		tuple{"babgbag", "bag"}:    5,
		tuple{"rabbbit", "rabbit"}: 3} {
		if r := subs(k.a, k.b); r != v {
			t.Errorf("failed: subs %s %s is %d, got %d",
				k.a, k.b, v, r)
		}
	}
}

func subs(s, t string) uint {
	if len(t) == 0 {
		return 1
	}
	for len(s) > 0 {
		if s[0] == t[0] {
			return subs(s[1:len(s)], t[1:len(t)]) + subs(s[1:len(s)], t)
		}
		s = s[1:len(s)]
	}
	return 0
}
