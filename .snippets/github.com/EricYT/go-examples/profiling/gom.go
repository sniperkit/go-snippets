package main

import (
	"log"
	"net/http"
	"sync"

	gomhttp "github.com/rakyll/gom/http"
)

func main() {

	var wg sync.WaitGroup

	signal := make(chan bool)

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			wg.Done()
			var count int64
			for i := 0; i < 10000000000; i++ {
				count++
				log.Println("value:", count)
			}
			<-signal
		}()
	}
	wg.Wait()

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/_gom", gomhttp.Handler())
	log.Println(http.ListenAndServe("localhost:6060", nil))
}
