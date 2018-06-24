package main

import (
	"fmt"
	"strings"
	"testing"
)

type tuple struct {
	p, q string
}

func TestValid(t *testing.T) {
	if valid(1, true, 2) {
		t.Error("failed: valid 1 true 2 is false, got true")
	}
	if !valid(3, true, 7) {
		t.Error("failed: valid 7 true 3 is true, got false")
	}
	if valid(3, false, 6) {
		t.Error("failed: valid 6 false 3 is false, got true")
	}
}

func TestBrokenLcd(t *testing.T) {
	for k, v := range map[tuple]bool{
		tuple{"10110001 11111000 11111110 11111111 11111111 11111111 11111111 11101101 11111111 01111111 11110010 10100111", "84.525784"}: true,
		tuple{"11111111 11110110 11101111 11110111 10111110 11110110 10111011 10100111 11111100 01100100 11111101 01011110", "5.57"}:      true,
		tuple{"11000010 00001111 11111111 10111111 11101011 11110011 01111110 11011111 11111111 11111111 11111001 01101110", "857.71284"}: true,
		tuple{"11111111 01110111 10111011 11001101 11111011 11101010 11110100 01001101 11011111 11111010 10010110 10111111", "66.92"}:     false,
		tuple{"11111011 10010001 11111011 11111101 10011111 10111110 01111100 11011101 10111001 11111110 11101111 11110110", "188.87"}:    false} {
		if r := brokenLcd(k.p, k.q); r != v {
			t.Errorf("failed: brokenLcd %s is %t, got %t",
				k.p, k.q, v, r)
		}
	}
}

func BenchmarkValid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		valid(i%256, (i/256)&1 == 1, (i/512)%256)
	}
}

var digits map[uint8]int

func init() {
	digits = map[uint8]int{'0': 252, '1': 96,
		'2': 218, '3': 242, '4': 102,
		'5': 182, '6': 190, '7': 224,
		'8': 254, '9': 246}
}

func valid(y int, d bool, x int) bool {
	if d && (x&1 == 0) {
		return false
	}
	return y&x == y
}

func brokenLcd(p, q string) bool {
	var v bool
	t := strings.Fields(p)
	for u := 0; u <= len(t)-len(q)+1; u++ {
		var j int
		v = true
		for i := 0; i < len(q); i++ {
			d := i < len(q)-1 && q[i+1] == '.'
			var x int
			fmt.Sscanf(t[u+j], "%b", &x)
			if !valid(digits[q[i]], d, x) {
				v = false
				break
			}
			if d {
				i++
			}
			j++
		}
		if v {
			return true
		}
	}
	return v
}
