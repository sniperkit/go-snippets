package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

var nsqProducer *nsq.Producer

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

func sendSensorData() {
	var values [3]uint32
	for i := range values {
		values[i] = rand.Uint32()
	}
	jsonValue, _ := json.Marshal(values)
	fmt.Println("Sending values", string(jsonValue))
	nsqProducer.Publish("sensor_values", jsonValue)
}

func main() {
	fmt.Println("Starting sensor...")
	connectNsqProducer()
	for range time.Tick(time.Second * 1) {
		sendSensorData()
	}
}
