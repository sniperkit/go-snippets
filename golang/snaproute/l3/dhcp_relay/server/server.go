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
	"errors"
	"fmt"
	"l3/dhcp_relay/infra"
	"l3/dhcp_relay/manager"
	"utils/asicdClient"
	"utils/commonDefs"
	"utils/dbutils"
	"utils/eventUtils"
	"utils/keepalive"
	"utils/logging"
)

type DmnServer struct {
	// store info related to server
	DbHdl          dbutils.DBIntf
	Logger         logging.LoggerIntf
	CfgFileName    string
	InitCompleteCh chan bool

	AsicdHdl         asicdClient.AsicdClientIntf
	AsicdSubSocketCh chan commonDefs.AsicdNotifyMsg
	EventDbHdl       *dbutils.DBUtil

	// API message request reply channels
	ReqChan   chan *ServerRequest
	ReplyChan chan *ServerReply

	// DHCP Relay Manager and Infrastructure Manager
	DMgr *manager.DRAMgr
	IMgr *infra.InfraMgr
}

type ServerInitParams struct {
	DmnName     string
	CfgFileName string
	DbHdl       dbutils.DBIntf
	Logger      logging.LoggerIntf
	// AsicdHdl   asicdClient.AsicdClientIntf
}

func NewDHCPRELAYDServer(initParams *ServerInitParams) *DmnServer {
	srvr := DmnServer{}
	srvr.DbHdl = initParams.DbHdl
	srvr.Logger = initParams.Logger
	srvr.CfgFileName = initParams.CfgFileName
	srvr.InitCompleteCh = make(chan bool)
	//Parse dhcprelayd manager file
	//	cfgFileInfo, err := parseCfgFile(initParams.CfgFileName)
	//	if err != nil {
	//		srvr.logger.Err("Failed to parse dhcprelayd manager file, using default values for all attributes")
	//	}
	return &srvr
}

func (srvr *DmnServer) initServer() error {
	// Allocate reply request channels (buffer len 1?)
	srvr.ReqChan = make(chan *ServerRequest)
	srvr.ReplyChan = make(chan *ServerReply)

	srvr.AsicdHdl = srvr.initAsicdHandler()
	srvr.AsicdSubSocketCh = make(chan commonDefs.AsicdNotifyMsg)
	err := srvr.initEventHandler()
	if err != nil {
		srvr.Logger.Err(fmt.Sprintln(
			"Unable to initialize events", err,
		),
		)
		return err
	}

	srvr.IMgr = infra.NewInfraMgr(srvr.Logger, srvr.AsicdHdl)
	srvr.DMgr = manager.NewDRAMgr(srvr.Logger, srvr.DbHdl, srvr.IMgr)
	if !srvr.DMgr.InitDRAMgr() {
		return errors.New("Unable to initialize DHCP Relay Manager")
	}
	return nil
}

func (srvr *DmnServer) initAsicdHandler() asicdClient.AsicdClientIntf {
	nHdl, nMap := srvr.NewNotificationHdl(srvr, srvr.Logger.(*logging.Writer))
	tmpHdl := commonDefs.AsicdClientStruct{
		Logger: srvr.Logger.(*logging.Writer),
		NHdl:   nHdl,
		NMap:   nMap,
	}
	return asicdClient.NewAsicdClientInit(
		"Flexswitch", srvr.CfgFileName, tmpHdl)
}

func (srvr *DmnServer) initEventHandler() error {
	srvr.EventDbHdl = dbutils.NewDBUtil(srvr.Logger)
	err := srvr.EventDbHdl.Connect()
	if err != nil {
		srvr.Logger.Err("Failed to create event DB handle")
		return err
	}

	return eventUtils.InitEvents("dhcprelayd", srvr.EventDbHdl, srvr.EventDbHdl, srvr.Logger, 1000)
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
		case req := <-srvr.ReqChan:
			srvr.Logger.Info("Server request received - ", *req)
			switch req.Op {
			case CREATE_DHCPRELAY_GLOBAL:
				success, err := srvr.DMgr.CreateDRAv4Global(
					req.Data.(*CreateDHCPRelayGlobalInArgs).
						DHCPRelayGlobal,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj:   success,
					Error: err,
				}
			case UPDATE_DHCPRELAY_GLOBAL:
				success, err := srvr.DMgr.UpdateDRAv4Global(
					req.Data.(*UpdateDHCPRelayGlobalInArgs).
						DHCPRelayGlobalOld,
					req.Data.(*UpdateDHCPRelayGlobalInArgs).
						DHCPRelayGlobalNew,
					req.Data.(*UpdateDHCPRelayGlobalInArgs).
						AttrSet,
					nil,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj:   success,
					Error: err,
				}
			case DELETE_DHCPRELAY_GLOBAL:
				success, err := srvr.DMgr.DeleteDRAv4Global(
					req.Data.(*DeleteDHCPRelayGlobalInArgs).
						Vrf)
				srvr.ReplyChan <- &ServerReply{
					Obj:   success,
					Error: err,
				}
			case CREATE_DHCPRELAY_INTF:
				success, err := srvr.DMgr.CreateDRAv4Interface(
					req.Data.(*CreateDHCPRelayIntfInArgs).
						DHCPRelayIntf,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj:   success,
					Error: err,
				}
			case UPDATE_DHCPRELAY_INTF:
				success, err := srvr.DMgr.UpdateDRAv4Interface(
					req.Data.(*UpdateDHCPRelayIntfInArgs).
						DHCPRelayIntfOld,
					req.Data.(*UpdateDHCPRelayIntfInArgs).
						DHCPRelayIntfNew,
					req.Data.(*UpdateDHCPRelayIntfInArgs).
						AttrSet,
					nil,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj:   success,
					Error: err,
				}
			case DELETE_DHCPRELAY_INTF:
				success, err := srvr.DMgr.DeleteDRAv4Interface(
					req.Data.(*DeleteDHCPRelayIntfInArgs).
						IntfRef)
				srvr.ReplyChan <- &ServerReply{
					Obj:   success,
					Error: err,
				}
			case GET_DHCPRELAY_CLIENT_STATE:
				state, err := srvr.DMgr.GetDRAv4ClientState(
					req.Data.(*GetDHCPRelayClientStateInArgs).
						MacAddr,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj: &GetDHCPRelayClientStateOutArgs{
						Obj: state,
						Err: err,
					},
					Error: nil,
				}
			case GETBLK_DHCPRELAY_CLIENT_STATE:
				blkState := srvr.DMgr.GetBulkDRAv4ClientState(
					req.Data.(*GetBulkInArgs).FromIdx,
					req.Data.(*GetBulkInArgs).Count,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj: &GetBulkDHCPRelayClientStateOutArgs{
						Obj: blkState,
						Err: nil,
					},
					Error: nil,
				}
			case GET_DHCPRELAY_INTF_STATE:
				state, err := srvr.DMgr.GetDRAv4IntfState(
					req.Data.(*GetDHCPRelayIntfStateInArgs).
						IntfRef,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj: &GetDHCPRelayIntfStateOutArgs{
						Obj: state,
						Err: err,
					},
					Error: nil,
				}
			case GETBLK_DHCPRELAY_INTF_STATE:
				blkState := srvr.DMgr.GetBulkDRAv4IntfState(
					req.Data.(*GetBulkInArgs).FromIdx,
					req.Data.(*GetBulkInArgs).Count,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj: &GetBulkDHCPRelayIntfStateOutArgs{
						Obj: blkState,
						Err: nil,
					},
					Error: nil,
				}
			case GET_DHCPRELAY_INTFSERVER_STATE:
				state, err := srvr.DMgr.GetDRAv4IntfServerState(
					req.Data.(*GetDHCPRelayIntfServerStateInArgs).
						IntfRef,
					req.Data.(*GetDHCPRelayIntfServerStateInArgs).
						ServerIp,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj: &GetDHCPRelayIntfServerStateOutArgs{
						Obj: state,
						Err: err,
					},
					Error: nil,
				}
			case GETBLK_DHCPRELAY_INTFSERVER_STATE:
				blkState := srvr.DMgr.GetBulkDRAv4IntfServerState(
					req.Data.(*GetBulkInArgs).FromIdx,
					req.Data.(*GetBulkInArgs).Count,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj: &GetBulkDHCPRelayIntfServerStateOutArgs{
						Obj: blkState,
						Err: nil,
					},
					Error: nil,
				}
			case CREATE_DHCPV6RELAY_GLOBAL:
				success, err := srvr.DMgr.CreateDRAv6Global(
					req.Data.(*CreateDHCPv6RelayGlobalInArgs).
						DHCPv6RelayGlobal,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj:   success,
					Error: err,
				}
			case UPDATE_DHCPV6RELAY_GLOBAL:
				success, err := srvr.DMgr.UpdateDRAv6Global(
					req.Data.(*UpdateDHCPv6RelayGlobalInArgs).
						DHCPv6RelayGlobalOld,
					req.Data.(*UpdateDHCPv6RelayGlobalInArgs).
						DHCPv6RelayGlobalNew,
					req.Data.(*UpdateDHCPv6RelayGlobalInArgs).
						AttrSet,
					nil,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj:   success,
					Error: err,
				}
			case DELETE_DHCPV6RELAY_GLOBAL:
				success, err := srvr.DMgr.DeleteDRAv6Global(
					req.Data.(*DeleteDHCPv6RelayGlobalInArgs).
						Vrf)
				srvr.ReplyChan <- &ServerReply{
					Obj:   success,
					Error: err,
				}
			case CREATE_DHCPV6RELAY_INTF:
				success, err := srvr.DMgr.CreateDRAv6Interface(
					req.Data.(*CreateDHCPv6RelayIntfInArgs).
						DHCPv6RelayIntf,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj:   success,
					Error: err,
				}
			case UPDATE_DHCPV6RELAY_INTF:
				success, err := srvr.DMgr.UpdateDRAv6Interface(
					req.Data.(*UpdateDHCPv6RelayIntfInArgs).
						DHCPv6RelayIntfOld,
					req.Data.(*UpdateDHCPv6RelayIntfInArgs).
						DHCPv6RelayIntfNew,
					req.Data.(*UpdateDHCPv6RelayIntfInArgs).
						AttrSet,
					nil,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj:   success,
					Error: err,
				}
			case DELETE_DHCPV6RELAY_INTF:
				success, err := srvr.DMgr.DeleteDRAv6Interface(
					req.Data.(*DeleteDHCPv6RelayIntfInArgs).
						IntfRef)
				srvr.ReplyChan <- &ServerReply{
					Obj:   success,
					Error: err,
				}
			case GET_DHCPV6RELAY_CLIENT_STATE:
				state, err := srvr.DMgr.GetDRAv6ClientState(
					req.Data.(*GetDHCPv6RelayClientStateInArgs).
						MacAddr,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj: &GetDHCPv6RelayClientStateOutArgs{
						Obj: state,
						Err: err,
					},
					Error: nil,
				}
			case GETBLK_DHCPV6RELAY_CLIENT_STATE:
				blkState := srvr.DMgr.GetBulkDRAv6ClientState(
					req.Data.(*GetBulkInArgs).FromIdx,
					req.Data.(*GetBulkInArgs).Count,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj: &GetBulkDHCPv6RelayClientStateOutArgs{
						Obj: blkState,
						Err: nil,
					},
					Error: nil,
				}
			case GET_DHCPV6RELAY_INTF_STATE:
				state, err := srvr.DMgr.GetDRAv6IntfState(
					req.Data.(*GetDHCPv6RelayIntfStateInArgs).
						IntfRef,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj: &GetDHCPv6RelayIntfStateOutArgs{
						Obj: state,
						Err: err,
					},
					Error: nil,
				}
			case GETBLK_DHCPV6RELAY_INTF_STATE:
				blkState := srvr.DMgr.GetBulkDRAv6IntfState(
					req.Data.(*GetBulkInArgs).FromIdx,
					req.Data.(*GetBulkInArgs).Count,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj: &GetBulkDHCPv6RelayIntfStateOutArgs{
						Obj: blkState,
						Err: nil,
					},
					Error: nil,
				}
			case GET_DHCPV6RELAY_INTFSERVER_STATE:
				state, err := srvr.DMgr.GetDRAv6IntfServerState(
					req.Data.(*GetDHCPv6RelayIntfServerStateInArgs).
						IntfRef,
					req.Data.(*GetDHCPv6RelayIntfServerStateInArgs).
						ServerIp,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj: &GetDHCPv6RelayIntfServerStateOutArgs{
						Obj: state,
						Err: err,
					},
					Error: nil,
				}
			case GETBLK_DHCPV6RELAY_INTFSERVER_STATE:
				blkState := srvr.DMgr.GetBulkDRAv6IntfServerState(
					req.Data.(*GetBulkInArgs).FromIdx,
					req.Data.(*GetBulkInArgs).Count,
				)
				srvr.ReplyChan <- &ServerReply{
					Obj: &GetBulkDHCPv6RelayIntfServerStateOutArgs{
						Obj: blkState,
						Err: nil,
					},
					Error: nil,
				}
			}
		case msg := <-srvr.AsicdSubSocketCh:
			srvr.Logger.Info("Notification", msg)
			srvr.DMgr.ProcessAsicdNotification(msg)
		case daemonStatus := <-daemonStatusListener.DaemonStatusCh:
			srvr.Logger.Info("Received daemon status: ",
				daemonStatus.Name, daemonStatus.Status)
		}
	}
}
