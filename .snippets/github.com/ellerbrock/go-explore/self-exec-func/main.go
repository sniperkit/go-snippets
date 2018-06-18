package main

import "fmt"

func main() {
	func() {
		fmt.Println("self executing function")
	}()
}
