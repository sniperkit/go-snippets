package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestRealFake(t *testing.T) {
	for k, v := range map[string]bool{
		"9999999999999999": false,
		"9999999999999993": true} {
		if r := realFake(k); r != v {
			t.Errorf("failed: realFake %s is %t, got %t",
				k, v, r)
		}
	}
}

func BenchmarkRealFake(b *testing.B) {
	for i := 1; i < b.N; i++ {
		realFake(fmt.Sprint(uint64(i) + 1000000000000000))
	}
}

func realFake(q string) bool {
	s := strings.Split(q, "")
	t := make([]int, len(s))
	for i := range s {
		fmt.Sscan(s[i], &t[i])
	}
	for i := len(t) - 2; i >= 0; i -= 2 {
		t[i] *= 2
	}
	var su int
	for i := range t {
		su += t[i]
	}
	return su%10 == 0
}
