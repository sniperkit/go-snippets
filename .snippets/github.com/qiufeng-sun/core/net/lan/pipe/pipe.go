/**
* 服务端有一个用来接受其他服务器消息的管道，和多个向其他服务器发送消息的管道
**/
package pipe

import (
	"errors"
	"math/rand"
	"sync"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/pull"
	"github.com/go-mangos/mangos/protocol/push"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"

	"util/logs"

	"core/net"
	. "core/net/lan"
)

var _ = logs.Debug

// 接收其他服务器消息
type Server struct {
	*LanCfg
	mangos.Socket

	SrvUrl string
}

func NewServer(cfg *LanCfg) *Server {
	//
	sock, _ := pull.NewSocket()

	//
	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	if e := sock.Listen(cfg.Addr); e != nil {
		logs.Panicln(e)
	}

	//
	logs.Info("server<%+v> start listen lan connect!", cfg)

	//
	srvUrl := net.GenUrl(cfg.ServerId(), "0")

	return &Server{LanCfg: cfg, Socket: sock, SrvUrl: srvUrl}
}

func (this *Server) Recv() ([]byte, error) {
	return this.Socket.Recv()
}

func (this *Server) Close() {
	this.Socket.Close()
}

func (this Server) GetUrl() string {
	return this.SrvUrl
}

//
type Client struct {
	*LanCfg
	mangos.Socket

	clear bool
}

func NewClient(cfg *LanCfg) *Client {
	//
	sock, _ := push.NewSocket()

	//
	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	if e := sock.Dial(cfg.Addr); e != nil {
		logs.Panicln(e)
	}

	return &Client{LanCfg: cfg, Socket: sock}
}

func (this *Client) Send(msg []byte) error {
	return this.Socket.Send(msg)
}

func (this *Client) Close() {
	this.Socket.Close()
}

// 向其他服务器发送消息
type Clients struct {
	SenderId    string               // 发送端服务器id
	NameClients map[string][]*Client // 接收端连接组serverName=>[]
	IdClient    map[string]*Client   // 接收端连接serverId=>*

	sync.Mutex
}

func NewClients(senderId string) *Clients {
	return &Clients{
		SenderId:    senderId,
		NameClients: map[string][]*Client{},
		IdClient:    map[string]*Client{},
	}
}

// update server addrs
func (this *Clients) Update(srv string, addrs []string) {
	logs.Info("%v update servers! %v:%v", this.SenderId, srv, addrs)

	this.Lock()
	defer this.Unlock()

	olds := this.NameClients[srv]
	delete(this.NameClients, srv)

	// set close 2 old 1st
	for _, v := range olds {
		v.clear = true
	}

	// update
	for _, addr := range addrs {
		cfg := NewLanCfg(srv, addr)
		id := cfg.ServerId()
		client, ok := this.IdClient[id]
		if ok {
			// use old
			client.clear = false
			this.NameClients[srv] = append(this.NameClients[srv], client)
		} else {
			// new one
			client := NewClient(cfg)
			this.NameClients[srv] = append(this.NameClients[srv], client)
			this.IdClient[id] = client
		}
	}

	// close invalid
	for _, v := range olds {
		if !v.clear {
			continue
		}

		delete(this.IdClient, v.ServerId())
		v.Close()
	}

	// log
	logs.Info("%v update name clients:%v", this.SenderId, this.NameClients)
	logs.Info("%v update id client:%v", this.SenderId, this.IdClient)
}

// select random client
func (this *Clients) SelectRand(srv string) string {
	srvs := this.NameClients[srv]
	num := len(srvs)
	if num <= 0 {
		return ""
	}
	s := srvs[0]
	if num > 1 {
		s = srvs[rand.Intn(num)]
	}
	return s.ServerId()
}

// send msg
func (this *Clients) SendMsg(srvId string, d []byte) error {
	// find
	c, ok := this.IdClient[srvId]
	if !ok {
		return errors.New("server not found! id:" + srvId)
	}

	return c.Send(d)
}

//
type Lan struct {
	*Server
	*Clients
}

func NewLan(cfg *LanCfg) *Lan {
	return &Lan{
		Server:  NewServer(cfg),
		Clients: NewClients(cfg.ServerId()),
	}
}
