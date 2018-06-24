package main

import (
	"fmt"
	"strings"
	"testing"
)

var tests map[string]string

func init() {
	tests = map[string]string{
		"40 40 40 40 29 29 29 29 29 29 29 29 57 57 92 92 92 92 92 86 86 86 86 86 86 86 86 86 86": "4 40 8 29 2 57 5 92 10 86",
		"73 73 73 73 41 41 41 41 41 41 41 41 41 41":                                              "4 73 10 41",
		"1 1 3 3 3 2 2 2 2 14 14 14 11 11 11 2":                                                  "2 1 3 3 4 2 3 14 3 11 1 2",
		"7": "1 7"}
}

func TestCompress(t *testing.T) {
	for k, v := range tests {
		if r := compress(k); r != v {
			t.Errorf("failed: compress %s is %s, got %s",
				k, v, r)
		}
	}
}

func TestCompressInt(t *testing.T) {
	for k, v := range tests {
		if r := compressInt(k); r != v {
			t.Errorf("failed: compressInt %s is %s, got %s",
				k, v, r)
		}
	}
}

func BenchmarkCompress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var s []string
		for j := i; j > 0; j /= 48 {
			s = append(s, strings.Repeat(fmt.Sprint(j%12+1), (j/12)*4+1))
		}
		compress(strings.Join(s, " "))
	}
}

func BenchmarkCompressInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var s []string
		for j := i; j > 0; j /= 48 {
			s = append(s, strings.Repeat(fmt.Sprint(j%12+1), (j/12)*4+1))
		}
		compressInt(strings.Join(s, " "))
	}
}

func compress(s string) string {
	var (
		r []string
		n string
		c int
	)
	for _, i := range strings.Fields(s) {
		if i == n {
			c++
			continue
		}
		if c > 0 {
			r = append(append(r, fmt.Sprint(c)), n)
		}
		c, n = 1, i
	}
	r = append(append(r, fmt.Sprint(c)), n)
	return strings.Join(r, " ")
}

func compressInt(s string) string {
	var (
		r       []string
		n, c, k int
	)
	for _, i := range strings.Fields(s) {
		fmt.Sscan(i, &k)
		if k == n {
			c++
			continue
		}
		if c > 0 {
			r = append(append(r, fmt.Sprint(c)), fmt.Sprint(n))
		}
		c, n = 1, k
	}
	r = append(append(r, fmt.Sprint(c)), fmt.Sprint(n))
	return strings.Join(r, " ")
}
