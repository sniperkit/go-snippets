package survey

import (
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/respondent"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"
	// "strings"

	"testing"
	"time"
)

const url = "tcp://127.0.0.1:600"

var surveyResponseChannel chan string
var tGlobal *testing.T

//Creates a new Survey Server and Tests a single Survey
func TestSingleSurvey(t *testing.T) {
	tGlobal = t
	// surveyResponseChannel = make(chan string)

	var s Server
	t.Log("Starting Survey Server")

	err := s.Listen(url, 100, 2, handleSurveyResponse)
	if err != nil {
		t.Errorf("Error at survey.TestSingleSurvey:  %v", err.Error())
	}

	var sock mangos.Socket

	if sock, err = respondent.NewSocket(); err != nil {
		t.Errorf("Error creating new Socket at survey.TestSingleSurvey:  %v", err.Error())
		return
	}

	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())

	t.Log("Connecting to Survey Server")

	if err = sock.Dial(url); err != nil {
		t.Errorf("Error Dialing at survey.TestSingleSurvey:  %v", err.Error())
		return
	}

	messages := make(chan string)

	go respondToSurvey(sock, t, messages, "TestSurvey", "HelloWorld")

	time.Sleep(1 * time.Second)

	err = s.Send([]byte("TestSurvey"))
	if err != nil {
		t.Errorf("Error sending survey message at survey.TestSingleSurvey:  %v", err.Error())
		return
	}

	time.Sleep(1 * time.Second)
	msg := <-messages
	t.Log(msg)
}

func handleSurveyResponse(msg []byte) {

	if string(msg) != "HelloWorld" {
		tGlobal.Errorf("Failed to match the survey response message at survey.TestSingleSurvey")
		return
	}

}

//Responds to the Survey
func respondToSurvey(sock mangos.Socket, t *testing.T, messages chan string, surveyQuestion string, surveyResponse string) {
	var err error
	var msg []byte

	if msg, err = sock.Recv(); err != nil {
		t.Errorf("Error Receiving at survey.respondToSurvey:  %v", err.Error())
		messages <- "Test Failed"
		return
	}

	if string(msg) != surveyQuestion {
		t.Errorf("Failed to respond to survey question.")
		messages <- "Test Failed"
		return
	}

	if err = sock.Send([]byte(surveyResponse)); err != nil {
		t.Errorf("Error Sending Survey Response at survey.respondToSurvey:  %v", err.Error())
		messages <- "Test Failed"
		return
	}

	messages <- "Test Completed"
}
