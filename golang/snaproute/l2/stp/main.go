//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

// main
package main

import (
	"flag"
	"l2/stp/asicdMgr"
	stp "l2/stp/protocol"
	"l2/stp/rpc"
	"l2/stp/server"
	"utils/asicdClient"
	"utils/commonDefs"
	"utils/keepalive"
)

func main() {

	var err error

	// lookup port
	paramsDir := flag.String("params", "./params", "Params directory")
	flag.Parse()
	path := *paramsDir
	if path[len(path)-1] != '/' {
		path = path + "/"
	}
	clientInfoFile := path + "clients.json"

	stpServer := server.NewSTPServer(stp.GetStpLogger())

	// lets setup north bound notifications
	nHdl, nMap := asicdMgr.NewNotificationHdl(stpServer)
	asicdHdl := commonDefs.AsicdClientStruct{
		Logger: stp.GetStpLogger(),
		NHdl:   nHdl,
		NMap:   nMap,
	}
	asicdPlugin := asicdClient.NewAsicdClientInit("Flexswitch", clientInfoFile, asicdHdl)

	// connect to any needed services
	// This must be called before StartSTPSConfigNotificationListener
	stp.SetAsicDPlugin(asicdPlugin)
	stp.SaveSwitchMac(asicdPlugin.GetSwitchMAC(path))

	// Start keepalive routine
	go keepalive.InitKeepAlive("stpd", path)

	// this must be called before creating service handler as
	// the service handle instanciation will start the read config from db
	stpServer.InitServer()
	confIface := rpc.NewSTPDServiceHandler(stpServer)
	stp.StpLogger("INFO", "Starting STP Thrift daemon")
	rpc.StartServer(stp.GetStpLogger(), confIface, *paramsDir)
	stp.StpLogger("ERROR", "ERROR server not started")
	panic(err)
}
