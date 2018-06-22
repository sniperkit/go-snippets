package utils

import (
	"github.com/Shopify/sarama"
	"github.com/howbazaar/loggo"
	"strings"
	"os"
	"os/signal"
	"encoding/json"
	"go-temp-project/model"
)

var producer sarama.AsyncProducer
var signals chan os.Signal
var kafkaUp bool

func KafkaInit() {

	config := sarama.NewConfig()
	var err error
	producer, err = sarama.NewAsyncProducer(strings.Split("localhost:9092", ","), config)
	if err == nil {
		logger := loggo.GetLogger("")
		logger.Infof("===> KAFKA UP (%s) <===")
		kafkaUp = true
	}
	// Trap SIGINT to trigger a graceful shutdown.
	signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go func() {
		<- signals
		producer.AsyncClose() // Trigger a shutdown of the producer.
	}()
}

func KafkaSend(topic string, value model.KafkaMonitoringLog) {
	if kafkaUp {
		go _KafkaSend(topic, value)
	}
}

func _KafkaSend(topic string, value model.KafkaMonitoringLog) {
	kafkaMessage, _ := json.Marshal(value)
	message := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(string(kafkaMessage))}
	producer.Input() <- message
}
