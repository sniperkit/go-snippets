package main

import (
	"strings"
	"testing"
)

func TestHiddenDigits(t *testing.T) {
	for k, v := range map[string]string{
		"abcdefghik":              "012345678",
		"Xa,}A#5N}{xOBwYBHIlH,#W": "05",
		"(ABW>'yy^'M{X-K}q,":      "NONE",
		"6240488":                 "6240488"} {
		if r := hiddenDigits(k); r != v {
			t.Errorf("failed: hiddenDigits %s is %s, got %s",
				k, v, r)
		}
	}
}

func hiddenDigits(q string) string {
	maph := func(r rune) rune {
		if strings.Contains("0123456789", string(r)) {
			return r
		} else if strings.Contains("abcdefghij", string(r)) {
			return r - 'a' + '0'
		}
		return -1
	}
	if r := strings.Map(maph, q); len(r) > 0 {
		return r
	} else {
		return "NONE"
	}
}
