package main

import (
	"fmt"
	"net"
	"net/rpc"
)

// if cfg.RpcPort is specified listenRpc opens up a
// Remote Procedure call listener to communicate with
// other servers
func listenRpc() error {
	if cfg.RpcPort == "" {
		log.Infoln("no rpc port specified, rpc disabled")
		return nil
	}

	if err := rpc.Register(UsersRequests); err != nil {
		log.Infof("register RPC Users error: %s", err)
		return err
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.RpcPort))
	if err != nil {
		log.Infof("listen on port %s error: %s", cfg.RpcPort, err)
		return err
	}

	log.Infof("accepting RPC requests on port %s", cfg.RpcPort)
	rpc.Accept(ln)
	return nil
}
