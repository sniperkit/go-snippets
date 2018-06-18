package main

//http://talks.golang.org/2012/chat.slide#29

//var sock = new WebSocket("ws://localhost:4000/");
//sock.onmessage = function(m) { console.log("Received:", m.data); }
//sock.send("Hello!\n")

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
)

func main() {
	http.Handle("/", websocket.Handler(handler))
	http.ListenAndServe("localhost:4000", nil)
}

func handler(c *websocket.Conn) {
	var s string
	fmt.Fscan(c, &s)
	fmt.Println("Received:", s)
	fmt.Fprint(c, "How do you do?")
}
