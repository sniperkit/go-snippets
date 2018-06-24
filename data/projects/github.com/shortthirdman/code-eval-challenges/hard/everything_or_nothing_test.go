package main

import (
	"strings"
	"testing"
)

func TestEverything(t *testing.T) {
	for k, v := range map[string]bool{
		"user_1=>file_1=>read user_2=>file_2=>write":               true,
		"user_1=>file_1=>grant=>read=>user_4 user_4=>file_1=>read": true,
		"user_4=>file_1=>read":                                     false} {
		if r := everything(k); r != v {
			t.Errorf("failed: everything %s is %t, got %t",
				k, v, r)
		}
	}
}

const (
	read  = 4
	write = 2
	grant = 1
)

func d(s uint8) uint8 {
	return s - 49
}

func everything(q string) bool {
	uperm := [][]int{[]int{7, 3, 0}, []int{6, 2, 4}, []int{5, 1, 5},
		[]int{3, 7, 1}, []int{6, 0, 2}, []int{4, 2, 6}}
	for _, i := range strings.Fields(q) {
		w, f := d(i[5]), d(i[13])
		switch i[16] {
		case 'r':
			if uperm[w][f]&read == 0 {
				return false
			}
		case 'w':
			if uperm[w][f]&write == 0 {
				return false
			}
		case 'g':
			if uperm[w][f]&grant == 0 {
				return false
			}
			switch i[23] {
			case 'r':
				uperm[d(i[34])][f] |= read
			case 'w':
				uperm[d(i[35])][f] |= write
			case 'g':
				uperm[d(i[35])][f] |= grant
			}
		}
	}
	return true
}
