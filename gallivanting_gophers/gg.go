package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gophergala/gallivanting_gophers/board"
	"github.com/gophergala/gallivanting_gophers/data"
	"github.com/gophergala/gallivanting_gophers/webserver"
	"github.com/julienschmidt/httprouter"
)

var factory *board.Factory
var server *webserver.Server
var database *data.DB

// This is a temporary global gameboard. This board will be removed to handle
// session based boards for real clients. For now this is just to make the API
// functional.
var gameboard *board.FixBoard

func main() {
	fmt.Println("Gallivanting Gophers Starting...")

	database = data.NewDB()
	factory = board.NewFactory()

	gameboard = factory.GetBoard()

	server = webserver.NewServer(database)

	server.Start()

	server.RegisterGet("/board", getBoard)
	server.RegisterPost("/board", postBoard)
	server.RegisterPost("/move", postMove)

	fmt.Print("PRESS ENTER TO SHUTDOWN...")
	fmt.Scanln()

	fmt.Println("Gallivanting Gophers Shutdown.")
}

// This will be moved to link the gameboard to the session.
func getBoard(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	o, err := json.Marshal(&gameboard)
	if err != nil {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(200)
		w.Write(o)
	}
}

// This will be moved to link the gameboard to the session.
func postBoard(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	gameboard = factory.GetBoard()

	o, err := json.Marshal(&gameboard)
	if err != nil {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(200)
		w.Write(o)
	}
}

// This will be moved to link the gameboard to the session.
func postMove(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var m movePost
	var err error
	var o []byte

	err = decoder.Decode(&m)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	gameboard.MoveGopher(m.Gopher, m.Direction)

	o, err = json.Marshal(&gameboard)
	if err != nil {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(200)
		w.Write(o)
	}
}

type movePost struct {
	Gopher    byte `json:"gopher"`
	Direction int  `json:"direction"`
}
