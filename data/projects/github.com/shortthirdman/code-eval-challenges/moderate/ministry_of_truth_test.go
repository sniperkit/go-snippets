package main

import (
	"strings"
	"testing"
)

type tuple struct {
	p, q string
}

func TestMinistry(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{"Higher meaning", "Hi mean"}:           "Hi____ mean___",
		tuple{"this is impossible", "im possible"}:   "I cannot fix history",
		tuple{"twenty   two minutes", "two minutes"}: "______ two minutes",
		tuple{"Higher meaning", "e"}:                 "____e_ _______"} {
		if r := ministry(k.p, k.q); r != v {
			t.Errorf("failed: ministry %s;%s is %s, got %s",
				k.p, k.q, v, r)
		}
	}
}

func contains(a, b string) (string, bool) {
	r := strings.Repeat("_", len(a))
	i := strings.Index(a, b)
	if i == -1 {
		return r, false
	}
	return r[:i] + b + r[i+len(b):], true
}

func ministry(p, q string) string {
	var (
		c int
		r []string
	)
	t, u := strings.Fields(p), strings.Fields(q)
	for i := range t {
		var v string
		if c < len(u) {
			var f bool
			if v, f = contains(t[i], u[c]); f {
				c++
			}
		} else {
			v = strings.Repeat("_", len(t[i]))
		}
		if len(v) > 0 {
			r = append(r, v)
		}
	}
	if c < len(u) {
		return "I cannot fix history"
	}
	return strings.Join(r, " ")
}
