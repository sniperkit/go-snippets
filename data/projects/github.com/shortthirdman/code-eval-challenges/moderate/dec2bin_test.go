package main

import (
	"fmt"
	"testing"
)

func TestDec2bin(t *testing.T) {
	for k, v := range map[uint]string{
		256:     "100000000",
		10:      "1010",
		0:       "0",
		67:      "1000011",
		347:     "101011011",
		65535:   "1111111111111111",
		1657934: "110010100110001001110"} {
		if r := dec2bin(k); r != v {
			t.Errorf("failed: dec2bin %d is %s, got %s",
				k, v, r)
		}
	}
}

func dec2bin(a uint) string {
	return fmt.Sprintf("%b", a)
}
