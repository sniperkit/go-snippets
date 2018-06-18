package main

import "fmt"
import "time"

func main() {
	fmt.Println("hello,world")
	done := make(chan bool)

	go func() {
		time.Sleep(600 * time.Second)
		done <- true
	}()

	<-done
}
