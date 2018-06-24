package main

import "testing"

type tuple struct {
	l, u float64
}

func TestLcase(t *testing.T) {
	for k, v := range map[string]tuple{
		"-":        tuple{0, 0},
		"thisTHIS": tuple{50, 50},
		"aBCD":     tuple{25, 75}} {
		if rl, ru := lcase(k); rl != v.l || ru != v.u {
			t.Errorf("failed: lcase %s is %f %f, got %f %f", k, v.l, v.u, rl, ru)
		}
	}
}

func BenchmarkLcase(b *testing.B) {
	h := []byte{32, 65, 66, 67, 97, 98, 99}
	for i := 0; i < b.N; i++ {
		var s []byte
		for j := i; j > 0; j /= 7 {
			s = append(s, h[i%7])
		}
		lcase(string(s))
	}
}

func lcase(s string) (float64, float64) {
	var c, l float64
	for _, i := range s {
		if i >= 'a' && i <= 'z' {
			l++
			c++
		} else if i >= 'A' && i <= 'Z' {
			c++
		}
	}
	if c == 0 {
		return 0, 0
	}
	r := 100 * l / c
	return r, 100 - r
}
