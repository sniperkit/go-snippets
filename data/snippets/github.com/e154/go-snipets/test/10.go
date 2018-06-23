package main

import (
	"fmt"
)

// 10. Какие типы данных можно использовать в конструкции for range

// Ответы:
// slice, array, map, string, chan, <-chan

func main() {

	str := "string"
	for k, v := range str  {
		fmt.Println(k, v)
	}
}
