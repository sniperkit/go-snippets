package main

import (
	"fmt"
)

func main() {

	var a string = "hello world"

	{
		a := "Hello china"
		fmt.Println(a)
	}

	{
		b := "Hello beijing"
		fmt.Println(b)
	}
	//  Blow error
	//	fmt.Println(b)
	fmt.Println(a)
}
