package test

import "fmt"

func Reverse(s string) string {
	runs := []rune(s)
	for i, j := 0, len(runs)-1; i < len(runs)/2; i, j = i+1, j-1 {
		runs[i], runs[j] = runs[j], runs[i]
	}
	return string(runs)
}

func test() {
	fmt.Println("vim-go")
}
