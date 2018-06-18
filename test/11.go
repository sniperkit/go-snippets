package main

import (
	"fmt"
)

// 11. Какие выражения считаются "истинными" (приводят к выполнению условия if)?

// Ответы:
// true

func main() {

	//if "0" {
	//	fmt.Print("'0'")
	//}

	//if 0 {
	//	fmt.Print(0)
	//}

	//if []int{} {
	//	fmt.Print("[]int{}")
	//}

	//if 1 {
	//	fmt.Print(1)
	//}

	//if nil {
	//	fmt.Print(nil)
	//}

	//if fmt.Printf {
	//	fmt.Print("fmt.Printf")
	//}

	if true {
		fmt.Print("true")
	}
}
