package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestM2last(t *testing.T) {
	for k, v := range map[string]string{
		"a b c d 4":     "a",
		"a b c d 5":     "",
		"e f g h 2":     "g",
		"12 23 34 45 1": "45",
		"0 1 2":         "0"} {
		if r := m2last(k); r != v {
			t.Errorf("failed: m2last %s is %s, got %s",
				k, v, r)
		}
	}
}

func m2last(q string) string {
	t := strings.Fields(q)
	var m uint
	n := uint(len(t))
	fmt.Sscan(t[n-1], &m)
	if m < n {
		return t[n-1-m]
	}
	return ""
}
