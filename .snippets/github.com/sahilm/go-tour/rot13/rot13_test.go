package rot13_test

import (
	"strings"
	"testing"

	"io"

	"io/ioutil"

	"github.com/sahilm/go-tour/rot13"
)

func TestRot13Reader(t *testing.T) {
	cases := []struct {
		in   io.Reader
		want string
	}{
		{strings.NewReader("Lbh penpxrq gur pbqr!"), "You cracked the code!"},
		{strings.NewReader("You cracked the code!"), "Lbh penpxrq gur pbqr!"},
		{strings.NewReader("   "), "   "},
	}

	for _, c := range cases {

		b, err := ioutil.ReadAll(rot13.Reader{R: c.in})
		got := string(b)

		if err != nil {
			t.Errorf("got error: %v, want no error", err)
		}

		if got != c.want {
			t.Errorf("got: %v, want: %v", got, c.want)
		}
	}
}
