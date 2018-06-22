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
 *     Initial: 2017/08/19        Feng Yifei
 */

package tcp

import (
	"errors"
	"net"
	"time"

	"github.com/fengyfei/nuts/scheduler"
	"github.com/mailru/easygo/netpoll"
)

var (
	// ErrEpollCreate means couldn't create a epoll struct.
	ErrEpollCreate = errors.New("couldn't create epoll")
)

// Server represents a generic TCP server.
type Server struct {
	ln        net.Listener
	poller    netpoll.Poller
	scheduler *scheduler.Pool
	handler   Handler
}

// StartServer starts a TCP server based on configuration.
func StartServer(c *Config, h Handler, pool *scheduler.Pool) (*Server, error) {
	var (
		p netpoll.Poller
	)
	ln, err := net.Listen("tcp", c.Address)
	if err != nil {
		return nil, err
	}

	p, err = netpoll.New(nil)
	if err != nil {
		ln.Close()
		return nil, ErrEpollCreate
	}

	s := &Server{
		ln:        ln,
		poller:    p,
		scheduler: pool,
		handler:   h,
	}

	if err = s.accept(); err != nil {
		ln.Close()
		return nil, err
	}

	return s, nil
}

func (s *Server) accept() error {
	// error channel
	acceptError := make(chan error, 1)

	// read event handler
	handler := func(conn net.Conn) {
		desc := netpoll.Must(netpoll.HandleRead(conn))

		s.poller.Start(desc, func(e netpoll.Event) {
			// Client connection closed.
			if e&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
				s.handler.OnClose(conn)
				return
			}

			// Read from connection
			s.scheduler.Schedule(scheduler.TaskFunc(func() error {
				if err := s.handler.OnReadMessage(conn); err != nil {
					s.poller.Stop(desc)
					s.handler.OnError(conn)
					return err
				}
				return nil
			}))
		})
	}

	// accept event handler
	acceptor := func() error {
		conn, err := s.ln.Accept()

		if err != nil {
			acceptError <- err
		}

		handler(conn)

		return nil
	}

	acceptDesc := netpoll.Must(netpoll.HandleListener(s.ln, netpoll.EventRead|netpoll.EventOneShot))
	s.poller.Start(acceptDesc, func(e netpoll.Event) {
		err := s.scheduler.ScheduleWithTimeout(2*time.Microsecond, scheduler.TaskFunc(acceptor))
		if err == nil {
			err = <-acceptError
		}

		if err != nil {
			time.Sleep(5 * time.Microsecond)
		}

		s.poller.Resume(acceptDesc)
	})

	return nil
}
