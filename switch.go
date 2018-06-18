package main

import "fmt"
import "time"

func main() {
	switch time.Now().Weekday() {
	case time.Saturday, time.Sunday:
		fmt.Println("It's the weekend")
	default:
		fmt.Println("It's a weekday")
	}

	whatAmI := func(i interface{}) {
		switch t := i.(type) {
		case bool:
			fmt.Println("It's a bool")
		case int:
			fmt.Println("It's a int")
		default:
			fmt.Println("Don't know type", t)
		}
	}

	whatAmI(true)
	whatAmI(1)
	whatAmI("hey")
}
