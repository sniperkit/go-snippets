package main

import (
	"fmt"
	"runtime"
)

func main() {
	names := []string{"Kitty", "Slice", "Namao"}
	for _, name := range names {
		go func() {
			fmt.Printf("Hello %s\n", name)
		}()
		runtime.Gosched() // for goroutine is wait for the go func goroutine
	}
	// runtime.Gosched()  the result is: all is "Hello Namao" because for code is faster then go func()
}
