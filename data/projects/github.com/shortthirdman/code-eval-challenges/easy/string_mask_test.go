package main

import "testing"

type tuple struct {
	s, t string
}

func TestMask(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{"hello", "11001"}: "HEllO",
		tuple{"world", "10000"}: "World",
		tuple{"cba", "111"}:     "CBA"} {
		if r := mask(k.s, k.t); r != v {
			t.Errorf("failed: mask %s %s is %s, got %s",
				k.s, k.t, v, r)
		}
	}
}

func BenchmarkMask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var s, t []byte
		for j := i; j > 0; j /= 26 {
			s = append(s, byte(j%26+97))
			t = append(t, '1')
		}
		mask(string(s), string(t))
	}
}

func upper(b byte) byte {
	return b & 223
}

func mask(s, t string) string {
	r := make([]byte, len(s))
	for ix, i := range s {
		if t[ix] == '1' {
			r[ix] = upper(byte(i))
		} else {
			r[ix] = byte(i)
		}
	}
	return string(r)
}
