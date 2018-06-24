package main

import (
	"strings"
	"testing"
)

func TestDeco(t *testing.T) {
	for k, v := range map[string]string{
		"mke": "try", "mh": "to", "lhsby": "solve", "pm": "it"} {
		if r := strings.Map(deco, k); r != v {
			t.Errorf("failed: deco %s is %s, got %s",
				k, v, r)
		}
	}
}

func deco(r rune) rune {
	switch {
	case r >= 'a' && r <= 'f':
		return r + 20
	case r >= 'g' && r <= 'm':
		return r + 7
	case r >= 'n' && r <= 't':
		return r - 7
	case r >= 'u' && r <= 'z':
		return r - 20
	}
	return r
}
