package main

import "fmt"
import "strconv"

type show interface {
	Show(int)
}

func worker(msg <-chan string, res chan<- string) {
	select {
	case work := <-msg:
		fmt.Println("Message from parent:", work)
		res <- work
	}
	worker(msg, res)
}

func receiver(msg <-chan string, end chan<- string) {
	var completeCount int = 0
	for {
		if completeCount < 10 {
			select {
			case complete := <-msg:
				fmt.Println("Worker complete from child:", complete)
        completeCount++
			}
      continue
		}
		break
	}
  end<- "Game over"
}

func main() {
	message := make(chan string)
	result := make(chan string)
	total := make(chan string)

	go worker(message, result)
	go receiver(result, total)

	for i := 0; i < 10; i++ {
		message <- strconv.Itoa(i)
	}

	select {
	case end := <-total:
		fmt.Println("The job is done:", end)
	}

}
