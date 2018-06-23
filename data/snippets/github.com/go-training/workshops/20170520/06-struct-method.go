package main

import (
	"fmt"

	"github.com/go-training/workshops/20170520/car"
)

func main() {
	car1 := car.Car{
		Name:  "lexus",
		Color: "white",
		Price: 800,
	}

	fmt.Println("car1", car1.Name)
	fmt.Println("car1", car1.Color)
	car1.UpdateColor("yellow")
	fmt.Println("car1", car1.Color)

	car2 := car.New("test", "red")

	fmt.Println("car2", car2.Name)
	fmt.Println("car2", car2.Color)
	car2.UpdateColor("yellow")
	fmt.Println("car2", car2.Color)
}
