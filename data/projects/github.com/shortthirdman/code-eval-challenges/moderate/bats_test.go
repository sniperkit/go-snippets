package main

import (
	"fmt"
	"testing"
)

func TestBats(t *testing.T) {
	if r := bats(22, 2, 2, []string{"9", "11"}); r != 3 {
		t.Errorf("failed: bats 22 2 2 9 11 is 3, got %d",
			r)
	}
	if r := bats(835, 125, 1, []string{"113"}); r != 5 {
		t.Errorf("failed: bats 835 125 1 113 is 5, got %d",
			r)
	}
	if r := bats(47, 5, 0, []string{}); r != 8 {
		t.Errorf("failed: bats 475 5 0 is 8, got %d",
			r)
	}
}

func bats(l, d, n int, s []string) (c int) {
	var t int
	tx := 6 - d
	for i := 6; i <= l-6; i += d {
		if i > tx-d {
			i = tx
			if t == n {
				tx = l - 6 + d
			} else {
				fmt.Sscanf(s[t], "%d", &tx)
				t++
			}
		} else {
			c++
		}
	}
	return c
}
