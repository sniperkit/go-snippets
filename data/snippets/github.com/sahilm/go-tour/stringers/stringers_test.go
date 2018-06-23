package stringers_test

import (
	"testing"

	"fmt"

	"github.com/sahilm/go-tour/stringers"
)

func TestIPAddrIsAStringer(t *testing.T) {
	cases := []struct {
		in   stringers.IPAddr
		want string
	}{
		{
			stringers.IPAddr{127, 0, 0, 1},
			"127.0.0.1",
		},
		{
			stringers.IPAddr{8, 8, 8, 8},
			"8.8.8.8",
		},
	}

	for _, c := range cases {
		got := fmt.Sprint(c.in)
		if got != c.want {
			t.Errorf("got: %v, want: %v", c.want, got)
		}
	}
}
