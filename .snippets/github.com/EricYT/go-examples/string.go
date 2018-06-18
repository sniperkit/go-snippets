package main

import "fmt"

var str string = "Not my business "

func main() {

	// go scope
	{
		var str string = `
	"Middle the strings"
                    world`

		if str == "helloworld" {
			fmt.Println("same line")
		}
		fmt.Println("different line")
		fmt.Println(str)
	}

	fmt.Println(str)

}
