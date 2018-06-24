package main

import (
	"strings"
	"testing"
)

func TestCapitalize(t *testing.T) {
	for k, v := range map[string]string{
		"Hello world":         "Hello World",
		"javaScript language": "JavaScript Language",
		"a letter":            "A Letter",
		"1st thing":           "1st Thing"} {
		if r := capitalize(k); r != v {
			t.Errorf("failed: capitalize %s is %s, got %s",
				k, v, r)
		}
	}
}

func capitalize(q string) string {
	return strings.Title(q)
}
