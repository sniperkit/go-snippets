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
 *     Initial: 2017/08/13        Feng Yifei
 */

package main

import (
	"fmt"

	"github.com/fengyfei/nuts/udp/packet"
	"github.com/fengyfei/nuts/udp/server"
)

var (
	udpServer *server.Server
)

func main() {
	newServer()

	close := make(chan struct{})

	select {
	case <-close:
		return
	}
}

func newServer() {
	var err error

	conf := server.Conf{
		Address:    "",
		Port:       "9527",
		PacketSize: 32,
		CacheCount: 12,
	}

	udpServer, err = server.NewServer(&conf, &handler{})

	if err != nil {
		panic(err)
	}
}

type handler struct {
}

func (h *handler) OnPacket(p *packet.Packet) error {
	fmt.Println(string(p.Payload), " from ", p.Remote)
	resp := make([]byte, p.Size)
	copy(resp, p.Payload[:p.Size])

	udpServer.Send(resp, p.Remote)
	return nil
}

func (h *handler) OnError(err error) error {
	fmt.Println(err)
	return nil
}

func (h *handler) OnClose() error {
	fmt.Println("OnClose triggered")
	return nil
}
