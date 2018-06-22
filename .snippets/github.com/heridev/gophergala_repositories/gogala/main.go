package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julien/gogala/lib"
	"golang.org/x/net/websocket"
)

const (
	defaultName = "U-"
)

var (
	listenAddr = flag.String("addr", os.Getenv("PORT"), "Listen address")
	clients    = make(map[string]lib.Client)
	debug      lib.Debug
	verbose    bool
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "Debug mode")

	if *listenAddr == "" {
		*listenAddr = "8080"
	}
}

func main() {
	flag.Parse()

	debug = lib.Debug(verbose)

	http.Handle("/", indexHandler())
	http.Handle("/static/", lib.GZipHandler(lib.CacheHandler(30, staticHandler())))
	http.Handle("/ws", websocket.Handler(wsHandler))

	debug.Printf("Listening on: %s\n", *listenAddr)
	log.Fatal(http.ListenAndServe(":"+*listenAddr, nil))
}

func indexHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/", http.StatusMovedPermanently)
	})
}

func staticHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
}

func wsHandler(ws *websocket.Conn) {
	registerClient(ws)

	for {
		var msg lib.Message
		var out lib.Message

		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			debug.Printf("Error reading message: %s\n", err)
			unregisterClient(ws)
			return
		}

		debug.Printf("Received message: %s\n", msg.Kind)

		switch msg.Kind {
		case "format":
			data, err := lib.Format([]byte(msg.Body))
			if err != nil {
				debug.Printf("Format Error: %s\n", err)

				out = lib.Message{
					Kind: "error",
					Body: err.Error(),
				}

				if err := sendToAll(ws, out); err != nil {
					debug.Printf("Error sending message: %s\n", err)
				}

			}

			if c := getClient(ws); c != nil {
				out = lib.Message{
					Kind: "code",
					Body: string(data),
					Args: lib.MakeArgs(c.Name),
				}
				if err := sendToAll(ws, out); err != nil {
					debug.Printf("Error sending message: %s\n", err)
				}
			}

		case "save":
			data, err := lib.CreateGist("GoGist", msg.Body)
			if err != nil {
				debug.Printf("Error creating gist: %v\n", err)
			}

			resp, err := lib.ParseResponse(data)
			if err != nil {
				debug.Printf("Error parsing Gists response: %s\n", err)
			}

			s := fmt.Sprintf("%s", resp["html_url"])

			out = lib.Message{
				Kind: "gist",
				Body: s,
			}

			if err := sendToAll(ws, out); err != nil {
				debug.Printf("Error sending message: %s\n", err)
			}

		case "compile":
			data, err := lib.Compile(msg.Body)
			if err != nil {
				debug.Printf("Error compiling code (remote): %s\n", err)
			}

			cr, err := lib.ParseCompileResponse(data)
			if err != nil {
				debug.Printf("Error parsing compile response: %s\n", err)
			}

			// WTF!
			if s := cr.Message(); s != nil {
				if ok := s.(string); ok != "" {
					out = lib.Message{
						Kind: "stdout",
						Body: s.(string),
					}

					sendToAll(ws, out)
				}
			}

		case "chat":
			t := time.Now().Format(time.Kitchen)

			if c := getClient(ws); c != nil {
				out = lib.Message{
					Kind: "chat",
					Body: lib.AppendString("[", t, "]", c.Name, ": ", msg.Body),
				}
				sendToAll(ws, out)
			}

		case "update":
			if c := getClient(ws); c != nil {
				out = lib.Message{
					Kind: "update",
					Body: msg.Body,
					Args: lib.MakeArgs(c.Name),
				}

				if err := sendToOthers(ws, out); err != nil {
					debug.Printf("Error sending message: %s\n", err)
				}
			}
		}

	}
}

func registerClient(ws *websocket.Conn) {

	n := len(clients)
	u := strconv.Itoa(n)
	if n < 10 {
		u = "0" + u
	}

	// Store client in clients map
	clients[u] = lib.Client{
		Id:   u,
		Name: defaultName + u,
		Conn: ws,
	}

	// Send welcome message
	var msg lib.Message

	msg = lib.Message{
		Kind: "info",
		Body: lib.AppendString("[", lib.PrintTimeStamp(), "] ", "Welcome, ", clients[u].Name),
		Args: lib.MakeArgs(clients[u].Name),
	}

	if err := sendToClient(ws, msg); err != nil {
		debug.Printf("Error sending message: %s\n", err)
	}

	msg = lib.Message{
		Kind: "info",
		Body: lib.AppendString("[", lib.PrintTimeStamp(), "] ", clients[u].Name, " joined"),
	}

	if err := sendToOthers(ws, msg); err != nil {
		debug.Printf("Error sending mesaage to others: %s\n", err)
	}
}

func unregisterClient(ws *websocket.Conn) {
	var n string

	// Remove client
	for k, v := range clients {
		if v.Conn == ws {
			n = k
			delete(clients, k)
		}
	}

	s := false
	if len(clients) == 1 {
		s = true
	}

	// Notify
	msg := lib.Message{
		Kind: "leave",
		Body: n,
		Args: lib.MakeArgs(s),
	}

	sendToOthers(ws, msg)
}

func getClient(ws *websocket.Conn) *lib.Client {
	var c *lib.Client
	for _, v := range clients {
		if v.Conn == ws {
			return &v
		}
	}
	return c
}

func sendToClient(ws *websocket.Conn, msg lib.Message) error {
	if err := websocket.JSON.Send(ws, msg); err != nil {
		return err
	}
	return nil
}

func sendToOthers(ws *websocket.Conn, msg lib.Message) error {
	for _, v := range clients {
		if v.Conn != ws {
			return sendToClient(v.Conn, msg)
		}
	}
	return nil
}

func sendToAll(ws *websocket.Conn, msg lib.Message) error {

	if err := sendToClient(ws, msg); err != nil {
		debug.Printf("Error sending message to client: %s\n", err)
		return err
	}

	if err := sendToOthers(ws, msg); err != nil {
		debug.Printf("Error sending message to others: %s\n", err)
		return err
	}

	return nil
}
