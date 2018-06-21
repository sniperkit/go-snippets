package pub

import (
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/sub"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"
	"strings"
	"testing"
)

const url = "tcp://127.0.0.1:600"
const urlTopic = "tcp://127.0.0.1:601"

//Creates a new Pub Server and broadcasts a plain message
func TestPubSubBroadCastAll(t *testing.T) {
	var s Server
	err := s.Listen(url)
	if err != nil {
		t.Errorf("Error at pub_test.TestPubSubBroadCastAll:  %v", err.Error())
	}

	var sock mangos.Socket
	var errSubScribe error

	if sock, errSubScribe = sub.NewSocket(); errSubScribe != nil {
		t.Errorf("Error Creating Socket at pub_test.TestPubSubBroadCastAll:  %v", errSubScribe.Error())
		return
	}
	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	if errSubScribe = sock.Dial(url); errSubScribe != nil {
		t.Errorf("Error Dialing at pub_test.TestPubSubBroadCastAll:  %v", errSubScribe.Error())
		return
	}

	messages := make(chan string)

	go subscribeToAll(sock, t, messages)

	s.Publish([]byte("TestSubscribeAll"))

	msg := <-messages
	t.Log(msg)

}

//Creates a new Pub Server and broadcasts a topic message
func TestPubSubBroadCastTopic(t *testing.T) {
	var s Server
	err := s.Listen(urlTopic)
	if err != nil {
		t.Errorf("Error at pubsub_test.TestPubSubBroadCastTopic:  %v", err.Error())
	}

	var sock mangos.Socket
	var errSubScribe error

	if sock, errSubScribe = sub.NewSocket(); errSubScribe != nil {
		t.Errorf("Error Creating Socket at pub_test.TestPubSubBroadCastTopic:  %v", errSubScribe.Error())
		return
	}
	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	if errSubScribe = sock.Dial(urlTopic); errSubScribe != nil {
		t.Errorf("Error Dialing at pub_test.TestPubSubBroadCastTopic:  %v", errSubScribe.Error())
		return
	}

	messages := make(chan string)

	go subscribeToTopic(sock, t, "TestTopic", messages)

	s.PublishTopic("TestTopic", "TestSubscribeTopic")

	msg := <-messages
	t.Log(msg)

}

//Subscribes to AllMessages
func subscribeToAll(sock mangos.Socket, t *testing.T, messages chan string) {

	var msg []byte
	// Empty byte array effectively subscribes to everything
	err := sock.SetOption(mangos.OptionSubscribe, []byte(""))
	if err != nil {
		t.Errorf("Error Subscribing at pub_test.subscribeToAll:  %v", err.Error())
		messages <- "Test Failed"
		return
	}

	if msg, err = sock.Recv(); err != nil {
		t.Errorf("Error Recieving data at pub_test.subscribeToAll:  %v", err.Error())
		messages <- "Test Failed"
		return
	}

	if string(msg) != "TestSubscribeAll" {
		t.Errorf("Error Asserting published data at pub_test.subscribeToAll")
		t.Log("Assert Message:  " + string(msg))
		messages <- "Test Failed"
		return
	}

	messages <- "Test Completed"
}

//Subscribes to a specific Topic
func subscribeToTopic(sock mangos.Socket, t *testing.T, topic string, messages chan string) {

	var msg []byte
	// Empty byte array effectively subscribes to everything
	err := sock.SetOption(mangos.OptionSubscribe, []byte(topic))
	if err != nil {
		t.Errorf("Error Subscribing at pub_test.subscribeToTopic:  %v", err.Error())
		messages <- "Test Failed"
		return
	}

	if msg, err = sock.Recv(); err != nil {
		t.Errorf("Error Recieving data at pub_test.subscribeToTopic:  %v", err.Error())
		messages <- "Test Failed"
		return
	}

	msgTopic := strings.Replace(string(msg), topic+"|", "", -1)

	if msgTopic != "TestSubscribeTopic" {
		t.Errorf("Error Asserting published data at pub_test.subscribeToTopic")
		t.Log("Assert Message:  " + msgTopic)
		messages <- "Test Failed"
		return
	}

	messages <- "Test Completed"
}
