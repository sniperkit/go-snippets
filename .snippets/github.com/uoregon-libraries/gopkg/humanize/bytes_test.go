package humanize

import (
	"math"
	"testing"
)

// TestBytes verifies various values get converted correctly
func TestBytes(t *testing.T) {
	var valueExpectations = []struct {
		size int64
		str  string
	}{
		{size: 50, str: "50 B"},
		{size: 1023, str: "1023 B"},
		{size: 1024, str: "1 KB"},
		{size: Terabyte - Gigabyte, str: "1023 GB"},
		{size: Terabyte - 1, str: "1023.99 GB"},
		{size: Terabyte, str: "1 TB"},
		{size: Terabyte + Megabyte, str: "1 TB"},
		{size: 5000000000, str: "4.65 GB"},
		{size: math.MaxInt64, str: "7.99 EB"},
	}

	for _, expect := range valueExpectations {
		var actual = Bytes(expect.size)
		if actual != expect.str {
			t.Errorf("%d should have given us %q, but instead we got %q", expect.size, expect.str, actual)
		}
	}
}

func BenchmarkBytes(b *testing.B) {
	var tests = []int64{0, 1, 2, 3, 4, 1024, 2048, 5000, 5000000000000, math.MaxInt64, 5000000}
	tests = []int64{2048, 5000, 5000000000000, math.MaxInt64, 5000000}
	b.ResetTimer()

	b.Run("Bytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, t := range tests {
				Bytes(t)
			}
		}
	})

	b.Run("bytesSimple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, t := range tests {
				bytesSimple(t)
			}
		}
	})
}
