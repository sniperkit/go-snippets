package main

import (
	"strings"
	"testing"
)

func TestCycle(t *testing.T) {
	for k, v := range map[string]string{
		"2 0 6 3 1 6 3 1 6 3 1":         "6 3 1",
		"2 0 6 3 1 5 7 6 3 1 8 9 6 3 1": "6 3 1 5 7",
		"99":                                  "",
		"":                                    "",
		"1 2 3 1 2 3":                         "1 2 3",
		"1 2 3 3":                             "3",
		"1 2 3 4 5 6":                         "",
		"12 14 16 15 17 12 14 16 15 17 12 14": "12 14 16 15 17",
		"1 1 2 3 1 2 3":                       "1"} {
		if r := cycle(k); r != v {
			t.Errorf("failed: cycle %s is %s, got %s",
				k, v, r)
		}
	}
}

func uniq(s []string, t string) bool {
	for _, i := range s {
		if i == t {
			return false
		}
	}
	return true
}

func cycle(q string) string {
	s := strings.Fields(q)
	for len(s) > 1 && uniq(s[1:len(s)], s[0]) {
		s = s[1:len(s)]
	}
	if len(s) <= 1 {
		return ""
	}
	t := []string{s[0]}
	for i := 1; i < len(s) && s[i] != s[0]; i++ {
		t = append(t, s[i])
	}
	return strings.Join(t, " ")
}
