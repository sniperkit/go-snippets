package main

import (
	"fmt"
)

func trace(s string) { fmt.Println("entering:", s) }

func untrace(s string) { fmt.Println("leaving:", s) }

func test() {
	trace("test")
	untrace("test")
	fmt.Println("test format")
}

func main() {
	trace("main")
	untrace("main")
	test()
}
