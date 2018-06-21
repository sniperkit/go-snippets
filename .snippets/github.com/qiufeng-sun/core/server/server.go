package server

import (
	"time"

	"util/logs"

	ctime "core/time"
)

// 服务器每次更新后sleep时间
const X_ServerSleep time.Duration = 10 * time.Millisecond

// server接口
type IServer interface {
	Init() bool
	Update()
	Destroy()
	PreQuit()
	String() string
}

//
type Server struct {
}

//
func (this *Server) Init() bool {
	return false
}

//
func (this *Server) Update() {}

//
func (this *Server) Destroy() {}

//
func (this *Server) PreQuit() {}

//
func (this Server) String() string {
	return "server"
}

// wrap server
type WrapServer struct {
	srv  IServer
	quit chan bool // 阻塞chan
}

// init
func (s *WrapServer) init() bool {
	logs.Infoln(s.srv, "init...")

	if s.srv.Init() {
		logs.Infoln(s.srv, "init ok.")

		return true
	}

	logs.Infoln(s.srv, "init failed!")

	return false
}

// run -- main loop
func (s *WrapServer) run() {
	logs.Infoln(s.srv, "running...")
	defer logs.Infoln(s.srv, "run end.")

	// signal
	chSig := WatchSignal()

	for {
		select {
		case strSig := <-chSig:
			logs.Infoln(s.srv, "receive signal:", strSig)

			s.srv.PreQuit()
			return

		case <-s.quit:
			logs.Infoln(s.srv, "run quit...")

			s.srv.PreQuit()
			return

		default:
			ctime.Update()
			s.srv.Update()
			time.Sleep(X_ServerSleep)
		}
	}
}

// destroy
func (s *WrapServer) destroy() {
	logs.Infoln(s.srv, "destroy...")
	defer logs.Infoln(s.srv, "destroy end.")

	s.srv.Destroy()
	close(s.quit)
}

// stop
func (s *WrapServer) stop() {
	logs.Infoln(s.srv, "stop...")
	defer logs.Infoln(s.srv, "stop end.")

	s.quit <- true
	<-s.quit
}

// new server
func newServer(s IServer) *WrapServer {
	return &WrapServer{quit: make(chan bool), srv: s}
}

// server obj
var g_server *WrapServer

// run server
func Run(s IServer) {
	g_server = newServer(s)

	if g_server.init() {
		g_server.run()
	}
	g_server.destroy()
}

// stop server
func Stop() {
	g_server.stop()
}
