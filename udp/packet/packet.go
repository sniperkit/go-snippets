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

package packet

import (
	"net"
)

// Packet represents a UDP packet, including remote addr.
type Packet struct {
	Payload []byte
	Size    int
	Remote  *net.UDPAddr
}

// NewPacket generates a Packet with a len([]byte) == cap.
func NewPacket(cap int) *Packet {
	return &Packet{
		Payload: make([]byte, cap),
		Size:    0,
		Remote:  nil,
	}
}

// Reset the underlying Buffer.
func (p *Packet) Reset() {
	p.Size = 0
	p.Remote = nil
}

// Read from UDP connection.
func (p *Packet) Read(conn *net.UDPConn) error {
	size, remote, err := conn.ReadFromUDP(p.Payload)

	if err != nil {
		return err
	}

	p.Size = size
	p.Remote = remote

	return nil
}

// Write to UDP connection.
func (p *Packet) Write(conn *net.UDPConn) error {
	_, err := conn.WriteToUDP(p.Payload[:p.Size], p.Remote)

	return err
}
