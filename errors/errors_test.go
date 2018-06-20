package errors_test

import (
	"math"
	"testing"

	"github.com/sahilm/go-tour/errors"
)

func TestSqrtErrors(t *testing.T) {
	cases := []struct {
		in    float64
		want  float64
		error string
	}{
		{-2, 0, "cannot sqrt negative number: -2"},
		{-100, 0, "cannot sqrt negative number: -100"},
		{2, math.Sqrt(2), ""},
	}

	for _, c := range cases {
		got, _, err := errors.Sqrt(c.in)

		if got != c.want {
			t.Errorf("got: %v, want: %v", got, c.want)
		}

		switch err {
		case nil:
			if c.error != "" {
				t.Errorf("got no error, want error: %v", c.error)
			}

		default:
			got := err.Error()
			if got != c.error {
				t.Errorf("got error: %v, want error: %v", got, c.error)
			}
		}
	}
}
