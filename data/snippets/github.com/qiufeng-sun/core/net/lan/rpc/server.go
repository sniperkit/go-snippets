package rpc

import (
	"sync"
	"sync/atomic"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/rep"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"

	"util/logs"

	"core/net"
	. "core/net/lan"
)

var _ = logs.Debug

//
type handleFunc func([]byte) []byte

// 接收其他服务器消息
type Server struct {
	*LanCfg
	mangos.Socket
	SrvUrl string

	handleFunc
	goNums int

	stop bool
	wg   sync.WaitGroup

	SendFailed int64
	SendOk     int64
}

func NewServer(cfg *LanCfg) *Server {
	//
	sock, _ := rep.NewSocket()

	//
	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	if e := sock.Listen(cfg.Addr); e != nil {
		logs.Panicln(e)
	}

	//
	logs.Info("server<%#v> start rpc listen!", cfg)

	//
	srvUrl := net.GenUrl(cfg.ServerId(), "0")

	return &Server{LanCfg: cfg, Socket: sock, SrvUrl: srvUrl}
}

func (this *Server) Init(goNums int, h handleFunc) {
	this.goNums = goNums
	this.handleFunc = h
}

func (this *Server) recvAndProc() {
	defer this.wg.Done()

	for {
		m, e := this.RecvMsg()
		if mangos.ErrRecvTimeout == e {
			continue
		}
		if e != nil {
			logs.Warn("receive msg failed! error=%v", e)
			return
		}
		if this.stop {
			return
		}

		m.Body = this.handleFunc(m.Body)

		if e := this.SendMsg(m); e != nil {
			atomic.AddInt64(&this.SendFailed, 1)
			logs.Warn("send msg failed!")
		} else {
			atomic.AddInt64(&this.SendOk, 1)
		}
	}
}

func (this *Server) Serve() {
	num := this.goNums

	this.wg.Add(num)
	for i := 0; i < num; i++ {
		go this.recvAndProc()
	}
}

func (this *Server) ServeAndWait() {
	this.Serve()
	this.Wait()
}

func (this *Server) Wait() {
	this.wg.Wait()
}
