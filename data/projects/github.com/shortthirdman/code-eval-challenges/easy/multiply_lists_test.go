package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestMultiplyLists(t *testing.T) {
	for k, v := range map[string]string{
		"9 0 6 | 15 14 9": "135 0 54",
		"5 | 8":           "40",
		"13 4 15 1 15 5 | 1 4 15 14 8 2": "13 16 225 14 120 10"} {
		if r := multiplyLists(k); r != v {
			t.Errorf("failed: multiplyLists %s is %s, got %s",
				k, v, r)
		}
	}
}

func multiplyLists(q string) string {
	var u, v int
	s := strings.Split(q, " | ")
	x, y := strings.Fields(s[0]), strings.Fields(s[1])
	r := make([]string, len(x))
	for i := range x {
		fmt.Sscan(x[i], &u)
		fmt.Sscan(y[i], &v)
		r[i] = fmt.Sprint(u * v)
	}
	return strings.Join(r, " ")
}
