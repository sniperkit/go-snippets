/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2017/08/12        Feng Yifei
 */

package server

import (
	"net"

	"github.com/fengyfei/nuts/udp/packet"
)

const (
	readBufferDefaultSize  = 65536
	writeBufferDefaultSize = 65536
)

// Server is a generic UDP server.
type Server struct {
	conf     *Conf
	conn     *net.UDPConn
	handler  Handler
	buffer   []*packet.Packet
	sender   chan *packet.Packet
	shutdown chan struct{}
}

// NewServer creates a UDP server instance.
func NewServer(conf *Conf, handler Handler) (*Server, error) {
	var (
		conn *net.UDPConn
	)

	addr, err := net.ResolveUDPAddr("udp", conf.Address+":"+conf.Port)
	if err != nil {
		return nil, err
	}

	conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	server := &Server{
		conf:     conf,
		conn:     conn,
		handler:  handler,
		buffer:   make([]*packet.Packet, conf.CacheCount),
		sender:   make(chan *packet.Packet, 64),
		shutdown: make(chan struct{}),
	}
	server.prepare()

	go server.start()

	return server, nil
}

func (server *Server) prepare() {
	server.conn.SetReadBuffer(readBufferDefaultSize)
	server.conn.SetWriteBuffer(writeBufferDefaultSize)

	for i := 0; i < server.conf.CacheCount; i++ {
		server.buffer[i] = packet.NewPacket(server.conf.PacketSize)
	}
}

// Shutdown the UDP server.
func (server *Server) Shutdown() {
	server.shutdown <- struct{}{}
}

// Start the messaging loop.
func (server *Server) start() {
	go server.receive()

	for {
		select {
		case <-server.shutdown:
			server.conn.Close()
			return

		case packet := <-server.sender:
			err := packet.Write(server.conn)

			if err != nil {
				server.handler.OnError(err)
			}
		}
	}
}

func (server *Server) receive() {
	index := 0
	cap := server.conf.CacheCount

	for {
		// Pick a packet and reset it for receiving.
		packet := server.buffer[index]
		packet.Reset()

		err := packet.Read(server.conn)

		if err != nil {
			server.handler.OnError(err)
		} else {
			server.handler.OnPacket(packet)
		}

		index = (index + 1) % cap
	}
}

// Send a message.
func (server *Server) Send(payload []byte, remote *net.UDPAddr) error {
	packet := &packet.Packet{
		Payload: payload,
		Size:    len(payload),
		Remote:  remote,
	}

	server.sender <- packet
	return nil
}
