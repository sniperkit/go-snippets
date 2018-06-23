# mangosServer Pubsub

To publish a broadcast to all subscribers:

###Server Code

	package main
	
	import (
		"github.com/DanielRenne/mangosServer/pub"
		"log"
		"time"
	)
	
	const url = "tcp://127.0.0.1:600"
	
	//Creates a new Pub Server and broadcasts a plain message
	func main() {
		var s pub.Server
		err := s.Listen(url)
		if err != nil {
			log.Printf("Error:  %v", err.Error())
		}
	
		//Code a forever loop to stop main from exiting.
		for {
			time.Sleep(3 * time.Second)
			go s.Publish([]byte("Publishing Message."))
		}
	
	}

###Client Code 

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


To publish a Topic to all topic subscribers.

###Server Code

	package main
	
	import (
		"github.com/DanielRenne/mangosServer/pub"
		"log"
		"time"
	)
	
	const url = "tcp://127.0.0.1:600"
	
	//Creates a new Pub Server and broadcasts a plain message
	func main() {
		var s pub.Server
		err := s.Listen(url)
		if err != nil {
			log.Printf("Error:  %v", err.Error())
		}
	
		//Code a forever loop to stop main from exiting.
		for {
			time.Sleep(3 * time.Second)
			go s.PublishTopic("TestTopic", "Publishing Message.")
		}
	
	}

###Client Code
	
	package main
	
	import (
		"github.com/go-mangos/mangos"
		"github.com/go-mangos/mangos/protocol/sub"
		"github.com/go-mangos/mangos/transport/ipc"
		"github.com/go-mangos/mangos/transport/tcp"
		"log"
		"strings"
	)
	
	const url = "tcp://127.0.0.1:600"
	
	func main() {
	
		var sock mangos.Socket
		var err error
	
		if sock, err = sub.NewSocket(); err != nil {
			log.Printf("Error Creating Socket at pubsub_test.TestPubSubBroadCastAll:  %v", err.Error())
			return
		}
		sock.AddTransport(ipc.NewTransport())
		sock.AddTransport(tcp.NewTransport())
		if err = sock.Dial(url); err != nil {
			log.Printf("Error Dialing at pubsub_test.TestPubSubBroadCastAll:  %v", err.Error())
			return
		}
	
		go subscribeToTopic(sock, "TestTopic")
	
		for {
	
		}
	
	}
	
	//Subscribes to a specific topic
	func subscribeToTopic(sock mangos.Socket, topic string) {
	
		var msg []byte
		// Empty byte array effectively subscribes to everything
		err := sock.SetOption(mangos.OptionSubscribe, []byte(topic))
		if err != nil {
			log.Printf("Error Subscribing at subscribeToTopic:  %v", err.Error())
			return
		}
	
		if msg, err = sock.Recv(); err != nil {
			log.Printf("Error Recieving data at subscribeToTopic:  %v", err.Error())
			return
		}
	
		msgTopic := strings.Replace(string(msg), topic+"|", "", -1)
	
		log.Printf(string(msgTopic))
		go subscribeToTopic(sock, topic)
	}

