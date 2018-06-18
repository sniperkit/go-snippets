package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

var nsqProducer *nsq.Producer

// exposes health check listener fo consul
func listenHealthCheck() {
	fmt.Println("Starting health check listener")
	http.HandleFunc("/health_check", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Doing fine.")
	})
	http.ListenAndServe(":9001", nil)
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
	go listenHealthCheck()
	connectNsqProducer()
	for range time.Tick(time.Second * 1) {
		sendSensorData()
	}
}
