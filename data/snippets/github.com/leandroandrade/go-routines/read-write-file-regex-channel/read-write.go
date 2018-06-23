package main

import (
	"sync"
	"regexp"
	"time"
	"fmt"
	"os"
	"bufio"
	"log"
)

const numberWorkers = 8

var waitGroup sync.WaitGroup
var regex = regexp.MustCompile("\"([0-9]+_[0-9]+)\"")

func worker(queue chan []byte, id int, logger *os.File) {
	defer waitGroup.Done()

	for {
		line, ok := <-queue
		if !ok {
			fmt.Printf("Worker: %d : Shutting Down\n", id)
			return
		}
		if regex.Match(line) {
			logger.WriteString(string(line) + "\n")
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

	logger, err := os.OpenFile("FILE_TO_WRITE", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Errorf("erro %v", err)
	}
	//defer logger.Close()

	waitGroup.Add(numberWorkers)

	// create workers to listening the channel
	queue := make(chan []byte)
	for gr := 1; gr <= numberWorkers; gr++ {
		go worker(queue, gr, logger)
	}

	// read each line from file
	reader := bufio.NewReader(file)
	for line, _, err := reader.ReadLine(); err == nil; line, _, err = reader.ReadLine() {
		queue <- line
	}

	// to error, use this
	if errClose := logger.Close(); err == nil {
		err = errClose
		if err != nil {
			log.Fatalf("error when write a file: %v", err)
		}
	}

	// or use this
	//defer func() {
	//	if errClose := logger.Close(); err == nil {
	//		err = errClose
	//		if err != nil {
	//			log.Fatalf("error when write a file: %v", err)
	//		}
	//	}
	//}()

	close(queue)
	waitGroup.Wait()

	fmt.Printf("Elapsed time: %v\n", time.Since(start))
}
