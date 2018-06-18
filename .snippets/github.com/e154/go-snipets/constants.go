package main

import "fmt"

const (
	Pi = 3.14
)

func needInt(x int) int {return x*10 + 1}
func needFloat(x float64) float64 {return x * 0.1}

func main() {
	const (
		World = "World"
		True = true
	)

	fmt.Println("Hello", World)
	fmt.Println("Happy", Pi, "Day")
	fmt.Println("Go rules?", True)

	//Numeric Constants
	const (
		Big		= 1 << 100
		Small	= Big >> 99
	)

	fmt.Println(needInt(Small))
	fmt.Println(needFloat(Small))
	fmt.Println(needFloat(Big))
}
