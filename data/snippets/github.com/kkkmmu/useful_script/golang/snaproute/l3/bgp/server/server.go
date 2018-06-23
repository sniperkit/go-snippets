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
//  _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

// server.go
package server

import (
	"bgpd"
	"encoding/json"
	"errors"
	"fmt"
	"l3/bgp/config"
	"l3/bgp/fsm"
	"l3/bgp/packet"
	bgppolicy "l3/bgp/policy"
	bgprib "l3/bgp/rib"
	"l3/bgp/utils"
	"net"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"utils/dbutils"
	"utils/eventUtils"
	"utils/logging"
	"utils/netUtils"
	"utils/patriciaDB"
	utilspolicy "utils/policy"
	"utils/policy/policyCommonDefs"
	"utils/statedbclient"
)

type GlobalUpdate struct {
	BGPConfig *bgpd.BGPGlobal
	OldConfig config.GlobalConfig
	NewConfig config.GlobalConfig
	AttrSet   []bool
	PatchOp   []*bgpd.PatchOpInfo
	Op        string
}

type PeerUpdate struct {
	BGPPeer  interface{}
	PeerType string
	OldPeer  config.NeighborConfig
	NewPeer  config.NeighborConfig
	AttrSet  []bool
	PatchOp  []*bgpd.PatchOpInfo
	Op       string
}

type PeerGroupUpdate struct {
	OldGroup config.PeerGroupConfig
	NewGroup config.PeerGroupConfig
	AttrSet  []bool
}

type AggUpdate struct {
	OldAgg  config.BGPAggregate
	NewAgg  config.BGPAggregate
	AttrSet []bool
}

type NSUpdate struct {
	OldNS   config.BGPNetworkStatement
	NewNS   config.BGPNetworkStatement
	AttrSet []bool
}

type PolicyParams struct {
	CreateType      int
	DeleteType      int
	route           *bgprib.Route
	dest            *bgprib.Destination
	updated         *(map[uint32]map[*bgprib.Path][]*bgprib.Destination)
	withdrawn       *([]*bgprib.Destination)
	updatedAddPaths *([]*bgprib.Destination)
}

type NSPolicyParams struct {
	CreateType       int
	DeleteType       int
	Traverse         bool
	TraverseRev      bool
	NS               *config.BGPNetworkStatement
	Accept           int
	AddPolicyStmt    []utilspolicy.PolicyStmt
	RemovePolicyStmt []utilspolicy.PolicyStmt
}

type IntfEntry struct {
	Name string
}

type BGPServer struct {
	logger           *logging.Writer
	policyManager    *bgppolicy.BGPPolicyManager
	locRibPE         map[uint32]*bgppolicy.LocRibPolicyEngine
	ribInPE          *bgppolicy.AdjRibPPolicyEngine
	ribOutPE         *bgppolicy.AdjRibPPolicyEngine
	networkStmtPE    *bgppolicy.NetworkStmtPolicyEngine
	listener         *net.TCPListener
	listenerIPv6     *net.TCPListener
	ifaceMgr         *utils.InterfaceMgr
	BgpConfig        config.Bgp
	GlobalConfigCh   chan GlobalUpdate
	AddPeerCh        chan PeerUpdate
	RemPeerCh        chan config.NeighborConfig
	AddPeerGroupCh   chan PeerGroupUpdate
	RemPeerGroupCh   chan config.PeerGroupConfig
	AddAggCh         chan AggUpdate
	RemAggCh         chan config.BGPAggregate
	AddNSCh          chan NSUpdate
	RemNSCh          chan config.BGPNetworkStatement
	PeerFSMConnCh    chan fsm.PeerFSMConn
	PeerConnEstCh    chan string
	PeerConnBrokenCh chan string
	PeerCommandCh    chan config.PeerCommand
	ReachabilityCh   chan config.ReachabilityInfo
	BGPPktSrcCh      chan *packet.BGPPktSrc
	BfdCh            chan config.BfdInfo
	IntfCh           chan config.IntfStateInfo
	IntfMapCh        chan config.IntfMapInfo
	RoutesCh         chan *config.RouteCh
	acceptCh         chan *net.TCPConn
	ServerUpCh       chan bool
	GlobalCfgDone    bool

	NeighborMutex     sync.RWMutex
	PeerMap           map[string]*Peer
	ifaceNeighbors    map[config.PeerAddressType]map[int32]*Peer
	Neighbors         []*Peer
	LocRib            *bgprib.LocRib
	ConnRoutesPath    *bgprib.Path
	DefaultRoutesPath *bgprib.Path
	NSRoutesPathMap   map[string]*bgprib.Path
	NSRoutes          []*config.BGPNetworkStatement
	NSRoutesMutex     sync.RWMutex
	IfIndexPeerMap    map[int32][]string
	IntfIdNameMap     map[int32]IntfEntry
	IfNameToIfIndex   map[string]int32
	RedistributionMap map[string]string
	ifaceIP           net.IP
	AddPathCount      int
	// all managers
	IntfMgr    config.IntfStateMgrIntf
	routeMgr   config.RouteMgrIntf
	bfdMgr     config.BfdMgrIntf
	stateDBMgr statedbclient.StateDBClient
	eventDbHdl *dbutils.DBUtil
}

func NewBGPServer(logger *logging.Writer, policyManager *bgppolicy.BGPPolicyManager, iMgr config.IntfStateMgrIntf,
	rMgr config.RouteMgrIntf, bMgr config.BfdMgrIntf, sDBMgr statedbclient.StateDBClient) *BGPServer {
	bgpServer := &BGPServer{}
	bgpServer.logger = logger
	bgpServer.policyManager = policyManager
	bgpServer.ifaceMgr = utils.NewInterfaceMgr(logger)
	bgpServer.GlobalCfgDone = false
	bgpServer.GlobalConfigCh = make(chan GlobalUpdate)
	bgpServer.AddPeerCh = make(chan PeerUpdate)
	bgpServer.RemPeerCh = make(chan config.NeighborConfig)
	bgpServer.AddPeerGroupCh = make(chan PeerGroupUpdate)
	bgpServer.RemPeerGroupCh = make(chan config.PeerGroupConfig)
	bgpServer.AddAggCh = make(chan AggUpdate)
	bgpServer.RemAggCh = make(chan config.BGPAggregate)
	bgpServer.AddNSCh = make(chan NSUpdate)
	bgpServer.RemNSCh = make(chan config.BGPNetworkStatement)
	bgpServer.PeerFSMConnCh = make(chan fsm.PeerFSMConn, 50)
	bgpServer.PeerConnEstCh = make(chan string)
	bgpServer.PeerConnBrokenCh = make(chan string)
	bgpServer.PeerCommandCh = make(chan config.PeerCommand)
	bgpServer.ReachabilityCh = make(chan config.ReachabilityInfo)
	bgpServer.BGPPktSrcCh = make(chan *packet.BGPPktSrc)
	bgpServer.BfdCh = make(chan config.BfdInfo)
	bgpServer.IntfCh = make(chan config.IntfStateInfo)
	bgpServer.IntfMapCh = make(chan config.IntfMapInfo)
	bgpServer.RoutesCh = make(chan *config.RouteCh)
	bgpServer.ServerUpCh = make(chan bool)

	bgpServer.NeighborMutex = sync.RWMutex{}
	bgpServer.PeerMap = make(map[string]*Peer)
	bgpServer.NSRoutesPathMap = make(map[string]*bgprib.Path)
	bgpServer.NSRoutes = make([]*config.BGPNetworkStatement, 0)
	bgpServer.NSRoutesMutex = sync.RWMutex{}
	bgpServer.ifaceNeighbors = make(map[config.PeerAddressType]map[int32]*Peer)
	bgpServer.ifaceNeighbors[config.PeerAddressV4] = make(map[int32]*Peer)
	bgpServer.ifaceNeighbors[config.PeerAddressV6] = make(map[int32]*Peer)

	bgpServer.Neighbors = make([]*Peer, 0)
	bgpServer.IntfMgr = iMgr
	bgpServer.routeMgr = rMgr
	bgpServer.bfdMgr = bMgr
	bgpServer.stateDBMgr = sDBMgr
	bgpServer.LocRib = bgprib.NewLocRib(logger, rMgr, sDBMgr, &bgpServer.BgpConfig.Global.Config)
	bgpServer.IfNameToIfIndex = make(map[string]int32)
	bgpServer.IntfIdNameMap = make(map[int32]IntfEntry)
	bgpServer.IfIndexPeerMap = make(map[int32][]string)
	bgpServer.RedistributionMap = make(map[string]string)
	bgpServer.ifaceIP = nil
	bgpServer.AddPathCount = 0
	bgpServer.initGlobalConfig()
	bgpServer.initPolicyEngines()
	return bgpServer
}

func (s *BGPServer) initGlobalConfig() {
	s.BgpConfig = config.Bgp{}
	s.BgpConfig.Afs = make(map[uint32]*config.AddressFamily)
	for _, pfNumber := range packet.ProtocolFamilyMap {
		s.BgpConfig.Afs[pfNumber] = &config.AddressFamily{}
	}
}

func (s *BGPServer) initPolicyEngines() {
	type TraverseFuncMap struct {
		ApplyFunc   utilspolicy.EntityTraverseAndApplyPolicyfunc
		ReverseFunc utilspolicy.EntityTraverseAndReversePolicyfunc
	}
	var traverseFuncMap = map[packet.AFI]TraverseFuncMap{
		packet.AfiIP:  TraverseFuncMap{s.TraverseAndApplyBGPRib, s.TrAndRevAggForIPv4},
		packet.AfiIP6: TraverseFuncMap{s.TraverseAndApplyBGPRib, s.TrAndRevAggForIPv6},
	}
	var actionFunc bgppolicy.PolicyActionFunc
	actionFuncMap := make(map[int]bgppolicy.PolicyActionFunc)
	s.locRibPE = make(map[uint32]*bgppolicy.LocRibPolicyEngine)

	actionFunc.ApplyFunc = s.ApplyAggregateAction
	actionFunc.UndoFunc = s.UndoAggregateAction
	actionFuncMap[policyCommonDefs.PolicyActionTypeAggregate] = actionFunc
	s.logger.Infof("BGPServer: actionfuncmap=%v", actionFuncMap)

	for _, pfNumber := range packet.ProtocolFamilyMap {
		locRibPE := bgppolicy.NewLocRibPolicyEngine(s.logger)
		locRibPE.SetEntityUpdateFunc(s.UpdateRouteAndPolicyDB)
		locRibPE.SetIsEntityPresentFunc(s.DoesRouteExist)
		locRibPE.SetActionFuncs(actionFuncMap)
		afi, _ := packet.GetAfiSafi(pfNumber)
		traverseFuncs, ok := traverseFuncMap[afi]
		if !ok {
			s.logger.Crit("BGPServer: traverse funcs not found for address family", afi)
		}
		locRibPE.SetTraverseFuncs(traverseFuncs.ApplyFunc, traverseFuncs.ReverseFunc)
		s.locRibPE[pfNumber] = locRibPE
		//s.policyManager.AddPolicyEngine(locRibPE)
	}

	s.ribInPE = bgppolicy.NewAdjRibPolicyEngine(s.logger)
	s.ribOutPE = bgppolicy.NewAdjRibPolicyEngine(s.logger)

	actionFunc.ApplyFunc = s.ApplyAdjRIBAction
	actionFunc.UndoFunc = s.UndoAdjRIBAction

	s.ribInPE.SetEntityUpdateFunc(s.UpdateAdjRIBInRouteAndPolicyDB)
	s.ribInPE.SetIsEntityPresentFunc(s.DoesAdjRIBInRouteExist)
	actionFuncMap = make(map[int]bgppolicy.PolicyActionFunc)
	actionFuncMap[policyCommonDefs.PolicyActionTypeRIBIn] = actionFunc
	s.ribInPE.SetActionFuncs(actionFuncMap)
	s.ribInPE.SetTraverseFuncs(s.TraverseAndApplyAdjRibIn, s.TraverseAndReverseAdjRIBIn)
	s.policyManager.AddPolicyEngine(s.ribInPE)

	s.ribOutPE.SetEntityUpdateFunc(s.UpdateAdjRIBOutRouteAndPolicyDB)
	s.ribOutPE.SetIsEntityPresentFunc(s.DoesAdjRIBOutRouteExist)
	actionFuncMap = make(map[int]bgppolicy.PolicyActionFunc)
	actionFuncMap[policyCommonDefs.PolicyActionTypeRIBOut] = actionFunc
	s.ribOutPE.SetActionFuncs(actionFuncMap)
	s.ribOutPE.SetTraverseFuncs(s.TraverseAndApplyAdjRibOut, s.TraverseAndReverseAdjRIBOut)
	s.policyManager.AddPolicyEngine(s.ribOutPE)

	s.networkStmtPE = bgppolicy.NewNetworkStmtPolicyEngine(s.logger)

	actionFunc.ApplyFunc = s.ApplyNSAction
	actionFunc.UndoFunc = s.UndoNSAction

	s.networkStmtPE.SetEntityUpdateFunc(s.UpdateNSAndPolicyDB)
	s.networkStmtPE.SetIsEntityPresentFunc(s.DoesNSExist)
	actionFuncMap = make(map[int]bgppolicy.PolicyActionFunc)
	actionFuncMap[policyCommonDefs.PolicyActionTypeNetworkStatementAdvertise] = actionFunc
	s.networkStmtPE.SetActionFuncs(actionFuncMap)
	s.networkStmtPE.SetTraverseFuncs(s.TraverseAndApplyNS, s.TraverseAndReverseNS)
	s.policyManager.AddPolicyEngine(s.networkStmtPE)
}

func (s *BGPServer) createListener(proto string) (*net.TCPListener, error) {
	addr := ":" + config.BGPPort
	s.logger.Infof("Listening for incomig connections on %s", addr)
	tcpAddr, err := net.ResolveTCPAddr(proto, addr)
	if err != nil {
		s.logger.Info("ResolveTCPAddr failed with", err)
		return nil, err
	}

	listener, err := net.ListenTCP(proto, tcpAddr)
	if err != nil {
		s.logger.Info("ListenTCP failed with", err)
		return nil, err
	}

	return listener, nil
}

func (s *BGPServer) setListener(listener *net.TCPListener, proto string) {
	switch proto {
	case "tcp4":
		s.listener = listener
	case "tcp6":
		s.listenerIPv6 = listener
	default:
		s.logger.Err("BGPServer:setListener - unknonn protocol type", proto)
	}
}

func (s *BGPServer) listenForPeers(listener *net.TCPListener, proto string, acceptCh chan *net.TCPConn) {
	for {
		s.logger.Info("Waiting for peer connections...")
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			s.logger.Info("AcceptTCP failed with", err)
			if strings.Contains(err.Error(), "use of closed network connection") {
				newListener, err2 := s.createListener(proto)
				if err2 != nil {
					ticker := time.NewTicker(time.Duration(5) * time.Second)
					for range ticker.C {
						newListener, err2 = s.createListener(proto)
						if err2 == nil {
							ticker.Stop()
							break
						}
						s.logger.Err("Create TCPListener for", proto, "failed with err", err)
					}
				}
				s.logger.Info("Created new TCPListener for", proto)
				listener = newListener
				s.setListener(listener, proto)
			}
			continue
		}
		s.logger.Info("Got a peer connection from %s", tcpConn.RemoteAddr())
		acceptCh <- tcpConn
	}
}

func (s *BGPServer) SendUpdate(updated map[uint32]map[*bgprib.Path][]*bgprib.Destination, withdrawn,
	updatedAddPaths []*bgprib.Destination) {
	for _, peer := range s.PeerMap {
		peer.SendUpdate(updated, withdrawn, updatedAddPaths)
	}
}

func (s *BGPServer) DoesRouteExist(params interface{}) bool {
	policyParams := params.(PolicyParams)
	dest := policyParams.dest
	if dest == nil {
		s.logger.Info("BGPServer:DoesRouteExist - dest not found for ip",
			policyParams.route.Dest.BGPRouteState.GetNetwork(), "prefix length",
			policyParams.route.Dest.BGPRouteState.GetCIDRLen())
		return false
	}

	locRibRoute := dest.GetLocRibPathRoute()
	if policyParams.route == locRibRoute {
		return true
	}

	return false
}

func (s *BGPServer) getAggPrefix(conditionsList []interface{}) *packet.IPPrefix {
	s.logger.Info("BGPServer:getAggPrefix")
	var ipPrefix *packet.IPPrefix
	var err error
	for _, condition := range conditionsList {
		switch condition.(type) {
		case utilspolicy.MatchPrefixConditionInfo:
			s.logger.Info("BGPServer:getAggPrefix - PolicyConditionTypeDstIpPrefixMatch case")
			matchPrefix := condition.(utilspolicy.MatchPrefixConditionInfo)
			s.logger.Info("BGPServer:getAggPrefix - exact prefix match conditiontype")
			ipPrefix, err = packet.ConstructIPPrefixFromCIDR(matchPrefix.Prefix.IpPrefix)
			if err != nil {
				s.logger.Info("BGPServer:getAggPrefix - ipPrefix invalid ")
				return nil
			}
			break
		default:
			s.logger.Info("BGPServer:getAggPrefix - Not a known condition type")
			break
		}
	}
	return ipPrefix
}

func (s *BGPServer) setUpdatedAddPaths(policyParams *PolicyParams,
	updatedAddPaths []*bgprib.Destination) {
	if len(updatedAddPaths) > 0 {
		addPathsMap := make(map[*bgprib.Destination]bool)
		for _, dest := range *(policyParams.updatedAddPaths) {
			addPathsMap[dest] = true
		}

		for _, dest := range updatedAddPaths {
			if !addPathsMap[dest] {
				(*policyParams.updatedAddPaths) =
					append((*policyParams.updatedAddPaths), dest)
			}
		}
	}
}

func (s *BGPServer) setWithdrawnWithAggPaths(policyParams *PolicyParams, withdrawn []*bgprib.Destination,
	sendSummaryOnly bool, updatedAddPaths []*bgprib.Destination) {
	destMap := make(map[*bgprib.Destination]bool)
	for _, dest := range *policyParams.withdrawn {
		destMap[dest] = true
	}

	aggDestMap := make(map[*bgprib.Destination]bool)
	for _, aggDestination := range withdrawn {
		aggDestMap[aggDestination] = true
		if !destMap[aggDestination] {
			s.logger.Infof("setWithdrawnWithAggPaths: add agg dest %+v to withdrawn",
				aggDestination.NLRI.GetPrefix())
			(*policyParams.withdrawn) = append((*policyParams.withdrawn), aggDestination)
		}
	}

	// There will be only one destination per aggregated path.
	// So, break out of the loop as soon as we find it.
	for protoFamily, pathDestMap := range *policyParams.updated {
		for path, destinations := range pathDestMap {
			for idx, dest := range destinations {
				if aggDestMap[dest] {
					(*policyParams.updated)[protoFamily][path][idx] = nil
					s.logger.Infof("setWithdrawnWithAggPaths: remove dest %+v from withdrawn",
						dest.NLRI.GetPrefix())
				}
			}
		}
	}

	if sendSummaryOnly {
		if policyParams.DeleteType == utilspolicy.Valid {
			for idx, dest := range *policyParams.withdrawn {
				if dest == policyParams.dest {
					s.logger.Infof("setWithdrawnWithAggPaths: remove dest %+v from withdrawn",
						dest.NLRI.GetPrefix())
					(*policyParams.withdrawn)[idx] = nil
				}
			}
		} else if policyParams.CreateType == utilspolicy.Invalid {
			if policyParams.dest != nil && policyParams.dest.LocRibPath != nil {
				found := false
				protoFamily := policyParams.dest.GetProtocolFamily()
				if destinations, ok :=
					(*policyParams.updated)[protoFamily][policyParams.dest.LocRibPath]; ok {
					for _, dest := range destinations {
						if dest == policyParams.dest {
							found = true
						}
					}
				} else {
					(*policyParams.updated)[protoFamily][policyParams.dest.LocRibPath] = make([]*bgprib.Destination, 0)
				}
				if !found {
					s.logger.Infof("setWithdrawnWithAggPaths: add dest %+v to update",
						policyParams.dest.NLRI.GetPrefix())
					(*policyParams.updated)[protoFamily][policyParams.dest.LocRibPath] = append(
						(*policyParams.updated)[protoFamily][policyParams.dest.LocRibPath], policyParams.dest)
				}
			}
		}
	}

	s.setUpdatedAddPaths(policyParams, updatedAddPaths)
}

func (s *BGPServer) setUpdatedWithAggPaths(policyParams *PolicyParams,
	updated map[uint32]map[*bgprib.Path][]*bgprib.Destination, sendSummaryOnly bool, ipPrefix *packet.IPPrefix,
	protoFamily uint32, updatedAddPaths []*bgprib.Destination) {
	var routeDest *bgprib.Destination
	var ok bool
	if routeDest, ok = s.LocRib.GetDest(ipPrefix, protoFamily, false); !ok {
		s.logger.Err("setUpdatedWithAggPaths: Did not find destination for ip", ipPrefix)
		if policyParams.dest != nil {
			routeDest = policyParams.dest
		} else {
			sendSummaryOnly = false
		}
	}

	withdrawMap := make(map[*bgprib.Destination]bool, len(*policyParams.withdrawn))
	if sendSummaryOnly {
		for _, dest := range *policyParams.withdrawn {
			withdrawMap[dest] = true
		}
	}

	for aggFamily, aggPathDestMap := range updated {
		for aggPath, aggDestinations := range aggPathDestMap {
			destMap := make(map[*bgprib.Destination]bool)
			ppUpdated := *policyParams.updated
			if _, ok := ppUpdated[aggFamily]; !ok {
				ppUpdated[aggFamily] = make(map[*bgprib.Path][]*bgprib.Destination)
			}
			if _, ok := ppUpdated[aggFamily][aggPath]; !ok {
				ppUpdated[aggFamily][aggPath] = make([]*bgprib.Destination, 0)
			} else {
				for _, dest := range ppUpdated[aggFamily][aggPath] {
					destMap[dest] = true
				}
			}

			for _, dest := range aggDestinations {
				if !destMap[dest] {
					s.logger.Infof("setUpdatedWithAggPaths: add agg dest %+v to updated", dest.NLRI.GetPrefix())
					ppUpdated[aggFamily][aggPath] = append(ppUpdated[aggFamily][aggPath], dest)
				}
			}

			if sendSummaryOnly {
				if policyParams.CreateType == utilspolicy.Valid {
					if pathDestMap, ok := ppUpdated[protoFamily]; ok {
						for path, destinations := range pathDestMap {
							for idx, dest := range destinations {
								if routeDest == dest {
									ppUpdated[protoFamily][path][idx] = nil
									s.logger.Infof("setUpdatedWithAggPaths: summaryOnly, remove dest %+v"+
										" from updated", dest.NLRI.GetPrefix())
								}
							}
						}
					}
				} else if policyParams.DeleteType == utilspolicy.Invalid {
					if !withdrawMap[routeDest] {
						s.logger.Infof("setUpdatedWithAggPaths: summaryOnly, add dest %+v to withdrawn",
							routeDest.NLRI.GetPrefix())
						(*policyParams.withdrawn) = append((*policyParams.withdrawn), routeDest)
					}
				}
			}
		}
	}

	s.setUpdatedAddPaths(policyParams, updatedAddPaths)
}

func (s *BGPServer) UndoAggregateAction(actionInfo interface{}, conditionList []interface{},
	policy utilspolicy.Policy, params interface{}, policyStmt utilspolicy.PolicyStmt) {
	policyParams := params.(PolicyParams)
	ipPrefix := packet.NewIPPrefix(net.ParseIP(policyParams.route.Dest.BGPRouteState.GetNetwork()),
		uint8(policyParams.route.Dest.BGPRouteState.GetCIDRLen()))
	protoFamily := policyParams.route.Dest.GetProtocolFamily()
	aggPrefix := s.getAggPrefix(conditionList)
	aggActions := actionInfo.(utilspolicy.PolicyAggregateActionInfo)
	bgpAgg := config.BGPAggregate{
		GenerateASSet:   aggActions.GenerateASSet,
		SendSummaryOnly: aggActions.SendSummaryOnly,
	}

	s.logger.Infof("UndoAggregateAction: ipPrefix=%+v, aggPrefix=%+v", ipPrefix.Prefix, aggPrefix.Prefix)
	var updated map[uint32]map[*bgprib.Path][]*bgprib.Destination
	var withdrawn []*bgprib.Destination
	var updatedAddPaths []*bgprib.Destination
	var origDest *bgprib.Destination
	if policyParams.dest != nil {
		origDest = policyParams.dest
	}
	updated, withdrawn, updatedAddPaths = s.LocRib.RemoveRouteFromAggregate(ipPrefix, aggPrefix,
		s.BgpConfig.Global.Config.RouterId.String(), protoFamily, &bgpAgg, origDest, s.AddPathCount)

	s.logger.Infof("UndoAggregateAction: aggregate result update=%+v, withdrawn=%+v", updated, withdrawn)
	s.setWithdrawnWithAggPaths(&policyParams, withdrawn, aggActions.SendSummaryOnly, updatedAddPaths)
	s.logger.Infof("UndoAggregateAction: after updating withdraw agg paths, update=%+v, withdrawn=%+v,"+
		"policyparams.update=%+v, policyparams.withdrawn=%+v", updated, withdrawn, *policyParams.updated,
		*policyParams.withdrawn)
	return
}

func (s *BGPServer) ApplyAggregateAction(actionInfo interface{}, conditionInfo []interface{},
	policy utilspolicy.Policy, params interface{}, policyStmt utilspolicy.PolicyStmt) {
	policyParams := params.(PolicyParams)
	ipPrefix := packet.NewIPPrefix(net.ParseIP(policyParams.route.Dest.BGPRouteState.GetNetwork()),
		uint8(policyParams.route.Dest.BGPRouteState.GetCIDRLen()))
	protoFamily := policyParams.route.Dest.GetProtocolFamily()
	aggPrefix := s.getAggPrefix(conditionInfo)
	aggActions := actionInfo.(utilspolicy.PolicyAggregateActionInfo)
	bgpAgg := config.BGPAggregate{
		GenerateASSet:   aggActions.GenerateASSet,
		SendSummaryOnly: aggActions.SendSummaryOnly,
	}

	s.logger.Infof("ApplyAggregateAction: ipPrefix=%+v, aggPrefix=%+v", ipPrefix.Prefix, aggPrefix.Prefix)
	var updated map[uint32]map[*bgprib.Path][]*bgprib.Destination
	var withdrawn []*bgprib.Destination
	var updatedAddPaths []*bgprib.Destination
	if (policyParams.CreateType == utilspolicy.Valid) ||
		(policyParams.DeleteType == utilspolicy.Invalid) {
		s.logger.Infof("ApplyAggregateAction: CreateType= Valid or DeleteType = Invalid")
		updated, withdrawn, updatedAddPaths = s.LocRib.AddRouteToAggregate(ipPrefix, aggPrefix,
			s.BgpConfig.Global.Config.RouterId.String(), protoFamily, s.ifaceIP, &bgpAgg, s.AddPathCount)
	} else if policyParams.DeleteType == utilspolicy.Valid {
		s.logger.Infof("ApplyAggregateAction: DeleteType = Valid")
		origDest := policyParams.dest
		updated, withdrawn, updatedAddPaths = s.LocRib.RemoveRouteFromAggregate(ipPrefix, aggPrefix,
			s.BgpConfig.Global.Config.RouterId.String(), protoFamily, &bgpAgg, origDest, s.AddPathCount)
	}

	s.logger.Infof("ApplyAggregateAction: aggregate result update=%+v, withdrawn=%+v", updated, withdrawn)
	s.setUpdatedWithAggPaths(&policyParams, updated, aggActions.SendSummaryOnly, ipPrefix, protoFamily,
		updatedAddPaths)
	s.logger.Infof("ApplyAggregateAction: after updating agg paths, update=%+v, withdrawn=%+v, "+
		"policyparams.update=%+v, policyparams.withdrawn=%+v", updated, withdrawn, *policyParams.updated,
		*policyParams.withdrawn)
	return
}

func (s *BGPServer) CheckForAggregation(updated map[uint32]map[*bgprib.Path][]*bgprib.Destination, withdrawn,
	updatedAddPaths []*bgprib.Destination) (map[uint32]map[*bgprib.Path][]*bgprib.Destination, []*bgprib.Destination,
	[]*bgprib.Destination) {
	s.logger.Infof("BGPServer:checkForAggregate - start, updated %v withdrawn %v", updated, withdrawn)

	for _, dest := range withdrawn {
		if dest == nil || dest.LocRibPath == nil || dest.LocRibPath.IsAggregate() {
			continue
		}

		protoFamily := dest.GetProtocolFamily()
		if len(s.BgpConfig.Afs[protoFamily].BgpAggs) == 0 {
			s.logger.Crit("BGPServer:checkForAggregate - withdrawn, aggs for family", protoFamily, "is not set")
			continue
		}

		pe, ok := s.locRibPE[protoFamily]
		if !ok {
			s.logger.Err("BGPServer:checkForAggregate - Agg policy engine not found for family", protoFamily)
			continue
		}

		route := dest.GetLocRibPathRoute()
		if route == nil {
			s.logger.Infof("BGPServer:checkForAggregate - route not found withdraw dest %s",
				dest.NLRI.GetCIDR())
			continue
		}
		peEntity := utilspolicy.PolicyEngineFilterEntityParams{
			DestNetIp:  route.Dest.BGPRouteState.GetNetwork() + "/" + strconv.Itoa(int(route.Dest.BGPRouteState.GetCIDRLen())),
			NextHopIp:  route.PathInfo.NextHop,
			DeletePath: true,
		}
		s.logger.Infof("BGPServer:checkForAggregate - withdraw dest %s policylist %v hit %v before ",
			"applying delete policy", dest.NLRI.GetCIDR(), route.PolicyList, route.PolicyHitCounter)
		callbackInfo := PolicyParams{
			CreateType:      utilspolicy.Invalid,
			DeleteType:      utilspolicy.Valid,
			route:           route,
			dest:            dest,
			updated:         &updated,
			withdrawn:       &withdrawn,
			updatedAddPaths: &updatedAddPaths,
		}
		pe.PolicyEngine.PolicyEngineFilter(peEntity, policyCommonDefs.PolicyPath_Export, callbackInfo)
	}

	for protoFamily, pathDestMap := range updated {
		if len(s.BgpConfig.Afs[protoFamily].BgpAggs) == 0 {
			s.logger.Crit("BGPServer:checkForAggregate - updated, aggs for family", protoFamily, "is not set")
			continue
		}

		pe, ok := s.locRibPE[protoFamily]
		if !ok {
			s.logger.Err("BGPServer:checkForAggregate - Agg policy engine not found for family", protoFamily)
			continue
		}

		for _, destinations := range pathDestMap {
			s.logger.Infof("BGPServer:checkForAggregate - update destinations %+v", destinations)
			for _, dest := range destinations {
				if dest == nil || dest.LocRibPath == nil || dest.LocRibPath.IsAggregate() {
					continue
				}
				route := dest.GetLocRibPathRoute()
				s.logger.Infof("BGPServer:checkForAggregate - update dest %s policylist %v hit %v before "+
					"applying create policy", dest.NLRI.GetCIDR(), route.PolicyList, route.PolicyHitCounter)
				if route != nil {
					peEntity := utilspolicy.PolicyEngineFilterEntityParams{
						DestNetIp:  route.Dest.BGPRouteState.GetNetwork() + "/" + strconv.Itoa(int(route.Dest.BGPRouteState.GetCIDRLen())),
						NextHopIp:  route.PathInfo.NextHop,
						CreatePath: true,
					}
					callbackInfo := PolicyParams{
						CreateType:      utilspolicy.Valid,
						DeleteType:      utilspolicy.Invalid,
						route:           route,
						dest:            dest,
						updated:         &updated,
						withdrawn:       &withdrawn,
						updatedAddPaths: &updatedAddPaths,
					}
					pe.PolicyEngine.PolicyEngineFilter(peEntity, policyCommonDefs.PolicyPath_Export, callbackInfo)
					s.logger.Infof("BGPServer:checkForAggregate - update dest %s policylist %v hit %v "+
						"after applying create policy", dest.NLRI.GetCIDR(), route.PolicyList,
						route.PolicyHitCounter)
				}
			}
		}
	}

	s.logger.Infof("BGPServer:checkForAggregate - complete, updated %v withdrawn %v", updated, withdrawn)
	return updated, withdrawn, updatedAddPaths
}

func (s *BGPServer) UpdateRouteAndPolicyDB(policyDetails utilspolicy.PolicyDetails, params interface{}) {
	var op int
	policyParams := params.(PolicyParams)
	dest := policyParams.dest
	protoFamily := dest.GetProtocolFamily()
	pe, ok := s.locRibPE[protoFamily]
	if !ok {
		s.logger.Err("UpdateRouteAndPolicyDB - Agg policy engine not found for family", protoFamily)
		return
	}

	if policyParams.DeleteType != bgppolicy.Invalid {
		op = bgppolicy.Del
	} else {
		if policyDetails.EntityDeleted == false {
			s.logger.Info("Reject action was not applied, so add this policy to the route")
			op = bgppolicy.Add
			bgppolicy.UpdateRoutePolicyState(policyParams.route, op, policyDetails.Policy, policyDetails.PolicyStmt)
		}
		policyParams.route.PolicyHitCounter++
	}
	pe.UpdatePolicyRouteMap(policyParams.route, policyDetails.Policy, op)
}

func (s *BGPServer) TraverseAndApplyBGPRib(data interface{}, updateFunc utilspolicy.PolicyApplyfunc) {
	s.logger.Infof("BGPServer:TraverseAndApplyBGPRib - start")
	updated := make(map[uint32]map[*bgprib.Path][]*bgprib.Destination, 10)
	withdrawn := make([]*bgprib.Destination, 0, 10)
	updatedAddPaths := make([]*bgprib.Destination, 0)
	locRib := s.LocRib.GetLocRib()
	for protoFamily, pathDestMap := range locRib {
		if _, ok := s.locRibPE[protoFamily]; !ok {
			s.logger.Err("TraverseAndApplyBGPRib - Agg policy engine not found for family", protoFamily)
			continue
		}
		for path, destinations := range pathDestMap {
			for _, dest := range destinations {
				if !path.IsAggregatePath() {
					route := dest.GetLocRibPathRoute()
					if route == nil {
						continue
					}
					peEntity := utilspolicy.PolicyEngineFilterEntityParams{
						DestNetIp: route.Dest.BGPRouteState.GetNetwork() + "/" +
							strconv.Itoa(int(route.Dest.BGPRouteState.GetCIDRLen())),
						NextHopIp:  route.PathInfo.NextHop,
						PolicyList: route.PolicyList,
					}
					callbackInfo := PolicyParams{
						CreateType:      utilspolicy.Invalid,
						DeleteType:      utilspolicy.Invalid,
						route:           route,
						dest:            dest,
						updated:         &updated,
						withdrawn:       &withdrawn,
						updatedAddPaths: &updatedAddPaths,
					}

					updateFunc(peEntity, data, callbackInfo)
				}
			}
		}
	}
	s.logger.Infof("BGPServer:TraverseAndApplyBGPRib - updated %v withdrawn %v", updated, withdrawn)
	s.SendUpdate(updated, withdrawn, updatedAddPaths)
}

func (s *BGPServer) TrAndRevAggForIPv4(policyData interface{}) {
	protoFamily := packet.GetProtocolFamily(packet.AfiIP, packet.SafiUnicast)
	pe, ok := s.locRibPE[protoFamily]
	if !ok {
		s.logger.Err("TrAndRevAggForIPv4 - Agg policy engine not found for family", protoFamily)
		return
	}
	s.TraverseAndReverseBGPRib(policyData, pe)
}

func (s *BGPServer) TrAndRevAggForIPv6(policyData interface{}) {
	protoFamily := packet.GetProtocolFamily(packet.AfiIP6, packet.SafiUnicast)
	pe, ok := s.locRibPE[protoFamily]
	if !ok {
		s.logger.Err("TrAndRevAggForIPv6 - Agg policy engine not found for family", protoFamily)
		return
	}
	s.TraverseAndReverseBGPRib(policyData, pe)
}

func (s *BGPServer) TraverseAndReverseBGPRib(policyData interface{}, pe *bgppolicy.LocRibPolicyEngine) {
	updateInfo := policyData.(utilspolicy.PolicyEngineApplyInfo)
	applyPolicyInfo := updateInfo.ApplyPolicy
	policy := applyPolicyInfo.ApplyPolicy
	//	policy := policyData.(utilspolicy.Policy)
	s.logger.Info("BGPServer:TraverseAndReverseBGPRib - policy", policy.Name)
	policyExtensions := policy.Extensions.(bgppolicy.PolicyExtensions)
	if len(policyExtensions.RouteList) == 0 {
		s.logger.Info("No route affected by this policy, so nothing to do")
		return
	}

	updated := make(map[uint32]map[*bgprib.Path][]*bgprib.Destination, 10)
	withdrawn := make([]*bgprib.Destination, 0, 10)
	updatedAddPaths := make([]*bgprib.Destination, 0)
	var route *bgprib.Route
	for idx := 0; idx < len(policyExtensions.RouteInfoList); idx++ {
		route = policyExtensions.RouteInfoList[idx]
		dest := s.LocRib.GetDestFromIPAndLen(route.Dest.GetProtocolFamily(), route.Dest.BGPRouteState.GetNetwork(),
			uint32(route.Dest.BGPRouteState.GetCIDRLen()))

		callbackInfo := PolicyParams{
			CreateType:      utilspolicy.Invalid,
			DeleteType:      utilspolicy.Invalid,
			route:           route,
			dest:            dest,
			updated:         &updated,
			withdrawn:       &withdrawn,
			updatedAddPaths: &updatedAddPaths,
		}
		peEntity := utilspolicy.PolicyEngineFilterEntityParams{
			DestNetIp: route.Dest.BGPRouteState.GetNetwork() + "/" + strconv.Itoa(int(route.Dest.BGPRouteState.GetCIDRLen())),
			NextHopIp: route.PathInfo.NextHop,
		}

		ipPrefix, err := bgppolicy.GetNetworkPrefixFromCIDR(route.Dest.BGPRouteState.GetNetwork() + "/" +
			strconv.Itoa(int(route.Dest.BGPRouteState.GetCIDRLen())))
		if err != nil {
			s.logger.Info("Invalid route ", ipPrefix)
			continue
		}
		pe.PolicyEngine.PolicyEngineUndoApplyPolicyForEntity(peEntity, updateInfo, callbackInfo)
		pe.DeleteRoutePolicyState(route, policy.Name)
		pe.PolicyEngine.DeletePolicyEntityMapEntry(peEntity, policy.Name)
	}
}

func (s *BGPServer) DoesAdjRIBRouteExist(params interface{}, adjRIBDir bgprib.AdjRIBDir) bool {
	policyParams := params.(*AdjRIBPolicyParams)
	peer := policyParams.Peer
	if peer == nil {
		s.logger.Info("BGPServer:DoesAdjRIBRouteExist - Neighbor not set for route, nlri =",
			policyParams.Route.NLRI.String())
		return false
	}

	adjRIB := peer.GetAdjRIB(adjRIBDir)
	if prefixRouteMap, ok := adjRIB[policyParams.Route.ProtocolFamily]; ok {
		if prefixRouteMap[policyParams.Route.NLRI.String()] != nil {
			return true
		}
	}

	return false
}

func (s *BGPServer) DoesAdjRIBInRouteExist(params interface{}) bool {
	return s.DoesAdjRIBRouteExist(params, bgprib.AdjRIBDirIn)
}

func (s *BGPServer) DoesAdjRIBOutRouteExist(params interface{}) bool {
	return s.DoesAdjRIBRouteExist(params, bgprib.AdjRIBDirOut)
}

func (s *BGPServer) ApplyAdjRIBAction(actionInfo interface{}, conditionInfo []interface{}, policy utilspolicy.Policy,
	params interface{}, policyStmt utilspolicy.PolicyStmt) {
	policyParams := params.(*AdjRIBPolicyParams)
	policyParams.PolicyStmt = policyStmt
	s.logger.Infof("BGPServer:ApplyAdjRIBAction - actionInfo=%+v, conditionInfo=%+v, policyParams=%+v, policyStmt=%+v\n",
		actionInfo, conditionInfo, policyParams, policyStmt)
	if len(policyStmt.Actions) > 0 {
		for _, action := range policyStmt.Actions {
			if action == "permit" {
				s.logger.Infof("BGPServer:ApplyAdjRIBAction - policyParams=%+v, policyStmt=%+v, action permit\n",
					policyParams, policyStmt, action)
				policyParams.Accept = Accept
				break
			} else if action == "deny" {
				s.logger.Infof("BGPServer:ApplyAdjRIBAction - policyParams=%+v, policyStmt=%+v, action deny\n",
					policyParams, policyStmt, action)
				policyParams.Accept = Reject
			} else {
				s.logger.Errf("BGPServer:ApplyAdjRIBAction - policyParams=%+v, policyStmt=%+v, unknown action=%s\n",
					policyParams, policyStmt, action)
			}
		}
	}
}

func (s *BGPServer) UndoAdjRIBAction(actionInfo interface{}, conditionInfo []interface{}, policy utilspolicy.Policy,
	params interface{}, policyStmt utilspolicy.PolicyStmt) {
	policyParams := params.(*AdjRIBPolicyParams)
	s.logger.Infof("BGPServer:UndoAdjRIBAction - actionInfo=%+v, conditionInfo=%+v, policyParams=%+v policyStmt=%+v\n",
		actionInfo, conditionInfo, policyParams, policyStmt)
	if len(policyStmt.Actions) > 0 {
		for _, action := range policyStmt.Actions {
			if action == "permit" {
				s.logger.Infof("BGPServer:UndoAdjRIBAction - policyParams=%+v, policyStmt=%+v, action permit\n",
					policyParams, policyStmt, action)
				policyParams.Accept = Accept
				break
			} else if action == "deny" {
				s.logger.Infof("BGPServer:UndoAdjRIBAction - policyParams=%+v, policyStmt=%+v, action deny\n",
					policyParams, policyStmt, action)
				policyParams.Accept = Reject
			} else {
				s.logger.Errf("BGPServer:UndoAdjRIBAction - policyParams=%+v, policyStmt=%+v, unknown action=%s\n",
					policyParams, policyStmt, action)
			}
		}
	}
}

func (s *BGPServer) UpdateAdjRIBRouteAndPolicyDB(policyDetails utilspolicy.PolicyDetails, params interface{},
	pe *bgppolicy.AdjRibPPolicyEngine) {
	var op int
	policyParams := params.(*AdjRIBPolicyParams)
	s.logger.Infof("UpdateAdjRIBRouteAndPolicyDB - policyDetails=%+v, route=%+v", policyDetails, policyParams.Route)

	if policyParams.DeleteType != bgppolicy.Invalid {
		s.logger.Infof("UpdateAdjRIBRouteAndPolicyDB - Route deleted")
		op = bgppolicy.Del
	} else {
		if policyDetails.EntityDeleted == false {
			s.logger.Info("Reject action was not applied, so add this policy to the route")
			op = bgppolicy.Add
			bgppolicy.UpdateAdjRIBRoutePolicyState(policyParams.Route, op, policyDetails.Policy,
				policyDetails.PolicyStmt)
		}
		policyParams.Route.PolicyHitCounter++
	}
	pe.UpdateAdjRIBPolicyRouteMap(policyParams.Route, policyDetails.Policy, op)
}

func (s *BGPServer) UpdateAdjRIBInRouteAndPolicyDB(policyDetails utilspolicy.PolicyDetails, params interface{}) {
	s.UpdateAdjRIBRouteAndPolicyDB(policyDetails, params, s.ribInPE)
}

func (s *BGPServer) UpdateAdjRIBOutRouteAndPolicyDB(policyDetails utilspolicy.PolicyDetails, params interface{}) {
	s.UpdateAdjRIBRouteAndPolicyDB(policyDetails, params, s.ribOutPE)
}

func (s *BGPServer) getPeerForPolicy(data interface{}, updateFunc utilspolicy.PolicyApplyfunc,
	pe *bgppolicy.AdjRibPPolicyEngine) *Peer {
	s.logger.Infof("BGPServer:TraverseAndApplyAdjRib - start")
	policyInfo := data.(utilspolicy.PolicyEngineApplyInfo)
	conditionsDB := pe.PolicyEngine.PolicyConditionsDB
	var neighborIP string
	var peer *Peer
	var ok bool

	for _, condition := range policyInfo.ApplyPolicy.Conditions {
		s.logger.Infof("BGPServer:TraverseAndApplyAdjRib - condition:%+v", condition)
		nodeGet := conditionsDB.Get(patriciaDB.Prefix(condition))
		if nodeGet == nil {
			s.logger.Err("Condition", condition, "not defined")
			return nil
		}
		node := nodeGet.(utilspolicy.PolicyCondition)
		if node.ConditionType == policyCommonDefs.PolicyConditionTypeNeighborMatch {
			neighborIP = node.ConditionInfo.(string)
		}
	}

	if peer, ok = s.PeerMap[neighborIP]; !ok {
		s.logger.Err("Can't apply policy... Neighbor %s not found", neighborIP)
		return nil
	}

	return peer
	/*
		adjRIB := s.GetAdjRIB(peer, adjRibDir)
		for _, prefixRouteMap := range adjRIB {
			for _, adjRoute := range prefixRouteMap {
				if adjRoute == nil {
					continue
				}

				s.logger.Debugf("Peer %s - NLRI %s policylist %v hit %v before applying create policy",
					adjRoute.NLRI.GetPrefix(), adjRoute.PolicyList, adjRoute.PolicyHitCounter)
				peEntity := utilspolicy.PolicyEngineFilterEntityParams{
					DestNetIp:  adjRoute.NLRI.GetCIDR(),
					Neighbor:   peer.NeighborConf.RunningConf.NeighborAddress.String(),
					PolicyList: adjRoute.PolicyList,
				}
				callbackInfo := &AdjRIBPolicyParams{
					CreateType: utilspolicy.Invalid,
					DeleteType: utilspolicy.Invalid,
					Peer:       peer,
					Route:      adjRoute,
				}

				updateFunc(peEntity, policyInfo, callbackInfo)
			}
		}
	*/
}

func (s *BGPServer) TraverseAndApplyAdjRibIn(data interface{}, updateFunc utilspolicy.PolicyApplyfunc) {
	peer := s.getPeerForPolicy(data, updateFunc, s.ribInPE)
	if peer == nil {
		s.logger.Infof("BGPServer:TraverseAndApplyAdjRibIn - peer not found")
		return
	}
	updated, withdrawn, updatedAddPaths := peer.AdjRIBInPolicyUpdated(bgprib.AdjRIBDirIn, data, updateFunc)
	updated, withdrawn, updatedAddPaths = s.CheckForAggregation(updated, withdrawn, updatedAddPaths)
	s.SendUpdate(updated, withdrawn, updatedAddPaths)
}

func (s *BGPServer) TraverseAndApplyAdjRibOut(data interface{}, updateFunc utilspolicy.PolicyApplyfunc) {
	peer := s.getPeerForPolicy(data, updateFunc, s.ribOutPE)
	if peer == nil {
		s.logger.Infof("BGPServer:TraverseAndApplyAdjRibOut - peer not found")
		return
	}

	peer.AdjRIBOutPolicyUpdated(data, updateFunc)
}

func (s *BGPServer) TraverseAndReverseAdjRIB(policyData interface{}, pe *bgppolicy.AdjRibPPolicyEngine) {
	updateInfo := policyData.(utilspolicy.PolicyEngineApplyInfo)
	applyPolicyInfo := updateInfo.ApplyPolicy //policyData.(utilspolicy.ApplyPolicyInfo)
	policy := applyPolicyInfo.ApplyPolicy     //policyItem.(policy.Policy)
	s.logger.Info("BGPServer:TraverseAndReverseAdjRIB - policy", policy.Name)
	policyExtensions := policy.Extensions.(bgppolicy.AdjRibPolicyExtensions)
	if len(policyExtensions.RouteList) == 0 {
		s.logger.Info("No route affected by this policy, so nothing to do")
		return
	}

	var route *bgprib.AdjRIBPathIdRoute
	var peer *Peer
	for idx := 0; idx < len(policyExtensions.RouteInfoList); idx++ {
		route = policyExtensions.RouteInfoList[idx]
		peerIP := route.Neighbor.String()
		if peer == nil {
			var ok bool
			if peer, ok = s.PeerMap[peerIP]; !ok {
				s.logger.Err("Peer not found for ip", peerIP, "for NLRI", route.NLRI.GetIPPrefix().String())
				continue
			}
		}

		callbackInfo := &AdjRIBPolicyParams{
			CreateType: utilspolicy.Invalid,
			DeleteType: utilspolicy.Invalid,
			Route:      route,
			Peer:       peer,
		}
		peEntity := utilspolicy.PolicyEngineFilterEntityParams{
			DestNetIp: route.NLRI.GetCIDR(),
			Neighbor:  route.Neighbor.String(),
		}

		success := pe.PolicyEngine.PolicyEngineUndoApplyPolicyForEntity(peEntity, updateInfo, callbackInfo)
		s.logger.Info("success value agyer undoapplypolicy:", success, " for policy:", policy.Name)
		if success {
			pe.AdjRIBDeleteRoutePolicyState(route, policy.Name)
			pe.PolicyEngine.DeletePolicyEntityMapEntry(peEntity, policy.Name)
		}
	}
}

func (s *BGPServer) TraverseAndReverseAdjRIBIn(policyData interface{}) {
	s.TraverseAndReverseAdjRIB(policyData, s.ribInPE)
}

func (s *BGPServer) TraverseAndReverseAdjRIBOut(policyData interface{}) {
	s.TraverseAndReverseAdjRIB(policyData, s.ribOutPE)
}

func (s *BGPServer) ProcessUpdate(pktInfo *packet.BGPPktSrc) {
	peer, ok := s.PeerMap[pktInfo.Src]
	if !ok {
		s.logger.Err("BgpServer:ProcessUpdate - Peer not found, address:", pktInfo.Src)
		return
	}

	updated, withdrawn, updatedAddPaths := peer.ReceiveUpdate(pktInfo)
	updated, withdrawn, updatedAddPaths = s.CheckForAggregation(updated, withdrawn, updatedAddPaths)
	s.SendUpdate(updated, withdrawn, updatedAddPaths)
}

func (s *BGPServer) convertDestIPToIPPrefix(routes []*config.RouteInfo) map[uint32][]packet.NLRI {
	pfNLRI := make(map[uint32][]packet.NLRI)
	var protoFamily uint32
	for _, r := range routes {
		ip := net.ParseIP(r.IPAddr)
		if ip == nil {
			s.logger.Errf("Connected route %s/%s is not a valid IP", r.IPAddr, r.Mask)
			continue
		}

		if ip.To4() != nil {
			s.logger.Info("convertDestIPToIPPrefix:ipv4")
			protoFamily = packet.GetProtocolFamily(packet.AfiIP, packet.SafiUnicast)
		} else {
			s.logger.Info("convertDestIPToIPPrefix:ipv6")
			protoFamily = packet.GetProtocolFamily(packet.AfiIP6, packet.SafiUnicast)
		}

		s.logger.Infof("Connected route: addr %s netmask %s, protoFamily:", r.IPAddr, r.Mask, protoFamily)
		if _, ok := pfNLRI[protoFamily]; !ok {
			pfNLRI[protoFamily] = make([]packet.NLRI, 0)
		}

		ipPrefix := packet.ConstructIPPrefix(r.IPAddr, r.Mask)
		pfNLRI[protoFamily] = append(pfNLRI[protoFamily], ipPrefix)
	}
	return pfNLRI
}

func (s *BGPServer) ProcessConnectedRoutes(installedRoutes, withdrawnRoutes []*config.RouteInfo) {
	s.logger.Info("valid routes:", installedRoutes, "invalid routes:", withdrawnRoutes)
	valid := s.convertDestIPToIPPrefix(installedRoutes)
	invalid := s.convertDestIPToIPPrefix(withdrawnRoutes)
	s.logger.Info("pfNLRI valid:", valid, "invalid:", invalid)
	routerId := s.BgpConfig.Global.Config.RouterId.String()
	updated, withdrawn, updatedAddPaths := s.LocRib.ProcessConnectedRoutes(routerId, s.ConnRoutesPath, valid,
		invalid, s.AddPathCount)
	updated, withdrawn, updatedAddPaths = s.CheckForAggregation(updated, withdrawn, updatedAddPaths)
	s.SendUpdate(updated, withdrawn, updatedAddPaths)
}

func (s *BGPServer) ProcessDefaultRoutes(add, remove map[uint32][]packet.NLRI) {
	s.logger.Info("pfNLRI add:", add, "remove:", remove)
	routerId := s.BgpConfig.Global.Config.RouterId.String()
	updated, withdrawn, updatedAddPaths := s.LocRib.ProcessConnectedRoutes(routerId, s.DefaultRoutesPath, add,
		remove, s.AddPathCount)
	updated, withdrawn, updatedAddPaths = s.CheckForAggregation(updated, withdrawn, updatedAddPaths)
	s.SendUpdate(updated, withdrawn, updatedAddPaths)
}

func (s *BGPServer) UpdateDefaultRoutes(update map[uint32]bool) {
	s.logger.Debug("UpdateDefaultRoutes - update", update)
	add, remove := packet.ConstructNLRIForDefaultRoutes(update)
	if len(add) > 0 || len(remove) > 0 {
		s.ProcessDefaultRoutes(add, remove)
	}
}

func (s *BGPServer) ProcessIntfStates(intfs []*config.IntfStateInfo) {
	for _, ifState := range intfs {
		if ifState.State == config.INTF_CREATED {
			s.ifaceMgr.AddIface(ifState.Idx, ifState.IPAddr)
		} else if ifState.State == config.INTF_DELETED {
			s.ifaceMgr.RemoveIface(ifState.Idx, ifState.IPAddr)
		} else if ifState.State == config.INTFV6_CREATED {
			s.ifaceMgr.AddV6Iface(ifState.Idx, ifState.IPAddr)
		} else if ifState.State == config.INTFV6_DELETED {
			s.ifaceMgr.RemoveV6Iface(ifState.Idx, ifState.IPAddr)
		} else if ifState.State == config.IPV6_NEIGHBOR_CREATED {
			s.ifaceMgr.AddLinkLocalIface(ifState.Idx, ifState.LinkLocalIP)
		} else if ifState.State == config.IPV6_NEIGHBOR_DELETED {
			s.ifaceMgr.RemoveLinkLocalIface(ifState.Idx, ifState.LinkLocalIP)
		}
	}
}

func (s *BGPServer) GetIfaceIP(ifIndex int32) (*utils.IPInfo, error) {
	ipInfo, err := s.ifaceMgr.GetIfaceIP(ifIndex)
	return ipInfo, err
}

func (s *BGPServer) ProcessRemoveNeighbor(peerIp string, peer *Peer) {
	updated, withdrawn, updatedAddPaths := s.LocRib.RemoveUpdatesFromNeighbor(peerIp, peer.NeighborConf,
		s.AddPathCount)
	s.logger.Infof("ProcessRemoveNeighbor - Neighbor %s, send updated paths %v, withdrawn paths %v",
		peerIp, updated, withdrawn)
	updated, withdrawn, updatedAddPaths = s.CheckForAggregation(updated, withdrawn, updatedAddPaths)
	s.SendUpdate(updated, withdrawn, updatedAddPaths)
}

func (s *BGPServer) SendAllRoutesToPeer(peer *Peer) {
	withdrawn := make([]*bgprib.Destination, 0)
	updatedAddPaths := make([]*bgprib.Destination, 0)
	updated := s.LocRib.GetLocRib()
	s.SendUpdate(updated, withdrawn, updatedAddPaths)
}

func (s *BGPServer) RemoveRoutesFromAllNeighbor() {
	s.LocRib.RemoveUpdatesFromAllNeighbors(s.AddPathCount)
}

func (s *BGPServer) addPeerToList(peer *Peer) {
	s.Neighbors = append(s.Neighbors, peer)
}

func (s *BGPServer) removePeerFromList(peer *Peer) {
	for idx, item := range s.Neighbors {
		if item == peer {
			s.Neighbors[idx] = s.Neighbors[len(s.Neighbors)-1]
			s.Neighbors[len(s.Neighbors)-1] = nil
			s.Neighbors = s.Neighbors[:len(s.Neighbors)-1]
			break
		}
	}
}

func (s *BGPServer) StopPeersByGroup(groupName string, peerAddrType config.PeerAddressType) []*Peer {
	peers := make([]*Peer, 0)
	for peerIP, peer := range s.PeerMap {
		if peer.NeighborConf.Group != nil && peer.NeighborConf.RunningConf.PeerAddressType == peerAddrType &&
			peer.NeighborConf.Group.Name == groupName {
			s.logger.Info("Clean up peer", peerIP)
			peer.Cleanup()
			s.ProcessRemoveNeighbor(peerIP, peer)
			peers = append(peers, peer)

			runtime.Gosched()
		}
	}

	return peers
}

func (s *BGPServer) UpdatePeerGroupInPeers(groupName string, peerAddrType config.PeerAddressType,
	peerGroup *config.PeerGroupConfig) {
	peers := s.StopPeersByGroup(groupName, peerAddrType)
	for _, peer := range peers {
		peer.UpdatePeerGroup(peerGroup)
		peer.Init()
	}
}

func (s *BGPServer) DeleteAgg(aggConf config.BGPAggregate) error {
	pe, ok := s.locRibPE[aggConf.AddressFamily]
	if ok {
		policyEngine := pe.GetPolicyEngine()
		policyDB := policyEngine.PolicyDB

		nodeGet := policyDB.Get(patriciaDB.Prefix(aggConf.IPPrefix))
		if nodeGet == nil {
			s.logger.Err("Policy ", aggConf, " not created yet")
			return errors.New(fmt.Sprintf("Policy %s not found in policy engine", aggConf.IPPrefix))
		}
		node := nodeGet.(utilspolicy.Policy)

		pe.ReleaseId(uint32(node.Precedence))
		pe.DeletePolicyDefinition(aggConf.IPPrefix)
		pe.DeletePolicyStmt(aggConf.IPPrefix)
		pe.DeletePolicyCondition(aggConf.IPPrefix)
	}
	if _, ok := s.BgpConfig.Afs[aggConf.AddressFamily]; ok {
		delete(s.BgpConfig.Afs[aggConf.AddressFamily].BgpAggs, aggConf.IPPrefix)
	}
	return nil
}

func (s *BGPServer) AddOrUpdateAgg(oldConf config.BGPAggregate, newConf config.BGPAggregate, attrSet []bool) error {
	s.logger.Info("AddOrUpdateAgg")
	var err error

	pe, ok := s.locRibPE[newConf.AddressFamily]
	if !ok {
		s.logger.Err("Aggregate policy engine not created for address family", newConf.AddressFamily)
		return errors.New(fmt.Sprintf("Aggregate policy engine not created for address family", newConf.AddressFamily))
	}

	bytes := packet.GetAddressLengthForFamily(newConf.AddressFamily)
	if bytes == -1 {
		s.logger.Err("Could not find number of bytes for aggregate prefix family", newConf.AddressFamily)
		return errors.New(fmt.Sprintf("Could not find number of bytes for aggregate prefix family",
			newConf.AddressFamily))
	}
	strBits := strconv.Itoa(bytes * 8)
	s.logger.Infof("AddOrUpdateAgg: agg = %s, bytes = %d, bits =%s", newConf.IPPrefix, bytes, strBits)

	if oldConf.IPPrefix != "" {
		// Delete the policy
		s.DeleteAgg(oldConf)
	}

	if newConf.IPPrefix != "" {
		// Create the policy
		name := newConf.IPPrefix
		tokens := strings.Split(newConf.IPPrefix, "/")
		prefixLen := tokens[1]
		_, err := strconv.Atoi(prefixLen)
		if err != nil {
			s.logger.Errf("Failed to convert prefex len %s to int with error %s", prefixLen, err)
			return err
		}

		cond := utilspolicy.PolicyConditionConfig{
			Name:          name,
			ConditionType: "MatchDstIpPrefix",
			MatchDstIpPrefixConditionInfo: utilspolicy.PolicyDstIpMatchPrefixSetCondition{
				Prefix: utilspolicy.PolicyPrefix{
					IpPrefix:        newConf.IPPrefix,
					MasklengthRange: prefixLen + "-" + strBits,
				},
			},
		}

		_, err = pe.CreatePolicyCondition(cond)
		if err != nil {
			s.logger.Errf("Failed to create policy condition for aggregate %s with error %s", name, err)
			return err
		}

		stmt := utilspolicy.PolicyStmtConfig{Name: name, MatchConditions: "all"}
		stmt.Conditions = make([]string, 1)
		stmt.Conditions[0] = name
		stmt.Actions = make([]string, 1)
		stmt.Actions[0] = "permit"
		err = pe.CreatePolicyStmt(stmt)
		if err != nil {
			s.logger.Errf("Failed to create policy statement for aggregate %s with error %s", name, err)
			pe.DeletePolicyCondition(name)
			return err
		}

		precedence := pe.GetNextId()
		def := utilspolicy.PolicyDefinitionConfig{Name: name, Precedence: int(precedence), MatchType: "any", PolicyType: "BGP"}
		def.PolicyDefinitionStatements = make([]utilspolicy.PolicyDefinitionStmtPrecedence, 1)
		policyDefinitionStatement := utilspolicy.PolicyDefinitionStmtPrecedence{
			Precedence: 1,
			Statement:  name,
		}
		def.PolicyDefinitionStatements[0] = policyDefinitionStatement
		def.Extensions = bgppolicy.PolicyExtensions{}
		err = pe.CreatePolicyDefinition(def)
		if err != nil {
			s.logger.Errf("Failed to create policy definition for aggregate %s with error %s", name, err)
			pe.ReleaseId(precedence)
			pe.DeletePolicyStmt(name)
			pe.DeletePolicyCondition(name)
			return err
		}

		err = s.UpdateAggPolicy(name, pe, newConf)
		if _, ok := s.BgpConfig.Afs[newConf.AddressFamily]; !ok {
			s.BgpConfig.Afs[newConf.AddressFamily] = &config.AddressFamily{}
		}
		if s.BgpConfig.Afs[newConf.AddressFamily].BgpAggs == nil {
			s.BgpConfig.Afs[newConf.AddressFamily].BgpAggs = make(map[string]*config.BGPAggregate)
		}
		s.BgpConfig.Afs[newConf.AddressFamily].BgpAggs[newConf.IPPrefix] = &newConf
		return err
	}
	return err
}

func (s *BGPServer) UpdateAggPolicy(policyName string, pe *bgppolicy.LocRibPolicyEngine,
	aggConf config.BGPAggregate) error {
	s.logger.Debug("UpdateApplyPolicy")
	var err error
	var policyAction utilspolicy.PolicyAction
	conditionNameList := make([]string, 0)

	policyEngine := pe.GetPolicyEngine()
	policyDB := policyEngine.PolicyDB

	nodeGet := policyDB.Get(patriciaDB.Prefix(policyName))
	if nodeGet == nil {
		s.logger.Err("Policy ", policyName, " not defined")
		return errors.New(fmt.Sprintf("Policy %s not found in policy engine", policyName))
	}
	node := nodeGet.(utilspolicy.Policy)

	aggregateActionInfo := utilspolicy.PolicyAggregateActionInfo{aggConf.GenerateASSet, aggConf.SendSummaryOnly}
	policyAction = utilspolicy.PolicyAction{
		Name:       aggConf.IPPrefix,
		ActionType: policyCommonDefs.PolicyActionTypeAggregate,
		ActionInfo: aggregateActionInfo,
	}

	s.logger.Debug("Calling applypolicy with conditionNameList: ", conditionNameList)
	pe.UpdateApplyPolicy(utilspolicy.ApplyPolicyInfo{node, policyAction, conditionNameList}, true)
	return err
}

func (s *BGPServer) UndoNSAction(actionInfo interface{}, conditionList []interface{}, policy utilspolicy.Policy,
	params interface{}, policyStmt utilspolicy.PolicyStmt) {
	policyParams := params.(*NSPolicyParams)
	s.logger.Infof("UndoNSAction: ipPrefix=%s, policy=%s, statement=%s", policyParams.NS.IPPrefix, policy.Name,
		policyStmt.Name)
	s.logger.Infof("UndoNSAction - actionInfo=%+v, conditionInfo=%+v, policyParams=%+v, policyStmt=%+v\n",
		actionInfo, conditionList, policyParams, policyStmt)
	if len(policyStmt.Actions) > 0 {
		for _, action := range policyStmt.Actions {
			if action == "permit" {
				s.logger.Infof("UndoNSAction - policyParams=%+v, policyStmt=%+v, action permit\n",
					policyParams, policyStmt)
				policyParams.Accept = Accept
				if policyParams.RemovePolicyStmt == nil {
					policyParams.RemovePolicyStmt = make([]utilspolicy.PolicyStmt, 0)
				}
				policyParams.RemovePolicyStmt = append(policyParams.RemovePolicyStmt, policyStmt)
				break
			} else if action == "deny" {
				s.logger.Infof("UndoNSAction - policyParams=%+v, policyStmt=%+v, action deny\n",
					policyParams, policyStmt)
				policyParams.Accept = Accept
				break
			} else {
				s.logger.Errf("UndoNSAction - policyParams=%+v, policyStmt=%+v, unknown action=%s\n",
					policyParams, policyStmt, action)
			}
		}
	}
}

func (s *BGPServer) ApplyNSAction(actionInfo interface{}, conditionInfo []interface{}, policy utilspolicy.Policy,
	params interface{}, policyStmt utilspolicy.PolicyStmt) {
	policyParams := params.(*NSPolicyParams)

	s.logger.Infof("ApplyNSAction: ipPrefix=%s, policy=%s, statement=%s", policyParams.NS.IPPrefix, policy.Name,
		policyStmt.Name)
	s.logger.Infof("ApplyNSAction - actionInfo=%+v, conditionInfo=%+v, policyParams=%+v, policyStmt=%+v\n",
		actionInfo, conditionInfo, policyParams, policyStmt)
	if len(policyStmt.Actions) > 0 {
		for _, action := range policyStmt.Actions {
			if action == "permit" {
				s.logger.Infof("ApplyNSAction - policyParams=%+v, policyStmt=%+v, action permit\n",
					policyParams, policyStmt)
				policyParams.Accept = Accept
				if policyParams.AddPolicyStmt == nil {
					policyParams.AddPolicyStmt = make([]utilspolicy.PolicyStmt, 0)
				}
				policyParams.AddPolicyStmt = append(policyParams.AddPolicyStmt, policyStmt)
				break
			} else if action == "deny" {
				s.logger.Infof("ApplyNSAction - policyParams=%+v, policyStmt=%+v, action deny\n",
					policyParams, policyStmt)
				policyParams.Accept = Reject
				break
			} else {
				s.logger.Errf("ApplyNSAction - policyParams=%+v, policyStmt=%+v, unknown action=%s\n",
					policyParams, policyStmt, action)
			}
		}
	}
}

func (s *BGPServer) UpdateNSAndPolicyDB(policyDetails utilspolicy.PolicyDetails, params interface{}) {
	var op int
	policyParams := params.(*NSPolicyParams)
	s.logger.Infof("UpdateNSAndPolicyDB - ns=%+v", policyParams.NS)

	if policyParams.DeleteType != bgppolicy.Invalid {
		op = bgppolicy.Del
	} else {
		if policyDetails.EntityDeleted == false {
			s.logger.Info("Reject action was not applied, so add this policy to the route")
			op = bgppolicy.Add
			bgppolicy.UpdateNSPolicyState(policyParams.NS, op, policyDetails.Policy, policyDetails.PolicyStmt)
		}
		policyParams.NS.PolicyHitCounter++
	}
	s.networkStmtPE.UpdatePolicyNSMap(policyParams.NS, policyDetails.Policy, op)
}

func (s *BGPServer) DoesNSExist(params interface{}) bool {
	policyParams := params.(*NSPolicyParams)
	s.logger.Infof("DoesNSExist - ns=%+v", policyParams.NS)
	if af, ok := s.BgpConfig.Afs[policyParams.NS.AddressFamily]; ok {
		if _, ok := af.BgpNetworkStmts[policyParams.NS.IPPrefix]; ok {
			return true
		}
	}

	return false
}

func (s *BGPServer) getPrefixFromNSPolicy(data interface{}) (string, uint32) {
	s.logger.Infof("getPrefixFromNSPolicy - start")
	policyInfo := data.(utilspolicy.PolicyEngineApplyInfo)
	conditionsDB := s.networkStmtPE.PolicyEngine.PolicyConditionsDB
	var ipPrefix string
	var protoFamily uint32

	for _, condition := range policyInfo.ApplyPolicy.Conditions {
		s.logger.Infof("getPrefixFromNSPolicy - condition:%+v", condition)
		nodeGet := conditionsDB.Get(patriciaDB.Prefix(condition))
		if nodeGet == nil {
			s.logger.Err("Condition", condition, "not defined")
			return ipPrefix, 0
		}
		node := nodeGet.(utilspolicy.PolicyCondition)
		if node.ConditionType == policyCommonDefs.PolicyConditionTypeDstIpPrefixMatch {
			destIPCondition := node.ConditionInfo.(utilspolicy.MatchPrefixConditionInfo)
			prefix := destIPCondition.Prefix
			items := strings.Split(prefix.MasklengthRange, "-")
			if len(items) != 2 {
				s.logger.Err("IP prefix", destIPCondition.Prefix, "mask", prefix.MasklengthRange, "is not valid")
				return ipPrefix, 0
			}

			/*
				if _, err := strconv.Atoi(items[0]); err != nil {
					s.logger.Errf("Failed to convert prefex len %s to int with error %s", items[0], err)
					return ipPrefix, 0
				}
			*/

			totalBits, err := strconv.Atoi(items[1])
			if err != nil {
				s.logger.Errf("Failed to convert prefix max len %s to int with error %s", items[1], err)
				return ipPrefix, 0
			}

			protoFamily = packet.GetProtocolFamilyFromAddressLen(totalBits / 8)
			if protoFamily == 0 {
				s.logger.Errf("Failed to get protocol family from prefix max bits %d", totalBits)
				return ipPrefix, 0
			}

			ipPrefix = prefix.IpPrefix
		}
	}

	return ipPrefix, protoFamily
}

func (s *BGPServer) SetActionToNS(policyParams *NSPolicyParams) (*bgprib.Path, bool) {
	var path *bgprib.Path
	var accept, ok bool

	if policyParams.Accept == Accept {
		policyParams.NS.Accept = true
		accept = true
		key := ""
		if len(policyParams.RemovePolicyStmt) > 0 {
			s.logger.Infof("SetActionToNS - remove stmt=%+v, ns.policy stmt=%+v", policyParams.RemovePolicyStmt,
				policyParams.NS.PolicyStmtList)
			for _, policyStmt := range policyParams.RemovePolicyStmt {
				for idx, nsPolicyStmt := range policyParams.NS.PolicyStmtList {
					if policyStmt.Name == nsPolicyStmt {
						policyParams.NS.PolicyStmtList = append(policyParams.NS.PolicyStmtList[:idx],
							policyParams.NS.PolicyStmtList[idx+1:]...)
						break
					}
				}
			}
			s.logger.Infof("SetActionToNS - after mod, ns.policy stmt=%+v", policyParams.NS.PolicyStmtList)
		}

		if len(policyParams.AddPolicyStmt) > 0 {
			if policyParams.NS.PolicyStmtList == nil {
				policyParams.NS.PolicyStmtList = make([]string, 0)
			}
			for _, policyStmt := range policyParams.AddPolicyStmt {
				found := false
				for _, nsPolicyStmt := range policyParams.NS.PolicyStmtList {
					if policyStmt.Name == nsPolicyStmt {
						found = true
						break
					}
				}
				if !found {
					policyParams.NS.PolicyStmtList = append(policyParams.NS.PolicyStmtList, policyStmt.Name)
				}
			}
		}

		first := true
		for _, nsPolicyStmt := range policyParams.NS.PolicyStmtList {
			if !first {
				key = key + "+" + nsPolicyStmt
			} else {
				key = nsPolicyStmt
				first = false
			}
		}

		path, ok = s.NSRoutesPathMap[key]
		if !ok {
			pa := s.ConstructPathAttrsForNSRoutes()
			for _, policyStmt := range policyParams.AddPolicyStmt {
				pa, _ = bgppolicy.ApplyActionsToPacket(pa, policyStmt)
			}
			protoFamily := packet.GetProtocolFamily(packet.AfiIP6, packet.SafiUnicast)
			ipv6MPReach := packet.ConstructIPv6MPReachNLRIForConnRoutes(protoFamily)
			path = bgprib.NewPath(s.LocRib, nil, pa, ipv6MPReach, bgprib.RouteTypeStatic)
			s.NSRoutesPathMap[key] = path
		}
	} else {
		policyParams.NS.Accept = false
		accept = false
		path, ok = s.NSRoutesPathMap[""]
		if !ok {
			path = s.ConstructPathsForNSRoutes(&s.BgpConfig.Global.Config)
		}
	}
	return path, accept
}

func (s *BGPServer) GetNLRIForNSRoute(ns *config.BGPNetworkStatement, accept bool) (add,
	remove map[uint32][]packet.NLRI) {
	nlri, err := packet.ConstructIPPrefixFromCIDR(ns.IPPrefix)
	if err != nil {
		s.logger.Err("ProcessNSRoute - Could not convert prefix", ns.IPPrefix, "to NLRI")
		return add, remove
	}
	afNLRI := make(map[uint32][]packet.NLRI)
	afNLRI[ns.AddressFamily] = make([]packet.NLRI, 1)
	afNLRI[ns.AddressFamily][0] = nlri
	if accept {
		add = afNLRI
	} else {
		remove = afNLRI
	}
	return add, remove
}

func (s *BGPServer) TraverseAndApplyNS(data interface{}, updateFunc utilspolicy.PolicyApplyfunc) {
	s.logger.Infof("TraverseAndApplyNS - start")

	ipPrefix, protoFamily := s.getPrefixFromNSPolicy(data)
	s.logger.Infof("TraverseAndApplyNS - ipPrefix=%v, protoFamily=%v, s.BgpConfig.Afs=%+v", ipPrefix, protoFamily, s.BgpConfig.Afs)
	if bgpAFconfigs, ok := s.BgpConfig.Afs[protoFamily]; ok {
		s.logger.Infof("TraverseAndApplyNS - bgpAFconfigs=%+v", bgpAFconfigs)
		if ns, ok := bgpAFconfigs.BgpNetworkStmts[ipPrefix]; ok {
			peEntity := utilspolicy.PolicyEngineFilterEntityParams{
				DestNetIp:  ns.IPPrefix,
				PolicyList: ns.PolicyList,
			}

			callbackInfo := &NSPolicyParams{
				CreateType: utilspolicy.Invalid,
				DeleteType: utilspolicy.Invalid,
				Traverse:   true,
				NS:         ns,
				Accept:     Accept,
			}

			updateFunc(peEntity, data, callbackInfo)
			path, accept := s.SetActionToNS(callbackInfo)
			add, remove := s.GetNLRIForNSRoute(ns, accept)
			routerId := s.BgpConfig.Global.Config.RouterId.String()
			updated, withdrawn, updatedAddPaths := s.LocRib.ProcessConnectedRoutes(routerId, path, add, remove,
				s.AddPathCount)
			updated, withdrawn, updatedAddPaths = s.CheckForAggregation(updated, withdrawn, updatedAddPaths)
			s.SendUpdate(updated, withdrawn, updatedAddPaths)
		}
	}
}

func (s *BGPServer) TraverseAndReverseNS(policyData interface{}) {
	updateInfo := policyData.(utilspolicy.PolicyEngineApplyInfo)
	applyPolicyInfo := updateInfo.ApplyPolicy
	policy := applyPolicyInfo.ApplyPolicy
	s.logger.Info("TraverseAndReverseNS - policy", policy.Name)
	policyExtensions := policy.Extensions.(bgppolicy.NetworkStmtPolicyExtensions)
	s.logger.Info("TraverseAndReverseNS - policy extensions = %+v", policyExtensions)

	//	policyInfo := s.networkStmtPE.PolicyEngine.PolicyDB.Get(patriciaDB.Prefix(policy.Name))
	//	if policyInfo == nil {
	//		utils.Logger.Info("Unexpected:policyInfo nil for policy ", policy.Name)
	//		return
	//	}
	//	tempPolicy := policyInfo.(utilspolicy.Policy)
	//	policyExtensions = tempPolicy.Extensions.(bgppolicy.NetworkStmtPolicyExtensions)
	//	s.logger.Info("TraverseAndReverseNS - policy extensions = %+v", policyExtensions)

	if len(policyExtensions.PrefixList) == 0 {
		s.logger.Info("No route affected by this policy, so nothing to do")
		return
	}

	var ns *config.BGPNetworkStatement

	nsPolicyExtensions := make([]*config.BGPNetworkStatement, len(policyExtensions.PrefixInfoList))
	copy(nsPolicyExtensions, policyExtensions.PrefixInfoList)
	for idx := 0; idx < len(nsPolicyExtensions); idx++ {
		ns = nsPolicyExtensions[idx]

		callbackInfo := &NSPolicyParams{
			CreateType:  utilspolicy.Invalid,
			DeleteType:  utilspolicy.Invalid,
			TraverseRev: true,
			NS:          ns,
		}

		peEntity := utilspolicy.PolicyEngineFilterEntityParams{
			DestNetIp:  ns.IPPrefix,
			PolicyList: ns.PolicyList,
		}

		s.networkStmtPE.PolicyEngine.PolicyEngineUndoApplyPolicyForEntity(peEntity, updateInfo, callbackInfo)
		path, accept := s.SetActionToNS(callbackInfo)
		if !s.DoesNSExist(callbackInfo) {
			accept = false
		}
		add, remove := s.GetNLRIForNSRoute(ns, accept)
		routerId := s.BgpConfig.Global.Config.RouterId.String()
		updated, withdrawn, updatedAddPaths := s.LocRib.ProcessConnectedRoutes(routerId, path, add, remove,
			s.AddPathCount)
		updated, withdrawn, updatedAddPaths = s.CheckForAggregation(updated, withdrawn, updatedAddPaths)
		s.SendUpdate(updated, withdrawn, updatedAddPaths)

		s.networkStmtPE.DeleteNSPolicyState(ns, policy.Name)
		s.networkStmtPE.PolicyEngine.DeletePolicyEntityMapEntry(peEntity, policy.Name)
	}
}

func (s *BGPServer) UpdateNSPolicy(nsConf *config.BGPNetworkStatement) error {
	s.logger.Debug("UpdateNSPolicy")
	var err error
	var policyAction utilspolicy.PolicyAction

	if len(nsConf.Policy) != 0 {
		policyEngine := s.networkStmtPE.GetPolicyEngine()
		policyDB := policyEngine.PolicyDB

		nodeGet := policyDB.Get(patriciaDB.Prefix(nsConf.Policy))
		if nodeGet == nil {
			s.logger.Err("Policy ", nsConf.Policy, " not defined")
			return errors.New(fmt.Sprintf("Policy %s not found in policy engine", nsConf.Policy))
		}
		node := nodeGet.(utilspolicy.Policy)

		bytes := packet.GetAddressLengthForFamily(nsConf.AddressFamily)
		if bytes == -1 {
			s.logger.Err("Could not find number of bytes for aggregate prefix family", nsConf.AddressFamily)
			return errors.New(fmt.Sprintf("Could not find number of bytes for aggregate prefix family",
				nsConf.AddressFamily))
		}
		strBits := strconv.Itoa(bytes * 8)
		s.logger.Infof("AddOrUpdateAgg: agg = %s, bytes = %d, bits =%s", nsConf.IPPrefix, bytes, strBits)

		name := "NS+" + nsConf.IPPrefix
		tokens := strings.Split(nsConf.IPPrefix, "/")
		prefixLen := tokens[1]
		_, err = strconv.Atoi(prefixLen)
		if err != nil {
			s.logger.Errf("Failed to convert prefex len %s to int with error %s", prefixLen, err)
			return err
		}

		ipPrefix := nsConf.IPPrefix
		cond := utilspolicy.PolicyConditionConfig{
			Name:          name,
			ConditionType: "MatchDstIpPrefix",
			MatchDstIpPrefixConditionInfo: utilspolicy.PolicyDstIpMatchPrefixSetCondition{
				Prefix: utilspolicy.PolicyPrefix{
					IpPrefix:        ipPrefix,
					MasklengthRange: prefixLen + "-" + strBits,
				},
			},
		}

		_, err = s.networkStmtPE.CreatePolicyCondition(cond)
		if err != nil {
			s.logger.Errf("Failed to create policy condition to match ip prefix %s with error %s", ipPrefix, err)
			return err
		}

		conditionNameList := make([]string, 1)
		conditionNameList[0] = name

		policyAction = utilspolicy.PolicyAction{
			Name:       name,
			ActionType: policyCommonDefs.PolicyActionTypeNetworkStatementAdvertise,
		}

		s.logger.Debug("Calling applypolicy with conditionNameList: ", conditionNameList)
		s.networkStmtPE.UpdateApplyPolicy(utilspolicy.ApplyPolicyInfo{node, policyAction, conditionNameList}, true)
	} else {
		nsConf.Accept = true
		accept := true
		path, ok := s.NSRoutesPathMap[""]
		if !ok {
			path = s.ConstructPathsForNSRoutes(&s.BgpConfig.Global.Config)
		}
		add, remove := s.GetNLRIForNSRoute(nsConf, accept)
		routerId := s.BgpConfig.Global.Config.RouterId.String()
		updated, withdrawn, updatedAddPaths := s.LocRib.ProcessConnectedRoutes(routerId, path, add,
			remove, s.AddPathCount)
		updated, withdrawn, updatedAddPaths = s.CheckForAggregation(updated, withdrawn, updatedAddPaths)
		s.SendUpdate(updated, withdrawn, updatedAddPaths)
	}
	return nil
}

func (s *BGPServer) DeleteNS(nsConf *config.BGPNetworkStatement) error {
	policyEngine := s.networkStmtPE.GetPolicyEngine()
	policyDB := policyEngine.PolicyDB

	if pfNSMap, ok := s.BgpConfig.Afs[nsConf.AddressFamily]; ok {
		if ns, ok := pfNSMap.BgpNetworkStmts[nsConf.IPPrefix]; ok {
			delete(s.BgpConfig.Afs[nsConf.AddressFamily].BgpNetworkStmts, nsConf.IPPrefix)
			s.NSRoutesMutex.Lock()
			s.NSRoutes[ns.Idx] = s.NSRoutes[len(s.NSRoutes)-1]
			s.NSRoutes[ns.Idx].Idx = ns.Idx
			s.NSRoutes[len(s.NSRoutes)-1] = nil
			s.NSRoutes = s.NSRoutes[:len(s.NSRoutes)-1]
			s.NSRoutesMutex.Unlock()
		}
	} else {
		s.logger.Err("DeleteNS: Address family", nsConf.AddressFamily, "not found BGP config address families")
	}

	nodeGet := policyDB.Get(patriciaDB.Prefix(nsConf.Policy))
	if nodeGet == nil {
		s.logger.Err("NS ", nsConf, " not created yet")
		return errors.New(fmt.Sprintf("Policy %s not found in policy engine", nsConf.IPPrefix))
	}
	node := nodeGet.(utilspolicy.Policy)

	name := "NS+" + nsConf.IPPrefix
	conditionNameList := make([]string, 1)
	conditionNameList[0] = name
	policyAction := utilspolicy.PolicyAction{
		Name:       name,
		ActionType: policyCommonDefs.PolicyActionTypeNetworkStatementAdvertise,
	}

	s.networkStmtPE.UpdateUndoApplyPolicy(utilspolicy.ApplyPolicyInfo{node, policyAction, conditionNameList}, true)
	s.networkStmtPE.DeletePolicyCondition(name)
	return nil
}

func (s *BGPServer) AddOrUpdateNS(oldConf config.BGPNetworkStatement, newConf config.BGPNetworkStatement,
	attrSet []bool) error {
	s.logger.Info("AddOrUpdateNS")
	var err error

	bytes := packet.GetAddressLengthForFamily(newConf.AddressFamily)
	if bytes == -1 {
		s.logger.Err("Could not find number of bytes for aggregate prefix family", newConf.AddressFamily)
		return errors.New(fmt.Sprintf("Could not find number of bytes for aggregate prefix family",
			newConf.AddressFamily))
	}
	strBits := strconv.Itoa(bytes * 8)
	s.logger.Infof("AddOrUpdateNS: agg = %s, bytes = %d, bits =%s", newConf.IPPrefix, bytes, strBits)

	if oldConf.IPPrefix != "" {
		// Delete the policy
		s.DeleteNS(&oldConf)
	}

	if newConf.IPPrefix != "" {
		if _, ok := s.BgpConfig.Afs[newConf.AddressFamily]; !ok {
			s.BgpConfig.Afs[newConf.AddressFamily] = &config.AddressFamily{}
		}
		if s.BgpConfig.Afs[newConf.AddressFamily].BgpNetworkStmts == nil {
			s.BgpConfig.Afs[newConf.AddressFamily].BgpNetworkStmts = make(map[string]*config.BGPNetworkStatement)
		}
		s.BgpConfig.Afs[newConf.AddressFamily].BgpNetworkStmts[newConf.IPPrefix] = &newConf
		s.NSRoutesMutex.Lock()
		idx := len(s.NSRoutes)
		s.NSRoutes = append(s.NSRoutes, &newConf)
		newConf.Idx = idx
		s.NSRoutesMutex.Unlock()
		err = s.UpdateNSPolicy(&newConf)
		return err
	}
	return err
}

func (s *BGPServer) copyGlobalConf(gConf config.GlobalConfig) {
	// Don't create a new Global object. Peers have reference to the global object.
	s.BgpConfig.Global.Config.Vrf = gConf.Vrf
	s.BgpConfig.Global.Config.AS = gConf.AS
	s.BgpConfig.Global.Config.RouterId = gConf.RouterId
	s.BgpConfig.Global.Config.Disabled = gConf.Disabled
	s.BgpConfig.Global.Config.UseMultiplePaths = gConf.UseMultiplePaths
	s.BgpConfig.Global.Config.EBGPMaxPaths = gConf.EBGPMaxPaths
	s.BgpConfig.Global.Config.EBGPAllowMultipleAS = gConf.EBGPAllowMultipleAS
	s.BgpConfig.Global.Config.IBGPMaxPaths = gConf.IBGPMaxPaths
	s.BgpConfig.Global.Config.Defaultv4Route = gConf.Defaultv4Route
	s.BgpConfig.Global.Config.Defaultv6Route = gConf.Defaultv6Route
	s.BgpConfig.Global.Config.DefaultMED = gConf.DefaultMED
}

func (s *BGPServer) handleBfdNotifications(oper config.Operation, DestIp string,
	State bool) {
	if peer, ok := s.PeerMap[DestIp]; ok {
		if !State && peer.NeighborConf.Neighbor.State.BfdNeighborState == "up" {
			peer.BfdFaultSet()
		}
		if State && peer.NeighborConf.Neighbor.State.BfdNeighborState == "down" {
			peer.BfdFaultCleared()
		}
		s.logger.Info("Bfd state of peer ", peer.NeighborConf.Neighbor.NeighborAddress, " is ",
			peer.NeighborConf.Neighbor.State.BfdNeighborState)
	}
}

func (s *BGPServer) setInterfaceMapForPeer(peerIP string, peer *Peer) {
	s.logger.Info("Server: setInterfaceMapForPeer Peer", peer, "calling GetRouteReachabilityInfo")
	reachInfo, err := s.routeMgr.GetNextHopInfo(peerIP, -1)
	s.logger.Info("Server: setInterfaceMapForPeer Peer", peer, "GetRouteReachabilityInfo returned", reachInfo)
	if err != nil {
		s.logger.Infof("Server: Peer %s is not reachable", peerIP)
	} else {
		// @TODO: jgheewala think of something better for ovsdb....
		ifIdx := s.IntfMgr.GetIfIndex(int(reachInfo.NextHopIfIndex), int(reachInfo.NextHopIfType))
		// ifIdx := asicdCommonDefs.GetIfIndexFromIntfIdAndIntfType(int(reachInfo.NextHopIfIndex),
		//	int(reachInfo.NextHopIfType))
		s.logger.Infof("Server: Peer %s IfIdx %d", peerIP, ifIdx)
		if _, ok := s.IfIndexPeerMap[ifIdx]; !ok {
			s.IfIndexPeerMap[ifIdx] = make([]string, 0)
		}
		s.IfIndexPeerMap[ifIdx] = append(s.IfIndexPeerMap[ifIdx], peerIP)
		peer.setIfIdx(ifIdx)
	}
}

func (s *BGPServer) clearInterfaceMapForPeer(peerIP string, peer *Peer) {
	ifIdx := peer.getIfIdx()
	s.logger.Infof("Server: Peer %s FSM connection broken ifIdx %v", peerIP, ifIdx)
	if peerList, ok := s.IfIndexPeerMap[ifIdx]; ok {
		for idx, ip := range peerList {
			if ip == peerIP {
				s.IfIndexPeerMap[ifIdx] = append(s.IfIndexPeerMap[ifIdx][:idx],
					s.IfIndexPeerMap[ifIdx][idx+1:]...)
				if len(s.IfIndexPeerMap[ifIdx]) == 0 {
					delete(s.IfIndexPeerMap, ifIdx)
				}
				break
			}
		}
	}
	peer.setIfIdx(-1)
}

func (s *BGPServer) constructBGPGlobalState(gConf *config.GlobalConfig) {
	s.BgpConfig.Global.State.Vrf = gConf.Vrf
	s.BgpConfig.Global.State.AS = gConf.AS
	s.BgpConfig.Global.State.RouterId = gConf.RouterId
	s.BgpConfig.Global.State.Disabled = gConf.Disabled
	s.BgpConfig.Global.State.UseMultiplePaths = gConf.UseMultiplePaths
	s.BgpConfig.Global.State.EBGPMaxPaths = gConf.EBGPMaxPaths
	s.BgpConfig.Global.State.EBGPAllowMultipleAS = gConf.EBGPAllowMultipleAS
	s.BgpConfig.Global.State.IBGPMaxPaths = gConf.IBGPMaxPaths
	s.BgpConfig.Global.State.Defaultv4Route = gConf.Defaultv4Route
	s.BgpConfig.Global.State.Defaultv6Route = gConf.Defaultv6Route
	s.BgpConfig.Global.State.DefaultMED = gConf.DefaultMED
}

func (s *BGPServer) SetupRedistribution(newConf config.GlobalConfig) {
	s.logger.Info("SetUpRedistribution")
	if s.RedistributionMap == nil {
		s.RedistributionMap = make(map[string]string)
	}
	currentRedistMap := make(map[string]bool)
	applyList := make([]*config.ApplyPolicyInfo, 0)
	undoApplyList := make([]*config.ApplyPolicyInfo, 0)
	applyIndex := 0
	undoIndex := 0
	source := ""

	for i := 0; i < len(newConf.Redistribution); i++ {
		s.logger.Info("Sources: ", newConf.Redistribution[i].Sources)
		sources := make([]string, 0)
		sources = strings.Split(newConf.Redistribution[i].Sources, ",")
		s.logger.Infof("Setting up %s as redistribution policy for source(s): ", newConf.Redistribution[i].Policy)
		for j := 0; j < len(sources); j++ {
			source = sources[j]
			s.logger.Info("source: ", source)
			var condition *config.ConditionInfo
			condition = nil
			if sources[j] != "" {
				condition = &config.ConditionInfo{ConditionType: "MatchProtocol", Protocol: source}
			}
			_, ok := s.RedistributionMap[source]
			currentRedistMap[source] = true
			if !ok {
				s.logger.Info("No policy applied for this source so far")
				s.RedistributionMap[source] = newConf.Redistribution[i].Policy
				applyList = append(applyList, &config.ApplyPolicyInfo{
					Protocol: "BGP",
					Policy:   newConf.Redistribution[i].Policy,
					Action:   "Redistribution"})
				if condition != nil {
					if applyList[applyIndex].Conditions == nil {
						applyList[applyIndex].Conditions = make([]*config.ConditionInfo, 0)
					}
					applyList[applyIndex].Conditions = append(applyList[applyIndex].Conditions, condition)
				}
				applyIndex++
			} else if s.RedistributionMap[source] == newConf.Redistribution[i].Policy {
				s.logger.Info("Policy unchanged for source ", source)
				continue
			} else {
				s.logger.Info("Another policy:", s.RedistributionMap[source], " already applied for source :", source)
				applyList = append(applyList, &config.ApplyPolicyInfo{
					Protocol: "BGP",
					Policy:   newConf.Redistribution[i].Policy,
					Action:   "Redistribution"})
				if condition != nil {
					if applyList[applyIndex].Conditions == nil {
						applyList[applyIndex].Conditions = make([]*config.ConditionInfo, 0)
					}
					applyList[applyIndex].Conditions = append(applyList[applyIndex].Conditions, condition)
				}
				applyIndex++
				if s.RedistributionMap[source] != "" {
					undoApplyList = append(undoApplyList, &config.ApplyPolicyInfo{
						Protocol: "BGP",
						Policy:   s.RedistributionMap[source],
						Action:   "Redistribution"})
					if condition != nil {
						if undoApplyList[undoIndex].Conditions == nil {
							undoApplyList[undoIndex].Conditions = make([]*config.ConditionInfo, 0)
						}
						undoApplyList[undoIndex].Conditions = append(undoApplyList[undoIndex].Conditions, condition)
					}
					undoIndex++
				}

				s.RedistributionMap[source] = newConf.Redistribution[i].Policy
			}
		}
	}

	for source, policy := range s.RedistributionMap {
		if !currentRedistMap[source] {
			var condition *config.ConditionInfo
			if policy != "" {
				condition = &config.ConditionInfo{ConditionType: "MatchProtocol", Protocol: source}
				undoApplyList = append(undoApplyList, &config.ApplyPolicyInfo{
					Protocol: "BGP",
					Policy:   s.RedistributionMap[source],
					Action:   "Redistribution"})
				if undoApplyList[undoIndex].Conditions == nil {
					undoApplyList[undoIndex].Conditions = make([]*config.ConditionInfo, 0)
				}
				undoApplyList[undoIndex].Conditions = append(undoApplyList[undoIndex].Conditions, condition)
				undoIndex++
			}
			delete(s.RedistributionMap, source)
		}
	}

	if len(applyList) > 0 || len(undoApplyList) > 0 {
		s.routeMgr.ApplyPolicy(applyList, undoApplyList)
	}
}

func (s *BGPServer) UpdateGlobalForPatchUpdate(oldConfig, newConfig config.GlobalConfig, op []*bgpd.PatchOpInfo) {
	s.logger.Info("UpdateGlobalForPatchUpdate")
	for idx := 0; idx < len(op); idx++ {
		s.logger.Debug("patch update")
		switch op[idx].Path {
		case "Redistribution":
			s.logger.Debug("Patch update for redistribution")
			applyList := make([]*config.ApplyPolicyInfo, 0)
			undoApplyList := make([]*config.ApplyPolicyInfo, 0)
			if len(op[idx].Value) == 0 {
				/*
					If redistribution update is trying to update redistribution, non zero value is expected
				*/
				s.logger.Err("Must specify sources")
				return
			}
			s.logger.Debug("value = ", op[idx].Value)
			valueObjArr := []bgpd.SourcePolicyList{}
			err := json.Unmarshal([]byte(op[idx].Value), &valueObjArr)
			if err != nil {
				s.logger.Err("error unmarshaling value:", err)
				return
			}
			s.logger.Debug("Number of redistribution soures:", len(valueObjArr))
			for _, val := range valueObjArr {
				policy := val.Policy
				sources := make([]string, 0)
				sources = strings.Split(val.Sources, ",")
				s.logger.Infof("Setting up %s as redistribution policy for source(s): ", policy)
				for j := 0; j < len(sources); j++ {
					var condition *config.ConditionInfo
					condition = nil
					source := sources[j]
					s.logger.Info("source: ", source)
					if sources[j] != "" {
						condition = &config.ConditionInfo{ConditionType: "MatchProtocol", Protocol: source}
					}
					switch op[idx].Op {
					case "add":
						s.logger.Debug("add op: source:", source, " policy:", policy)
						_, ok := s.RedistributionMap[source]
						if !ok {
							s.logger.Info("No policy applied for this source so far")
							s.RedistributionMap[source] = policy
							if len(applyList) == 0 {
								applyList = append(applyList, &config.ApplyPolicyInfo{
									Protocol: "BGP",
									Policy:   policy,
									Action:   "Redistribution"})
							}
							if condition != nil {
								if applyList[0].Conditions == nil {
									applyList[0].Conditions = make([]*config.ConditionInfo, 0)
								}
								applyList[0].Conditions = append(applyList[0].Conditions, condition)
							}
							if len(applyList) > 0 || len(undoApplyList) > 0 {
								s.routeMgr.ApplyPolicy(applyList, undoApplyList)
							}
						} else {
							s.logger.Err("Cannot add policy for source:", source, " there is already a policy ,",
								s.RedistributionMap[source], " applied")
						}
					case "remove":
						s.logger.Debug("remove op: source:", source, " policy:", policy)
						_, ok := s.RedistributionMap[source]
						if !ok {
							s.logger.Err("No policy applied for source:", source, " nothing to be removed")
						} else if policy != "" && s.RedistributionMap[source] != policy {
							s.logger.Err("Policy applied", s.RedistributionMap,
								"is not the same as policy being removed:", policy)
						} else {
							if len(undoApplyList) == 0 {
								undoApplyList = append(undoApplyList, &config.ApplyPolicyInfo{
									Protocol: "BGP",
									Policy:   s.RedistributionMap[source],
									Action:   "Redistribution"})
							}
							if condition != nil {
								if undoApplyList[0].Conditions == nil {
									undoApplyList[0].Conditions = make([]*config.ConditionInfo, 0)
								}
								undoApplyList[0].Conditions = append(undoApplyList[0].Conditions, condition)
							}
							delete(s.RedistributionMap, source)
							if len(applyList) > 0 || len(undoApplyList) > 0 {
								s.routeMgr.ApplyPolicy(applyList, undoApplyList)
							}
						}
					default:
						s.logger.Err("operation ", op[idx].Op, " not supported")
						return
					}
				}
			}
		default:
			s.logger.Err("Patch update for attribute:", op[idx].Path, " not supported")
			return
		}
	}
}

func (s *BGPServer) UpdateGlobal(bgpGlobal *bgpd.BGPGlobal, oldConfig, newConfig config.GlobalConfig, attrSet []bool) {
	s.logger.Info("UpdateGlobal")
	if bgpGlobal == nil {
		s.logger.Err("bgpglobal nil in update")
		return
	}

	defaultRoutes := make(map[uint32]bool)

	if attrSet != nil {
		objTyp := reflect.TypeOf(*bgpGlobal)
		restart := false
		for i := 0; i < objTyp.NumField(); i++ {
			objName := objTyp.Field(i).Name
			if attrSet[i] {
				s.logger.Debug("UpdateGlobal: changed ", objName)
				if objName == "Redistribution" {
					s.SetupRedistribution(newConfig)
				} else if objName == "Defaultv4Route" {
					protoFamily := packet.GetProtocolFamily(packet.AfiIP, packet.SafiUnicast)
					defaultRoutes[protoFamily] = newConfig.Defaultv4Route
				} else if objName == "Defaultv6Route" {
					protoFamily := packet.GetProtocolFamily(packet.AfiIP6, packet.SafiUnicast)
					defaultRoutes[protoFamily] = newConfig.Defaultv6Route
				} else {
					restart = true
				}
			}
		}

		if restart {
			s.Restart(newConfig)
		} else if len(defaultRoutes) > 0 {
			s.UpdateDefaultRoutes(defaultRoutes)
		}
	}
}

func (s *BGPServer) isBGPGlobalDisabled() bool {
	return s.BgpConfig.Global.Config.Disabled
}

func (s *BGPServer) Restart(cfg config.GlobalConfig) {
	s.logger.Info("Restart BGP")
	for peerIP, peer := range s.PeerMap {
		s.logger.Infof("Cleanup peer %s", peerIP)
		peer.Cleanup()
	}
	s.logger.Infof("Giving up CPU so that all peer FSMs will get cleaned up")
	runtime.Gosched()

	s.RemoveRoutesFromAllNeighbor()

	gConf := cfg
	s.copyGlobalConf(gConf)
	s.constructBGPGlobalState(&gConf)

	for _, peer := range s.PeerMap {
		peer.UpdateGlobal(&s.BgpConfig.Global.Config)
	}

	if s.isBGPGlobalDisabled() {
		s.logger.Info("BGP global for Vrf", gConf.Vrf, "is disabled, not bringing the neighbors up.")
		return
	}

	s.ConstructPathsForLocalRoutes(&s.BgpConfig.Global.Config)
	s.ConstructDefaultRoutes(&s.BgpConfig.Global.Config)

	add, remove := s.routeMgr.GetRoutes()
	if add != nil && remove != nil {
		s.ProcessConnectedRoutes(add, remove)
	}

	for _, peer := range s.PeerMap {
		peer.Init()
	}
	//s.SetupRedistribution(gConf)
	// Get routes from the route manager
}

func (s *BGPServer) updateGlobalConfig(bgpGlobal *bgpd.BGPGlobal, oldConfig, newConfig config.GlobalConfig,
	attrSet []bool, op []*bgpd.PatchOpInfo) {
	s.logger.Info("updateGlobalConfig")
	if op == nil || len(op) == 0 {
		s.UpdateGlobal(bgpGlobal, oldConfig, newConfig, attrSet)
	} else {
		s.UpdateGlobalForPatchUpdate(oldConfig, newConfig, op)
	}
}

func (s *BGPServer) getIfaceIP(ifIndex int32, peerAddrType config.PeerAddressType) net.IP {
	ipInfo, err := s.GetIfaceIP(ifIndex)
	s.logger.Info("ipInfo:", ipInfo, " err:", err)
	if err != nil {
		s.logger.Errf("IP not configured on interface %d yet", ifIndex)
		return nil
	}

	if peerAddrType == config.PeerAddressV4 {
		if ipInfo.IpAddr == nil {
			return nil
		}
		ip := ipInfo.IpAddr
		ifIP := make(net.IP, len(ip))
		copy(ifIP, ip)
		ipMask := ipInfo.IpMask
		if ipMask[len(ipMask)-1] < 252 {
			s.logger.Err("IPv4Addr", ifIP, "of the interface", ifIndex, "is not /30 or /31 address")
			return nil
		}
		s.logger.Info("IPv4Addr of the v4Neighbor local interface", ifIndex, "is", ifIP)
		ifIP[len(ifIP)-1] = ifIP[len(ifIP)-1] ^ (^ipMask[len(ipMask)-1])
		s.logger.Info("Peer IPv4Addr of the v4Neighbor interface", ifIndex, "is", ifIP)
		return ifIP
	} else if peerAddrType == config.PeerAddressV6 {
		if ipInfo.LinklocalIpAddr != "" {
			return net.ParseIP(ipInfo.LinklocalIpAddr)
		}
	} else {
		s.logger.Err("getIfaceIP - Unknown peer address type", peerAddrType, "ifIndex", ifIndex)
	}
	return nil
}

func (s *BGPServer) handleIntfCreate(ifIndex int32, peerAddrType config.PeerAddressType) {
	s.logger.Infof("handleIntfCreate - ifIndex:%d, peerAddrType:%d", ifIndex, peerAddrType)
	ip := s.getIfaceIP(ifIndex, peerAddrType)
	if ip != nil {
		if _, ok := s.ifaceNeighbors[peerAddrType]; ok {
			if peer, ok := s.ifaceNeighbors[peerAddrType][ifIndex]; ok {
				peer.SetNeighborAddress(ip)
				s.PeerMap[ip.String()] = peer
				peer.Init()
			} else {
				s.logger.Infof("handleIntfCreate - ifIndex not found in neighbors mape:%+v, ifIndex:%d, peerAddrType:%d",
					s.ifaceNeighbors[peerAddrType], ifIndex, peerAddrType)
			}
		} else {
			s.logger.Infof("handleIntfCreate - addr type not found in neighbors mape:%+v, ifIndex:%d, peerAddrType:%d",
				s.ifaceNeighbors, ifIndex, peerAddrType)
		}
	}
}

func (s *BGPServer) handleIntfDelete(ifIndex int32, peerAddrType config.PeerAddressType) {
	s.logger.Infof("handleIntfDelete - ifIndex:%d, peerAddrType:%d", ifIndex, peerAddrType)
	if _, ok := s.ifaceNeighbors[peerAddrType]; ok {
		if peer, ok := s.ifaceNeighbors[peerAddrType][ifIndex]; ok {
			ip := peer.NeighborConf.RunningConf.NeighborAddress
			delete(s.PeerMap, ip.String())
			peer.Cleanup()
			if ip != nil {
				s.ProcessRemoveNeighbor(ip.String(), peer)
			}
			peer.ResetNeighborAddress()
		} else {
			s.logger.Infof("handleIntfDelete - ifIndex not found in neighbors mape:%+v, ifIndex:%d, peerAddrType:%d",
				s.ifaceNeighbors[peerAddrType], ifIndex, peerAddrType)
		}
	} else {
		s.logger.Infof("handleIntfDelete - addr type not found in neighbors mape:%+v, ifIndex:%d, peerAddrType:%d",
			s.ifaceNeighbors, ifIndex, peerAddrType)
	}
}

func (s *BGPServer) CreatePeer(newPeer config.NeighborConfig) {
	s.logger.Infof("CreatePeer %+v", newPeer)
	var ok bool
	var peer *Peer

	if newPeer.NeighborAddress != nil {
		if _, ok = s.PeerMap[newPeer.NeighborAddress.String()]; ok {
			s.logger.Infof("Failed to add neighbor. Neighbor at address %s already exists", newPeer.NeighborAddress)
			return
		}
	}

	if newPeer.IfIndex != -1 {
		if _, ok = s.ifaceNeighbors[newPeer.PeerAddressType]; !ok {
			s.logger.Infof("Failed to add neighbor. Peer address type", newPeer.PeerAddressType, "not supported")
			return
		}

		if _, ok = s.ifaceNeighbors[newPeer.PeerAddressType][newPeer.IfIndex]; ok {
			s.logger.Infof("Failed to add neighbor. Neighbor at interface %d already exists", newPeer.IfIndex)
			return
		}

		s.logger.Info("Add iface neighbor, ip:", newPeer.NeighborAddress.String(), "ifIndex:", newPeer.IfIndex,
			"ifName:", newPeer.IfName)
		newPeer.NeighborAddress = s.getIfaceIP(newPeer.IfIndex, newPeer.PeerAddressType)
	}

	var groupConfig *config.PeerGroupConfig
	if newPeer.PeerGroup != "" {
		protoFamily, _ := packet.GetProtocolFamilyFromPeerAddrType(newPeer.PeerAddressType)
		if _, ok := s.BgpConfig.PeerGroups[protoFamily]; !ok {
			s.logger.Info("No peer groups with Peer address type", newPeer.PeerAddressType, "exists")
		} else {
			if group, ok := s.BgpConfig.PeerGroups[protoFamily][newPeer.PeerGroup]; !ok {
				s.logger.Info("Peer group", newPeer.PeerGroup, "not created yet, creating peer",
					newPeer.NeighborAddress.String(), "without the group")
			} else {
				groupConfig = &group.Config
			}
		}
	}

	s.logger.Info("Add neighbor, ip:", newPeer.NeighborAddress.String(), "ifIndex:", newPeer.IfIndex)
	peer = NewPeer(s, s.LocRib, &s.BgpConfig.Global.Config, groupConfig, newPeer)
	if peer.NeighborConf.RunningConf.NeighborAddress.To4() != nil &&
		peer.NeighborConf.RunningConf.AuthPassword != "" {
		err := netUtils.SetTCPListenerMD5(s.listener, newPeer.NeighborAddress.String(),
			peer.NeighborConf.RunningConf.AuthPassword)
		if err != nil {
			s.logger.Info("Failed to add MD5 authentication for neighbor",
				newPeer.NeighborAddress.String(), "with error", err)
		}
	}

	if newPeer.NeighborAddress != nil {
		s.PeerMap[newPeer.NeighborAddress.String()] = peer
	}

	if newPeer.IfIndex != -1 {
		s.ifaceNeighbors[newPeer.PeerAddressType][newPeer.IfIndex] = peer
	}

	s.NeighborMutex.Lock()
	s.addPeerToList(peer)
	s.NeighborMutex.Unlock()
	if s.isBGPGlobalDisabled() {
		s.logger.Info("BGP global", s.BgpConfig.Global.Config.Vrf, "is disabled, not activating neighbor",
			newPeer.NeighborAddress)
		return
	}
	peer.Init()
}

func (s *BGPServer) getPeer(neighbor config.NeighborConfig) *Peer {
	var peer *Peer
	var ok bool
	if neighbor.NeighborAddress != nil {
		if peer, ok = s.PeerMap[neighbor.NeighborAddress.String()]; !ok {
			s.logger.Err("Peer not found for address", neighbor.NeighborAddress.String())
		}
	} else if neighbor.IfIndex != -1 {
		if _, ok = s.ifaceNeighbors[neighbor.PeerAddressType]; ok {
			if peer, ok = s.ifaceNeighbors[neighbor.PeerAddressType][neighbor.IfIndex]; !ok {
				s.logger.Err("Peer not found for ifIndex", neighbor.IfIndex)
			}
		}
	}

	return peer
}

func (s *BGPServer) updatePeerConf(oldPeer, newPeer config.NeighborConfig, peer *Peer) {
	s.logger.Info("Clean up peer, ip:", oldPeer.NeighborAddress.String(), "ifIndex:", oldPeer.IfIndex)
	peer.Cleanup()
	if peer.NeighborConf.RunningConf.NeighborAddress != nil {
		s.ProcessRemoveNeighbor(peer.NeighborConf.RunningConf.NeighborAddress.String(), peer)
		if peer.NeighborConf.RunningConf.NeighborAddress.To4() != nil &&
			peer.NeighborConf.RunningConf.AuthPassword != "" {
			err := netUtils.SetTCPListenerMD5(s.listener, peer.NeighborConf.RunningConf.NeighborAddress.String(), "")
			if err != nil {
				s.logger.Info("Failed to add MD5 authentication for old neighbor",
					newPeer.NeighborAddress.String(), "with error", err)
			}
		}
	}
	peer.UpdateNeighborConf(newPeer, &s.BgpConfig)

	runtime.Gosched()

	if newPeer.PeerGroup != "" {
		protoFamily, _ := packet.GetProtocolFamilyFromPeerAddrType(newPeer.PeerAddressType)
		if _, ok := s.BgpConfig.PeerGroups[protoFamily]; !ok {
			s.logger.Info("No peer groups with Peer address type", newPeer.PeerAddressType, "exists")
		} else {
			if group, ok := s.BgpConfig.PeerGroups[protoFamily][newPeer.PeerGroup]; !ok {
				s.logger.Info("Peer group", newPeer.PeerGroup, "not created yet, creating peer",
					newPeer.NeighborAddress.String(), "without the group")
			} else {
				peer.UpdatePeerGroup(&group.Config)
			}
		}
	}

	if peer.NeighborConf.RunningConf.NeighborAddress.To4() != nil &&
		peer.NeighborConf.RunningConf.AuthPassword != "" {
		err := netUtils.SetTCPListenerMD5(s.listener, newPeer.NeighborAddress.String(),
			peer.NeighborConf.RunningConf.AuthPassword)
		if err != nil {
			s.logger.Info("Failed to add MD5 authentication for neighbor",
				newPeer.NeighborAddress.String(), "with error", err)
		}
	}

	if s.isBGPGlobalDisabled() {
		s.logger.Info("BGP global", s.BgpConfig.Global.Config.Vrf, "is disabled, not activating neighbor",
			newPeer.NeighborAddress)
		return
	}

	peer.Init()
}

func (s *BGPServer) Updatev4Peer(bgpPeer *bgpd.BGPv4Neighbor, oldPeer, newPeer config.NeighborConfig, attrSet []bool) {
	s.logger.Info("Updatev4Peer")
	var peer *Peer

	if peer = s.getPeer(oldPeer); peer == nil {
		s.logger.Err("Updatev4Peer - Peer not found for ip:", oldPeer.NeighborAddress, "ifIndex:", oldPeer.IfIndex)
		return
	}

	if attrSet != nil {
		updateConfig := false
		s.logger.Info("attrSet:", attrSet)
		objTyp := reflect.TypeOf(*bgpPeer)
		s.logger.Info("numfield:", objTyp.NumField())
		for i := 0; i < objTyp.NumField(); i++ {
			objName := objTyp.Field(i).Name
			if attrSet[i] {
				s.logger.Debug("UpdatePeer : changed ", objName)
				if objName == "AdjRIBInFilter" {
					s.logger.Info("Update AdjRIBInFilter to ", newPeer.AdjRIBInFilter)
					if oldPeer.AdjRIBInFilter != "" {
						s.logger.Info("old filter applied was :", oldPeer.AdjRIBInFilter,
							"p.NeighborConf.RunningConf.AdjRIBInFilter:", peer.NeighborConf.RunningConf.AdjRIBInFilter)
						peer.RemoveAdjRIBFilter(peer.server.ribInPE, peer.NeighborConf.RunningConf.AdjRIBInFilter,
							bgprib.AdjRIBDirIn)
					}
					peer.UpdateNeighborConf(newPeer, &s.BgpConfig)
					if newPeer.AdjRIBInFilter != "" {
						s.logger.Info("new filter being applied is:", newPeer.AdjRIBInFilter,
							"p.NeighborConf.RunningConf.AdjRIBInFilter:", peer.NeighborConf.RunningConf.AdjRIBInFilter)
						peer.AddAdjRIBFilter(peer.server.ribInPE, peer.NeighborConf.RunningConf.AdjRIBInFilter,
							bgprib.AdjRIBDirIn)
					}
				} else { // this needs to be changed to if cases of each and every settable attribute
					updateConfig = true
				}
			}
		}

		if updateConfig {
			s.updatePeerConf(oldPeer, newPeer, peer)
		}
	}
}

func (s *BGPServer) Updatev6Peer(bgpPeer *bgpd.BGPv6Neighbor, oldPeer, newPeer config.NeighborConfig, attrSet []bool) {
	s.logger.Info("Updatev4Peer")
	var peer *Peer

	if peer = s.getPeer(oldPeer); peer == nil {
		s.logger.Err("Updatev4Peer - Peer not found for ip:", oldPeer.NeighborAddress, "ifIndex:", oldPeer.IfIndex)
		return
	}

	if attrSet != nil {
		updateConfig := false
		s.logger.Info("attrSet:", attrSet)
		objTyp := reflect.TypeOf(*bgpPeer)
		s.logger.Info("numfield:", objTyp.NumField())
		for i := 0; i < objTyp.NumField(); i++ {
			objName := objTyp.Field(i).Name
			if attrSet[i] {
				s.logger.Debug("UpdatePeer : changed ", objName)
				if objName == "AdjRIBInFilter" {
					s.logger.Info("Update AdjRIBInFilter to ", newPeer.AdjRIBInFilter)
					if oldPeer.AdjRIBInFilter != "" {
						s.logger.Info("old filter applied was :", oldPeer.AdjRIBInFilter,
							"p.NeighborConf.RunningConf.AdjRIBInFilter:", peer.NeighborConf.RunningConf.AdjRIBInFilter)
						peer.RemoveAdjRIBFilter(peer.server.ribInPE, peer.NeighborConf.RunningConf.AdjRIBInFilter,
							bgprib.AdjRIBDirIn)
					}
					peer.UpdateNeighborConf(newPeer, &s.BgpConfig)
					if newPeer.AdjRIBInFilter != "" {
						s.logger.Info("new filter being applied is:", newPeer.AdjRIBInFilter,
							"p.NeighborConf.RunningConf.AdjRIBInFilter:", peer.NeighborConf.RunningConf.AdjRIBInFilter)
						peer.AddAdjRIBFilter(peer.server.ribInPE, peer.NeighborConf.RunningConf.AdjRIBInFilter,
							bgprib.AdjRIBDirIn)
					}
				} else { // this needs to be changed to if cases of each and every settable attribute
					updateConfig = true
				}
			}
		}

		if updateConfig {
			s.updatePeerConf(oldPeer, newPeer, peer)
		}
	}
}

func (s *BGPServer) removePeer(neighbor config.NeighborConfig) {
	s.logger.Info("Remove Peer, ip:", neighbor.NeighborAddress, "ifIndex:", neighbor.IfIndex)
	var peerIP string
	var ifacePeer *Peer
	var ok bool

	if neighbor.NeighborAddress != nil {
		peerIP = neighbor.NeighborAddress.String()
	}

	if neighbor.IfIndex != -1 {
		if _, ok = s.ifaceNeighbors[neighbor.PeerAddressType]; ok {
			if ifacePeer, ok = s.ifaceNeighbors[neighbor.PeerAddressType][neighbor.IfIndex]; ok {
				if ifacePeer.NeighborConf.RunningConf.NeighborAddress != nil {
					peerIP = ifacePeer.NeighborConf.RunningConf.NeighborAddress.String()
				}
				delete(s.ifaceNeighbors[neighbor.PeerAddressType], neighbor.IfIndex)
			}
		}
	}

	if peer, ok := s.PeerMap[peerIP]; ok {
		s.NeighborMutex.Lock()
		s.removePeerFromList(peer)
		s.NeighborMutex.Unlock()
		delete(s.PeerMap, peerIP)
		peer.Cleanup()
		s.ProcessRemoveNeighbor(peerIP, peer)
	} else if ifacePeer != nil {
		s.NeighborMutex.Lock()
		s.removePeerFromList(ifacePeer)
		s.NeighborMutex.Unlock()
		ifacePeer.Cleanup()
	}
}

func (s *BGPServer) listenChannelUpdates() {
	for {
		select {
		case globalUpdate := <-s.GlobalConfigCh:
			if globalUpdate.Op == "create" {
				s.Restart(globalUpdate.NewConfig)
			} else if globalUpdate.Op == "update" {
				s.updateGlobalConfig(globalUpdate.BGPConfig, globalUpdate.OldConfig, globalUpdate.NewConfig,
					globalUpdate.AttrSet, globalUpdate.PatchOp)
			}

		case peerUpdate := <-s.AddPeerCh:
			s.logger.Info("message received on AddPeerCh")
			oldPeer := peerUpdate.OldPeer
			newPeer := peerUpdate.NewPeer
			if peerUpdate.Op == "create" {
				s.CreatePeer(newPeer)
			} else if peerUpdate.Op == "update" {
				if peerUpdate.PeerType == "v4" {
					peer := peerUpdate.BGPPeer.(*bgpd.BGPv4Neighbor)
					s.Updatev4Peer(peer, oldPeer, newPeer, peerUpdate.AttrSet)
				} else {
					peer := peerUpdate.BGPPeer.(*bgpd.BGPv6Neighbor)
					s.Updatev6Peer(peer, oldPeer, newPeer, peerUpdate.AttrSet)
				}
			}
			/*
				var peer *Peer
				var ok bool
				if oldPeer.NeighborAddress != nil {
					if peer, ok = s.PeerMap[oldPeer.NeighborAddress.String()]; ok {
						s.logger.Info("Clean up peer", oldPeer.NeighborAddress.String())
						peer.Cleanup()
						s.ProcessRemoveNeighbor(oldPeer.NeighborAddress.String(), peer)
						if peer.NeighborConf.RunningConf.NeighborAddress.To4() != nil &&
							peer.NeighborConf.RunningConf.AuthPassword != "" {
							err := netUtils.SetTCPListenerMD5(s.listener, oldPeer.NeighborAddress.String(), "")
							if err != nil {
								s.logger.Info("Failed to add MD5 authentication for old neighbor",
									newPeer.NeighborAddress.String(), "with error", err)
							}
						}
						peer.UpdateNeighborConf(newPeer, &s.BgpConfig)

						runtime.Gosched()
					} else {
						s.logger.Info("Can't find neighbor with old address", oldPeer.NeighborAddress.String())
					}
				}

				if !ok {
					_, ok = s.PeerMap[newPeer.NeighborAddress.String()]
					if ok {
						s.logger.Info("Failed to add neighbor. Neighbor at that address already exists,",
							newPeer.NeighborAddress.String())
						break
					}

					var groupConfig *config.PeerGroupConfig
					if newPeer.PeerGroup != "" {
						protoFamily, _ := packet.GetProtocolFamilyFromPeerAddrType(newPeer.PeerAddressType)
						if _, ok := s.BgpConfig.PeerGroups[protoFamily]; !ok {
							s.logger.Info("No peer groups with Peer address type", newPeer.PeerAddressType, "exists")
						} else {
							if group, ok := s.BgpConfig.PeerGroups[protoFamily][newPeer.PeerGroup]; !ok {
								s.logger.Info("Peer group", newPeer.PeerGroup, "not created yet, creating peer",
									newPeer.NeighborAddress.String(), "without the group")
							} else {
								groupConfig = &group.Config
							}
						}
					}
					s.logger.Info("Add neighbor, ip:", newPeer.NeighborAddress.String())
					peer = NewPeer(s, s.LocRib, &s.BgpConfig.Global.Config, groupConfig, newPeer)
					if peer.NeighborConf.RunningConf.NeighborAddress.To4() != nil &&
						peer.NeighborConf.RunningConf.AuthPassword != "" {
						err := netUtils.SetTCPListenerMD5(s.listener, newPeer.NeighborAddress.String(),
							peer.NeighborConf.RunningConf.AuthPassword)
						if err != nil {
							s.logger.Info("Failed to add MD5 authentication for neighbor",
								newPeer.NeighborAddress.String(), "with error", err)
						}
					}
					s.PeerMap[newPeer.NeighborAddress.String()] = peer
					s.NeighborMutex.Lock()
					s.addPeerToList(peer)
					s.NeighborMutex.Unlock()
				}
				peer.Init()
			*/

		case remPeer := <-s.RemPeerCh:
			s.removePeer(remPeer)

		case groupUpdate := <-s.AddPeerGroupCh:
			oldGroupConf := groupUpdate.OldGroup
			newGroupConf := groupUpdate.NewGroup
			s.logger.Info("Peer group update old:", oldGroupConf, "new:", newGroupConf)
			var ok bool

			if oldGroupConf.Name != "" {
				protoFamily, _ := packet.GetProtocolFamilyFromPeerAddrType(oldGroupConf.PeerAddressType)
				if _, ok = s.BgpConfig.PeerGroups[protoFamily]; ok {
					if _, ok = s.BgpConfig.PeerGroups[protoFamily][oldGroupConf.Name]; !ok {
						s.logger.Err("Could not find peer group", oldGroupConf.Name)
					}
				}
			}

			protoFamily, _ := packet.GetProtocolFamilyFromPeerAddrType(newGroupConf.PeerAddressType)
			if _, ok = s.BgpConfig.PeerGroups[protoFamily]; !ok {
				s.BgpConfig.PeerGroups[protoFamily] = make(map[string]*config.PeerGroup)
			}
			if _, ok = s.BgpConfig.PeerGroups[protoFamily][newGroupConf.Name]; !ok {
				s.logger.Info("Add new peer group with name", newGroupConf.Name)
				peerGroup := config.PeerGroup{
					Config: newGroupConf,
				}
				s.BgpConfig.PeerGroups[protoFamily][newGroupConf.Name] = &peerGroup
			} else {
				s.logger.Info("Update peer group", newGroupConf.Name)
				s.BgpConfig.PeerGroups[protoFamily][newGroupConf.Name].Config = newGroupConf
			}
			s.UpdatePeerGroupInPeers(newGroupConf.Name, newGroupConf.PeerAddressType, &newGroupConf)

		case group := <-s.RemPeerGroupCh:
			s.logger.Info("Remove Peer group:", group.Name)
			protoFamily, _ := packet.GetProtocolFamilyFromPeerAddrType(group.PeerAddressType)
			if _, ok := s.BgpConfig.PeerGroups[protoFamily]; !ok {
				s.logger.Err("Peer group address type", group.PeerAddressType, "not found in map")
				break
			}

			if _, ok := s.BgpConfig.PeerGroups[protoFamily][group.Name]; !ok {
				s.logger.Err("Peer group", group.Name, "not found in map")
				break
			}
			delete(s.BgpConfig.PeerGroups[protoFamily], group.Name)
			s.UpdatePeerGroupInPeers(group.Name, group.PeerAddressType, nil)

		case aggUpdate := <-s.AddAggCh:
			oldAgg := aggUpdate.OldAgg
			newAgg := aggUpdate.NewAgg
			if newAgg.IPPrefix != "" {
				s.AddOrUpdateAgg(oldAgg, newAgg, aggUpdate.AttrSet)
			}

		case aggConf := <-s.RemAggCh:
			s.DeleteAgg(aggConf)

		case nsUpdate := <-s.AddNSCh:
			oldNS := nsUpdate.OldNS
			newNS := nsUpdate.NewNS
			if newNS.IPPrefix != "" {
				s.AddOrUpdateNS(oldNS, newNS, nsUpdate.AttrSet)
			}

		case nsConf := <-s.RemNSCh:
			s.DeleteNS(&nsConf)

		case tcpConn := <-s.acceptCh:
			s.logger.Info("Connected to", tcpConn.RemoteAddr().String())
			host, _, _ := net.SplitHostPort(tcpConn.RemoteAddr().String())
			hostSplit := strings.Split(host, "%")
			host = hostSplit[0]
			peer, ok := s.PeerMap[host]
			if !ok {
				s.logger.Info("Can't accept connection. Peer is not configured yet", host)
				tcpConn.Close()
				s.logger.Info("Closed connection from", host)
				break
			}
			peer.AcceptConn(tcpConn)

		case peerCommand := <-s.PeerCommandCh:
			s.logger.Info("Peer Command received", peerCommand)
			peer, ok := s.PeerMap[peerCommand.IP.String()]
			if !ok {
				s.logger.Infof("Failed to apply command %s. Peer at that address does not exist, %v",
					peerCommand.Command, peerCommand.IP)
			}
			peer.Command(peerCommand.Command, fsm.BGPCmdReasonNone)

		case peerFSMConn := <-s.PeerFSMConnCh:
			s.logger.Infof("Server: Peer %s FSM established/broken channel", peerFSMConn.PeerIP)
			peer, ok := s.PeerMap[peerFSMConn.PeerIP]
			if !ok {
				s.logger.Infof("Failed to process FSM connection success, Peer %s does not exist",
					peerFSMConn.PeerIP)
				break
			}

			if peerFSMConn.Established {
				peer.PeerConnEstablished(peerFSMConn.Conn)
				addPathsMaxTx := peer.getAddPathsMaxTx()
				if addPathsMaxTx > s.AddPathCount {
					s.AddPathCount = addPathsMaxTx
				}
				s.setInterfaceMapForPeer(peerFSMConn.PeerIP, peer)
				s.SendAllRoutesToPeer(peer)
			} else {
				peer.PeerConnBroken(true)
				addPathsMaxTx := peer.getAddPathsMaxTx()
				if addPathsMaxTx < s.AddPathCount {
					s.AddPathCount = 0
					for _, otherPeer := range s.PeerMap {
						addPathsMaxTx = otherPeer.getAddPathsMaxTx()
						if addPathsMaxTx > s.AddPathCount {
							s.AddPathCount = addPathsMaxTx
						}
					}
				}
				s.clearInterfaceMapForPeer(peerFSMConn.PeerIP, peer)
				s.ProcessRemoveNeighbor(peerFSMConn.PeerIP, peer)
			}

		case peerIP := <-s.PeerConnEstCh:
			s.logger.Infof("Server: Peer %s FSM connection established", peerIP)
			peer, ok := s.PeerMap[peerIP]
			if !ok {
				s.logger.Infof("Failed to process FSM connection success, Peer %s does not exist", peerIP)
				break
			}
			reachInfo, err := s.routeMgr.GetNextHopInfo(peerIP, -1)
			if err != nil {
				s.logger.Infof("Server: Peer %s is not reachable", peerIP)
			} else {
				// @TODO: jgheewala think of something better for ovsdb....
				ifIdx := s.IntfMgr.GetIfIndex(int(reachInfo.NextHopIfIndex),
					int(reachInfo.NextHopIfType))
				s.logger.Infof("Server: Peer %s IfIdx %d", peerIP, ifIdx)
				if _, ok := s.IfIndexPeerMap[ifIdx]; !ok {
					s.IfIndexPeerMap[ifIdx] = make([]string, 0)
					//ifIdx := asicdCommonDefs.GetIfIndexFromIntfIdAndIntfType(int(reachInfo.NextHopIfIndex),
					//	int(reachInfo.NextHopIfType))
				}
				s.IfIndexPeerMap[ifIdx] = append(s.IfIndexPeerMap[ifIdx],
					peerIP)
				peer.setIfIdx(ifIdx)
			}

			s.SendAllRoutesToPeer(peer)

		case peerIP := <-s.PeerConnBrokenCh:
			s.logger.Infof("Server: Peer %s FSM connection broken", peerIP)
			peer, ok := s.PeerMap[peerIP]
			if !ok {
				s.logger.Infof("Failed to process FSM connection failure, Peer %s does not exist", peerIP)
				break
			}
			ifIdx := peer.getIfIdx()
			s.logger.Infof("Server: Peer %s FSM connection broken ifIdx %v", peerIP, ifIdx)
			if peerList, ok := s.IfIndexPeerMap[ifIdx]; ok {
				for idx, ip := range peerList {
					if ip == peerIP {
						s.IfIndexPeerMap[ifIdx] =
							append(s.IfIndexPeerMap[ifIdx][:idx],
								s.IfIndexPeerMap[ifIdx][idx+1:]...)
						if len(s.IfIndexPeerMap[ifIdx]) == 0 {
							delete(s.IfIndexPeerMap, ifIdx)
						}
						break
					}
				}
			}
			peer.setIfIdx(-1)
			s.ProcessRemoveNeighbor(peerIP, peer)

		case pktInfo := <-s.BGPPktSrcCh:
			s.logger.Info("Received BGP message from peer %s", pktInfo.Src)
			s.ProcessUpdate(pktInfo)

		case reachabilityInfo := <-s.ReachabilityCh:
			s.logger.Info("Server: Get reachability info for ip", reachabilityInfo.IP)

			nhInfo, err := s.routeMgr.GetNextHopInfo(reachabilityInfo.IP, reachabilityInfo.IfIndex)
			s.logger.Infof("Server: Reachability info for ip is %+v", nhInfo)
			reachabilityInfo.ReachableCh <- config.ReachabilityResult{Err: err, NextHopInfo: nhInfo}

		case bfdNotify := <-s.BfdCh:
			s.handleBfdNotifications(bfdNotify.Oper, bfdNotify.DestIp, bfdNotify.State)

		case ifState := <-s.IntfCh:
			s.logger.Info("Received message on ItfCh")
			if ifState.State == config.INTF_STATE_DOWN {
				if peerList, ok := s.IfIndexPeerMap[ifState.Idx]; ok {
					for _, peerIP := range peerList {
						if peer, ok := s.PeerMap[peerIP]; ok {
							peer.StopFSM("Interface Down")
						}
					}
				}
			} else if ifState.State == config.INTF_CREATED {
				s.ifaceMgr.AddIface(ifState.Idx, ifState.IPAddr)
				s.handleIntfCreate(ifState.Idx, config.PeerAddressV4)
			} else if ifState.State == config.INTF_DELETED {
				s.ifaceMgr.RemoveIface(ifState.Idx, ifState.IPAddr)
				s.handleIntfDelete(ifState.Idx, config.PeerAddressV4)
			} else if ifState.State == config.INTFV6_CREATED {
				s.ifaceMgr.AddV6Iface(ifState.Idx, ifState.IPAddr)
				s.handleIntfCreate(ifState.Idx, config.PeerAddressV6)
			} else if ifState.State == config.INTFV6_DELETED {
				s.ifaceMgr.RemoveV6Iface(ifState.Idx, ifState.IPAddr)
				s.handleIntfDelete(ifState.Idx, config.PeerAddressV6)
			} else if ifState.State == config.IPV6_NEIGHBOR_CREATED {
				s.logger.Info("IPV6_NEIGHBOR_CREATED message")
				s.ifaceMgr.AddLinkLocalIface(ifState.Idx, ifState.LinkLocalIP)
				s.handleIntfCreate(ifState.Idx, config.PeerAddressV6)
			} else if ifState.State == config.IPV6_NEIGHBOR_DELETED {
				s.ifaceMgr.RemoveLinkLocalIface(ifState.Idx, ifState.LinkLocalIP)
				s.handleIntfDelete(ifState.Idx, config.PeerAddressV6)
			}

		case ifMap := <-s.IntfMapCh:
			s.logger.Info("Received interface map")
			s.ProcessIntfMapUpdates([]config.IntfMapInfo{config.IntfMapInfo{Idx: ifMap.Idx, IfName: ifMap.IfName}})

		case routeInfo := <-s.RoutesCh:
			s.ProcessConnectedRoutes(routeInfo.Add, routeInfo.Remove)
		}
	}

}

func (s *BGPServer) ProcessIntfMapUpdates(cfg []config.IntfMapInfo) {
	s.logger.Infof("ProcessIntfMapUpdates, cfg = %+v", cfg)
	if s.IntfIdNameMap == nil {
		s.IntfIdNameMap = make(map[int32]IntfEntry)
	}
	if s.IfNameToIfIndex == nil {
		s.IfNameToIfIndex = make(map[string]int32)
	}
	for _, ifMap := range cfg {
		intfEntry := IntfEntry{Name: ifMap.IfName}
		s.IntfIdNameMap[int32(ifMap.Idx)] = intfEntry
		s.IfNameToIfIndex[ifMap.IfName] = ifMap.Idx
	}
}

func (s *BGPServer) InitBGPEvent() {
	// Start DB Util
	s.eventDbHdl = dbutils.NewDBUtil(s.logger)
	err := s.eventDbHdl.Connect()
	if err != nil {
		s.logger.Errf("DB connect failed with error %s. Exiting!!", err)
		return
	}
	err = eventUtils.InitEvents("BGPD", s.eventDbHdl, s.eventDbHdl, s.logger, 1000)
	if err != nil {
		s.logger.Err("Unable to initialize events", err)
	}
}

func (s *BGPServer) GetIntfObjects() {
	intfs := s.IntfMgr.GetIPv4Intfs()
	s.ProcessIntfStates(intfs)
	s.logger.Info("After ProcessIntfStates for intfs")

	v6intfs := s.IntfMgr.GetIPv6Intfs()
	s.ProcessIntfStates(v6intfs)
	s.logger.Info("After ProcessIntfStates for v6Intfs")

	v6Neighbors := s.IntfMgr.GetIPv6Neighbors()
	s.ProcessIntfStates(v6Neighbors)
	s.logger.Info("After ProcessIntfStates for v6Neighbors")

	portIntfMap := s.IntfMgr.GetPortInfo()
	s.ProcessIntfMapUpdates(portIntfMap)
	s.logger.Info("After ProcessIntfMapUpdates for ports")

	vlanIntfMap := s.IntfMgr.GetVlanInfo()
	s.ProcessIntfMapUpdates(vlanIntfMap)
	s.logger.Info("After ProcessIntfMapUpdates for vlans")

	logicalIntfMap := s.IntfMgr.GetLogicalIntfInfo()
	s.ProcessIntfMapUpdates(logicalIntfMap)
	s.logger.Info("After ProcessIntfMapUpdates for logicalIntfs")
}

func (s *BGPServer) ConstructPathAttrsForNSRoutes() []packet.BGPPathAttr {
	defaultPath, ok := s.NSRoutesPathMap[""]
	if !ok {
		defaultPath = s.ConstructPathsForNSRoutes(&s.BgpConfig.Global.Config)
	}
	newPathAttrs := packet.ClonePathAttrs(defaultPath.PathAttrs)
	return newPathAttrs
}

func (s *BGPServer) ConstructPathsForNSRoutes(gConf *config.GlobalConfig) *bgprib.Path {
	path, ok := s.NSRoutesPathMap[""]
	if !ok {
		pathAttrs := packet.ConstructPathAttrForDefaultRoute(gConf.AS, gConf.DefaultMED)
		protoFamily := packet.GetProtocolFamily(packet.AfiIP6, packet.SafiUnicast)
		ipv6MPReach := packet.ConstructIPv6MPReachNLRIForConnRoutes(protoFamily)
		path = bgprib.NewPath(s.LocRib, nil, pathAttrs, ipv6MPReach, bgprib.RouteTypeStatic)
		s.NSRoutesPathMap[""] = path
	}

	return path
}

func (s *BGPServer) ConstructPathsForLocalRoutes(gConf *config.GlobalConfig) {
	pathAttrs := packet.ConstructPathAttrForConnRoutes(gConf.AS, gConf.DefaultMED)
	protoFamily := packet.GetProtocolFamily(packet.AfiIP6, packet.SafiUnicast)
	ipv6MPReach := packet.ConstructIPv6MPReachNLRIForConnRoutes(protoFamily)
	s.ConnRoutesPath = bgprib.NewPath(s.LocRib, nil, pathAttrs, ipv6MPReach, bgprib.RouteTypeConnected)

	pathAttrs = packet.ConstructPathAttrForDefaultRoute(gConf.AS, gConf.DefaultMED)
	protoFamily = packet.GetProtocolFamily(packet.AfiIP6, packet.SafiUnicast)
	ipv6MPReach = packet.ConstructIPv6MPReachNLRIForConnRoutes(protoFamily)
	s.DefaultRoutesPath = bgprib.NewPath(s.LocRib, nil, pathAttrs, ipv6MPReach, bgprib.RouteTypeDefault)

	s.ConstructPathsForNSRoutes(gConf)
}

func (s *BGPServer) ConstructDefaultRoutes(gConf *config.GlobalConfig) {
	defaultRoutes := make(map[uint32]bool)
	if gConf.Defaultv4Route {
		defaultRoutes[packet.GetProtocolFamily(packet.AfiIP, packet.SafiUnicast)] = true
	}
	if gConf.Defaultv6Route {
		defaultRoutes[packet.GetProtocolFamily(packet.AfiIP6, packet.SafiUnicast)] = true
	}

	if len(defaultRoutes) > 0 {
		s.UpdateDefaultRoutes(defaultRoutes)
	}
}

func (s *BGPServer) StartServer() {
	// Initialize Event Handler
	s.InitBGPEvent()
	//read the intfMgr objects before the global conf - this is the case during restart
	s.GetIntfObjects()
	s.ServerUpCh <- true
	s.logger.Info("Setting serverup to true")

	globalUpdate := <-s.GlobalConfigCh
	gConf := globalUpdate.NewConfig
	s.GlobalCfgDone = true
	s.logger.Info("Recieved global conf:", gConf)
	s.BgpConfig.Global.Config = gConf
	s.constructBGPGlobalState(&gConf)
	s.BgpConfig.PeerGroups = make(map[uint32]map[string]*config.PeerGroup)

	s.ConstructPathsForLocalRoutes(&s.BgpConfig.Global.Config)
	s.ConstructDefaultRoutes(&s.BgpConfig.Global.Config)

	s.logger.Info("Setting up Peer connections")
	// channel for accepting connections
	s.acceptCh = make(chan *net.TCPConn)

	s.listener, _ = s.createListener("tcp4")
	go s.listenForPeers(s.listener, "tcp4", s.acceptCh)

	s.listenerIPv6, _ = s.createListener("tcp6")
	go s.listenForPeers(s.listenerIPv6, "tcp6", s.acceptCh)

	s.logger.Info("Start all managers and initialize API Layer")
	s.IntfMgr.Start()
	s.routeMgr.Start()
	s.bfdMgr.Start()
	s.SetupRedistribution(gConf)

	/*  ALERT: StartServer is a go routine and hence do not have any other go routine where
	 *	   you are making calls to other client. FlexSwitch uses thrift for rpc and hence
	 *	   on return it will not know which go routine initiated the thrift call.
	 */
	// Get routes from the route manager
	add, remove := s.routeMgr.GetRoutes()
	if add != nil && remove != nil {
		s.ProcessConnectedRoutes(add, remove)
	}
	s.GetIntfObjects()
	s.listenChannelUpdates()
}

func (s *BGPServer) GetBGPGlobalState() config.GlobalState {
	routesCount := s.LocRib.GetRoutesCount()
	s.BgpConfig.Global.State.Totalv4Prefixes = 0
	s.BgpConfig.Global.State.Totalv6Prefixes = 0
	for protoFamily, count := range routesCount {
		switch protoFamily {
		case packet.ProtocolFamilyMap["ipv4-unicast"]:
			s.BgpConfig.Global.State.Totalv4Prefixes = count

		case packet.ProtocolFamilyMap["ipv6-unicast"]:
			s.BgpConfig.Global.State.Totalv6Prefixes = count

		default:
			s.logger.Err("Unknown protocol family type", protoFamily)
		}
	}
	return s.BgpConfig.Global.State
}

func (s *BGPServer) GetBGPNeighborState(neighborIP string) *config.NeighborState {
	peer, ok := s.PeerMap[neighborIP]
	if !ok {
		s.logger.Errf("GetBGPNeighborState - Neighbor not found for address:%s", neighborIP)
		return nil
	}
	return &peer.NeighborConf.Neighbor.State
}

func (s *BGPServer) bulkGetBGPNeighbors(index int, count int, addrType config.PeerAddressType) (int, int,
	[]*config.NeighborState) {
	defer s.NeighborMutex.RUnlock()

	s.NeighborMutex.RLock()

	num := 0
	result := make([]*config.NeighborState, 0)
	for i := index; i < len(s.Neighbors); i++ {
		if s.Neighbors[i+index].NeighborConf.Neighbor.Config.PeerAddressType == addrType {
			num++
			if num <= count {
				result = append(result, &s.Neighbors[i+index].NeighborConf.Neighbor.State)
			} else {
				break
			}
		}
	}

	if num > count {
		index += count
	} else {
		index = 0
		count = num
	}
	return index, count, result
}

func (s *BGPServer) BulkGetBGPv4Neighbors(index int, count int) (int, int, []*config.NeighborState) {
	return s.bulkGetBGPNeighbors(index, count, config.PeerAddressV4)
}

func (s *BGPServer) BulkGetBGPv6Neighbors(index int, count int) (int, int, []*config.NeighborState) {
	return s.bulkGetBGPNeighbors(index, count, config.PeerAddressV6)
}

func (s *BGPServer) GetBGPNetworkStatementState(ipPrefix string) *config.BGPNetworkStatement {
	ip, _, err := net.ParseCIDR(ipPrefix)
	if err != nil {
		s.logger.Info("BGPNetworkStatement: IP", ipPrefix, "is not valid")
		return nil
	}

	var afi packet.AFI
	if ip.To4() != nil {
		afi = packet.AfiIP
	} else if ip.To16() != nil {
		afi = packet.AfiIP6
	} else {
		s.logger.Info("BGPNetworkStatement: IP", ipPrefix, "is not a valid v4 or v6 address")
		return nil
	}

	protoFamily := packet.GetProtocolFamily(afi, packet.SafiUnicast)
	if pfNSMap, ok := s.BgpConfig.Afs[protoFamily]; ok {
		if ns, ok := pfNSMap.BgpNetworkStmts[ipPrefix]; ok {
			return ns
		}
	}
	return nil
}

func (s *BGPServer) BulkGetBGPNetworkStatements(index int, count int) (int, int, []*config.BGPNetworkStatement) {
	s.NSRoutesMutex.RLock()
	defer s.NSRoutesMutex.RUnlock()

	num := count
	if index+count > len(s.NSRoutes) {
		num = len(s.NSRoutes) - index
	}

	result := make([]*config.BGPNetworkStatement, num)
	copy(result, s.NSRoutes[index:index+num])

	if index+count >= len(s.NSRoutes) {
		index = 0
		count = num
	} else {
		index += count
	}
	return index, count, result
}

func (s *BGPServer) VerifyBgpGlobalConfig() bool {
	return s.GlobalCfgDone
}

func (s *BGPServer) ConvertIntfStrToIfIndex(intfString string) (ifIndex int32, ifName string, err error) {
	if val, err := strconv.Atoi(intfString); err == nil {
		s.logger.Info("ConvertIntfStrToIfIndex - intfString", intfString, "is ifIndex")
		//Verify ifIndex is valid
		ifIndex = int32(val)
		s.logger.Info("IfIndex = ", val)
		ifEntry, ok := s.IntfIdNameMap[ifIndex]
		if !ok || ifEntry.Name == "" {
			s.logger.Errf("ConvertIntfStrToIfIndex - Did not find ifIndex %d in IntfIdNameMap map %+v", ifIndex,
				s.IntfIdNameMap)
			return ifIndex, ifName, errors.New(fmt.Sprintf("Did not ifIndex %d in interface map", ifIndex))
		}
		ifName = ifEntry.Name
		s.logger.Info("ConvertIntfStrToIfIndex - ifIndex =", ifIndex, "ifEntry =", ifEntry, "ifName =", ifName)
	} else {
		//Verify ifName is valid
		s.logger.Info("ConvertIntfStrToIfIndex - intfString", intfString, "is ifName")
		var ok bool
		if ifIndex, ok = s.IfNameToIfIndex[intfString]; !ok {
			s.logger.Errf("ConvertIntfStrToIfIndex - Did not find ifName %s in ifnametoifindex map %+v", intfString,
				s.IfNameToIfIndex)
			return ifIndex, ifName, errors.New(fmt.Sprintf("Invalid ifName %d", intfString))
		}
		ifName = intfString
		s.logger.Err("ConvertIntfStrToIfIndex - ifName =", ifName, "ifIndex =", ifIndex)
	}
	return ifIndex, ifName, nil
}
