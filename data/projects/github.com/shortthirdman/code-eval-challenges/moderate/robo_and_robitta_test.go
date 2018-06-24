package main

import (
	"fmt"
	"testing"
)

func TestRobo(t *testing.T) {
	var a, b, x, y uint
	for k, v := range map[string]uint{
		"3x2 | 2 1":    5,
		"4x4 | 3 3":    14,
		"50x50 | 1 50": 1,
		"50x50 | 50 1": 99} {
		fmt.Sscanf(k, "%dx%d | %d %d", &a, &b, &x, &y)
		if r := robo(a, b, x, y); r != v {
			t.Errorf("failed: robo %s is %d, got %d",
				k, v, r)
		}
	}
}

func robo(a, b, x, y uint) uint {
	var r uint
	for b != y {
		r, a, b, x, y = r+a, b-1, a, b-y, x
	}
	return r + x
}
