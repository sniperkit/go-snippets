package main

import (
	"fmt"

	"github.com/shijuvar/go-recipes/binarypkg"
)

func main() {
	str := "Golang"
	// Convert to upper case
	fmt.Println("To Upper Case:", binarypkg.ToUpperCase(str))

	// Convert to lower case
	fmt.Println("To Lower Case:", binarypkg.ToLowerCase(str))

}
