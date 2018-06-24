package main

import (
	"strings"
	"testing"
)

func TestPangram(t *testing.T) {
	for k, v := range map[string]string{
		"A quick brown fox jumps over the lazy dog":        "NULL",
		"A slow Yellow fox crawls Under the proactive dog": "bjkmqz"} {
		if r := pangram(k); r != v {
			t.Errorf("failed: pangram %s is %s, got %s",
				k, v, r)
		}
	}
}

func BenchmarkPangram(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var s []byte
		for j := i; j > 0; j /= 26 {
			s = append(s, byte(j%26+97))
		}
		pangram(string(s))
	}
}

func pangram(q string) string {
	var r []byte
	q = strings.ToLower(q)
	for i := byte('a'); i <= 'z'; i++ {
		if !strings.Contains(q, string(i)) {
			r = append(r, i)
		}
	}
	if len(r) == 0 {
		return "NULL"
	}
	return string(r)
}
