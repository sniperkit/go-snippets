package main

import (
	"fmt"
)

const (
	s = 19
	t = 1000
)

func zsum(a int) (r int) {
	for a > 0 {
		r, a = r+a%10, a/10
	}
	return r
}

func contains(a []int, b int) bool {
	for _, i := range a {
		if i == b {
			return true
		}
	}
	return false
}

func main() {
	var k int
	h := []int{0}
	for i := 0; i < len(h); i++ {
		c := h[i]
		if zsum(c+t) <= s && !contains(h, c+t) {
			h = append(h, c+t)
		}
		if zsum(c+1) <= s && !contains(h, c+1) {
			h = append(h, c+1)
		}
		if c >= t {
			k++
		}
	}
	fmt.Println(k*4 + 1)
}
