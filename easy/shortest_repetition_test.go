package main

import "testing"

func TestShortestRep(t *testing.T) {
	for k, v := range map[string]int{
		"abcabcabcabc":                 3,
		"bcbcbcbcbcbcbcbcbcbcbcbcbcbc": 2,
		"dddddddddddddddddddd":         1,
		"abcdefg":                      7,
		"b":                            1,
		"ll":                           1} {
		if r := shortestRep(k); r != v {
			t.Errorf("failed: shortestRep %s is %d, got %d",
				k, v, r)
		}
	}
}

func shortestRep(q string) int {
	r := 1
	for i := 1; i < len(q); i++ {
		if q[i] != q[i-r] {
			r = i + 1
		}
	}
	return r
}
