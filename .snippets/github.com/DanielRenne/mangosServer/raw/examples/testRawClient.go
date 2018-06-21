package main

import (
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/req"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"
	"log"
)

const url = "tcp://127.0.0.1:600"

func main() {

	go makeRequest(url, "Hello World")

	for {

	}
}

func makeRequest(url string, message string) {
	var sock mangos.Socket
	var err error

	if sock, err = req.NewSocket(); err != nil {
		log.Printf("Error creating new Socket at raw.TestSingleRawMessage:  %v", err.Error())
		return
	}

	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())

	log.Println("Connecting to Raw Server")

	if err = sock.Dial(url); err != nil {
		log.Printf("Error Dialing at raw.TestSingleRawMessage:  %v", err.Error())
		return
	}

	body := []byte(message)
	m := mangos.NewMessage(len(body) + 2)
	m.Body = body

	log.Println("Sending Request to Raw Server")
	if err = sock.SendMsg(m); err != nil {
		log.Printf("Error sending raw message at raw.TestSingleRawMessage:  %v", err.Error())
		return
	}

	respondToRaw(sock)

}

//Responds to the Raw Socket Message
func respondToRaw(sock mangos.Socket) {

	msg, err := sock.RecvMsg()

	if err != nil {
		return
	}

	log.Println(string(msg.Body))

}
