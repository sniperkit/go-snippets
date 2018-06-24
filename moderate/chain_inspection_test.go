package main

import (
	"strings"
	"testing"
)

func TestChain(t *testing.T) {
	for k, v := range map[string]bool{
		"4-2;BEGIN-3;3-4;2-END": true,
		"4-2;BEGIN-3;3-4;2-3":   false,
		"4-7;3-4;BEGIN-3":       false} {
		if r := chain(k); r != v {
			t.Errorf("failed: chain %s is %t, got %t",
				k, v, r)
		}
	}
}

func chain(q string) bool {
	var (
		been []string
		f    bool
	)
	s := strings.Split(q, ";")
	m, b := make(map[string]string, len(s)), "BEGIN"
	for _, i := range s {
		t := strings.Split(i, "-")
		if t[0] == t[1] {
			return false
		}
		m[t[0]] = t[1]
	}
	for _ = range s {
		if b, f = m[b]; !f {
			return false
		}
		for _, j := range been {
			if b == j {
				return false
			}
		}
		been = append(been, b)
	}
	if b == "END" {
		return true
	}
	return false
}
