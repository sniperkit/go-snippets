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
	"l3/bfd/asicdMgr"
	"l3/bfd/rpc"
	"l3/bfd/server"
	"utils/clntUtils/clntIntfs"
	"utils/clntUtils/clntIntfs/asicdClntIntfs"
	"utils/dbutils"
	"utils/keepalive"
	"utils/logging"
)

func main() {
	fmt.Println("Starting bfd daemon")
	paramsDir := flag.String("params", "./params", "Params directory")
	flag.Parse()
	fileName := *paramsDir
	if fileName[len(fileName)-1] != '/' {
		fileName = fileName + "/"
	}

	fmt.Println("Start logger")
	logger, err := logging.NewLogger("bfdd", "BFD", true)
	if err != nil {
		fmt.Println("Failed to start the logger. Nothing will be logged...")
	}
	logger.Info("Started the logger successfully.")

	dbHdl := dbutils.NewDBUtil(logger)
	err = dbHdl.Connect()
	if err != nil {
		logger.Err("Failed to dial out to Redis server")
		return
	}

	clientsFileName := fileName + "clients.json"

	logger.Info("Starting BFD Server...")
	bfdServer := server.NewBFDServer(logger)

	asicdNHdl := asicdMgr.NewNotificationHdl(bfdServer, logger)
	asicdClntInitParams, err := clntIntfs.NewBaseClntInitParams("asicd", logger, asicdNHdl, fileName)
	if err != nil {
		logger.Err("BFDD: Error Initializing base clnt for asicd")
		panic(err)
	}
	bfdServer.AsicdPlugin, err = asicdClntIntfs.NewAsicdClntInit(asicdClntInitParams)
	if err != nil {
		logger.Err("BFDD: Error Initializing new Asicd Clnt")
		panic(err)
	}

	// Start signal handler
	go bfdServer.SigHandler(dbHdl)
	// Start bfd server
	go bfdServer.StartServer(clientsFileName, dbHdl)

	<-bfdServer.ServerStartedCh
	logger.Info("BFD Server started")

	// Start keepalive routine
	go keepalive.InitKeepAlive("bfdd", fileName)

	logger.Info("Starting Config listener...")
	confIface := rpc.NewBFDHandler(logger, bfdServer)
	// Read BFD configurations already present in DB
	confIface.ReadConfigFromDB(dbHdl)
	rpc.StartServer(logger, confIface, clientsFileName)
}
