//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package main

import (
	"l3/ospfv2/api"
	"l3/ospfv2/rpc"
	"l3/ospfv2/server"
	"strconv"
	"utils/dmnBase"
)

const (
	DMN_NAME = "ospfv2d"
)

type ospfv2Daemon struct {
	*dmnBase.FSBaseDmn
	server    *server.OSPFV2Server
	rpcServer *rpc.RPCServer
}

var dmn ospfv2Daemon

func main() {
	// Get base daemon handle and initialize
	dmn.FSBaseDmn = dmnBase.NewBaseDmn(DMN_NAME, DMN_NAME)
	ok := dmn.Init()
	if ok == false {
		panic("OSPF v2 Daemon: Base Daemon Initialization failed")
	}

	initParams := server.InitParams{
		Logger:    dmn.FSBaseDmn.Logger,
		DbHdl:     dmn.DbHdl,
		ParamsDir: dmn.ParamsDir,
		DmnName:   DMN_NAME,
	}

	// Get server handle and start server
	var err error
	dmn.server, err = server.NewOspfv2Server(initParams)
	if err != nil {
		panic("Unable to initilize ospfv2 Daemon")
	}
	go dmn.server.StartOspfv2Server()

	//Initialize API layer
	api.InitApiLayer(dmn.server)

	// Start Keep Alive for watchdog
	dmn.StartKeepAlive()

	_ = <-dmn.server.InitCompleteCh

	//Get RPC server handle
	var rpcServerAddr string
	for _, value := range dmn.FSBaseDmn.ClientsList {
		if value.Name == "ospfv2d" {
			rpcServerAddr = "localhost:" + strconv.Itoa(value.Port)
			break
		}
	}

	if rpcServerAddr == "" {
		panic("Ospf v2 Daemon is not part of system profile")
	}

	dmn.rpcServer = rpc.NewRPCServer(rpcServerAddr, dmn.FSBaseDmn.Logger, dmn.DbHdl)

	//Start RPC server
	dmn.FSBaseDmn.Logger.Info("Ospf V2 Daemon server started")
	dmn.rpcServer.Serve()
	panic("Ospf V2 Daemon RPC server terminated")
}
