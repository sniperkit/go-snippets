package main

import "fmt"

func main() {
	defer func() {
		str := recover()
		fmt.Println(str)
	}()
	panic("shit, crash")
	panic("shit, crash111")
}
