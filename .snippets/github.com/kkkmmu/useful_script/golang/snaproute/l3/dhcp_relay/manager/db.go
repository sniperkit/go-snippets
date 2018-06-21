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

package manager

import (
	"dhcprelayd"
	"models/objects"
)

func convertDRAv4IntfObjToThriftType(obj *objects.DHCPRelayIntf) *dhcprelayd.DHCPRelayIntf {
	thriftObj := &dhcprelayd.DHCPRelayIntf{
		IntfRef:  obj.IntfRef,
		Enable:   obj.Enable,
		ServerIp: obj.ServerIp,
	}
	return thriftObj
}

func convertDRAv4GlobalObjToThriftType(obj *objects.DHCPRelayGlobal) *dhcprelayd.DHCPRelayGlobal {
	thriftObj := &dhcprelayd.DHCPRelayGlobal{
		Vrf:           obj.Vrf,
		Enable:        obj.Enable,
		HopCountLimit: obj.HopCountLimit,
	}
	return thriftObj
}

func convertDRAv6IntfObjToThriftType(obj *objects.DHCPv6RelayIntf) *dhcprelayd.DHCPv6RelayIntf {
	thriftObj := &dhcprelayd.DHCPv6RelayIntf{
		IntfRef:  obj.IntfRef,
		Enable:   obj.Enable,
		ServerIp: obj.ServerIp,
	}
	return thriftObj
}

func convertDRAv6GlobalObjToThriftType(obj *objects.DHCPv6RelayGlobal) *dhcprelayd.DHCPv6RelayGlobal {
	thriftObj := &dhcprelayd.DHCPv6RelayGlobal{
		Vrf:           obj.Vrf,
		Enable:        obj.Enable,
		HopCountLimit: obj.HopCountLimit,
	}
	return thriftObj
}

func (draMgr *DRAMgr) readDRAv4GlobalConfig() (*dhcprelayd.DHCPRelayGlobal, error) {
	draMgr.Logger.Info("Reading DRAv4Global from db")
	var dbObj objects.DHCPRelayGlobal
	objList, err := draMgr.DbHdl.GetAllObjFromDb(dbObj)
	if err != nil {
		return nil, err
	}
	if len(objList) == 0 {
		return nil, nil
	}
	draMgr.Logger.Info("Objects from db are", objList)
	dbEntry := objList[0].(objects.DHCPRelayGlobal)
	thriftObj := convertDRAv4GlobalObjToThriftType(&dbEntry)
	return thriftObj, nil
}

func (draMgr *DRAMgr) readDRAv4IntfConfig() ([]*dhcprelayd.DHCPRelayIntf, error) {
	draMgr.Logger.Info("Reading DRAv4Intf from db")
	var dbObj objects.DHCPRelayIntf
	result := []*dhcprelayd.DHCPRelayIntf{}
	objList, err := draMgr.DbHdl.GetAllObjFromDb(dbObj)
	if err != nil {
		return nil, err
	}
	draMgr.Logger.Info("Objects from db are", objList)
	for _, obj := range objList {
		dbEntry := obj.(objects.DHCPRelayIntf)
		// (TODO): Remove this later and store as model object type
		// instead of thrift type
		thriftObj := convertDRAv4IntfObjToThriftType(&dbEntry)
		result = append(result, thriftObj)
	}
	return result, nil
}

func (draMgr *DRAMgr) readDRAv6GlobalConfig() (*dhcprelayd.DHCPv6RelayGlobal, error) {
	draMgr.Logger.Info("Reading DRAv6Global from db")
	var dbObj objects.DHCPv6RelayGlobal
	objList, err := draMgr.DbHdl.GetAllObjFromDb(dbObj)
	if err != nil {
		return nil, err
	}
	if len(objList) == 0 {
		return nil, nil
	}
	draMgr.Logger.Info("Objects from db are", objList)
	dbEntry := objList[0].(objects.DHCPv6RelayGlobal)
	thriftObj := convertDRAv6GlobalObjToThriftType(&dbEntry)
	return thriftObj, nil
}

func (draMgr *DRAMgr) readDRAv6IntfConfig() ([]*dhcprelayd.DHCPv6RelayIntf, error) {
	draMgr.Logger.Info("Reading DRAv6Intf from db")
	var dbObj objects.DHCPv6RelayIntf
	result := []*dhcprelayd.DHCPv6RelayIntf{}
	objList, err := draMgr.DbHdl.GetAllObjFromDb(dbObj)
	if err != nil {
		return nil, err
	}
	draMgr.Logger.Info("Objects from db are", objList)
	for _, obj := range objList {
		dbEntry := obj.(objects.DHCPv6RelayIntf)
		// (TODO): Remove this later and store as model object type
		// instead of thrift type
		thriftObj := convertDRAv6IntfObjToThriftType(&dbEntry)
		result = append(result, thriftObj)
	}
	return result, nil
}
