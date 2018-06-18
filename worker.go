package main

import (
	"fmt"
	"time"
)

func Worker(id int, jobs <-chan int, res chan<- int) {
	for i := range jobs {
		fmt.Println("Worker ", id, " job ", i)
		time.Sleep(time.Millisecond * 500)
		res <- i
	}
}

func main() {
	jobs := make(chan int, 100)
	res := make(chan int, 100)

	for i := 0; i < 3; i++ {
		go Worker(i, jobs, res)
	}

	for i := 0; i < 10; i++ {
		jobs <- i
	}
	close(jobs)

	for i := 0; i < 10; i++ {
		result := <-res
		fmt.Println("Worker receive res ", result)
	}
}
