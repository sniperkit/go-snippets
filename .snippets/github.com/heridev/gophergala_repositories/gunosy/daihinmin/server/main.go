package main

import (
	"log"
	"net/http"

	"code.google.com/p/go.net/websocket"
	"encoding/json"
	dh "github.com/gophergala/gunosy/daihinmin"
)

func main() {
	bind := "localhost:3000"

	// Create default game
	m := dh.NewMatch("New Game 1")
	log.Printf("Created the initial game: %s\n", m)
	http.Handle("/ws", websocket.Handler(handle))
	http.HandleFunc("/match", handleMatch)

	log.Println("running:", bind)
	http.ListenAndServe(bind, nil)
}

func handleMatch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		js, err := json.Marshal(dh.ListMatches())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func handle(ws *websocket.Conn) {
	c := dh.NewClient(ws)
	go c.Write()
	c.Run()
}
