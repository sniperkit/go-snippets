package main

import "fmt"

func main() {
	justDoIt([]int{10, 20, 30}, func(n int) {
		fmt.Println(n)
	})
}

func justDoIt(numbers []int, callback func(int)) {
	for _, v := range numbers {
		callback(v)
	}
}
