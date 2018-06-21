package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

type Test struct {
	Data [32758]int32
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	time.AfterFunc(time.Second*2, func() {
		for i := 0; i < 10; i++ {
			var t = Test{Data: [32758]int32{1}}
			go fmt.Sprint(&t)
		}
	})

	<-time.Tick(time.Minute * 10)
}
