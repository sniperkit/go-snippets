package main

import "testing"

func TestIsBalanced(t *testing.T) {
	for k, v := range map[string]bool{
		":))": false, "(:)": true, "j": true, "": true, ":):": true,
		"( :) ))": false, "asdf(asdf:)ghjk:ghjk(ghjk)": true,
		"(((::))()": false} {
		if r := isBalanced(trimRight(k), 0); r != v {
			t.Errorf("failed: isBalanced2 %s is %t, got %t",
				k, v, r)
		}
	}
}

func BenchmarkIsBalanced(b *testing.B) {
	ch := []string{"(", ")", ":", "j", " "}
	for i := 0; i < b.N; i++ {
		var (
			s string
			j int
		)
		for j = i; j > 0; j /= 5 {
			s += ch[j%5]
		}
		isBalanced(trimRight(s), 0)
	}
}

func isBalanced(s string, c int) bool {
	for ; s != "" && c >= 0; s = s[1:] {
		switch s[0] {
		case '(':
			c++
		case ')':
			c--
		case ':':
			if s[1] == '(' {
				return isBalanced(s[2:], c) || isBalanced(s[2:], c+1)
			} else if s[1] == ')' {
				return isBalanced(s[2:], c) || isBalanced(s[2:], c-1)
			}
		}
	}
	return c == 0
}

func trimRight(s string) string {
	r := len(s)
	for ; r > 0 && s[r-1] != '(' && s[r-1] != ')'; r-- {
	}
	return s[:r]
}
