package main

import (
	"strings"
	"testing"
)

func TestToLower(t *testing.T) {
	for k, v := range map[string]string{
		"HELLO CODEEVAL":    "hello codeeval",
		"This is some Text": "this is some text",
		"Bash SM\\ASH":      "bash sm\\ash"} {
		if r := toLower(k); r != v {
			t.Errorf("failed: toLower %s is %s, got %s")
		}
	}
}

func toLower(q string) string {
	return strings.ToLower(q)
}
