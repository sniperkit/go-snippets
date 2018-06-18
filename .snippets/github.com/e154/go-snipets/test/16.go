package main

import "fmt"

// 16. Каким образом получить доступ к полю X типа Point из переменной circle ?

// Ответы:
// circle.Point.X

type Point struct {
	X, Y float64
}

type Circle struct {
	Point
	X float64
}

func main() {

	circle := &Circle{}
	circle.X = 12.2

	fmt.Println(circle.X)
	fmt.Println(circle.Point.X)
}

