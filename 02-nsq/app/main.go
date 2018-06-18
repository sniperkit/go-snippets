package main

import (
	"encoding/json"
	"fmt"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

var nsqConsumer *nsq.Consumer

func connectNsqConsumer() {
	for true {
		fmt.Println("Connecting Consumer to nsqd...")
		err := nsqConsumer.ConnectToNSQD("127.0.0.1:4150")
		if err == nil {
			return
		}
		time.Sleep(time.Second * 1)
	}
}

func main() {
	fmt.Printf("Starting app...")
	nsqConsumer, _ = nsq.NewConsumer("worker_result", "app", nsq.NewConfig())
	nsqConsumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		var sum int
		json.Unmarshal(message.Body, &sum)
		fmt.Println("Got worker result", sum)
		return nil
	}))
	connectNsqConsumer()

	// wait forever
	select {}
}
