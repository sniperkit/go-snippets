package main

import (
	"strings"
	"testing"
)

type tuple struct {
	l, r uint8
}

func TestPos(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{0, 0}: "a1",
		tuple{7, 7}: "h8"} {
		if r := pos(k.l, k.r); r != v {
			t.Errorf("failed: pos %d %d is %s, got %s",
				k.l, k.r, v, r)
		}
	}
}

func TestPosConcat(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{0, 0}: "a1",
		tuple{7, 7}: "h8"} {
		if r := posConcat(k.l, k.r); r != v {
			t.Errorf("failed: posConcat %d %d is %s, got %s",
				k.l, k.r, v, r)
		}
	}
}

func TestKnight(t *testing.T) {
	for k, v := range map[string]string{
		"g2": "e1 e3 f4 h4",
		"a1": "b3 c2",
		"d6": "b5 b7 c4 c8 e4 e8 f5 f7",
		"e5": "c4 c6 d3 d7 f3 f7 g4 g6",
		"b1": "a3 c3 d2"} {
		if r := knight(k); r != v {
			t.Errorf("failed: knight %s is %s, got %s",
				k, v, r)
		}
	}
}

func BenchmarkPos(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pos(uint8(i%8), uint8((i/8)%8))
	}
}

func BenchmarkPosConcat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		posConcat(uint8(i%8), uint8((i/8)%8))
	}
}

func BenchmarkKnight(b *testing.B) {
	for i := 0; i < b.N; i++ {
		knight(pos(uint8(i%8), uint8((i/8)%8)))
	}
}

func pos(l, r uint8) string {
	return string([]byte{'a' + l, '1' + r})
}

func posConcat(l, r uint8) string {
	return string('a'+l) + string('1'+r)
}

func knight(q string) string {
	var m []string
	l, r := q[0]-'a', q[1]-'1'
	if l > 1 {
		if r > 0 {
			m = append(m, pos(l-2, r-1))
		}
		if r < 7 {
			m = append(m, pos(l-2, r+1))
		}
	}
	if l > 0 {
		if r > 1 {
			m = append(m, pos(l-1, r-2))
		}
		if r < 6 {
			m = append(m, pos(l-1, r+2))
		}
	}
	if l < 7 {
		if r > 1 {
			m = append(m, pos(l+1, r-2))
		}
		if r < 6 {
			m = append(m, pos(l+1, r+2))
		}
	}
	if l < 6 {
		if r > 0 {
			m = append(m, pos(l+2, r-1))
		}
		if r < 7 {
			m = append(m, pos(l+2, r+1))
		}
	}
	return strings.Join(m, " ")
}
