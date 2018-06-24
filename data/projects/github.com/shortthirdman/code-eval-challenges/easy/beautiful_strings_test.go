package main

import (
	"sort"
	"strings"
	"testing"
)

func TestBeauty(t *testing.T) {
	for k, v := range map[string]int{"ABbCcc": 152, "<3": 0,
		"The quick brown fox jumps over the lazy dog":     569,
		"Good luck in the Facebook Hacker Cup this year!": 754,
		"Ignore punctuation, please :)":                   491,
		"Sometimes test cases are hard to make up.":       729,
		"So I just go consult Professor Dalves":           646} {
		if r := beauty(k); r != v {
			t.Errorf("failed: beauty %s is %d, got %d",
				k, v, r)
		}
	}
}

func BenchmarkBeauty(b *testing.B) {
	var i uint64
	for i = 0; i < uint64(b.N); i++ {
		var s string
		for n := i + (1 << 63); n > 0; n /= 58 {
			s += string(rune(n%58 + 64))
		}
		beauty(strings.Repeat(s, 10))
	}
}

var c []int

func init() {
	c = make([]int, 26)
}

func clean(char rune) rune {
	if char >= 'a' && char <= 'z' {
		return char
	}
	return -1
}

func beauty(s string) (r int) {
	t := strings.Map(clean, strings.ToLower(s))
	for _, i := range t {
		c[i-'a']++
	}
	sort.Ints(c)
	for i := 26; i > 0 && c[i-1] > 0; i-- {
		r += i * c[i-1]
		c[i-1] = 0
	}
	return r
}
