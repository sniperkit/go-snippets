package main

import (
	"strings"
	"testing"
)

func TestWordToDigit(t *testing.T) {
	for k, v := range map[string]string{
		"zero;two;five;seven;eight;four": "025784",
		"three;seven;eight;nine;two":     "37892"} {
		if r := wordToDigit(k); r != v {
			t.Errorf("failed: wordToDigit %s is %s, got %s",
				k, v, r)
		}
	}
}

var digits map[string]string

func init() {
	digits = map[string]string{"zero": "0", "one": "1", "two": "2",
		"three": "3", "four": "4", "five": "5", "six": "6",
		"seven": "7", "eight": "8", "nine": "9"}
}

func wordToDigit(q string) string {
	t := strings.Split(q, ";")
	r := make([]string, len(t))
	for ix, i := range t {
		r[ix] = digits[i]
	}
	return strings.Join(r, "")
}
