package main

import "fmt"

func main() {
	fmt.Println("Please input a number:")
	var input float64
	fmt.Scanf("%f", &input)
	var output float64
	output = input * 2
	fmt.Println("Output is ", output)
}
