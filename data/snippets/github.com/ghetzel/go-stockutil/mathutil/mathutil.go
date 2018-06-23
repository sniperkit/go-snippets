package mathutil

import (
	"math"
)

const (
	uvone    = 0x3FF0000000000000
	mask     = 0x7FF
	shift    = 64 - 11 - 1
	bias     = 1023
	signMask = 1 << 63
	fracMask = 1<<shift - 1
)

func Clamp(value float64, lower float64, upper float64) float64 {
	value = ClampLower(value, lower)
	value = ClampUpper(value, upper)
	return value
}

func ClampLower(value float64, lower float64) float64 {
	if value < lower {
		value = lower
	}

	return value
}

func ClampUpper(value float64, upper float64) float64 {
	if value > upper {
		value = upper
	}

	return value
}

// The Golang 1.10 math.Round implementation.
// (see: https://golang.org/pkg/math/#Round)
func Round(x float64) float64 {
	bits := math.Float64bits(x)
	e := uint(bits>>shift) & mask
	if e < bias {
		// Round abs(x) < 1 including denormals.
		bits &= signMask // +-0
		if e == bias-1 {
			bits |= uvone // +-1
		}
	} else if e < bias+shift {
		// Round any abs(x) >= 1 containing a fractional component [0,1).
		//
		// Numbers with larger exponents are returned unchanged since they
		// must be either an integer, infinity, or NaN.
		const half = 1 << (shift - 1)
		e -= bias
		bits += half >> e
		bits &^= fracMask >> e
	}
	return math.Float64frombits(bits)
}

func RoundPlaces(x float64, places int) float64 {
	multi := math.Pow(10, Clamp(float64(places), 0, 16))
	return (Round(x*multi) / multi)
}
