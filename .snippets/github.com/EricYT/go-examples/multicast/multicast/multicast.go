package multicast

import (
	"errors"
	"log"
	"net"
	"sync"

	tomb "gopkg.in/tomb.v1"
)

const (
	UDP_PROTOCOL    string = "udp"
	maxDatagramSize int    = 8192
)

var (
	ErrorMulticastNotFound        error = errors.New("multicast: not found")
	ErrorMulticastHandlerExists   error = errors.New("multicast: handler already exists")
	ErrorMulticastHandlerNotFound error = errors.New("multicast: handler not found")
)

// default
var _center *multicastCenter

func init() {
	_center = new(multicastCenter)
	_center.multicasters = make(map[string]*multicast)
}

// multicastCenter package
type multicastCenter struct {
	mutex sync.Mutex // protecting remaining fields

	multicasters map[string]*multicast
}

func New(addr string) (Multicaster, error) {
	if _center == nil {
		log.Fatal("multicastCenter multicast not initialize")
	}
	return _center.NewMulticast(addr)
}

func Remove(addr string) error {
	if _center == nil {
		log.Fatal("multicastCenter multicast not initialize")
	}
	return _center.RemoveMulticast(addr)
}

func (m *multicastCenter) NewMulticast(addr string) (Multicaster, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if mul, ok := m.multicasters[addr]; ok {
		return mul, nil
	}
	mul, err := NewMulticast(addr)
	if err != nil {
		return nil, err
	}
	m.multicasters[mul.Addr()] = mul
	return mul, nil
}

func (m *multicastCenter) RemoveMulticast(addr string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	mul, ok := m.multicasters[addr]
	if !ok {
		return ErrorMulticastNotFound
	}
	delete(m.multicasters, addr)
	mul.Close()
	return nil
}

// handler
type Handler func(src *net.UDPAddr, count int, data []byte)
type HandlerKey string

type Multicaster interface {
	RegisterHandler(key HandlerKey, hnd Handler) error
	UnregisterHandler(key HandlerKey) error
	Notify(data []byte) (int, error)
}

type multicast struct {
	mutex sync.Mutex
	tomb  *tomb.Tomb

	addr     *net.UDPAddr
	handlers map[HandlerKey]Handler

	listener *net.UDPConn
	sender   *net.UDPConn
}

func NewMulticast(a string) (*multicast, error) {
	addr, err := net.ResolveUDPAddr(UDP_PROTOCOL, a)
	if err != nil {
		return nil, err
	}
	m := new(multicast)
	m.tomb = new(tomb.Tomb)
	m.addr = addr
	m.handlers = make(map[HandlerKey]Handler)
	// initialize sender
	c, err := net.DialUDP(UDP_PROTOCOL, nil, addr)
	if err != nil {
		return nil, err
	}
	m.sender = c
	// reciver
	go func() {
		defer m.tomb.Done()
		m.tomb.Kill(m.serve())
	}()
	return m, nil
}

func (m *multicast) Close() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.tomb.Kill(nil)
	m.cleanAll()
}

func (m *multicast) Addr() string {
	return m.addr.String()
}

func (m *multicast) RegisterHandler(key HandlerKey, hnd Handler) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.handlers[key]; ok {
		return ErrorMulticastHandlerExists
	}
	m.handlers[key] = hnd
	return nil
}

func (m *multicast) UnregisterHandler(key HandlerKey) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.handlers[key]; !ok {
		return ErrorMulticastHandlerNotFound
	}
	delete(m.handlers, key)
	return nil
}

func (m *multicast) cleanAll() {
	m.handlers = nil
}

func (m *multicast) serve() error {
	l, err := net.ListenMulticastUDP(UDP_PROTOCOL, nil, m.addr)
	if err != nil {
		return err
	}
	m.listener = l
	l.SetReadBuffer(maxDatagramSize)

	for {
		b := make([]byte, maxDatagramSize)
		n, src, err := l.ReadFromUDP(b)
		if err != nil {
			return err
		}

		// dispatch
		for _, hnd := range m.handlers {
			go hnd(src, n, b)
		}

		// judge
		select {
		case <-m.tomb.Dying():
			return nil
		default:
			// continue
		}
	}
}

func (m *multicast) Notify(data []byte) (int, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	n, err := m.sender.Write(data)
	if err != nil {
		return 0, err
	}
	return n, nil
}
