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

package server

import (
	"fmt"
	"github.com/op/go-nanomsg"
	defs "l3/rib/ribdCommonDefs"
	"models/objects"
	"net"
	"ribd"
	"ribdInt"
	"strconv"
	"utils/clntUtils/clntIntfs"
	"utils/clntUtils/clntIntfs/arpdClntIntfs"
	"utils/clntUtils/clntIntfs/asicdClntIntfs"
	"utils/dbutils"
	"utils/logging"
	"utils/patriciaDB"
	"utils/policy"
	"utils/policy/policyCommonDefs"
)

type RIBdServerConfig struct {
	OrigConfigObject          interface{}
	NewConfigObject           interface{}
	OrigBulkRouteConfigObject []*ribdInt.IPv4RouteConfig
	Bulk                      bool
	BulkEnd                   bool
	AttrSet                   []bool
	Op                        defs.OpType //"add"/"del"/"update/get"
	PatchOp                   []*ribd.PatchOpInfo
	PolicyList                ApplyPolicyList
	AdditionalParams          interface{}
}

type V4IntfGetInfo struct {
	Count        int
	IPv4IntfList []*objects.IPv4IntfState
}
type V6IntfGetInfo struct {
	Count        int
	IPv6IntfList []*objects.IPv6IntfState
}

/*type PatchUpdateRouteInfo struct {
	OrigRoute *ribd.IPv4Route
	NewRoute  *ribd.IPv4Route
	Op        []*ribd.PatchOpInfo
}*/
type RIBDServer struct {
	Logger               *logging.Writer
	PolicyEngineDB       *policy.PolicyEngineDB
	GlobalPolicyEngineDB *policy.PolicyEngineDB
	TrackReachabilityCh  chan TrackReachabilityInfo
	RouteConfCh          chan RIBdServerConfig
	AsicdRouteCh         chan RIBdServerConfig
	ArpdRouteCh          chan RIBdServerConfig
	NotificationChannel  chan NotificationMsg
	NextHopInfoMap       map[NextHopInfoKey]NextHopInfo
	/*PolicyConditionConfCh  chan RIBdServerConfig
	PolicyActionConfCh     chan RIBdServerConfig
	PolicyStmtConfCh       chan RIBdServerConfig*/
	PolicyConfCh chan RIBdServerConfig
	//PolicyApplyCh       chan ApplyPolicyList
	PolicyUpdateApplyCh chan ApplyPolicyList
	DBRouteCh           chan RIBdServerConfig
	AcceptConfig        bool
	ServerUpCh          chan bool
	DBReadDone          chan bool
	V4IntfsGetDone      chan V4IntfGetInfo
	V6IntfsGetDone      chan V6IntfGetInfo
	PolicyConfDone      chan error
	DbHdl               *dbutils.DBUtil
	Clients             map[string]ClientIf
	//RouteInstallCh                 chan RouteParams
	ArpdClntPlugin   arpdClntIntfs.ArpdClntIntf
	AsicdPlugin      asicdClntIntfs.AsicdClntIntf
	AsicdSubSocketCh chan clntIntfs.NotifyMsg
}

const (
	PROTOCOL_NONE      = -1
	PROTOCOL_CONNECTED = 0
	PROTOCOL_STATIC    = 1
	PROTOCOL_OSPF      = 2
	PROTOCOL_BGP       = 3
	PROTOCOL_LAST      = 4
)

const (
	add = iota
	del
	delAll
	invalidate
)
const (
	Invalid   = -1
	FIBOnly   = 0
	FIBAndRIB = 1
	RIBOnly   = 2
)
const (
	SUB_ASICD = 0
)

type localDB struct {
	prefix     patriciaDB.Prefix
	isValid    bool
	precedence int
	nextHopIp  string
}
type IntfEntry struct {
	name string
}

var count int
var ConnectedRoutes []*ribdInt.Routes
var logger *logging.Writer
var AsicdSub *nanomsg.SubSocket
var RouteServiceHandler *RIBDServer
var IntfIdNameMap map[int32]IntfEntry
var IfNameToIfIndex map[string]int32
var GlobalPolicyEngineDB *policy.PolicyEngineDB
var PolicyEngineDB *policy.PolicyEngineDB
var PARAMSDIR string
var v4rtCount int
var v4routeCreatedTimeMap map[int]string
var v6rtCount int
var v6routeCreatedTimeMap map[int]string

var dbReqCount = 0
var dbReqCountLimit = 1
var dbReqCheckCount = 0
var dbReqCheckCountLimit = 5

/*
   Handle Interface down event
*/
func (ribdServiceHandler *RIBDServer) ProcessIPv4IntfDownEvent(ipAddr string, ifIndex int32) {
	logger.Debug("processIPv4IntfDownEvent")
	var ipMask net.IP
	ip, ipNet, err := net.ParseCIDR(ipAddr)
	if err != nil {
		return
	}
	ipMask = make(net.IP, 4)
	copy(ipMask, ipNet.Mask)
	ipAddrStr := ip.String()
	ipMaskStr := net.IP(ipMask).String()
	logger.Info(" processIPv4IntfDownEvent for  ipaddr ", ipAddrStr, " mask ", ipMaskStr)
	//deleteIPRoute(ConnectedRoutes[i].Ipaddr, defs.IPv4, ConnectedRoutes[i].Mask, "CONNECTED", ConnectedRoutes[i].NextHopIp, ribd.Int(ConnectedRoutes[i].IfIndex), FIBOnly, defs.RoutePolicyStateChangeNoChange)
	cfg := ribd.IPv4Route{
		DestinationNw: ipAddrStr,
		Protocol:      "CONNECTED",
		Cost:          0,
		NetworkMask:   ipMaskStr,
	}
	nextHop := ribd.NextHopInfo{
		NextHopIp:     "0.0.0.0",
		NextHopIntRef: strconv.Itoa(int(ifIndex)),
	}
	cfg.NextHop = make([]*ribd.NextHopInfo, 0)
	cfg.NextHop = append(cfg.NextHop, &nextHop)
	ribdServiceHandler.RouteConfCh <- RIBdServerConfig{
		OrigConfigObject: &cfg,
		Op:               defs.DelFIBOnly,
	}
}
func (ribdServiceHandler *RIBDServer) ProcessIPv6IntfDownEvent(ipAddr string, ifIndex int32) {
	logger.Debug("processIPv6IntfDownEvent")
	var ipMask net.IP
	ip, ipNet, err := net.ParseCIDR(ipAddr)
	if err != nil {
		return
	}
	ipMask = make(net.IP, 16)
	copy(ipMask, ipNet.Mask)
	ipAddrStr := ip.String()
	ipMaskStr := net.IP(ipMask).String()
	logger.Info(" processIPv6IntfDownEvent for  ipaddr ", ipAddrStr, " mask ", ipMaskStr)
	cfg := ribd.IPv6Route{
		DestinationNw: ipAddrStr,
		Protocol:      "CONNECTED",
		Cost:          0,
		NetworkMask:   ipMaskStr,
	}
	nextHop := ribd.NextHopInfo{
		NextHopIp:     "::",
		NextHopIntRef: strconv.Itoa(int(ifIndex)),
	}
	cfg.NextHop = make([]*ribd.NextHopInfo, 0)
	cfg.NextHop = append(cfg.NextHop, &nextHop)
	ribdServiceHandler.RouteConfCh <- RIBdServerConfig{
		OrigConfigObject: &cfg,
		Op:               defs.Delv6FIBOnly,
	}
}

/*
   Handle Interface up event
*/
func (ribdServiceHandler *RIBDServer) ProcessIPv4IntfUpEvent(ipAddr string, ifIndex int32) {
	logger.Debug("processIPv4IntfUpEvent")
	var ipMask net.IP
	ip, ipNet, err := net.ParseCIDR(ipAddr)
	if err != nil {
		return
	}
	ipMask = make(net.IP, 4)
	copy(ipMask, ipNet.Mask)
	ipAddrStr := ip.String()
	ipMaskStr := net.IP(ipMask).String()
	logger.Info(" processIPv4IntfUpEvent for  ipaddr ", ipAddrStr, " mask ", ipMaskStr)
	for i := 0; i < len(ConnectedRoutes); i++ {
		//logger.Info("Current state of this connected route is ", ConnectedRoutes[i].IsValid)
		if ConnectedRoutes[i].Ipaddr == ipAddrStr && ConnectedRoutes[i].Mask == ipMaskStr && ConnectedRoutes[i].IsValid == false {
			if ifIndex != -1 && ConnectedRoutes[i].IfIndex != ribdInt.Int(ifIndex) {
				continue
			}
			logger.Info("Add this route with destAddress = ", ConnectedRoutes[i].Ipaddr, " nwMask = ", ConnectedRoutes[i].Mask)

			ConnectedRoutes[i].IsValid = true
			//			policyRoute := ribdInt.Routes{Ipaddr: ConnectedRoutes[i].Ipaddr, IPAddrType: ribdInt.Int(defs.IPv4), Mask: ConnectedRoutes[i].Mask, NextHopIp: ConnectedRoutes[i].NextHopIp, IfIndex: ConnectedRoutes[i].IfIndex, Metric: ConnectedRoutes[i].Metric, Prototype: ConnectedRoutes[i].Prototype}
			//			params := RouteParams{destNetIp: ConnectedRoutes[i].Ipaddr, ipType: defs.IPv4, networkMask: ConnectedRoutes[i].Mask, nextHopIp: ConnectedRoutes[i].NextHopIp, nextHopIfIndex: ribd.Int(ConnectedRoutes[i].IfIndex), metric: ribd.Int(ConnectedRoutes[i].Metric), routeType: ribd.Int(ConnectedRoutes[i].Prototype), sliceIdx: ribd.Int(ConnectedRoutes[i].SliceIdx), createType: FIBOnly, deleteType: Invalid}
			//			PolicyEngineFilter(policyRoute, policyCommonDefs.PolicyPath_Import, params)
			cfg := ribd.IPv4Route{
				DestinationNw: ipAddrStr,
				Protocol:      "CONNECTED",
				Cost:          0,
				NetworkMask:   ipMaskStr,
			}
			nextHop := ribd.NextHopInfo{
				NextHopIp:     "0.0.0.0",
				NextHopIntRef: strconv.Itoa(int(ifIndex)),
			}
			cfg.NextHop = make([]*ribd.NextHopInfo, 0)
			cfg.NextHop = append(cfg.NextHop, &nextHop)

			ribdServiceHandler.RouteConfCh <- RIBdServerConfig{
				OrigConfigObject: &cfg,
				Op:               defs.AddFIBOnly,
				AdditionalParams: ribd.Int(ConnectedRoutes[i].SliceIdx),
			}
		}
	}
}
func (ribdServiceHandler *RIBDServer) ProcessIPv6IntfUpEvent(ipAddr string, ifIndex int32) {
	logger.Debug("processIPv6IntfUpEvent")
	var ipMask net.IP
	ip, ipNet, err := net.ParseCIDR(ipAddr)
	if err != nil {
		return
	}
	ipMask = make(net.IP, 16)
	copy(ipMask, ipNet.Mask)
	ipAddrStr := ip.String()
	ipMaskStr := net.IP(ipMask).String()
	logger.Info(" processIPv6IntfUpEvent for  ipaddr ", ipAddrStr, " mask ", ipMaskStr)
	for i := 0; i < len(ConnectedRoutes); i++ {
		//logger.Info("Current state of this connected route is ", ConnectedRoutes[i].IsValid)
		if ConnectedRoutes[i].Ipaddr == ipAddrStr && ConnectedRoutes[i].Mask == ipMaskStr && ConnectedRoutes[i].IsValid == false {
			if ifIndex != -1 && ConnectedRoutes[i].IfIndex != ribdInt.Int(ifIndex) {
				continue
			}
			logger.Info("Add this route with destAddress = ", ConnectedRoutes[i].Ipaddr, " nwMask = ", ConnectedRoutes[i].Mask)

			ConnectedRoutes[i].IsValid = true
			//			policyRoute := ribdInt.Routes{Ipaddr: ConnectedRoutes[i].Ipaddr, IPAddrType: ribdInt.Int(defs.IPv6), Mask: ConnectedRoutes[i].Mask, NextHopIp: ConnectedRoutes[i].NextHopIp, IfIndex: ConnectedRoutes[i].IfIndex, Metric: ConnectedRoutes[i].Metric, Prototype: ConnectedRoutes[i].Prototype}
			//			params := RouteParams{destNetIp: ConnectedRoutes[i].Ipaddr, ipType: defs.IPv6, networkMask: ConnectedRoutes[i].Mask, nextHopIp: ConnectedRoutes[i].NextHopIp, nextHopIfIndex: ribd.Int(ConnectedRoutes[i].IfIndex), metric: ribd.Int(ConnectedRoutes[i].Metric), routeType: ribd.Int(ConnectedRoutes[i].Prototype), sliceIdx: ribd.Int(ConnectedRoutes[i].SliceIdx), createType: FIBOnly, deleteType: Invalid}
			//			PolicyEngineFilter(policyRoute, policyCommonDefs.PolicyPath_Import, params)
			cfg := ribd.IPv6Route{
				DestinationNw: ipAddrStr,
				Protocol:      "CONNECTED",
				Cost:          0,
				NetworkMask:   ipMaskStr,
			}
			nextHop := ribd.NextHopInfo{
				NextHopIp:     "::",
				NextHopIntRef: strconv.Itoa(int(ifIndex)),
			}
			cfg.NextHop = make([]*ribd.NextHopInfo, 0)
			cfg.NextHop = append(cfg.NextHop, &nextHop)

			ribdServiceHandler.RouteConfCh <- RIBdServerConfig{
				OrigConfigObject: &cfg,
				Op:               defs.Addv6FIBOnly,
				AdditionalParams: ribd.Int(ConnectedRoutes[i].SliceIdx),
			}
		}
	}
}

func (server *RIBDServer) GetLogicalIntfInfo() {
	logger.Debug("Getting Logical Interfaces from asicd")
	var currMarker int
	var count int
	count = 100
	for {
		logger.Info("Getting ", count, "GetBulkLogicalIntf objects from currMarker:", currMarker)
		bulkInfo, err := server.AsicdPlugin.GetBulkLogicalIntfState(currMarker, count)
		if err != nil {
			logger.Info("GetBulkLogicalIntfState with err ", err)
			return
		}
		if bulkInfo.Count == 0 {
			logger.Info("0 objects returned from GetBulkLogicalIntfState")
			return
		}
		logger.Info("len(bulkInfo.GetBulkLogicalIntfState)  = ", len(bulkInfo.LogicalIntfStateList), " num objects returned = ", bulkInfo.Count)
		for i := 0; i < int(bulkInfo.Count); i++ {
			ifId := (bulkInfo.LogicalIntfStateList[i].IfIndex)
			logger.Info("logical interface = ", bulkInfo.LogicalIntfStateList[i].Name, "ifId = ", ifId)
			if IntfIdNameMap == nil {
				IntfIdNameMap = make(map[int32]IntfEntry)
			}
			intfEntry := IntfEntry{name: bulkInfo.LogicalIntfStateList[i].Name}
			IntfIdNameMap[ifId] = intfEntry
			if IfNameToIfIndex == nil {
				IfNameToIfIndex = make(map[string]int32)
			}
			IfNameToIfIndex[bulkInfo.LogicalIntfStateList[i].Name] = ifId
		}
		if bulkInfo.More == false {
			logger.Info("more returned as false, so no more get bulks")
			return
		}
		currMarker = int(bulkInfo.EndIdx)
	}
}
func (server *RIBDServer) GetVlanInfo() {
	logger.Debug("Getting vlans from asicd")
	var currMarker int
	var count int
	count = 100
	for {
		logger.Info("Getting ", count, "GetBulkVlan objects from currMarker:", currMarker)
		bulkInfo, err := server.AsicdPlugin.GetBulkVlanState(currMarker, count)
		if err != nil {
			logger.Info("GetBulkVlan with err ", err)
			return
		}
		if bulkInfo.Count == 0 {
			logger.Info("0 objects returned from GetBulkVlan")
			return
		}
		logger.Info("len(bulkInfo.GetBulkVlan)  = ", len(bulkInfo.VlanStateList), " num objects returned = ", bulkInfo.Count)
		for i := 0; i < int(bulkInfo.Count); i++ {
			ifId := (bulkInfo.VlanStateList[i].IfIndex)
			logger.Info("vlan = ", bulkInfo.VlanStateList[i].VlanId, "ifId = ", ifId)
			if IntfIdNameMap == nil {
				IntfIdNameMap = make(map[int32]IntfEntry)
			}
			intfEntry := IntfEntry{name: bulkInfo.VlanStateList[i].VlanName}
			IntfIdNameMap[ifId] = intfEntry
			if IfNameToIfIndex == nil {
				IfNameToIfIndex = make(map[string]int32)
			}
			IfNameToIfIndex[bulkInfo.VlanStateList[i].VlanName] = ifId
		}
		if bulkInfo.More == false {
			logger.Info("more returned as false, so no more get bulks")
			return
		}
		currMarker = int(bulkInfo.EndIdx)
	}
}
func (server *RIBDServer) GetPortInfo() {
	logger.Debug("Getting ports from asicd")
	var currMarker int
	var count int
	count = 100
	for {
		logger.Info("Getting ", count, "objects from currMarker:", currMarker)
		bulkInfo, err := server.AsicdPlugin.GetBulkPortState(currMarker, count)
		if err != nil {
			logger.Info("GetBulkPortState with err ", err)
			return
		}
		if bulkInfo.Count == 0 {
			logger.Info("0 objects returned from GetBulkPortState")
			return
		}
		logger.Info("len(bulkInfo.PortStateList)  = ", len(bulkInfo.PortStateList), " num objects returned = ", bulkInfo.Count)
		for i := 0; i < int(bulkInfo.Count); i++ {
			ifId := bulkInfo.PortStateList[i].IfIndex
			if IntfIdNameMap == nil {
				IntfIdNameMap = make(map[int32]IntfEntry)
			}
			intfEntry := IntfEntry{name: bulkInfo.PortStateList[i].Name}
			IntfIdNameMap[ifId] = intfEntry
			if IfNameToIfIndex == nil {
				IfNameToIfIndex = make(map[string]int32)
			}
			IfNameToIfIndex[bulkInfo.PortStateList[i].Name] = ifId
			logger.Info("ifId = ", ifId, "IntfIdNameMap[", ifId, "] = ", IntfIdNameMap[ifId], "IfNameToIfIndex[", bulkInfo.PortStateList[i].Name, "] = ", IfNameToIfIndex[bulkInfo.PortStateList[i].Name])
		}
		if bulkInfo.More == false {
			logger.Info("more returned as false, so no more get bulks")
			return
		}
		currMarker = int(bulkInfo.EndIdx)
	}
}
func (server *RIBDServer) GetIntfInfo() {
	server.GetPortInfo()
	server.GetVlanInfo()
	server.GetLogicalIntfInfo()
}
func (ribdServiceHandler *RIBDServer) AcceptConfigActions() {
	logger.Info("AcceptConfigActions: Setting AcceptConfig to true")
	RouteServiceHandler.AcceptConfig = true
	ribdServiceHandler.GetIntfInfo()
	logger.Info("adding fetchv4 to asicdroutech")
	ribdServiceHandler.AsicdRouteCh <- RIBdServerConfig{Op: defs.AsicdFetchv4}
	v4IntfsGetDone := <-ribdServiceHandler.V4IntfsGetDone
	logger.Info("adding fetchv6 to asicdroutech")
	//getV4ConnectedRoutes()
	ribdServiceHandler.AsicdRouteCh <- RIBdServerConfig{Op: defs.AsicdFetchv6}
	v6IntfsGetDone := <-ribdServiceHandler.V6IntfsGetDone
	//getV6ConnectedRoutes()
	logger.Info("creating v4 and v6 routes")
	CreateV4ConnectedRoutes(v4IntfsGetDone.Count, v4IntfsGetDone.IPv4IntfList)
	CreateV6ConnectedRoutes(v6IntfsGetDone.Count, v6IntfsGetDone.IPv6IntfList)
	//update dbRouteCh to fetch route data
	ribdServiceHandler.DBRouteCh <- RIBdServerConfig{Op: defs.DBFetch}
	dbRead := <-ribdServiceHandler.DBReadDone
	logger.Info("Received dbread: ")
	if dbRead != true {
		logger.Err("DB read failed")
	}
	//	go ribdServiceHandler.SetupEventHandler(AsicdSub, asicdCommonDefs.PUB_SOCKET_ADDR, SUB_ASICD)
	logger.Info("All set to signal start the RIBd server")
	ribdServiceHandler.ServerUpCh <- true
}

func (ribdServiceHandler *RIBDServer) InitializeGlobalPolicyDB() *policy.PolicyEngineDB {
	ribdServiceHandler.GlobalPolicyEngineDB = policy.NewPolicyEngineDB(logger)
	ribdServiceHandler.GlobalPolicyEngineDB.SetDefaultImportPolicyActionFunc(defaultImportPolicyEngineActionFunc)
	ribdServiceHandler.GlobalPolicyEngineDB.SetDefaultExportPolicyActionFunc(defaultExportPolicyEngineActionFunc)
	ribdServiceHandler.GlobalPolicyEngineDB.SetIsEntityPresentFunc(DoesRouteExist)
	ribdServiceHandler.GlobalPolicyEngineDB.SetEntityUpdateFunc(UpdateRouteAndPolicyDB)
	ribdServiceHandler.GlobalPolicyEngineDB.SetActionFunc(policyCommonDefs.PolicyActionTypeRouteDisposition, policyEngineRouteDispositionAction)
	ribdServiceHandler.GlobalPolicyEngineDB.SetActionFunc(policyCommonDefs.PolicyActionTypeRouteRedistribute, policyEngineActionRedistribute)
	ribdServiceHandler.GlobalPolicyEngineDB.SetActionFunc(policyCommonDefs.PolicyActionTypeNetworkStatementAdvertise, policyEngineActionNetworkStatementAdvertise)
	ribdServiceHandler.GlobalPolicyEngineDB.SetActionFunc(policyCommonDefs.PoilcyActionTypeSetAdminDistance, policyEngineActionSetAdminDistance)
	ribdServiceHandler.GlobalPolicyEngineDB.SetUndoActionFunc(policyCommonDefs.PolicyActionTypeRouteDisposition, policyEngineUndoRouteDispositionAction)
	ribdServiceHandler.GlobalPolicyEngineDB.SetUndoActionFunc(policyCommonDefs.PolicyActionTypeRouteRedistribute, policyEngineActionUndoRedistribute)
	ribdServiceHandler.GlobalPolicyEngineDB.SetUndoActionFunc(policyCommonDefs.PoilcyActionTypeSetAdminDistance, policyEngineActionUndoSetAdminDistance)
	ribdServiceHandler.GlobalPolicyEngineDB.SetUndoActionFunc(policyCommonDefs.PolicyActionTypeNetworkStatementAdvertise, policyEngineActionUndoNetworkStatemenAdvertiseAction)
	ribdServiceHandler.GlobalPolicyEngineDB.SetTraverseAndApplyPolicyFunc(policyEngineTraverseAndApply)
	ribdServiceHandler.GlobalPolicyEngineDB.SetTraverseAndReversePolicyFunc(policyEngineTraverseAndReverse)
	ribdServiceHandler.GlobalPolicyEngineDB.SetGetPolicyEntityMapIndexFunc(getPolicyRouteMapIndex)
	ribdServiceHandler.GlobalPolicyEngineDB.Global = true //this policy engine does not apply the policies
	return ribdServiceHandler.GlobalPolicyEngineDB
}

func (ribdServiceHandler *RIBDServer) InitializePolicyDB() *policy.PolicyEngineDB {
	ribdServiceHandler.PolicyEngineDB = policy.NewPolicyEngineDB(logger)
	ribdServiceHandler.PolicyEngineDB.SetDefaultImportPolicyActionFunc(defaultImportPolicyEngineActionFunc)
	ribdServiceHandler.PolicyEngineDB.SetDefaultExportPolicyActionFunc(defaultExportPolicyEngineActionFunc)
	ribdServiceHandler.PolicyEngineDB.SetIsEntityPresentFunc(DoesRouteExist)
	ribdServiceHandler.PolicyEngineDB.SetEntityUpdateFunc(UpdateRouteAndPolicyDB)
	ribdServiceHandler.PolicyEngineDB.SetActionFunc(policyCommonDefs.PolicyActionTypeRouteDisposition, policyEngineRouteDispositionAction)
	ribdServiceHandler.PolicyEngineDB.SetActionFunc(policyCommonDefs.PolicyActionTypeRouteRedistribute, policyEngineActionRedistribute)
	ribdServiceHandler.PolicyEngineDB.SetActionFunc(policyCommonDefs.PolicyActionTypeNetworkStatementAdvertise, policyEngineActionNetworkStatementAdvertise)
	ribdServiceHandler.PolicyEngineDB.SetActionFunc(policyCommonDefs.PoilcyActionTypeSetAdminDistance, policyEngineActionSetAdminDistance)
	ribdServiceHandler.PolicyEngineDB.SetUndoActionFunc(policyCommonDefs.PolicyActionTypeRouteDisposition, policyEngineUndoRouteDispositionAction)
	ribdServiceHandler.PolicyEngineDB.SetUndoActionFunc(policyCommonDefs.PolicyActionTypeRouteRedistribute, policyEngineActionUndoRedistribute)
	ribdServiceHandler.PolicyEngineDB.SetUndoActionFunc(policyCommonDefs.PoilcyActionTypeSetAdminDistance, policyEngineActionUndoSetAdminDistance)
	ribdServiceHandler.PolicyEngineDB.SetUndoActionFunc(policyCommonDefs.PolicyActionTypeNetworkStatementAdvertise, policyEngineActionUndoNetworkStatemenAdvertiseAction)
	ribdServiceHandler.PolicyEngineDB.SetTraverseAndApplyPolicyFunc(policyEngineTraverseAndApply)
	ribdServiceHandler.PolicyEngineDB.SetTraverseAndReversePolicyFunc(policyEngineTraverseAndReverse)
	ribdServiceHandler.PolicyEngineDB.SetGetPolicyEntityMapIndexFunc(getPolicyRouteMapIndex)
	return ribdServiceHandler.PolicyEngineDB
}
func NewRIBDServicesHandler(dbHdl *dbutils.DBUtil, loggerC *logging.Writer) *RIBDServer {
	V4RouteInfoMap = patriciaDB.NewTrie()
	V6RouteInfoMap = patriciaDB.NewTrie()
	ribdServicesHandler := &RIBDServer{}
	ribdServicesHandler.Logger = loggerC
	logger = loggerC
	localRouteEventsDB = make([]RouteEventInfo, 0)
	RedistributeRouteMap = make(map[string][]RedistributeRouteInfo)
	ribdServicesHandler.Clients = make(map[string]ClientIf)
	TrackReachabilityMap = make(map[string][]string)
	v4routeCreatedTimeMap = make(map[int]string)
	v6routeCreatedTimeMap = make(map[int]string)
	RouteProtocolTypeMapDB = make(map[string]int)
	ReverseRouteProtoTypeMapDB = make(map[int]string)
	ProtocolAdminDistanceMapDB = make(map[string]RouteDistanceConfig)
	PublisherInfoMap = make(map[string]PublisherMapInfo)
	ribdServicesHandler.NextHopInfoMap = make(map[NextHopInfoKey]NextHopInfo)
	ribdServicesHandler.AsicdSubSocketCh = make(chan clntIntfs.NotifyMsg)
	ribdServicesHandler.TrackReachabilityCh = make(chan TrackReachabilityInfo, 1000)
	ribdServicesHandler.RouteConfCh = make(chan RIBdServerConfig, 200000)
	ribdServicesHandler.AsicdRouteCh = make(chan RIBdServerConfig, 200000)
	ribdServicesHandler.ArpdRouteCh = make(chan RIBdServerConfig, 5000)
	ribdServicesHandler.NotificationChannel = make(chan NotificationMsg, 5000)
	/*	ribdServicesHandler.PolicyConditionConfCh = make(chan RIBdServerConfig, 5000)
		ribdServicesHandler.PolicyActionConfCh = make(chan RIBdServerConfig, 5000)
		ribdServicesHandler.PolicyStmtConfCh = make(chan RIBdServerConfig, 5000)*/
	ribdServicesHandler.PolicyConfCh = make(chan RIBdServerConfig, 5000)
	//ribdServicesHandler.PolicyApplyCh = make(chan ApplyPolicyList, 100)
	ribdServicesHandler.PolicyUpdateApplyCh = make(chan ApplyPolicyList, 100)
	ribdServicesHandler.DBRouteCh = make(chan RIBdServerConfig, 200000)
	ribdServicesHandler.ServerUpCh = make(chan bool)
	ribdServicesHandler.DBReadDone = make(chan bool)
	ribdServicesHandler.V4IntfsGetDone = make(chan V4IntfGetInfo)
	ribdServicesHandler.V6IntfsGetDone = make(chan V6IntfGetInfo)
	ribdServicesHandler.PolicyConfDone = make(chan error)
	ribdServicesHandler.DbHdl = dbHdl
	RouteServiceHandler = ribdServicesHandler
	//ribdServicesHandler.RouteInstallCh = make(chan RouteParams)
	BuildRouteProtocolTypeMapDB()
	BuildProtocolAdminDistanceMapDB()
	BuildPublisherMap()
	PolicyEngineDB = ribdServicesHandler.InitializePolicyDB()
	GlobalPolicyEngineDB = ribdServicesHandler.InitializeGlobalPolicyDB()
	return ribdServicesHandler
}
func (s *RIBDServer) InitServer() {
	/*	sigChan := make(chan os.Signal, 1)
		signalList := []os.Signal{syscall.SIGHUP}
		signal.Notify(sigChan, signalList...)*/
	go s.ListenToClientStateChanges()
	//go s.SigHandler(sigChan)
	go s.StartRouteProcessServer()
	go s.StartDBServer()
	go s.StartPolicyServer()
	go s.NotificationServer()
	go s.StartAsicdServer()
	go s.StartArpdServer()

}
func (ribdServiceHandler *RIBDServer) StartServer(paramsDir string) {
	ribdServiceHandler.InitServer()
	logger.Info("Starting RIB server comment out logger. calls")
	DummyRouteInfoRecord.protocol = PROTOCOL_NONE
	configFile := paramsDir + "/clients.json"
	logger.Info(fmt.Sprintln("configfile = ", configFile))
	PARAMSDIR = paramsDir
	ribdServiceHandler.UpdatePolicyObjectsFromDB() //(paramsDir)
	//ribdServiceHandler.ConnectToClients(configFile)
	ribdServiceHandler.AcceptConfigActions()
	logger.Info("Starting the server loop")
	count := 0
	for {
		if !RouteServiceHandler.AcceptConfig {
			if count%10000 == 0 {
				//				logger.Debug("RIBD not ready to accept config")
			}
			count++
			continue
		}
		select {
		/*case list := <-ribdServiceHandler.PolicyApplyCh:
		logger.Debug("received message on PolicyApplyCh channel")
		//update the local policyEngineDB
		ribdServiceHandler.UpdateApplyPolicyList(list.ApplyList, list.UndoList, true, PolicyEngineDB)
		ribdServiceHandler.PolicyUpdateApplyCh <- list)*/
		case info := <-ribdServiceHandler.TrackReachabilityCh:
			logger.Debug("received message on TrackReachabilityCh channel")
			ribdServiceHandler.TrackReachabilityStatus(info.IpAddr, info.Protocol, info.Op)
		case msg := <-ribdServiceHandler.AsicdSubSocketCh:
			ribdServiceHandler.processAsicdNotification(msg)
		}
	}
}
func (s *RIBDServer) StopServer() {
	logger.Debug("StopServer")
	//clean up IPv4RouteState* from DB
	s.DbHdl.DeleteObjectWithKeyFromDb("IPv4RouteState*")
	//clean up IPv6RouteState* from DB
	s.DbHdl.DeleteObjectWithKeyFromDb("IPv6RouteState*")
}
