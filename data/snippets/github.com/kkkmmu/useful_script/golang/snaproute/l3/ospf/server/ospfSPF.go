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
	"fmt"
	"l3/ospf/config"
	"sort"
)

type VertexKey struct {
	Type   uint8
	ID     uint32
	AdvRtr uint32
}

type Path []VertexKey

type TreeVertex struct {
	Paths      []Path
	Distance   uint16
	NumOfPaths int
}

type StubVertex struct {
	NbrVertexKey  VertexKey
	NbrVertexCost uint16
	LinkData      uint32
	LsaKey        LsaKey
	AreaId        uint32
	LinkStateId   uint32
}

type Vertex struct {
	NbrVertexKey  []VertexKey
	NbrVertexCost []uint16
	LinkData      map[VertexKey]uint32
	LsaKey        LsaKey
	AreaId        uint32
	Visited       bool
	LinkStateId   uint32
	NetMask       uint32
}

const (
	RouterVertex   uint8 = 0
	SNetworkVertex uint8 = 1 // Stub
	TNetworkVertex uint8 = 2 // Transit
)

type VertexData struct {
	vKey     VertexKey
	distance uint16
}

var check bool = true

type VertexDataArr []VertexData

func (v VertexDataArr) Len() int {
	return len(v)
}

func (v VertexDataArr) Less(i, j int) bool {
	//return v[i].distance < v[j].distance
	if v[i].distance == v[j].distance {
		if v[i].vKey.Type == v[j].vKey.Type {
			return false
		} else if v[i].vKey.Type == RouterVertex &&
			v[j].vKey.Type == TNetworkVertex {
			return false
		} else if v[i].vKey.Type == TNetworkVertex &&
			v[j].vKey.Type == RouterVertex {
			return true
		}
	}
	return v[i].distance < v[j].distance
}

func (v VertexDataArr) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func findSelfOrigRouterLsaKey(ent map[LsaKey]bool) (LsaKey, error) {
	var key LsaKey
	for key, _ := range ent {
		if key.LSType == RouterLSA {
			return key, nil
		}
	}
	err := errors.New("No Self Orignated Router LSA found")
	return key, err
}

func (server *OSPFServer) checkRouterLsaConsistency(areaId uint32, rtrId uint32, network uint32, netmask uint32) bool {
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, exist := server.AreaLsdb[lsdbKey]
	if !exist {
		server.logger.Err(fmt.Sprintln("No LS Database found for areaId:", areaId))
		return false
	}
	lsaKey := LsaKey{
		LSType:    RouterLSA,
		LSId:      rtrId,
		AdvRouter: rtrId,
	}
	lsaEnt, exist := lsDbEnt.RouterLsaMap[lsaKey]
	if !exist {
		server.logger.Err(fmt.Sprintln("Unable to find router lsa for ", rtrId))
		return false
	}
	for i := 0; i < int(lsaEnt.NumOfLinks); i++ {
		if (lsaEnt.LinkDetails[i].LinkId & netmask) == network {
			if lsaEnt.LinkDetails[i].LinkType == StubLink {
				server.logger.Err(fmt.Sprintln("Have router lsa which still has a stub link in the network, hence ignoring it, lsaKey:", lsaKey))
				return false
			} else if lsaEnt.LinkDetails[i].LinkType == TransitLink {
				server.logger.Err(fmt.Sprintln("Have router lsa which has a transit link in the network, hence processing it, lsaKey:", lsaKey))
				return true
			}
		}
	}
	return false
}

func (server *OSPFServer) UpdateAreaGraphNetworkLsa(lsaEnt NetworkLsa, lsaKey LsaKey, areaId uint32) error {
	server.logger.Info(fmt.Sprintln("2: Using Lsa with key as:", dumpLsaKey(lsaKey), "for SPF calc"))
	vertexKey := VertexKey{
		Type:   TNetworkVertex,
		ID:     lsaKey.LSId,
		AdvRtr: lsaKey.AdvRouter,
	}
	ent, exist := server.AreaGraph[vertexKey]
	if exist {
		server.logger.Info(fmt.Sprintln("Entry already exists in SPF Graph for vertexKey:", vertexKey))
		server.logger.Info(fmt.Sprintln("SPF Graph:", server.AreaGraph))
		return nil
	}
	netmask := lsaEnt.Netmask
	network := lsaKey.LSId & netmask
	server.logger.Info(fmt.Sprintln("netmask:", netmask, "network:", network))
	ent.NbrVertexKey = make([]VertexKey, 0)
	ent.NbrVertexCost = make([]uint16, 0)
	ent.LinkData = make(map[VertexKey]uint32)
	for i := 0; i < len(lsaEnt.AttachedRtr); i++ {
		Rtr := lsaEnt.AttachedRtr[i]
		ret := server.checkRouterLsaConsistency(areaId, Rtr, network, netmask)
		if ret == false {
			server.logger.Info(fmt.Sprintln("Rtr: ", Rtr, "has a stub link in the network", network, "hence skiping it"))
			continue
		}
		server.logger.Info(fmt.Sprintln("Attached Router at index:", i, "is:", Rtr))
		var vKey VertexKey
		var cost uint16
		vKey = VertexKey{
			Type:   RouterVertex,
			ID:     Rtr,
			AdvRtr: Rtr,
		}
		cost = 0
		ent.NbrVertexKey = append(ent.NbrVertexKey, vKey)
		ent.NbrVertexCost = append(ent.NbrVertexCost, cost)
		ent.LinkData[vKey] = lsaEnt.Netmask
	}
	if len(ent.NbrVertexKey) == 0 {
		server.logger.Err(fmt.Sprintln("No router has transit link in this network", network, " hence skiping this network"))
		return nil
	}
	ent.AreaId = areaId
	ent.LsaKey = lsaKey
	ent.Visited = false
	//ent.NetMask = lsaEnt.NetMask
	ent.LinkStateId = lsaKey.LSId
	server.AreaGraph[vertexKey] = ent
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, exist := server.AreaLsdb[lsdbKey]
	if !exist {
		server.logger.Err(fmt.Sprintln("No LS Database found for areaId:", areaId))
		err := errors.New(fmt.Sprintln("No LS Database found for areaId:", areaId))
		return err
	}
	for _, vKey := range ent.NbrVertexKey {
		_, exist := server.AreaGraph[vKey]
		if exist {
			server.logger.Info(fmt.Sprintln("Entry already exists in SPF Graph for vertexKey:", vertexKey))
			continue
		}
		lsaKey := LsaKey{
			LSType:    RouterLSA,
			LSId:      vKey.ID,
			AdvRouter: vKey.AdvRtr,
		}
		lsaEnt, exist := lsDbEnt.RouterLsaMap[lsaKey]
		if !exist {
			server.logger.Err(fmt.Sprintln("Router LSA with LsaKey:", lsaKey, "not found in areaId:", areaId))
			err := errors.New(fmt.Sprintln("Router LSA with LsaKey:", lsaKey, "not found in areaId:", areaId))
			// continue
			if check == true {
				continue
			}
			return err
		}
		err := server.UpdateAreaGraphRouterLsa(lsaEnt, lsaKey, areaId)
		if err != nil {
			if check == true {
				continue
			}
			return err
		}
	}
	return nil
}

func (server *OSPFServer) findNetworkLsa(areaId uint32, LSId uint32) (lsaKey LsaKey, err error) {
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, exist := server.AreaLsdb[lsdbKey]
	if !exist {
		server.logger.Err(fmt.Sprintln("No LS Database found for areaId:", areaId))
		return
	}

	for key, _ := range lsDbEnt.NetworkLsaMap {
		if key.LSId == LSId &&
			key.LSType == NetworkLSA {
			return key, nil
		}
	}

	err = errors.New("Network LSA not found")
	return lsaKey, err
}

func (server *OSPFServer) findRouterLsa(areaId uint32, LSId uint32) (err error) {
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, exist := server.AreaLsdb[lsdbKey]
	if !exist {
		server.logger.Err(fmt.Sprintln("No LS Database found for areaId:", areaId))
		return
	}

	lsaKey := LsaKey{
		LSType:    RouterLSA,
		LSId:      LSId,
		AdvRouter: LSId,
	}

	_, exist = lsDbEnt.RouterLsaMap[lsaKey]
	if !exist {
		err = errors.New("Router LSA not found")
		return err
	}
	return nil
}

func (server *OSPFServer) UpdateAreaGraphRouterLsa(lsaEnt RouterLsa, lsaKey LsaKey, areaId uint32) error {
	server.logger.Info(fmt.Sprintln("1: Using Lsa with key as:", dumpLsaKey(lsaKey), "for SPF calc"))
	vertexKey := VertexKey{
		Type:   RouterVertex,
		ID:     lsaKey.LSId,
		AdvRtr: lsaKey.AdvRouter,
	}
	ent, exist := server.AreaGraph[vertexKey]
	if exist {
		server.logger.Info(fmt.Sprintln("Entry already exists in SPF Graph for vertexKey:", vertexKey))
		server.logger.Info(fmt.Sprintln("SPF Graph:", server.AreaGraph))
		return nil
	}
	if lsaEnt.BitV == true {
		area := config.AreaId(convertUint32ToIPv4(areaId))
		areaConfKey := AreaConfKey{
			AreaId: area,
		}
		aEnt, exist := server.AreaConfMap[areaConfKey]
		if exist && aEnt.TransitCapability == false {
			aEnt.TransitCapability = true
			server.AreaConfMap[areaConfKey] = aEnt
		}
	}
	ent.NbrVertexKey = make([]VertexKey, 0)
	ent.NbrVertexCost = make([]uint16, 0)
	ent.LinkData = make(map[VertexKey]uint32)
	for i := 0; i < int(lsaEnt.NumOfLinks); i++ {
		server.logger.Info(fmt.Sprintln("Link Detail at index", i, "is:", lsaEnt.LinkDetails[i]))
		linkDetail := lsaEnt.LinkDetails[i]
		var vKey VertexKey
		var cost uint16
		var lData uint32
		if linkDetail.LinkType == TransitLink {
			server.logger.Info("===It is TransitLink===")
			vKey = VertexKey{
				Type:   TNetworkVertex,
				ID:     linkDetail.LinkId,
				AdvRtr: 0,
			}
			nLsaKey, err := server.findNetworkLsa(areaId, vKey.ID)
			if err != nil {
				server.logger.Info(fmt.Sprintln("Err:", err, vKey.ID))
				if check == true {
					continue
				}
				return err
				//continue
			}
			vKey.AdvRtr = nLsaKey.AdvRouter
			cost = linkDetail.LinkMetric
			lData = linkDetail.LinkData
			ent.NbrVertexKey = append(ent.NbrVertexKey, vKey)
			ent.NbrVertexCost = append(ent.NbrVertexCost, cost)
			ent.LinkData[vKey] = lData
		} else if linkDetail.LinkType == StubLink {
			server.logger.Info("===It is StubLink===")
			vKey = VertexKey{
				Type:   SNetworkVertex,
				ID:     linkDetail.LinkId,
				AdvRtr: lsaKey.AdvRouter,
			}
			cost = linkDetail.LinkMetric
			lData = linkDetail.LinkData
			sentry, _ := server.AreaStubs[vKey]
			sentry.NbrVertexKey = vertexKey
			sentry.NbrVertexCost = cost
			sentry.LinkData = lData
			sentry.AreaId = areaId
			sentry.LsaKey = lsaKey
			sentry.LinkStateId = lsaKey.LSId
			server.AreaStubs[vKey] = sentry
		} else if linkDetail.LinkType == P2PLink {
			server.logger.Info("===It is P2PLink===")
			vKey = VertexKey{
				Type:   RouterVertex,
				ID:     linkDetail.LinkId,
				AdvRtr: linkDetail.LinkId,
			}
			err := server.findRouterLsa(areaId, vKey.ID)
			if err != nil {
				server.logger.Info(fmt.Sprintln("Err:", err, vKey.ID))
				return err
			}
			cost = linkDetail.LinkMetric
			lData = linkDetail.LinkData
			ent.NbrVertexKey = append(ent.NbrVertexKey, vKey)
			ent.NbrVertexCost = append(ent.NbrVertexCost, cost)
			ent.LinkData[vKey] = lData
		}
	}
	if len(ent.NbrVertexKey) == 0 {
		// If all the links in self originated Router Lsa is stub
		// Add Self router to AreaGraph
		if lsaKey.AdvRouter == convertIPv4ToUint32(server.ospfGlobalConf.RouterId) {
			ent.AreaId = areaId
			ent.LsaKey = lsaKey
			ent.Visited = false
			ent.LinkStateId = lsaKey.LSId
			server.AreaGraph[vertexKey] = ent
			return nil
		}
		err := errors.New(fmt.Sprintln("None of the Network LSA are found"))
		return err
	}
	ent.AreaId = areaId
	ent.LsaKey = lsaKey
	ent.Visited = false
	ent.LinkStateId = lsaKey.LSId
	server.AreaGraph[vertexKey] = ent
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, exist := server.AreaLsdb[lsdbKey]
	if !exist {
		server.logger.Err(fmt.Sprintln("No LS Database found for areaId:", areaId))
		err := errors.New(fmt.Sprintln("No LS Database found for areaId:", areaId))
		return err
	}
	for _, vKey := range ent.NbrVertexKey {
		_, exist := server.AreaGraph[vKey]
		if exist {
			server.logger.Info(fmt.Sprintln("Entry for Vertex:", vKey, "already exist in Area Graph"))
			continue
		}
		lsaKey := LsaKey{
			LSType:    0,
			LSId:      vKey.ID,
			AdvRouter: vKey.AdvRtr,
		}
		if vKey.Type == TNetworkVertex {
			lsaKey.LSType = NetworkLSA
			lsaEnt, exist := lsDbEnt.NetworkLsaMap[lsaKey]
			if !exist {
				server.logger.Err(fmt.Sprintln("Network LSA with LsaKey:", lsaKey, "not found in LS Database of areaId:", areaId))
				err := errors.New(fmt.Sprintln("Network LSA with LsaKey:", lsaKey, "not found in LS Database of areaId:", areaId))
				if check == true {
					continue
				}
				//continue
				return err
			}
			err := server.UpdateAreaGraphNetworkLsa(lsaEnt, lsaKey, areaId)
			if err != nil {
				if check == true {
					continue
				}
				return err
			}
		} else if vKey.Type == RouterVertex {
			lsaKey.LSType = RouterLSA
			lsaEnt, exist := lsDbEnt.RouterLsaMap[lsaKey]
			if !exist {
				server.logger.Err(fmt.Sprintln("Router LSA with LsaKey:", lsaKey, "not found in LS Database of areaId:", areaId))
				err := errors.New(fmt.Sprintln("Router LSA with LsaKey:", lsaKey, "not found in LS Database of areaId:", areaId))
				//continue
				if check == true {
					continue
				}
				return err
			}
			err := server.UpdateAreaGraphRouterLsa(lsaEnt, lsaKey, areaId)
			if err != nil {
				if check == true {
					continue
				}
				return err
			}
		}
	}
	return nil
}

func (server *OSPFServer) CreateAreaGraph(areaId uint32) (VertexKey, error) {
	var vKey VertexKey
	server.logger.Info(fmt.Sprintln("Create SPF Graph for: areaId:", areaId))
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, exist := server.AreaLsdb[lsdbKey]
	if !exist {
		server.logger.Err(fmt.Sprintln("No LS Database found for areaId:", areaId))
		err := errors.New(fmt.Sprintln("No LS Database found for areaId:", areaId))
		return vKey, err
	}

	selfOrigLsaEnt, exist := server.AreaSelfOrigLsa[lsdbKey]
	if !exist {
		server.logger.Err(fmt.Sprintln("No Self Originated LSAs found for areaId:", areaId))
		err := errors.New(fmt.Sprintln("No Self Originated LSAs found for areaId:", areaId))
		return vKey, err
	}
	selfRtrLsaKey, err := findSelfOrigRouterLsaKey(selfOrigLsaEnt)
	if err != nil {
		server.logger.Err(fmt.Sprintln("No Self Originated Router LSA Key found for areaId:", areaId))
		err := errors.New(fmt.Sprintln("No Self Originated Router LSA Key found for areaId:", areaId))
		return vKey, err
	}
	server.logger.Info(fmt.Sprintln("Self Orginated Router LSA Key:", selfRtrLsaKey))
	lsaEnt, exist := lsDbEnt.RouterLsaMap[selfRtrLsaKey]
	if !exist {
		server.logger.Err(fmt.Sprintln("No Self Originated Router LSA found for areaId:", areaId))
		err := errors.New(fmt.Sprintln("No Self Originated Router LSA found for areaId:", areaId))
		return vKey, err
	}

	err = server.UpdateAreaGraphRouterLsa(lsaEnt, selfRtrLsaKey, areaId)
	if check == true {
		err = nil
	}
	vKey = VertexKey{
		Type:   RouterVertex,
		ID:     selfRtrLsaKey.LSId,
		AdvRtr: selfRtrLsaKey.AdvRouter,
	}
	return vKey, err
}

func (server *OSPFServer) ExecuteDijkstra(vKey VertexKey, areaId uint32) error {
	//var treeVSlice []VertexKey = make([]VertexKey, 0)
	var treeVSlice []VertexData = make([]VertexData, 0)

	//treeVSlice = append(treeVSlice, vKey)
	vData := VertexData{
		vKey:     vKey,
		distance: 0,
	}
	treeVSlice = append(treeVSlice, vData)
	ent, exist := server.SPFTree[vKey]
	if !exist {
		ent.Distance = 0
		ent.NumOfPaths = 1
		ent.Paths = make([]Path, 1)
		var path Path
		path = make(Path, 0)
		ent.Paths[0] = path
		server.SPFTree[vKey] = ent
	}

	for j := 0; j < len(treeVSlice); j++ {
		verArr := make([]VertexData, 0)
		server.logger.Debug(fmt.Sprintln("treeVSlice:", treeVSlice))
		server.logger.Debug(fmt.Sprintln("The value of j:", j, "treeVSlice:", treeVSlice[j].vKey))
		//ent, exist := server.AreaGraph[treeVSlice[j]]
		ent, exist := server.AreaGraph[treeVSlice[j].vKey]
		if !exist {
			server.logger.Info(fmt.Sprintln("No entry found for:", treeVSlice[j].vKey))
			err := errors.New(fmt.Sprintln("No entry found for:", treeVSlice[j].vKey))
			return err
		}

		server.logger.Debug(fmt.Sprintln("Number of neighbor of treeVSlice:", treeVSlice[j].vKey, "ent:", ent, "is:", len(ent.NbrVertexKey)))
		for i := 0; i < len(ent.NbrVertexKey); i++ {
			verKey := ent.NbrVertexKey[i]
			cost := ent.NbrVertexCost[i]
			entry, exist := server.AreaGraph[verKey]
			server.logger.Debug(fmt.Sprintln("Neighboring Vertex Number :", i, "verKey", verKey, "cost:", cost, "entry:", entry))
			if !exist {
				server.logger.Err("Something is wrong in SPF Calculation: Entry should exist in Area Graph")
				err := errors.New("Something is wrong in SPF Calculation: Entry should exist in Area Graph")
				return err
			}
			tEnt, exist := server.SPFTree[verKey]
			if !exist {
				server.logger.Debug("Entry doesnot exist for the neighbor in SPF hence adding it")
				tEnt.Paths = make([]Path, 1)
				var path Path
				path = make(Path, 0)
				tEnt.Paths[0] = path
				tEnt.Distance = 0xff00 // LSInfinity
				tEnt.NumOfPaths = 1
			}
			tEntry, exist := server.SPFTree[treeVSlice[j].vKey]
			if !exist {
				server.logger.Err("Something is wrong is SPF Calculation")
				err := errors.New("Something is wrong is SPF Calculation")
				return err
			}
			server.logger.Debug(fmt.Sprintln("Parent Node:", treeVSlice[j].vKey, tEntry))
			server.logger.Debug(fmt.Sprintln("Child Node:", verKey, tEnt))
			if tEnt.Distance > tEntry.Distance+cost {
				server.logger.Debug(fmt.Sprintln("We have lower cost path via", tEntry))
				tEnt.Distance = tEntry.Distance + cost
				for l := 0; l < tEnt.NumOfPaths; l++ {
					tEnt.Paths[l] = nil
				}
				tEnt.Paths = tEnt.Paths[:0]
				tEnt.NumOfPaths = 0
				tEnt.Paths = nil
				tEnt.Paths = make([]Path, tEntry.NumOfPaths)
				for l := 0; l < tEntry.NumOfPaths; l++ {
					var path Path
					path = make(Path, len(tEntry.Paths[l])+1)
					copy(path, tEntry.Paths[l])
					path[len(tEntry.Paths[l])] = treeVSlice[j].vKey
					tEnt.Paths[l] = path
				}
				tEnt.NumOfPaths = tEntry.NumOfPaths
			} else if tEnt.Distance == tEntry.Distance+cost {
				server.logger.Debug(fmt.Sprintln("We have equal cost path via:", tEntry))
				server.logger.Debug(fmt.Sprintln("tEnt:", tEnt, "tEntry:", tEntry))
				server.logger.Debug(fmt.Sprintln("tEnt.NumOfPaths:", tEnt.NumOfPaths, "tEntry.NumOfPaths:", tEntry.NumOfPaths))
				paths := make([]Path, (tEntry.NumOfPaths + tEnt.NumOfPaths))
				for l := 0; l < tEnt.NumOfPaths; l++ {
					var path Path
					path = make(Path, len(tEnt.Paths[l]))
					copy(path, tEnt.Paths[l])
					paths[l] = path
					tEnt.Paths[l] = nil
				}
				server.logger.Debug(fmt.Sprintln("1. paths:", paths))
				tEnt.Paths = tEnt.Paths[:0]
				tEnt.Paths = nil
				for l := 0; l < tEntry.NumOfPaths; l++ {
					var path Path
					path = make(Path, len(tEntry.Paths[l])+1)
					copy(path, tEntry.Paths[l])
					path[len(tEntry.Paths[l])] = treeVSlice[j].vKey
					paths[tEnt.NumOfPaths+l] = path
				}
				server.logger.Debug(fmt.Sprintln("2. paths:", paths))
				tEnt.Paths = paths
				tEnt.NumOfPaths = tEntry.NumOfPaths + tEnt.NumOfPaths
			}
			if _, ok := server.SPFTree[verKey]; !ok {
				server.logger.Debug(fmt.Sprintln("Adding verKey:", verKey, "to treeVSlice"))
				vData := VertexData{
					vKey:     verKey,
					distance: tEnt.Distance,
				}
				verArr = append(verArr, vData)
			}
			server.SPFTree[verKey] = tEnt
		}
		treeVSlice = append(treeVSlice, verArr...)
		if len(treeVSlice[j+1:]) > 0 {
			sort.Sort(VertexDataArr(treeVSlice[j+1:]))
		}
		verArr = verArr[:0]
		verArr = nil
		ent.Visited = true
		server.AreaGraph[treeVSlice[j].vKey] = ent
	}

	return nil
}

func (server *OSPFServer) HandleStubs(vKey VertexKey, areaId uint32) {
	server.logger.Info("Handle Stub Networks")
	for key, entry := range server.AreaStubs {
		//Finding the Vertex(Router) to which this stub is connected to
		vertexKey := entry.NbrVertexKey
		parent, exist := server.SPFTree[vertexKey]
		if !exist {
			continue
		}
		ent, _ := server.SPFTree[key]
		ent.Distance = parent.Distance + entry.NbrVertexCost
		ent.Paths = make([]Path, parent.NumOfPaths)
		for i := 0; i < parent.NumOfPaths; i++ {
			var path Path
			path = make(Path, len(parent.Paths[i])+1)
			copy(path, parent.Paths[i])
			path[len(parent.Paths[i])] = vertexKey
			ent.Paths[i] = path
		}
		ent.NumOfPaths = parent.NumOfPaths
		server.UpdateRoutingTblWithStub(areaId, key, ent, parent, vertexKey, vKey)
	}
}

func dumpVertexKey(key VertexKey) string {
	var Type string
	if key.Type == RouterVertex {
		Type = "Router"
	} else if key.Type == SNetworkVertex {
		Type = "Stub"
	} else if key.Type == TNetworkVertex {
		Type = "Transit"
	}
	ID := convertUint32ToIPv4(key.ID)
	AdvRtr := convertUint32ToIPv4(key.AdvRtr)
	return fmt.Sprintln("Vertex Key[Type:", Type, "ID:", ID, "AdvRtr:", AdvRtr)

}

func dumpLsaKey(key LsaKey) string {
	var Type string
	if key.LSType == RouterLSA {
		Type = "Router LSA"
	} else if key.LSType == NetworkLSA {
		Type = "Network LSA"
	}

	LSId := convertUint32ToIPv4(key.LSId)
	AdvRtr := convertUint32ToIPv4(key.AdvRouter)

	return fmt.Sprintln("LSA Type:", Type, "LSId:", LSId, "AdvRtr:", AdvRtr)
}

/*
func (server *OSPFServer) dumpAreaStubs() {
	server.logger.Info("=======================Dump Area Stubs======================")
	for key, ent := range server.AreaStubs {
		server.logger.Info("==================================================")
		server.logger.Info(fmt.Sprintln("Vertex Keys:", dumpVertexKey(key)))
		server.logger.Info("==================================================")
		LData := convertUint32ToIPv4(ent.LinkData)
		server.logger.Info(fmt.Sprintln("VertexKeys:", dumpVertexKey(ent.NbrVertexKey), "Cost:", ent.NbrVertexCost, "LinkData:", LData))
		server.logger.Info("==================================================")
		server.logger.Info(fmt.Sprintln("Lsa Key:", dumpLsaKey(ent.LsaKey)))
		server.logger.Info(fmt.Sprintln("AreaId:", ent.AreaId))
		server.logger.Info(fmt.Sprintln("LinkStateId:", ent.LinkStateId))
	}
	server.logger.Info("==================================================")

}


func (server *OSPFServer) dumpAreaGraph() {
	server.logger.Info("=======================Dump Area Graph======================")
	for key, ent := range server.AreaGraph {
		server.logger.Info("==================================================")
		server.logger.Info(fmt.Sprintln("Vertex Keys:", dumpVertexKey(key)))
		server.logger.Info("==================================================")
		if len(ent.NbrVertexKey) != 0 {
			server.logger.Info("List of Neighbor Vertices(except stub)")
		} else {
			server.logger.Info("No Neighbor Vertices(except stub)")
		}
		for i := 0; i < len(ent.NbrVertexKey); i++ {
			LData := convertUint32ToIPv4(ent.LinkData[ent.NbrVertexKey[i]])
			server.logger.Info(fmt.Sprintln("VertexKeys:", dumpVertexKey(ent.NbrVertexKey[i]), "Cost:", ent.NbrVertexCost[i], "LinkData:", LData))
		}
		server.logger.Info("==================================================")
		server.logger.Info(fmt.Sprintln("Lsa Key:", dumpLsaKey(ent.LsaKey)))
		server.logger.Info(fmt.Sprintln("AreaId:", ent.AreaId))
		server.logger.Info(fmt.Sprintln("Visited:", ent.Visited))
		server.logger.Info(fmt.Sprintln("LinkStateId:", ent.LinkStateId))
	}
	server.logger.Info("==================================================")
}

func (server *OSPFServer) dumpSPFTree() {
	server.logger.Info("=======================Dump SPF Tree======================")
	for key, ent := range server.SPFTree {
		server.logger.Info("==================================================")
		server.logger.Info(fmt.Sprintln("Vertex Keys:", dumpVertexKey(key)))
		server.logger.Info("==================================================")
		server.logger.Info(fmt.Sprintln("Distance:", ent.Distance))
		server.logger.Info(fmt.Sprintln("NumOfPaths:", ent.NumOfPaths))
		for i := 0; i < ent.NumOfPaths; i++ {
			var paths string
			paths = fmt.Sprintln("Path[", i, "]")
			for j := 0; j < len(ent.Paths[i]); j++ {
				paths = paths + fmt.Sprintln("[", dumpVertexKey(ent.Paths[i][j]), "]")
			}
			server.logger.Info(fmt.Sprintln(paths))
		}

	}
}
*/

func (server *OSPFServer) UpdateRoutingTbl(vKey VertexKey, areaId uint32) {
	areaIdKey := AreaIdKey{
		AreaId: areaId,
	}
	for key, ent := range server.SPFTree {
		if vKey == key {
			server.logger.Info("It's own vertex")
			//continue
		}
		switch key.Type {
		case RouterVertex:
			//TODO: If Bit V is set in corresponding RtrLsa set Transit Capability
			// for the given Area
			//rEnt, exist := lsDbEnt.Summary3LsaMap[ent.LsaKey]
			server.UpdateRoutingTblForRouter(areaIdKey, key, ent, vKey)
		//case SNetworkVertex:
		//	server.UpdateRoutingTblForSNetwork(areaIdKey, key, ent, vKey)
		case TNetworkVertex:
			server.UpdateRoutingTblForTNetwork(areaIdKey, key, ent, vKey)
		}
	}
}

func (server *OSPFServer) initRoutingTbl(areaId uint32) {
	server.GlobalRoutingTbl = make(map[RoutingTblEntryKey]GlobalRoutingTblEntry)
}

func (server *OSPFServer) initialiseSPFStructs() {
	server.AreaGraph = make(map[VertexKey]Vertex)
	server.AreaStubs = make(map[VertexKey]StubVertex)
	server.SPFTree = make(map[VertexKey]TreeVertex)

}

func (server *OSPFServer) spfCalculation() {
	for {
		msg := <-server.StartCalcSPFCh
		server.logger.Info(fmt.Sprintln("Recevd SPF Calculation Notification for:", msg))
		server.logger.Info(fmt.Sprintln("Area LS Database:", server.AreaLsdb))
		// Create New Routing table
		// Invalidate Old Routing table
		// Backup Old Routing table
		// Have Per Area Routing Tbl
		// Initialize Algorithm's Data Structure
		server.OldGlobalRoutingTbl = nil
		server.OldGlobalRoutingTbl = make(map[RoutingTblEntryKey]GlobalRoutingTblEntry)
		server.OldGlobalRoutingTbl = server.GlobalRoutingTbl
		server.TempAreaRoutingTbl = nil
		server.TempAreaRoutingTbl = make(map[AreaIdKey]AreaRoutingTbl)
		for key, aEnt := range server.AreaConfMap {

			//server.logger.Info(fmt.Sprintln("===========Area Id : ", key.AreaId, "Area Bdr Status:", server.ospfGlobalConf.isABR, "======================================================="))
			if len(aEnt.IntfListMap) == 0 {
				continue
			}
			aEnt.TransitCapability = false
			areaId := convertAreaOrRouterIdUint32(string(key.AreaId))
			server.initialiseSPFStructs()
			areaIdKey := AreaIdKey{
				AreaId: areaId,
			}

			tempRoutingTbl := server.TempAreaRoutingTbl[areaIdKey]
			tempRoutingTbl.RoutingTblMap = make(map[RoutingTblEntryKey]RoutingTblEntry)
			server.TempAreaRoutingTbl[areaIdKey] = tempRoutingTbl

			vKey, err := server.CreateAreaGraph(areaId)
			if err != nil {
				server.logger.Err(fmt.Sprintln("Error while creating graph for areaId:", areaId))
				//flag = true
				continue
			}
			//server.logger.Info("=========================Start before Dijkstra=================")
			//server.dumpAreaGraph()
			//server.dumpAreaStubs()
			//server.logger.Info("=========================End before Dijkstra=================")
			//server.printRouterLsa()
			err = server.ExecuteDijkstra(vKey, areaId)
			if err != nil {
				server.logger.Err(fmt.Sprintln("Error while executing Dijkstra for areaId:", areaId))
				//flag = true
				continue
			}
			server.logger.Info("=========================Start after Dijkstra=================")
			//	server.dumpAreaGraph()
			//	server.dumpAreaStubs()
			//	server.dumpSPFTree()
			server.logger.Info("=========================End after Dijkstra=================")
			server.UpdateRoutingTbl(vKey, areaId)
			server.logger.Info("==============Handling Stub links...====================")
			server.HandleStubs(vKey, areaId)
			server.HandleSummaryLsa(areaId)
			server.AreaGraph = nil
			server.AreaStubs = nil
			server.SPFTree = nil
		}
		/*
			server.dumpRoutingTbl()
		*/
		server.TempGlobalRoutingTbl = nil
		server.TempGlobalRoutingTbl = make(map[RoutingTblEntryKey]GlobalRoutingTblEntry)
		/* Summarize and Install/Delete Routes In Routing Table */
		server.InstallRoutingTbl()
		// Copy the Summarize Routing Table in Global Routing Table
		server.GlobalRoutingTbl = nil
		server.GlobalRoutingTbl = make(map[RoutingTblEntryKey]GlobalRoutingTblEntry)
		server.GlobalRoutingTbl = server.TempGlobalRoutingTbl
		//server.dumpGlobalRoutingTbl()
		for key, _ := range server.AreaConfMap {
			areaId := convertAreaOrRouterIdUint32(string(key.AreaId))
			areaIdKey := AreaIdKey{
				AreaId: areaId,
			}
			tempAreaRoutingTbl := server.TempAreaRoutingTbl[areaIdKey]
			tempAreaRoutingTbl.RoutingTblMap = nil
			server.TempAreaRoutingTbl[areaIdKey] = tempAreaRoutingTbl
		}
		server.TempAreaRoutingTbl = nil
		server.OldGlobalRoutingTbl = nil
		server.TempGlobalRoutingTbl = nil
		//server.dumpGlobalRoutingTbl()
		if server.ospfGlobalConf.AreaBdrRtrStatus == true {
			server.logger.Info("Examine transit areas, Summary LSA...")
			server.HandleTransitAreaSummaryLsa()
			server.logger.Info("Generate Summary LSA...")
			server.GenerateSummaryLsa()
			server.logger.Info(fmt.Sprintln("========", server.SummaryLsDb, "=========="))
		}
		server.DoneCalcSPFCh <- true
	}
}
