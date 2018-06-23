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

import ()

func (server *OSPFV2Server) HandleSummaryType4Lsa(areaId uint32) {
	server.logger.Info("Inside Handling Summary Type 4 LSA for area:", areaId)
	server.CalcASBorderRoutes(areaId)
}

func (server *OSPFV2Server) HandleASExternalLsa(areaId uint32) {
	// TODO: External Path Preference RFC Section 16.4.1
	// TODO: RFC1583Compatibility
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, exist := server.LsdbData.AreaLsdb[lsdbKey]
	if !exist {
		server.logger.Err("Unable to find Area Lsdb entry")
		return
	}

	for lsaKey, lsaEnt := range lsDbEnt.ASExternalLsaMap {
		server.logger.Info("AS External LSAKey:", lsaKey, "lsaENt:", lsaEnt)
		if lsaEnt.Metric == LSInfinity ||
			lsaEnt.LsaMd.LSAge == MAX_AGE {
			server.logger.Info("Ignoring AS External LSA...")
			continue
		}
		rtrId := server.globalData.RouterId
		if lsaKey.AdvRouter == rtrId {
			server.logger.Info("Self originated AS External LSA, so no need to process for routing table calc")
			continue
		}

		areaIdKey := AreaIdKey{
			AreaId: areaId,
		}

		var rKey RoutingTblEntryKey
		var rEnt RoutingTblEntry
		var exist bool
		tempAreaRoutingTbl := server.RoutingTblData.TempAreaRoutingTbl[areaIdKey]
		if lsaEnt.FwdAddr == 0 {
			//Packet should be sent to ASBr
			rKey = RoutingTblEntryKey{
				DestId:   lsaKey.AdvRouter,
				AddrMask: 0,
				DestType: ASBdrRouter,
			}
			rEnt, exist = tempAreaRoutingTbl.RoutingTblMap[rKey]
			if !exist {
				server.logger.Info("AS Border Router routing table entry doesnot exists for AS External Lsa Advertising Router")
				rKey = RoutingTblEntryKey{
					DestId:   lsaKey.AdvRouter,
					AddrMask: 0,
					DestType: ASAreaBdrRouter,
				}
				rEnt, exist = tempAreaRoutingTbl.RoutingTblMap[rKey]
				if !exist {
					server.logger.Info("AS Area Border Router routing table entry doesnot exists for AS External Lsa Advertising Router")
					continue
				}
			}
		} else {
			// Packet should be sent to forwarding address
			rKey = RoutingTblEntryKey{
				DestId:   lsaEnt.FwdAddr,
				AddrMask: 0,
				DestType: ASBdrRouter,
			}
			rEnt, exist = tempAreaRoutingTbl.RoutingTblMap[rKey]
			if !exist {
				server.logger.Info("AS Border Router routing table entry doesnot exists for AS External Lsa Advertising Router")
				rKey = RoutingTblEntryKey{
					DestId:   lsaEnt.FwdAddr,
					AddrMask: 0,
					DestType: ASAreaBdrRouter,
				}
				rEnt, exist = tempAreaRoutingTbl.RoutingTblMap[rKey]
				if !exist {
					server.logger.Info("AS Area Border Router routing table entry doesnot exists for AS External Lsa Advertising Router")
					continue
				}
			}
		}
		if rEnt.NumOfPaths == 0 {
			continue
		}

		cost := rEnt.Cost + uint16(lsaEnt.Metric)
		nextHopMap := rEnt.NextHops
		numOfNextHops := rEnt.NumOfPaths
		rKey = RoutingTblEntryKey{
			DestId:   lsaKey.LSId & lsaEnt.Netmask,
			AddrMask: lsaEnt.Netmask,
			DestType: Network, // TODO: Need to be revisited
		}

		tempAreaRoutingTbl = server.RoutingTblData.TempAreaRoutingTbl[areaIdKey]
		rEnt, exist = tempAreaRoutingTbl.RoutingTblMap[rKey]
		if exist {
			if rEnt.PathType == IntraArea ||
				rEnt.PathType == InterArea {
				//IntraArea or InterArea Paths are always preferred
				continue
			}
			if rEnt.PathType == Type1Ext &&
				lsaEnt.BitE == true {
				//Type1Ext path is always preferred over Type2Ext
				continue
			}
			var pathType PathType
			if lsaEnt.BitE == true {
				pathType = Type2Ext
			} else {
				pathType = Type1Ext
			}
			if rEnt.Cost < cost &&
				rEnt.PathType == pathType {
				//Routing table entry cost is less and path type is same
				server.logger.Info("Route already exists with lesser cost")
				continue
			} else if (rEnt.Cost > cost &&
				rEnt.PathType == pathType) ||
				(rEnt.Cost < cost &&
					rEnt.PathType == Type2Ext) {
				rEnt.OptCapabilities = 0 //TODO
				//rEnt.PathType = InterArea
				rEnt.PathType = pathType
				rEnt.Cost = cost
				rEnt.Type2Cost = uint16(lsaEnt.Metric)
				//rEnt.LSOrigin = lsaKey
				rEnt.NumOfPaths = numOfNextHops
				rEnt.NextHops = make(map[NextHop]bool)
				for key, _ := range nextHopMap {
					key.AdvRtr = lsaKey.AdvRouter
					rEnt.NextHops[key] = true
				}
			} else {
				cnt := 0
				for key, _ := range nextHopMap {
					_, exist = rEnt.NextHops[key]
					if !exist {
						key.AdvRtr = lsaKey.AdvRouter
						rEnt.NextHops[key] = true
						cnt++
					}
				}
				rEnt.NumOfPaths = numOfNextHops + cnt
			}
		} else {
			rEnt.OptCapabilities = 0 //TODO
			if lsaEnt.BitE == true {
				rEnt.PathType = Type2Ext
			} else {
				rEnt.PathType = Type1Ext
			}
			rEnt.Cost = cost
			rEnt.Type2Cost = uint16(lsaEnt.Metric)
			//rEnt.LSOrigin = lsaKey
			rEnt.NumOfPaths = numOfNextHops
			rEnt.NextHops = make(map[NextHop]bool)
			for key, _ := range nextHopMap {
				key.AdvRtr = lsaKey.AdvRouter
				rEnt.NextHops[key] = true
			}
		}
		tempAreaRoutingTbl.RoutingTblMap[rKey] = rEnt
		server.RoutingTblData.TempAreaRoutingTbl[areaIdKey] = tempAreaRoutingTbl
	}
}

func (server *OSPFV2Server) CalcASBorderRoutes(areaId uint32) {
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, exist := server.LsdbData.AreaLsdb[lsdbKey]
	if !exist {
		server.logger.Err("Unable to find Area Lsdb entry")
		return
	}

	for lsaKey, lsaEnt := range lsDbEnt.Summary4LsaMap {
		server.logger.Info("Summary LSAKey:", lsaKey, "lsaENt:", lsaEnt)
		if lsaEnt.Metric == LSInfinity ||
			lsaEnt.LsaMd.LSAge == MAX_AGE {
			server.logger.Info("Ignoring Summary 4 LSA...")
			continue
		}
		rtrId := server.globalData.RouterId
		if lsaKey.AdvRouter == rtrId {
			server.logger.Info("Self originated summary 4 LSA, so no need to process for routing table calc")
			continue
		}

		areaIdKey := AreaIdKey{
			AreaId: areaId,
		}
		// TODO: Handle Area Range Section 16.2 Point 3
		//Network := lsaKey.LSId & lsaEnt.Netmask
		//Mask := lsaEnt.Netmask
		rKey := RoutingTblEntryKey{
			DestId:   lsaKey.AdvRouter,
			AddrMask: 0,
			DestType: AreaBdrRouter,
		}

		tempAreaRoutingTbl := server.RoutingTblData.TempAreaRoutingTbl[areaIdKey]
		rEnt, exist := tempAreaRoutingTbl.RoutingTblMap[rKey]
		if !exist {
			server.logger.Info("Area Border Router routing table entry doesnot exists for Summary Lsa Advertising Router")
			rKey = RoutingTblEntryKey{
				DestId:   lsaKey.AdvRouter,
				AddrMask: 0,
				DestType: ASAreaBdrRouter,
			}
			rEnt, exist = tempAreaRoutingTbl.RoutingTblMap[rKey]
			if !exist {
				server.logger.Info("AS Area Border Router routing table entry doesnot exists for Summary Lsa Advertising Router")
				continue

			}
		}
		if rEnt.NumOfPaths == 0 {
			continue
		}
		cost := rEnt.Cost + uint16(lsaEnt.Metric)
		nextHopMap := rEnt.NextHops
		numOfNextHops := rEnt.NumOfPaths
		rKey = RoutingTblEntryKey{
			DestId:   lsaKey.LSId,
			AddrMask: 0,
			DestType: ASBdrRouter, // TODO: Need to be revisited
		}

		tempAreaRoutingTbl = server.RoutingTblData.TempAreaRoutingTbl[areaIdKey]
		rEnt, exist = tempAreaRoutingTbl.RoutingTblMap[rKey]
		if exist {
			if rEnt.PathType == IntraArea {
				continue
			}
			if rEnt.Cost < cost {
				server.logger.Info("Route already exists with lesser cost")
				continue
			} else if rEnt.Cost > cost {
				rEnt.OptCapabilities = 0 //TODO
				//rEnt.PathType = InterArea
				rEnt.Cost = cost
				//rEnt.Type2Cost = 0
				//rEnt.LSOrigin = lsaKey
				rEnt.NumOfPaths = numOfNextHops
				rEnt.NextHops = make(map[NextHop]bool)
				for key, _ := range nextHopMap {
					key.AdvRtr = 0
					rEnt.NextHops[key] = true
				}
			} else {
				cnt := 0
				for key, _ := range nextHopMap {
					_, exist = rEnt.NextHops[key]
					if !exist {
						key.AdvRtr = 0
						rEnt.NextHops[key] = true
						cnt++
					}
				}
				rEnt.NumOfPaths = numOfNextHops + cnt
			}
		} else {
			rEnt.OptCapabilities = 0 //TODO
			rEnt.PathType = InterArea
			rEnt.Cost = cost
			rEnt.Type2Cost = 0
			//rEnt.LSOrigin = lsaKey
			rEnt.NumOfPaths = numOfNextHops
			rEnt.NextHops = make(map[NextHop]bool)
			for key, _ := range nextHopMap {
				key.AdvRtr = 0
				rEnt.NextHops[key] = true
			}
		}
		tempAreaRoutingTbl.RoutingTblMap[rKey] = rEnt
		server.RoutingTblData.TempAreaRoutingTbl[areaIdKey] = tempAreaRoutingTbl
	}
}

func (server *OSPFV2Server) GenerateType4SummaryLSA(rKey RoutingTblEntryKey, rEnt GlobalRoutingTblEntry, lsDbKey LsdbKey) (LsaKey, SummaryLsa) {
	var summaryLsa SummaryLsa
	seq_num := int(InitialSequenceNum)

	AdvRouter := server.globalData.RouterId
	lsaKey := LsaKey{
		LSType:    Summary4LSA,
		LSId:      rKey.DestId,
		AdvRouter: AdvRouter,
	}
	//check if summaryLSA exist
	if lsdbEnt, ok := server.SummaryLsDb[lsDbKey]; ok {
		if summaryLsa, ok = lsdbEnt[lsaKey]; ok {
			seq_num = summaryLsa.LsaMd.LSSequenceNum + 1
			server.logger.Info("SUMMARY: Refreshed summary LSA ", lsaKey)
		}
	}
	summaryLsa.LsaMd.Options = uint8(2) // Need to be re-visited
	summaryLsa.LsaMd.LSAge = 0
	summaryLsa.LsaMd.LSSequenceNum = seq_num
	summaryLsa.LsaMd.LSLen = uint16(OSPF_LSA_HEADER_SIZE + 8)
	summaryLsa.Netmask = rKey.AddrMask
	summaryLsa.Metric = uint32(rEnt.RoutingTblEnt.Cost)

	return lsaKey, summaryLsa
}
