package main

import "testing"

func TestRomanArabic(t *testing.T) {
	for k, v := range map[string]int{
		"3M1D2C":     3700,
		"2I3I2X9V1X": -16} {
		if r := romanArabic(k); r != v {
			t.Errorf("failed: romanArabic %s is %d, got %d",
				k, v, r)
		}
	}
}

func romanArabic(q string) int {
	var n, b, r, s int
	for i := 0; i < len(q); i += 2 {
		a := int(q[i] - '0')
		switch q[i+1] {
		case 'M':
			r = 1000
		case 'D':
			r = 500
		case 'C':
			r = 100
		case 'L':
			r = 50
		case 'X':
			r = 10
		case 'V':
			r = 5
		case 'I':
			r = 1
		}
		if r > s {
			n -= b * s
		} else {
			n += b * s
		}
		b, s = a, r
	}
	return n + b*s
}
