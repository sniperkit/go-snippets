package main

import (
	"fmt"
)

var x interface{}

func main() {
	x = 5
	i, ok := x.(int)
	var cells []interface{} = []interface{}{"x", 124, true}
	if ok {
		fmt.Printf("x holds the integer value %d\n\n", i)
	}

	x = "test"
	var s = x.(string)
	fmt.Printf("The string value of x is now %s\n\n", s)

	x = 5
	fmt.Printf("The runtime value of x is now %v\n\n", x)

	for i := 0; i < len(cells); i++ {
		fmt.Printf("Item %d ", i)
		switch cells[i].(type) {
		case int:
			fmt.Printf("int :  %d\n", cells[i].(int))
			break
		case string:
			fmt.Printf("string : %s\n", cells[i].(string))
			break
		case bool:
			fmt.Printf("bool : %t\n", cells[i].(bool))
			break
		}
	}
}
