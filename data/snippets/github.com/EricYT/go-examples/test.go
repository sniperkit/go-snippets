package main

import "fmt"

func main() {
	var x string
	x = "hello world"
	fmt.Println(x)
	test_name()
}

func init() {
	var cal int32
	cal = 12
	fmt.Println(cal * cal)
}

func test() {
	cal := "string"
	var calValue = 123
	fmt.Println(cal)
	fmt.Println(calValue)
}

func test_name() {
	name := "Max"
	fmt.Println(name)
	name = "Little"
	fmt.Println(name)
	const constString string = "const string"
	fmt.Println(constString)
}
