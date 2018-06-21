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

// portmap
package utils

import (
	"asicd/asicdCommonDefs"
	"fmt"
	"models/objects"
	"net"
	"utils/dbutils"
)

var PortConfigMap map[int32]PortConfig
var AggConfigMap map[int32]string

type PortConfig struct {
	Name         string
	HardwareAddr net.HardwareAddr
	Speed        int32
	IfIndex      int32
	Mtu          int32
	Duplex       string
}

func ConstructPortConfigMap() {
	currMarker := int(asicdCommonDefs.MIN_SYS_PORTS)
	count := 100
	for _, client := range GetAsicDPluginList() {
		GlobalLogger.Info("Calling asicd for port config")
		for {
			bulkInfo, err := client.GetBulkPortState(currMarker, count)
			if err != nil {
				GlobalLogger.Err(fmt.Sprintf("GetBulkPortState Error: %s", err))
				return
			}
			GlobalLogger.Info(fmt.Sprintf("Length of GetBulkPortState: %d", bulkInfo.Count))

			bulkCfgInfo, err := client.GetBulkPort(currMarker, count)
			if err != nil {
				GlobalLogger.Err(fmt.Sprintf("Error: %s", err))
				return
			}

			GlobalLogger.Info(fmt.Sprintf("Length of GetBulkPortConfig: %d", bulkCfgInfo.Count))
			objCount := int(bulkInfo.Count)
			more := bool(bulkInfo.More)
			currMarker = int(bulkInfo.EndIdx)
			for i := 0; i < objCount; i++ {
				ifindex := bulkInfo.PortStateList[i].IfIndex
				ent := PortConfigMap[ifindex]
				ent.IfIndex = ifindex
				ent.Name = bulkInfo.PortStateList[i].Name
				ent.HardwareAddr, _ = net.ParseMAC(bulkCfgInfo.PortList[i].MacAddr)
				PortConfigMap[ifindex] = ent
				GlobalLogger.Info(fmt.Sprintf("Found Port IfIndex %d Name %s\n", ent.IfIndex, ent.Name))
			}
			if !more {
				return
			}
		}
	}

	// lets read from db the rest of the info from db
	// MTU/Duplex/Speed
	dbHdl := dbutils.NewDBUtil(GetLaLogger())
	err := dbHdl.Connect()
	if err != nil {
		GlobalLogger.Info(fmt.Sprintf("Failed to open connection to read Port Info from the DB with error %s", err))
		return
	}
	defer dbHdl.Disconnect()

	var dbObj objects.Port
	objList, err := dbObj.GetAllObjFromDb(dbHdl)
	if err != nil {
		fmt.Println(fmt.Sprintf("DB Query failed when retrieving Port objects", err))
		return
	}
	for idx := 0; idx < len(objList); idx++ {
		dbObject := objList[idx].(objects.Port)
		ifindex := dbObject.IfIndex
		ent := PortConfigMap[ifindex]
		ent.Mtu = dbObject.Mtu
		ent.Duplex = dbObject.Duplex
		ent.Speed = dbObject.Speed
		PortConfigMap[ifindex] = ent
	}
}

func AddAggConfigMap(ifindex int32, intfref string) {
	if _, ok := AggConfigMap[ifindex]; !ok {
		AggConfigMap[ifindex] = intfref
	}
}

func DelAggConfigMap(ifindex int32, intfref string) {
	if _, ok := AggConfigMap[ifindex]; ok {
		delete(AggConfigMap, ifindex)
	}
}

func GetAggIfIndexFromName(name string) int32 {
	for ifindex, intfref := range AggConfigMap {
		if name == intfref {
			return ifindex
		}
	}
	return 0
}

func GetAggNameFromIfIndex(ifindex int32) string {
	if intfref, ok := AggConfigMap[ifindex]; ok {
		return intfref
	}
	return ""
}

func GetIfIndexFromName(name string) int32 {

	for _, portcfg := range PortConfigMap {
		if portcfg.Name == name {
			return portcfg.IfIndex
		}
	}
	return 0
}

func GetNameFromIfIndex(ifindex int32) string {
	for _, portcfg := range PortConfigMap {
		if portcfg.IfIndex == ifindex {
			return portcfg.Name
		}
	}
	return ""
}
