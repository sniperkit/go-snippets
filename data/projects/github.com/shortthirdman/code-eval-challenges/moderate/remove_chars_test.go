package main

import (
	"strings"
	"testing"
)

type tuple struct {
	p, q string
}

func TestRemoveChars(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{"how are you", "abc"}: "how re you",
		tuple{"hello world", "def"}: "hllo worl"} {
		if r := removeChars(k.p, k.q); r != v {
			t.Errorf("failed: removeChars %s, %s is %s, got %s",
				k.p, k.q, v, r)
		}
	}
}

func removeChars(p, q string) string {
	mapf := func(r rune) rune {
		if strings.Contains(q, string(r)) {
			return -1
		}
		return r
	}
	return strings.Map(mapf, p)
}
