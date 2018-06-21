package main

import (
	"github.com/DanielRenne/mangosServer/raw"
	"github.com/go-mangos/mangos"
	"log"
)

const url = "tcp://127.0.0.1:600"

func main() {
	var s raw.Server

	err := s.Listen(url, 2, handleRawRequest)
	if err != nil {
		log.Printf("Error:  %v", err.Error)
	}

	//Code a forever loop to stop main from exiting.
	for {

	}

}

func handleRawRequest(s *raw.Server, m *mangos.Message) {

	log.Printf(string(m.Body))
	m.Body = []byte("Custom Response to the Request")
	s.Reply(m)
}
