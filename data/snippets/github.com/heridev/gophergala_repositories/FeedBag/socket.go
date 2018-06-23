package main

import (
	"log"

	"github.com/googollee/go-socket.io"
)

func SetupSocketIO() (*socketio.Server, error) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}

	checkErr(err, "Problem making socket.io server:")
	server.On("connection", func(so socketio.Socket) {
		log.Println("New socket.io connection:", so.Id())
		so.Join("feedbag")
		so.On("disconnection", func() {
			// no op
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("Socket.io server error:", err.Error())
	})

	return server, nil
}

func StartSocketPusher(s *socketio.Server, c <-chan []Activity) error {
	go func() {
		for {
			activities := <-c
			s.BroadcastTo("feedbag", "activity", activities)
		}
	}()

	return nil
}
