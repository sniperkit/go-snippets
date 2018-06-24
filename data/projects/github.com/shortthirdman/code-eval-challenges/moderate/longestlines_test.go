package main

import (
	"strings"
	"testing"
)

type tuple struct {
	p, q string
}

func TestLongest(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{"Hello World\nCodeEval\n\n", "Quick Fox"}:               "Hello World\nQuick Fox\nCodeEval\n",
		tuple{"Hello World\nQuick Fox\nCodeEval\nA", "San Francisco"}: "San Francisco\nHello World\nQuick Fox\nCodeEval"} {
		s := strings.Split(k.p, "\n")
		if r := strings.Join(longest(s, k.q), "\n"); r != v {
			t.Errorf("failed: longest %s, %s is %s, got %s",
				k.p, k.q, v, r)
		}
	}
}

func longest(p []string, q string) []string {
	for ix, i := range p {
		if len(q) > len(i) {
			p = append(p[:ix], append([]string{q}, p[ix:len(p)-1]...)...)
			break
		}
	}
	return p
}
