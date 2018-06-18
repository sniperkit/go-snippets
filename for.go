package main

import (
	"fmt"
	_ "math"
)

func main() {
	var i int = 1
	for i <= 10 {
		fmt.Println(i)
		i += 1
	}
	fmt.Println("------------------->")
	for i = 1; i < 10; i++ {
		if i%2 == 0 {
			fmt.Println("even")
		} else {
			fmt.Println("odd")
			//break
		}
	}
	fmt.Println("------------------->")
	i = 1
	for i = 1; i < 10; i++ {
		switch i {
		case 0:
			fmt.Println(i, "first")
		case 1:
			if i > 0 {
				fmt.Println("scope")
			}
		default:
			fmt.Println("other")
		}
	}
	fmt.Println("------------------->")
	var arr [5]int
	for i = 0; i < 5; i++ {
		fmt.Println(arr[i])
	}
	fmt.Println("------------------->")
	var float float64
	float = 12.3
	//	fmt.Println(float / len(arr))
	fmt.Println(float / float64(len(arr)))

}
