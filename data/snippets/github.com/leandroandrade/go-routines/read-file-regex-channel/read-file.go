package main

import (
	"sync"
	"regexp"
	"time"
	"fmt"
	"os"
	"bufio"
)

const numberWorkers = 4

var waitGroup sync.WaitGroup
var regex = regexp.MustCompile("\"([0-9]+_[0-9]+)\"")

func worker(queue chan []byte, id int) {
	defer waitGroup.Done()

	for {
		line, ok := <-queue
		if !ok {
			fmt.Printf("Worker: %d : Shutting Down\n", id)
			return
		}
		if regex.Match(line) {
			fmt.Printf("[%v] Match: %v\n", id, string(line))
		}
	}
}

func main() {
	start := time.Now()
	fmt.Printf("Started : %v\n", start)

	file, err := os.Open("FILE_TO_READ")
	if err != nil {
		fmt.Errorf("erro %v", err)
	}
	defer file.Close()

	waitGroup.Add(numberWorkers)

	// create workers to listening the channel
	queue := make(chan []byte)
	for gr := 1; gr <= numberWorkers; gr++ {
		go worker(queue, gr)
	}

	// read each line from file
	reader := bufio.NewReader(file)
	for line, _, err := reader.ReadLine(); err == nil; line, _, err = reader.ReadLine() {
		queue <- line
	}

	close(queue)
	waitGroup.Wait()

	fmt.Printf("Elapsed time: %v\n", time.Since(start))
}
