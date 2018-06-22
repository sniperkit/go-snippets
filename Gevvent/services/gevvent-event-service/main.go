package main

import (
	"log"

	"github.com/asim/go-micro/server"

	"github.com/gophergala/Gevvent/services/gevvent-event-service/handler"
	"github.com/gophergala/Gevvent/services/gevvent-lib/monitoring"
)

func main() {
	server.Name = "gevvent-event-service"

	monitoring.Init(server.Name)
	server.Init()

	// Register Handlers
	server.Register(server.NewReceiver(new(handler.Attendees)))
	server.Register(server.NewReceiver(new(handler.Create)))
	server.Register(server.NewReceiver(new(handler.Delete)))
	server.Register(server.NewReceiver(new(handler.Invite)))
	server.Register(server.NewReceiver(new(handler.List)))
	server.Register(server.NewReceiver(new(handler.Newest)))
	server.Register(server.NewReceiver(new(handler.Read)))
	server.Register(server.NewReceiver(new(handler.ReadUser)))
	server.Register(server.NewReceiver(new(handler.RSVP)))
	server.Register(server.NewReceiver(new(handler.Search)))

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
