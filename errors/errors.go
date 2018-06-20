package errors

import (
	"fmt"
	"math"
)

// ErrNegativeSqrt represents an error which occurs when
// the square root of a negative number is requested.
type ErrNegativeSqrt float64

// The error message returned by ErrNegativeSqrt
func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot sqrt negative number: %v", float64(e))
}

// Sqrt calculates the square root of x using the Newton-Raphson method.
// For positive values of x, the square root, number of iterations to process the result
// and nil error is returned. A ErrNegativeSqrt is returned if x is negative.
func Sqrt(x float64) (float64, int, error) {
	if x < 0 {
		return 0, 0, ErrNegativeSqrt(x)
	}
	const minDelta = 1e-8
	z := 1.0
	delta := 1.0
	iterations := 0
	for ; delta > minDelta; iterations++ {
		next := z - (z*z-x)/(2*z)
		delta = math.Abs(z - next)
		z = next
	}
	return z, iterations, nil
}
