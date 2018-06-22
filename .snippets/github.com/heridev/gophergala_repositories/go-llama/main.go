package main

import (
	"code.google.com/p/websocket"
	"flag"
	"fmt"
	"github.com/gophergala/go-llama/intchess"
	"net/http"
)

var (
	addr   = flag.String("addr", ":8080", "http service address")
	makeDB = flag.Bool("makeDB", false, "Create the databases for the system")
)

func main() {
	//this parses any command line arguments to the program
	flag.Parse()

	if err := intchess.ConnectToDatabase(); err != nil {
		fmt.Println("Could not connect to database: " + err.Error())
		return
	}

	if *makeDB {
		intchess.CreateDatabaseTables()
		return
	}

	intchess.StartGameServer()

	http.Handle("/ws", websocket.Handler(intchess.WsHandler))
	http.Handle("/", http.FileServer(http.Dir("./static/")))
	if err := http.ListenAndServe(*addr, nil); err != nil {
		fmt.Println("Fatal error: ListenAndServe:", err)
	}
}
