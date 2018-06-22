package main

import (
	"log"

	"github.com/asim/go-micro/server"

	"github.com/gophergala/Gevvent/services/gevvent-lib/monitoring"
	"github.com/gophergala/Gevvent/services/gevvent-user-service/handler"
)

func main() {
	server.Name = "gevvent-user-service"

	monitoring.Init(server.Name)
	server.Init()

	// Register Handlers
	server.Register(server.NewReceiver(new(handler.Authorised)))
	server.Register(server.NewReceiver(new(handler.Login)))
	server.Register(server.NewReceiver(new(handler.Logout)))
	server.Register(server.NewReceiver(new(handler.ReadUser)))
	server.Register(server.NewReceiver(new(handler.Register)))

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
