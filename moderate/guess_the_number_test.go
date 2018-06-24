package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestGuessNumber(t *testing.T) {
	for k, v := range map[string]int{
		"100 Lower Lower Higher Lower Lower Lower Yay!": 13,
		"948 Higher Lower Yay!":                         593} {
		if r := guessNumber(k); r != v {
			t.Errorf("failed: guessNumber %s is %d, got %d",
				k, v, r)
		}
	}
}

func guessNumber(q string) int {
	s := strings.Fields(q)
	var c, lo, hi int
	fmt.Sscan(s[0], &hi)
	s = s[1 : len(s)-1]
	for _, i := range s {
		if i == "Lower" {
			hi = (lo+hi)/2 + c - 1
		} else {
			lo = (lo+hi)/2 + c + 1
		}
		c = (lo + hi) % 2
	}
	return (lo+hi)/2 + c
}
