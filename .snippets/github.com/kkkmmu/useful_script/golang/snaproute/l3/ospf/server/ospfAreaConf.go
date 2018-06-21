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
	"l3/ospf/config"
	"time"
)

type AreaConfKey struct {
	AreaId config.AreaId
}

/* TODO - Add list of interfaces for this Area */
type AreaConf struct {
	AuthType          config.AuthType
	ImportAsExtern    config.ImportAsExtern
	AreaSummary       config.AreaSummary
	StubDefaultCost   int32
	TransitCapability bool
	IntfListMap       map[IntfConfKey]bool
}

type AreaState struct {
	SpfRuns                  int32
	AreaBdrRtrCount          int32
	AsBdrRtrCount            int32
	AreaLsaCount             int32
	AreaLsaCksumSum          int32
	AreaNssaTranslatorState  config.NssaTranslatorState
	AreaNssaTranslatorEvents int32
}

func (server *OSPFServer) processAreaConfig(areaConf config.AreaConf) error {
	areaConfKey := AreaConfKey{
		AreaId: areaConf.AreaId,
	}

	ent, _ := server.AreaConfMap[areaConfKey]
	ent.AuthType = areaConf.AuthType
	ent.ImportAsExtern = areaConf.ImportAsExtern
	ent.AreaSummary = areaConf.AreaSummary
	ent.StubDefaultCost = areaConf.StubDefaultCost
	ent.IntfListMap = make(map[IntfConfKey]bool)
	server.AreaConfMap[areaConfKey] = ent
	server.initAreaStateSlice(areaConfKey)
	areaId := convertAreaOrRouterIdUint32(string(areaConf.AreaId))
	server.initLSDatabase(areaId)
	server.initRoutingTbl(areaId)
	return nil
}

func (server *OSPFServer) initAreaConfDefault() {
	server.logger.Info("Initializing default area config")
	areaConfKey := AreaConfKey{
		AreaId: "0.0.0.0",
	}
	ent, exist := server.AreaConfMap[areaConfKey]
	if !exist {
		ent.AuthType = config.NoAuth
		ent.ImportAsExtern = config.ImportExternal
		ent.AreaSummary = config.NoAreaSummary
		ent.StubDefaultCost = 20
		ent.IntfListMap = make(map[IntfConfKey]bool)
		server.AreaConfMap[areaConfKey] = ent
	}
	server.initAreaStateSlice(areaConfKey)
	server.areaStateRefresh()
}

func (server *OSPFServer) initAreaStateSlice(key AreaConfKey) {
	//server.AreaStateMutex.Lock()
	server.logger.Debug(fmt.Sprintln("Initializing area slice", key))
	ent, exist := server.AreaStateMap[key]
	ent.SpfRuns = 0
	ent.AreaBdrRtrCount = 0
	ent.AsBdrRtrCount = 0
	ent.AreaLsaCount = 0
	ent.AreaLsaCksumSum = 0
	ent.AreaNssaTranslatorState = config.NssaTranslatorDisabled
	ent.AreaNssaTranslatorEvents = 0
	server.AreaStateMap[key] = ent
	if !exist {
		server.AreaStateSlice = append(server.AreaStateSlice, key)
		server.AreaConfKeyToSliceIdxMap[key] = len(server.AreaStateSlice) - 1
	}
	//server.AreaStateMutex.Unlock()
}

func (server *OSPFServer) areaStateRefresh() {
	var areaStateRefFunc func()
	areaStateRefFunc = func() {
		//server.AreaStateMutex.Lock()
		server.AreaStateSlice = []AreaConfKey{}
		server.AreaConfKeyToSliceIdxMap = nil
		server.AreaConfKeyToSliceIdxMap = make(map[AreaConfKey]int)
		for key, _ := range server.AreaStateMap {
			server.AreaStateSlice = append(server.AreaStateSlice, key)
			server.AreaConfKeyToSliceIdxMap[key] = len(server.AreaStateSlice) - 1
		}
		//server.AreaStateMutex.Unlock()
		server.AreaStateTimer.Reset(server.RefreshDuration)
	}
	server.AreaStateTimer = time.AfterFunc(server.RefreshDuration, areaStateRefFunc)
}

func (server *OSPFServer) updateIntfToAreaMap(key IntfConfKey, oldAreaId string, newAreaId string, enable bool) {

	server.logger.Debug(fmt.Sprintln("===========1. updateIntfToAreaMap============", server.AreaConfMap, "oldAreaId:", oldAreaId, "newAreaId:", newAreaId, "IntfConfKey:", key))
	if oldAreaId != "none" && newAreaId != "none" {
		oldAreaConfKey := AreaConfKey{
			AreaId: config.AreaId(oldAreaId),
		}
		oldAreaConfEnt, exist := server.AreaConfMap[oldAreaConfKey]
		if !exist {
			server.logger.Err(fmt.Sprintln("Update: No such area configuration exist.", oldAreaConfKey))
			return
		}

		delete(oldAreaConfEnt.IntfListMap, key)
		server.AreaConfMap[oldAreaConfKey] = oldAreaConfEnt

		newAreaConfKey := AreaConfKey{
			AreaId: config.AreaId(newAreaId),
		}
		newAreaConfEnt, exist := server.AreaConfMap[newAreaConfKey]
		if !exist {
			server.logger.Err(fmt.Sprintln("Update: No such area configuration exist.", newAreaConfKey))
			return
		}
		if enable {
			newAreaConfEnt.IntfListMap[key] = true
		} else {
			newAreaConfEnt.IntfListMap[key] = false
		}
		server.AreaConfMap[newAreaConfKey] = newAreaConfEnt
	} else if oldAreaId == "none" {
		newAreaConfKey := AreaConfKey{
			AreaId: config.AreaId(newAreaId),
		}
		newAreaConfEnt, exist := server.AreaConfMap[newAreaConfKey]
		if !exist {
			server.logger.Err(fmt.Sprintln("Update: No such area configuration exist.", newAreaConfKey))
			return
		}

		newAreaConfEnt.IntfListMap[key] = true
		server.AreaConfMap[newAreaConfKey] = newAreaConfEnt

	} else if newAreaId == "none" {
		oldAreaConfKey := AreaConfKey{
			AreaId: config.AreaId(oldAreaId),
		}
		oldAreaConfEnt, exist := server.AreaConfMap[oldAreaConfKey]
		if !exist {
			server.logger.Err("No such area configuration exist.")
			return
		}

		delete(oldAreaConfEnt.IntfListMap, key)
		server.AreaConfMap[oldAreaConfKey] = oldAreaConfEnt
	} else {
		server.logger.Err("Invalid Argument for updating Intf List in Area Conf Map")
		return
	}
	server.updateIfABR()
	server.logger.Info(fmt.Sprintln("AreaConf Map:", server.AreaConfMap, "AreaBdr Status:", server.ospfGlobalConf.isABR))
}

func (server *OSPFServer) updateIfABR() {
	index := 0
	var validIntf bool
	for _, areaEnt := range server.AreaConfMap {
		validIntf = false
		for _, enable := range areaEnt.IntfListMap {
			if enable {
				validIntf = true
				break
			}
		}
		if validIntf {
			index++
		}
	}
	if index > 1 {
		server.ospfGlobalConf.isABR = true
		server.ospfGlobalConf.AreaBdrRtrStatus = true
	} else {
		server.ospfGlobalConf.isABR = false
		server.ospfGlobalConf.AreaBdrRtrStatus = false
	}
}

func (server *OSPFServer) isStubArea(areaid config.AreaId) bool {

	areaConfKey := AreaConfKey{
		AreaId: areaid,
	}

	conf, exist := server.AreaConfMap[areaConfKey]
	if !exist {
		return false
	}
	if conf.ImportAsExtern == config.ImportNoExternal {
		return true
	}
	return false
}
