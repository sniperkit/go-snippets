package main

import "testing"

type tuple struct {
	x1, y1, x2, y2 int
}

func TestComparePoints(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{0, 0, 1, 5}:     "NE",
		tuple{12, 13, 12, 13}: "here",
		tuple{0, 1, 0, 5}:     "N",
		tuple{0, 0, 1, -5}:    "SE"} {
		if r := comparePoints(k.x1, k.y1, k.x2, k.y2); r != v {
			t.Errorf("failed: comparePoints %d %d %d %d is %s, got %s",
				k.x1, k.y1, k.x2, k.y2, v, r)
		}
	}
}

func BenchmarkComparePoints(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x1 := i%21 - 10
		y1 := (i/21)%21 - 10
		x2 := (i/441)%21 - 10
		y2 := (i/9261)%21 - 10
		comparePoints(x1, y1, x2, y2)
	}
}

func comparePoints(x1, y1, x2, y2 int) string {
	switch {
	case x1 == x2:
		if y1 == y2 {
			return "here"
		} else if y1 < y2 {
			return "N"
		} else {
			return "S"
		}
	case y1 == y2:
		if x1 < x2 {
			return "E"
		} else {
			return "W"
		}
	case x1 < x2:
		if y1 < y2 {
			return "NE"
		} else {
			return "SE"
		}
	default:
		if y1 < y2 {
			return "NW"
		} else {
			return "SW"
		}
	}
}
