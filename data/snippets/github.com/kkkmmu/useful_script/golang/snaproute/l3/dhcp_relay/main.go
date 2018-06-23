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
//   This is a auto-generated file, please do not edit!
// _______   __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----  \   \/    \/   /  |  |  ---|  |----    ,---- |  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |        |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |        `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package main

import (
	"l3/dhcp_relay/rpc"
	"l3/dhcp_relay/server"
	"strconv"
	"strings"
	"utils/dmnBase"
)

const (
	DMN_NAME = "dhcprelayd"
)

type Daemon struct {
	*dmnBase.FSBaseDmn
	daemonServer *server.DmnServer
	rpcServer    *rpc.RPCServer
}

var dmn Daemon

func main() {
	//	dmn.FSBaseDmn = dmnBase.NewBaseDmn(DMN_NAME, DMN_NAME)
	dmnB := dmnBase.NewBaseDmn(DMN_NAME, DMN_NAME)
	dmn.FSBaseDmn = dmnB

	ok := dmnB.Init()
	if !ok {
		panic("Daemon Base initialization failed for dhcprelayd")
	}

	serverInitParams := &server.ServerInitParams{
		DmnName:     DMN_NAME,
		CfgFileName: dmn.ParamsDir + "clients.json",
		DbHdl:       dmn.DbHdl,
		Logger:      dmn.FSBaseDmn.Logger,
	}
	dmn.daemonServer = server.NewDHCPRELAYDServer(serverInitParams)
	go dmn.daemonServer.Serve()

	var rpcServerAddr string
	for _, value := range dmn.FSBaseDmn.ClientsList {
		if value.Name == strings.ToLower(DMN_NAME) {
			rpcServerAddr = "localhost:" + strconv.Itoa(value.Port)
			break
		}
	}
	if rpcServerAddr == "" {
		panic("Daemon dhcprelayd is not part of the system profile")
	}
	dmn.rpcServer = rpc.NewRPCServer(rpcServerAddr, dmn.FSBaseDmn.Logger, dmn.daemonServer)

	dmn.StartKeepAlive()

	// Wait for server started msg before opening up RPC port to accept calls
	_ = <-dmn.daemonServer.InitCompleteCh

	//Start RPC server
	dmn.FSBaseDmn.Logger.Info("Daemon Server started for dhcprelayd")
	dmn.rpcServer.Serve()
	panic("Daemon RPC Server terminated dhcprelayd")
}
