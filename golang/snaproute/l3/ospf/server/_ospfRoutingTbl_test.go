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
	"fmt"
	"testing"
)

func initrTableTestParams() {
	fmt.Println("\n Get Server object")
	ospf = getServerObject()
	initAttr()
	ospf.GlobalRoutingTbl[rKey] = rEntry
	go startDummyChannels(ospf)
}

func TestOspfrTable(t *testing.T) {
	fmt.Println("\n**************** ROUTING TABLE ************\n")
	initrTableTestParams()
	for index := 1; index < 21; index++ {
		err := rTableTestLogic(index)
		if err != SUCCESS {
			fmt.Println("Failed test  for routing table. ")
		}
	}
}

func rTableTestLogic(tNum int) int {
	switch tNum {
	case 1:
		ospf.initLSDatabase(lsdbKey.AreaId)

		lsDbEnt, _ := ospf.AreaLsdb[lsdbKey]
		lsDbEnt.RouterLsaMap[routerKey] = routerLsa
		lsDbEnt.NetworkLsaMap[networkKey] = networkLsa

		ospf.AreaLsdb[lsdbKey] = lsDbEnt

		ospf.initRoutingTbl(lsdbKey.AreaId)
		//go ospf.spfCalculation()
		//ospf.StartCalcSPFCh <- true
		ospf.AreaGraph = make(map[VertexKey]Vertex)
		ospf.SPFTree = make(map[VertexKey]TreeVertex)
		ospf.AreaGraph[vKeyR] = vertexR

		ospf.checkRouterLsaConsistency(lsdbKey.AreaId, routerKey.LSId, uint32(8), uint32(10))
		fmt.Println("Running Update graph for network LSA")
		ospf.UpdateAreaGraphNetworkLsa(networkLsa, networkKey, lsdbKey.AreaId)

		fmt.Println("Running search APIs for all LSAs")
		ospf.findNetworkLsa(lsdbKey.AreaId, networkKey.LSId)
		ospf.findRouterLsa(lsdbKey.AreaId, routerKey.LSId)

		fmt.Println("Running Area graph updates")
		ospf.CreateAreaGraph(lsdbKey.AreaId)
		ospf.UpdateAreaGraphRouterLsa(routerLsa, routerKey, lsdbKey.AreaId)

		fmt.Println("Running all display APIs.")
		dumpVertexKey(vKeyR)
		dumpVertexKey(vKeyT)
		dumpVertexKey(vKeyN)

		fmt.Println("Running ExecuteDijkstra")
		ospf.SPFTree[vKeyR] = treeVertex
		ospf.SPFTree[vKeyN] = treeVertex
		ospf.ExecuteDijkstra(vKeyR, lsdbKey.AreaId)
		ospf.ExecuteDijkstra(vKeyN, lsdbKey.AreaId)
		ospf.ExecuteDijkstra(vKeyT, lsdbKey.AreaId)

		fmt.Println(" Running SPF calculations.")
		go ospf.spfCalculation()
		ospf.StartCalcSPFCh <- true

		ospf.TempAreaRoutingTbl = make(map[AreaIdKey]AreaRoutingTbl)
		ospf.TempAreaRoutingTbl[areaidkey] = areaRoutingTable
		ospf.UpdateRoutingTbl(vKeyR, lsdbKey.AreaId)
		ospf.UpdateRoutingTbl(vKeyN, lsdbKey.AreaId)
		ospf.UpdateRoutingTbl(vKeyT, lsdbKey.AreaId)
		fmt.Println(" Running routing table display.")
		//ospf.dumpGlobalRoutingTbl()

		fmt.Println("Running install routing table.")
		ospf.InstallRoutingTbl()

		ospf.TempGlobalRoutingTbl = make(map[RoutingTblEntryKey]GlobalRoutingTblEntry)
		//init area config and install routes
		ospf.initAreaConfDefault()
		ospf.processAreaConfig(areaConf)
		ospf.ConsolidatingRoutingTbl()

		ospf.GenerateSummaryLsa()
		lsakey, lsa := ospf.GenerateDefaultSummary3LSA(lsdbKey)
		fmt.Println("Generated summary LSAs ", "lsaKey ", lsakey, " lsa ", lsa)
		lsakey, slsa := ospf.GenerateType3SummaryLSA(rKey, rEntry, lsdbKey)
		fmt.Println("Generated type 3 summary lsa ", "lsakey ", lsakey, " lsa ", slsa)

		/* install routing table */
		ospf.InstallRoute(rKey)
		ospf.UpdateRoute(rKey)
		ospf.CompareRoutes(rKey)
		ospf.DeleteRoute(rKey)

		ospf.AreaGraph[vKeyR] = vertexR
		ospf.AreaGraph[vKeyN] = vertexN
		ospf.AreaGraph[vKeyT] = vertexT

		ospf.UpdateRoutingTblForTNetwork(areaidkey, vKeyR, treeVertex, vKeyN)
		ospf.UpdateRoutingTblForTNetwork(areaidkey, vKeyR, treeVertex, vKeyT)

		ospf.UpdateRoutingTblWithStub(lsdbKey.AreaId, vKeyN, treeVertex, treeVertex, vKeyT, vKeyT)
		ospf.UpdateRoutingTblWithStub(lsdbKey.AreaId, vKeyT, treeVertex, treeVertex, vKeyT, vKeyN)
		ospf.AddIPv4RoutesState(rKey)
		ospf.DelIPv4RoutesState(rKey)

		checkRibdAPIs()
	}
	return SUCCESS
}

func checkRibdAPIs() {
	err := ospf.startRibdUpdates()
	fmt.Println("started ribd updates with error code ", err)
	ospf.verifyOspfRoute(1112, 122)
	ospf.listenForRIBUpdates("ribd")
	ospf.processRibdNotification(hello)
	ospf.getRibdRoutes()
	ospf.ProcessRibdRoutes(route, uint16(1))
}
