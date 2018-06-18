package main

import "fmt"

func vals() (int, int) {
	return 3, 6
}

func plus(a int, b int) int {
	return a + b
}

func plusPlus(a, b, c int) int {
	return a + b + c
}

func sum(nums ...int) {
	fmt.Println(nums)

	total := 0

	for _, num := range nums {
		total += num
	}

	fmt.Println("total:", total)
}

func main() {
	res := plus(1, 2)
	fmt.Println("1+2=", res)

	res = plusPlus(1, 2, 3)
	fmt.Println("1+2+3=", res)

	a, b := vals()
	fmt.Println(a, b)

	sum(1, 2)
	sum(1, 2, 3)

	nums := []int{1, 2, 3}
	sum(nums...)
}
