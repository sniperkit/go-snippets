package main

import "fmt"
import (
	"runtime"
	"time"
)

func main() {
	fmt.Println("Go runs on ")
	switch os := runtime.GOOS; os{
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	default:
		// freebsd, openbsd
		// plan9, windows
		fmt.Printf("%s.", os)
	}

	//////////////////
	//Switch evaluation order
	//////////////////
	fmt.Println("")
	fmt.Println("Switch evaluation order")
	fmt.Println("When's Saturday?")
	today := time.Now().Weekday()
	switch time.Saturday{
	case today + 0:
		fmt.Println("Today.")
	case today + 1:
		fmt.Println("Tomorrow.")
	case today + 2:
		fmt.Println("In two days.")
	default:
		fmt.Println("Too far away.")
	}

	//////////////////
	//Switch with no condition
	//////////////////
	fmt.Println("")
	fmt.Println("Switch with no condition")
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("Good morning!")
	case t.Hour() < 17:
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}
}
