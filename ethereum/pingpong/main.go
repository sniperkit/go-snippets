package main

import (
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/smallnest/log"
)

type MsgType string

type Message struct {
	Data    string    `json:"data,omitempty"`
	Created time.Time `json:"created,omitempty"`
}

func main() {
	// start server 1
	key1, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("failed to generate private key: %v", err)
	}
	srv1 := NewServer(key1, "p2p_server_1", ":9981")
	err = srv1.Start()
	if err != nil {
		log.Fatalf("failed to start server1: %v", err)
	}
	defer srv1.Stop()

	// event1C := make(chan *p2p.PeerEvent)
	// sub1 := srv1.srv.SubscribeEvents(event1C)
	// go func() {
	// 	for {
	// 		select {
	// 		case e := <-event1C:
	// 			if e.Type == "add" {
	// 				log.Infof("#1 received peer add event: %v", e.Peer)
	// 			} else if e.Type == "msgrecv" {
	// 				log.Infof("#1 received message event: %v", e)
	// 			}
	// 		case <-sub1.Err():
	// 			return
	// 		}
	// 	}
	// }()

	node1Self := srv1.srv.Self()
	log.Infof("node1: IP=%s, TCP Port:%d, UDP Port:%d, ID: %s", node1Self.IP, node1Self.TCP, node1Self.UDP, node1Self.ID)

	// start server 2
	key2, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("failed to generate private key: %v", err)
	}
	srv2 := NewServer(key2, "p2p_server_2", ":9982")
	err = srv2.Start()
	if err != nil {
		log.Fatalf("failed to start server2: %v", err)
	}
	defer srv2.Stop()

	// event2C := make(chan *p2p.PeerEvent)
	// sub2 := srv2.srv.SubscribeEvents(event2C)
	// go func() {
	// 	for {
	// 		select {
	// 		case e := <-event2C:
	// 			if e.Type == "add" {
	// 				log.Infof("#2 received peer add event: %+v", e.Peer)
	// 			} else if e.Type == "msgrecv" {
	// 				log.Infof("#2 received message event: %+v", e)
	// 			}
	// 		case <-sub2.Err():
	// 			return
	// 		}
	// 	}
	// }()

	node2Self := srv2.srv.Self()
	log.Infof("node2: IP=%s, TCP Port:%d, UDP Port:%d, ID: %s", node2Self.IP, node2Self.TCP, node2Self.UDP, node2Self.ID)

	// add peers
	node2 := &discover.Node{
		IP:  node2Self.IP,
		TCP: node2Self.TCP,
		UDP: node2Self.UDP,
		ID:  node2Self.ID,
	}
	srv1.srv.AddPeer(node2)

	time.Sleep(5 * time.Second)
	log.Infof("server1 peers: %v, server2 peers: %v", srv1.srv.Peers(), srv2.srv.Peers())

	select {}
}
