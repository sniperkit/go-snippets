package chanbuf

import (
	"errors"
	"net"

	"core/net/msg"
)

var ESenderFull = errors.New("sender chan is full!")

//
type Msg struct {
	h []byte // header
	b []byte // body
}

//
type ChanSender struct {
	chMsg  chan *Msg
	chSend chan bool

	msgNum int
}

//
func NewChanSender(sz int) *ChanSender {
	s := &ChanSender{msgNum: sz}
	s.reset()

	return s
}

//
func (this *ChanSender) reset() {
	this.chMsg = make(chan *Msg, this.msgNum)
	this.chSend = make(chan bool, 1)
}

//
func (this *ChanSender) Send(conn net.Conn) error {
	for {
		select {
		case d := <-this.chMsg:
			// 消息大小
			sz := uint32(len(d.h) + len(d.b))
			b := msg.Uint32Bytes(sz)
			if _, e := conn.Write(b); e != nil {
				return e
			}

			// 消息
			if _, e := conn.Write(d.h); e != nil {
				return e
			}

			if len(d.b) > 0 {
				if _, e := conn.Write(d.b); e != nil {
					return e
				}
			}

		default:
			return nil
		}
	}

	return nil
}

//
func (this *ChanSender) Write(h, b []byte) error {
	msg := &Msg{h: h, b: b}
	select {
	case this.chMsg <- msg:

	default:
		return ESenderFull
	}

	select {
	case this.chSend <- true:
	default:
	}

	return nil
}

//
func (this *ChanSender) WatchSend() <-chan bool {
	return this.chSend
}

//
func (this *ChanSender) Clear() {
	this.reset()
}
