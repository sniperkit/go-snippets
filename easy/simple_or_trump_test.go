package main

import (
	"testing"
)

type tuple struct {
	l, r string
	t    byte
}

func TestWin(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{"AD", "2H", 'H'}:  "2H",
		tuple{"KD", "KH", 'C'}:  "KD KH",
		tuple{"JH", "10S", 'C'}: "JH"} {
		if r := win(k.l, k.r, k.t); r != v {
			t.Errorf("failed: win %s %s | %c is %s, got %s",
				k.l, k.r, k.t, v, r)
		}
	}
}

var v map[byte]int

func init() {
	v = map[byte]int{'2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7,
		'8': 8, '9': 9, '1': 10, 'J': 11, 'Q': 12, 'K': 13, 'A': 14}
}

func rank(c string, t byte) int {
	r := v[c[0]]
	if c[len(c)-1] == t {
		r += 13
	}
	return r
}

func win(l, r string, t byte) string {
	lr, rr := rank(l, t), rank(r, t)
	if lr > rr {
		return l
	} else if lr < rr {
		return r
	}
	return l + " " + r
}
