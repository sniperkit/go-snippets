package main

import "fmt"
import "strconv"

func main() {
	jobs := make(chan string)
	done := make(chan bool)

	go func() {
		for {
			j, more := <-jobs
			if more {
				fmt.Println("receive message ", j)
			} else {
				fmt.Println("receive all ")
				done <- true
				return
			}
		}
	}()

	for i := 0; i < 3; i++ {
		jobs <- strconv.Itoa(i)
	}
	close(jobs)

	fmt.Println(<-done)
}
