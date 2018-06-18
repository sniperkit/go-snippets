package main

import (
	"encoding/json"
	"fmt"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

var nsqProducer *nsq.Producer
var nsqConsumer *nsq.Consumer

func connectNsqProducer() {
	for true {
		fmt.Println("Connecting Producer to nsqd...")
		var err error
		nsqProducer, err = nsq.NewProducer("127.0.0.1:4150", nsq.NewConfig())
		if err == nil {
			err = nsqProducer.Ping()
			if err == nil {
				return
			}
		}
		time.Sleep(time.Second * 1)
	}
}

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
	fmt.Println("Starting worker...")
	nsqConsumer, _ = nsq.NewConsumer("sensor_values", "worker", nsq.NewConfig())
	nsqConsumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		var numbers []int
		json.Unmarshal(message.Body, &numbers)
		fmt.Println("Got request", numbers)
		sum := 0
		for _, value := range numbers {
			sum += value
		}
		fmt.Println("Sending result", sum)
		msgBody, _ := json.Marshal(sum)
		nsqProducer.Publish("worker_result", msgBody)
		return nil
	}))
	connectNsqProducer()
	connectNsqConsumer()

	// wait forever
	select {}
}
