package memcache

import (
	"bufio"
	"net"
	"strings"
	"time"
)

type Connection struct {
	conn     net.Conn
	buffered bufio.ReadWriter
	timeout  time.Duration
}

type Result struct {
	Key   string
	Value []byte
	Flags uint16
	Cas   uint64
}

func Connect(address string, timeout time.Duration) (conn *Connection, err error) {
	var network string
	if strings.Contains(address, "/") {
		network = "unix"
	} else {
		network = "tcp"
	}
	var nc net.Conn
	nc, err = net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, err
	}
	// code omitted ...
}
