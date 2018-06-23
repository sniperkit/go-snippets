package loopsandfunctions

import "math"

// Sqrt calculates the square root of x using the Newton-Raphson method.
// For positive values of x, the square root and number of iterations to process the result
// are returned. The behaviour is undefined for negative values of x.
func Sqrt(x float64) (float64, int) {
	const minDelta = 1e-8
	z := 1.0
	delta := 1.0
	iterations := 0
	for ; delta > minDelta; iterations++ {
		next := z - (z*z-x)/(2*z)
		delta = math.Abs(z - next)
		z = next
	}
	return z, iterations
}
