package main

import (
	"github.com/nsqio/go-nsq"
	"log"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func simpleProducer() {
	config := nsq.NewConfig()
	producer, _ := nsq.NewProducer("127.0.0.1:4150", config)

	err := producer.Publish("topic_test", []byte("Hello NSQ"))
	if err != nil {
		log.Panic("Could not connect")
	}

	producer.Stop()
}

func main() {
	config := nsq.NewConfig()
	producer, _ := nsq.NewProducer("127.0.0.1:4150", config)

	for i := 0; i < 10; i++ {
		producer.Publish("topic_string", []byte(fmt.Sprintf("string%d", i)))
	}

	var jsonData []string
	jsonData = append(jsonData, `
		{
			"func_name":"BarFuncAdd",
			"params":[
				0.5,
				0.51
			]
		}
	`)
	jsonData = append(jsonData, `
		{
			"func_name":"FooFuncSwap",
			"params":[
				"a",
				"b"
			]
		}
	`)
	jsonData = append(jsonData, `
		{
			"func_name":"FooFuncSwap",
			"params":[
				1,
				2
			]
		}
	`)
	jsonData = append(jsonData, `
		{
			"func_name":"FakeMethod",
			"params":[
				"a",
				"b"
			]
		}
	`)

	for _, meta := range jsonData {
		producer.Publish("topic_json", []byte(meta))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	producer.Stop()
}
