package bits

func Sum(x, y int) int {
	return (x & y) + (x | y)
	// return (x ^ y) + 2*(x&y)
}

func Avg(x, y int) int {
	return Sum(x, y) >> 1
}

func Abs(x int64) int64 {
	y := x >> 63
	return (x + y) ^ y
}

func Reverse(x int) int {
	x = ((x & 0xaaaaaaaa) >> 1) | ((x & 0x55555555) << 1)
	x = ((x & 0xcccccccc) >> 2) | ((x & 0x33333333) << 2)
	x = ((x & 0xf0f0f0f0) >> 4) | ((x & 0x0f0f0f0f) << 4)
	x = ((x & 0xff00ff00) >> 8) | ((x & 0x00ff00ff) << 8)
	return (x >> 16) | (x << 16)
}

func ToGray(n int) int {
	return n ^ (n >> 1)
}

func FromGray(g int) int {
	n := 0
	for ; g > 0; g >>= 1 {
		n ^= g
	}
	return n

	// gray ^= (gray >> 16)
	// gray ^= (gray >> 8)
	// gray ^= (gray >> 4)
	// gray ^= (gray >> 2)
	// gray ^= (gray >> 1)
	// return gray
}

func IsPowerOf2(x int) bool {
	return (x & (x - 1)) == 0
}

func LSB(x int) int {
	return x & -x
}

func MSB(x int) int {
	x |= (x >> 1)
	x |= (x >> 2)
	x |= (x >> 4)
	x |= (x >> 8)
	x |= (x >> 16)
	return x & ^(x >> 1)
}
