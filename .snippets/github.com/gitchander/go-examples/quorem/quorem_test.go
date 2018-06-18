package quorem

import (
	"testing"
)

func TestQuoRem(t *testing.T) {

	const n = 10000
	for x := 1; x < n; x++ {
		for y := 1; y < n; y++ {

			quo1, rem1 := quoRem1(x, y)
			quo2, rem2 := quoRem2(x, y)

			if quo1 != quo2 {
				t.Fatalf("quoRem(%d, %d): quo1(%d) != quo2(%d)", x, y, quo1, quo2)
			}

			if rem1 != rem2 {
				t.Fatalf("quoRem(%d, %d): rem1(%d) != rem2(%d)", x, y, rem1, rem2)
			}

			if quo1*y+rem1 != x {
				t.Fatalf("quoRem(%d, %d): quo(%d), rem(%d)", x, y, quo1, rem1)
			}
		}
	}
}

// Mersenne primes
const (
	m17 = 131071
	m31 = 2147483647
)

func BenchmarkQuoRem1(b *testing.B) {
	var x, y = m31, m17
	for i := 0; i < b.N; i++ {
		quoRem1(x, y)
	}
}

func BenchmarkQuoRem2(b *testing.B) {
	var x, y = m31, m17
	for i := 0; i < b.N; i++ {
		quoRem2(x, y)
	}
}

func BenchmarkQuoRem3(b *testing.B) {
	var x, y = m31, m17
	for i := 0; i < b.N; i++ {
		quo := x / y
		rem := x - quo*y
		_ = rem
	}
}
