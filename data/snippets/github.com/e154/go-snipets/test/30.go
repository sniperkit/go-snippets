package main

import "fmt"

// 30. Что будет содержаться в слайсе a2 после выполнения copy?

// Ответы:
// 1,2

func main() {
	a1 := []int{1,2,3,4}
	a2 := []int{5,6}

	copy(a2, a1)

	fmt.Println(a2)
}

