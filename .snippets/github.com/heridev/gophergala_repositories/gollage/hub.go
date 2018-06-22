// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type SocketData struct {
	Message []byte
	Channel string
}

type RegisterConn struct {
	Conn    *connection
	Channel string
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections, organized by channel
	connections map[string]map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan SocketData

	// Register requests from the connections.
	register chan RegisterConn

	// Unregister requests from connections.
	unregister chan RegisterConn
}

var h = hub{
	broadcast:   make(chan SocketData),
	register:    make(chan RegisterConn),
	unregister:  make(chan RegisterConn),
	connections: make(map[string]map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			if _, ok := h.connections[c.Channel]; !ok {
				h.connections[c.Channel] = make(map[*connection]bool)
			}
			h.connections[c.Channel][c.Conn] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c.Channel][c.Conn]; ok {
				delete(h.connections[c.Channel], c.Conn)
				close(c.Conn.send)
			}
		case m := <-h.broadcast:
			for c := range h.connections[m.Channel] {
				select {
				case c.send <- m.Message:
				default:
					close(c.send)
					delete(h.connections[m.Channel], c)
				}
			}
		}
	}
}
