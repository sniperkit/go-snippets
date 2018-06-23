package closures

// Fibonacci Returns a closure which yields numbers from the
// fibonacci sequence in order.
func Fibonacci() func() int {
	x, y := -1, 1
	return func() int {
		x, y = y, x+y
		return y
	}
}
