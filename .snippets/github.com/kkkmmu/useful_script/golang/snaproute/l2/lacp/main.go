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
	"l2/lacp/asicdMgr"
	"l2/lacp/protocol/utils"
	"l2/lacp/rpc"
	"l2/lacp/server"
	"utils/asicdClient"
	"utils/commonDefs"
	"utils/keepalive"
	"utils/logging"
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

	logger, _ := logging.NewLogger("lacpd", "LA", true)
	utils.SetLaLogger(logger)
	laServer := server.NewLAServer(logger)

	// lets setup north bound notifications
	nHdl, nMap := asicdMgr.NewNotificationHdl(laServer)
	asicdHdl := commonDefs.AsicdClientStruct{
		Logger: logger,
		NHdl:   nHdl,
		NMap:   nMap,
	}
	asicdPlugin := asicdClient.NewAsicdClientInit("Flexswitch", clientInfoFile, asicdHdl)

	utils.SetAsicDPlugin(asicdPlugin)
	utils.SaveSwitchMac(asicdPlugin.GetSwitchMAC(path))

	// Start keepalive routine
	go keepalive.InitKeepAlive("lacpd", path)

	laServer.InitServer()
	confIface := rpc.NewLACPDServiceHandler(laServer)
	logger.Info("Starting LACP Thrift daemon")
	rpc.StartServer(utils.GetLaLogger(), confIface, *paramsDir)
	logger.Err("ERROR server not started")
	panic(err)
}
