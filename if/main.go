package main

import "fmt"

func main() {
	counter := 42

	// example 1 (if)
	if counter > 10 {
		fmt.Println(counter, " is greater then 10")
	}

	// example 2 (if, else if, else)
	if counter > 100 {
		fmt.Println(counter, " is greater then 100")
	} else if counter > 50 {
		fmt.Println(counter, " is greater then 50")
	} else {
		fmt.Println(counter, " not greater then 100 or 50")
	}

	// example 3 (if with initialisation)
	if err := check(counter); err != 0 {
		// err is only available inside the if statement
		fmt.Println("ERROR: ", err)
	}
}

func check(n int) int {
	if n > 100 {
		return 0
	}
	return 1
}
