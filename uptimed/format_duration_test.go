package main

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	cases := []struct {
		Input    time.Duration
		Expected string
	}{
		{time.Duration(0), "00:00"},
		{time.Second, "00:00"},
		{59 * time.Second, "00:00"},
		{time.Minute, "00:01"},
		{time.Minute + (59 * time.Second), "00:01"},
		{time.Hour, "01:00"},
		{time.Hour + (59 * time.Minute), "01:59"},
		{time.Hour + (59 * time.Minute) + (39 * time.Second), "01:59"},
	}
	var given string
	for _, c := range cases {
		given = formatDuration(c.Input)
		if given != c.Expected {
			t.Errorf("Unexpected result for duration %q. Expected: %s, Given: %s", c.Input, c.Expected, given)
		}
	}
}
