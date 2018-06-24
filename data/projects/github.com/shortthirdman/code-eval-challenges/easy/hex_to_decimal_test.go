package main

import (
	"fmt"
	"testing"
)

func TestHex2dec(t *testing.T) {
	for k, v := range map[string]int{
		"9f":       159,
		"11":       17,
		"0":        0,
		"499602d2": 1234567890} {
		if r := hex2dec(k); r != v {
			t.Errorf("failed: hex2dec %s is %d, got %d",
				k, v, r)
		}
	}
}

func BenchmarkHex2dec(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hex2dec(fmt.Sprintf("%x", i))
	}
}

func hex2dec(q string) (r int) {
	fmt.Sscanf(q, "%x", &r)
	return r
}
