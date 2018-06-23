package devilc

import (
	"math"
)

func devilCalcFunc(aa, bb float64, t float64) Point {

	sin, cos := math.Sincos(t)

	var (
		s2 = sin * sin
		c2 = cos * cos
	)

	d := (aa*s2 - bb*c2) / (s2 - c2)
	if d < 0 {
		return Point{}
	}

	sqrt := math.Sqrt(d)

	return Point{
		X: cos * sqrt,
		Y: sin * sqrt,
	}
}

func getDevilCalcFunc(a, b float64) CalcFunc {
	var (
		aa = a * a
		bb = b * b
	)
	return func(t float64) Point {
		return devilCalcFunc(aa, bb, t)
	}
}
