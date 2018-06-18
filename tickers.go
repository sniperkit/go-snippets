package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()

	time.Sleep(time.Millisecond * 1600)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}

//Tickers use a similar mechanism to timers: a channel that is sent values.
// Here we’ll use the range builtin on the channel to iterate over the values as they arrive every 500ms.
//Tickers can be stopped like timers. Once a ticker is stopped it won’t receive any
// more values on its channel. We’ll stop ours after 1600ms
