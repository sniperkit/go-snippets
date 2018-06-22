package rethinkdb

import (
	"github.com/asim/go-micro/store"
	"github.com/bitly/go-nsq"
)

func Producer() (*nsq.Producer, error) {
	var address string

	item, err := store.Get("nsq/address")
	if err != nil {
		address = "localhost:4160"
	} else {
		address = string(item.Value())
	}

	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(address, config)
	if err != nil {
		return nil, err
	}

	return producer, err
}

func Consumer(topic, channel string) (*nsq.Consumer, error) {
	var err error

	var address string

	item, err := store.Get("nsq/address")
	if err != nil {
		address = "localhost:4160"
	} else {
		address = string(item.Value())
	}

	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		return nil, err
	}

	consumer.ConnectToNSQD(address)

	return consumer, err
}
