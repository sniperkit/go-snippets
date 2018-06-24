package main

import (
	"sync"
)

func smallBufferFunc(ch chan Graph, in chan Graph) {
	var buffer []Graph
	for {
		if len(buffer) == 0 {
			x, ok := <-in
			if !ok {
				return
			}
			buffer = append(buffer, x)
		}
		select {
		case x, ok := <-in:
			if !ok {
				return
			}
			buffer = append(buffer, x)

		case ch <- buffer[0]:
			buffer = buffer[1:]

		case <-DefaultDone:
			return
		}
	}
}

func smallBuffer(in chan Graph) chan Graph {
	ch := make(chan Graph)

	var wg sync.WaitGroup
	wg.Add(DefaultNumBuffers)
	for i := 0; i < DefaultNumBuffers; i++ {
		go smallBufferFunc(ch, in)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}




