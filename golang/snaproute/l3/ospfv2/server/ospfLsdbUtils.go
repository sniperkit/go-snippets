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

func (server *OSPFV2Server) getRouterLsaFromLsdb(areaId uint32, lsaKey LsaKey) (lsa RouterLsa, retVal bool) {
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, _ := server.LsdbData.AreaLsdb[lsdbKey]
	lsa, exist := lsDbEnt.RouterLsaMap[lsaKey]
	if !exist {
		return lsa, LsdbEntryNotFound
	}
	return lsa, LsdbEntryFound
}

func (server *OSPFV2Server) getNetworkLsaFromLsdb(areaId uint32, lsaKey LsaKey) (lsa NetworkLsa, retVal bool) {
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, _ := server.LsdbData.AreaLsdb[lsdbKey]
	lsa, exist := lsDbEnt.NetworkLsaMap[lsaKey]
	if !exist {
		return lsa, LsdbEntryNotFound
	}
	return lsa, LsdbEntryFound
}

func (server *OSPFV2Server) getSummaryLsaFromLsdb(areaId uint32, lsaKey LsaKey) (lsa SummaryLsa, retVal bool) {
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, exist := server.LsdbData.AreaLsdb[lsdbKey]
	if !exist {
		return lsa, LsdbEntryNotFound
	}
	if lsaKey.LSType == Summary3LSA {
		lsa, exist = lsDbEnt.Summary3LsaMap[lsaKey]
		if !exist {
			return lsa, LsdbEntryNotFound
		}
	} else if lsaKey.LSType == Summary4LSA {
		lsa, exist = lsDbEnt.Summary4LsaMap[lsaKey]
		if !exist {
			return lsa, LsdbEntryNotFound
		}
	} else {
		return lsa, LsdbEntryNotFound
	}
	return lsa, LsdbEntryFound
}

func (server *OSPFV2Server) getASExternalLsaFromLsdb(areaId uint32, lsaKey LsaKey) (lsa ASExternalLsa, retVal bool) {
	lsdbKey := LsdbKey{
		AreaId: areaId,
	}
	lsDbEnt, _ := server.LsdbData.AreaLsdb[lsdbKey]
	lsa, exist := lsDbEnt.ASExternalLsaMap[lsaKey]
	if !exist {
		return lsa, LsdbEntryNotFound
	}
	return lsa, LsdbEntryFound
}
