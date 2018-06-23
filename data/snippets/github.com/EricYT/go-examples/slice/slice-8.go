package main

import (
	"fmt"
	"strings"
)

func AppendByte(slice []byte, data ...byte) []byte {
	m := len(slice)
	n := m + len(data)
	if n > cap(slice) {
		// allocate double what's need, for future growth.
		newSlice := make([]byte, (n+1)*2)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0:n]
	copy(slice[m:n], data)
	return slice

}

func Filter(s []int, fn func(int) bool) []int {
	// p == nil
	var p []int
	for _, v := range s {
		if fn(v) {
			p = append(p, v)
		}
	}
	return p
}

func main() {

	s1 := []int{1, 2, 3, 4}
	s2 := make([]int, len(s1)-2)

	// a deep copy
	copy(s2, s1)

	fmt.Println("s1:", s1)
	fmt.Println("s2:", s2)

	s2[0] = 0
	fmt.Println("s1:", s1)
	fmt.Println("s2:", s2)
	fmt.Println("res:", strings.Contains("a.b.c...", "..."))

}
