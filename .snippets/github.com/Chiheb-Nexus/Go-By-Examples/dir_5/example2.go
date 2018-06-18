package main

import (
	"fmt"
	"strings"
)

type MyAlias string

func (s MyAlias) ToUpperCase() string {
	return strings.ToUpper(string(s))
}

func main() {
	fmt.Println(MyAlias("test").ToUpperCase())
}
