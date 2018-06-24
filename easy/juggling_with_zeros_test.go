package main

import (
	"strings"
	"testing"
)

func TestJuggle(t *testing.T) {
	for k, v := range map[string]uint64{
		"00 0 0 0 00 0 0 00 00 0000 0 00 00 0 0 00 00 0 0 000 00 0 0 0 00 0 0 0 00 0 0 0 00 0 0 0 00 0 0 00 00 00 0 0": 2811374246,
		"00 0 0 00 00 0": 9,
		"00 0":           1,
		"00 0 0 000 00 0000000 0 000":                                         9208,
		"0 000000000 00 00":                                                   3,
		"00 0000000000000000000000000000000000000000000000000000000000000000": 18446744073709551615} {
		if r := juggle(k); r != v {
			t.Errorf("failed: juggle %s is %d, got %d",
				k, v, r)
		}
	}
}

func BenchmarkJuggle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		juggle(unjuggle(uint64(i)))
	}
}

func juggle(q string) (r uint64) {
	t := strings.Fields(q)
	for i := 0; i < len(t); i += 2 {
		r <<= uint(len(t[i+1]))
		if t[i] == "00" {
			r += (1 << uint(len(t[i+1]))) - 1
		}
	}
	return r
}

func unjuggle(n uint64) (s string) {
	if n == 0 {
		s = "0 0"
		return s
	}
	m0, m1 := 0, 0
	for n > 0 {
		if n&1 == 0 {
			if m1 == 0 {
				m0++
			} else {
				if s == "" {
					s = "00 " + strings.Repeat("0", m1)
				} else {
					s = "00 " + strings.Repeat("0", m1) + " " + s
				}
				m0, m1 = 1, 0
			}
		} else {
			if m0 == 0 {
				m1++
			} else {
				if s == "" {
					s = "0 " + strings.Repeat("0", m0)
				} else {
					s = "0 " + strings.Repeat("0", m0) + " " + s
				}
				m0, m1 = 0, 1
			}
		}
		n >>= 1
	}
	if m1 > 0 {
		if s == "" {
			s = "00 " + strings.Repeat("0", m1)
		} else {
			s = "00 " + strings.Repeat("0", m1) + " " + s
		}
	}
	return s
}
