package main

import (
	"math"
	"testing"
)

type tuple struct {
	x1, y1, x2, y2 float64
}

func TestDist(t *testing.T) {
	for k, v := range map[tuple]int{
		tuple{0, 0, 10000, 10000}: 14142135,
		tuple{3, 7, 1, 5}:         2828,
		tuple{1, 9, 8, 3}:         9219} {
		if r := int(1000 * dist(point{k.x1, k.y1}, point{k.x2, k.y2})); r != v {
			t.Errorf("failed: dist %f %f, %f %f is %f, got %f",
				k.x1, k.y1, k.x2, k.y2, float64(v)/1000, r)
		}
	}
}

func BenchmarkDist(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dist(point{float64(i%11) * 1000, float64((i/11)%11) * 1000}, point{float64((i/121)%11) * 1000, float64((i/1331)%11) * 1000})
	}
}

type point struct {
	x, y float64
}

func dist(a, b point) float64 {
	x, y := a.x-b.x, a.y-b.y
	return math.Sqrt(x*x + y*y)
}
