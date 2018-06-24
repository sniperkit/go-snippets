package main

import (
	"strings"
	"testing"
)

func TestPenultimate(t *testing.T) {
	for k, v := range map[string]string{
		"some line with text":               "with",
		"Each line has more than one word.": "one",
		"another line":                      "another"} {
		if r := penultimate(k); r != v {
			t.Errorf("failed: penultimate %s is %s, got %s",
				k, v, r)
		}
	}
}

func penultimate(q string) string {
	t := strings.Fields(q)
	if len(t) < 2 {
		return ""
	}
	return t[len(t)-2]
}
