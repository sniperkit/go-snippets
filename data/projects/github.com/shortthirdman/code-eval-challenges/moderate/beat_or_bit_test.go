package main

import (
	"fmt"
	"strconv" // remove as soon as Go supports binary literals
	"strings"
	"testing"
)

func TestG2d(t *testing.T) {
	for k, v := range map[string]uint{
		"10": 3, "101": 6, "1111": 10,
		"1110": 11, "1100001": 65} {
		b, _ := strconv.ParseInt(k, 2, 0)
		if r := g2d(uint(b)); r != v {
			t.Errorf("failed: g2d %s is %d, got %d",
				k, v, r)
		}
	}
}

func TestBeatOrBit(t *testing.T) {
	for k, v := range map[string]string{
		"1111 | 1110":        "10 | 11",
		"10 | 1100001 | 101": "3 | 65 | 6"} {
		if r := beatOrBit(k); r != v {
			t.Errorf("failed: beatOrBit %s is %s, got %s",
				k, v, r)
		}
	}
}

func BenchmarkG2d(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g2d(uint(i % 128))
	}
}

func g2d(a uint) uint {
	a ^= a >> 4
	a ^= a >> 2
	a ^= a >> 1
	return a
}

func beatOrBit(q string) string {
	var v uint
	s := strings.Split(q, " | ")
	for i := range s {
		fmt.Sscanf(s[i], "%b", &v)
		s[i] = fmt.Sprint(g2d(v))
	}
	return strings.Join(s, " | ")
}
