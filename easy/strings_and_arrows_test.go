package main

import "testing"

func TestStringsAndArrows(t *testing.T) {
	for k, v := range map[string]uint{
		"<--<<--<<":               2,
		"<<>>--><--<<--<<>>>--><": 4,
		"<-->>":                   0} {
		if r := stringsAndArrows(k); r != v {
			t.Errorf("failed: stringsAndArrows %s is %d, got %d",
				k, v, r)
		}
	}
}

func stringsAndArrows(s string) (r uint) {
	for i := 0; i < len(s)-4; i++ {
		if s[i:i+5] == ">>-->" || s[i:i+5] == "<--<<" {
			r++
		}
	}
	return r
}
