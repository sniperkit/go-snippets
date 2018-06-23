package main

import "fmt"

type sayInterface interface {
	say()
	talk()
}

type print1 struct {
	name string
}

func (pPtr *print1) say() {
	fmt.Println("print1 first ", pPtr.name)
}

func (pPtr *print1) talk() {
	fmt.Println("print1 talk ", pPtr.name)
}

type print2 struct {
	name string
}

func (pPtr *print2) say() {
	fmt.Println("print2 second ", pPtr.name)
}

func (pPtr *print2) talk() {
	fmt.Println("print2 talk ", pPtr.name)
}

func say(print ...sayInterface) {
	for _, p := range print {
		p.say()
		p.talk()
	}
}

func sayAny(any ...interface{}) {
	for _, arg := range any {
		switch arg.(type) {
		case int:
			fmt.Println("Int")
		case float64:
			fmt.Println("Float")
		case string:
			fmt.Println("String")
		case byte:
			fmt.Println("byte")
		default:
			fmt.Println("shit")
		}
	}
}

func main() {
	p1 := new(print1)
	p1.name = "eric"

	p2 := new(print2)
	p2.name = "yutao"

	say(p1, p2)

	var sayInters sayInterface = new(print1)
	sayInters.say()

	if _, ok := sayInters.(sayInterface); ok {
		fmt.Println("-------> Interface")
	}

  sayAny(123, p1, 123.2, "string")

}
