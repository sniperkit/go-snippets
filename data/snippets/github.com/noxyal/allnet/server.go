package allnet

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/noxyal/allsoc"
	"golang.org/x/net/websocket"
)

// server stores a ptr to the Server struct
var server *Server

// onConnect is a callback called on each new connection
var onConnect func(*allsoc.Socket)

// Server example fields: Server{":9090", ":8080", "public/fileserver"}
type Server struct {
	TCPPort, HTTPPort, StaticPath string
	OnConnect                     func(*allsoc.Socket)
}

// StartServer starts server and listens. Default settings: Server{":9090", ":8080", "public/fileserver"}
func StartServer(s *Server) {

	if s == nil {
		s = &Server{}
	}
	if s.TCPPort == "" {
		s.TCPPort = ":9090"
	}
	if s.HTTPPort == "" {
		s.HTTPPort = ":8080"
	}
	if s.StaticPath == "" {
		s.StaticPath = "public/fileserver"
	}
	if s.OnConnect == nil {
		s.OnConnect = newConn
	}

	allsoc.SetupAllsoc()
	server = s
	onConnect = s.OnConnect

	go netListen(s.TCPPort)
	http.Handle("/ws", websocket.Handler(wsHandler))
	http.Handle("/", http.FileServer(http.Dir(s.StaticPath)))
	err := http.ListenAndServe(s.HTTPPort, nil)
	if err != nil {
		fmt.Println("HTTP Error:", err.Error())
	}
}

func netListen(port string) {
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("TCP Error:", err.Error())
		return
	}
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("TCP Error:", err.Error())
		} else {
			go socketHandler(c)
		}
	}
}
func wsHandler(ws *websocket.Conn) {
	socketHandler(ws)
}
func socketHandler(iorw io.ReadWriter) {
	rw := &iorw
	soc := allsoc.NewSocket(rw)
	soc.Join("")
	onConnect(soc)
}

// newConn is the default OnConnect. It is an echo and broadcast server.
func newConn(soc *allsoc.Socket) {
	fmt.Println("a user connected!")
	msg := make([]byte, 32*1024)
	for {
		n, err := soc.Read(msg)
		if err != nil {
			return
		}
		if n > 0 {
			soc.Broadcast("", msg[0:n])
			soc.Write(msg[0:n])
		}
	}
}
