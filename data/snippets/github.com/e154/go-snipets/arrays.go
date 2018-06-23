package main

import "fmt"

func removeDuplicates(a []int) []int {
	result := []int{}
	seen := map[int]int{}
	for _, val := range a {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = val
		}
	}
	return result
}

func main() {
	var a [2]string
	a[0] = "Hello"
	a[1] = "World"
	fmt.Println(a[0], a[1])
	fmt.Println(a)
	fmt.Println("array length", len(a))

	// Duplicates
	x := []int{9, 1, 2, 3, 5, 6, 9, 2, 3, 3, 6, 8, 1}
	x = removeDuplicates(x)
	fmt.Println(x)
}
