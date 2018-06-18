package main

import (
	"fmt"
)

// 06. Что происходит при передачи слайса в функцию

// Ответы:
// Слайс копируется но ссылается на исходные данные. Функция может изменить данные но не исходный слайс

func main() {

	a := []int{1,2,3}
	b := []int{4,5,6}

	test61(&a)
	test62(b)

	fmt.Println("a",a)
	fmt.Println("b",b)

}

func test61(s *[]int) {
	*s = append(*s, 2)
	//fmt.Println("a", *s)
}

func test62(s []int) {
	s[0] = 7
	s = append(s, 2)
	//fmt.Println("b",s)
}