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
	"errors"
	"fmt"
	"l3/ospfv2/objects"
	"net"
	"time"
)

type IntfConfKey struct {
	IpAddr  uint32
	IntfIdx uint32
}

type IntfConf struct {
	AdminState      bool
	AreaId          uint32
	Type            uint8 //Broadcast, P2P
	RtrPriority     uint8
	TransitDelay    uint16
	RetransInterval uint16
	HelloInterval   uint16
	RtrDeadInterval uint32
	Cost            uint32
	Mtu             uint32
	AuthType        uint16
	AuthKey         uint64

	DRIpAddr  uint32
	DRtrId    uint32
	BDRIpAddr uint32
	BDRtrId   uint32

	OperState           bool
	FSMState            uint8
	NumOfStateChange    uint32
	TimeOfStateChange   string
	FSMCtrlCh           chan bool
	FSMCtrlReplyCh      chan bool
	HelloIntervalTicker *time.Ticker
	WaitTimer           *time.Timer

	BackupSeenCh chan BackupSeenMsg
	NbrCreateCh  chan NbrCreateMsg
	NbrChangeCh  chan NbrChangeMsg
	//NbrStateChangeCh chan NbrStateChangeMsg

	NbrMap map[NbrConfKey]NbrData //Nbrs IP Address in case of Broadcast

	LsaCount  uint32
	IfName    string
	IpAddr    uint32
	IfMacAddr net.HardwareAddr
	IfType    uint32 //Loopback/Vlan/Lag/Port
	Netmask   uint32
	txHdl     IntfTxHandle
	rxHdl     IntfRxHandle
}

func getOspfv2IntfUpdateMask(attrset []bool) uint32 {
	var mask uint32 = 0

	if attrset == nil {
		mask = objects.OSPFV2_INTF_UPDATE_ADMIN_STATE |
			objects.OSPFV2_INTF_UPDATE_AREA_ID |
			objects.OSPFV2_INTF_UPDATE_TYPE |
			objects.OSPFV2_INTF_UPDATE_RTR_PRIORITY |
			objects.OSPFV2_INTF_UPDATE_TRANSIT_DELAY |
			objects.OSPFV2_INTF_UPDATE_RETRANS_INTERVAL |
			objects.OSPFV2_INTF_UPDATE_HELLO_INTERVAL |
			objects.OSPFV2_INTF_UPDATE_RTR_DEAD_INTERVAL |
			objects.OSPFV2_INTF_UPDATE_METRIC_VALUE
	} else {
		for idx, val := range attrset {
			if true == val {
				switch idx {
				case 0:
					// IPAddress
				case 1:
					//AddressLessIfIdx
				case 2:
					mask |= objects.OSPFV2_INTF_UPDATE_ADMIN_STATE
				case 3:
					mask |= objects.OSPFV2_INTF_UPDATE_AREA_ID
				case 4:
					mask |= objects.OSPFV2_INTF_UPDATE_TYPE
				case 5:
					mask |= objects.OSPFV2_INTF_UPDATE_RTR_PRIORITY
				case 6:
					mask |= objects.OSPFV2_INTF_UPDATE_TRANSIT_DELAY
				case 7:
					mask |= objects.OSPFV2_INTF_UPDATE_RETRANS_INTERVAL
				case 8:
					mask |= objects.OSPFV2_INTF_UPDATE_HELLO_INTERVAL
				case 9:
					mask |= objects.OSPFV2_INTF_UPDATE_RTR_DEAD_INTERVAL
				case 10:
					mask |= objects.OSPFV2_INTF_UPDATE_METRIC_VALUE
				}
			}
		}
	}
	return mask
}

func (server *OSPFV2Server) updateIntf(newCfg, oldCfg *objects.Ospfv2Intf, attrset []bool) (bool, error) {
	server.logger.Info("Intf configuration update")
	intfConfKey := IntfConfKey{
		IpAddr:  newCfg.IpAddress,
		IntfIdx: newCfg.AddressLessIfIdx,
	}
	intfConfEnt, exist := server.IntfConfMap[intfConfKey]
	if !exist {
		server.logger.Err("Ospf Interface configuration doesnot exist")
		return false, errors.New("Ospf Interface configuration doesnot exist")
	}
	areaEnt, _ := server.AreaConfMap[intfConfEnt.AreaId]
	if intfConfEnt.AdminState == true &&
		server.globalData.AdminState == true &&
		areaEnt.AdminState == true &&
		intfConfEnt.OperState == true {
		server.StopIntfFSM(intfConfKey)
	}
	intfConfEnt, _ = server.IntfConfMap[intfConfKey]
	oldIntfConfEnt := intfConfEnt
	mask := getOspfv2IntfUpdateMask(attrset)
	if mask&objects.OSPFV2_INTF_UPDATE_ADMIN_STATE == objects.OSPFV2_INTF_UPDATE_ADMIN_STATE {
		intfConfEnt.AdminState = newCfg.AdminState
	}
	if mask&objects.OSPFV2_INTF_UPDATE_AREA_ID == objects.OSPFV2_INTF_UPDATE_AREA_ID {
		_, exist := server.AreaConfMap[newCfg.AreaId]
		if !exist {
			server.logger.Err("Area doesnot exist")
			return false, errors.New("Area doesnot exist")
		}
		intfConfEnt.AreaId = newCfg.AreaId
	}
	if mask&objects.OSPFV2_INTF_UPDATE_TYPE == objects.OSPFV2_INTF_UPDATE_TYPE {
		intfConfEnt.Type = newCfg.Type
	}
	if mask&objects.OSPFV2_INTF_UPDATE_RTR_PRIORITY == objects.OSPFV2_INTF_UPDATE_RTR_PRIORITY {
		intfConfEnt.RtrPriority = newCfg.RtrPriority
	}
	if mask&objects.OSPFV2_INTF_UPDATE_TRANSIT_DELAY == objects.OSPFV2_INTF_UPDATE_TRANSIT_DELAY {
		intfConfEnt.TransitDelay = newCfg.TransitDelay
	}
	if mask&objects.OSPFV2_INTF_UPDATE_RETRANS_INTERVAL == objects.OSPFV2_INTF_UPDATE_RETRANS_INTERVAL {
		intfConfEnt.RetransInterval = newCfg.RetransInterval
	}
	if mask&objects.OSPFV2_INTF_UPDATE_HELLO_INTERVAL == objects.OSPFV2_INTF_UPDATE_HELLO_INTERVAL {
		intfConfEnt.HelloInterval = newCfg.HelloInterval
	}
	if mask&objects.OSPFV2_INTF_UPDATE_RTR_DEAD_INTERVAL == objects.OSPFV2_INTF_UPDATE_RTR_DEAD_INTERVAL {
		intfConfEnt.RtrDeadInterval = newCfg.RtrDeadInterval
	}
	if mask&objects.OSPFV2_INTF_UPDATE_METRIC_VALUE == objects.OSPFV2_INTF_UPDATE_METRIC_VALUE {
		intfConfEnt.Cost = uint32(newCfg.MetricValue)
	}
	areaEnt, _ = server.AreaConfMap[oldIntfConfEnt.AreaId]
	delete(areaEnt.IntfMap, intfConfKey)
	server.AreaConfMap[oldIntfConfEnt.AreaId] = areaEnt

	server.IntfConfMap[intfConfKey] = intfConfEnt

	//Add interface to new area
	areaEnt, _ = server.AreaConfMap[intfConfEnt.AreaId]
	areaEnt.IntfMap[intfConfKey] = true
	server.AreaConfMap[intfConfEnt.AreaId] = areaEnt

	if intfConfEnt.AdminState == true &&
		server.globalData.AdminState == true &&
		areaEnt.AdminState == true &&
		intfConfEnt.OperState == true {
		server.StartIntfFSM(intfConfKey)
		server.logger.Info("Started Intf FSM successfully", intfConfKey, intfConfEnt)
	}
	return true, nil
}

func (server *OSPFV2Server) createIntf(cfg *objects.Ospfv2Intf) (bool, error) {
	server.logger.Info("Intf configuration create")
	intfConfKey := IntfConfKey{
		IpAddr:  cfg.IpAddress,
		IntfIdx: cfg.AddressLessIfIdx,
	}

	intfConfEnt, exist := server.IntfConfMap[intfConfKey]
	if exist {
		server.logger.Err("Ospf Interface configuration already exist")
		return false, errors.New("Ospf Interface configuration already exist")
	}

	l3IfIdx, exist := server.infraData.ipToIfIdxMap[cfg.IpAddress]
	if !exist {
		// TODO: May be un numbered
		/*
			intfConfEnt.Mtu = uint32(1500) // Revisit
			intfConfEnt.IfName = ipEnt.IfName
			intfConfEnt.IfMacAddr = ipEnt.MacAddr
			intfConfEnt.Netmask = ipEnt.NetMask
		*/
		server.logger.Err("Unknown L3 Interface", cfg.IpAddress, cfg.AddressLessIfIdx)
		return false, errors.New("Unable to create Interface config: since no such L3 Interface exist")
	} else {
		ipEnt, _ := server.infraData.ipPropertyMap[l3IfIdx]
		intfConfEnt.OperState = ipEnt.State
		intfConfEnt.Mtu = uint32(ipEnt.Mtu)
		intfConfEnt.IfName = ipEnt.IfName
		intfConfEnt.IfMacAddr = ipEnt.MacAddr
		intfConfEnt.Netmask = ipEnt.NetMask
		intfConfEnt.IpAddr = ipEnt.IpAddr
		intfConfEnt.IfType = ipEnt.IfType
	}
	intfConfEnt.AdminState = cfg.AdminState
	areaEnt, exist := server.AreaConfMap[cfg.AreaId]
	if !exist {
		server.logger.Err("Area doesnot exist")
		return false, errors.New("Area doesnot exist")
	}
	intfConfEnt.AreaId = cfg.AreaId
	intfConfEnt.Type = cfg.Type
	intfConfEnt.RtrPriority = cfg.RtrPriority
	intfConfEnt.TransitDelay = cfg.TransitDelay
	intfConfEnt.RetransInterval = cfg.RetransInterval
	intfConfEnt.HelloInterval = cfg.HelloInterval
	intfConfEnt.RtrDeadInterval = cfg.RtrDeadInterval
	intfConfEnt.Cost = uint32(cfg.MetricValue)

	//intfConfEnt.DRIpAddr = 0
	//intfConfEnt.DRtrId = 0
	//intfConfEnt.BDRIpAddr = 0
	//intfConfEnt.BDRtrId = 0
	intfConfEnt.AuthKey = 0

	intfConfEnt.FSMState = objects.INTF_FSM_STATE_DOWN

	//intfConfEnt.FSMCtrlCh = make(chan bool)
	//intfConfEnt.FSMCtrlReplyCh = make(chan bool)
	//intfConfEnt.HelloIntervalTicker = nil
	//intfConfEnt.WaitTimer = nil

	//intfConfEnt.BackupSeenCh = make(chan BackupSeenMsg)
	//intfConfEnt.NbrCreateCh = make(chan NbrCreateMsg)
	//intfConfEnt.NbrChangeCh = make(chan NbrChangeMsg)

	intfConfEnt.LsaCount = 0
	server.IntfConfMap[intfConfKey] = intfConfEnt

	areaEnt.IntfMap[intfConfKey] = true
	server.MessagingChData.NbrToIntfFSMChData.NbrDownMsgChMap[intfConfKey] = make(chan NbrDownMsg)
	server.AreaConfMap[cfg.AreaId] = areaEnt
	if intfConfEnt.AdminState == true &&
		server.globalData.AdminState == true &&
		areaEnt.AdminState == true &&
		intfConfEnt.OperState == true {
		server.StartIntfFSM(intfConfKey)
	}
	//Adding Interface Conf Key To Slice
	server.GetBulkData.IntfConfSlice = append(server.GetBulkData.IntfConfSlice, intfConfKey)
	return true, nil
}

func (server *OSPFV2Server) deleteIntf(cfg *objects.Ospfv2Intf) (bool, error) {
	server.logger.Info("Intf configuration delete")
	intfConfKey := IntfConfKey{
		IpAddr:  cfg.IpAddress,
		IntfIdx: cfg.AddressLessIfIdx,
	}
	intfConfEnt, exist := server.IntfConfMap[intfConfKey]
	if !exist {
		server.logger.Err("Ospf Interface configuration doesnot exist")
		return false, errors.New("Ospf Interface configuration doesnot exist")
	}

	server.logger.Info("Intf Conf Ent", intfConfEnt)
	areaEnt, _ := server.AreaConfMap[intfConfEnt.AreaId]
	if intfConfEnt.AdminState == true &&
		server.globalData.AdminState == true &&
		areaEnt.AdminState == true &&
		intfConfEnt.OperState == true {
		server.StopIntfFSM(intfConfKey)
	}

	delete(areaEnt.IntfMap, intfConfKey)
	server.AreaConfMap[intfConfEnt.AreaId] = areaEnt
	delete(server.MessagingChData.NbrToIntfFSMChData.NbrDownMsgChMap, intfConfKey)
	delete(server.IntfConfMap, intfConfKey)
	return true, nil
}

func (server *OSPFV2Server) getIntfState(ipAddr, addressLessIfIdx uint32) (*objects.Ospfv2IntfState, error) {
	var retObj objects.Ospfv2IntfState
	server.logger.Info("ipAddr:", ipAddr, "addressLessIfIdx:", addressLessIfIdx, server.IntfConfMap)
	intfKey := IntfConfKey{
		IpAddr:  ipAddr,
		IntfIdx: addressLessIfIdx,
	}
	intfEnt, exist := server.IntfConfMap[intfKey]
	if !exist {
		server.logger.Err("Get Intf State: Interface does not exist", intfKey)
		return nil, errors.New("Interface does not exist")
	}
	retObj.IpAddress = ipAddr
	retObj.AddressLessIfIdx = addressLessIfIdx
	retObj.State = intfEnt.FSMState
	retObj.DesignatedRouterId = intfEnt.DRtrId
	retObj.BackupDesignatedRouterId = intfEnt.BDRtrId
	retObj.DesignatedRouter = intfEnt.DRIpAddr
	retObj.BackupDesignatedRouter = intfEnt.BDRIpAddr
	retObj.NumOfNbrs = uint32(len(intfEnt.NbrMap))
	retObj.Mtu = intfEnt.Mtu
	retObj.Cost = intfEnt.Cost
	lsdbKey := LsdbKey{
		AreaId: intfEnt.AreaId,
	}
	lsdbEnt, exist := server.LsdbData.AreaLsdb[lsdbKey]
	if exist {
		retObj.NumOfRouterLSA = uint32(len(lsdbEnt.RouterLsaMap))
		retObj.NumOfNetworkLSA = uint32(len(lsdbEnt.NetworkLsaMap))
		retObj.NumOfSummary3LSA = uint32(len(lsdbEnt.Summary3LsaMap))
		retObj.NumOfSummary4LSA = uint32(len(lsdbEnt.Summary4LsaMap))
		retObj.NumOfASExternalLSA = uint32(len(lsdbEnt.ASExternalLsaMap))
		retObj.NumOfLSA = retObj.NumOfRouterLSA + retObj.NumOfNetworkLSA +
			retObj.NumOfSummary3LSA + retObj.NumOfSummary4LSA +
			retObj.NumOfASExternalLSA
	}
	//TODO: NumOfRoutes
	retObj.NumOfStateChange = intfEnt.NumOfStateChange
	retObj.TimeOfStateChange = intfEnt.TimeOfStateChange
	return &retObj, nil
}

func (server *OSPFV2Server) getBulkIntfState(fromIdx, cnt int) (*objects.Ospfv2IntfStateGetInfo, error) {
	var retObj objects.Ospfv2IntfStateGetInfo
	server.logger.Info(server.IntfConfMap)
	//numIntfConf := len(server.IntfConfMap)
	count := 0
	idx := fromIdx
	sliceLen := len(server.GetBulkData.IntfConfSlice)
	if fromIdx >= sliceLen {
		return nil, errors.New("Invalid Range")
	}
	for count < cnt {
		if idx == sliceLen {
			break
		}
		intfKey := server.GetBulkData.IntfConfSlice[idx]
		intfEnt, exist := server.IntfConfMap[intfKey]
		if !exist {
			idx++
			continue
		}
		var obj objects.Ospfv2IntfState
		obj.IpAddress = intfKey.IpAddr
		obj.AddressLessIfIdx = intfKey.IntfIdx
		obj.State = intfEnt.FSMState
		obj.DesignatedRouterId = intfEnt.DRtrId
		obj.BackupDesignatedRouterId = intfEnt.BDRtrId
		obj.DesignatedRouter = intfEnt.DRIpAddr
		obj.BackupDesignatedRouter = intfEnt.BDRIpAddr
		obj.NumOfNbrs = uint32(len(intfEnt.NbrMap))
		obj.Mtu = intfEnt.Mtu
		obj.Cost = intfEnt.Cost
		lsdbKey := LsdbKey{
			AreaId: intfEnt.AreaId,
		}
		lsdbEnt, exist := server.LsdbData.AreaLsdb[lsdbKey]
		if exist {
			obj.NumOfRouterLSA = uint32(len(lsdbEnt.RouterLsaMap))
			obj.NumOfNetworkLSA = uint32(len(lsdbEnt.NetworkLsaMap))
			obj.NumOfSummary3LSA = uint32(len(lsdbEnt.Summary3LsaMap))
			obj.NumOfSummary4LSA = uint32(len(lsdbEnt.Summary4LsaMap))
			obj.NumOfASExternalLSA = uint32(len(lsdbEnt.ASExternalLsaMap))
			obj.NumOfLSA = obj.NumOfRouterLSA + obj.NumOfNetworkLSA +
				obj.NumOfSummary3LSA + obj.NumOfSummary4LSA +
				obj.NumOfASExternalLSA
		}
		//TODO: NumOfRoutes
		obj.NumOfStateChange = intfEnt.NumOfStateChange
		obj.TimeOfStateChange = intfEnt.TimeOfStateChange
		retObj.List = append(retObj.List, &obj)
		count++
		idx++
	}

	retObj.EndIdx = idx
	if retObj.EndIdx == sliceLen {
		retObj.More = false
		retObj.Count = 0
	} else {
		retObj.More = true
		retObj.Count = sliceLen - retObj.EndIdx + 1
	}
	return &retObj, nil
}

func (server *OSPFV2Server) GetIntfConfForGivenIntfKey(intfConfKey IntfConfKey) (IntfConf, error) {
	intfConfEnt, exist := server.IntfConfMap[intfConfKey]
	if !exist {
		return intfConfEnt, errors.New(fmt.Sprintln("Error: Unable to get interface config for", intfConfKey))
	}
	return intfConfEnt, nil
}

/*
func (server *OSPFV2Server) GetIntfNbrList(intfEnt IntfConf) (nbrList []NbrConfKey) {
	for nbrKey, _ := range intfEnt.NbrMap {
		nbrList = append(nbrList, nbrKey)
	}
	return nbrList
}

*/

func (server *OSPFV2Server) RefreshIntfConfSlice() {
	if len(server.GetBulkData.AreaConfSlice) == 0 {
		return
	}
	server.GetBulkData.IntfConfSlice = server.GetBulkData.IntfConfSlice[:len(server.GetBulkData.AreaConfSlice)-1]
	server.GetBulkData.IntfConfSlice = nil
	for intfKey, _ := range server.IntfConfMap {
		server.GetBulkData.IntfConfSlice = append(server.GetBulkData.IntfConfSlice, intfKey)
	}
}
