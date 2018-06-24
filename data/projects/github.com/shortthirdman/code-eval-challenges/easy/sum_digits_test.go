package main

import (
	"fmt"
	"testing"
)

func TestSumDigits(t *testing.T) {
	for k, v := range map[string]uint{
		"23": 5, "496": 19, "120345": 15} {
		if r := sumDigits(k); r != v {
			t.Errorf("failed: sumDigits %s is %d, got %d",
				k, v, r)
		}
	}
}

func BenchmarkSumDigits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sumDigits(fmt.Sprint(i))
	}
}

func sumDigits(s string) (r uint) {
	for _, i := range s {
		if i >= '0' && i <= '9' {
			r += uint(i - '0')
		}
	}
	return r
}
