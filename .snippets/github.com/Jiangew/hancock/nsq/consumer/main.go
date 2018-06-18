package main

import (
	"sync"
	"github.com/nsqio/go-nsq"
	"log"
	"fmt"
	"github.com/jiangew/hancock/utils/reflect_invoke"
	"encoding/json"
	"time"
	"os/signal"
	"syscall"
	"os"
)

func simpleConsumer() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	config := nsq.NewConfig()
	consumer, _ := nsq.NewConsumer("topic_test", "ch", config)

	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Printf("Got a message: %v", message)
		wg.Done()
		return nil
	}))

	//err := consumer.ConnectToNSQD("127.0.0.1:4150")
	err := consumer.ConnectToNSQLookupd("127.0.0.1:4161")
	if err != nil {
		log.Panic("Could not connect")
	}

	wg.Wait()
}

type Bar struct {
}

type Foo struct {
}

func (b *Bar) BarFuncAdd(argOne, argTwo float64) float64 {
	return argOne + argTwo
}

func (b *Foo) FooFuncSwap(argOne, argTwo string) (string, string) {
	return argTwo, argOne
}

func HandleStringMessage(message *nsq.Message) error {
	fmt.Printf("HandleStringMessage get a message: %v\n\n", string(message.Body))
	return nil
}

func HandleJsonMessage(message *nsq.Message) error {
	resultJson := reflect_invoke.InvokeByJson([]byte(message.Body))
	result := reflect_invoke.Response{}

	err := json.Unmarshal(resultJson, &result)
	if err != nil {
		return err
	}

	var info []string
	info = append(info, "HandleJsonMessage got a result\n")
	info = append(info, "raw:\n"+string(resultJson)+"\n")
	info = append(info, "function: "+result.FuncName+" \n")
	info = append(info, fmt.Sprintf("result: %v\n", result.Data))
	info = append(info, fmt.Sprintf("error: %d,%s\n\n", result.ErrorCode, reflect_invoke.ErrorMsg(result.ErrorCode)))

	fmt.Println(info)
	return nil
}

func MakeConsumer(topic, channel string, config *nsq.Config, handle func(message *nsq.Message) error) {
	consumer, _ := nsq.NewConsumer(topic, channel, config)
	consumer.AddHandler(nsq.HandlerFunc(handle))

	err := consumer.ConnectToNSQLookupd("127.0.0.1:4161")
	if err != nil {
		log.Panic("Could not connect")
	}
}

func init() {
	bar := &Bar{}
	foo := &Foo{}

	reflect_invoke.RegisterMethod(bar)
	reflect_invoke.RegisterMethod(foo)
}

func main() {
	config := nsq.NewConfig()
	config.DefaultRequeueDelay = 0
	config.MaxBackoffDuration = 20 * time.Microsecond
	config.LookupdPollInterval = 1000 * time.Millisecond
	config.RDYRedistributeInterval = 1000 * time.Millisecond
	config.MaxInFlight = 2500

	MakeConsumer("topic_string", "ch", config, HandleStringMessage)
	MakeConsumer("topic_json", "ch", config, HandleJsonMessage)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	fmt.Println("exit")
}
