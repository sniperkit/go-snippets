package main

import "testing"

func TestHappy(t *testing.T) {
	for k, v := range map[uint]uint{
		1: 1, 22: 8, 8: 64, 64: 52, 52: 29, 29: 85,
		85: 89, 89: 145, 145: 42, 42: 20, 20: 4, 4: 16, 16: 37,
		37: 58, 58: 89, 7: 49, 49: 97, 97: 130, 130: 10, 10: 1} {
		if r := happy(k); r != v {
			t.Errorf("failed: happy %d is %d, got %d",
				k, v, r)
		}
	}
}

func TestHappyNumbers(t *testing.T) {
	for k, v := range map[uint]bool{
		1:                    true,
		7:                    true,
		22:                   false,
		12345678901234567890: false} {
		if r := happyNumbers(k); r != v {
			t.Errorf("failed: happyNumbers %d is %t, got %t",
				k, v, r)
		}
	}
}

func BenchmarkHappy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		happy(uint(i))
	}
}

func BenchmarkHappyNumbers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		happyNumbers(uint(i))
	}
}

func happy(a uint) (r uint) {
	for a > 0 {
		b := a % 10
		r += b * b
		a /= 10
	}
	return r
}

func happyNumbers(a uint) bool {
	b := []uint{a}
	for i := 0; a > 1 && i < 127; i++ {
		a = happy(a)
		for _, j := range b {
			if j == a {
				return false
			}
		}
		b = append(b, a)
	}
	return a == 1
}
