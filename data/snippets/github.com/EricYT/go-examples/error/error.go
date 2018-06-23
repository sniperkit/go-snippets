package main

import "fmt"

type ErrorFoo struct {
	msg string
}

func (f ErrorFoo) Error() string {
	return f.msg
}

type ErrorBar struct {
	msg string
}

func (b ErrorBar) Error() string {
	return b.msg
}

func ReturnError() error {
	return &ErrorFoo{msg: "foo error message"}
}

func main() {
	fmt.Println("vim-go")

	err := ReturnError()

	if foo, ok := err.(*ErrorFoo); ok {
		fmt.Println("error message:", foo)
	}

	switch err.(type) {
	case *ErrorFoo:
		fmt.Println("this is a foo error")
	case *ErrorBar:
		fmt.Println("this is a bar error")
	default:
		fmt.Println("what the hell is this")
	}
}
