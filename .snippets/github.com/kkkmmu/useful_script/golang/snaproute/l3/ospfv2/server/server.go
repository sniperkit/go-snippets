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
	"encoding/json"
	"errors"
	"git.apache.org/thrift.git/lib/go/thrift"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
	"utils/dbutils"
	"utils/logging"
)

type ClientJson struct {
	Name string `json:Name`
	Port int    `json:Port`
}

type OspfClientBase struct {
	Address            string
	Transport          thrift.TTransport
	PtrProtocolFactory *thrift.TBinaryProtocolFactory
	IsConnected        bool
}

type InitParams struct {
	Logger    logging.LoggerIntf
	DbHdl     dbutils.DBIntf
	ParamsDir string
	DmnName   string
}

type OSPFV2Server struct {
	paramsDir      string
	dmnName        string
	logger         logging.LoggerIntf
	dbHdl          dbutils.DBIntf
	ReqChan        chan *ServerRequest
	ReplyChan      chan interface{}
	InitCompleteCh chan bool

	ribdComm  RibdCommStruct
	asicdComm AsicdCommStruct

	infraData InfraStruct

	globalData      GlobalStruct
	IntfConfMap     map[IntfConfKey]IntfConf
	NbrConfMap      map[NbrConfKey]NbrConf
	AreaConfMap     map[uint32]AreaConf //Key AreaId
	MessagingChData MessagingChStruct

	NbrConfData    NbrStruct
	LsdbData       LsdbStruct
	FloodData      FloodStruct
	SPFData        SPFStruct
	RoutingTblData RoutingTblStruct
	SummaryLsDb    map[LsdbKey]SummaryLsaMap

	GetBulkData GetBulkStruct
}

func NewOspfv2Server(initParams InitParams) (*OSPFV2Server, error) {
	var server OSPFV2Server

	server.logger = initParams.Logger
	server.dbHdl = initParams.DbHdl
	server.dmnName = initParams.DmnName
	server.paramsDir = initParams.ParamsDir
	server.ReqChan = make(chan *ServerRequest)
	server.ReplyChan = make(chan interface{})
	server.InitCompleteCh = make(chan bool)
	server.IntfConfMap = make(map[IntfConfKey]IntfConf)
	server.AreaConfMap = make(map[uint32]AreaConf)
	return &server, nil
}

func (server *OSPFV2Server) ConnectToServers() {
	var clientsList []ClientJson

	paramsFile := server.paramsDir + "/clients.json"

	bytes, err := ioutil.ReadFile(paramsFile)
	if err != nil {
		server.logger.Info("Error in reading configuration file")
		return
	}

	err = json.Unmarshal(bytes, &clientsList)
	if err != nil {
		server.logger.Info("Error in Unmarshalling Json")
		return
	}
	for _, client := range clientsList {
		if client.Name == "ribd" {
			server.ConnectToRibdServer(client.Port)
		} else if client.Name == "asicd" {
			server.ConnectToAsicdServer(client.Port)
		}
	}
}

func (server *OSPFV2Server) StartSubscribers() {
	server.StartAsicdSubscriber()
	server.StartRibdSubscriber()
}

func (server *OSPFV2Server) initMessagingChData() {
	server.MessagingChData.IntfToNbrFSMChData.NbrHelloEventCh = make(chan NbrHelloEventMsg)
	server.MessagingChData.IntfToNbrFSMChData.DeleteNbrCh = make(chan DeleteNbrMsg)
	server.MessagingChData.IntfToNbrFSMChData.NetworkDRChangeCh = make(chan NetworkDRChangeMsg)
	server.MessagingChData.IntfFSMToLsdbChData.GenerateRouterLSACh = make(chan GenerateRouterLSAMsg)
	server.MessagingChData.NbrToIntfFSMChData.NbrDownMsgChMap = make(map[IntfConfKey]chan NbrDownMsg)
	server.MessagingChData.NbrFSMToLsdbChData.RecvdLsaMsgCh = make(chan RecvdLsaMsg)
	server.MessagingChData.NbrFSMToLsdbChData.RecvdSelfLsaMsgCh = make(chan RecvdSelfLsaMsg)
	server.MessagingChData.NbrFSMToLsdbChData.UpdateSelfNetworkLSACh = make(chan UpdateSelfNetworkLSAMsg)
	server.MessagingChData.NbrFSMToLsdbChData.NbrDeadMsgCh = make(chan NbrDeadMsg)
	server.MessagingChData.LsdbToFloodChData.LsdbToFloodLSACh = make(chan []LsdbToFloodLSAMsg)
	server.MessagingChData.NbrFSMToFloodChData.LsaFloodCh = make(chan NbrToFloodMsg)
	server.MessagingChData.LsdbToSPFChData.StartSPF = make(chan bool)
	server.MessagingChData.SPFToLsdbChData.DoneSPF = make(chan bool)
	server.MessagingChData.ServerToLsdbChData.RefreshLsdbSliceCh = make(chan bool)
	server.MessagingChData.ServerToLsdbChData.RouteInfoDataUpdateCh = make(chan RouteInfoDataUpdateMsg)
	server.MessagingChData.ServerToLsdbChData.InitAreaLsdbCh = make(chan uint32)
	server.MessagingChData.LsdbToServerChData.InitAreaLsdbDoneCh = make(chan bool)
	server.MessagingChData.LsdbToServerChData.RefreshLsdbSliceDoneCh = make(chan bool)
	server.MessagingChData.RouteTblToDBClntChData.RouteAddMsgCh = make(chan RouteAddMsg, 100)
	server.MessagingChData.RouteTblToDBClntChData.RouteDelMsgCh = make(chan RouteDelMsg, 100)
	server.MessagingChData.ServerToDBClntChData.FlushRouteFromDBCh = make(chan bool)
	server.MessagingChData.DBClntToServerChData.FlushRouteFromDBDoneCh = make(chan bool)
}

func (server *OSPFV2Server) SigHandler(sigChan <-chan os.Signal) {
	server.logger.Debug("Inside sigHandler....")
	signal := <-sigChan
	switch signal {
	case syscall.SIGHUP:
		server.logger.Debug("Received SIGHUP signal")
		server.SendFlushRouteMsgToDBClnt()
		<-server.MessagingChData.DBClntToServerChData.FlushRouteFromDBDoneCh
		debug.PrintStack()
		var memStat runtime.MemStats
		runtime.ReadMemStats(&memStat)
		server.logger.Info("===Memstat===", memStat)
	default:
		server.logger.Err("Unhandled signal : ", signal)
	}
}

func (server *OSPFV2Server) initServer() error {
	server.logger.Info("Starting OspfV2 server")
	sigChan := make(chan os.Signal, 1)
	signalList := []os.Signal{syscall.SIGHUP}
	signal.Notify(sigChan, signalList...)
	go server.SigHandler(sigChan)
	server.initMessagingChData()
	server.initAsicdComm()
	server.initRibdComm()
	server.ConnectToServers()
	server.StartSubscribers()
	server.initInfra()
	server.buildInfra()
	//TODO:server.DeleteRouteFromDB()
	server.InitGetBulkSliceRefresh()
	go server.GetBulkSliceRefresh()
	if server.dbHdl == nil {
		server.logger.Err("DB Handle is nil")
		return errors.New("DB Handle is nil")
	}
	go server.StartDBClient()
	return nil
}

func (server *OSPFV2Server) handleRPCRequest(req *ServerRequest) {
	server.logger.Info("Handle RPC Request:", *req)
	switch req.Op {
	case CREATE_OSPFV2_AREA:
		var retObj CreateConfigOutArgs
		if val, ok := req.Data.(*CreateOspfv2AreaInArgs); ok {
			retObj.RetVal, retObj.Err = server.createArea(val.Cfg)
		}
		server.ReplyChan <- interface{}(&retObj)
	case UPDATE_OSPFV2_AREA:
		var retObj UpdateConfigOutArgs
		if val, ok := req.Data.(*UpdateOspfv2AreaInArgs); ok {
			retObj.RetVal, retObj.Err = server.updateArea(val.NewCfg, val.OldCfg, val.AttrSet)
		}
		server.ReplyChan <- interface{}(&retObj)
	case DELETE_OSPFV2_AREA:
		var retObj DeleteConfigOutArgs
		if val, ok := req.Data.(*DeleteOspfv2AreaInArgs); ok {
			retObj.RetVal, retObj.Err = server.deleteArea(val.Cfg)
		}
		server.ReplyChan <- interface{}(&retObj)
	case GET_OSPFV2_AREA_STATE:
		var retObj GetOspfv2AreaStateOutArgs
		if val, ok := req.Data.(*GetOspfv2AreaStateInArgs); ok {
			retObj.Obj, retObj.Err = server.getAreaState(val.AreaId)
		}
		server.ReplyChan <- interface{}(&retObj)
	case GET_BULK_OSPFV2_AREA_STATE:
		var retObj GetBulkOspfv2AreaStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = server.getBulkAreaState(val.FromIdx, val.Count)
		}
		server.ReplyChan <- interface{}(&retObj)
	case CREATE_OSPFV2_GLOBAL:
		var retObj CreateConfigOutArgs
		if val, ok := req.Data.(*CreateOspfv2GlobalInArgs); ok {
			retObj.RetVal, retObj.Err = server.createGlobal(val.Cfg)
		}
		server.ReplyChan <- interface{}(&retObj)
	case UPDATE_OSPFV2_GLOBAL:
		var retObj UpdateConfigOutArgs
		if val, ok := req.Data.(*UpdateOspfv2GlobalInArgs); ok {
			retObj.RetVal, retObj.Err = server.updateGlobal(val.NewCfg, val.OldCfg, val.AttrSet)
		}
		server.ReplyChan <- interface{}(&retObj)
	case DELETE_OSPFV2_GLOBAL:
		var retObj DeleteConfigOutArgs
		if val, ok := req.Data.(*DeleteOspfv2GlobalInArgs); ok {
			retObj.RetVal, retObj.Err = server.deleteGlobal(val.Cfg)
		}
		server.ReplyChan <- interface{}(&retObj)
	case GET_OSPFV2_GLOBAL_STATE:
		var retObj GetOspfv2GlobalStateOutArgs
		if val, ok := req.Data.(*GetOspfv2GlobalStateInArgs); ok {
			retObj.Obj, retObj.Err = server.getGlobalState(val.Vrf)
		}
		server.ReplyChan <- interface{}(&retObj)
	case GET_BULK_OSPFV2_GLOBAL_STATE:
		var retObj GetBulkOspfv2GlobalStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = server.getBulkGlobalState(val.FromIdx, val.Count)
		}
		server.ReplyChan <- interface{}(&retObj)
	case CREATE_OSPFV2_INTF:
		var retObj CreateConfigOutArgs
		if val, ok := req.Data.(*CreateOspfv2IntfInArgs); ok {
			retObj.RetVal, retObj.Err = server.createIntf(val.Cfg)
		}
		server.ReplyChan <- interface{}(&retObj)
	case UPDATE_OSPFV2_INTF:
		var retObj UpdateConfigOutArgs
		if val, ok := req.Data.(*UpdateOspfv2IntfInArgs); ok {
			retObj.RetVal, retObj.Err = server.updateIntf(val.NewCfg, val.OldCfg, val.AttrSet)
		}
		server.ReplyChan <- interface{}(&retObj)
	case DELETE_OSPFV2_INTF:
		var retObj DeleteConfigOutArgs
		if val, ok := req.Data.(*DeleteOspfv2IntfInArgs); ok {
			retObj.RetVal, retObj.Err = server.deleteIntf(val.Cfg)
		}
		server.ReplyChan <- interface{}(&retObj)
	case GET_OSPFV2_INTF_STATE:
		var retObj GetOspfv2IntfStateOutArgs
		if val, ok := req.Data.(*GetOspfv2IntfStateInArgs); ok {
			retObj.Obj, retObj.Err = server.getIntfState(val.IpAddr, val.AddressLessIfIdx)
		}
		server.ReplyChan <- interface{}(&retObj)
	case GET_BULK_OSPFV2_INTF_STATE:
		var retObj GetBulkOspfv2IntfStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = server.getBulkIntfState(val.FromIdx, val.Count)
		}
		server.ReplyChan <- interface{}(&retObj)
	case GET_OSPFV2_NBR_STATE:
		var retObj GetOspfv2NbrStateOutArgs
		if val, ok := req.Data.(*GetOspfv2NbrStateInArgs); ok {
			retObj.Obj, retObj.Err = server.getNbrState(val.IpAddr, val.AddressLessIfIdx)
		}
		server.ReplyChan <- interface{}(&retObj)
	case GET_BULK_OSPFV2_NBR_STATE:
		var retObj GetBulkOspfv2NbrStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = server.getBulkNbrState(val.FromIdx, val.Count)
		}
		server.ReplyChan <- interface{}(&retObj)
	case GET_OSPFV2_LSDB_STATE:
		var retObj GetOspfv2LsdbStateOutArgs
		if val, ok := req.Data.(*GetOspfv2LsdbStateInArgs); ok {
			retObj.Obj, retObj.Err = server.getLsdbState(val.LSType, val.LSId, val.AreaId, val.AdvRtrId)
		}
		server.ReplyChan <- interface{}(&retObj)
	case GET_BULK_OSPFV2_LSDB_STATE:
		var retObj GetBulkOspfv2LsdbStateOutArgs
		if val, ok := req.Data.(*GetBulkInArgs); ok {
			retObj.BulkInfo, retObj.Err = server.getBulkLsdbState(val.FromIdx, val.Count)
		}
		server.ReplyChan <- interface{}(&retObj)
	default:
		server.logger.Err("Error: Server received unrecognized request -", req.Op)
	}
}

func (server *OSPFV2Server) StartOspfv2Server() {
	err := server.initServer()
	if err != nil {
		panic(err)
	}
	server.InitCompleteCh <- true

	for {
		select {
		case req := <-server.ReqChan:
			server.logger.Debug("Handling RPC Req", req)
			server.handleRPCRequest(req)
			server.logger.Debug("Done Handling RPC Req", req)
		case asicdRxBuf := <-server.asicdComm.asicdSubSocketCh:
			server.logger.Debug("Process Asicd Rx Buf", asicdRxBuf)
			server.processAsicdNotification(asicdRxBuf)
			server.logger.Debug("Done Process Asicd Rx Buf", asicdRxBuf)
		case <-server.asicdComm.asicdSubSocketErrCh:
			server.logger.Err("Invalid Message from Asicd")
		case ribRxBuf := <-server.ribdComm.ribdSubSocketCh:
			server.logger.Debug("Process Rib Rx Buf", ribRxBuf)
			server.processRibdNotification(ribRxBuf)
			server.logger.Debug("Done Process Rib Rx Buf", ribRxBuf)
		case <-server.ribdComm.ribdSubSocketErrCh:
			server.logger.Err("Invalid Message from Ribd")
		case <-server.GetBulkData.SliceRefreshCh:
			server.logger.Debug("Refresh IntfConf Slice")
			server.RefreshIntfConfSlice()
			server.logger.Debug("Refresh NbrConf Slice")
			server.RefreshNbrConfSlice()
			server.logger.Debug("Refresh AreaConf Slice")
			server.RefreshAreaConfSlice()
			server.logger.Debug("Refresh Lsdb Slice")
			server.SendMsgToLsdbToRefreshSlice()
			<-server.MessagingChData.LsdbToServerChData.RefreshLsdbSliceDoneCh
			server.logger.Info("Ospf GetBulk Slice Refresh in progress")
			server.GetBulkData.SliceRefreshDoneCh <- true
			server.logger.Info("Ospf GetBulk Slice Refresh in done")
		}
	}
}
