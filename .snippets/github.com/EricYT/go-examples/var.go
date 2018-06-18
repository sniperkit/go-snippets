package main

import (
	"fmt"
)

const (
	constValue  = 123
	constString = "const string"
)

func test() string {
	return "test"
}

func main() {
	fmt.Println(constValue, constString)
	fmt.Println(test())

  // varibles
  var intA int = 900
  fmt.Println(intA)
}
