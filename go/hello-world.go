package main

import (
	"os"
	"fmt"
)

func main() {
	var args = os.Args
	var argsWithoutProg = os.Args[1:]

	fmt.Println("Hello!")
	fmt.Println(args)
	fmt.Println(argsWithoutProg)
}
