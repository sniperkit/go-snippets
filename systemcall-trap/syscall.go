package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"sync"
	"syscall"
	"time"
)

var pipefile = "pipe.log"

func main() {
	// spawn others
	for i := 0; i < 3; i++ {
		go func() {
			for {
				time.Sleep(time.Millisecond * 3)
			}
		}()
	}

	// clear pipe file
	os.Remove(pipefile)

	// create pipe
	err := syscall.Mkfifo(pipefile, 0666)
	if err != nil {
		panic(err)
	}

	log.Println("pipe file created")

	var wg sync.WaitGroup

	readPipeFunc := func() {
		defer wg.Done()

		file, err := os.OpenFile(pipefile, os.O_CREATE, os.ModeNamedPipe)
		if err != nil {
			panic(err)
		}

		log.Println("pipe file opened")
		reader := bufio.NewReader(file)

		// block read
		for {
			buf := make([]byte, 512)
			n, err := reader.Read(buf)
			if err != nil && err != io.EOF {
				panic(err)
			}

			log.Printf("pipe read %d byte: %s\n", n, string(buf))
		}
	}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go readPipeFunc()
	}

	wg.Wait()

}
