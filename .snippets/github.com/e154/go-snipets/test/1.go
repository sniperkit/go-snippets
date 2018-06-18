package main

import "fmt"

// 01. Какие синтаксические конструкции есть в GO?

// Ответы:
// Сокращенные операторы присвоения а += 3
// Множественное присвоение a,b = b,a

func main() {
	i := 0
	i += 3

	fmt.Println(i)

	y := 0
	i, y = 0,3

	fmt.Println(i,y)

}