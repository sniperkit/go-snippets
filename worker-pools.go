package main

import (
	"fmt"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, " processing job", j)
		time.Sleep(time.Second)
		results <- j * 2
	}
}
func main() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	// In order to use our pool of workers we need to send them work and collect their results. We make 2 channels for this.

	// This starts up 3 workers, initially blocked because there are no jobs yet.
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	//Here we send 9 jobs and then close that channel to indicate that’s all the work we have.
	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs)

	//Finally we collect all the results of the work
	for a := 1; a <= 9; a++ {
		fmt.Println(<-results)
	}
}
