package server

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/websocket"
	"net/http"
	"text/template"
	"time"
)

const (
	// Time allowed to write the nodes to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll nodes for changes with this period.
	nodesPeriod = 10 * time.Second
)

var (
	nodesH    map[string]string
	edgesH    map[string][]string
	homeTempl = template.Must(template.New("").Parse(homeHTML))
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func Run(nodes *map[string]string, edges *map[string][]string) {
	mux := http.NewServeMux()

	nodesH = *nodes
	edgesH = *edges
	mux.HandleFunc("/", serveHome)
	mux.HandleFunc("/ws", serveWs)

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3007")
}

func serveHome(rw http.ResponseWriter, r *http.Request) {

	//TODO format nodes and edges to json to use with js graph library
	p := []byte(fmt.Sprintf("nodes %v\n", nodesH))
	//rw.Write([]byte(fmt.Sprintf("edges %v\n", edgesH)))

	var v = struct {
		Host string
		Data string
	}{
		r.Host,
		string(p),
	}
	homeTempl.Execute(rw, &v)
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func writer(ws *websocket.Conn) {
	pingTicker := time.NewTicker(pingPeriod)
	nodesTicker := time.NewTicker(nodesPeriod)
	defer func() {
		pingTicker.Stop()
		nodesTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case <-nodesTicker.C:

			p := []byte(fmt.Sprintf("%v", nodesH))

			if p != nil {
				ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err := ws.WriteMessage(websocket.TextMessage, p); err != nil {
					return
				}
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			fmt.Println(err)
		}
		return
	}

	go writer(ws)
	reader(ws)
}

const homeHTML = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>Honeybee panel</title>
    </head>
    <body>
        <pre id="data">{{.Data}}</pre>
        <script type="text/javascript">
            (function() {
                var data = document.getElementById("data");
                var conn = new WebSocket("ws://{{.Host}}/ws");
                conn.onclose = function(evt) {
                    data.textContent = 'Connection closed';
                }
                conn.onmessage = function(evt) {
                    console.log('data updated');
                    data.textContent = evt.data;
                }
            })();
        </script>
    </body>
</html>
`
