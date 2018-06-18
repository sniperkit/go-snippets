package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

func printText() {
	for i := range []int{2, 3, 4} {
		fmt.Println("this is the output we want to catch #", i)
		time.Sleep(1 * time.Second)
	}
}

func main() {

	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printText()

	outC := make(chan string)

	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	// reading our temp stdout
	fmt.Println("previous output: ", out)

}
