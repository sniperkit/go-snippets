package main

import "testing"

type tuple struct {
	y, x int
}

func TestPowi(t *testing.T) {
	for k, v := range map[tuple]int{
		tuple{9, 19}: 1350851717672992089,
		tuple{0, 8}:  0,
		tuple{0, 0}:  1} {
		if r := powi(k.y, k.x); r != v {
			t.Errorf("failed: %d^%d = %d, got %d",
				k.y, k.x, v, r)
		}
	}
}

func TestArmstrong(t *testing.T) {
	for k, v := range map[int]bool{
		6:   true,
		153: true,
		351: false,
		0:   true} {
		if r := armstrong(k); r != v {
			t.Errorf("failed: armstrong %d is %t, got %t",
				k, v, r)
		}
	}
}

func BenchmarkPowi(b *testing.B) {
	for i := 0; i < b.N; i++ {
		powi(i%10, i%20)
	}
}

func BenchmarkPowiNaive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		powiNaive(i%10, i%20)
	}
}

func BenchmarkArmstrong(b *testing.B) {
	for i := 0; i < b.N; i++ {
		armstrong(i)
	}
}

func powi(y, x int) int {
	ret := 1
	for x > 0 {
		if x&1 > 0 {
			ret *= y
		}
		y *= y
		x >>= 1
	}
	return ret
}

func powiNaive(a, b int) (ret int) {
	ret = 1
	for i := 0; i < b; i++ {
		ret *= a
	}
	return ret
}

func armstrong(n int) bool {
	var e, t int
	for a := n; a > 0; a /= 10 {
		e++
	}
	for a := n; a > 0; a /= 10 {
		t += powi(a%10, e)
	}
	return n == t
}
