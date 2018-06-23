package main

import (
	"fmt"
)

// 09. Что произойдет в результате выполнения кода?

// Ответы:
// ошибка на этапе компиляции

func main() {

	var i int8 = -1
	var j uint32 = 2

	fmt.Printf("%d\n", i + j)
}
