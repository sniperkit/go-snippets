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

var nsqConsumer *nsq.Consumer

func listenHealthCheck() {
	fmt.Println("Starting health check listener")
	http.HandleFunc("/health_check", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Doing fine.")
	})
	http.ListenAndServe(":9000", nil)
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
	fmt.Printf("Starting app...")
	nsqConsumer, _ = nsq.NewConsumer("worker_result", "app", nsq.NewConfig())
	nsqConsumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		var sum int
		json.Unmarshal(message.Body, &sum)
		fmt.Println("Got worker result", sum)
		return nil
	}))
	go listenHealthCheck()
	connectNsqConsumer()

	// wait forever
	select {}
}
