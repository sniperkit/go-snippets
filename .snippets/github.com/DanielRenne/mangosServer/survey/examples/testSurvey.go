package main

import (
	"github.com/DanielRenne/mangosServer/survey"
	"log"
	"time"
)

const url = "tcp://127.0.0.1:600"

func main() {
	var s survey.Server

	err := s.Listen(url, 500, 2, handleSurveyResponse)
	if err != nil {
		log.Printf("Error:  %v", err.Error)
	}

	//Code a forever loop to stop main from exiting.
	for {
		time.Sleep(3 * time.Second)
		go s.Send([]byte("Sending Survey"))
	}

}

func handleSurveyResponse(msg []byte) {
	//Process Survey Results.
	log.Printf(string(msg))
}
