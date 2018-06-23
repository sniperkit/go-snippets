package channel

import (
	"fmt"
	"time"
)

func ExpChannel() {
	fmt.Println("Channel exp calıstırıldı...")
	data := make(chan bool)
	go func() {
		time.Sleep(4 * time.Second)
		data <- true
	}()
	<-data
	fmt.Println("End!")
}
