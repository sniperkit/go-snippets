package main

import "fmt"

func init() {
	fmt.Println("init")
}

func main() {
	defer end()
	fmt.Println("main")
}

func end() {
	fmt.Println("end")
}
