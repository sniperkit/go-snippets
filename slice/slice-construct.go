package main

import "fmt"

func main() {
	var slice []string = append([]string(nil), "hello,world")
	fmt.Println(slice)
}
