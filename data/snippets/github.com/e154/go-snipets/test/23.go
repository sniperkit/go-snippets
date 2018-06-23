package main

import "fmt"

// 23. Что произойдет в результате выполнения следующего кода?

// Ответы:
// напечатает It`s fine, затем вызовет панику


func keepCallmAndPanic(msg string)  {
	panic(msg)
}

func main() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err.(error).Error())
		}
	}()

	go keepCallmAndPanic("Aargh!")

	fmt.Println("It`s fine")
}

