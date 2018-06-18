// Package ch provides a handy wrapper mechanism for using a net.Conn via
// channels. This is currently only for text based streams where input is
// terminated via \r\n.  See echoServer/main.go for an example of using this
// package
package ch

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/textproto"
)

// Conn wraps a net.Conn and provides convenient in, out, and error channels.
type Conn struct {
	c   net.Conn
	In  chan string // Client input strings come down this channel
	Out chan string // Strings sent to this channel are relayed to the client
	Err chan string // This channel is closed when the connection is closed.
}

func (c *Conn) handleInput() {
	var r = textproto.NewReader(bufio.NewReader(c.c))
	for {
		input, err := r.ReadLine()
		if err != nil {
			if err != io.EOF {
				log.Printf("Got unexpected error reading from connection: %s", err.Error())
			}
			close(c.Err)
			return
		}
		c.In <- input
	}
}

func (c *Conn) handleOutput() {
	var w = textproto.NewWriter(bufio.NewWriter(c.c))
	for {
		select {
		case out := <-c.Out:
			w.PrintfLine(out)
		case <-c.Err:
			return
		}
	}
}

func (c *Conn) init() {
	go c.handleInput()
	go c.handleOutput()
}

// New returns a pointer to a new ch.Conn structure initialized and working.
func New(c net.Conn) *Conn {
	rval := &Conn{
		c:   c,
		In:  make(chan string),
		Out: make(chan string),
		Err: make(chan string),
	}
	go rval.init()
	return rval
}
