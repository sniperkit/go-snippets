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
	"l3/ndp/debug"
	"models/objects"
	"utils/dbutils"
)

func (svr *NDPServer) readNdpGblCfg(dbHdl *dbutils.DBUtil) {
	if dbHdl == nil {
		debug.Logger.Err("DB cannot be read as dbHDL is nil")
		return
	}
	var dbGblObj objects.NDPGlobal
	objList, err := dbHdl.GetAllObjFromDb(dbGblObj)
	if err != nil {
		debug.Logger.Err("DB Querry failed for NDPGlobal Config", err)
		return
	}
	debug.Logger.Info("Global Object reterived from DB are", objList)
	for _, obj := range objList {
		dbEntry := obj.(objects.NDPGlobal)
		svr.NdpConfig.Vrf = dbEntry.Vrf
		if dbEntry.RouterAdvertisementInterval == 0 {
			debug.Logger.Warning("Invalid Router Advertisment and hence setting default value",
				NDP_DEFAULT_RTR_ADVERTISEMENT_INTERVAL)
			svr.NdpConfig.RaRestransmitTime = NDP_DEFAULT_RTR_ADVERTISEMENT_INTERVAL
		} else {
			svr.NdpConfig.RaRestransmitTime = uint8(dbEntry.RouterAdvertisementInterval)
		}
		if dbEntry.ReachableTime == 0 {
			debug.Logger.Warning("Invalid Reachable Interval and hence setting default value",
				NDP_DEFAULT_REACHABLE_INTERVAL)
			svr.NdpConfig.ReachableTime = NDP_DEFAULT_REACHABLE_INTERVAL
		} else {
			svr.NdpConfig.ReachableTime = uint32(dbEntry.ReachableTime)
		}
		if dbEntry.RetransmitInterval == 0 {
			debug.Logger.Warning("Invalid ReTransmit Interval and hence setting default value",
				NDP_DEFAULT_RETRANSMIT_INTERVAL)
			svr.NdpConfig.RetransTime = NDP_DEFAULT_RETRANSMIT_INTERVAL
		} else {
			svr.NdpConfig.RetransTime = uint32(dbEntry.RetransmitInterval)
		}
		debug.Logger.Info("Done with reading NDPGlobal config from DB")
	}
}

func (svr *NDPServer) ReadDB() {
	if svr.dmnBase == nil {
		return
	}
	dbHdl := svr.dmnBase.DbHdl
	if dbHdl == nil {
		debug.Logger.Err("DB Handler is nil and hence cannot read anything from DATABASE")
		return
	}
	debug.Logger.Info("Reading Config from DB")
	svr.readNdpGblCfg(dbHdl)
}
