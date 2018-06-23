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

package server

import (
        "utils/dbutils"
        "utils/keepalive"
        "utils/logging"
)

type DmnServer struct {
        // store info related to server
        DbHdl          dbutils.DBIntf
        Logger         logging.LoggerIntf
        InitCompleteCh chan bool
}

type ServerInitParams struct {
        DmnName     string
        ParamsDir   string
        CfgFileName string
        DbHdl       dbutils.DBIntf
        Logger      logging.LoggerIntf
}

func NewISISDServer(initParams *ServerInitParams) *DmnServer {
        srvr := DmnServer{}
        srvr.DbHdl = initParams.DbHdl
        srvr.Logger = initParams.Logger
        srvr.InitCompleteCh = make(chan bool)

        // setup whatever you need for your server

        return &srvr
}

func (srvr *DmnServer) initServer() error {
        // initize the daemon server here
        return nil
}

func (srvr *DmnServer) Serve() {
        srvr.Logger.Info("Server initialization started")
        err := srvr.initServer()
        if err != nil {
              panic(err)
        }
        daemonStatusListener := keepalive.InitDaemonStatusListener()
        if daemonStatusListener != nil {
                go daemonStatusListener.StartDaemonStatusListner()
        }
        srvr.InitCompleteCh <- true
        srvr.Logger.Info("Server initialization complete, starting cfg/state listerner")
        for {
                select {
                //case req := <-srvr.ReqChan:
                //      srvr.Logger.Info("Server request received - ", *req)
                //      srvr.handleRPCRequest(req)
                case daemonStatus := <-daemonStatusListener.DaemonStatusCh:
                        srvr.Logger.Info("Received daemon status: ", daemonStatus.Name, daemonStatus.Status)
                }
        }
}
