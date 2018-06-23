package pages

import (
	"bitodd/model"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const websocketURL = "/ws"

func websocketHandler(w http.ResponseWriter, r *http.Request) {

	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println("Error while upgrading ws connection: ", err)
		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(w, "Not a websocket handshake", http.StatusBadRequest)
		} else if err != nil {
			http.Error(w, "Error", http.StatusBadRequest)
		}
		return
	}

	c := &model.Connection{Send: make(chan *model.Envelope, 256), Ws: ws}
	model.Hub.Register <- c
	defer func() { model.Hub.Unregister <- c }()
	go c.Writer()
	// Block this handler so that the websocket connection remains open
	c.Reader()
}
