package plugin

import (
	"testing"
)

func TestExitCode(t *testing.T) {
	tests := []struct {
		in   Status
		code int
	}{
		{OK, 0},
		{WARNING, 1},
		{CRITICAL, 2},
		{UNKNOWN, 3},
	}

	for _, test := range tests {
		out := test.in.ExitCode()
		if test.code != out {
			t.Errorf("Got %d, expected %d", out, test.code)
		}
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		in   Status
		text string
	}{
		{OK, "OK"},
		{WARNING, "WARNING"},
		{CRITICAL, "CRITICAL"},
		{UNKNOWN, "UNKNOWN"},
	}

	for _, test := range tests {
		out := test.in.String()
		if test.text != out {
			t.Errorf("Got %s, expected %s", out, test.text)
		}
	}
}
