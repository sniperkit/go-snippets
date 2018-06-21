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

package main

import (
	"flag"
	"fmt"
	"l2/lldp/api"
	"l2/lldp/flexswitch"
	"l2/lldp/server"
	"l2/lldp/utils"
	"utils/dbutils"
	"utils/keepalive"
	"utils/logging"
)

func main() {
	fmt.Println("Starting lldp daemon")
	paramsDir := flag.String("params", "./params", "Params directory")
	flag.Parse()
	fileName := *paramsDir
	if fileName[len(fileName)-1] != '/' {
		fileName = fileName + "/"
	}

	fmt.Println("Start logger")
	logger, err := logging.NewLogger("lldpd", "LLDP", true)
	if err != nil {
		fmt.Println("Failed to start the logger. Nothing will be logged...")
	}
	debug.SetLogger(logger)
	debug.Logger.Info("Started the logger successfully.")

	debug.Logger.Info("Starting LLDP server....")
	name := ""
	switch name {
	case "ovsdb":

	default:
		lldpDbHdl := dbutils.NewDBUtil(debug.Logger)
		aPlugin, err := flexswitch.NewAsicPlugin(fileName)
		if err != nil {
			return
		}
		sPlugin, err := flexswitch.NewSystemPlugin(fileName, lldpDbHdl)
		if err != nil {
			return
		}
		// Create lldp rpc handler
		lldpHdl := flexswitch.NewConfigHandler()
		lPlugin := flexswitch.NewNBPlugin(lldpHdl, fileName)

		// Create lldp server handler
		lldpSvr := server.LLDPNewServer(aPlugin, lPlugin, sPlugin, lldpDbHdl)
		// Start Api Layer
		api.Init(lldpSvr)

		// Until Server is connected to clients do not start with RPC
		lldpSvr.LLDPStartServer(*paramsDir)

		// Start keepalive routine
		go keepalive.InitKeepAlive("lldpd", fileName)
		debug.Logger.Info("Starting LLDP RPC listener....")
		err = lldpSvr.CfgPlugin.Start()
		if err != nil {
			debug.Logger.Err(fmt.Sprintln("Cannot start lldp server", err))
			return
		}
	}
}
