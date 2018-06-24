package main

import (
	"fmt"
	"strings"
	"testing"
)

type tuple struct {
	a, b int
}

func TestIsMagic(t *testing.T) {
	m := []int{1, 9, 35, 37, 174, 1267, 3562, 6712, 6392, 9263, 9627}
	n := []int{10, 12, 18, 175, 1624, 2715, 3261, 6372, 7216, 9876}
	for _, i := range m {
		if !isMagic(i) {
			t.Errorf("failed: isMagic %d is true, got false", i)
		}
	}
	for _, i := range n {
		if isMagic(i) {
			t.Errorf("failed: isMagic %d is false, got true", i)
		}
	}
}

func TestMagicNumbers(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{10, 100}:    "13 15 17 19 31 35 37 39 51 53 57 59 71 73 75 79 91 93 95 97",
		tuple{8382, 8841}: "-1",
		tuple{1, 10000}:   "1 2 3 4 5 6 7 8 9 13 15 17 19 31 35 37 39 51 53 57 59 71 73 75 79 91 93 95 97 147 174 258 285 417 471 528 582 714 741 825 852 1263 1267 1623 1627 2316 2356 2396 2631 2635 2639 2671 2675 2679 2716 2756 2796 3126 3162 3526 3562 3926 3962 5263 5267 5623 5627 6231 6235 6239 6271 6275 6279 6312 6352 6392 6712 6752 6792 7126 7162 7526 7562 7926 7962 9263 9267 9623 9627"} {
		if r := magicNumbers(k.a, k.b); r != v {
			t.Errorf("failed: magicNumbers %d, %d is %s, got %s",
				k.a, k.b, v, r)
		}
	}
}

func BenchmarkIsMagic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isMagic(i%10000 + 1)
	}
}

func BenchmarkMagicNumbers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n := i%10000 + 1
		magicNumbers(n, n+(i/10000)%(10001-n))
	}
}

func isMagic(a int) bool {
	var (
		dig, r uint
		ns     []uint
	)
	for a > 0 {
		r = uint(a % 10)
		if r == 0 || dig&(1<<r) > 0 {
			return false
		}
		dig |= 1 << r
		ns = append(ns, r)
		a /= 10
	}
	dig, r = 0, 0
	for _ = range ns {
		r = (r + ns[(uint(len(ns))-1-r)]) % uint(len(ns))
		if dig&(1<<r) > 0 {
			return false
		}
		dig |= 1 << r
	}
	return r == 0
}

var magic []int

func init() {
	for i := 1; i <= 9876; i++ {
		if isMagic(i) {
			magic = append(magic, i)
		}
	}
}

func magicNumbers(a, b int) string {
	var r []string
	for i := 0; i < len(magic) && magic[i] <= b; i++ {
		if magic[i] >= a {
			r = append(r, fmt.Sprint(magic[i]))
		}
	}
	if len(r) == 0 {
		return "-1"
	}
	return strings.Join(r, " ")
}
