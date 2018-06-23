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

// server.go
package rpc

import (
	"bgpd"
	"encoding/json"
	"errors"
	"fmt"
	"l3/bgp/config"
	"l3/bgp/fsm"
	"l3/bgp/packet"
	bgppolicy "l3/bgp/policy"
	"l3/bgp/server"
	bgputils "l3/bgp/utils"
	"math"
	"models/objects"
	"net"
	"reflect"
	"strings"
	"time"
	"utils/dbutils"
	"utils/logging"
)

const DBName string = "UsrConfDb.db"

type PeerConfigCommands struct {
	IP      net.IP
	Command int
}

type BGPHandler struct {
	PeerCommandCh chan PeerConfigCommands
	server        *server.BGPServer
	bgpPolicyMgr  *bgppolicy.BGPPolicyManager
	logger        *logging.Writer
	dbUtil        *dbutils.DBUtil
	globalASMap   map[string]uint32
}

func NewBGPHandler(server *server.BGPServer, policyMgr *bgppolicy.BGPPolicyManager, logger *logging.Writer,
	dbUtil *dbutils.DBUtil, filePath string) *BGPHandler {
	h := new(BGPHandler)
	h.PeerCommandCh = make(chan PeerConfigCommands)
	h.server = server
	h.bgpPolicyMgr = policyMgr
	h.logger = logger
	h.dbUtil = dbUtil
	h.globalASMap = make(map[string]uint32)
	return h
}

func (h *BGPHandler) convertModelToBGPGlobalConfig(obj objects.BGPGlobal) (gConf config.GlobalConfig, err error) {
	asnum, err := bgputils.GetAsNum(obj.ASNum)
	if err != nil {
		h.logger.Err("Invalid asnum")
		return gConf, err
	}

	gConf = config.GlobalConfig{
		GlobalBase: config.GlobalBase{
			Vrf:                 obj.Vrf,
			AS:                  uint32(asnum),
			RouterId:            h.convertStrIPToNetIP(obj.RouterId),
			Disabled:            obj.Disabled,
			UseMultiplePaths:    obj.UseMultiplePaths,
			EBGPMaxPaths:        obj.EBGPMaxPaths,
			EBGPAllowMultipleAS: obj.EBGPAllowMultipleAS,
			IBGPMaxPaths:        obj.IBGPMaxPaths,
			Defaultv4Route:      obj.Defaultv4Route,
			Defaultv6Route:      obj.Defaultv6Route,
			DefaultMED:          obj.DefaultMED,
		},
	}

	if obj.Redistribution != nil {
		gConf.Redistribution = make([]config.SourcePolicyMap, 0)
		for i := 0; i < len(obj.Redistribution); i++ {
			redistribution := config.SourcePolicyMap{obj.Redistribution[i].Sources, obj.Redistribution[i].Policy}
			gConf.Redistribution = append(gConf.Redistribution, redistribution)
		}
	}

	if gConf.RouterId == nil {
		h.logger.Err("convertModelToBGPGlobalConfig - IP is not valid:", obj.RouterId)
		err = config.IPError{obj.RouterId}
	}

	return gConf, err
}

func (h *BGPHandler) handleGlobalConfig() error {
	var obj objects.BGPGlobal
	objList, err := h.dbUtil.GetAllObjFromDb(obj)
	if err != nil {
		h.logger.Errf("GetAllObjFromDb failed for BGPGlobal with error %s", err)
		return err
	}

	for _, confObj := range objList {
		obj = confObj.(objects.BGPGlobal)

		gConf, err := h.convertModelToBGPGlobalConfig(obj)
		if err != nil {
			h.logger.Err("handleGlobalConfig - Failed to convert Model object BGP Global, error:", err)
			return err
		}
		h.globalASMap[gConf.Vrf] = gConf.AS
		h.server.GlobalConfigCh <- server.GlobalUpdate{nil, config.GlobalConfig{}, gConf, make([]bool, 0), nil, "create"}
	}
	return nil
}

func (h *BGPHandler) convertModelToBGPv4PeerGroup(obj objects.BGPv4PeerGroup) (group config.PeerGroupConfig,
	err error) {
	peerAS, err := bgputils.GetAsNum(obj.PeerAS)
	if err != nil {
		h.logger.Err("Invalid peer asnum")
		return group, err
	}
	localAS, err := bgputils.GetAsNum(obj.LocalAS)
	if err != nil {
		h.logger.Err("Invalid local asnum")
		return group, err
	}
	group = config.PeerGroupConfig{
		BaseConfig: config.BaseConfig{
			PeerAS:                  uint32(peerAS),
			LocalAS:                 uint32(localAS),
			PeerAddressType:         config.PeerAddressV4,
			UpdateSource:            obj.UpdateSource,
			NextHopSelf:             obj.NextHopSelf,
			AuthPassword:            obj.AuthPassword,
			Description:             obj.Description,
			RouteReflectorClusterId: uint32(obj.RouteReflectorClusterId),
			RouteReflectorClient:    obj.RouteReflectorClient,
			MultiHopEnable:          obj.MultiHopEnable,
			MultiHopTTL:             uint8(obj.MultiHopTTL),
			ConnectRetryTime:        uint32(obj.ConnectRetryTime),
			HoldTime:                uint32(obj.HoldTime),
			KeepaliveTime:           uint32(obj.KeepaliveTime),
			AddPathsRx:              obj.AddPathsRx,
			AddPathsMaxTx:           uint8(obj.AddPathsMaxTx),
			MaxPrefixes:             uint32(obj.MaxPrefixes),
			MaxPrefixesThresholdPct: uint8(obj.MaxPrefixesThresholdPct),
			MaxPrefixesDisconnect:   obj.MaxPrefixesDisconnect,
			MaxPrefixesRestartTimer: uint8(obj.MaxPrefixesRestartTimer),
			AdjRIBInFilter:          obj.AdjRIBInFilter,
			AdjRIBOutFilter:         obj.AdjRIBOutFilter,
		},
		Name: obj.Name,
	}
	return group, err
}

func (h *BGPHandler) handleV4PeerGroup() error {
	var obj objects.BGPv4PeerGroup
	objList, err := h.dbUtil.GetAllObjFromDb(obj)
	if err != nil {
		h.logger.Errf("GetAllObjFromDb for BGPv4PeerGroup failed with error %s", err)
		return err
	}

	for _, confObj := range objList {
		obj = confObj.(objects.BGPv4PeerGroup)

		group, err := h.convertModelToBGPv4PeerGroup(obj)
		if err != nil {
			h.logger.Err("handlePeerGroup - Failed to convert Model object to BGP Peer group, error:", err)
			return err
		}

		h.server.AddPeerGroupCh <- server.PeerGroupUpdate{config.PeerGroupConfig{}, group, make([]bool, 0)}
	}

	return nil
}

func (h *BGPHandler) convertModelToBGPv6PeerGroup(obj objects.BGPv6PeerGroup) (group config.PeerGroupConfig,
	err error) {
	peerAS, err := bgputils.GetAsNum(obj.PeerAS)
	if err != nil {
		h.logger.Err("Invalid peer asnum")
		return group, err
	}
	localAS, err := bgputils.GetAsNum(obj.LocalAS)
	if err != nil {
		h.logger.Err("Invalid local asnum")
		return group, err
	}
	group = config.PeerGroupConfig{
		BaseConfig: config.BaseConfig{
			PeerAS:                  uint32(peerAS),
			LocalAS:                 uint32(localAS),
			PeerAddressType:         config.PeerAddressV6,
			UpdateSource:            obj.UpdateSource,
			NextHopSelf:             obj.NextHopSelf,
			Description:             obj.Description,
			RouteReflectorClusterId: uint32(obj.RouteReflectorClusterId),
			RouteReflectorClient:    obj.RouteReflectorClient,
			MultiHopEnable:          obj.MultiHopEnable,
			MultiHopTTL:             uint8(obj.MultiHopTTL),
			ConnectRetryTime:        uint32(obj.ConnectRetryTime),
			HoldTime:                uint32(obj.HoldTime),
			KeepaliveTime:           uint32(obj.KeepaliveTime),
			AddPathsRx:              obj.AddPathsRx,
			AddPathsMaxTx:           uint8(obj.AddPathsMaxTx),
			MaxPrefixes:             uint32(obj.MaxPrefixes),
			MaxPrefixesThresholdPct: uint8(obj.MaxPrefixesThresholdPct),
			MaxPrefixesDisconnect:   obj.MaxPrefixesDisconnect,
			MaxPrefixesRestartTimer: uint8(obj.MaxPrefixesRestartTimer),
			AdjRIBInFilter:          obj.AdjRIBInFilter,
			AdjRIBOutFilter:         obj.AdjRIBOutFilter,
		},
		Name: obj.Name,
	}
	return group, err
}

func (h *BGPHandler) handleV6PeerGroup() error {
	var obj objects.BGPv6PeerGroup
	objList, err := h.dbUtil.GetAllObjFromDb(obj)
	if err != nil {
		h.logger.Errf("GetAllObjFromDb for BGPv6PeerGroup failed with error %s", err)
		return err
	}

	for _, confObj := range objList {
		obj = confObj.(objects.BGPv6PeerGroup)

		group, err := h.convertModelToBGPv6PeerGroup(obj)
		if err != nil {
			h.logger.Err("handlePeerGroup - Failed to convert Model object to BGP Peer group, error:", err)
			return err
		}

		h.server.AddPeerGroupCh <- server.PeerGroupUpdate{config.PeerGroupConfig{}, group, make([]bool, 0)}
	}

	return nil
}

func (h *BGPHandler) convertModelToBGPv4Neighbor(obj objects.BGPv4Neighbor) (neighbor config.NeighborConfig,
	err error) {
	var ip net.IP
	var ifIndex int32
	ip, ifIndex, _, err = h.getIPAndIfIndexForV4Neighbor(obj.NeighborAddress, obj.IntfRef)
	if err != nil {
		h.logger.Info("convertModelToBGPv4Neighbor: getIPAndIfIndexForV4Neighbor",
			"failed for neighbor address", obj.NeighborAddress, "and ifIndex", obj.IntfRef)
		return neighbor, err
	}

	peerAS, err := bgputils.GetAsNum(obj.PeerAS)
	if err != nil {
		h.logger.Err("Invalid peer asnum")
		return neighbor, err
	}
	localAS, err := bgputils.GetAsNum(obj.LocalAS)
	if err != nil {
		h.logger.Err("Invalid local asnum")
		return neighbor, err
	}

	neighbor = config.NeighborConfig{
		BaseConfig: config.BaseConfig{
			PeerAS:                  uint32(peerAS),
			LocalAS:                 uint32(localAS),
			PeerAddressType:         config.PeerAddressV4,
			UpdateSource:            obj.UpdateSource,
			NextHopSelf:             obj.NextHopSelf,
			AuthPassword:            obj.AuthPassword,
			Description:             obj.Description,
			RouteReflectorClusterId: uint32(obj.RouteReflectorClusterId),
			RouteReflectorClient:    obj.RouteReflectorClient,
			MultiHopEnable:          obj.MultiHopEnable,
			MultiHopTTL:             uint8(obj.MultiHopTTL),
			ConnectRetryTime:        uint32(obj.ConnectRetryTime),
			HoldTime:                uint32(obj.HoldTime),
			KeepaliveTime:           uint32(obj.KeepaliveTime),
			BfdEnable:               obj.BfdEnable,
			BfdSessionParam:         obj.BfdSessionParam,
			AddPathsRx:              obj.AddPathsRx,
			AddPathsMaxTx:           uint8(obj.AddPathsMaxTx),
			MaxPrefixes:             uint32(obj.MaxPrefixes),
			MaxPrefixesThresholdPct: uint8(obj.MaxPrefixesThresholdPct),
			MaxPrefixesDisconnect:   obj.MaxPrefixesDisconnect,
			MaxPrefixesRestartTimer: uint8(obj.MaxPrefixesRestartTimer),
			AdjRIBInFilter:          obj.AdjRIBInFilter,
			AdjRIBOutFilter:         obj.AdjRIBOutFilter,
		},
		NeighborAddress: ip,
		IfIndex:         ifIndex,
		PeerGroup:       obj.PeerGroup,
		Disabled:        obj.Disabled,
	}
	return neighbor, err
}

func (h *BGPHandler) handleV4NeighborConfig() error {
	var obj objects.BGPv4Neighbor
	objList, err := h.dbUtil.GetAllObjFromDb(obj)
	if err != nil {
		h.logger.Errf("GetAllObjFromDb for BGPNeighbor failed with error %s", err)
		return err
	}

	for _, confObj := range objList {
		obj = confObj.(objects.BGPv4Neighbor)

		neighbor, err := h.convertModelToBGPv4Neighbor(obj)
		if err != nil {
			h.logger.Err("handleV4NeighborConfig - Failed to convert Model object to BGP neighbor, error:", err)
			return err
		}

		h.server.AddPeerCh <- server.PeerUpdate{nil, "v4", config.NeighborConfig{}, neighbor, make([]bool, 0), nil, "create"}
	}

	return nil
}

func (h *BGPHandler) convertModelToBGPv6Neighbor(obj objects.BGPv6Neighbor) (neighbor config.NeighborConfig,
	err error) {
	var ip net.IP
	var ifIndex int32
	var ifName string
	ip, ifIndex, ifName, err = h.getIPAndIfIndexForV6Neighbor(obj.NeighborAddress, obj.IntfRef)
	if err != nil {
		h.logger.Info("convertModelToBGPv6Neighbor: getIPAndIfIndexForV6Neighbor",
			"failed for neighbor address", obj.NeighborAddress, "and ifIndex", obj.IntfRef)
		return neighbor, err
	}

	peerAS, err := bgputils.GetAsNum(obj.PeerAS)
	if err != nil {
		h.logger.Err("Invalid peer asnum")
		return neighbor, err
	}
	localAS, err := bgputils.GetAsNum(obj.LocalAS)
	if err != nil {
		h.logger.Err("Invalid local asnum")
		return neighbor, err
	}

	neighbor = config.NeighborConfig{
		BaseConfig: config.BaseConfig{
			PeerAS:                  uint32(peerAS),
			LocalAS:                 uint32(localAS),
			PeerAddressType:         config.PeerAddressV6,
			UpdateSource:            obj.UpdateSource,
			NextHopSelf:             obj.NextHopSelf,
			Description:             obj.Description,
			RouteReflectorClusterId: uint32(obj.RouteReflectorClusterId),
			RouteReflectorClient:    obj.RouteReflectorClient,
			MultiHopEnable:          obj.MultiHopEnable,
			MultiHopTTL:             uint8(obj.MultiHopTTL),
			ConnectRetryTime:        uint32(obj.ConnectRetryTime),
			HoldTime:                uint32(obj.HoldTime),
			KeepaliveTime:           uint32(obj.KeepaliveTime),
			BfdEnable:               obj.BfdEnable,
			BfdSessionParam:         obj.BfdSessionParam,
			AddPathsRx:              obj.AddPathsRx,
			AddPathsMaxTx:           uint8(obj.AddPathsMaxTx),
			MaxPrefixes:             uint32(obj.MaxPrefixes),
			MaxPrefixesThresholdPct: uint8(obj.MaxPrefixesThresholdPct),
			MaxPrefixesDisconnect:   obj.MaxPrefixesDisconnect,
			MaxPrefixesRestartTimer: uint8(obj.MaxPrefixesRestartTimer),
			AdjRIBInFilter:          obj.AdjRIBInFilter,
			AdjRIBOutFilter:         obj.AdjRIBOutFilter,
		},
		NeighborAddress: ip,
		IfIndex:         ifIndex,
		IfName:          ifName,
		PeerGroup:       obj.PeerGroup,
		Disabled:        obj.Disabled,
	}
	return neighbor, err
}

func (h *BGPHandler) handleV6NeighborConfig() error {
	var obj objects.BGPv6Neighbor
	objList, err := h.dbUtil.GetAllObjFromDb(obj)
	if err != nil {
		h.logger.Errf("GetAllObjFromDb for BGPNeighbor failed with error %s", err)
		return err
	}

	for _, confObj := range objList {
		obj = confObj.(objects.BGPv6Neighbor)

		neighbor, err := h.convertModelToBGPv6Neighbor(obj)
		if err != nil {
			h.logger.Err("handleV6NeighborConfig - Failed to convert Model object to BGP neighbor, error:", err)
			return err
		}

		h.server.AddPeerCh <- server.PeerUpdate{nil, "v6", config.NeighborConfig{}, neighbor, make([]bool, 0), nil, "create"}
	}

	return nil
}

func (h *BGPHandler) convertModelToBGPv4Aggregate(obj objects.BGPv4Aggregate) (config.BGPAggregate, error) {
	aggConf := config.BGPAggregate{
		IPPrefix:        obj.IpPrefix,
		GenerateASSet:   obj.GenerateASSet,
		SendSummaryOnly: obj.SendSummaryOnly,
		AddressFamily:   packet.GetProtocolFamily(packet.AfiIP, packet.SafiUnicast),
	}

	return aggConf, nil
}

func (h *BGPHandler) handleBGPv4Aggregate() error {
	var obj objects.BGPv4Aggregate
	objList, err := h.dbUtil.GetAllObjFromDb(obj)
	if err != nil {
		h.logger.Errf("GetAllObjFromDb failed for BGPv4Aggregate with error %s", err)
		return err
	}

	for _, confObj := range objList {
		obj = confObj.(objects.BGPv4Aggregate)

		aggConf, err := h.convertModelToBGPv4Aggregate(obj)
		if err != nil {
			h.logger.Err("handleBGPv4Aggregate - Failed to convert Model object BGPv4Aggregate, error:", err)
			return err
		}
		h.server.AddAggCh <- server.AggUpdate{config.BGPAggregate{}, aggConf, make([]bool, 0)}
	}
	return nil
}

func (h *BGPHandler) convertModelToBGPv6Aggregate(obj objects.BGPv6Aggregate) (config.BGPAggregate, error) {
	aggConf := config.BGPAggregate{
		IPPrefix:        obj.IpPrefix,
		GenerateASSet:   obj.GenerateASSet,
		SendSummaryOnly: obj.SendSummaryOnly,
		AddressFamily:   packet.GetProtocolFamily(packet.AfiIP6, packet.SafiUnicast),
	}

	return aggConf, nil
}

func (h *BGPHandler) handleBGPv6Aggregate() error {
	var obj objects.BGPv6Aggregate
	objList, err := h.dbUtil.GetAllObjFromDb(obj)
	if err != nil {
		h.logger.Errf("GetAllObjFromDb failed for BGPv6Aggregate with error %s", err)
		return err
	}

	for _, confObj := range objList {
		obj = confObj.(objects.BGPv6Aggregate)

		aggConf, err := h.convertModelToBGPv6Aggregate(obj)
		if err != nil {
			h.logger.Err("handleBGPv6Aggregate - Failed to convert Model object BGPv6Aggregate, error:", err)
			return err
		}
		h.server.AddAggCh <- server.AggUpdate{config.BGPAggregate{}, aggConf, make([]bool, 0)}
	}
	return nil
}

func (h *BGPHandler) ReadBGPConfigFromDB() error {
	var err error
	if err = h.handleGlobalConfig(); err != nil {
		return err
	}

	if err = h.handleBGPv4Aggregate(); err != nil {
		return err
	}

	if err = h.handleBGPv6Aggregate(); err != nil {
		return err
	}

	if err = h.handleV4PeerGroup(); err != nil {
		return err
	}
	if err = h.handleV4NeighborConfig(); err != nil {
		return err
	}

	if err = h.handleV6PeerGroup(); err != nil {
		return err
	}

	if err = h.handleV6NeighborConfig(); err != nil {
		return err
	}

	return nil
}

func (h *BGPHandler) convertStrIPToNetIP(ip string) net.IP {
	if ip == "localhost" {
		ip = "127.0.0.1"
	}

	netIP := net.ParseIP(ip)
	return netIP
}

func (h *BGPHandler) validateBGPGlobal(bgpGlobal *bgpd.BGPGlobal) (gConf config.GlobalConfig, err error) {
	if bgpGlobal == nil {
		return gConf, err
	}

	asNum, err := bgputils.GetAsNum(bgpGlobal.ASNum) //uint32(bgpGlobal.ASNum)
	if err != nil {
		return gConf, err
	}
	h.logger.Info("Setting ASnum = ", asNum, " for ", bgpGlobal.ASNum)
	if asNum == (math.MaxUint16) || asNum == (math.MaxUint32) || asNum == int(packet.BGPASTrans) {
		err = errors.New(fmt.Sprintf("BGPGlobal: AS number %d is not valid", bgpGlobal.ASNum))
		h.logger.Info("SendBGPGlobal: AS number", bgpGlobal.ASNum, "is a reserved AS number")
		return gConf, err
	}

	ip := h.convertStrIPToNetIP(bgpGlobal.RouterId)
	if ip == nil {
		err = errors.New(fmt.Sprintf("BGPGlobal: Router id %s is not valid", bgpGlobal.RouterId))
		h.logger.Info("SendBGPGlobal: Router id", bgpGlobal.RouterId, "is not valid")
		return gConf, err
	}

	gConf = config.GlobalConfig{
		GlobalBase: config.GlobalBase{
			Vrf:                 bgpGlobal.Vrf,
			AS:                  uint32(asNum),
			RouterId:            ip,
			Disabled:            bgpGlobal.Disabled,
			UseMultiplePaths:    bgpGlobal.UseMultiplePaths,
			EBGPMaxPaths:        uint32(bgpGlobal.EBGPMaxPaths),
			EBGPAllowMultipleAS: bgpGlobal.EBGPAllowMultipleAS,
			IBGPMaxPaths:        uint32(bgpGlobal.IBGPMaxPaths),
			Defaultv4Route:      bgpGlobal.Defaultv4Route,
			Defaultv6Route:      bgpGlobal.Defaultv6Route,
			DefaultMED:          uint32(bgpGlobal.DefaultMED),
		},
	}

	if bgpGlobal.Redistribution != nil {
		gConf.Redistribution = make([]config.SourcePolicyMap, 0)
		for i := 0; i < len(bgpGlobal.Redistribution); i++ {
			redistribution := config.SourcePolicyMap{bgpGlobal.Redistribution[i].Sources, bgpGlobal.Redistribution[i].Policy}
			gConf.Redistribution = append(gConf.Redistribution, redistribution)
		}
	}

	return gConf, nil
}

func (h *BGPHandler) validateBGPGlobalForPatchUpdate(oldConfig *bgpd.BGPGlobal, newConfig *bgpd.BGPGlobal, op []*bgpd.PatchOpInfo) (gConf config.GlobalConfig, err error) {
	h.logger.Info("validateBGPGlobalForPatchUpdate")
	if oldConfig == nil || newConfig == nil {
		err = errors.New(fmt.Sprintf("validateBGPGlobalForUpdate: oldConfig %+v or newConfig %+v is nil", oldConfig.RouterId))
		return gConf, err
	}

	ip := h.convertStrIPToNetIP(oldConfig.RouterId)
	if ip == nil {
		err = errors.New(fmt.Sprintf("BGPGlobal: Router id %s is not valid", oldConfig.RouterId))
		h.logger.Info("SendBGPGlobal: Router id", oldConfig.RouterId, "is not valid")
		return gConf, err
	}

	oldAsnum, err := bgputils.GetAsNum(oldConfig.ASNum)
	if err != nil {
		return gConf, err
	}

	gConf = config.GlobalConfig{
		GlobalBase: config.GlobalBase{
			Vrf:                 oldConfig.Vrf,
			AS:                  uint32(oldAsnum),
			RouterId:            ip,
			Disabled:            oldConfig.Disabled,
			UseMultiplePaths:    oldConfig.UseMultiplePaths,
			EBGPMaxPaths:        uint32(oldConfig.EBGPMaxPaths),
			EBGPAllowMultipleAS: oldConfig.EBGPAllowMultipleAS,
			IBGPMaxPaths:        uint32(oldConfig.IBGPMaxPaths),
			Defaultv4Route:      oldConfig.Defaultv4Route,
			Defaultv6Route:      oldConfig.Defaultv6Route,
			DefaultMED:          uint32(oldConfig.DefaultMED),
		},
	}

	for idx := 0; idx < len(op); idx++ {
		h.logger.Debug("patch update")
		switch op[idx].Path {
		case "Redistribution":
			h.logger.Debug("Patch update for redistribution")
			if len(op[idx].Value) == 0 {
				/*
					If redistribution update is trying to update redistribution, non zero value is expected
				*/
				h.logger.Err("Must specify sources")
				return gConf, errors.New("Redistribution update list not specified")
			}
			h.logger.Debug("value = ", op[idx].Value)
			valueObjArr := []bgpd.SourcePolicyList{}
			err = json.Unmarshal([]byte(op[idx].Value), &valueObjArr)
			if err != nil {
				h.logger.Err("error unmarshaling value:", err)
				return gConf, errors.New(fmt.Sprintln("error unmarshaling value:", err))
			}
			h.logger.Debug("Number of redistribution soures:", len(valueObjArr))
			for _, _ = range valueObjArr {
				switch op[idx].Op {
				case "add":
					h.logger.Debug("add op")
				case "remove":
					h.logger.Debug("remove op")
				default:
					h.logger.Err("operation ", op[idx].Op, " not supported")
					return gConf, errors.New(fmt.Sprintln("operation ", op[idx].Op, " not supported"))
				}
			}
		default:
			h.logger.Err("Patch update for attribute:", op[idx].Path, " not supported")
			return gConf, errors.New("Invalid attribute for patch update")
		}
	}
	return gConf, err
}

func (h *BGPHandler) validateBGPGlobalForUpdate(oldConfig *bgpd.BGPGlobal, newConfig *bgpd.BGPGlobal, attrSet []bool) (gConf config.GlobalConfig, err error) {
	h.logger.Info("validateBGPGlobalForUpdate")
	if oldConfig == nil || newConfig == nil {
		err = errors.New(fmt.Sprintf("validateBGPGlobalForUpdate: oldConfig %+v or newConfig %+v is nil", oldConfig.RouterId))
		return gConf, err
	}

	ip := h.convertStrIPToNetIP(newConfig.RouterId)
	if ip == nil {
		err = errors.New(fmt.Sprintf("BGPGlobal: Router id %s is not valid", oldConfig.RouterId))
		h.logger.Info("SendBGPGlobal: Router id", oldConfig.RouterId, "is not valid")
		return gConf, err
	}

	newASNum, err := bgputils.GetAsNum(newConfig.ASNum)
	if err != nil {
		return gConf, err
	}

	gConf = config.GlobalConfig{
		GlobalBase: config.GlobalBase{
			Vrf:                 newConfig.Vrf,
			AS:                  uint32(newASNum),
			RouterId:            ip,
			Disabled:            newConfig.Disabled,
			UseMultiplePaths:    newConfig.UseMultiplePaths,
			EBGPMaxPaths:        uint32(newConfig.EBGPMaxPaths),
			EBGPAllowMultipleAS: newConfig.EBGPAllowMultipleAS,
			IBGPMaxPaths:        uint32(newConfig.IBGPMaxPaths),
			Defaultv4Route:      newConfig.Defaultv4Route,
			Defaultv6Route:      newConfig.Defaultv6Route,
			DefaultMED:          uint32(newConfig.DefaultMED),
		},
	}

	if newConfig.Redistribution != nil {
		gConf.Redistribution = make([]config.SourcePolicyMap, 0)
		for i := 0; i < len(newConfig.Redistribution); i++ {
			redistribution := config.SourcePolicyMap{newConfig.Redistribution[i].Sources, newConfig.Redistribution[i].Policy}
			gConf.Redistribution = append(gConf.Redistribution, redistribution)
		}
	}

	if attrSet != nil {
		objTyp := reflect.TypeOf(*newConfig)
		for i := 0; i < objTyp.NumField(); i++ {
			objName := objTyp.Field(i).Name
			if attrSet[i] {
				h.logger.Debug("validateBGPGlobalForUpdate : changed ", objName)
				if objName == "ASNum" {
					if (newASNum == (0)) || newASNum == (math.MaxUint16) ||
						newASNum == (math.MaxUint32) || newASNum == int(packet.BGPASTrans) {
						err = errors.New(fmt.Sprintf("BGPGlobal: AS number %d is not valid", newConfig.ASNum))
						h.logger.Info("SendBGPGlobal: AS number", newConfig.ASNum, "is not valid")
						return gConf, err
					}
					gConf.AS = uint32(newASNum)
				}
			}
		}
	}
	return gConf, err
}

func (h *BGPHandler) SendBGPGlobal(oldConfig *bgpd.BGPGlobal, newConfig *bgpd.BGPGlobal, attrSet []bool, patchOp []*bgpd.PatchOpInfo, op string) (bool, error) {
	var newGlobal config.GlobalConfig
	oldGlobal, err := h.validateBGPGlobal(oldConfig)
	if err != nil {
		return false, err
	}
	if op != "update" {
		newGlobal, err = h.validateBGPGlobal(newConfig)
		if err != nil {
			return false, err
		}
	} else {
		if patchOp == nil || len(patchOp) == 0 {
			newGlobal, err = h.validateBGPGlobalForUpdate(oldConfig, newConfig, attrSet)
			if err != nil {
				h.logger.Err("validateBGPGlobalForUpdate failed with err:", err)
				return false, err
			}
		} else {
			newGlobal, err = h.validateBGPGlobalForPatchUpdate(oldConfig, newConfig, patchOp)
			if err != nil {
				h.logger.Err("validateBGPGlobalForPatchUpdate failed with err:", err)
				return false, err
			}
		}
	}

	h.globalASMap[newGlobal.Vrf] = newGlobal.AS
	h.server.GlobalConfigCh <- server.GlobalUpdate{oldConfig, oldGlobal, newGlobal, attrSet, patchOp, op}
	return true, err
}

func (h *BGPHandler) CreateBGPGlobal(bgpGlobal *bgpd.BGPGlobal) (bool, error) {
	h.logger.Info("Create global config attrs:", bgpGlobal)
	return h.SendBGPGlobal(nil, bgpGlobal, make([]bool, 0), make([]*bgpd.PatchOpInfo, 0), "create")
}

func (h *BGPHandler) GetBGPGlobalState(vrfId string) (*bgpd.BGPGlobalState, error) {
	bgpGlobal := h.server.GetBGPGlobalState()
	bgpGlobalResponse := bgpd.NewBGPGlobalState()
	bgpGlobalResponse.Vrf = bgpGlobal.Vrf
	bgpGlobalResponse.AS, _ = bgputils.GetAsDot(int(bgpGlobal.AS)) //int32(bgpGlobal.AS)
	bgpGlobalResponse.RouterId = bgpGlobal.RouterId.String()
	bgpGlobalResponse.Disabled = bgpGlobal.Disabled
	bgpGlobalResponse.UseMultiplePaths = bgpGlobal.UseMultiplePaths
	bgpGlobalResponse.EBGPMaxPaths = int32(bgpGlobal.EBGPMaxPaths)
	bgpGlobalResponse.EBGPAllowMultipleAS = bgpGlobal.EBGPAllowMultipleAS
	bgpGlobalResponse.IBGPMaxPaths = int32(bgpGlobal.IBGPMaxPaths)
	bgpGlobalResponse.Defaultv4Route = bgpGlobal.Defaultv4Route
	bgpGlobalResponse.Defaultv6Route = bgpGlobal.Defaultv6Route
	bgpGlobalResponse.DefaultMED = int32(bgpGlobal.DefaultMED)
	bgpGlobalResponse.TotalPaths = int32(bgpGlobal.TotalPaths)
	bgpGlobalResponse.Totalv4Prefixes = int32(bgpGlobal.Totalv4Prefixes)
	bgpGlobalResponse.Totalv6Prefixes = int32(bgpGlobal.Totalv6Prefixes)
	return bgpGlobalResponse, nil
}

func (h *BGPHandler) GetBulkBGPGlobalState(index bgpd.Int,
	count bgpd.Int) (*bgpd.BGPGlobalStateGetInfo, error) {
	bgpGlobalStateBulk := bgpd.NewBGPGlobalStateGetInfo()
	bgpGlobalStateBulk.EndIdx = bgpd.Int(0)
	bgpGlobalStateBulk.Count = bgpd.Int(1)
	bgpGlobalStateBulk.More = false
	bgpGlobalStateBulk.BGPGlobalStateList = make([]*bgpd.BGPGlobalState, 1)
	bgpGlobalStateBulk.BGPGlobalStateList[0], _ = h.GetBGPGlobalState("default")

	return bgpGlobalStateBulk, nil
}

func (h *BGPHandler) UpdateBGPGlobal(origG *bgpd.BGPGlobal, updatedG *bgpd.BGPGlobal,
	attrSet []bool, op []*bgpd.PatchOpInfo) (bool, error) {
	h.logger.Info("Update global config attrs:", updatedG, "old config:", origG, " new config:", updatedG)
	return h.SendBGPGlobal(origG, updatedG, attrSet, op, "update")
}

func (h *BGPHandler) DeleteBGPGlobal(bgpGlobal *bgpd.BGPGlobal) (bool, error) {
	h.logger.Info("Delete global config attrs:", bgpGlobal)
	return false, errors.New(fmt.Sprintf("Can't delete BGP global object"))
}

func (h *BGPHandler) checkBGPGlobal() error {
	vrf := "default"

	if as, ok := h.globalASMap[vrf]; !ok || as == 0 {
		return errors.New(fmt.Sprintf("The default BGP AS number 0 is not updated yet."))
	}

	return nil
}

func (h *BGPHandler) getIPAndIfIndexForV4Neighbor(neighborIP string, neighborIntfRef string) (ip net.IP, ifIndex int32,
	ifName string, err error) {
	if strings.TrimSpace(neighborIP) != "" {
		ip = net.ParseIP(strings.TrimSpace(neighborIP))
		ifIndex = -1
		if ip == nil {
			err = errors.New(fmt.Sprintf("Neighbor address %s not valid", neighborIP))
		}
	} else if neighborIntfRef != "" {
		//neighbor address is a intfRef
		ifIndex, _, err = h.server.ConvertIntfStrToIfIndex(neighborIntfRef)
		if err != nil {
			h.logger.Err("Invalid intfref:", neighborIntfRef)
			return ip, ifIndex, ifName, err
		}

		var ipInfo *bgputils.IPInfo
		ipInfo, err = h.server.GetIfaceIP(ifIndex)
		h.logger.Info("getIPAndIfIndexForV4Neighbor - ipInfo:", ipInfo, " err:", err)
		if err == nil {
			if ipInfo.IpAddr == nil {
				return ip, ifIndex, ifName, errors.New(fmt.Sprint("IPv4 address is not set on interface ", neighborIntfRef))
			}

			ifIP := make(net.IP, len(ipInfo.IpAddr))
			copy(ifIP, ipInfo.IpAddr)
			ipMask := ipInfo.IpMask
			if ipMask[len(ipMask)-1] < 252 {
				h.logger.Err("IPv4Addr", ifIP, "of the interface", ifIndex, "is not /30 or /31 address")
				err = errors.New(fmt.Sprintln("IPv4Addr", ifIP, "of the interface", ifIndex,
					"is not /30 or /31 address"))
				return ip, ifIndex, ifName, err
			}
			h.logger.Info("IPv4Addr of the v4Neighbor local interface", ifIndex, "is", ifIP)
			ifIP[len(ifIP)-1] = ifIP[len(ifIP)-1] ^ (^ipMask[len(ipMask)-1])
			h.logger.Info("IPv4Addr of the v4Neighbor peer interface is", ifIP)
			ip = ifIP
		} else {
			h.logger.Err("v4Neighbor IP", neighborIP, "or interface", ifIndex, "not configured ")
		}
	}
	return ip, ifIndex, ifName, err
}

func (h *BGPHandler) isValidIP(ip string) bool {
	if strings.TrimSpace(ip) != "" {
		netIP := net.ParseIP(strings.TrimSpace(ip))
		if netIP == nil {
			return false
		}
	}

	return true
}

func (h *BGPHandler) ConvertV4NeighborFromThrift(bgpNeighbor *bgpd.BGPv4Neighbor, ip net.IP, ifIndex int32) (
	pConf config.NeighborConfig, err error) {

	peerAS, err := bgputils.GetAsNum(bgpNeighbor.PeerAS)
	if err != nil {
		h.logger.Err("peer AS invalid, err:", err)
		return pConf, err
	}
	localAS, err := bgputils.GetAsNum(bgpNeighbor.LocalAS)
	if err != nil {
		h.logger.Err("local AS invalid, err:", err)
		return pConf, err
	}
	pConf = config.NeighborConfig{
		BaseConfig: config.BaseConfig{
			PeerAS:                  uint32(peerAS),
			LocalAS:                 uint32(localAS),
			PeerAddressType:         config.PeerAddressV4,
			UpdateSource:            bgpNeighbor.UpdateSource,
			NextHopSelf:             bgpNeighbor.NextHopSelf,
			AuthPassword:            bgpNeighbor.AuthPassword,
			Description:             bgpNeighbor.Description,
			RouteReflectorClusterId: uint32(bgpNeighbor.RouteReflectorClusterId),
			RouteReflectorClient:    bgpNeighbor.RouteReflectorClient,
			MultiHopEnable:          bgpNeighbor.MultiHopEnable,
			MultiHopTTL:             uint8(bgpNeighbor.MultiHopTTL),
			ConnectRetryTime:        uint32(bgpNeighbor.ConnectRetryTime),
			HoldTime:                uint32(bgpNeighbor.HoldTime),
			KeepaliveTime:           uint32(bgpNeighbor.KeepaliveTime),
			BfdEnable:               bgpNeighbor.BfdEnable,
			BfdSessionParam:         bgpNeighbor.BfdSessionParam,
			AddPathsRx:              bgpNeighbor.AddPathsRx,
			AddPathsMaxTx:           uint8(bgpNeighbor.AddPathsMaxTx),
			MaxPrefixes:             uint32(bgpNeighbor.MaxPrefixes),
			MaxPrefixesThresholdPct: uint8(bgpNeighbor.MaxPrefixesThresholdPct),
			MaxPrefixesDisconnect:   bgpNeighbor.MaxPrefixesDisconnect,
			MaxPrefixesRestartTimer: uint8(bgpNeighbor.MaxPrefixesRestartTimer),
			AdjRIBInFilter:          bgpNeighbor.AdjRIBInFilter,
			AdjRIBOutFilter:         bgpNeighbor.AdjRIBOutFilter,
		},
		NeighborAddress: ip,
		IfIndex:         ifIndex,
		PeerGroup:       bgpNeighbor.PeerGroup,
		Disabled:        bgpNeighbor.Disabled,
	}
	return pConf, err
}

func (h *BGPHandler) ValidateV4Neighbor(bgpNeighbor *bgpd.BGPv4Neighbor) (pConf config.NeighborConfig, err error) {
	if bgpNeighbor == nil {
		return pConf, errors.New("NeighborConfig nil")
	}

	var ip net.IP
	var ifIndex int32
	ip, ifIndex, _, err = h.getIPAndIfIndexForV4Neighbor(bgpNeighbor.NeighborAddress, bgpNeighbor.IntfRef)
	if err != nil {
		h.logger.Info("ValidateV4Neighbor: getIPAndIfIndexForNeighbor failed for neighbor address",
			bgpNeighbor.NeighborAddress, "and ifIndex", bgpNeighbor.IntfRef)
		return pConf, err
	}
	h.logger.Info("ValidateV4Neighbor: getIPAndIfIndexForNeighbor returned ip", ip, "ifIndex", ifIndex)

	if !h.isValidIP(bgpNeighbor.UpdateSource) {
		err = errors.New(fmt.Sprintf("Update source %s not a valid IP", bgpNeighbor.UpdateSource))
		return pConf, err
	}
	pConf, _ = h.ConvertV4NeighborFromThrift(bgpNeighbor, ip, ifIndex)
	return pConf, err
}

func (h *BGPHandler) ValidateV4NeighborForUpdate(oldNeigh *bgpd.BGPv4Neighbor, oldNeighConfig config.NeighborConfig,
	newNeigh *bgpd.BGPv4Neighbor, attrSet []bool) (pConf config.NeighborConfig, err error) {
	pConf, _ = h.ConvertV4NeighborFromThrift(newNeigh, oldNeighConfig.NeighborAddress, oldNeighConfig.IfIndex)
	h.logger.Info("ValidateV4NeighborForUpdate: AttrSet", attrSet)
	if attrSet != nil {
		objTyp := reflect.TypeOf(*oldNeigh)
		for i := 0; i < objTyp.NumField(); i++ {
			objName := objTyp.Field(i).Name
			if attrSet[i] {
				h.logger.Debug("ValidateV4NeighborForUpdate : changed ", objName)
				if objName == "UpdateSource" {
					if !h.isValidIP(newNeigh.UpdateSource) {
						err = errors.New(fmt.Sprintf("Update source %s not a valid IP", newNeigh.UpdateSource))
						return pConf, err
					}
				}
			}
		}
	}
	return pConf, err
}

func (h *BGPHandler) SendBGPv4Neighbor(oldNeigh *bgpd.BGPv4Neighbor, newNeigh *bgpd.BGPv4Neighbor, attrSet []bool,
	patchOp []*bgpd.PatchOpInfo, op string) (
	bool, error) {
	h.logger.Info("SendBGPv4Neighbor, op:", op)
	if err := h.checkBGPGlobal(); err != nil {
		h.logger.Err("checkBGPGlobal failed with err:", err)
		return false, err
	}

	oldNeighConf, err := h.ValidateV4Neighbor(oldNeigh)
	if err != nil && op != "create" {
		h.logger.Err("validation of old neighbor failed with err:", err)
		return false, err
	}
	var newNeighConf config.NeighborConfig
	if op != "update" {
		h.logger.Info("not an update op")
		newNeighConf, err = h.ValidateV4Neighbor(newNeigh)
		if err != nil {
			h.logger.Err("validation of newNeigh failed with err:", err)
			return false, err
		}
	} else if patchOp == nil || len(patchOp) == 0 {
		h.logger.Info("update op on v4 neighbor")
		newNeighConf, err = h.ValidateV4NeighborForUpdate(oldNeigh, oldNeighConf, newNeigh, attrSet)
		if err != nil {
			h.logger.Err("validation of v4beighborforupdate failed with err:", err)
			return false, err
		}
	} else {
		//no-op because there are no list objects for neighbor
		h.logger.Info("patch update of v4 neighbor")
	}

	h.server.AddPeerCh <- server.PeerUpdate{oldNeigh, "v4", oldNeighConf, newNeighConf, attrSet, patchOp, op}
	return true, nil
}

func (h *BGPHandler) CreateBGPv4Neighbor(bgpNeighbor *bgpd.BGPv4Neighbor) (bool, error) {
	h.logger.Info("Create BGP neighbor attrs:", bgpNeighbor)
	return h.SendBGPv4Neighbor(nil, bgpNeighbor, make([]bool, 0), nil, "create")
}

func (h *BGPHandler) convertToThriftV4Neighbor(neighborState *config.NeighborState) *bgpd.BGPv4NeighborState {
	bgpNeighborResponse := bgpd.NewBGPv4NeighborState()
	bgpNeighborResponse.NeighborAddress = neighborState.NeighborAddress.String()
	//bgpNeighborResponse.IfIndex = neighborState.IfIndex
	bgpNeighborResponse.IntfRef = "" //strconv.Itoa(int(neighborState.IfIndex))
	intfEntry, ok := h.server.IntfIdNameMap[int32(neighborState.IfIndex)]
	if ok {
		h.logger.Info("Map foud for ifndex : ", neighborState.IfIndex, "Name = ", intfEntry.Name)
		bgpNeighborResponse.IntfRef = intfEntry.Name
	}
	peerAS, err := bgputils.GetAsDot(int(neighborState.PeerAS))
	if err != nil {
		h.logger.Err("peer AS invalid, err:", err)
		return bgpNeighborResponse
	}
	localAS, err := bgputils.GetAsDot(int(neighborState.LocalAS))
	if err != nil {
		h.logger.Err("local AS invalid, err:", err)
		return bgpNeighborResponse
	}
	bgpNeighborResponse.Disabled = neighborState.Disabled
	bgpNeighborResponse.PeerAS = peerAS   //int32(neighborState.PeerAS)
	bgpNeighborResponse.LocalAS = localAS //int32(neighborState.LocalAS)
	bgpNeighborResponse.UpdateSource = neighborState.UpdateSource
	bgpNeighborResponse.NextHopSelf = neighborState.NextHopSelf
	bgpNeighborResponse.AuthPassword = neighborState.AuthPassword
	bgpNeighborResponse.PeerType = int8(neighborState.PeerType)
	bgpNeighborResponse.Description = neighborState.Description
	bgpNeighborResponse.SessionState = int32(neighborState.SessionState)
	bgpNeighborResponse.SessionStateDuration = string(time.Since(neighborState.SessionStateUpdatedTime).String())
	bgpNeighborResponse.RouteReflectorClusterId = int32(neighborState.RouteReflectorClusterId)
	bgpNeighborResponse.RouteReflectorClient = neighborState.RouteReflectorClient
	bgpNeighborResponse.MultiHopEnable = neighborState.MultiHopEnable
	bgpNeighborResponse.MultiHopTTL = int8(neighborState.MultiHopTTL)
	bgpNeighborResponse.ConnectRetryTime = int32(neighborState.ConnectRetryTime)
	bgpNeighborResponse.HoldTime = int32(neighborState.HoldTime)
	bgpNeighborResponse.KeepaliveTime = int32(neighborState.KeepaliveTime)
	bgpNeighborResponse.BfdNeighborState = neighborState.BfdNeighborState
	bgpNeighborResponse.PeerGroup = neighborState.PeerGroup
	bgpNeighborResponse.AddPathsRx = neighborState.AddPathsRx
	bgpNeighborResponse.AddPathsMaxTx = int8(neighborState.AddPathsMaxTx)

	bgpNeighborResponse.MaxPrefixes = int32(neighborState.MaxPrefixes)
	bgpNeighborResponse.MaxPrefixesThresholdPct = int8(neighborState.MaxPrefixesThresholdPct)
	bgpNeighborResponse.MaxPrefixesDisconnect = neighborState.MaxPrefixesDisconnect
	bgpNeighborResponse.MaxPrefixesRestartTimer = int8(neighborState.MaxPrefixesRestartTimer)
	bgpNeighborResponse.TotalPrefixes = int32(neighborState.TotalPrefixes)
	bgpNeighborResponse.AdjRIBInFilter = neighborState.AdjRIBInFilter
	bgpNeighborResponse.AdjRIBOutFilter = neighborState.AdjRIBOutFilter

	received := bgpd.NewBGPCounters()
	received.Notification = int64(neighborState.Messages.Received.Notification)
	received.Update = int64(neighborState.Messages.Received.Update)
	sent := bgpd.NewBGPCounters()
	sent.Notification = int64(neighborState.Messages.Sent.Notification)
	sent.Update = int64(neighborState.Messages.Sent.Update)
	messages := bgpd.NewBGPMessages()
	messages.Received = received
	messages.Sent = sent
	bgpNeighborResponse.Messages = messages

	queues := bgpd.NewBGPQueues()
	queues.Input = int32(neighborState.Queues.Input)
	queues.Output = int32(neighborState.Queues.Output)
	bgpNeighborResponse.Queues = queues

	return bgpNeighborResponse
}

func (h *BGPHandler) GetBGPv4NeighborState(neighborAddr string, intfref string) (*bgpd.BGPv4NeighborState, error) {
	ip, _, _, err := h.getIPAndIfIndexForV4Neighbor(neighborAddr, intfref)
	if err != nil {
		h.logger.Info("GetBGPv4NeighborState: getIPAndIfIndexForV4Neighbor failed for neighbor address", neighborAddr,
			"and ifIndex", intfref)
		return bgpd.NewBGPv4NeighborState(), err
	}

	bgpNeighborState := h.server.GetBGPNeighborState(ip.String())
	if bgpNeighborState == nil {
		return bgpd.NewBGPv4NeighborState(), errors.New(fmt.Sprintf("GetBGPNeighborState: Neighbor %s not configured", ip))
	}
	bgpNeighborResponse := h.convertToThriftV4Neighbor(bgpNeighborState)
	return bgpNeighborResponse, nil
}

func (h *BGPHandler) GetBulkBGPv4NeighborState(index bgpd.Int, count bgpd.Int) (*bgpd.BGPv4NeighborStateGetInfo, error) {
	nextIdx, currCount, bgpNeighbors := h.server.BulkGetBGPv4Neighbors(int(index), int(count))
	bgpNeighborsResponse := make([]*bgpd.BGPv4NeighborState, len(bgpNeighbors))
	for idx, item := range bgpNeighbors {
		bgpNeighborsResponse[idx] = h.convertToThriftV4Neighbor(item)
	}

	bgpNeighborStateBulk := bgpd.NewBGPv4NeighborStateGetInfo()
	bgpNeighborStateBulk.EndIdx = bgpd.Int(nextIdx)
	bgpNeighborStateBulk.Count = bgpd.Int(currCount)
	bgpNeighborStateBulk.More = (nextIdx != 0)
	bgpNeighborStateBulk.BGPv4NeighborStateList = bgpNeighborsResponse

	return bgpNeighborStateBulk, nil
}

func (h *BGPHandler) UpdateBGPv4Neighbor(origN *bgpd.BGPv4Neighbor, updatedN *bgpd.BGPv4Neighbor, attrSet []bool,
	op []*bgpd.PatchOpInfo) (bool, error) {
	h.logger.Info("Update peer attrs:", updatedN)
	return h.SendBGPv4Neighbor(origN, updatedN, attrSet, op, "update")
}

func (h *BGPHandler) DeleteBGPv4Neighbor(bgpNeighbor *bgpd.BGPv4Neighbor) (bool, error) {
	h.logger.Info("Delete BGPv4 neighbor:", bgpNeighbor.NeighborAddress)
	if err := h.checkBGPGlobal(); err != nil {
		return false, err
	}

	neighConf, err := h.ValidateV4Neighbor(bgpNeighbor)
	if err != nil {
		h.logger.Info("DeleteBGPv4Neighbor: ValidateV6Neighbor failed for neighbor address",
			bgpNeighbor.NeighborAddress, "and ifIndex", bgpNeighbor.IntfRef)
		return false, err
	}

	h.server.RemPeerCh <- neighConf
	return true, nil
}

func (h *BGPHandler) getIPAndIfIndexForV6Neighbor(neighborIP string, neighborIntfRef string) (ip net.IP, ifIndex int32,
	ifName string, err error) {
	if strings.TrimSpace(neighborIP) != "" {
		ifIndex = -1
		ip = net.ParseIP(strings.TrimSpace(neighborIP))
		if ip == nil {
			err = errors.New(fmt.Sprintf("v6Neighbor address %s not valid", neighborIP))
		}
	} else if neighborIntfRef != "" {
		//neighbor address is a intfRef
		ifIndex, ifName, err = h.server.ConvertIntfStrToIfIndex(neighborIntfRef)
		h.logger.Info("getIPAndIfIndexForV6Neighbor - ifIndex:", ifIndex, "ifName:", ifName)

		if err != nil {
			h.logger.Err("Invalid intfref:", neighborIntfRef)
			return ip, ifIndex, ifName, err
		}

		var ipInfo *bgputils.IPInfo
		ipInfo, err = h.server.GetIfaceIP(ifIndex)
		h.logger.Info("ipInfo:", ipInfo, " err:", err)
		if err == nil {
			h.logger.Info("getIPAndIfIndexForV6Neighbor - ipInfo.LinkLocalIpAddr:", ipInfo.LinklocalIpAddr,
				"after GetIfaceIP of neighborIfIndex:", ifIndex)
			if ipInfo.LinklocalIpAddr != "" {
				ip = net.ParseIP(ipInfo.LinklocalIpAddr)
			} else {
				return ip, ifIndex, ifName, errors.New(fmt.Sprint("Link-local IP is not set on interface ", neighborIntfRef))
			}
		}
	}
	return ip, ifIndex, ifName, err
}

func (h *BGPHandler) ConvertV6NeighborFromThrift(bgpNeighbor *bgpd.BGPv6Neighbor, ip net.IP, ifIndex int32,
	ifName string) (pConf config.NeighborConfig, err error) {
	peerAS, err := bgputils.GetAsNum(bgpNeighbor.PeerAS)
	if err != nil {
		h.logger.Err("peer AS invalid, err:", err)
		return pConf, err
	}
	localAS, err := bgputils.GetAsNum(bgpNeighbor.LocalAS)
	if err != nil {
		h.logger.Err("local AS invalid, err:", err)
		return pConf, err
	}
	pConf = config.NeighborConfig{
		BaseConfig: config.BaseConfig{
			PeerAS:                  uint32(peerAS),
			LocalAS:                 uint32(localAS),
			PeerAddressType:         config.PeerAddressV6,
			UpdateSource:            bgpNeighbor.UpdateSource,
			NextHopSelf:             bgpNeighbor.NextHopSelf,
			Description:             bgpNeighbor.Description,
			RouteReflectorClusterId: uint32(bgpNeighbor.RouteReflectorClusterId),
			RouteReflectorClient:    bgpNeighbor.RouteReflectorClient,
			MultiHopEnable:          bgpNeighbor.MultiHopEnable,
			MultiHopTTL:             uint8(bgpNeighbor.MultiHopTTL),
			ConnectRetryTime:        uint32(bgpNeighbor.ConnectRetryTime),
			HoldTime:                uint32(bgpNeighbor.HoldTime),
			KeepaliveTime:           uint32(bgpNeighbor.KeepaliveTime),
			BfdEnable:               bgpNeighbor.BfdEnable,
			BfdSessionParam:         bgpNeighbor.BfdSessionParam,
			AddPathsRx:              bgpNeighbor.AddPathsRx,
			AddPathsMaxTx:           uint8(bgpNeighbor.AddPathsMaxTx),
			MaxPrefixes:             uint32(bgpNeighbor.MaxPrefixes),
			MaxPrefixesThresholdPct: uint8(bgpNeighbor.MaxPrefixesThresholdPct),
			MaxPrefixesDisconnect:   bgpNeighbor.MaxPrefixesDisconnect,
			MaxPrefixesRestartTimer: uint8(bgpNeighbor.MaxPrefixesRestartTimer),
			AdjRIBInFilter:          bgpNeighbor.AdjRIBInFilter,
			AdjRIBOutFilter:         bgpNeighbor.AdjRIBOutFilter,
		},
		NeighborAddress: ip,
		IfIndex:         ifIndex,
		IfName:          ifName,
		PeerGroup:       bgpNeighbor.PeerGroup,
		Disabled:        bgpNeighbor.Disabled,
	}
	return pConf, err
}

func (h *BGPHandler) ValidateV6Neighbor(bgpNeighbor *bgpd.BGPv6Neighbor) (pConf config.NeighborConfig, err error) {
	if bgpNeighbor == nil {
		return pConf, err
	}

	var ip net.IP
	var ifIndex int32
	var ifName string
	ifIndex = -1
	ip, ifIndex, ifName, err = h.getIPAndIfIndexForV6Neighbor(bgpNeighbor.NeighborAddress, bgpNeighbor.IntfRef)
	if err != nil {
		h.logger.Info("ValidateV6Neighbor: getIPAndIfIndexForNeighbor failed for neighbor address",
			bgpNeighbor.NeighborAddress, "and ifIndex", bgpNeighbor.IntfRef)
		return pConf, err
	}
	h.logger.Info("ValidateV6Neighbor: getIPAndIfIndexForNeighbor returned ip", ip, "ifIndex", ifIndex, "ifName", ifName)

	if !h.isValidIP(bgpNeighbor.UpdateSource) {
		err = errors.New(fmt.Sprintf("Update source %s not a valid IP", bgpNeighbor.UpdateSource))
		return pConf, err
	}

	pConf, _ = h.ConvertV6NeighborFromThrift(bgpNeighbor, ip, ifIndex, ifName)
	return pConf, err
}

func (h *BGPHandler) ValidateV6NeighborForUpdate(oldNeigh *bgpd.BGPv6Neighbor, oldNeighConfig config.NeighborConfig,
	newNeigh *bgpd.BGPv6Neighbor, attrSet []bool) (pConf config.NeighborConfig, err error) {
	pConf, _ = h.ConvertV6NeighborFromThrift(newNeigh, oldNeighConfig.NeighborAddress, oldNeighConfig.IfIndex,
		oldNeighConfig.IfName)
	if attrSet != nil {
		objTyp := reflect.TypeOf(*oldNeigh)
		for i := 0; i < objTyp.NumField(); i++ {
			objName := objTyp.Field(i).Name
			if attrSet[i] {
				h.logger.Debug("ValidateV6NeighborForUpdate : changed ", objName)
				if objName == "UpdateSource" {
					if !h.isValidIP(newNeigh.UpdateSource) {
						err = errors.New(fmt.Sprintf("Update source %s not a valid IP", newNeigh.UpdateSource))
						return pConf, err
					}
				}
			}
		}
	}
	return pConf, err
}

func (h *BGPHandler) SendBGPv6Neighbor(oldNeigh *bgpd.BGPv6Neighbor, newNeigh *bgpd.BGPv6Neighbor, attrSet []bool, patchOp []*bgpd.PatchOpInfo, op string) (
	bool, error) {
	if err := h.checkBGPGlobal(); err != nil {
		return false, err
	}

	oldNeighConf, err := h.ValidateV6Neighbor(oldNeigh)
	if err != nil {
		return false, err
	}
	var newNeighConf config.NeighborConfig
	if op != "update" {
		newNeighConf, err = h.ValidateV6Neighbor(newNeigh)
		if err != nil {
			return false, err
		}
	} else if patchOp == nil || len(patchOp) == 0 {
		//no-op because there are no list objects for neighbor
	} else {
		newNeighConf, err = h.ValidateV6NeighborForUpdate(oldNeigh, newNeighConf, newNeigh, attrSet)
		if err != nil {
			return false, err
		}
	}
	h.server.AddPeerCh <- server.PeerUpdate{oldNeigh, "v6", oldNeighConf, newNeighConf, attrSet, patchOp, op}
	return true, nil
}

func (h *BGPHandler) CreateBGPv6Neighbor(bgpNeighbor *bgpd.BGPv6Neighbor) (bool, error) {
	h.logger.Info("Create BGP neighbor attrs:", bgpNeighbor)
	return h.SendBGPv6Neighbor(nil, bgpNeighbor, make([]bool, 0), nil, "create")
}

func (h *BGPHandler) convertToThriftV6Neighbor(neighborState *config.NeighborState) *bgpd.BGPv6NeighborState {
	bgpNeighborResponse := bgpd.NewBGPv6NeighborState()
	bgpNeighborResponse.NeighborAddress = neighborState.NeighborAddress.String()
	//bgpNeighborResponse.IfIndex = neighborState.IfIndex
	bgpNeighborResponse.IntfRef = "" //strconv.Itoa(int(neighborState.IfIndex))
	intfEntry, ok := h.server.IntfIdNameMap[int32(neighborState.IfIndex)]
	if ok {
		h.logger.Info("Map foud for ifndex : ", neighborState.IfIndex, "Name = ", intfEntry.Name)
		bgpNeighborResponse.IntfRef = intfEntry.Name
	}
	peerAS, err := bgputils.GetAsDot(int(neighborState.PeerAS))
	if err != nil {
		h.logger.Err("peer AS invalid, err:", err)
		return bgpNeighborResponse
	}
	localAS, err := bgputils.GetAsDot(int(neighborState.LocalAS))
	if err != nil {
		h.logger.Err("local AS invalid, err:", err)
		return bgpNeighborResponse
	}
	bgpNeighborResponse.Disabled = neighborState.Disabled
	bgpNeighborResponse.PeerAS = peerAS   // int32(neighborState.PeerAS)
	bgpNeighborResponse.LocalAS = localAS //int32(neighborState.LocalAS)
	bgpNeighborResponse.UpdateSource = neighborState.UpdateSource
	bgpNeighborResponse.NextHopSelf = neighborState.NextHopSelf
	bgpNeighborResponse.PeerType = int8(neighborState.PeerType)
	bgpNeighborResponse.Description = neighborState.Description
	bgpNeighborResponse.SessionState = int32(neighborState.SessionState)
	bgpNeighborResponse.SessionStateDuration = string(time.Since(neighborState.SessionStateUpdatedTime).String())
	bgpNeighborResponse.RouteReflectorClusterId = int32(neighborState.RouteReflectorClusterId)
	bgpNeighborResponse.RouteReflectorClient = neighborState.RouteReflectorClient
	bgpNeighborResponse.MultiHopEnable = neighborState.MultiHopEnable
	bgpNeighborResponse.MultiHopTTL = int8(neighborState.MultiHopTTL)
	bgpNeighborResponse.ConnectRetryTime = int32(neighborState.ConnectRetryTime)
	bgpNeighborResponse.HoldTime = int32(neighborState.HoldTime)
	bgpNeighborResponse.KeepaliveTime = int32(neighborState.KeepaliveTime)
	bgpNeighborResponse.BfdNeighborState = neighborState.BfdNeighborState
	bgpNeighborResponse.PeerGroup = neighborState.PeerGroup
	bgpNeighborResponse.AddPathsRx = neighborState.AddPathsRx
	bgpNeighborResponse.AddPathsMaxTx = int8(neighborState.AddPathsMaxTx)

	bgpNeighborResponse.MaxPrefixes = int32(neighborState.MaxPrefixes)
	bgpNeighborResponse.MaxPrefixesThresholdPct = int8(neighborState.MaxPrefixesThresholdPct)
	bgpNeighborResponse.MaxPrefixesDisconnect = neighborState.MaxPrefixesDisconnect
	bgpNeighborResponse.MaxPrefixesRestartTimer = int8(neighborState.MaxPrefixesRestartTimer)
	bgpNeighborResponse.TotalPrefixes = int32(neighborState.TotalPrefixes)
	bgpNeighborResponse.AdjRIBInFilter = neighborState.AdjRIBInFilter
	bgpNeighborResponse.AdjRIBOutFilter = neighborState.AdjRIBOutFilter

	received := bgpd.NewBGPCounters()
	received.Notification = int64(neighborState.Messages.Received.Notification)
	received.Update = int64(neighborState.Messages.Received.Update)
	sent := bgpd.NewBGPCounters()
	sent.Notification = int64(neighborState.Messages.Sent.Notification)
	sent.Update = int64(neighborState.Messages.Sent.Update)
	messages := bgpd.NewBGPMessages()
	messages.Received = received
	messages.Sent = sent
	bgpNeighborResponse.Messages = messages

	queues := bgpd.NewBGPQueues()
	queues.Input = int32(neighborState.Queues.Input)
	queues.Output = int32(neighborState.Queues.Output)
	bgpNeighborResponse.Queues = queues

	return bgpNeighborResponse
}

func (h *BGPHandler) GetBGPv6NeighborState(neighborAddr string, intfref string) (*bgpd.BGPv6NeighborState, error) {
	ip, _, _, err := h.getIPAndIfIndexForV6Neighbor(neighborAddr, intfref)
	if err != nil {
		h.logger.Info("GetBGPv4NeighborState: getIPAndIfIndexForV4Neighbor failed for neighbor address", neighborAddr,
			"and ifIndex", intfref)
		return bgpd.NewBGPv6NeighborState(), err
	}

	bgpNeighborState := h.server.GetBGPNeighborState(ip.String())
	if bgpNeighborState == nil {
		return bgpd.NewBGPv6NeighborState(), errors.New(fmt.Sprintf("GetBGPNeighborState: Neighbor %s not configured", ip))
	}
	bgpNeighborResponse := h.convertToThriftV6Neighbor(bgpNeighborState)
	return bgpNeighborResponse, nil
}

func (h *BGPHandler) GetBulkBGPv6NeighborState(index bgpd.Int, count bgpd.Int) (*bgpd.BGPv6NeighborStateGetInfo,
	error) {
	nextIdx, currCount, bgpNeighbors := h.server.BulkGetBGPv6Neighbors(int(index), int(count))
	bgpNeighborsResponse := make([]*bgpd.BGPv6NeighborState, len(bgpNeighbors))
	for idx, item := range bgpNeighbors {
		bgpNeighborsResponse[idx] = h.convertToThriftV6Neighbor(item)
	}

	bgpNeighborStateBulk := bgpd.NewBGPv6NeighborStateGetInfo()
	bgpNeighborStateBulk.EndIdx = bgpd.Int(nextIdx)
	bgpNeighborStateBulk.Count = bgpd.Int(currCount)
	bgpNeighborStateBulk.More = (nextIdx != 0)
	bgpNeighborStateBulk.BGPv6NeighborStateList = bgpNeighborsResponse

	return bgpNeighborStateBulk, nil
}

func (h *BGPHandler) UpdateBGPv6Neighbor(origN *bgpd.BGPv6Neighbor, updatedN *bgpd.BGPv6Neighbor, attrSet []bool,
	op []*bgpd.PatchOpInfo) (bool, error) {
	h.logger.Info("Update peer attrs:", updatedN)
	return h.SendBGPv6Neighbor(origN, updatedN, attrSet, op, "update")
}

func (h *BGPHandler) DeleteBGPv6Neighbor(bgpNeighbor *bgpd.BGPv6Neighbor) (bool, error) {
	h.logger.Info("Delete BGPv6 neighbor:", bgpNeighbor.NeighborAddress)
	if err := h.checkBGPGlobal(); err != nil {
		return false, err
	}

	neighConf, err := h.ValidateV6Neighbor(bgpNeighbor)
	if err != nil {
		return false, err
	}

	h.server.RemPeerCh <- neighConf
	return true, nil
}

func (h *BGPHandler) PeerCommand(in *PeerConfigCommands, out *bool) error {
	h.PeerCommandCh <- *in
	h.logger.Info("Good peer command:", in)
	*out = true
	return nil
}

func (h *BGPHandler) ValidateBGPv4PeerGroup(peerGroup *bgpd.BGPv4PeerGroup) (group config.PeerGroupConfig,
	err error) {
	if peerGroup == nil {
		return group, err
	}

	if !h.isValidIP(peerGroup.UpdateSource) {
		err = errors.New(fmt.Sprintf("Update source %s not a valid IP", peerGroup.UpdateSource))
		return group, err
	}
	peerAS, err := bgputils.GetAsNum(peerGroup.PeerAS)
	if err != nil {
		h.logger.Err("peer AS invalid, err:", err)
		return group, err
	}
	localAS, err := bgputils.GetAsNum(peerGroup.LocalAS)
	if err != nil {
		h.logger.Err("local AS invalid, err:", err)
		return group, err
	}

	group = config.PeerGroupConfig{
		BaseConfig: config.BaseConfig{
			PeerAS:                  uint32(peerAS),
			LocalAS:                 uint32(localAS),
			PeerAddressType:         config.PeerAddressV4,
			UpdateSource:            peerGroup.UpdateSource,
			NextHopSelf:             peerGroup.NextHopSelf,
			AuthPassword:            peerGroup.AuthPassword,
			Description:             peerGroup.Description,
			RouteReflectorClusterId: uint32(peerGroup.RouteReflectorClusterId),
			RouteReflectorClient:    peerGroup.RouteReflectorClient,
			MultiHopEnable:          peerGroup.MultiHopEnable,
			MultiHopTTL:             uint8(peerGroup.MultiHopTTL),
			ConnectRetryTime:        uint32(peerGroup.ConnectRetryTime),
			HoldTime:                uint32(peerGroup.HoldTime),
			KeepaliveTime:           uint32(peerGroup.KeepaliveTime),
			AddPathsRx:              peerGroup.AddPathsRx,
			AddPathsMaxTx:           uint8(peerGroup.AddPathsMaxTx),
			MaxPrefixes:             uint32(peerGroup.MaxPrefixes),
			MaxPrefixesThresholdPct: uint8(peerGroup.MaxPrefixesThresholdPct),
			MaxPrefixesDisconnect:   peerGroup.MaxPrefixesDisconnect,
			MaxPrefixesRestartTimer: uint8(peerGroup.MaxPrefixesRestartTimer),
			AdjRIBInFilter:          peerGroup.AdjRIBInFilter,
			AdjRIBOutFilter:         peerGroup.AdjRIBOutFilter,
		},
		Name: peerGroup.Name,
	}

	return group, err
}

func (h *BGPHandler) SendBGPv4PeerGroup(oldGroup *bgpd.BGPv4PeerGroup, newGroup *bgpd.BGPv4PeerGroup, attrSet []bool) (
	bool, error) {
	if err := h.checkBGPGlobal(); err != nil {
		return false, err
	}

	oldGroupConf, err := h.ValidateBGPv4PeerGroup(oldGroup)
	if err != nil {
		return false, err
	}

	newGroupConf, err := h.ValidateBGPv4PeerGroup(newGroup)
	if err != nil {
		return false, err
	}

	h.server.AddPeerGroupCh <- server.PeerGroupUpdate{oldGroupConf, newGroupConf, attrSet}
	return true, nil
}

func (h *BGPHandler) CreateBGPv4PeerGroup(peerGroup *bgpd.BGPv4PeerGroup) (bool, error) {
	h.logger.Infof("Create BGP v4 peer group:%+v", peerGroup)
	return h.SendBGPv4PeerGroup(nil, peerGroup, make([]bool, 0))
}

func (h *BGPHandler) UpdateBGPv4PeerGroup(origG *bgpd.BGPv4PeerGroup, updatedG *bgpd.BGPv4PeerGroup, attrSet []bool,
	op []*bgpd.PatchOpInfo) (bool, error) {
	h.logger.Infof("Update BGP v4 peer group:%+v", updatedG)
	return h.SendBGPv4PeerGroup(origG, updatedG, attrSet)
}

func (h *BGPHandler) DeleteBGPv4PeerGroup(peerGroup *bgpd.BGPv4PeerGroup) (bool, error) {
	h.logger.Info("Delete BGP v4 peer group:%+v", peerGroup.Name)
	if err := h.checkBGPGlobal(); err != nil {
		return false, err
	}

	gConf, err := h.ValidateBGPv4PeerGroup(peerGroup)
	if err != nil {
		return false, err
	}
	h.server.RemPeerGroupCh <- gConf
	return true, nil
}

func (h *BGPHandler) ValidateBGPv6PeerGroup(peerGroup *bgpd.BGPv6PeerGroup) (group config.PeerGroupConfig,
	err error) {
	if peerGroup == nil {
		return group, err
	}

	if !h.isValidIP(peerGroup.UpdateSource) {
		err = errors.New(fmt.Sprintf("Update source %s not a valid IP", peerGroup.UpdateSource))
		return group, err
	}
	peerAS, err := bgputils.GetAsNum(peerGroup.PeerAS)
	if err != nil {
		h.logger.Err("peer AS invalid, err:", err)
		return group, err
	}
	localAS, err := bgputils.GetAsNum(peerGroup.LocalAS)
	if err != nil {
		h.logger.Err("local AS invalid, err:", err)
		return group, err
	}

	group = config.PeerGroupConfig{
		BaseConfig: config.BaseConfig{
			PeerAS:                  uint32(peerAS),
			LocalAS:                 uint32(localAS),
			PeerAddressType:         config.PeerAddressV6,
			UpdateSource:            peerGroup.UpdateSource,
			NextHopSelf:             peerGroup.NextHopSelf,
			Description:             peerGroup.Description,
			RouteReflectorClusterId: uint32(peerGroup.RouteReflectorClusterId),
			RouteReflectorClient:    peerGroup.RouteReflectorClient,
			MultiHopEnable:          peerGroup.MultiHopEnable,
			MultiHopTTL:             uint8(peerGroup.MultiHopTTL),
			ConnectRetryTime:        uint32(peerGroup.ConnectRetryTime),
			HoldTime:                uint32(peerGroup.HoldTime),
			KeepaliveTime:           uint32(peerGroup.KeepaliveTime),
			AddPathsRx:              peerGroup.AddPathsRx,
			AddPathsMaxTx:           uint8(peerGroup.AddPathsMaxTx),
			MaxPrefixes:             uint32(peerGroup.MaxPrefixes),
			MaxPrefixesThresholdPct: uint8(peerGroup.MaxPrefixesThresholdPct),
			MaxPrefixesDisconnect:   peerGroup.MaxPrefixesDisconnect,
			MaxPrefixesRestartTimer: uint8(peerGroup.MaxPrefixesRestartTimer),
			AdjRIBInFilter:          peerGroup.AdjRIBInFilter,
			AdjRIBOutFilter:         peerGroup.AdjRIBOutFilter,
		},
		Name: peerGroup.Name,
	}

	return group, err
}

func (h *BGPHandler) SendBGPv6PeerGroup(oldGroup *bgpd.BGPv6PeerGroup, newGroup *bgpd.BGPv6PeerGroup, attrSet []bool) (
	bool, error) {
	if err := h.checkBGPGlobal(); err != nil {
		return false, err
	}

	oldGroupConf, err := h.ValidateBGPv6PeerGroup(oldGroup)
	if err != nil {
		return false, err
	}

	newGroupConf, err := h.ValidateBGPv6PeerGroup(newGroup)
	if err != nil {
		return false, err
	}

	h.server.AddPeerGroupCh <- server.PeerGroupUpdate{oldGroupConf, newGroupConf, attrSet}
	return true, nil
}

func (h *BGPHandler) CreateBGPv6PeerGroup(peerGroup *bgpd.BGPv6PeerGroup) (bool, error) {
	h.logger.Info("Create BGP v6 peer group:", peerGroup)
	return h.SendBGPv6PeerGroup(nil, peerGroup, make([]bool, 0))
}

func (h *BGPHandler) UpdateBGPv6PeerGroup(origG *bgpd.BGPv6PeerGroup, updatedG *bgpd.BGPv6PeerGroup, attrSet []bool,
	op []*bgpd.PatchOpInfo) (bool, error) {
	h.logger.Info("Update BGP v6 peer group:", updatedG)
	return h.SendBGPv6PeerGroup(origG, updatedG, attrSet)
}

func (h *BGPHandler) DeleteBGPv6PeerGroup(peerGroup *bgpd.BGPv6PeerGroup) (bool, error) {
	h.logger.Info("Delete BGP v6 peer group:", peerGroup.Name)
	if err := h.checkBGPGlobal(); err != nil {
		return false, err
	}

	gConf, err := h.ValidateBGPv6PeerGroup(peerGroup)
	if err != nil {
		return false, err
	}
	h.server.RemPeerGroupCh <- gConf
	return true, nil
}

func (h *BGPHandler) GetBGPv4RouteState(network string, cidrLen int16) (*bgpd.BGPv4RouteState, error) {
	bgpRoute := h.server.LocRib.GetBGPv4Route(network)
	var err error = nil
	if bgpRoute == nil {
		err = errors.New(fmt.Sprintf("Route not found for destination %s", network))
	}
	return bgpRoute, err
}

func (h *BGPHandler) GetBulkBGPv4RouteState(index bgpd.Int, count bgpd.Int) (*bgpd.BGPv4RouteStateGetInfo, error) {
	nextIdx, currCount, bgpRoutes := h.server.LocRib.BulkGetBGPv4Routes(int(index), int(count))

	bgpRoutesBulk := bgpd.NewBGPv4RouteStateGetInfo()
	bgpRoutesBulk.EndIdx = bgpd.Int(nextIdx)
	bgpRoutesBulk.Count = bgpd.Int(currCount)
	bgpRoutesBulk.More = (nextIdx != 0)
	bgpRoutesBulk.BGPv4RouteStateList = bgpRoutes

	return bgpRoutesBulk, nil
}

func (h *BGPHandler) GetBGPv6RouteState(network string, cidrLen int16) (*bgpd.BGPv6RouteState, error) {
	bgpRoute := h.server.LocRib.GetBGPv6Route(network)
	var err error = nil
	if bgpRoute == nil {
		err = errors.New(fmt.Sprintf("Route not found for destination %s", network))
	}
	return bgpRoute, err
}

func (h *BGPHandler) GetBulkBGPv6RouteState(index bgpd.Int, count bgpd.Int) (*bgpd.BGPv6RouteStateGetInfo, error) {
	nextIdx, currCount, bgpRoutes := h.server.LocRib.BulkGetBGPv6Routes(int(index), int(count))

	bgpRoutesBulk := bgpd.NewBGPv6RouteStateGetInfo()
	bgpRoutesBulk.EndIdx = bgpd.Int(nextIdx)
	bgpRoutesBulk.Count = bgpd.Int(currCount)
	bgpRoutesBulk.More = (nextIdx != 0)
	bgpRoutesBulk.BGPv6RouteStateList = bgpRoutes

	return bgpRoutesBulk, nil
}

func (h *BGPHandler) validateBGPAggregate(bgpAgg *bgpd.BGPv4Aggregate) (aggConf config.BGPAggregate, err error) {
	if bgpAgg == nil {
		return aggConf, err
	}
	var ip net.IP

	ip, _, err = net.ParseCIDR(bgpAgg.IpPrefix)
	if err != nil {
		err = errors.New(fmt.Sprintf("BGPAggregate: IP %s is not valid", bgpAgg.IpPrefix))
		h.logger.Info("SendBGPAggregate: IP", bgpAgg.IpPrefix, "is not valid")
		return aggConf, err
	}

	if ip.To4() == nil {
		err = errors.New(fmt.Sprintf("BGPAggregate: IP %s is not a v4 address", bgpAgg.IpPrefix))
		h.logger.Info("SendBGPAggregate: IP", bgpAgg.IpPrefix, "is not a v4 address")
		return aggConf, err
	}

	aggConf = config.BGPAggregate{
		IPPrefix:        bgpAgg.IpPrefix,
		GenerateASSet:   bgpAgg.GenerateASSet,
		SendSummaryOnly: bgpAgg.SendSummaryOnly,
		AddressFamily:   packet.GetProtocolFamily(packet.AfiIP, packet.SafiUnicast),
	}
	return aggConf, nil
}

func (h *BGPHandler) SendBGPAggregate(oldConfig *bgpd.BGPv4Aggregate, newConfig *bgpd.BGPv4Aggregate, attrSet []bool) (
	bool, error) {
	if err := h.checkBGPGlobal(); err != nil {
		return false, err
	}

	oldAgg, err := h.validateBGPAggregate(oldConfig)
	if err != nil {
		return false, err
	}

	newAgg, err := h.validateBGPAggregate(newConfig)
	if err != nil {
		return false, err
	}

	h.server.AddAggCh <- server.AggUpdate{oldAgg, newAgg, attrSet}
	return true, err
}

func (h *BGPHandler) CreateBGPv4Aggregate(bgpAgg *bgpd.BGPv4Aggregate) (bool, error) {
	h.logger.Info("Create BGP v4 aggregate:", bgpAgg)
	return h.SendBGPAggregate(nil, bgpAgg, make([]bool, 0))
}

func (h *BGPHandler) UpdateBGPv4Aggregate(origA *bgpd.BGPv4Aggregate, updatedA *bgpd.BGPv4Aggregate, attrSet []bool,
	op []*bgpd.PatchOpInfo) (bool, error) {
	h.logger.Info("Update BGP v4 aggregate:", updatedA, "old:", origA)
	return h.SendBGPAggregate(origA, updatedA, attrSet)
}

func (h *BGPHandler) DeleteBGPv4Aggregate(bgpAgg *bgpd.BGPv4Aggregate) (bool, error) {
	h.logger.Info("Delete BGP v4 aggregate:", bgpAgg)
	if err := h.checkBGPGlobal(); err != nil {
		return false, err
	}

	agg, _ := h.validateBGPAggregate(bgpAgg)
	h.server.RemAggCh <- agg
	return true, nil
}

func (h *BGPHandler) validateBGPv6Aggregate(bgpAgg *bgpd.BGPv6Aggregate) (aggConf config.BGPAggregate, err error) {
	if bgpAgg == nil {
		return aggConf, err
	}

	var ip net.IP
	ip, _, err = net.ParseCIDR(bgpAgg.IpPrefix)
	if err != nil {
		err = errors.New(fmt.Sprintf("BGPAggregate: IP %s is not valid", bgpAgg.IpPrefix))
		h.logger.Info("SendBGPAggregate: IP", bgpAgg.IpPrefix, "is not valid")
		return aggConf, err
	}

	if ip.To4() != nil {
		err = errors.New(fmt.Sprintf("BGPAggregate: IP %s is not a v6 address", bgpAgg.IpPrefix))
		h.logger.Info("SendBGPAggregate: IP", bgpAgg.IpPrefix, "is not a v6 address")
		return aggConf, err
	}

	aggConf = config.BGPAggregate{
		IPPrefix:        bgpAgg.IpPrefix,
		GenerateASSet:   bgpAgg.GenerateASSet,
		SendSummaryOnly: bgpAgg.SendSummaryOnly,
		AddressFamily:   packet.GetProtocolFamily(packet.AfiIP6, packet.SafiUnicast),
	}
	return aggConf, nil
}

func (h *BGPHandler) SendBGPv6Aggregate(oldConfig *bgpd.BGPv6Aggregate, newConfig *bgpd.BGPv6Aggregate,
	attrSet []bool) (bool, error) {
	if err := h.checkBGPGlobal(); err != nil {
		return false, err
	}

	oldAgg, err := h.validateBGPv6Aggregate(oldConfig)
	if err != nil {
		return false, err
	}

	newAgg, err := h.validateBGPv6Aggregate(newConfig)
	if err != nil {
		return false, err
	}

	h.server.AddAggCh <- server.AggUpdate{oldAgg, newAgg, attrSet}
	return true, err
}

func (h *BGPHandler) CreateBGPv6Aggregate(bgpAgg *bgpd.BGPv6Aggregate) (bool, error) {
	h.logger.Info("Create BGP IPv6 aggregate:", bgpAgg)
	return h.SendBGPv6Aggregate(nil, bgpAgg, make([]bool, 0))
}

func (h *BGPHandler) UpdateBGPv6Aggregate(origA *bgpd.BGPv6Aggregate, updatedA *bgpd.BGPv6Aggregate, attrSet []bool,
	op []*bgpd.PatchOpInfo) (bool, error) {
	h.logger.Info("Update BGP IPv6 aggregate:", updatedA, "old:", origA)
	return h.SendBGPv6Aggregate(origA, updatedA, attrSet)
}

func (h *BGPHandler) DeleteBGPv6Aggregate(bgpAgg *bgpd.BGPv6Aggregate) (bool, error) {
	h.logger.Info("Delete BGP IPv6 aggregate:", bgpAgg)
	if err := h.checkBGPGlobal(); err != nil {
		return false, err
	}

	agg, _ := h.validateBGPv6Aggregate(bgpAgg)
	h.server.RemAggCh <- agg
	return true, nil
}

func (h *BGPHandler) validateBGPNetworkStatement(bgpNS *bgpd.BGPNetworkStatement) (nsConf config.BGPNetworkStatement,
	err error) {
	if bgpNS == nil {
		return nsConf, err
	}
	var ip net.IP

	ip, _, err = net.ParseCIDR(bgpNS.IpPrefix)
	if err != nil {
		err = errors.New(fmt.Sprintf("BGPNetworkStatement: IP %s is not valid", bgpNS.IpPrefix))
		h.logger.Info("BGPNetworkStatement: IP", bgpNS.IpPrefix, "is not valid")
		return nsConf, err
	}

	var afi packet.AFI
	if ip.To4() != nil {
		afi = packet.AfiIP
	} else if ip.To16() != nil {
		afi = packet.AfiIP6
	} else {
		err = errors.New(fmt.Sprintf("BGPNetworkStatement: IP %s is not a valid v4 or v6 address", bgpNS.IpPrefix))
		h.logger.Info("BGPNetworkStatement: IP", bgpNS.IpPrefix, "is not a valid v4 or v6 address")
		return nsConf, err
	}

	nsConf = config.BGPNetworkStatement{
		IPPrefix:      bgpNS.IpPrefix,
		Policy:        bgpNS.Policy,
		AddressFamily: packet.GetProtocolFamily(afi, packet.SafiUnicast),
	}
	return nsConf, nil
}

func (h *BGPHandler) SendBGPNetworkStatement(oldConfig *bgpd.BGPNetworkStatement, newConfig *bgpd.BGPNetworkStatement,
	attrSet []bool) (
	bool, error) {
	if err := h.checkBGPGlobal(); err != nil {
		return false, err
	}

	oldAgg, err := h.validateBGPNetworkStatement(oldConfig)
	if err != nil {
		return false, err
	}

	newAgg, err := h.validateBGPNetworkStatement(newConfig)
	if err != nil {
		return false, err
	}

	h.server.AddNSCh <- server.NSUpdate{oldAgg, newAgg, attrSet}
	return true, err
}

func (h *BGPHandler) CreateBGPNetworkStatement(bgpNS *bgpd.BGPNetworkStatement) (bool, error) {
	h.logger.Info("Create BGP network statement:", bgpNS)
	return h.SendBGPNetworkStatement(nil, bgpNS, make([]bool, 0))
}

func (h *BGPHandler) UpdateBGPNetworkStatement(bgpNS *bgpd.BGPNetworkStatement, updatedNS *bgpd.BGPNetworkStatement,
	attrSet []bool, op []*bgpd.PatchOpInfo) (bool, error) {
	h.logger.Info("Update BGP network statement:", updatedNS, "old:", bgpNS)
	return h.SendBGPNetworkStatement(bgpNS, updatedNS, attrSet)
}

func (h *BGPHandler) DeleteBGPNetworkStatement(bgpNS *bgpd.BGPNetworkStatement) (bool, error) {
	h.logger.Info("Delete BGP network statement:", bgpNS)
	if err := h.checkBGPGlobal(); err != nil {
		return false, err
	}

	nsConf, _ := h.validateBGPNetworkStatement(bgpNS)
	h.server.RemNSCh <- nsConf
	return true, nil
}

func (h *BGPHandler) ConvertBGPNetworkStateToThriftObject(ns *config.BGPNetworkStatement) *bgpd.BGPNetworkStatementState {
	return &bgpd.BGPNetworkStatementState{
		IpPrefix:       ns.IPPrefix,
		Policy:         ns.Policy,
		PolicyStmtList: ns.PolicyStmtList,
	}
}

func (h *BGPHandler) GetBGPNetworkStatementState(ipPrefix string) (*bgpd.BGPNetworkStatementState, error) {
	var err error = nil
	var bgpNS *bgpd.BGPNetworkStatementState

	ns := h.server.GetBGPNetworkStatementState(ipPrefix)
	if ns == nil {
		err = errors.New(fmt.Sprintf("Network statement not found for prefix %s", ipPrefix))
	} else {
		bgpNS = h.ConvertBGPNetworkStateToThriftObject(ns)
	}
	return bgpNS, err
}

func (h *BGPHandler) GetBulkBGPNetworkStatementState(index bgpd.Int, count bgpd.Int) (
	*bgpd.BGPNetworkStatementStateGetInfo, error) {
	nextIdx, currCount, nsRoutes := h.server.BulkGetBGPNetworkStatements(int(index), int(count))
	bgpdNSRoutes := make([]*bgpd.BGPNetworkStatementState, len(nsRoutes))
	for idx, nsRoute := range nsRoutes {
		bgpdNSRoutes[idx] = h.ConvertBGPNetworkStateToThriftObject(nsRoute)
	}
	bgpRoutesBulk := bgpd.NewBGPNetworkStatementStateGetInfo()
	bgpRoutesBulk.EndIdx = bgpd.Int(nextIdx)
	bgpRoutesBulk.Count = bgpd.Int(currCount)
	bgpRoutesBulk.More = (nextIdx != 0)
	bgpRoutesBulk.BGPNetworkStatementStateList = bgpdNSRoutes

	return bgpRoutesBulk, nil
}

func (h *BGPHandler) ExecuteActionResetBGPv4NeighborByIPAddr(resetIP *bgpd.ResetBGPv4NeighborByIPAddr) (bool, error) {
	h.logger.Info("Reset BGP v4 neighbor by IP address", resetIP.IPAddr)
	if err := h.checkBGPGlobal(); err != nil {
		return false, err
	}

	ip := net.ParseIP(strings.TrimSpace(resetIP.IPAddr))
	if ip == nil {
		return false, errors.New(fmt.Sprintf("IPv4 Neighbor address %s is not a valid IP", resetIP.IPAddr))
	}
	h.server.PeerCommandCh <- config.PeerCommand{IP: ip, Command: int(fsm.BGPEventManualStop)}
	return true, nil
}

func (h *BGPHandler) ExecuteActionResetBGPv4NeighborByInterface(resetIf *bgpd.ResetBGPv4NeighborByInterface) (bool,
	error) {
	h.logger.Info("Reset BGP v4 neighbor by interface", resetIf.IntfRef)
	if err := h.checkBGPGlobal(); err != nil {
		return false, err
	}

	ifIndexInt, _, err := h.server.ConvertIntfStrToIfIndex(resetIf.IntfRef)
	if err != nil {
		h.logger.Err("Invalid intfref:", resetIf.IntfRef)
		return false, errors.New(fmt.Sprint("Invalid IntfRef ", resetIf.IntfRef))
	}

	ipInfo, err := h.server.GetIfaceIP(int32(ifIndexInt))
	h.logger.Info("ipInfo:", ipInfo, " err:", err)
	if err != nil {
		h.logger.Err("Failed to get IP for interface", resetIf.IntfRef)
		return false, errors.New(fmt.Sprintf("Failed to get IP for interface %s", resetIf.IntfRef))
	}

	if ipInfo.IpAddr == nil && ipInfo.IpMask == nil {
		return false, errors.New(fmt.Sprintln("IP address", ipInfo.IpAddr, "or netmask", ipInfo.IpMask,
			"of the interface", resetIf.IntfRef, "is not valid"))
	}

	ifIP := ipInfo.IpAddr
	ipMask := ipInfo.IpMask
	if ipMask[len(ipMask)-1] < 252 {
		h.logger.Err("IPv4Addr", ifIP, "of the interface", ifIndexInt, "is not /30 or /31 address")
		return false, errors.New(fmt.Sprintln("IPv4Addr", ifIP, "of the interface", resetIf.IntfRef,
			"is not /30 or /31 address"))
	}
	h.logger.Info("IPv4Addr of the v4Neighbor local interface", ifIndexInt, "is", ifIP)
	ifIP[len(ifIP)-1] = ifIP[len(ifIP)-1] ^ (^ipMask[len(ipMask)-1])
	h.logger.Info("IPv4Addr of the v4Neighbor remote interface is", ifIP)
	ip := ifIP

	h.server.PeerCommandCh <- config.PeerCommand{IP: ip, Command: int(fsm.BGPEventManualStop)}
	return true, nil
}

func (h *BGPHandler) ExecuteActionResetBGPv6NeighborByIPAddr(resetIP *bgpd.ResetBGPv6NeighborByIPAddr) (bool, error) {
	h.logger.Info("Reset BGP v6 neighbor by IP address", resetIP.IPAddr)
	if err := h.checkBGPGlobal(); err != nil {
		return false, err
	}

	ip := net.ParseIP(strings.TrimSpace(resetIP.IPAddr))
	if ip == nil {
		return false, errors.New(fmt.Sprintf("IPv6 Neighbor address %s is not a valid IP", resetIP.IPAddr))
	}
	h.server.PeerCommandCh <- config.PeerCommand{IP: ip, Command: int(fsm.BGPEventManualStop)}
	return true, nil
}

func (h *BGPHandler) ExecuteActionResetBGPv6NeighborByInterface(resetIf *bgpd.ResetBGPv6NeighborByInterface) (bool,
	error) {
	h.logger.Info("Reset BGP v6 neighbor by interface", resetIf.IntfRef)
	if err := h.checkBGPGlobal(); err != nil {
		return false, err
	}

	ifIndexInt, _, err := h.server.ConvertIntfStrToIfIndex(resetIf.IntfRef)
	if err != nil {
		h.logger.Err("Invalid intfref:", resetIf.IntfRef)
		return false, errors.New(fmt.Sprintf("Invalid IntfRef", resetIf.IntfRef))
	}

	ipInfo, err := h.server.GetIfaceIP(int32(ifIndexInt))
	h.logger.Info("ipInfo:", ipInfo, " err:", err)
	if err != nil {
		h.logger.Err("Failed to get IP for interface", resetIf.IntfRef)
		return false, errors.New(fmt.Sprintf("Failed to get IP for interface %s", resetIf.IntfRef))
	}

	h.logger.Info("Reset IPv6 neighbor by interface - ipInfo:%+v", ipInfo)
	if ipInfo.LinklocalIpAddr == "" {
		return false, errors.New(fmt.Sprintln("Link-local IP is not set on interface", resetIf.IntfRef))
	}

	ip := net.ParseIP(ipInfo.LinklocalIpAddr)

	if ip == nil {
		h.logger.Errf("IPv6 Neighbor address %s for interface %s is not a valid IP", ipInfo.LinklocalIpAddr,
			resetIf.IntfRef)
		return false, errors.New(fmt.Sprintf("IPv6 Neighbor address %s for interface %s is not a valid IP",
			ipInfo.LinklocalIpAddr, resetIf.IntfRef))
	}
	h.server.PeerCommandCh <- config.PeerCommand{IP: ip, Command: int(fsm.BGPEventManualStop)}
	return true, nil
}
