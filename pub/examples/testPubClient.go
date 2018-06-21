package main

import (
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/sub"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"
	"log"
)

const url = "tcp://127.0.0.1:600"

func main() {

	var sock mangos.Socket
	var err error

	if sock, err = sub.NewSocket(); err != nil {
		log.Printf("Error Creating Socket at TestPubSubBroadCastAll:  %v", err.Error())
		return
	}
	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	if err = sock.Dial(url); err != nil {
		log.Printf("Error Dialing at TestPubSubBroadCastAll:  %v", err.Error())
		return
	}

	go subscribeToAll(sock)

	for {

	}

}

//Subscribes to AllMessages
func subscribeToAll(sock mangos.Socket) {

	var msg []byte
	// Empty byte array effectively subscribes to everything
	err := sock.SetOption(mangos.OptionSubscribe, []byte(""))
	if err != nil {
		log.Printf("Error Subscribing at subscribeToAll:  %v", err.Error())
		return
	}

	if msg, err = sock.Recv(); err != nil {
		log.Printf("Error Recieving data at subscribeToAll:  %v", err.Error())
		return
	}

	log.Printf(string(msg))
	go subscribeToAll(sock)
}
