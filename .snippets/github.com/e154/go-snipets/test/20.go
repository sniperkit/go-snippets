package main

import "fmt"

// 20. Какой метод должет содержать тип что бы удовлетворять интерфейсу error?

// Ответы:
// Error() string

type Error string

func (t Error) Error() string {
	return string(t)
}

func SendError(e error) {
	fmt.Println("error:", e)
}

func main() {
	err := Error("custom error text")
	SendError(err)
}

