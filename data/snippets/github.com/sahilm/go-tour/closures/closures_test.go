package closures_test

import (
	"testing"

	"github.com/sahilm/go-tour/closures"
)

func TestFibonacci(t *testing.T) {
	cases := []int{0, 1, 1, 2, 3, 5, 8, 13, 21}
	f := closures.Fibonacci()
	for i, want := range cases {
		got := f()
		if got != want {
			t.Errorf("got: %d, want: %d, at: %d", got, want, i)
		}
	}
}
