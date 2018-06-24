package main

import (
	"fmt"
	"strings"
	"testing"
)

type tuple struct {
	f, b, n uint
}

func TestFizzbuzz(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{3, 5, 21}:    "1 2 F 4 B F 7 8 F B 11 F 13 14 FB 16 17 F 19 B F",
		tuple{2, 7, 25}:    "1 F 3 F 5 F B F 9 F 11 F 13 FB 15 F 17 F 19 F B F 23 F 25",
		tuple{1, 4, 32}:    "F F F FB F F F FB F F F FB F F F FB F F F FB F F F FB F F F FB F F F FB",
		tuple{20, 20, 100}: "1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 FB 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 FB 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55 56 57 58 59 FB 61 62 63 64 65 66 67 68 69 70 71 72 73 74 75 76 77 78 79 FB 81 82 83 84 85 86 87 88 89 90 91 92 93 94 95 96 97 98 99 FB"} {
		if r := fizzbuzz(k.f, k.b, k.n); r != v {
			t.Errorf("failed: fizzbuzz %d %d %d is %s, got %s",
				k.f, k.b, k.n)
		}
	}
}

type node struct {
	value string
	next  *node
}

func (a node) isEmpty() bool {
	return len(a.value) == 0
}

func fizzbuzz(f, b, n uint) string {
	var r []string
	tail := node{value: "FB"}
	list := &tail
	for i, j := f-1, b-1; i > 0 || j > 0; i, j = (i+f-1)%f, (j+b-1)%b {
		var v string
		if i == 0 {
			v = "F"
		} else if j == 0 {
			v = "B"
		}
		list = &node{v, list}
	}
	tail.next = list

	for i := uint(1); i <= n; i++ {
		if list.isEmpty() {
			r = append(r, fmt.Sprint(i))
		} else {
			r = append(r, list.value)
		}
		list = list.next
	}
	return strings.Join(r, " ")
}
