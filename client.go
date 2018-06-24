// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	// Careful with this one, if messages exceed this size it seems the default
	// behaviour is to close the connection.
	// TODO - properly handle messages that are too large client side
	// previous values: 2048
	maxMessageSize = 32786
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Room
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Infof("error: %v", err)
			}

			if err := c.conn.Close(); err != nil {
				log.Infof("close connection error: %s", err.Error())
			}
			break
		}

		c.HandleAction(message)

		// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		// c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// c.conn.WriteJSON()
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	w.Write(newline)
			// 	w.Write(<-c.send)
			// }

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *Client) HandleAction(data []byte) {
	action := struct {
		Type        string
		RequestId   string
		SilentError bool
		Data        json.RawMessage
	}{}
	if err := json.Unmarshal(data, &action); err != nil {
		log.Infof("error parsing action JSON: %s", err.Error())
		c.SendResponse(&ClientResponse{
			Type:  "PARSE_ERROR",
			Error: fmt.Sprintf("action parsing error type: %s", err.Error()),
		})
		return
	}
	// TODO - This looks a lot like a muxer...
	// if action.Type == "URL_ARCHIVE_REQUEST" {
	// 	act := struct {
	// 		Url string
	// 	}{}
	// 	if err := json.Unmarshal(action.Data, &act); err != nil {
	// 		c.SendResponse(&ClientResponse{
	// 			Type:  "PARSE_ERROR",
	// 			Error: fmt.Sprintf("action parsing error: %s", err.Error()),
	// 		})
	// 		return
	// 	}
	// 	c.ArchiveUrl(appDB, action.RequestId, act.Url)
	// } else

	if strings.HasSuffix(action.Type, "REQUEST") {
		log.Infof("%s: %s", action.RequestId, action.Type)
		c.HandleRequestAction(action.Type, action.RequestId, action.SilentError, action.Data)
		return
	}
	log.Infof("unrecognized action: %s", action.Type)
}

func (c *Client) SendResponse(res *ClientResponse) {
	// TODO - switch client to use "conn.SendJSON" for this stuff
	data, err := json.Marshal(res)
	if err != nil {
		// TODO - handle "internal server parsing error" here
		// sending a response
		log.Info(err.Error())
		return
	}
	c.send <- data
	// if err := c.conn.WriteJSON(res); err != nil {
	// 	log.Info(err.Error())
	// }
}

func (c *Client) HandleRequestAction(req string, reqId string, silentError bool, data json.RawMessage) {
	for _, t := range ClientReqActions {
		if t.Type() == req {
			res := t.Parse(reqId, data).Exec()
			res.SilentError = silentError
			c.SendResponse(res)
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Room, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Info(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	go client.writePump()
	client.readPump()
}
