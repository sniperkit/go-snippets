package bar

import (
	"flag"
	"fmt"
	"runtime"
	"time"
)

var (
	bar *string = flag.String("bar", "default bar", "sub package flag init")
)

func int() {
	go Display()
}

func Display() string {
	for !flag.Parsed() {
		// Defer execution of this goroutine
		runtime.Gosched()

		time.Sleep(time.Second * 1)
	}
	fmt.Println("bar package bar:", *bar)
	return *bar
}
