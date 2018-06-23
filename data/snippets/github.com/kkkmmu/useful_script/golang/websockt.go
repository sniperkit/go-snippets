// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

/*
The websocket protocol has two parts: a handshake and the data transfer.

The handshake from the client looks as follows:

GET /chat HTTP/1.1
Host: server.example.com
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
Origin: http://example.com
Sec-WebSocket-Protocol: chat, superchat
Sec-WebSocket-Version: 13

The handshake from the server looks as follows:

HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
Sec-WebSocket-Protocol: chat

The leading line from the client follows the Request-Line format.
The leading line from the server follows the Status-Line format.
An unordered set of header fields comes after the leading line in both cases.

Once the client and server have both sent their handshakes, and if
the handshake was successful, then the data transfer part starts.
This is a two-way communication channel where each side can, independently from the other, send data at will.
On the wire, a message is composed of one or more frames.
A frame has an associated type.  Each frame belonging to the same message contains the same type of data.

Broadly speaking, there are types for textual data (which is interpreted as UTF-8 [RFC3629]
text), binary data (whose interpretation is left up to the
application), and control frames (which are not intended to carry
data for the application but instead for protocol-level signaling,
such as to signal that the connection should be closed).  This
version of the protocol defines six frame types and leaves ten
reserved for future use.

The "Request-URI" of the GET method [RFC2616] is used to identify the
endpoint of the WebSocket connection, both to allow multiple domains
to be served from one IP address and to allow multiple WebSocket
endpoints to be served by a single server.

Additional header fields are used to select options in the WebSocket
Protocol.  Typical options available in this version are the
subprotocol selector (|Sec-WebSocket-Protocol|), list of extensions
support by the client (|Sec-WebSocket-Extensions|), |Origin| header
field, etc.  The |Sec-WebSocket-Protocol| request-header field can be
used to indicate what subprotocols (application-level protocols
layered over the WebSocket Protocol) are acceptable to the client.
The server selects one or none of the acceptable protocols and echoes
that value in its handshake to indicate that it has selected that
protocol.

Any status code other than 101 indicates that the WebSocket handshake
has not completed and that the semantics of HTTP still apply.  The
headers follow the status code.

The |Connection| and |Upgrade| header fields complete the HTTP
Upgrade.  The |Sec-WebSocket-Accept| header field indicates whether
the server is willing to accept the connection.  If present, this
header field must include a hash of the client's nonce sent in
|Sec-WebSocket-Key| along with a predefined GUID.  Any other value
must not be interpreted as an acceptance of the connection by the
server.

After sending a control frame indicating the connection should be
closed, a peer does not send any further data; after receiving a
control frame indicating the connection should be closed, a peer
discards any further data received.
*/

/*
Conceptually, WebSocket is really just a layer on top of TCP that
does the following:

o  adds a web origin-based security model for browsers

o  adds an addressing and protocol naming mechanism to support
multiple services on one port and multiple host names on one IP
address

o  layers a framing mechanism on top of TCP to get back to the IP
packet mechanism that TCP is built on, but without length limits

o  includes an additional closing handshake in-band that is designed
to work in the presence of proxies and other intermediaries
*/

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	filePeriod = 10 * time.Second
)

var (
	addr      = flag.String("addr", ":8080", "http service address")
	homeTempl = template.Must(template.New("").Parse(homeHTML))
	filename  string
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func readFileIfModified(lastMod time.Time) ([]byte, time.Time, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, lastMod, err
	}
	if !fi.ModTime().After(lastMod) {
		return nil, lastMod, nil
	}
	p, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fi.ModTime(), err
	}
	return p, fi.ModTime(), nil
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

func writer(ws *websocket.Conn, lastMod time.Time) {
	lastError := ""
	pingTicker := time.NewTicker(pingPeriod)
	fileTicker := time.NewTicker(filePeriod)
	defer func() {
		pingTicker.Stop()
		fileTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case <-fileTicker.C:
			var p []byte
			var err error

			p, lastMod, err = readFileIfModified(lastMod)

			if err != nil {
				if s := err.Error(); s != lastError {
					lastError = s
					p = []byte(lastError)
				}
			} else {
				lastError = ""
			}

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
			log.Println(err)
		}
		return
	}

	var lastMod time.Time
	if n, err := strconv.ParseInt(r.FormValue("lastMod"), 16, 64); err == nil {
		lastMod = time.Unix(0, n)
	}

	go writer(ws, lastMod)
	reader(ws)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	p, lastMod, err := readFileIfModified(time.Time{})
	if err != nil {
		p = []byte(err.Error())
		lastMod = time.Unix(0, 0)
	}
	var v = struct {
		Host    string
		Data    string
		LastMod string
	}{
		r.Host,
		string(p),
		strconv.FormatInt(lastMod.UnixNano(), 16),
	}
	homeTempl.Execute(w, &v)
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("filename not specified")
	}
	filename = flag.Args()[0]
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

const homeHTML = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>WebSocket Example</title>
    </head>
    <body>
        <pre id="fileData">{{.Data}}</pre>
        <script type="text/javascript">
            (function() {
                var data = document.getElementById("fileData");
                var conn = new WebSocket("ws://{{.Host}}/ws?lastMod={{.LastMod}}");
                conn.onclose = function(evt) {
                    data.textContent = 'Connection closed';
                }
                conn.onmessage = function(evt) {
                    console.log('file updated');
                    data.textContent = evt.data;
                }
            })();
        </script>
    </body>
</html>
`
