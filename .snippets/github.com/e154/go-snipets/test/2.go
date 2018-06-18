package main

import "fmt"

// 02. Выберите корректные варианты описания циклов

// Ответы:
// 2,3

func main() {

	i := 0
	for {
		fmt.Println(i)
		i++
		if i >= 10 {
			break
		}
	}

	i = 0
	for i < 10 {
		fmt.Println(i)
		i++
	}

}
