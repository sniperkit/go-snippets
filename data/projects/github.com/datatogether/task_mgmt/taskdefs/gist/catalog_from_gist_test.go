package gist

import (
	"testing"
)

func TestGistIdFromUrl(t *testing.T) {
	cases := []struct {
		in, out string
		err     string
	}{
		{"https://api.github.com/gists/foo", "foo", "invalid gist url or id: 'foo'"},
		{"https://api.github.com/gists/5a309263a23b911f159a0ca8c8436496", "5a309263a23b911f159a0ca8c8436496", ""},
		{"5a309263a23b911f159a0ca8c8436496", "5a309263a23b911f159a0ca8c8436496", ""},
	}

	for i, c := range cases {
		got, err := GistIdFromUrl(c.in)
		if !(err == nil && c.err == "" || err != nil && err.Error() == c.err) {
			t.Errorf("case %d error mismatch: expected: %s, got: %s", i, c.err, err)
			continue
		}

		if got != c.out {
			t.Errorf("case %d mismatch. expected: %s, got: %s", i, c.out, got)
			continue
		}
	}
}
