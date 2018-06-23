package main

import (
	"container/list"
	"fmt"
	"time"
)

func main() {
	bufferReuse()
}

func stopSignal() {

	c := make(chan struct{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c)
	}()
	<-c
	fmt.Println("done")
}

func startSignal() {

	worker := func(start chan struct{}, id int) {
		<-start
		fmt.Println("worker", id)
	}

	start := make(chan struct{})

	for i := 0; i < 100; i++ {
		go worker(start, i)
	}

	close(start)

	time.Sleep(2 * time.Second)
}

func stopWorkers() {

	worker := func(stop chan struct{}, id int) {
		for {
			select {
			case <-stop:
				fmt.Println("stop worker", id)
				return
			}
		}
	}

	stop := make(chan struct{})

	for i := 0; i < 100; i++ {
		go worker(stop, i)
	}

	close(stop)

	time.Sleep(2 * time.Second)
}

func stopAndWaitWorker() {

	worker := func(stop chan bool) {
		for {
			select {

			case <-stop:
				fmt.Println("stop worker")
				stop <- true
				return
			}
		}
	}

	stop := make(chan bool)

	go worker(stop)

	stop <- true
	<-stop
}

func uniqueId() {

	id := make(chan string)

	go func() {
		var counter int
		for {
			s := fmt.Sprintf("%04x", counter)
			counter++
			id <- s
		}
	}()

	for i := 0; i < 100; i++ {
		fmt.Printf("%s\n", <-id)
	}
}

func bufferReuse() {

	const bufferSize = 1024

	var (
		get = make(chan []byte)
		put = make(chan []byte)
	)

	go func() {
		l := list.New()

		for {
			if l.Len() == 0 {
				l.PushBack(make([]byte, bufferSize))
			}

			e := l.Front()

			select {

			case b := <-put:
				l.PushBack(b)

			case get <- e.Value.([]byte):
				l.Remove(e)
			}
		}
	}()

	buffer := <-get
	put <- buffer

	buffer = <-get
	put <- buffer
}
