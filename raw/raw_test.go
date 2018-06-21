package raw

import (
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/req"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"
	// "strings"

	"testing"
	// "time"
)

const url1 = "tcp://127.0.0.1:600"
const url2 = "tcp://127.0.0.1:700"

var surveyResponseChannel chan string
var tGlobal *testing.T

//Creates a new Survey Server and Tests a single Survey
func TestSingleRawMessage(t *testing.T) {
	tGlobal = t
	// surveyResponseChannel = make(chan string)

	var s Server
	t.Log("Starting Raw Server")

	err := s.Listen(url1, 2, handleRawResponse)
	if err != nil {
		t.Errorf("Error at raw.TestSingleRawMessage:  %v", err.Error())
	}

	makeRequest(url1, t, "TestRaw")

}

//Creates a new Survey Server and Tests a single Survey
func TestTwoRawMessages(t *testing.T) {
	tGlobal = t
	// surveyResponseChannel = make(chan string)

	var s Server
	t.Log("Starting Raw Server")

	err := s.Listen(url2, 2, handleRawResponse)
	if err != nil {
		t.Errorf("Error at raw.TestTwoRawMessages:  %v", err.Error())
	}

	makeRequest(url2, t, "Connection1")
	makeRequest(url2, t, "Connection2")

}

func handleRawResponse(s *Server, m *mangos.Message) {
	s.Reply(m)
}

func makeRequest(url string, t *testing.T, message string) {
	var sock mangos.Socket
	var err error

	if sock, err = req.NewSocket(); err != nil {
		t.Errorf("Error creating new Socket at raw.TestSingleRawMessage:  %v", err.Error())
		return
	}

	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())

	t.Log("Connecting to Raw Server")

	if err = sock.Dial(url); err != nil {
		t.Errorf("Error Dialing at raw.TestSingleRawMessage:  %v", err.Error())
		return
	}

	body := []byte(message)
	m := mangos.NewMessage(len(body) + 2)
	m.Body = body

	t.Log("Sending Request to Raw Server")
	if err = sock.SendMsg(m); err != nil {
		t.Errorf("Error sending raw message at raw.TestSingleRawMessage:  %v", err.Error())
		return
	}

	respondToRaw(sock, t, message)
}

//Responds to the Raw Socket Message
func respondToRaw(sock mangos.Socket, t *testing.T, rawData string) {

	msg, err := sock.RecvMsg()

	if err != nil {
		t.Errorf("Error Receiving at raw.respondToRaw:  %v", err.Error())
		return
	}

	if string(msg.Body) != rawData {
		t.Errorf("Failed to respond to raw question.")
		return
	}
}
