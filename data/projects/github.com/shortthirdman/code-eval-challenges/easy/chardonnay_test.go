package main

import (
	"strings"
	"testing"
)

type tuple struct {
	a, b string
}

func TestContainsAll(t *testing.T) {
	for k, v := range map[tuple]bool{
		tuple{"ddccbbaa", "abcd"}:  true,
		tuple{"Chardonnay", "ann"}: true,
		tuple{"qwer", "ee"}:        false,
		tuple{"Cabernet", "ot"}:    false} {
		if r := containsAll(k.a, k.b); r != v {
			t.Errorf("failed: containsAll %s %s is %t, got %t",
				k.a, k.b, v, r)
		}
	}
}

func BenchmarkContainsAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		containsAll("Chardonnay", "ann")
		containsAll("Cabernet", "ot")
	}
}

func containsAll(a, b string) bool {
	for _, i := range b {
		ix := strings.IndexRune(a, i)
		if ix == -1 {
			return false
		}
		a = a[:ix] + a[ix+1:]
	}
	return true
}
