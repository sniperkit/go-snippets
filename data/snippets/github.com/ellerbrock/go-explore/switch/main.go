package main

import "fmt"

func main() {
	i := 42
	s := "a"
	b := true
	p := &i

	info(i)
	info(s)
	info(b)
	info(p)
}

func info(t interface{}) {
	switch t.(type) {
	case int:
		fmt.Println("int")
	case string:
		fmt.Println("string")
	case bool:
		fmt.Println("bool")
	default:
		fmt.Println("something else ...")
	}

}
