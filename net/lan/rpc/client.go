package rpc

import (
	"sync"
	"math/rand"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/req"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"

	"util/logs"

	"core/net/lan"
	"errors"
)

var _ = logs.Debug

//
var ErrEmptyClient = errors.New("empty rpc client!")

//
type Client struct {
	mangos.Socket

	SrvAddr string
}

func NewClient(srvAddr string) *Client {
	//
	sock, _ := req.NewSocket()

	//
	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())

	//
	srvAddr = lan.FormatTcpAddr(srvAddr)
	if e := sock.Dial(srvAddr); e != nil {
		logs.Warnln(e)
		return nil
	}

	return &Client{Socket: sock, SrvAddr: srvAddr}
}

func (this *Client) Close() {
	this.Socket.Close()
}

func (this *Client) Call(raw []byte) ([]byte, error) {
	e := this.Send(raw)
	if e != nil {
		return  nil, e
	}
	return this.Recv()
}

//
type PoolConfig struct {
	Name                     string
	InitNum, IdleNum, MaxNum int
}

func (this PoolConfig) Check() bool {
	if this.InitNum < 0 || this.InitNum > this.IdleNum || this.IdleNum > this.MaxNum {
		return false
	}
	return true
}

//
type ClientPool struct {
	*PoolConfig

	sync.Mutex
	SrvAddrs   map[string]bool
	SrvAddrArr []string

	IdleClients []*Client
	CurNum int
}

func NewClientPool(cfg *PoolConfig) *ClientPool {
	if !cfg.Check() {
		logs.Panicln("invalid client pool config:%+v", cfg)
		return nil
	}

	p := &ClientPool{PoolConfig:cfg}

	return p
}

// update server addrs
func (this *ClientPool) Update(srv string, addrs []string) {
	logs.Info("%v update servers! %v:%v, old:%v", this.Name, srv, addrs, this.SrvAddrArr)

	this.Lock()
	defer this.Unlock()
	defer func() {
		logs.Info("%v addrs -- array:%v, map:%v", this.Name, this.SrvAddrArr, this.SrvAddrs)
	}()

	// unexpected output -- params be push to stack, print the old addr params at the end of call
	//defer logs.Info("%v addrs -- array:%v, map:%v", this.Name, this.SrvAddrArr, this.SrvAddrs)

	//
	if len(addrs) == 0 {
		this.SrvAddrs = nil
		this.SrvAddrArr = nil
		return
	}

	//
	nm := make(map[string]bool, len(addrs))
	for _, addr := range addrs {
		nm[addr] = true
	}
	this.SrvAddrs = nm
	this.SrvAddrArr = addrs

	//
	if this.IdleClients != nil {
		return
	}

	this.IdleClients = make([]*Client, this.InitNum, this.IdleNum)
	for i := 0; i < this.InitNum; i++ {
		addr := this.randAddr()
		this.IdleClients[i] = NewClient(addr)
	}
}

// get client
func (this *ClientPool) GetClient() *Client {
	this.Lock()
	defer this.Unlock()

	var client *Client = nil
	num := len(this.IdleClients)
	if num > 0 {
		num--
		client = this.IdleClients[num]
		this.IdleClients = this.IdleClients[:num]
		if client != nil && !this.checkClient(client) {
			client.Close()
			client = nil
		}
	}

	if nil == client && this.CurNum < this.MaxNum {
		client = this.NewClient()
		if client != nil {
			this.CurNum++
		}
	}

	return client
}

// return client
func (this *ClientPool) ReturnClient(client *Client) {
	if nil == client {
		return
	}

	this.Lock()
	defer this.Unlock()

	ok := this.checkClient(client)
	num := len(this.IdleClients)
	if !ok || num >= this.IdleNum {
		client.Close()
		this.CurNum--
		return
	}

	this.IdleClients = append(this.IdleClients, client)
}

// select random client
func (this *ClientPool) randAddr() string {
	addrs := this.SrvAddrArr
	num := len(addrs)
	if num <= 0 {
		return ""
	}
	addr := addrs[0]
	if num > 1 {
		addr = addrs[rand.Intn(num)]
	}
	return addr
}

func (this *ClientPool) NewClient() *Client {
	addr := this.randAddr()
	if "" == addr {
		return nil
	}
	return NewClient(addr)
}

func (this *ClientPool) checkClient(client *Client) bool {
	if nil == client {
		logs.Panicln("please check nil before call this function!")
	}
	_, ok := this.SrvAddrs[client.SrvAddr]
	return ok
}

func (this *ClientPool) Call(raw []byte) ([]byte, error) {
	c := this.GetClient()
	if nil == c {
		return nil, ErrEmptyClient
	}
	defer this.ReturnClient(c)

	return c.Call(raw)
}
