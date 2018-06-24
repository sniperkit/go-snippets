package main

import "testing"

func TestSwapCase(t *testing.T) {
	for k, v := range map[byte]byte{
		'a': 'A', 'z': 'Z', 'A': 'a', 'Z': 'z',
		'+': '+', '.': '.', ' ': ' ', '\n': '\n'} {
		if r := swapCase(k); r != v {
			t.Errorf("failed: swapCase %c is %c, got %c",
				k, v, r)
		}
	}
}

func BenchmarkSwapCase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		swapCase(byte(i))
	}
}

func swapCase(c byte) byte {
	switch {
	case c >= 'a' && c <= 'z':
		return c & 223
	case c >= 'A' && c <= 'Z':
		return c | 32
	}
	return c
}
