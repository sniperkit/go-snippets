package main

import (
	"testing"
)

type tuple struct {
	e, s int
}

func TestFloors(t *testing.T) {
	for k, v := range map[tuple]int{
		tuple{2, 100}:    14,
		tuple{90, 10000}: 14} {
		var r int
		for floors(k.e, r) < k.s {
			r++
		}
		if r != v {
			t.Errorf("failed: floors %d %d is %d, got %d",
				k.e, k.s, v, r)
		}
	}
}

var h map[int]int

func init() {
	h = make(map[int]int)
}

func floors(e, s int) int {
	switch {
	case e == 0 || s == 0:
		return 0
	case e == 1:
		return s
	case s == 1:
		return 1
	case e > s:
		return floors(s, s)
	}
	if hf, f := h[e<<16+s]; f {
		return hf
	}
	h[e<<16+s] = floors(e-1, s-1) + floors(e, s-1) + 1
	return h[e<<16+s]
}
