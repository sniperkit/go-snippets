package main

import "testing"

func TestAlter(t *testing.T) {
	tests := map[int]int{4: 1, 17: 6, 100: 292}
	for k, v := range tests {
		if r := alterMap(k, 50); r != v {
			t.Errorf("failed: alterMap %d, 50 is %d, got %d",
				k, v, r)
		}
	}
	for k, v := range tests {
		if r := alter(k, 4); r != v {
			t.Errorf("failed: alter %d, 4 is %d, got %d",
				k, v, r)
		}
	}
}

func BenchmarkAlterMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		alterMap(i%100+1, 50)
	}
}

func BenchmarkAlter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		alter(i%100+1, 4)
	}
}

var (
	mMap map[int]int
	m    []int
)

func init() {
	mMap = map[int]int{50: 25, 25: 10, 10: 5, 5: 1}
	m = []int{1, 5, 10, 25, 50}
}

func alterMap(n, c int) int {
	if n < 5 || c == 1 {
		return 1
	}
	if c > n {
		return alterMap(n, mMap[c])
	}
	return alterMap(n-c, c) + alterMap(n, mMap[c])
}

func alter(n, c int) int {
	if n < m[1] || c == 0 {
		return 1
	}
	if m[c] > n {
		return alter(n, c-1)
	}
	return alter(n-m[c], c) + alter(n, c-1)
}
