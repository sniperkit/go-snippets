package main

import (
	_ "fmt"
	"math/rand"
)

type Order struct {
	Id    string
	Value string
	Cur   string
}

func main() {
	orders := []Order{}
	for i := 0; i < 100000; i++ {
		orders = append(orders, Order{Id: string(rand.Intn(1000)), Value: string(rand.Intn(1000)), Cur: string(rand.Intn(1000))})
	}

	i := 0
	for _, o := range orders {
		for _, o2 := range orders {
			i++
			o.Id = o2.Id
		}
	}

}
