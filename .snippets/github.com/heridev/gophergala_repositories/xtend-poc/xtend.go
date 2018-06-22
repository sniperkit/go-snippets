package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/websocket"
)

var players []Player

type Event struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

type Player struct {
	Name  string          `json:"name"`
	Conn  *websocket.Conn `json:"-"`
	X     int             `json:"x"`
	Y     int             `json:"y"`
	Color string          `json:"color"`
}

func broadcast(evt Event, players []Player) {
	for _, p := range players {
		if err := websocket.JSON.Send(p.Conn, evt); err != nil {
			log.Fatal(err)
		}
	}
}

func Start(ws *websocket.Conn) {
	var event Event
	for websocket.JSON.Receive(ws, &event) != io.EOF {

		switch event.Action {
		case "new":
			var player Player
			err := mapstructure.Decode(event.Data, &player)
			if err != nil {
				log.Fatal(err)
			}
			player.Conn = ws
			players = append(players, player)
			fmt.Println(player.Name, " logged in")

			if len(players) >= 2 {
				fmt.Println("Total players is ", len(players), " game start")

				players[0].X = 100
				players[0].Y = 300
				players[0].Color = "0xFFFF0B"
				players[1].X = 700
				players[1].Y = 300
				players[1].Color = "0xBC1C22"

				event := Event{
					Action: "render_base",
					Data: map[string]interface{}{
						"players": players,
					},
				}

				broadcast(event, players)
			} else {
				websocket.JSON.Send(player.Conn, Event{Action: "init"})
			}
		case "":
		}

	}
}

func optionalEnv(key, defaultValue string) string {
	var v = os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}

func main() {
	http.Handle("/api/start", websocket.Handler(Start))

	staticPath := optionalEnv("STATIC_PATH", "web")
	http.Handle("/", http.FileServer(http.Dir(staticPath)))

	port := optionalEnv("PORT", "12345")
	log.Println("Listen on port", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
