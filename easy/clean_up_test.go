package main

import (
	"strings"
	"testing"
)

func TestClean(t *testing.T) {
	for k, v := range map[string]string{
		"(--9Hello----World...--)":         "hello world",
		"Can 0$9 ---you~":                  "can you",
		"13What213are;11you-123+138doing7": "what are you doing"} {
		if r := clean(k); r != v {
			t.Errorf("failed: clean %s is %s, got %s",
				k, v, r)
		}
	}
}

func isAlpha(a rune) bool {
	return (a >= 'a' && a <= 'z') || (a >= 'A' && a <= 'Z')
}

func clean(q string) string {
	var (
		c []byte
		r []string
	)
	w := false
	for _, i := range q {
		if isAlpha(i) {
			w = true
			c = append(c, byte(i)|32)
		} else if w {
			w = false
			if len(c) > 0 {
				r = append(r, string(c))
				c = []byte{}
			}
		}
	}
	if len(c) > 0 {
		r = append(r, string(c))
	}
	return strings.Join(r, " ")
}
