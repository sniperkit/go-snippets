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
	"asicd/rpc"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"utils/keepalive"
	"utils/logging"
)

var logger *logging.Writer
var asicdServer *rpc.AsicDaemonServerInfo

func sigHandler() {
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGUSR1)
	signal := <-sigChan
	switch signal {
	case syscall.SIGHUP,
		syscall.SIGUSR1:
		//Stop thrift server
		asicdServer.Server.Stop()
		os.Exit(0)
	default:
		logger.Err(fmt.Sprintln("Unhandled signal : ", signal))
	}
}

type ClientJson struct {
	Name string `json:Name`
	Port int    `json:Port`
}

type initCfgParams struct {
	thriftServerPort int
}

func parseConfigFile(paramsDir string) *initCfgParams {
	var initCfg initCfgParams
	var clientsList []ClientJson

	//Retrieve thrift port number
	initCfg.thriftServerPort = 4000
	bytes, err := ioutil.ReadFile(paramsDir + "clients.json")
	if err != nil {
		logger.Err("Error retrieving thrift server port number using default port 4000")
	} else {
		err := json.Unmarshal(bytes, &clientsList)
		if err != nil {
			logger.Err("Error retrieving thrift server port number using default port 4000")
		} else {
			for _, client := range clientsList {
				if client.Name == "asicd" {
					initCfg.thriftServerPort = client.Port
				}
			}
		}
	}
	return &initCfg
}

func main() {
	var err error
	fmt.Println("Starting asicd daemon")
	paramsDirStr := flag.String("params", "", "Directory Location for config file")
	flag.Parse()
	paramsDir := *paramsDirStr
	if paramsDir[len(paramsDir)-1] != '/' {
		paramsDir = paramsDir + "/"
	}
	//Initialize logger
	logger, err = logging.NewLogger("asicd", "ASICD :", true)
	if err != nil {
		fmt.Println("Failed to start the logger. Nothing will be logged...")
	}
	//Parse cfg file and instantiate/initialize the plugin manager
	cfgFileInfo := parseConfigFile(paramsDir)
	// Start keepalive routine
	go keepalive.InitKeepAlive("asicd", paramsDir)
	//Start rpc server
	asicdServer = rpc.NewAsicdServer("localhost:" + strconv.Itoa(cfgFileInfo.thriftServerPort))
	go sigHandler()
	logger.Info("ASICD: server started")
	asicdServer.Server.Serve()
}
