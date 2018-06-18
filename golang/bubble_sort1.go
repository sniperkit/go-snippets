package main

import "fmt"

func main() {
	a := []int{5, -1, 0, 12, 3, 5}
	fmt.Printf("unsorted %v\n", a)
	n := bubblesort(5, -1, 0, 12, 3, 5)
	fmt.Printf("sorted %v\n", n)
}

func bubblesort(n ...int) (a []int) {
	for i := 0; i < len(n)-1; i++ {
		for j := i + 1; j < len(n); j++ {
			if n[j] < n[i] {
				n[i], n[j] = n[j], n[i]
			}
		}
	}
	return n
}
