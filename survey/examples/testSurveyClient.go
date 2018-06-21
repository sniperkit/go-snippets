package main

import (
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/respondent"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"
	"log"
)

const url = "tcp://127.0.0.1:600"

func main() {

	go makeSurveyConnection(url)

	for {

	}
}

func makeSurveyConnection(url string) {
	var err error
	var sock mangos.Socket

	if sock, err = respondent.NewSocket(); err != nil {
		log.Printf("Error creating new Socket at survey.TestSingleSurvey:  %v", err.Error())
		return
	}

	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())

	log.Println("Connecting to Survey Server")

	if err = sock.Dial(url); err != nil {
		log.Printf("Error Dialing at survey.TestSingleSurvey:  %v", err.Error())
		return
	}

	go respondToSurvey(sock, "Hello Survey.")

}

//Responds to the Survey
func respondToSurvey(sock mangos.Socket, surveyResponse string) {
	var err error
	var msg []byte

	if msg, err = sock.Recv(); err != nil {
		log.Printf("Error Receiving at respondToSurvey:  %v", err.Error())
		return
	}

	log.Println(string(msg))

	if err = sock.Send([]byte(surveyResponse)); err != nil {
		log.Printf("Error Sending Survey Response at respondToSurvey:  %v", err.Error())
		return
	}

	go respondToSurvey(sock, "Hello Survey.")
}
