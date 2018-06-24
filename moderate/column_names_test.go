package main

import "testing"

func TestColumnNames(t *testing.T) {
	for k, v := range map[int]string{
		52:    "AZ",
		3702:  "ELJ",
		18278: "ZZZ"} {
		if r := columnNames(k); r != v {
			t.Errorf("failed: columnNames %d is %s, got %s",
				k, v, r)
		}
	}
}

func BenchmarkColumnNames(b *testing.B) {
	for i := 0; i < b.N; i++ {
		columnNames(i%18278 + 1)
	}
}

func columnNames(c int) (r string) {
	for c > 0 {
		c--
		r = string('A'+c%26) + r
		c /= 26
	}
	return r
}
