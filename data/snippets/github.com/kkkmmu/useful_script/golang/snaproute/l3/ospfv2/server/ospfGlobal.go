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
	"errors"
	"l3/ospfv2/objects"
)

type GlobalStruct struct {
	Vrf                string
	RouterId           uint32
	AdminState         bool
	ASBdrRtrStatus     bool
	ReferenceBandwidth uint32
	AreaBdrRtrStatus   bool
	//isABR             bool
}

func genOspfv2GlobalUpdateMask(attrset []bool) uint32 {
	var mask uint32 = 0

	if attrset == nil {
		mask = objects.OSPFV2_GLOBAL_UPDATE_ROUTER_ID |
			objects.OSPFV2_GLOBAL_UPDATE_ADMIN_STATE |
			objects.OSPFV2_GLOBAL_UPDATE_AS_BDR_RTR_STATUS |
			objects.OSPFV2_GLOBAL_UPDATE_REFERENCE_BANDWIDTH
	} else {
		for idx, val := range attrset {
			if true == val {
				switch idx {
				case 0:
					// Vrf
				case 1:
					mask |= objects.OSPFV2_GLOBAL_UPDATE_ROUTER_ID
				case 2:
					mask |= objects.OSPFV2_GLOBAL_UPDATE_ADMIN_STATE
				case 3:
					mask |= objects.OSPFV2_GLOBAL_UPDATE_AS_BDR_RTR_STATUS
				case 4:
					mask |= objects.OSPFV2_GLOBAL_UPDATE_REFERENCE_BANDWIDTH
				}
			}
		}
	}
	return mask
}

func (server *OSPFV2Server) updateGlobal(newCfg, oldCfg *objects.Ospfv2Global, attrset []bool) (bool, error) {
	server.logger.Info("Global configuration update")
	if server.globalData.AdminState == true {
		server.StopAllIntfFSM()
		//Stop Rx Pkt
		//Deinit Rx Pkt
		//Deinit Tx Pkt
		// TODO: Stop Nbr FSM
		server.StopNbrFSM()

		// Deinit LSDB Data Structure
		// Stop LSDB
		server.StopLsdbRoutine()

		//TODO: Stop Flooding

		server.StopFlooding()
		// Flush all the routes
		// Deinit Routing Tbl
		// Deinit SPF Structures
		// Stop SPF
		server.StopSPF()

		//TODO: Stop Ribd updates
		err := server.deinitAsicdForRxMulticastPkt()
		if err != nil {
			server.logger.Err("Unable to initialize ASIC for recving Multicast Packets", err)
			return false, err
		}
	}

	mask := genOspfv2GlobalUpdateMask(attrset)
	if mask&objects.OSPFV2_GLOBAL_UPDATE_ADMIN_STATE == objects.OSPFV2_GLOBAL_UPDATE_ADMIN_STATE {
		server.globalData.AdminState = newCfg.AdminState
	}
	if mask&objects.OSPFV2_GLOBAL_UPDATE_ROUTER_ID == objects.OSPFV2_GLOBAL_UPDATE_ROUTER_ID {
		server.globalData.RouterId = newCfg.RouterId
	}
	if mask&objects.OSPFV2_GLOBAL_UPDATE_AS_BDR_RTR_STATUS == objects.OSPFV2_GLOBAL_UPDATE_AS_BDR_RTR_STATUS {
		server.globalData.ASBdrRtrStatus = newCfg.ASBdrRtrStatus
	}
	if mask&objects.OSPFV2_GLOBAL_UPDATE_REFERENCE_BANDWIDTH == objects.OSPFV2_GLOBAL_UPDATE_REFERENCE_BANDWIDTH {
		server.globalData.ReferenceBandwidth = newCfg.ReferenceBandwidth
	}

	if server.globalData.AdminState == true {
		err := server.initAsicdForRxMulticastPkt()
		if err != nil {
			server.logger.Err("Unable to initialize ASIC for recving Multicast Packets", err)
			return false, err
		}
		// Init SPF Structures
		// Init Routing Tbl
		// Start SPF
		server.StartSPF()
		server.logger.Info("Successfully started SPF")

		server.StartFlooding()
		server.logger.Info("Successfully started Flooding")
		// Init LSDB Data Structure
		// Start LSDB
		server.StartLsdbRoutine()
		server.logger.Info("Successfully started Lsdb Routine")
		server.StartNbrFSM()
		server.logger.Info("Successfully started Nbr FSM")

		// Init Rx
		// Init Tx
		// Start Rx
		// Init Ospf Intf FSM
		//Start All Interface FSM
		server.StartAllIntfFSM()
		server.logger.Info("Successfully started All Intf FSM")
	}

	return true, nil
}

func (server *OSPFV2Server) createGlobal(cfg *objects.Ospfv2Global) (bool, error) {
	server.logger.Info("Global configuration create")
	if cfg.Vrf != "default" {
		server.logger.Err("Vrf other than default is not supported")
		return false, errors.New("Vrf other than default is not supported")
	}
	server.globalData.Vrf = cfg.Vrf
	server.globalData.AdminState = cfg.AdminState
	server.globalData.RouterId = cfg.RouterId
	server.globalData.ASBdrRtrStatus = cfg.ASBdrRtrStatus
	server.globalData.ReferenceBandwidth = cfg.ReferenceBandwidth
	if server.globalData.AdminState == true {
		err := server.initAsicdForRxMulticastPkt()
		if err != nil {
			server.logger.Err("Unable to initialize ASIC for recving Multicast Packets", err)
			return false, err
		}
		server.StartSPF()
		server.logger.Info("Successfully started SPF")
		server.StartFlooding()
		server.logger.Info("Successfully started Flooding")
		server.StartLsdbRoutine()
		server.logger.Info("Successfully started Lsdb Routine")
		server.StartNbrFSM()
		server.logger.Info("Successfully started Nbr FSM")
		server.StartAllIntfFSM()
		server.logger.Info("Successfully started All Intf FSM")
	}
	return true, nil
}

func (server *OSPFV2Server) deleteGlobal(cfg *objects.Ospfv2Global) (bool, error) {
	server.logger.Info("Global configuration delete")
	server.logger.Err("Global Configuration delete not supported")
	return false, errors.New("Global Configuration delete not supported")
}

func (server *OSPFV2Server) getGlobalState(vrf string) (*objects.Ospfv2GlobalState, error) {
	var retObj objects.Ospfv2GlobalState
	retObj.Vrf = vrf
	retObj.AreaBdrRtrStatus = server.globalData.AreaBdrRtrStatus
	retObj.NumOfAreas = uint32(len(server.AreaConfMap))
	retObj.NumOfIntfs = uint32(len(server.IntfConfMap))
	numOfNbrs := 0
	for _, intfEnt := range server.IntfConfMap {
		numOfNbrs += len(intfEnt.NbrMap)
	}
	for _, lsdbEnt := range server.LsdbData.AreaLsdb {
		retObj.NumOfRouterLSA += uint32(len(lsdbEnt.RouterLsaMap))
		retObj.NumOfNetworkLSA += uint32(len(lsdbEnt.NetworkLsaMap))
		retObj.NumOfSummary3LSA += uint32(len(lsdbEnt.Summary3LsaMap))
		retObj.NumOfSummary4LSA += uint32(len(lsdbEnt.Summary4LsaMap))
		retObj.NumOfASExternalLSA += uint32(len(lsdbEnt.ASExternalLsaMap))
	}
	retObj.NumOfLSA = retObj.NumOfRouterLSA + retObj.NumOfNetworkLSA +
		retObj.NumOfSummary3LSA + retObj.NumOfSummary4LSA +
		retObj.NumOfASExternalLSA
	//TODO: num of routes
	return &retObj, nil
}

func (server *OSPFV2Server) getBulkGlobalState(fromIdx, cnt int) (*objects.Ospfv2GlobalStateGetInfo, error) {
	var retObj objects.Ospfv2GlobalStateGetInfo
	if fromIdx > 0 {
		return nil, errors.New("Invalid range.")
	}
	retObj.EndIdx = 1
	retObj.More = false
	retObj.Count = 1
	for idx := fromIdx; idx < retObj.EndIdx; idx++ {
		obj, err := server.getGlobalState("default")
		if err != nil {
			server.logger.Err("Error getting the Ospfv2GlobalState for vrf=default")
			return nil, err
		}
		retObj.List = append(retObj.List, obj)
	}
	return &retObj, nil
}
