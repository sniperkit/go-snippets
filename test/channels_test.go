package test

import (
	"fmt"
	"testing"
	"time"
)

func TestChannelsPingPong(t *testing.T) {
	// 1. unbuffered channels
	pingChannel := make(chan string)
	pongChannel := make(chan string)
	// 2. buffered channel of 10 messages
	// pingChannel := make(chan string, 10)
	// pongChannel := make(chan string, 10)
	go ping(pingChannel)
	go pong(pongChannel)
	go print(pingChannel, pongChannel)
	var input string
	fmt.Scanln(&input)
}

func ping(pingChannel chan string) {
	for i := 0; i < 10; i++ {
		s := fmt.Sprintf("ping %d", i)
		fmt.Printf("send %s\n", s)
		pingChannel <- s
	}
}

func pong(pongChannel chan string) {
	for i := 0; i < 10; i++ {
		s := fmt.Sprintf("pong %d", i)
		fmt.Printf("send %s\n", s)
		pongChannel <- s
	}
}

func print(pingChannel chan string, pongChannel chan string) {
	for {
		select {
		case pingMessage := <-pingChannel:
			fmt.Printf("received %s\n", pingMessage)
			time.Sleep(time.Second * 1)
		case pongMessage := <-pongChannel:
			fmt.Printf("received %s\n", pongMessage)
			time.Sleep(time.Second * 1)
		}
	}
}

func TestChannelsWorkerJobs(t *testing.T) {

	start := time.Now()

	jobs := make(chan int, 100)
	results := make(chan int, 100)

	// start 3 workers
	for i := 1; i <= 3; i++ {
		go worker(i, jobs, results)
	}

	// schedule 5 jobs
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	// collect the results
	for k := 1; k <= 5; k++ {
		<-results
	}

	elapsed := time.Since(start)
	fmt.Printf("completing jobs all jobs took %s", elapsed)
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}
