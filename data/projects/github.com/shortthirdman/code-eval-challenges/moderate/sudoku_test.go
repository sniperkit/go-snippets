package main

import (
	"fmt"
	"strings"
	"testing"
)

type tuple struct {
	n int
	q string
}

func TestSudoku(t *testing.T) {
	for k, v := range map[tuple]bool{
		tuple{4, "1,4,2,3,2,3,1,4,4,2,3,1,3,1,4,2"}: true,
		tuple{4, "2,1,3,2,3,2,1,4,1,4,2,3,2,3,4,1"}: false} {
		if r := sudoku(k.n, k.q); r != v {
			t.Errorf("failed: sudoku %d %s is %t, got %t",
				k.n, k.q, v, r)
		}
	}
}

func sq(a int) (ret int) {
	switch a {
	case 4:
		ret = 2
	case 9:
		ret = 3
	}
	return ret
}

func sudoku(n int, q string) bool {
	var d int
	t := strings.Split(q, ",")
	csqu, col := make([]int, sq(n)), make([]int, n)
	for i := 0; i < n; i++ {
		var crow int
		for j := 0; j < n; j++ {
			fmt.Sscan(t[i*n+j], &d)
			d = 1 << uint(d-1)
			crow |= d
			csqu[j/sq(n)] |= d
			col[j] |= d
		}
		d = (1 << uint(n)) - 1
		if crow != d {
			return false
		}
		if ((i + 1) % sq(n)) == 0 {
			for j := 0; j < sq(n); j++ {
				if csqu[j] != d {
					return false
				}
				csqu[j] = 0
			}
		}
	}
	for i := 0; i < n; i++ {
		if col[i] != (1<<uint(n))-1 {
			return false
		}
	}
	return true
}
