package main

import (
	"golang.org/x/net/websocket"
)

type hub struct {
	connections map[*connection]bool
	register    chan *connection
	unregister  chan *connection
}

var h = hub{
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
		}
	}
}

type connection struct {
	ws   *websocket.Conn
	send chan string
}

func (c *connection) reader() {
	for {
		var name string
		err := websocket.Message.Receive(c.ws, &name)
		if err != nil {
			break
		}
		ExecuteScript(name, c.send)
	}
}

func (c *connection) writer() {
	for message := range c.send {
		err := websocket.Message.Send(c.ws, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

func wsHandler(ws *websocket.Conn) {
	c := &connection{send: make(chan string, 256), ws: ws}
	h.register <- c
	defer func() { h.unregister <- c }()
	go c.writer()
	c.reader()
}
