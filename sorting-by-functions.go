package main

import (
	"fmt"
	"sort"
)

type ByLength []string

func (s ByLength) Len() int {
	return len(s)
}

func (s ByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

func main() {
	fruits := []string{"peach", "banana", "kiwi"}
	fmt.Println(ByLength(fruits))
	sort.Sort(ByLength(fruits))
	fmt.Println(fruits)
}

// We implement sort.Interface - Len, Less, and Swap - on our type so we can use the sort package’s
// generic Sort function. Len and Swap will usually be similar across types and Less will hold the actual
// custom sorting logic. In our case we want to sort in order of increasing string length,
// so we use len(s[i]) and len(s[j]) here.
