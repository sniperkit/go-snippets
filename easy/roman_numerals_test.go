package main

import "testing"

func TestRoman(t *testing.T) {
	for k, v := range map[uint]string{
		159:  "CLIX",
		296:  "CCXCVI",
		3992: "MMMCMXCII",
		1:    "I",
		3999: "MMMCMXCIX",
		347:  "CCCXLVII"} {
		if r := roman(k); r != v {
			t.Errorf("failed: roman %d is %s, got %s",
				k, v, r)
		}
	}
}

var (
	ronum []uint
	rostr []string
)

func init() {
	ronum, rostr = []uint{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1},
		[]string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
}

func roman(a uint) (r string) {
	for i := 0; a > 0; {
		if a >= ronum[i] {
			r += rostr[i]
			a -= ronum[i]
		} else {
			i++
		}
	}
	return r
}
