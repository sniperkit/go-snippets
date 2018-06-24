package main

import "testing"

func TestIsPalindrome(t *testing.T) {
	for k, v := range map[int]bool{
		0:        true,
		2:        true,
		10:       false,
		11:       true,
		347:      false,
		12621:    true,
		3490943:  true,
		12345670: false} {
		if r := isPalindrome(k); r != v {
			t.Errorf("failed: isPalindrome %d is %t, got %t",
				k, v, r)
		}
	}
}

func BenchmarkIsPalindrome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isPalindrome(i)
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
