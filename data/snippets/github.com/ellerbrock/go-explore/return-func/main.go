package main

import "fmt"

func main() {
	fmt.Println(wrapper("Hello")())
}

func wrapper(msg string) func() string {
	return func() string {
		return msg + " World!"
	}
}
