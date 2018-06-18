package main

import (
	. "fmt"
)

// 05. Возможно ли опрежделить функцию принимающую не обязательные параметры?

// Ответы:
// func(x int, y ...int)

func main() {

	test5(1)
	test5(1, 2)
	test5(1, 2, 3)

}

func test5(x int, y ...int) {
	Println(x, y)
}