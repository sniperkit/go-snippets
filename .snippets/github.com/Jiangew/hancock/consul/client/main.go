package main

import (
	"fmt"
	"log"
	consulapi "github.com/hashicorp/consul/api"
	"net"
	"time"
)

const RECV_BUF_LEN = 1024

// sendData
func sendData(service *consulapi.AgentService) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", service.Address, service.Port))
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	buf := make([]byte, RECV_BUF_LEN)
	i := 0
	for {
		i++
		msg := fmt.Sprintf("Hello Go, %03d", i)

		n, err := conn.Write([]byte(msg))
		if err != nil {
			println("Write Buffer Err:", err.Error())
			break
		}

		n, err = conn.Read(buf)
		if err != nil {
			println("Read Buffer Err:", err.Error())
			break
		}

		log.Println("get:", string(buf[0:n]))

		time.Sleep(time.Second)
	}
}

func main() {
	client, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		log.Fatal("consul client error : ", err)
	}

	for {
		time.Sleep(time.Second * 3)

		var services map[string]*consulapi.AgentService
		var err error

		services, err = client.Agent().Services()
		if err != nil {
			log.Println("in consual list Services:", err)
			continue
		}

		if _, found := services["serverNode_1"]; !found {
			log.Println("serverNode_1 not found")
			continue
		}

		sendData(services["serverNode_1"])
	}
}
