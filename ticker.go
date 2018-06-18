package main

import "fmt"
import "time"

func main() {
	ticker := time.NewTicker(time.Millisecond * 50)

	go func() {
		for t := range ticker.C {
			fmt.Println("ticker kick ", t)
		}
	}()

	time.Sleep(time.Second * 3)
	ticker.Stop()
	fmt.Println("Ticker stop")
}
