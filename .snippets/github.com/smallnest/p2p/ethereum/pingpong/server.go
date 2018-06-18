package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/smallnest/log"
)

type Server struct {
	name        string
	srv         *p2p.Server
	sendMsgChan chan *Message
	recvMsgChan chan *Message

	done chan struct{}
}

func NewServer(key *ecdsa.PrivateKey, name string, addr string) *Server {
	sendMsgChan := make(chan *Message)
	recvMsgChan := make(chan *Message)
	done := make(chan struct{})

	var proto = p2p.Protocol{
		Name:    "hello",
		Version: 1,
		Length:  1,
		Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
			go func() {
				for {
					select {
					case <-done:
						return
					case msg := <-sendMsgChan:
						data, _ := json.Marshal(msg)

						// send the message
						err := p2p.Send(rw, 0, data)
						if err != nil {
							log.Infof("failed to send: %v", err)
							continue
						}
						log.Infof("sent msg to %+v", p)
					}
				}
			}()

			for {
				pmsg, err := rw.ReadMsg()
				if err != nil {
					return err
				}

				var data []byte
				err = pmsg.Decode(&data)
				if err != nil {
					log.Infof("failed to decode msg: %v", err)
					return err
				}

				var msg Message
				json.Unmarshal(data, &msg)
				log.Infof("received message from %v: %+v", p, msg)

				select {
				case <-done:
					return nil
				case recvMsgChan <- &msg:
				}

			}
		},
	}

	cfg := p2p.Config{
		MaxPeers:   10,
		PrivateKey: key,
		Name:       name,
		ListenAddr: addr,
		Protocols:  []p2p.Protocol{proto},
	}
	srv := &p2p.Server{
		Config: cfg,
	}

	return &Server{
		name:        name,
		srv:         srv,
		sendMsgChan: sendMsgChan,
		recvMsgChan: recvMsgChan,
		done:        done,
	}
}

func (s *Server) Start() error {
	go s.sendMsg()
	go s.processMsg()
	return s.srv.Start()
}

func (s *Server) sendMsg() {
	for {
		select {
		case <-s.done:
			return
		case <-time.After(10 * time.Second):
			pingMsg := &Message{
				Data:    time.Now().String(),
				Created: time.Now(),
			}

			s.sendMsgChan <- pingMsg
		}
	}
}

func (s *Server) processMsg() {
	for {
		select {
		case <-s.done:
			return
		case msg := <-s.recvMsgChan:
			_ = msg
			// log.Infof("received %+v", msg)
		}
	}
}

func (s *Server) Stop() {
	s.srv.Stop()
	close(s.done)
}
