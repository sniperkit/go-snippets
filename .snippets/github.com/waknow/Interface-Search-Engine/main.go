package main

import (
	"fmt"

	"interface/lib"
)

func main() {
	path := "test.go"

	fmt.Println("scan start")
	var tokens lib.Tokens
	tokens.Scan(path)

	interfaces := tokens.Types.GetInterfaces()
	structs := tokens.Types.GetStructs()
	funcations := tokens.Funcations
	fmt.Println("interface\n", interfaces.String())
	fmt.Println("struct\n", structs.String())
	fmt.Println("func\n", funcations.String())
}
