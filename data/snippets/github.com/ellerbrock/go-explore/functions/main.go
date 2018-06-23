package main

import "fmt"

func main() {
	numbers := []int{10, 20, 30, 40, 50}

	printArray(&numbers)
}

func printArray(numbers *[]int) {
	// fmt.Println(num)
	for index, value := range *numbers {
		fmt.Println("index: ", index, " | value: ", value)
	}
}
