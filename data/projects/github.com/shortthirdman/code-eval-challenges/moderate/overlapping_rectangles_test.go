package main

import "testing"

type point struct {
	x, y int
}

type tuple struct {
	a, b, c, d point
}

func TestOverlapping(t *testing.T) {
	for k, v := range map[tuple]bool{
		tuple{point{-3, 3}, point{-1, 1}, point{1, -1}, point{3, -3}}: false,
		tuple{point{-3, 3}, point{-1, 1}, point{-2, 4}, point{2, 2}}:  true} {
		if r := overlapping(k.a, k.b, k.c, k.d); r != v {
			t.Errorf("failed: overlapping %d %d %d %d %d %d %d %d is %t, got %t",
				k.a.x, k.a.y, k.b.x, k.b.y, k.c.x, k.c.y, k.d.x, k.d.y, v, r)
		}
	}
}

func within(a, b, c point) bool {
	return c.x >= a.x && c.x <= b.x && c.y >= b.y && c.y <= a.y
}

func overlapping(a, b, c, d point) bool {
	for _, i := range []point{a, point{a.x, b.y}, point{b.x, a.y}, b} {
		if within(c, d, i) {
			return true
		}
	}
	for _, i := range []point{c, point{c.x, d.y}, point{d.x, c.y}, d} {
		if within(a, b, i) {
			return true
		}
	}
	return false
}
