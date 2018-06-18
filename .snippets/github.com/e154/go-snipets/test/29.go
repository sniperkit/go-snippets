package main

import "fmt"

// 29. Выберите фрагмент кода который выполнится без ошибок

// Ответы:
//

func main() {
	//var x map[int]int
	//x[1] = 1

	var x [8]int
	x[1] = 1

	//var y []int
	//y[1] = 1

	var z map[int]int
	l := len(z)

	fmt.Println(l)
}

