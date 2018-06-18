package main

import "fmt"
import "time"

func main() {
	var today int64
	today = 5

	fmt.Print("Today is ", today, " as \n")
	switch today {
	case 1:
		fmt.Print("1\n")
	case 2:
		fmt.Print("2\n")
	default:
		fmt.Print("Other\n")
	}

	switch time.Now().Weekday() {
	case time.Friday:
		fmt.Println("today is Friday")
	default:
		fmt.Println("I don't know")
	}

	switch {
	case time.Now().Hour() < 20:
		fmt.Println("not to late")
	default:
		fmt.Println("time to go home")
	}
}
