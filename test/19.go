package main

import (
	"strconv"
	"errors"
	"fmt"
)

// 19. Что вернет следующая функция при вызове findSquare("5")

// Ответы:
// 0 <nil>

var squares map[int]int = map[int]int{0:0,1:1,2:4,3:9}

func findSquare(n string) (res int, err error) {
	if n, err := strconv.Atoi(n); err == nil {
		fmt.Println(1)
		if n, ok := squares[n]; !ok {
			fmt.Println(2)
			err = errors.New("not found")
		} else {
			res = n
		}
	}

	fmt.Println(3)
	return
}

func main() {
	fmt.Println(findSquare("5"))
}

