package main

import "fmt"

import "package/math"

func main() {
	xs := []float64{1, 2, 3, 12.2}
	Res := math.Average(xs)
	fmt.Println(Res)
}
