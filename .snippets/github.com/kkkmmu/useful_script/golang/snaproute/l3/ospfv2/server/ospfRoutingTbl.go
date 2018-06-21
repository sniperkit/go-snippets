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
	//"asicd/asicdCommonDefs"
	"errors"
	"fmt"
	"ribd"
	"strconv"
)

type RoutingTblStruct struct {
	TempAreaRoutingTbl   map[AreaIdKey]AreaRoutingTbl
	GlobalRoutingTbl     map[RoutingTblEntryKey]GlobalRoutingTblEntry
	OldGlobalRoutingTbl  map[RoutingTblEntryKey]GlobalRoutingTblEntry
	TempGlobalRoutingTbl map[RoutingTblEntryKey]GlobalRoutingTblEntry
}

type DestType uint8

const (
	Network         DestType = 0
	InternalRouter  DestType = 1
	ASBdrRouter     DestType = 2
	AreaBdrRouter   DestType = 3
	ASAreaBdrRouter DestType = 4
)

type PathType int

const (
	/* Decreasing order of Precedence */
	IntraArea PathType = 4
	InterArea PathType = 3
	Type1Ext  PathType = 2
	Type2Ext  PathType = 1
)

type IfData struct {
	IfIpAddr uint32
	IfIdx    uint32
}

type NbrIP uint32

type NextHop struct {
	IfIPAddr  uint32
	IfIdx     uint32
	NextHopIP uint32
	AdvRtr    uint32 // Nbr Router Id
}

type AreaIdKey struct {
	AreaId uint32
}

type RoutingTblEntryKey struct {
	DestId   uint32   // IP address(Network Type) RouterID(Router Type)
	AddrMask uint32   // Only For Network Type
	DestType DestType // true: Network, false: Router
}

type AreaRoutingTbl struct {
	RoutingTblMap map[RoutingTblEntryKey]RoutingTblEntry
}

type RoutingTblEntry struct {
	OptCapabilities uint8    // Optional Capabilities
	PathType        PathType // Path Type
	Cost            uint16
	Type2Cost       uint16
	LSOrigin        LsaKey
	NumOfPaths      int
	NextHops        map[NextHop]bool // Next Hop
}

type GlobalRoutingTblEntry struct {
	AreaId        uint32 // Area
	RoutingTblEnt RoutingTblEntry
}

func (server *OSPFV2Server) dumpRoutingTbl() {
	server.logger.Info("=============Routing Table============")
	server.logger.Info("DestId      AddrMask        DestType        OprCapabilities Area    PathType        Cost    Type2Cost       LSOrigin        NumOfPaths      NextHops")
	for areaIdKey, areaEnt := range server.RoutingTblData.TempAreaRoutingTbl {
		server.logger.Info("=============Area Id:", areaIdKey.AreaId, " =====================")
		for key, ent := range areaEnt.RoutingTblMap {
			DestId := convertUint32ToDotNotation(key.DestId)
			AddrMask := convertUint32ToDotNotation(key.AddrMask)
			var DestType string
			if key.DestType == Network {
				DestType = "Network"
			} else if key.DestType == InternalRouter {
				DestType = "Internal Router"
			} else if key.DestType == AreaBdrRouter {
				DestType = "Area Border Router"
			} else if key.DestType == ASBdrRouter {
				DestType = "ASBdrRouter"
			} else if key.DestType == ASAreaBdrRouter {
				DestType = "AS Area Border Router"
			} else {
				DestType = "Invalid"
			}
			Area := convertUint32ToDotNotation(areaIdKey.AreaId)
			var PathType string
			if ent.PathType == IntraArea {
				PathType = "IntraArea"
			} else if ent.PathType == InterArea {
				PathType = "InterArea"
			} else if ent.PathType == Type1Ext {
				PathType = "Type1Ext"
			} else {
				PathType = "Type2Ext"
			}
			var LsaType string
			var LsaLSId string
			var LsaAdvRouter string
			if ent.PathType == IntraArea {
				if ent.LSOrigin.LSType == RouterLSA {
					LsaType = "RouterLSA"
				} else if ent.LSOrigin.LSType == NetworkLSA {
					LsaType = "NetworkLSA"
				}
				LsaLSId = convertUint32ToDotNotation(ent.LSOrigin.LSId)
				LsaAdvRouter = convertUint32ToDotNotation(ent.LSOrigin.AdvRouter)
			}
			var NextHops string = "["
			for nxtHopKey, _ := range ent.NextHops {
				NextHops = NextHops + "{"
				IfIPAddr := convertUint32ToDotNotation(nxtHopKey.IfIPAddr)
				NextHopIP := convertUint32ToDotNotation(nxtHopKey.NextHopIP)
				AdvRtr := convertUint32ToDotNotation(nxtHopKey.AdvRtr)
				nextHops := fmt.Sprint("IfIpAddr:", IfIPAddr, "IfIdx:", nxtHopKey.IfIdx, "NextHopIP:", NextHopIP, "AdvRtr:", AdvRtr)
				NextHops = NextHops + nextHops
				NextHops = NextHops + "}"
			}
			NextHops = NextHops + "]"
			if ent.PathType == IntraArea {
				server.logger.Info(DestId, AddrMask, DestType, ent.OptCapabilities, Area, PathType, ent.Cost, ent.Type2Cost, "[", LsaType, LsaLSId, LsaAdvRouter, "]", ent.NumOfPaths, NextHops)
			} else {
				server.logger.Info(DestId, AddrMask, DestType, ent.OptCapabilities, Area, PathType, ent.Cost, ent.Type2Cost, "[ ---------------------------------- ]", ent.NumOfPaths, NextHops)
			}
		}
	}
	server.logger.Info("==============End of Routing Table================")
}

func (server *OSPFV2Server) dumpGlobalRoutingTbl() {
	server.logger.Info("=============Routing Table============")
	server.logger.Info("DestId      AddrMask        DestType        OprCapabilities Area    PathType        Cost    Type2Cost       LSOrigin        NumOfPaths      NextHops")
	for key, ent := range server.RoutingTblData.GlobalRoutingTbl {
		DestId := convertUint32ToDotNotation(key.DestId)
		AddrMask := convertUint32ToDotNotation(key.AddrMask)
		var DestType string
		if key.DestType == Network {
			DestType = "Network"
		} else if key.DestType == InternalRouter {
			DestType = "Internal Router"
		} else if key.DestType == AreaBdrRouter {
			DestType = "Area Border Router"
		} else if key.DestType == ASBdrRouter {
			DestType = "ASBdrRouter"
		} else if key.DestType == ASAreaBdrRouter {
			DestType = "AS Area Border Router"
		} else {
			DestType = "None"
		}
		Area := convertUint32ToDotNotation(ent.AreaId)
		var PathType string
		if ent.RoutingTblEnt.PathType == IntraArea {
			PathType = "IntraArea"
		} else if ent.RoutingTblEnt.PathType == InterArea {
			PathType = "InterArea"
		} else if ent.RoutingTblEnt.PathType == Type1Ext {
			PathType = "Type1Ext"
		} else {
			PathType = "Type2Ext"
		}
		var LsaType string
		var LsaLSId string
		var LsaAdvRouter string
		if ent.RoutingTblEnt.PathType == IntraArea {
			if ent.RoutingTblEnt.LSOrigin.LSType == RouterLSA {
				LsaType = "RouterLSA"
			} else if ent.RoutingTblEnt.LSOrigin.LSType == NetworkLSA {
				LsaType = "NetworkLSA"
			}
			LsaLSId = convertUint32ToDotNotation(ent.RoutingTblEnt.LSOrigin.LSId)
			LsaAdvRouter = convertUint32ToDotNotation(ent.RoutingTblEnt.LSOrigin.AdvRouter)
		}
		var NextHops string = "["
		for nxtHopKey, _ := range ent.RoutingTblEnt.NextHops {
			NextHops = NextHops + "{"
			IfIPAddr := convertUint32ToDotNotation(nxtHopKey.IfIPAddr)
			NextHopIP := convertUint32ToDotNotation(nxtHopKey.NextHopIP)
			AdvRtr := convertUint32ToDotNotation(nxtHopKey.AdvRtr)
			nextHops := fmt.Sprint("IfIpAddr:", IfIPAddr, "IfIdx:", nxtHopKey.IfIdx, "NextHopIP:", NextHopIP, "AdvRtr:", AdvRtr)
			NextHops = NextHops + nextHops
			NextHops = NextHops + "}"
		}
		NextHops = NextHops + "]"
		if ent.RoutingTblEnt.PathType == IntraArea {
			server.logger.Info(DestId, AddrMask, DestType, ent.RoutingTblEnt.OptCapabilities, Area, PathType, ent.RoutingTblEnt.Cost, ent.RoutingTblEnt.Type2Cost, "[", LsaType, LsaLSId, LsaAdvRouter, "]", ent.RoutingTblEnt.NumOfPaths, NextHops)
		} else {
			server.logger.Info(DestId, AddrMask, DestType, ent.RoutingTblEnt.OptCapabilities, Area, PathType, ent.RoutingTblEnt.Cost, ent.RoutingTblEnt.Type2Cost, "[ ---------------------------------- ]", ent.RoutingTblEnt.NumOfPaths, NextHops)
		}
	}
	server.logger.Info("==============End of Routing Table================")
}

func (server *OSPFV2Server) findP2PNextHopIP(vFirst VertexKey, vSecond VertexKey, areaIdKey AreaIdKey) (ifIPAddr uint32, nextHopIP uint32, err error) {
	// Our link is P2P
	lsDbKey := LsdbKey{
		AreaId: areaIdKey.AreaId,
	}
	lsDbEnt, exist := server.LsdbData.AreaLsdb[lsDbKey]
	if !exist {
		server.logger.Err("No LS Database found for areaId:", areaIdKey.AreaId)
		err = errors.New("No LS Database found")
		return 0, 0, err
	}
	firstLsaKey := LsaKey{
		LSType:    RouterLSA,
		LSId:      vFirst.ID,
		AdvRouter: vFirst.AdvRtr,
	}
	firstLsa, exist := lsDbEnt.RouterLsaMap[firstLsaKey]
	if !exist {
		server.logger.Err("Unable to find the Router Lsa for first node")
		err = errors.New("Unable to find the Router Lsa for first node")
		return 0, 0, err
	}
	if firstLsa.LsaMd.LSAge == MAX_AGE {
		server.logger.Err("Router Lsa for first node with MAX_AGE")
		return 0, 0, errors.New("Router Lsa for first node with MAX_AGE")
	}
	secondLsaKey := LsaKey{
		LSType:    RouterLSA,
		LSId:      vSecond.ID,
		AdvRouter: vSecond.AdvRtr,
	}
	secondLsa, exist := lsDbEnt.RouterLsaMap[secondLsaKey]
	if !exist {
		server.logger.Err("Unable to find the Router Lsa for second node")
		err = errors.New("Unable to find the Router Lsa for second node")
		return 0, 0, err
	}
	if secondLsa.LsaMd.LSAge == MAX_AGE {
		server.logger.Err("Router Lsa for second node with MAX_AGE")
		return 0, 0, errors.New("Router Lsa for second node with MAX_AGE")
	}
	var firstLink LinkDetail
	flag := false
	var secondLink LinkDetail
	for _, link := range firstLsa.LinkDetails {
		if link.LinkId == vSecond.AdvRtr {
			firstLink = link
			flag = true
			break
		}
	}
	if flag == false {
		err = errors.New("Unable to find the Link for second vertex")
		return 0, 0, err
	} else {
		flag = false
	}
	for _, link := range secondLsa.LinkDetails {
		if link.LinkId == vFirst.AdvRtr &&
			link.LinkType == P2P_LINK {
			secondLink = link
			flag = true
			break
		}
	}

	if flag == false {
		err = errors.New("Unable to find the Link for first vertex")
		return 0, 0, err
	}
	ifIPAddr = firstLink.LinkData
	nextHopIP = secondLink.LinkData
	return ifIPAddr, nextHopIP, nil

}

func (server *OSPFV2Server) UpdateRoutingTblForRouter(areaIdKey AreaIdKey, vKey VertexKey, tVertex TreeVertex, rootVKey VertexKey) {
	server.logger.Info("Updating Routing Table for Router Vertex", vKey, tVertex)

	gEnt, exist := server.SPFData.AreaGraph[vKey]
	if !exist {
		server.logger.Err("Entry doesn't exist in Area Graph for:", vKey)
		return
	}
	lsDbKey := LsdbKey{
		AreaId: areaIdKey.AreaId,
	}
	lsDbEnt, exist := server.LsdbData.AreaLsdb[lsDbKey]
	if !exist {
		server.logger.Err("No LS Database found for areaId:", areaIdKey.AreaId)
		return
	}
	lsaEnt, exist := lsDbEnt.RouterLsaMap[gEnt.LsaKey]
	if !exist {
		server.logger.Err("No LS Database Entry found for lsaKey:", gEnt.LsaKey)
		return
	}
	if lsaEnt.LsaMd.LSAge == MAX_AGE {
		server.logger.Err("Router LSA with MAX_AGE", gEnt.LsaKey)
		return
	}
	var destType DestType
	if lsaEnt.BitB == true &&
		lsaEnt.BitE == true {
		destType = ASAreaBdrRouter
	} else if lsaEnt.BitB == true {
		destType = AreaBdrRouter
	} else if lsaEnt.BitE == true {
		destType = ASBdrRouter
	} else {
		destType = InternalRouter
	}
	rKey := RoutingTblEntryKey{
		DestType: destType,
		AddrMask: 0, //TODO
		DestId:   vKey.ID,
	}

	tempAreaRoutingTbl := server.RoutingTblData.TempAreaRoutingTbl[areaIdKey]
	rEnt, exist := tempAreaRoutingTbl.RoutingTblMap[rKey]
	if exist {
		server.logger.Info("Routing Tbl entry already exist for:", rKey)
		return
	}

	rEnt.OptCapabilities = 0 //TODO
	//rEnt.Area = gEnt.AreaId
	rEnt.PathType = IntraArea
	rEnt.Cost = tVertex.Distance
	rEnt.Type2Cost = 0 //TODO
	rEnt.LSOrigin = gEnt.LsaKey
	rEnt.NumOfPaths = tVertex.NumOfPaths
	rEnt.NextHops = make(map[NextHop]bool, tVertex.NumOfPaths)
	if rootVKey == vKey {
		//rEnt.AdvRtr = vKey.AdvRtr
		tempAreaRoutingTbl.RoutingTblMap[rKey] = rEnt
		server.RoutingTblData.TempAreaRoutingTbl[areaIdKey] = tempAreaRoutingTbl
		return
	}
	for i := 0; i < tVertex.NumOfPaths; i++ {
		pathlen := len(tVertex.Paths[i])
		if tVertex.Paths[i][0] != rootVKey {
			server.logger.Info("Starting vertex is not our router, hence ignoring this path")
			continue
		}
		//if pathlen < 2 {
		if pathlen < 1 {
			server.logger.Info("Connected Route so no next hops")
			continue
		}
		var vFirst VertexKey
		var vSecond VertexKey
		if pathlen == 1 {
			vFirst = tVertex.Paths[i][0]
			vSecond = vKey
		} else {
			vFirst = tVertex.Paths[i][0]
			vSecond = tVertex.Paths[i][1]
		}
		var ifIPAddr uint32
		var nextHopIP uint32
		if vFirst.Type == RouterVertex &&
			vSecond.Type == RouterVertex {
			var err error
			ifIPAddr, nextHopIP, err = server.findP2PNextHopIP(vFirst, vSecond, areaIdKey)
			if err != nil {
				server.logger.Err("Error in find P2P Next HOP IP:", err)
				continue
			}
			server.logger.Info("P2P ifIPAddr:", ifIPAddr, "nextHopIP:", nextHopIP)
		} else {
			var vThird VertexKey
			if pathlen == 2 {
				vThird = vKey
			} else {
				vThird = tVertex.Paths[i][2]
			}
			gFirst, exist := server.SPFData.AreaGraph[vFirst]
			if !exist {
				server.logger.Info("1. Entry does not exist for:", vFirst, "in Area Graph")
				continue
			}
			gThird, exist := server.SPFData.AreaGraph[vThird]
			if !exist {
				server.logger.Info("3. Entry does not exist for:", vThird, "in Area Graph")
				continue
			}
			ifIPAddr = gFirst.LinkData[vSecond]
			nextHopIP = gThird.LinkData[vSecond]
		}
		nextHop := NextHop{
			IfIPAddr:  ifIPAddr,
			IfIdx:     0, //TODO
			NextHopIP: nextHopIP,
			AdvRtr:    0,
		}
		rEnt.NextHops[nextHop] = true
	}
	//rEnt.AdvRtr = vKey.AdvRtr
	tempAreaRoutingTbl.RoutingTblMap[rKey] = rEnt
	server.RoutingTblData.TempAreaRoutingTbl[areaIdKey] = tempAreaRoutingTbl
}

func (server *OSPFV2Server) UpdateRoutingTblWithStub(areaId uint32, vKey VertexKey, tVertex TreeVertex, parent TreeVertex, parentKey VertexKey, rootVKey VertexKey) {
	areaIdKey := AreaIdKey{
		AreaId: areaId,
	}

	server.logger.Info("Fetching Routing Table for Router Vertex", parentKey)

	if parentKey == rootVKey {
		server.logger.Info("Parent Key is same as root Key")
	}
	pEnt, exist := server.SPFData.AreaGraph[parentKey]
	if !exist {
		server.logger.Err("Entry doesn't exist in Area Graph for:", parentKey)
		return
	}
	lsDbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, exist := server.LsdbData.AreaLsdb[lsDbKey]
	if !exist {
		server.logger.Err("No LS Database found for areaId:", areaId)
		return
	}
	lsaEnt, exist := lsDbEnt.RouterLsaMap[pEnt.LsaKey]
	if !exist {
		server.logger.Err("No LS Database Entry found for lsaKey:", pEnt.LsaKey)
		return
	}
	if lsaEnt.LsaMd.LSAge == MAX_AGE {
		server.logger.Err("Router LSA with MAX_AGE", pEnt.LsaKey)
		return
	}
	var destType DestType
	if lsaEnt.BitB == true &&
		lsaEnt.BitE == true {
		destType = ASAreaBdrRouter
	} else if lsaEnt.BitB == true {
		destType = AreaBdrRouter
	} else if lsaEnt.BitE == true {
		destType = ASBdrRouter
	} else {
		destType = InternalRouter
	}
	pKey := RoutingTblEntryKey{
		DestType: destType,
		AddrMask: 0, //TODO
		DestId:   parentKey.ID,
	}

	tempAreaRoutingTbl := server.RoutingTblData.TempAreaRoutingTbl[areaIdKey]
	pREnt, exist := tempAreaRoutingTbl.RoutingTblMap[pKey]
	if !exist {
		server.logger.Info("Routing Tbl doesnot exist for:", pKey)
		return
	}

	sEnt, exist := server.SPFData.AreaStubs[vKey]
	if !exist {
		server.logger.Err("Entry doesn't exist in Area Stubs for:", vKey)
		return
	}

	rKey := RoutingTblEntryKey{
		DestType: Network,
		AddrMask: sEnt.LinkData,
		DestId:   vKey.ID,
	}
	rEnt, exist := tempAreaRoutingTbl.RoutingTblMap[rKey]
	if exist {
		server.logger.Info("Routing Tbl entry for Stub already exist for:", rKey)
		return
	}
	rEnt.OptCapabilities = pREnt.OptCapabilities //TODO
	rEnt.PathType = IntraArea                    //TODO
	rEnt.Cost = tVertex.Distance
	rEnt.Type2Cost = 0 //TODO
	rEnt.LSOrigin = sEnt.LsaKey
	rEnt.NumOfPaths = tVertex.NumOfPaths
	rEnt.NextHops = make(map[NextHop]bool, tVertex.NumOfPaths)
	for key, _ := range pREnt.NextHops {
		rEnt.NextHops[key] = true
	}
	//rEnt.AdvRtr = vKey.AdvRtr
	tempAreaRoutingTbl.RoutingTblMap[rKey] = rEnt
	server.RoutingTblData.TempAreaRoutingTbl[areaIdKey] = tempAreaRoutingTbl
}

/*
func (server *OSPFV2Server) UpdateRoutingTblForSNetwork(areaIdKey AreaIdKey, vKey VertexKey, tVertex TreeVertex, rootVKey VertexKey) {
	server.logger.Info("Updating Routing Table for Stub Network Vertex", vKey, tVertex)

	sEnt, exist := server.SPFData.AreaStubs[vKey]
	if !exist {
		server.logger.Err("Entry doesn't exist in Area Stubs for:", vKey)
		return
	}
	rKey := RoutingTblEntryKey{
		DestType: Network,
		AddrMask: sEnt.LinkData, //TODO
		DestId:   vKey.ID,
	}

	tempAreaRoutingTbl := server.RoutingTblData.TempAreaRoutingTbl[areaIdKey]
	rEnt, exist := tempAreaRoutingTbl.RoutingTblMap[rKey]
	if exist {
		server.logger.Info("Routing Tbl entry already exist for:", rKey)
		return
	}

	rEnt.OptCapabilities = 0  //TODO
	rEnt.PathType = IntraArea //TODO
	rEnt.Cost = tVertex.Distance
	rEnt.Type2Cost = 0 //TODO
	rEnt.LSOrigin = sEnt.LsaKey
	rEnt.NumOfPaths = tVertex.NumOfPaths
	if rEnt.NumOfPaths == 0 {
		server.logger.Info("==============Hello1===========")
	}
	rEnt.NextHops = make(map[NextHop]bool, tVertex.NumOfPaths)
	for i := 0; i < tVertex.NumOfPaths; i++ {
		pathlen := len(tVertex.Paths[i])
		if tVertex.Paths[i][0] != rootVKey {
			server.logger.Info("Starting vertex is not our router, hence ignoring this path")
			continue
		}
		if pathlen < 3 { //Path Example {R1}, {R1, N1, R2} -- TODO
			server.logger.Info("Connected Route so no next hops")
			continue
		}
		vFirst := tVertex.Paths[i][0]
		vSecond := tVertex.Paths[i][1]
		vThird := tVertex.Paths[i][2]
		gFirst, exist := server.SPFData.AreaGraph[vFirst]
		if !exist {
			server.logger.Info("1. Entry does not exist for:", vFirst, "in Area Graph")
			continue
		}
		gThird, exist := server.SPFData.AreaGraph[vThird]
		if !exist {
			server.logger.Info("3. Entry does not exist for:", vThird, "in Area Graph")
			continue
		}
		ifIPAddr := gFirst.LinkData[vSecond]
		nextHopIP := gThird.LinkData[vSecond]
		nextHop := NextHop{
			IfIPAddr:  ifIPAddr,
			IfIdx:     0, //TODO
			NextHopIP: nextHopIP,
			AdvRtr:    0,
		}
		rEnt.NextHops[nextHop] = true
	}
	//rEnt.AdvRtr = vKey.AdvRtr
	tempAreaRoutingTbl.RoutingTblMap[rKey] = rEnt
	server.RoutingTblData.TempAreaRoutingTbl[areaIdKey] = tempAreaRoutingTbl
}
*/

func (server *OSPFV2Server) UpdateRoutingTblForTNetwork(areaIdKey AreaIdKey, vKey VertexKey, tVertex TreeVertex, rootVKey VertexKey) {
	var err error
	server.logger.Info("Updating Routing Table for Transit Network Vertex", vKey, tVertex)

	gEnt, exist := server.SPFData.AreaGraph[vKey]
	if !exist {
		server.logger.Err("Entry doesn't exist in Area Graph for:", vKey)
		return
	}

	//Need to add check for len of gEnt.NbrVertexKey
	if len(gEnt.NbrVertexKey) < 1 {
		server.logger.Info("Vertex", vKey, "is listed as Transit but doesn't have any Neighboring routers")
		return
	}
	addrMask, exist := gEnt.LinkData[gEnt.NbrVertexKey[0]]
	if !exist {
		server.logger.Err("Vertex", vKey, "has neighboring router but no corresponding linkdata")
	}
	rKey := RoutingTblEntryKey{
		DestType: Network,
		AddrMask: addrMask, //TODO
		DestId:   vKey.ID & addrMask,
	}

	tempAreaRoutingTbl := server.RoutingTblData.TempAreaRoutingTbl[areaIdKey]
	rEnt, exist := tempAreaRoutingTbl.RoutingTblMap[rKey]
	if exist {
		server.logger.Info("Routing Tbl entry already exist for:", rKey)
		return
	}

	rEnt.OptCapabilities = 0 //TODO
	//rEnt.Area = gEnt.AreaId
	rEnt.PathType = IntraArea //TODO
	rEnt.Cost = tVertex.Distance
	rEnt.Type2Cost = 0 //TODO
	rEnt.LSOrigin = gEnt.LsaKey
	rEnt.NumOfPaths = tVertex.NumOfPaths
	rEnt.NextHops = make(map[NextHop]bool, tVertex.NumOfPaths)

	for i := 0; i < tVertex.NumOfPaths; i++ {
		pathlen := len(tVertex.Paths[i])
		if tVertex.Paths[i][0] != rootVKey {
			server.logger.Info("Starting vertex is not our router, hence ignoring this path")
			continue
		}
		if pathlen < 3 { //Path Example {R1}, {R1, N1, R2} -- TODO
			if pathlen == 2 { // Path Example {R1, R2} P2P
				vFirst := tVertex.Paths[i][0]
				vSecond := tVertex.Paths[i][1]
				if vFirst.Type != RouterVertex ||
					vSecond.Type != RouterVertex {
					server.logger.Info("Connected Route so no next hops")
					continue
				}
			} else {
				server.logger.Info("Connected Route so no next hops")
				continue
			}
		}
		vFirst := tVertex.Paths[i][0]
		vSecond := tVertex.Paths[i][1]
		var ifIPAddr uint32
		var nextHopIP uint32
		if vFirst.Type == RouterVertex &&
			vSecond.Type == RouterVertex {
			ifIPAddr, nextHopIP, err = server.findP2PNextHopIP(vFirst, vSecond, areaIdKey)
			if err != nil {
				server.logger.Err("Error in find P2P Next HOP IP:", err)
				continue
			}
			server.logger.Info("P2P ifIPAddr:", ifIPAddr, "nextHopIP:", nextHopIP)
		} else {
			vThird := tVertex.Paths[i][2]
			gFirst, exist := server.SPFData.AreaGraph[vFirst]
			if !exist {
				server.logger.Info("1. Entry does not exist for:", vFirst, "in Area Graph")
				continue
			}
			gThird, exist := server.SPFData.AreaGraph[vThird]
			if !exist {
				server.logger.Info("3. Entry does not exist for:", vThird, "in Area Graph")
				continue
			}
			ifIPAddr = gFirst.LinkData[vSecond]
			nextHopIP = gThird.LinkData[vSecond]
		}
		nextHop := NextHop{
			IfIPAddr:  ifIPAddr,
			IfIdx:     0, //TODO
			NextHopIP: nextHopIP,
			AdvRtr:    0,
		}
		rEnt.NextHops[nextHop] = true
	}
	//rEnt.AdvRtr = vKey.AdvRtr
	tempAreaRoutingTbl.RoutingTblMap[rKey] = rEnt
	server.RoutingTblData.TempAreaRoutingTbl[areaIdKey] = tempAreaRoutingTbl
}

// Compare Old and New Route
func (server *OSPFV2Server) CompareRoutes(rKey RoutingTblEntryKey) bool {
	oldEnt, exist := server.RoutingTblData.OldGlobalRoutingTbl[rKey]
	if !exist {
		server.logger.Err("No Route with", rKey, "was there in Old Routing Table")
		return true
	}
	newEnt, exist := server.RoutingTblData.TempGlobalRoutingTbl[rKey]
	if !exist {
		server.logger.Err("No Route with", rKey, "is there in New Routing Table")
		return true
	}
	if oldEnt.RoutingTblEnt.Cost != newEnt.RoutingTblEnt.Cost {
		return false
	}
	if len(oldEnt.RoutingTblEnt.NextHops) != len(newEnt.RoutingTblEnt.NextHops) {
		return false
	}

	for key, _ := range oldEnt.RoutingTblEnt.NextHops {
		_, exist := newEnt.RoutingTblEnt.NextHops[key]
		if !exist {
			return false
		}
	}
	return true
}

func (server *OSPFV2Server) DeleteRoute(rKey RoutingTblEntryKey) {
	server.logger.Info("Deleting route for rKey:", rKey)
	oldEnt, exist := server.RoutingTblData.OldGlobalRoutingTbl[rKey]
	if !exist {
		server.logger.Info("No route installed for rKey:", rKey, "hence, not deleting it")
		return
	}
	destNetIp := convertUint32ToDotNotation(rKey.DestId)
	networkMask := convertUint32ToDotNotation(rKey.AddrMask)
	routeType := "OSPF"
	for key, _ := range oldEnt.RoutingTblEnt.NextHops {
		nextHopIp := convertUint32ToDotNotation(key.NextHopIP)
		server.logger.Info("Deleting Route: destNetIp:", destNetIp, "networkMask:", networkMask, "nextHopIp:", nextHopIp, "routeType:", routeType)
		cfg := ribd.IPv4Route{
			DestinationNw: destNetIp,
			Protocol:      routeType,
			Cost:          0,
			NetworkMask:   networkMask,
		}
		nextHopInfo := ribd.NextHopInfo{
			NextHopIp: nextHopIp,
		}
		cfg.NextHop = make([]*ribd.NextHopInfo, 0)
		cfg.NextHop = append(cfg.NextHop, &nextHopInfo)
		if server.ribdComm.ribdClient.ClientHdl == nil {
			server.logger.Err("Nil ribd handle. Can not delete route. ")
			return
		}

		ret, err := server.ribdComm.ribdClient.ClientHdl.DeleteIPv4Route(&cfg)
		if err != nil {
			server.logger.Err("Error Deleting Route:", err)
		}
		server.logger.Info("Return Value for RIB DeleteV4Route call: ", ret)
	}
	msg := RouteDelMsg{
		RTblKey: rKey,
	}
	server.SendRouteDelMsgToDBClnt(msg)
}

func (server *OSPFV2Server) InstallRoute(rKey RoutingTblEntryKey) {
	server.logger.Info("Installing new route for rKey", rKey)
	newEnt, exist := server.RoutingTblData.TempGlobalRoutingTbl[rKey]
	if !exist {
		server.logger.Info("No new routing table entry exist for rkey:", rKey, "hence not installing it")
		return
	}
	destNetIp := convertUint32ToDotNotation(rKey.DestId)
	networkMask := convertUint32ToDotNotation(rKey.AddrMask)
	metric := ribd.Int(newEnt.RoutingTblEnt.Cost)
	routeType := "OSPF"
	for key, _ := range newEnt.RoutingTblEnt.NextHops {
		nextHopIp := convertUint32ToDotNotation(key.NextHopIP)
		ifIdx, exist := server.infraData.ipToIfIdxMap[key.IfIPAddr]
		if !exist {
			server.logger.Err("Unable to find entry for ip:", key.IfIPAddr, "in ipToIfIdxMap")
			continue
		}
		//nextHopIfIndex := asicdCommonDefs.GetIfIndexFromIntfIdAndIntfType(int(ipProp.IfId), int(ipProp.IfType))
		nextHopIfIndex := ifIdx
		server.logger.Info("Installing Route: destNetIp:", destNetIp, "networkMask:", networkMask, "metric:", metric, "nextHopIp:", nextHopIp, "nextHopIfIndex:", nextHopIfIndex, "routeType:", routeType)
		cfg := ribd.IPv4Route{
			DestinationNw: destNetIp,
			Protocol:      routeType,
			Cost:          int32(metric),
			NetworkMask:   networkMask,
		}
		nextHopInfo := ribd.NextHopInfo{
			NextHopIp:     nextHopIp,
			NextHopIntRef: strconv.Itoa(int(nextHopIfIndex)),
		}
		cfg.NextHop = make([]*ribd.NextHopInfo, 0)
		cfg.NextHop = append(cfg.NextHop, &nextHopInfo)
		ret, err := server.ribdComm.ribdClient.ClientHdl.CreateIPv4Route(&cfg)
		if err != nil {
			server.logger.Err("Error Installing Route:", err, ret)
			continue
		}
	}
	msg := RouteAddMsg{
		RTblKey:   rKey,
		RTblEntry: newEnt,
	}
	server.SendRouteAddMsgToDBClnt(msg)
}

func (server *OSPFV2Server) ConsolidatingRoutingTbl() {
	for areaId, _ := range server.AreaConfMap {
		areaIdKey := AreaIdKey{
			AreaId: areaId,
		}
		tempAreaRoutingTbl, exist := server.RoutingTblData.TempAreaRoutingTbl[areaIdKey]
		if !exist {
			continue
		}

		for rKey, rEnt := range tempAreaRoutingTbl.RoutingTblMap {
			ent, exist := server.RoutingTblData.TempGlobalRoutingTbl[rKey]
			if exist {

			} else {
				ent.AreaId = areaId
				ent.RoutingTblEnt = rEnt
			}
			server.RoutingTblData.TempGlobalRoutingTbl[rKey] = ent
		}
	}
}

func (server *OSPFV2Server) UpdateRoute(rKey RoutingTblEntryKey) {
	server.logger.Info("Updating route for rKey:", rKey)
	// Delete Old Route
	server.DeleteRoute(rKey)
	// Install New Route
	server.InstallRoute(rKey)
}

func (server *OSPFV2Server) InstallRoutingTbl() {
	server.logger.Info("Routing Table Consolidation:")
	server.ConsolidatingRoutingTbl()
	server.logger.Info("Installing Routing Table ")

	OldRoutingTblKeys := make(map[RoutingTblEntryKey]bool)
	NewRoutingTblKeys := make(map[RoutingTblEntryKey]bool)

	for rKey, rEnt := range server.RoutingTblData.OldGlobalRoutingTbl {
		if rKey.DestType != Network {
			continue
		}
		if len(rEnt.RoutingTblEnt.NextHops) > 0 {
			OldRoutingTblKeys[rKey] = false
		}
	}

	for rKey, rEnt := range server.RoutingTblData.TempGlobalRoutingTbl {
		if rKey.DestType != Network {
			continue
		}
		if len(rEnt.RoutingTblEnt.NextHops) > 0 {
			NewRoutingTblKeys[rKey] = false
		}
	}
	for rKey, _ := range NewRoutingTblKeys {
		_, exist := OldRoutingTblKeys[rKey]
		if exist {
			ret := server.CompareRoutes(rKey)
			if ret == false { // Old Routes and New Routes are not same
				server.UpdateRoute(rKey)
				OldRoutingTblKeys[rKey] = true
				NewRoutingTblKeys[rKey] = true
			} else { // Old Routes and New Routes are same
				OldRoutingTblKeys[rKey] = true
				NewRoutingTblKeys[rKey] = true
			}
		}
	}

	for rKey, ent := range OldRoutingTblKeys {
		if ent == false {
			server.DeleteRoute(rKey)
		}
		OldRoutingTblKeys[rKey] = true
	}

	for rKey, ent := range NewRoutingTblKeys {
		if ent == false {
			server.InstallRoute(rKey)
		}
		NewRoutingTblKeys[rKey] = true
	}
}

func (server *OSPFV2Server) FlushRoutingTbl() {
	for rKey, rEnt := range server.RoutingTblData.GlobalRoutingTbl {
		server.logger.Info("Deleting route for rKey:", rKey)
		destNetIp := convertUint32ToDotNotation(rKey.DestId)
		networkMask := convertUint32ToDotNotation(rKey.AddrMask)
		routeType := "OSPF"
		for key, _ := range rEnt.RoutingTblEnt.NextHops {
			nextHopIp := convertUint32ToDotNotation(key.NextHopIP)
			cfg := ribd.IPv4Route{
				DestinationNw: destNetIp,
				Protocol:      routeType,
				Cost:          0,
				NetworkMask:   networkMask,
			}
			nextHopInfo := ribd.NextHopInfo{
				NextHopIp: nextHopIp,
			}
			cfg.NextHop = make([]*ribd.NextHopInfo, 0)
			cfg.NextHop = append(cfg.NextHop, &nextHopInfo)
			if server.ribdComm.ribdClient.ClientHdl == nil {
				server.logger.Err("Nil ribd handle. Can not delete route. ")
				continue
			}

			ret, err := server.ribdComm.ribdClient.ClientHdl.DeleteIPv4Route(&cfg)
			if err != nil {
				server.logger.Err("Error Deleting Route:", err)
			}
			server.logger.Info("Return Value for RIB DeleteV4Route call: ", ret)
			/*
			   err = server.DelIPv4RoutesState(rKey)
			   if err != nil {
			           server.logger.Info(fmt.Sprintln("DB: Failed to delete route from db. route , err ", rKey, err))
			   }
			*/
		}
		msg := RouteDelMsg{
			RTblKey: rKey,
		}
		server.SendRouteDelMsgToDBClnt(msg)
	}
}
