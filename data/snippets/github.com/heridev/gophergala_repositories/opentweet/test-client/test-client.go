package main

import (
	"github.com/gophergala/opentweet/protocol"
	"time"
	"fmt"
	"os"
)

func main() {
	server := os.Args[1]
	user := os.Args[2]
	tweets, err := protocol.GetTweets(
		server,
		user,
		time.Now().Add(-24 * time.Hour),
		time.Now(),
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	for _, val := range(tweets) {
		fmt.Printf("Time: %v\n", val.Time)
		fmt.Printf("Tweet: %v\n", val.Text)
	}
}
		
