package cassandredis

import (
	"bufio"
	"io"
	"log"
	"net"
	"sync"
)

type Server struct {
	addr     string
	listener net.Listener
	proxy    *proxy

	mu    sync.Mutex
	conns []net.Conn
}

func NewServer(addr, hosts, keyspace string) (*Server, error) {
	proxy, err := newProxy(hosts, keyspace)
	if err != nil {
		return nil, err
	}

	return &Server{
		addr:  addr,
		proxy: proxy,
	}, nil
}

func (s *Server) Run() (err error) {
	s.listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("error while accepting connection: %v", err)
			continue
		}

		proxySession := makeProxySession(conn)

		s.mu.Lock()

		s.conns = append(s.conns, conn)
		s.proxy.addSession(proxySession)
		go handleConnection(conn, proxySession)

		s.mu.Unlock()
	}

	return
}

func handleConnection(conn net.Conn, proxySession *proxySession) {
	br := bufio.NewReader(conn)
	protocol := protocol{br}

	for {
		respValue, err := protocol.read()
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("could not parse: %v", err)
			continue
		}

		if respValue.respType != respArray {
			log.Println("should have received a array value")
			continue
		}

		array := respValue.value.(*respArrayValue)
		cmd, err := newCommandFromRespArray(array)
		if err != nil {
			log.Printf("could not parse command: %v", err)
			continue
		}

		proxySession.req <- cmd

		resp := <-proxySession.resp

		if err := resp.value.(respSerializable).SerializeTo(conn); err != nil {
			log.Printf("could not write: %v", err)
		}
	}
}
