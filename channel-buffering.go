package main

import (
	"fmt"
)

func main() {
	messages := make(chan string, 2)

	messages <- "buffered"
	fmt.Println(<-messages)
	messages <- "channel"
	messages <- "channel"

	fmt.Println(<-messages)
	fmt.Println(<-messages)
}
