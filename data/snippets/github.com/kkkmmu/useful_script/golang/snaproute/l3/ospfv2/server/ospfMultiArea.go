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
//"fmt"
)

type SummaryLsaMap map[LsaKey]SummaryLsa

// Handling Summary LSA in case of FS is internal router
func (server *OSPFV2Server) CalcInterAreaRoutes(areaId uint32) {
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, exist := server.LsdbData.AreaLsdb[lsdbKey]
	if !exist {
		server.logger.Err("Unable to find Area Lsdb entry")
		return
	}

	for lsaKey, lsaEnt := range lsDbEnt.Summary3LsaMap {
		server.logger.Info("Summary LSAKey:", lsaKey, "lsaENt:", lsaEnt)
		if lsaEnt.Metric == LSInfinity ||
			lsaEnt.LsaMd.LSAge == MAX_AGE {
			server.logger.Info("Ignoring Summary LSA...")
			continue
		}
		rtrId := server.globalData.RouterId
		if lsaKey.AdvRouter == rtrId {
			server.logger.Info("Self originated summary 3 LSA, so no need to process for routing table calc")
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
			server.logger.Info("Area Router routing table entry doesnot exists for Summary Lsa Advertising Router")
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
			DestId:   lsaKey.LSId & lsaEnt.Netmask,
			AddrMask: lsaEnt.Netmask,
			DestType: Network,
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
			rEnt.PathType = InterArea
			rEnt.Cost = cost
			rEnt.Type2Cost = 0
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

func (server *OSPFV2Server) HandleSummaryType3Lsa(areaId uint32) {
	server.logger.Info("Inside Handling Summary Type 3 LSA for area:", areaId)
	server.CalcInterAreaRoutes(areaId)
}

func (server *OSPFV2Server) HandleSummaryLsa(areaId uint32) {
	server.logger.Info("Inside Handling Summary LSA for area:", areaId)
	server.HandleSummaryType3Lsa(areaId)
	server.HandleSummaryType4Lsa(areaId)
	server.HandleASExternalLsa(areaId)
}

func (server *OSPFV2Server) HandleTransitAreaSummaryLsa() {
	//TODO
}

func (server *OSPFV2Server) GenerateSummaryLsa() {
	server.logger.Info("Generating Summary LSA")
	server.SummaryLsDb = make(map[LsdbKey]SummaryLsaMap)
	for areaId, aEnt := range server.AreaConfMap {
		if len(aEnt.IntfMap) == 0 {
			continue
		}

		lsDbKey := LsdbKey{
			AreaId: areaId,
		}
		isStub, _ := server.isStubArea(areaId)

		sEnt, _ := server.SummaryLsDb[lsDbKey]
		sEnt = make(map[LsaKey]SummaryLsa)
		for rKey, rEnt := range server.RoutingTblData.GlobalRoutingTbl {
			if rKey.DestType == AreaBdrRouter ||
				rEnt.RoutingTblEnt.PathType == Type1Ext ||
				rEnt.RoutingTblEnt.PathType == Type2Ext ||
				rEnt.AreaId == areaId ||
				rKey.DestType == InternalRouter ||
				uint32(rEnt.RoutingTblEnt.Cost) >= LSInfinity {
				// 1. If Dest Type is Area Border Router
				// 2. If Path Type is Type 1 Ext
				// 3. If Path Type is Type 2 Ext
				// 4. If Area associated with set of paths is Area itself
				// 5. If Cost >= LSInfinity
				// 6. TODO: Distance Vector Split Horizon Problem
				continue
			}

			// TODO: AS External Routes
			// If DestType == ASBdrRouter
			// If Routing Table Entry Describes the preferred path to
			// AS Boundary Router
			/*
			           RFC 2328 section 12.4.3
			           if the destination of this route is an AS boundary
			           router, a summary-LSA should be originated if and only
			   if the routing table entry describes the preferred path
			   to the AS boundary router (see Step 3 of Section 16.4).
			   If so, a Type 4 summary-LSA is originated for the
			   destination, with Link State ID equal to the AS boundary
			   router's Router ID and metric equal to the routing table
			   entry's cost. Note: these LSAs should not be generated
			    if Area A has been configured as a stub area.
			*/

			if (rKey.DestType == ASAreaBdrRouter ||
				rKey.DestType == ASBdrRouter) && !isStub {
				lsaKey, summaryLsa := server.GenerateType4SummaryLSA(rKey, rEnt, lsDbKey)
				sEnt[lsaKey] = summaryLsa
			}

			// Dest Type Network, Inter Area Routes
			if rKey.DestType == Network &&
				rEnt.RoutingTblEnt.PathType == InterArea {
				// Generate Type 3 Summary LSA for the desitnation
				// LSId = networks's address
				// Metric = Routing Table cost
				lsaKey, summaryLsa := server.GenerateType3SummaryLSA(rKey, rEnt, lsDbKey)
				sEnt[lsaKey] = summaryLsa
			} else if rKey.DestType == Network &&
				rEnt.RoutingTblEnt.PathType == IntraArea {
				//TODO: Address Range
				// By default LSId = network's address
				// Metric = Routing Table cost
				server.logger.Debug("Summary : generated summary 3 lsa ", rKey)
				lsaKey, summaryLsa := server.GenerateType3SummaryLSA(rKey, rEnt, lsDbKey)
				sEnt[lsaKey] = summaryLsa
			}
		}

		server.SummaryLsDb[lsDbKey] = sEnt
		if isStub {
			lsaKey, defsummaryLsa := server.GenerateDefaultSummary3LSA(lsDbKey)
			sEnt[lsaKey] = defsummaryLsa
		}
	}
}

/*
@fn GenerateDefaultSummary3LSA

 RFC 2328 A 4.4
 For stub areas, Type 3 summary-LSAs can also be used to describe a
 (per-area) default route.  Default summary routes are used in stub
  areas instead of flooding a complete set of external routes.  When
 describing a default summary route, the summary-LSA's Link State ID
 is always set to DefaultDestination (0.0.0.0) and the Network Mask
 is set to 0.0.0.0.

*/
func (server *OSPFV2Server) GenerateDefaultSummary3LSA(lsDbKey LsdbKey) (LsaKey, SummaryLsa) {
	var summaryLsa SummaryLsa
	seq_num := int(InitialSequenceNum)
	metric := int32(20)
	/* TODO
	   areaConfKey := AreaConfKey{
	           AreaId: config.AreaId(convertUint32ToIPv4(lsDbKey.AreaId)),
	   }

	   conf, exist := server.AreaConfMap[lsDbKey.AreaId]
	   if exist {
	           metric = conf.StubDefaultCost
	   }
	*/
	AdvRouter := server.globalData.RouterId
	lsaKey := LsaKey{
		LSType:    Summary3LSA,
		LSId:      0,
		AdvRouter: AdvRouter,
	}

	//check if summaryLSA exist
	if lsdbEnt, ok := server.SummaryLsDb[lsDbKey]; ok {
		if summaryLsa, ok = lsdbEnt[lsaKey]; ok {
			seq_num = summaryLsa.LsaMd.LSSequenceNum + 1
			server.logger.Info("SUMMARY: Refreshed summary LSA ", lsaKey)
		}
	}

	summaryLsa.LsaMd.Options = uint8(0)
	summaryLsa.LsaMd.LSAge = 0
	summaryLsa.LsaMd.LSSequenceNum = seq_num
	summaryLsa.LsaMd.LSLen = uint16(OSPF_LSA_HEADER_SIZE + 8)
	summaryLsa.Netmask = 0
	summaryLsa.Metric = uint32(metric)

	return lsaKey, summaryLsa
}

func (server *OSPFV2Server) GenerateType3SummaryLSA(rKey RoutingTblEntryKey, rEnt GlobalRoutingTblEntry, lsDbKey LsdbKey) (LsaKey, SummaryLsa) {
	var summaryLsa SummaryLsa
	seq_num := int(InitialSequenceNum)

	//AdvRouter := convertIPv4ToUint32(server.ospfGlobalConf.RouterId)
	AdvRouter := server.globalData.RouterId
	lsaKey := LsaKey{
		LSType:    Summary3LSA,
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

	summaryLsa.LsaMd.Options = uint8(2) //Need to be re-visited
	summaryLsa.LsaMd.LSAge = 0
	summaryLsa.LsaMd.LSSequenceNum = seq_num
	summaryLsa.LsaMd.LSLen = uint16(OSPF_LSA_HEADER_SIZE + 8)
	summaryLsa.Netmask = rKey.AddrMask
	summaryLsa.Metric = uint32(rEnt.RoutingTblEnt.Cost)

	return lsaKey, summaryLsa
}
