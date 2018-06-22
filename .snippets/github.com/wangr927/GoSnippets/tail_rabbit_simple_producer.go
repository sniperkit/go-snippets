package main

import (
	"os"
	"github.com/hpcloud/tail"
	"sync"
	"fmt"
	"log"
	"github.com/streadway/amqp"
)

const (
	queueName = "hl"
	exchange = ""
	mqurl = "amqp://127.0.0.1:5672"
)

var conn *amqp.Connection
var channel *amqp.Channel

func main(){
	var wg sync.WaitGroup
	ch := make(chan string,1000)
	go func(ch chan string,wg sync.WaitGroup){
		wg.Add(1)
		defer wg.Done()
		for {
			ctx := <-ch
			push(ctx)
		}
	}(ch,wg)
	tailwork("collector.json",ch,true)
	wg.Wait()
}


func tailwork(fn string, outchan chan string, isfollow bool) error {
	_, err := os.Stat(fn)
	if err != nil{
		panic(err)
	}
	t, err := tail.TailFile(fn, tail.Config{Follow:isfollow})
	if err != nil{
		panic(err)
	}
	defer t.Stop()
	for line := range t.Lines{
		outchan <- line.Text
	}
	return nil
}

func failOnErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:%s", msg, err)
		panic(fmt.Sprintf("%s:%s", msg, err))
	}
}

func mqConnect() {
	var err error
	conn, err = amqp.Dial(mqurl)
	failOnErr(err, "failed to connect tp rabbitmq")

	channel, err = conn.Channel()
	failOnErr(err, "failed to open a channel")
}

func push(pushmsg string) {
	if channel == nil {
		mqConnect()
	}
	//msgContent := "hello"

	channel.Publish(exchange,queueName,false,false,amqp.Publishing{
		ContentType:"text/plain",
		Body: []byte(pushmsg),
	})
}