package model

import (
	"github.com/gorilla/websocket"
	"log"
)

type Connection struct {
	// The websocket connection.
	Ws *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan *Envelope
}

func (c *Connection) Reader() {
	for {
		msg := &Envelope{}
		err := c.Ws.ReadJSON(msg)
		if err != nil {
			log.Println("Error while reading json", err)
			break
		}
	}
	c.Ws.Close()
}

func (c *Connection) Writer() {
	for message := range c.Send {
		err := c.Ws.WriteJSON(message)
		if err != nil {
			log.Println("Error while writing to websocket", err)
			break
		}
	}
	c.Ws.Close()
}
