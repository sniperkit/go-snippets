package main

import (
	"strings"
	"testing"
)

type tuple struct {
	p, q string
}

func TestStringRotation(t *testing.T) {
	for k, v := range map[tuple]bool{
		tuple{"Hello", "lloHe"}:       true,
		tuple{"Hello", "Hello"}:       true,
		tuple{"Basefont", "tBasefon"}: true,
		tuple{"asdfasdf", "fdsa"}:     false} {
		if r := stringRotation(k.p, k.q); r != v {
			t.Errorf("failed: stringRotation %s %s is %t, got %t",
				k.p, k.q, v, r)
		}
	}
}

func stringRotation(p, q string) bool {
	var n int
	for i := strings.Index(q, string(p[0])); i >= 0; i += n {
		if strings.HasSuffix(q, string(p[0:len(p)-i])) && (i == 0 || strings.HasPrefix(q, string(p[len(p)-i:len(p)]))) {
			return true
		}
		n = 1 + strings.Index(q[i+1:len(q)], string(p[0]))
		if n == 0 {
			break
		}
	}
	return false
}
