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
	"models/objects"
	"ospfv2d"
	"utils/typeConv"
)

func (server *OSPFV2Server) StartDBClient() {
	for {
		select {
		case msg := <-server.MessagingChData.RouteTblToDBClntChData.RouteAddMsgCh:
			server.RouteAddToDB(msg)
		case msg := <-server.MessagingChData.RouteTblToDBClntChData.RouteDelMsgCh:
			server.RouteDelFromDB(msg)
		case _ = <-server.MessagingChData.ServerToDBClntChData.FlushRouteFromDBCh:
			server.FlushAllRoutesFromDB()
			server.SendFlushRouteDoneMsgToServer()
		}
	}
}

func (server *OSPFV2Server) FlushAllRoutesFromDB() {
	ospfv2RouteKeys, err := server.dbHdl.GetAllKeys("Ospfv2RouteState#*")
	if err != nil {
		return
	}
	keys, err := typeConv.ConvertToStrings(ospfv2RouteKeys, nil)
	if err != nil {
		return
	}
	for _, ospfv2RouteKey := range keys {
		err := server.dbHdl.DeleteValFromDb(ospfv2RouteKey)
		if err != nil {
			server.logger.Err("Error Deleting Route from DB", ospfv2RouteKey)
		}
	}
	return
}

func (server *OSPFV2Server) RouteAddToDB(msg RouteAddMsg) {
	var dbObj objects.Ospfv2RouteState
	obj := ospfv2d.NewOspfv2RouteState()

	obj.DestId = convertUint32ToDotNotation(msg.RTblKey.DestId)
	obj.AddrMask = convertUint32ToDotNotation(msg.RTblKey.AddrMask)
	switch msg.RTblKey.DestType {
	case Network:
		obj.DestType = "Network"
	case InternalRouter:
		obj.DestType = "Internal Router"
	case AreaBdrRouter:
		obj.DestType = "Area Border Router"
	case ASBdrRouter:
		obj.DestType = "AS Border Router"
	case ASAreaBdrRouter:
		obj.DestType = "AS and Area Border Router"
	default:
		obj.DestType = "Invalid"
	}
	obj.OptCapabilities = int32(msg.RTblEntry.RoutingTblEnt.OptCapabilities)
	obj.AreaId = convertUint32ToDotNotation(msg.RTblEntry.AreaId)
	switch msg.RTblEntry.RoutingTblEnt.PathType {
	case IntraArea:
		obj.PathType = "Intra Area"
	case InterArea:
		obj.PathType = "Inter Area"
	case Type1Ext:
		obj.PathType = "Type1 External"
	case Type2Ext:
		obj.PathType = "Type2 External"
	default:
		obj.PathType = "Invalid"
	}
	obj.Cost = int32(msg.RTblEntry.RoutingTblEnt.Cost)
	obj.Type2Cost = int32(msg.RTblEntry.RoutingTblEnt.Type2Cost)
	obj.NumOfPaths = int16(msg.RTblEntry.RoutingTblEnt.NumOfPaths)
	nh_list := make([]ospfv2d.Ospfv2NextHop, len(msg.RTblEntry.RoutingTblEnt.NextHops))
	idx := 0
	for nxtHop, _ := range msg.RTblEntry.RoutingTblEnt.NextHops {
		nh_list[idx].IntfIPAddr = convertUint32ToDotNotation(nxtHop.IfIPAddr)
		nh_list[idx].IntfIdx = int32(nxtHop.IfIdx)
		nh_list[idx].NextHopIPAddr = convertUint32ToDotNotation(nxtHop.NextHopIP)
		nh_list[idx].AdvRtrId = convertUint32ToDotNotation(nxtHop.AdvRtr)
		obj.NextHops = append(obj.NextHops, &nh_list[idx])
		idx++
	}
	if msg.RTblEntry.RoutingTblEnt.PathType == IntraArea {
		switch msg.RTblEntry.RoutingTblEnt.LSOrigin.LSType {
		case RouterLSA:
			obj.LSOriginLSType = "Router LSA"
		case NetworkLSA:
			obj.LSOriginLSType = "Network LSA"
		}
		obj.LSOriginLSId = convertUint32ToDotNotation(msg.RTblEntry.RoutingTblEnt.LSOrigin.LSId)
		obj.LSOriginAdvRouter = convertUint32ToDotNotation(msg.RTblEntry.RoutingTblEnt.LSOrigin.AdvRouter)
	}
	objects.ConvertThriftToospfv2dOspfv2RouteStateObj(obj, &dbObj)
	if server.dbHdl == nil {
		server.logger.Err("Db Handler is nil")
		return
	}
	err := server.dbHdl.StoreObjectInDb(dbObj)
	if err != nil {
		server.logger.Err(fmt.Sprintln("Failed to add route in db:", err))
	}
	return
}

func (server *OSPFV2Server) RouteDelFromDB(msg RouteDelMsg) {
	var dbObj objects.Ospfv2RouteState
	obj := ospfv2d.NewOspfv2RouteState()

	obj.DestId = convertUint32ToDotNotation(msg.RTblKey.DestId)
	obj.AddrMask = convertUint32ToDotNotation(msg.RTblKey.AddrMask)
	switch msg.RTblKey.DestType {
	case Network:
		obj.DestType = "Network"
	case InternalRouter:
		obj.DestType = "Internal Router"
	case AreaBdrRouter:
		obj.DestType = "Area Border Router"
	case ASBdrRouter:
		obj.DestType = "AS Border Router"
	case ASAreaBdrRouter:
		obj.DestType = "AS and Area Border Router"
	default:
		obj.DestType = "Invalid"
	}

	objects.ConvertThriftToospfv2dOspfv2RouteStateObj(obj, &dbObj)
	if server.dbHdl == nil {
		server.logger.Err("Db Handler is nil")
		return
	}
	err := server.dbHdl.DeleteObjectFromDb(dbObj)
	if err != nil {
		server.logger.Err(fmt.Sprintln("Failed to delete route in db:", err))
	}
	return
}
