package cassandredis

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/gocql/gocql"
)

type proxySession struct {
	client  net.Conn
	req     chan *Command
	resp    chan *respValue
	lastErr error
	// currentQueries Queries
	currentKey queryKey
	scratch    interface{}
}

type proxy struct {
	hosts    []string
	session  *gocql.Session
	sessions []*proxySession
	md       *metadata

	mu       sync.RWMutex
	ddlState map[string]bool
}

func newProxy(hosts, keyspace string) (*proxy, error) {
	realHosts := strings.Split(hosts, ",")
	if len(hosts) < 1 {
		return nil, errors.New("no cassandra hosts defined")
	}

	cluster := gocql.NewCluster(realHosts...)
	cluster.ProtoVersion = 2
	cluster.Keyspace = keyspace

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &proxy{
		hosts:    realHosts,
		session:  session,
		md:       newMetadata(),
		ddlState: make(map[string]bool),
	}, nil
}

func (p *proxy) addSession(session *proxySession) {
	p.sessions = append(p.sessions, session)

	go p.listen(session)
}

func (p *proxy) listen(session *proxySession) {
	for {
		req := <-session.req

		p.pre(req, session)
		p.ddl(req, session)
		p.process(req, session)

		if session.lastErr != nil {
			log.Println(session.lastErr)
			sendErrorResponse(session, session.lastErr)
		}
	}
}

func (p *proxy) pre(req *Command, session *proxySession) {
	switch req.Type {
	case CommandLPUSH:
		p.preLPUSH(req, session)
	case CommandLRANGE:
		p.preLRANGE(req, session)
	default:
		session.lastErr = fmt.Errorf("command %s not supported", req.Name)
	}
}

func (p *proxy) ddl(req *Command, session *proxySession) {
	if session.lastErr != nil {
		return
	}

	var queries Queries
	var err error

	switch req.Type {
	case CommandLPUSH:
		queries, err = p.ddlMapLPUSH(req, session)
	}

	if err != nil {
		session.lastErr = err
		return
	}

	for _, v := range queries {
		p.mu.RLock()
		tmp, ok := p.ddlState[v.Stmt]
		run := !ok || !tmp
		p.mu.RUnlock()

		if run {
			log.Printf("executing DDL query: %s", v)
			if err := p.session.Query(v.Stmt, v.Args).Exec(); err != nil {
				session.lastErr = err
				return
			}

			p.mu.Lock()
			p.ddlState[v.Stmt] = true
			p.mu.Unlock()
		}
	}

	return
}

func (p *proxy) process(req *Command, session *proxySession) {
	if session.lastErr != nil {
		return
	}

	switch req.Type {
	case CommandLPUSH:
		p.processLPUSH(req, session)
	case CommandLRANGE:
		p.processLRANGE(req, session)
	}
}

func sendArrayResponse(session *proxySession, values [][]byte) {
	resp := &respArrayValue{
		length: len(values),
	}
	for _, v := range values {
		r := &respBulkStringValue{v}
		resp.values = append(resp.values, r)
	}

	session.resp <- &respValue{respArray, resp}
}

func sendBulkStringResponse(session *proxySession, value []byte) {
	resp := &respBulkStringValue{value: value}
	session.resp <- &respValue{respBulkString, resp}
}

func sendIntegerResponse(session *proxySession, val int64) {
	resp := &respIntegerValue{value: val}
	session.resp <- &respValue{respInteger, resp}
}

func sendErrorResponse(session *proxySession, err error, args ...interface{}) {
	resp := &respErrorValue{message: err.Error()}
	session.resp <- &respValue{respError, resp}
}

func makeProxySession(conn net.Conn) *proxySession {
	return &proxySession{
		client: conn,
		req:    make(chan *Command),
		resp:   make(chan *respValue),
	}
}
