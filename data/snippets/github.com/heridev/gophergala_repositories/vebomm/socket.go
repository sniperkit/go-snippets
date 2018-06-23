package main

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gophergala/vebomm/core"

	"github.com/gorilla/websocket"
)

type hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan message

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection

	onlineUsers map[int64]*connection
}

type message struct {
	msg  []byte
	conn *connection
}

var h = hub{
	broadcast:   make(chan message),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
	onlineUsers: map[int64]*connection{},
}

func findPartner(cons core.Constraints, subj *core.User) *connection {
	minAgeDiff := 999
	var bestMatch *connection
	for _, c := range h.onlineUsers {
		if c.user.Id == subj.Id {
			continue
		}
		age := subj.Age()
		if age >= cons.MinAge && age <= cons.MaxAge &&
			(cons.Gender == -1 || cons.Gender == c.user.Gender) {
			diff := int(math.Abs(float64(age - subj.Age())))
			if diff < minAgeDiff {
				minAgeDiff = diff
				bestMatch = c
			}
		}
	}

	return bestMatch
}

func (h *hub) run() {
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
			if _, ok := h.onlineUsers[c.user.Id]; !ok {
				h.onlineUsers[c.user.Id] = c
			}

		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}

			if _, ok := h.onlineUsers[c.user.Id]; ok {
				delete(h.onlineUsers, c.user.Id)
			}
			c.ws.Close()

		case m := <-h.broadcast:
			println(string(m.msg))
			var cons core.Constraints
			err := json.Unmarshal(m.msg, &cons)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			result := findPartner(cons, m.conn.user)
			if result != nil {
				go func() {
					b1, _ := json.Marshal(m.conn.user)
					b2, _ := json.Marshal(result.user)
					m.conn.send <- b2
					result.send <- b1
					h.unregister <- m.conn
					h.unregister <- result
				}()
			}

		case <-ticker.C:
			for c := range h.connections {
				go func() {
					err := c.ws.WriteControl(websocket.PingMessage, []byte{0}, time.Now().Add(2*time.Second))
					if err != nil {
						h.unregister <- c
					}

					h.onlineUsers[c.user.Id] = c
				}()
			}
		}
	}
}

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	user *core.User
}

func (c *connection) reader() {
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		h.broadcast <- message{msg, c}
	}
	c.ws.Close()
}

func (c *connection) writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func wsHandler(user *core.User, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws, user: user}
	h.register <- c
	defer func() { h.unregister <- c }()
	go c.writer()
	c.reader()

}
