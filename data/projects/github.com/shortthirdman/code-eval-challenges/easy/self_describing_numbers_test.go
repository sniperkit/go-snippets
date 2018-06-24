package main

import (
	"fmt"
	"testing"
)

func TestSelf(t *testing.T) {
	for k, v := range map[string]bool{
		"2020":          true,
		"22":            false,
		"1210":          true,
		"123":           false,
		"3211000":       true,
		"6210001000":    true,
		"72100001000":   true,
		"821000001000":  true,
		"9210000001000": true} {
		if r := self(k); r != v {
			t.Errorf("failed: self %s is %t, got %t",
				k, v, r)
		}
	}
}

func BenchmarkSelf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		self(fmt.Sprint(i))
	}
}

func self(s string) bool {
	var n, r int64
	fmt.Sscan(s, &n)
	p := make([]int, len(s))
	for m := n; m > 0; m /= 10 {
		if v := int(m % 10); v < len(p) {
			p[v]++
		} else {
			return false
		}
	}
	for _, i := range p {
		r = r*10 + int64(i)
	}
	return n == r
}
