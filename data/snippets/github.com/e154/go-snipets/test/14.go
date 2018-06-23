package main

import "fmt"

// 14. В какой строке кода ошибка?

// Ответы:
// ошибок нет

type Y int

func (y Y) meth()  {
	fmt.Println(y)
}

func (y *Y) meth2()  {
	fmt.Println(y)
}

func main() {

	y := Y(1)
	yp := new(Y)

	fmt.Println("Y(1)")
	y.meth()
	y.meth2()

	fmt.Println("new(Y)")
	yp.meth()
	yp.meth2()
}

