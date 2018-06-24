package main

import "testing"

func TestIsPalindrome(t *testing.T) {
	tests := map[int]bool{
		1:        true,
		2:        true,
		10:       false,
		11:       true,
		347:      false,
		12621:    true,
		3490943:  true,
		12345670: false}
	for k, v := range tests {
		if r := isPalindrome(k); r != v {
			t.Errorf("failed: isPalindrome %d is %t, got %t",
				k, v, r)
		}
		if r := palindromeSlice(k); r != v {
			t.Errorf("failed: palindromeSlice %d is %t, got %t",
				k, v, r)
		}
	}
}

func BenchmarkPalindromeSlice(b *testing.B) {
	for i := 1; i < b.N; i++ {
		palindromeSlice(i)
	}
}

func BenchmarkIsPalindrome(b *testing.B) {
	for i := 1; i < b.N; i++ {
		isPalindrome(i)
	}
}

var t []int

func init() {
	t = make([]int, 10)
}

func palindromeSlice(a int) bool {
	y := 0
	for c := a; c > 0; c /= 10 {
		t[y] = c % 10
		y++
	}
	for x := 0; ; x, y = x+1, y-1 {
		if t[x] != t[y-1] {
			return false
		}
		if y-x < 3 {
			return true
		}
	}
}

func isPalindrome(a int) bool {
	if a%10 == 0 {
		return a == 0
	}
	var rev int
	for ; a > rev; a /= 10 {
		rev = 10*rev + a%10
		if a == rev {
			return true
		}
	}
	return rev == a
}
