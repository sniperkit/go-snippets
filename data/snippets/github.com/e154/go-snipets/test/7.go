package main

import (
	"fmt"
)

// 07. Что будет напечатано в результате работы следующего кода

// Ответы:
// ничего

func main() {


	var a interface{} = int8(10)

	switch a := a.(type) {
	case int8:
	case int16:
	case int32:
	case int64:
		fmt.Println(a)
	default:
		fmt.Println("unknown", a)
	}
}
