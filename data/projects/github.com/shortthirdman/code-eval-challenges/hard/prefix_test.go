package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestPrefix(t *testing.T) {
	for k, v := range map[string]int{
		"* + 2 3 4":      20,
		"42":             42,
		"* / * 2 10 2 7": 70} {
		if r := prefix(k); r != v {
			t.Errorf("failed: prefix %s is %d, got %d",
				k, v, r)
		}
	}
}

func pol(o byte, a, b float32) float32 {
	switch o {
	case '*':
		return a * b
	case '/':
		return a / b
	default: // '+'
		return a + b
	}
}

func prefix(q string) int {
	var r, v float32
	s := strings.Fields(q)
	n := len(s)
	fmt.Sscan(s[n/2], &r)
	for i := 1; i <= n/2; i++ {
		fmt.Sscan(s[n/2+i], &v)
		r = pol(s[n/2-i][0], r, v)
	}
	return int(r + 0.001)
}
