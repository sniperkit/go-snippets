package main

import (
	"fmt"
	"time"
)

func main() {
	jobs := make(chan int, 5)
	done := make(chan bool)
	go func() {
		for {
			j, more := <-jobs
			if more {
				fmt.Println("received job", j)
			} else {
				fmt.Println("received all jobs")
				done <- true // chan 要在不同的goroutine中
				return
			}
		}
	}()

	for j := 1; j <= 3; j++ {
		jobs <- j
		fmt.Println("sent job", j)
		time.Sleep(time.Second * 1)
	}
	//This sends 3 jobs to the worker over the jobs channel, then closes it.
	close(jobs)
	fmt.Println("sent all jobs")
	<-done // chan 要在不同的goroutine中
}
