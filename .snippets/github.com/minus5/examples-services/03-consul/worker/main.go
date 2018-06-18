package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

var nsqProducer *nsq.Producer
var nsqConsumer *nsq.Consumer

func listenHealthCheck() {
	fmt.Println("Starting health check listener")
	http.HandleFunc("/health_check", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Doing fine.")
	})
	http.ListenAndServe(":9002", nil)
}

// resolves service address and port by consul
func consulResolveService(name string) string {
	type consulResponse struct {
		ServiceAddress string
		ServicePort    int
	}
	var crs []consulResponse
	for true {
		consulAddr := os.Getenv("CONSUL_ADDR")
		if consulAddr == "" {
			consulAddr = "127.0.0.1"
		}
		query := "http://" + consulAddr + ":8500/v1/catalog/service/" + name
		fmt.Println("Resolving service by Consul:", query)
		req, err := http.Get(query)
		if err == nil {
			body, err := ioutil.ReadAll(req.Body)
			if err == nil {
				json.Unmarshal(body, &crs)
				if len(crs) > 0 {
					cr := crs[0]
					result := fmt.Sprintf("%s:%d", cr.ServiceAddress, cr.ServicePort)
					fmt.Println("Consul resloved", name, result)
					return result
				}
			}
		}
		time.Sleep(time.Second * 1)
	}
	return ""
}

func connectNsqProducer() {
	for true {
		fmt.Println("Connecting Producer to nsqd...")
		var err error
		nsqProducer, err = nsq.NewProducer(consulResolveService("nsqd"), nsq.NewConfig())
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
		err := nsqConsumer.ConnectToNSQD(consulResolveService("nsqd"))
		if err == nil {
			return
		}
		time.Sleep(time.Second * 1)
	}
}

func main() {
	fmt.Println("Starting worker...")
	go listenHealthCheck()
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
