package model

import (
	"github.com/darkhelmet/twitterstream"
	"log"
)

type hub struct {
	// Registered connections.
	connections map[*Connection]bool

	// Inbound messages from the connections.
	broadcast chan *Envelope

	// Register requests from the connections.
	Register chan *Connection

	// Unregister requests from connections.
	Unregister chan *Connection
}

type msgContainer struct {
}

var Hub = hub{
	broadcast:   make(chan *Envelope),
	Register:    make(chan *Connection),
	Unregister:  make(chan *Connection),
	connections: make(map[*Connection]bool),
}

func (h *hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.registerConnection(c)
		case c := <-h.Unregister:
			h.closeConnection(c)
		case m := <-h.broadcast:
			h.handleMessage(m)
		}
	}
}

func (h *hub) registerConnection(c *Connection) {
	h.connections[c] = true

	log.Printf("Websocket register (count: %v)", len(h.connections))

	// Send info message to everybody
	postInfoMessage()
}

func (h *hub) closeConnection(c *Connection) {
	delete(h.connections, c)
	close(c.Send)

	log.Printf("Websocket close (count: %v)", len(h.connections))

	// Send info message to everybody
	postInfoMessage()
}

func (h *hub) handleMessage(msg *Envelope) {

	for c := range h.connections {
		select {
		case c.Send <- msg:
		default:
			delete(h.connections, c)
			close(c.Send)
			go c.Ws.Close()
		}
	}
}

func postInfoMessage() {
	// Create info payload
	info := &Info{UserCount: len(Hub.connections)}

	// Create message envelope
	Hub.handleMessage(&Envelope{Action: INFO, Info: info})
}

func PostLocationTweet(tweet *twitterstream.Tweet) {

	temp := &Tweet{Picture: tweet.User.ProfileImageUrl, DisplayName: tweet.User.Name, ScreenName: tweet.User.ScreenName, Message: tweet.Text, Lat: tweet.Coordinates.Lat.Float64(), Long: tweet.Coordinates.Long.Float64()}

	Hub.handleMessage(&Envelope{Action: TWEET, Tweet: temp})
}
