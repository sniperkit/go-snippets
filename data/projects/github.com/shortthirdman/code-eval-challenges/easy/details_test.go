package main

import (
	"strings"
	"testing"
)

func TestDetails(t *testing.T) {
	for k, v := range map[string]int{
		"XX.YY,XXX.Y,X..YY,XX..Y":         1,
		"XX...YY,X....YY,XX..YYY,X..YYYY": 2,
		"XXYY,X..Y,XX.Y":                  0,
		"X........Y,X........Y":           8,
		"X.Y..Y,X..Y":                     1,
		"XY.Y,XX.Y":                       0} {
		if r := details(k); r != v {
			t.Errorf("failed: details %s is %d, got %d",
				k, v, r)
		}
	}
}

func details(s string) int {
	t := strings.Split(s, ",")
	m := len(t[0])
	for _, i := range t {
		if n := strings.Index(i, "Y") - strings.LastIndex(i, "X"); n < m {
			m = n
		}
	}
	return m - 1
}
