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
	"l3/ospf/config"
)

type GlobalConf struct {
	Vrf                      string
	RouterId                 []byte
	AdminStat                config.Status
	ASBdrRtrStatus           bool
	TOSSupport               bool
	ExtLsdbLimit             int32
	MulticastExtensions      int32
	ExitOverflowInterval     config.PositiveInteger
	DemandExtensions         bool
	RFC1583Compatibility     bool
	ReferenceBandwidth       uint32
	RestartSupport           config.RestartSupport
	RestartInterval          int32
	RestartStrictLsaChecking bool
	StubRouterAdvertisement  config.AdvertiseAction
	Version                  uint8
	AreaBdrRtrStatus         bool
	ExternLsaCount           int32
	ExternLsaChecksum        int32
	OriginateNewLsas         int32
	RxNewLsas                int32
	OpaqueLsaSupport         bool
	RestartStatus            config.RestartStatus
	RestartAge               int32
	RestartExitReason        config.RestartExitReason
	AsLsaCount               int32
	AsLsaCksumSum            int32
	StubRouterSupport        bool
	//DiscontinuityTime        string
	DiscontinuityTime int32 // This should be string
	isABR             bool
}

func (server *OSPFServer) updateGlobalConf(gConf config.GlobalConf) {
	routerId := convertAreaOrRouterId(string(gConf.RouterId))
	if routerId == nil {
		server.logger.Err("Invalid Router Id")
		return
	}
	server.ospfGlobalConf.Vrf = gConf.Vrf
	server.ospfGlobalConf.RouterId = routerId
	server.ospfGlobalConf.AdminStat = gConf.AdminStat
	server.ospfGlobalConf.ASBdrRtrStatus = gConf.ASBdrRtrStatus
	server.ospfGlobalConf.TOSSupport = gConf.TOSSupport
	server.ospfGlobalConf.RestartSupport = gConf.RestartSupport
	server.ospfGlobalConf.RestartInterval = gConf.RestartInterval
	server.ospfGlobalConf.ReferenceBandwidth = uint32(gConf.ReferenceBandwidth)
	server.logger.Err("Global configuration updated")
}

func (server *OSPFServer) initOspfGlobalConfDefault() {
	routerId := convertAreaOrRouterId("0.0.0.0")
	if routerId == nil {
		server.logger.Err("Invalid Router Id")
		return
	}
	server.ospfGlobalConf.RouterId = routerId
	server.ospfGlobalConf.AdminStat = config.Disabled
	server.ospfGlobalConf.ASBdrRtrStatus = false
	server.ospfGlobalConf.TOSSupport = false
	server.ospfGlobalConf.ExtLsdbLimit = -1
	server.ospfGlobalConf.MulticastExtensions = 0
	server.ospfGlobalConf.ExitOverflowInterval = 0
	server.ospfGlobalConf.RFC1583Compatibility = false
	server.ospfGlobalConf.ReferenceBandwidth = 100000 // Default value 100 Gbps
	server.ospfGlobalConf.RestartSupport = config.None
	server.ospfGlobalConf.RestartInterval = 0
	server.ospfGlobalConf.RestartStrictLsaChecking = false
	server.ospfGlobalConf.StubRouterAdvertisement = config.DoNotAdvertise
	server.ospfGlobalConf.Version = uint8(OSPF_VERSION_2)
	server.ospfGlobalConf.AreaBdrRtrStatus = false
	server.ospfGlobalConf.ExternLsaCount = 0
	server.ospfGlobalConf.ExternLsaChecksum = 0
	server.ospfGlobalConf.OriginateNewLsas = 0
	server.ospfGlobalConf.RxNewLsas = 0
	server.ospfGlobalConf.OpaqueLsaSupport = false
	server.ospfGlobalConf.RestartStatus = config.NotRestarting
	server.ospfGlobalConf.RestartAge = 0
	server.ospfGlobalConf.RestartExitReason = config.NoAttempt
	server.ospfGlobalConf.AsLsaCount = 0
	server.ospfGlobalConf.AsLsaCksumSum = 0
	server.ospfGlobalConf.StubRouterSupport = false
	//server.ospfGlobalConf.DiscontinuityTime = "0"
	server.ospfGlobalConf.DiscontinuityTime = 0 //This should be string
	server.ospfGlobalConf.isABR = false
	server.logger.Err("Global configuration initialized")
}

func (server *OSPFServer) processASBdrRtrStatus(isASBR bool) {
	if isASBR {
		server.logger.Info(fmt.Sprintln("GLOBAL: Router is ASBR. Listen to RIBD updates."))
		//get ribd routes
		//server.testASExternal()
		server.getRibdRoutes()
		server.startRibdUpdates()
	}
}
func (server *OSPFServer) processGlobalConfig(gConf config.GlobalConf) error {
	var localIntfStateMap = make(map[IntfConfKey]config.Status)
	for key, ent := range server.IntfConfMap {
		localIntfStateMap[key] = ent.IfAdminStat
		if ent.IfAdminStat == config.Enabled &&
			server.ospfGlobalConf.AdminStat == config.Enabled {
			server.StopSendRecvPkts(key)
		}
	}

	if server.ospfGlobalConf.AdminStat == config.Enabled {
		server.nbrFSMCtrlCh <- false
		//	server.neighborConfStopCh <- true
		//server.NeighborListMap = nil
		server.StopLSDatabase()
		server.ospfRxNbrPktStopCh <- true
		server.ospfTxNbrPktStopCh <- true
		server.neighborFSMCtrlCh <- false
		server.neighborConfStopCh <- true

	}
	server.logger.Info(fmt.Sprintln("Received call for performing Global Configuration", gConf))
	server.updateGlobalConf(gConf)

	if server.ospfGlobalConf.AdminStat == config.Enabled {
		//server.NeighborListMap = make(map[IntfConfKey]list.List)
		server.logger.Info(fmt.Sprintln("Spawn Neighbor state machine"))
		server.InitNeighborStateMachine()
		go server.UpdateNeighborConf()
		go server.ProcessNbrStateMachine()
		go server.ProcessTxNbrPkt()
		go server.ProcessRxNbrPkt()
		server.StartLSDatabase()

	}
	server.processASBdrRtrStatus(server.ospfGlobalConf.AreaBdrRtrStatus)
	for key, ent := range localIntfStateMap {
		if ent == config.Enabled &&
			server.ospfGlobalConf.AdminStat == config.Enabled {
			server.logger.Info(fmt.Sprintln("Start rx/tx thread."))
			server.StartSendRecvPkts(key)
		}
	}

	return nil
}
