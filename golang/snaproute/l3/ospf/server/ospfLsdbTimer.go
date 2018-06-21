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
	"encoding/binary"
	"fmt"
	"l3/ospf/config"

	"time"
)

/*@fn
LSDB aging API.
*/
func (server *OSPFServer) lsdbStateRefresh() {
	var lsdbStateRefFunc func()
	lsdbStateRefFunc = func() {
		//server.logger.Info(fmt.Sprintln("The old Lsdb Slice after refresh", server.LsdbSlice))
		server.LsdbSlice = []LsdbSliceEnt{}
		for lsdbkey, lsdbEnt := range server.AreaLsdb {
			for lsakey, _ := range lsdbEnt.RouterLsaMap {
				var val LsdbSliceEnt
				val.AreaId = lsdbkey.AreaId
				val.LSType = lsakey.LSType
				val.LSId = lsakey.LSId
				val.AdvRtr = lsakey.AdvRouter
				server.LsdbSlice = append(server.LsdbSlice, val)
			}
			for lsakey, _ := range lsdbEnt.NetworkLsaMap {
				var val LsdbSliceEnt
				val.AreaId = lsdbkey.AreaId
				val.LSType = lsakey.LSType
				val.LSId = lsakey.LSId
				val.AdvRtr = lsakey.AdvRouter
				server.LsdbSlice = append(server.LsdbSlice, val)
			}
			for lsakey, _ := range lsdbEnt.Summary3LsaMap {
				var val LsdbSliceEnt
				val.AreaId = lsdbkey.AreaId
				val.LSType = lsakey.LSType
				val.LSId = lsakey.LSId
				val.AdvRtr = lsakey.AdvRouter
				server.LsdbSlice = append(server.LsdbSlice, val)
			}
			for lsakey, _ := range lsdbEnt.Summary4LsaMap {
				var val LsdbSliceEnt
				val.AreaId = lsdbkey.AreaId
				val.LSType = lsakey.LSType
				val.LSId = lsakey.LSId
				val.AdvRtr = lsakey.AdvRouter
				server.LsdbSlice = append(server.LsdbSlice, val)
			}
			for lsakey, _ := range lsdbEnt.ASExternalLsaMap {
				var val LsdbSliceEnt
				val.AreaId = lsdbkey.AreaId
				val.LSType = lsakey.LSType
				val.LSId = lsakey.LSId
				val.AdvRtr = lsakey.AdvRouter
				server.LsdbSlice = append(server.LsdbSlice, val)
			}
		}
		server.LsdbStateTimer.Reset(server.RefreshDuration)
	}
	server.LsdbStateTimer = time.AfterFunc(server.RefreshDuration, lsdbStateRefFunc)
}

/*
@fn lsdbSelfLsaRefresh
 This API will refresh self generated LSAs after every
LSARefreshTime .
From RFC 2328
      Whenever a new instance of an LSA is originated, its LS sequence
        number is incremented, its LS age is set to 0, its LS checksum
        is calculated, and the LSA is added to the link state database
        and flooded out the appropriate interfaces.  See Section 13.2
        for details concerning the installation of the LSA into the link
        state database.
*/
func (server *OSPFServer) lsdbSelfLsaRefresh() {
	server.logger.Info(fmt.Sprintln("REFRESH: LSDB refresh started..."))
	floodAsExt := 0
	ifkey := IntfConfKey{}
	nbr := NeighborConfKey{}
	//get areaId for self originated LSAs
	for lsdbKey, selfOrigLsaEnt := range server.AreaSelfOrigLsa {
		for lsaKey, valid := range selfOrigLsaEnt {
			if valid {
				err := server.regenerateLsa(lsdbKey, lsaKey)
				if floodAsExt == 0 && lsaKey.LSType == ASExternalLSA {
					server.sendLsdbToNeighborEvent(ifkey, nbr, 0, 0, 0, lsaKey, LSAEXTFLOOD)
				}
				if err != nil {
					server.logger.Warning(fmt.Sprintln("LSDB: Failed to regenerate LSA ", lsaKey, " Area ", lsdbKey))
				}
			}
		}
		floodAsExt += 1
	}
	// generate Summary LSAs
	server.GenerateSummaryLsa()
	lsaKey := LsaKey{}

	for entKey, ent := range server.IntfConfMap {
		areaid := convertIPv4ToUint32(ent.IfAreaId)
		server.sendLsdbToNeighborEvent(entKey, nbr, areaid, 0, 0, lsaKey, LSAFLOOD)
	}
}

func (server *OSPFServer) regenerateLsa(lsdbKey LsdbKey, lsaKey LsaKey) error {

	switch lsaKey.LSType {
	case RouterLSA:
		server.generateRouterLSA(lsdbKey.AreaId)

	case NetworkLSA:
		isDR := false
		ipAddr := config.IpAddress(convertUint32ToIPv4(lsaKey.LSId))
		ifIdx := 0
		rtr_id := binary.BigEndian.Uint32(server.ospfGlobalConf.RouterId)
		intfKey := IntfConfKey{
			IPAddr:  ipAddr,
			IntfIdx: config.InterfaceIndexOrZero(ifIdx),
		}
		intConf, exist := server.IntfConfMap[intfKey]
		if !exist {
			return nil
		}
		if intConf.IfDRtrId == rtr_id {
			isDR = true
		}

		server.generateNetworkLSA(lsdbKey.AreaId, intfKey, isDR)

	case ASExternalLSA:
		server.updateAsExternalLSA(lsdbKey, lsaKey)

	}
	return nil
}

/* @processLSDatabaseTicker
Visited every time ticker is fired to
check expired LSAs and send message to
flood LSA
*/
func (server *OSPFServer) processLSDatabaseTicker() {
	/* scan through LSDB. Flood expired LSAs and
	   delete from LSDB */
	for lsdbKey, lsDbEnt := range server.AreaLsdb {
		server.processMaxAgeLSA(lsdbKey, lsDbEnt)

	}

}

/*
@fn processMaxAgeLsaMsg
Once the flooding thread floods Max age LSA , this message is sent to
update MaxAgeLSA map.
*/
func (server *OSPFServer) processMaxAgeLsaMsg(msg maxAgeLsaMsg) {
	switch msg.msg_type {
	case addMaxAgeLsa:
		maxAgeLsaMap[msg.lsaKey] = msg.pkt
	case delMaxAgeLsa:
		delete(maxAgeLsaMap, msg.lsaKey)
	}
}
