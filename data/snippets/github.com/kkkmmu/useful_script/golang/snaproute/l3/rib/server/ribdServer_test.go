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

package server

import (
	"fmt"
	"infra/sysd/sysdCommonDefs"
	"log/syslog"
	"utils/clntUtils/clntIntfs"
	"utils/clntUtils/clntIntfs/arpdClntIntfs"
	"utils/clntUtils/clntIntfs/asicdClntIntfs"
	"utils/dbutils"
	"utils/logging"
)

var server *RIBDServer

func RIBdNewLogger(name string, tag string) (*logging.Writer, error) {
	var err error
	srLogger := new(logging.Writer)
	srLogger.MyComponentName = name

	srLogger.SysLogger, err = syslog.New(syslog.LOG_INFO|syslog.LOG_DAEMON, tag)
	if err != nil {
		fmt.Println("Failed to initialize syslog - ", err)
		return srLogger, err
	}

	srLogger.MyLogLevel = sysdCommonDefs.DEBUG
	return srLogger, err
}

func getServerObject() *RIBDServer {
	logger, err := RIBdNewLogger("ribd", "RIBDTEST")
	if err != nil {
		fmt.Println("ribdtest: creating logger failed")
	}
	dbHdl := dbutils.NewDBUtil(logger)
	err = dbHdl.Connect()
	if err != nil {
		logger.Err("Failed to dial out to Redis server")
		return nil
	}
	testserver := NewRIBDServicesHandler(dbHdl, logger)
	if testserver == nil {
		fmt.Sprintln("ribd server object is null ")
		return nil
	}
	fileName := "/opt/flexswitch/params/"
	arpdClntInitParams, err := clntIntfs.NewBaseClntInitParams("arpd", logger, nil, fileName)
	if err != nil {
		logger.Err("RIBD: Error Initializing base clnt for arpd")
		panic(err)
	}
	testserver.ArpdClntPlugin, err = arpdClntIntfs.NewArpdClntInit(arpdClntInitParams)
	if err != nil {
		logger.Err("RIBD: Error Initializing new Arpd clnt")
		panic(err)
	}

	asicdClntInitParams, err := clntIntfs.NewBaseClntInitParams("asicd", logger, nil, fileName)
	if err != nil {
		logger.Err("RIBD: Error Initializing base clnt for asicd")
		panic(err)
	}
	testserver.AsicdPlugin, err = asicdClntIntfs.NewAsicdClntInit(asicdClntInitParams)
	if err != nil {
		logger.Err("RIBD: Error Initializing new Asicd Clnt")
		panic(err)
	}

	return testserver
}
func InitTestServer() {
	fmt.Println("Init server ")
	rtserver := getServerObject()
	if rtserver == nil {
		logger.Println("server nil")
		return
	}
	server = rtserver
}
func StartTestServer() {
	if server != nil {
		fmt.Println("server already initialized, return")
		return
	}
	InitTestServer()
	server.AcceptConfig = true
	go server.StartServer("/opt/flexswitch/params")
	InitLogicalIntfList()
	InitVlanList()
	InitIPv4IntfList()
	InitIPv6IntfList()
	InitIpv4AddrInfoList()
	InitIpv4RouteList()
	InitIpv6AddrInfoList()
	InitIpv6RouteList()
	InitPatchOpList()

	fmt.Println("server started")
}
