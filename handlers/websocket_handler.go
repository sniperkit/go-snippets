package handlers

import "golang.org/x/net/websocket"

type WebSocketHandler struct {
	ws *websocket.Conn
}

func (h *WebSocketHandler) Conn() *websocket.Conn {
	return h.ws
}

func (h *WebSocketHandler) Handler(c *websocket.Conn) {
	h.ws = c
	<-make(chan struct{})
}
