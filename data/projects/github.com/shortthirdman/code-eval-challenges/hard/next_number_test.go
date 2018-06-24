package main

import "testing"

func TestNextNumber(t *testing.T) {
	for k, v := range map[int]int{
		199999: 919999,
		987654: 4056789,
		123456: 123465,
		1:      10,
		999999: 9099999,
		100001: 100010,
		42:     204} {
		if r := nextNumber(k); r != v {
			t.Errorf("failed: nextNumber %d is %d, got %d",
				k, v, r)
		}
	}
}

func BenchmarkNextNumber(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nextNumber(i%999999 + 1)
	}
}

func dsig(a int) (r int) {
	for a > 0 {
		if a%10 > 0 {
			r += 1 << uint(3*(a%10))
		}
		a /= 10
	}
	return r
}

func nextNumber(d int) int {
	e := d + 9
	for ds := dsig(d); dsig(e) != ds; e += 9 {
	}
	return e
}
