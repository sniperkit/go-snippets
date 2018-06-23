package main

import (
	"fmt"

	"github.com/go-training/workshops/20170520/lexus"
	"github.com/go-training/workshops/20170520/toyota"
)

type cart interface {
	Process(...int)
	Total() float32
}

func Output(c cart, price ...int) {
	c.Process(price...)
	fmt.Println(c.Total())
}

func main() {

	car1 := toyota.NewCart(1000)
	car2 := lexus.NewCart(2000)

	Output(car1, []int{100, 200, 300}...)
	Output(car2, []int{100, 400, 600}...)
}
