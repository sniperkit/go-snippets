package main

import "testing"

func TestDoubleTrouble(t *testing.T) {
	for k, v := range map[string]uint{
		"ABA*": 1,
		"BAA*": 0,
		"A*A*": 2} {
		if r := doubleTrouble(k); r != v {
			t.Errorf("failed: doubleTrouble %s is %d, got %d",
				k, v, r)
		}
	}
}

func doubleTrouble(q string) uint {
	var r uint = 1
	n := len(q) / 2
	for i := 0; i < n; i++ {
		if ((q[i] == 'A') && (q[i+n] == 'B')) ||
			((q[i] == 'B') && (q[i+n] == 'A')) {
			r = 0
			break
		} else if (q[i] == '*') && (q[i+n] == '*') {
			r *= 2
		}
	}
	return r
}
