package main

import "testing"

type tuple struct {
	p, q string
}

func TestWithoutRepetitions(t *testing.T) {
	for k, v := range map[tuple]tuple{
		tuple{"a", "b"}:  tuple{"b", "b"},
		tuple{"a", "a"}:  tuple{"a", ""},
		tuple{"\n", " "}: tuple{" ", " "}} {
		if rp, rq := withoutRepetitions(k.p, k.q); rp != v.p || rq != v.q {
			t.Errorf("failed: withoutRepetitions %s %s is %s %s, got %s %s",
				k.p, k.q, v.p, v.q, rp, rq)
		}
	}
}

func withoutRepetitions(p, q string) (string, string) {
	if p != q {
		return q, q
	}
	return p, ""
}
