package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	for k, v := range map[string]uint{
		"12": 2, "123": 3,
		"4651236715636226712123765123": 432,
		"66620": 1} {
		if r := decode(k); r != v {
			t.Errorf("failed: decode %s is %d, got %d",
				k, v, r)
		}
	}
}

func BenchmarkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var s []string
		for j := i + 1; j > 0; j /= 26 {
			s = append(s, fmt.Sprint(j%26+1))
		}
		decode(strings.Join(s, ""))
	}
}

func decode(s string) uint {
	for len(s) > 1 {
		switch {
		case s[0] != '1' && s[0] != '2':
			s = s[1:len(s)]
		case s[1] == '0' || (s[0] == '2' && s[1] > '6'):
			s = s[2:len(s)]
		default:
			if len(s) == 2 {
				return 2
			}
			return decode(s[1:len(s)]) + decode(s[2:len(s)])
		}
	}
	return 1
}
