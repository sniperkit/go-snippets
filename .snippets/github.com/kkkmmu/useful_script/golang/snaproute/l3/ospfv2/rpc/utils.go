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

package rpc

import (
	"errors"
	"fmt"
	"l3/ospfv2/objects"
	"net"
	"ospfv2d"
	"strconv"
	"strings"
)

func convertDotNotationToUint32(str string) (uint32, error) {
	var val uint32
	ip := net.ParseIP(str)
	if ip == nil {
		return 0, errors.New("Invalid string format")
	}
	ipBytes := ip.To4()
	val = val + uint32(ipBytes[0])
	val = (val << 8) + uint32(ipBytes[1])
	val = (val << 8) + uint32(ipBytes[2])
	val = (val << 8) + uint32(ipBytes[3])
	return val, nil
}

func convertUint32ToDotNotation(val uint32) string {
	p0 := int(val & 0xFF)
	p1 := int((val >> 8) & 0xFF)
	p2 := int((val >> 16) & 0xFF)
	p3 := int((val >> 24) & 0xFF)
	str := strconv.Itoa(p3) + "." + strconv.Itoa(p2) + "." +
		strconv.Itoa(p1) + "." + strconv.Itoa(p0)

	return str
}

func convertFromRPCFmtOspfv2Area(config *ospfv2d.Ospfv2Area) (*objects.Ospfv2Area, error) {
	areaId, err := convertDotNotationToUint32(config.AreaId)
	if err != nil {
		return nil, errors.New(fmt.Sprintln("Invalid AreaId", err))
	}
	var adminState bool
	switch strings.ToLower(config.AdminState) {
	case objects.AREA_ADMIN_STATE_UP_STR:
		adminState = objects.AREA_ADMIN_STATE_UP
	case objects.AREA_ADMIN_STATE_DOWN_STR:
		adminState = objects.AREA_ADMIN_STATE_DOWN
	default:
		return nil, errors.New("Invalid AdminState")
	}
	var authType uint8
	switch strings.ToLower(config.AuthType) {
	case objects.AUTH_TYPE_NONE_STR:
		authType = objects.AUTH_TYPE_NONE
	case objects.AUTH_TYPE_SIMPLE_PASSWORD_STR:
		authType = objects.AUTH_TYPE_SIMPLE_PASSWORD
	case objects.AUTH_TYPE_MD5_STR:
		authType = objects.AUTH_TYPE_MD5
	default:
		return nil, errors.New("Invalid Auth Type")
	}
	return &objects.Ospfv2Area{
		AreaId:         areaId,
		AdminState:     adminState,
		AuthType:       authType,
		ImportASExtern: config.ImportASExtern,
	}, nil
}

func convertToRPCFmtOspfv2AreaState(obj *objects.Ospfv2AreaState) *ospfv2d.Ospfv2AreaState {
	areaId := convertUint32ToDotNotation(obj.AreaId)
	return &ospfv2d.Ospfv2AreaState{
		AreaId: areaId,
		//NumSpfRuns:       int32(obj.NumSpfRuns),
		//NumBdrRtr:        int32(obj.NumBdrRtr),
		//NumAsBdrRtr:      int32(obj.NumAsBdrRtr),
		NumOfRouterLSA:     int32(obj.NumOfRouterLSA),
		NumOfNetworkLSA:    int32(obj.NumOfNetworkLSA),
		NumOfSummary3LSA:   int32(obj.NumOfSummary3LSA),
		NumOfSummary4LSA:   int32(obj.NumOfSummary4LSA),
		NumOfASExternalLSA: int32(obj.NumOfASExternalLSA),
		NumOfIntfs:         int32(obj.NumOfIntfs),
		NumOfNbrs:          int32(obj.NumOfNbrs),
		NumOfLSA:           int32(obj.NumOfLSA),
		NumOfRoutes:        int32(obj.NumOfRoutes),
	}
}

func convertFromRPCFmtOspfv2Global(config *ospfv2d.Ospfv2Global) (*objects.Ospfv2Global, error) {
	routerId, err := convertDotNotationToUint32(config.RouterId)
	if err != nil {
		return nil, errors.New(fmt.Sprintln("Invalid RouterId", err))
	}
	var adminState bool
	switch strings.ToLower(config.AdminState) {
	case objects.GLOBAL_ADMIN_STATE_UP_STR:
		adminState = objects.GLOBAL_ADMIN_STATE_UP
	case objects.GLOBAL_ADMIN_STATE_DOWN_STR:
		adminState = objects.GLOBAL_ADMIN_STATE_DOWN
	default:
		return nil, errors.New("Invalid AdminState")
	}
	// Skipping VRF for now
	if config.Vrf != "default" {
		return nil, errors.New("Invalid Vrf")
	}
	return &objects.Ospfv2Global{
		Vrf:                "default",
		RouterId:           routerId,
		AdminState:         adminState,
		ASBdrRtrStatus:     config.ASBdrRtrStatus,
		ReferenceBandwidth: uint32(config.ReferenceBandwidth),
	}, nil
}

func convertToRPCFmtOspfv2GlobalState(obj *objects.Ospfv2GlobalState) *ospfv2d.Ospfv2GlobalState {
	return &ospfv2d.Ospfv2GlobalState{
		Vrf:                "default",
		AreaBdrRtrStatus:   obj.AreaBdrRtrStatus,
		NumOfAreas:         int32(obj.NumOfAreas),
		NumOfIntfs:         int32(obj.NumOfIntfs),
		NumOfNbrs:          int32(obj.NumOfNbrs),
		NumOfLSA:           int32(obj.NumOfLSA),
		NumOfRouterLSA:     int32(obj.NumOfRouterLSA),
		NumOfNetworkLSA:    int32(obj.NumOfNetworkLSA),
		NumOfSummary3LSA:   int32(obj.NumOfSummary3LSA),
		NumOfSummary4LSA:   int32(obj.NumOfSummary4LSA),
		NumOfASExternalLSA: int32(obj.NumOfASExternalLSA),
		NumOfRoutes:        int32(obj.NumOfRoutes),
	}
}

func convertFromRPCFmtOspfv2Intf(config *ospfv2d.Ospfv2Intf) (*objects.Ospfv2Intf, error) {
	ipAddr, err := convertDotNotationToUint32(config.IpAddress)
	if err != nil {
		return nil, errors.New(fmt.Sprintln("Invalid IpAddress", err))
	}
	areaId, err := convertDotNotationToUint32(config.AreaId)
	if err != nil {
		return nil, errors.New(fmt.Sprintln("Invalid AreaId", err))
	}
	var adminState bool
	switch strings.ToLower(config.AdminState) {
	case objects.INTF_ADMIN_STATE_UP_STR:
		adminState = objects.INTF_ADMIN_STATE_UP
	case objects.INTF_ADMIN_STATE_DOWN_STR:
		adminState = objects.INTF_ADMIN_STATE_DOWN
	default:
		return nil, errors.New("Invalid AdminState")
	}
	var intfType uint8
	switch strings.ToLower(config.Type) {
	case objects.INTF_TYPE_POINT2POINT_STR:
		intfType = objects.INTF_TYPE_POINT2POINT
	case objects.INTF_TYPE_BROADCAST_STR:
		intfType = objects.INTF_TYPE_BROADCAST
	default:
		return nil, errors.New("Invalid Interface Type")
	}
	return &objects.Ospfv2Intf{
		IpAddress:        ipAddr,
		AddressLessIfIdx: uint32(config.AddressLessIfIdx),
		AdminState:       adminState,
		AreaId:           areaId,
		Type:             intfType,
		RtrPriority:      uint8(config.RtrPriority),
		TransitDelay:     uint16(config.TransitDelay),
		RetransInterval:  uint16(config.RetransInterval),
		HelloInterval:    uint16(config.HelloInterval),
		RtrDeadInterval:  uint32(config.RtrDeadInterval),
		MetricValue:      uint16(config.MetricValue),
	}, nil
}

func convertToRPCFmtOspfv2IntfState(obj *objects.Ospfv2IntfState) *ospfv2d.Ospfv2IntfState {
	ipAddr := convertUint32ToDotNotation(obj.IpAddress)
	var state string
	switch obj.State {
	case objects.INTF_FSM_STATE_OTHER_DR:
		state = strings.ToUpper(objects.INTF_FSM_STATE_OTHER_DR_STR)
	case objects.INTF_FSM_STATE_DR:
		state = strings.ToUpper(objects.INTF_FSM_STATE_DR_STR)
	case objects.INTF_FSM_STATE_BDR:
		state = strings.ToUpper(objects.INTF_FSM_STATE_BDR_STR)
	case objects.INTF_FSM_STATE_LOOPBACK:
		state = strings.ToUpper(objects.INTF_FSM_STATE_LOOPBACK_STR)
	case objects.INTF_FSM_STATE_DOWN:
		state = strings.ToUpper(objects.INTF_FSM_STATE_DOWN_STR)
	case objects.INTF_FSM_STATE_WAITING:
		state = strings.ToUpper(objects.INTF_FSM_STATE_WAITING_STR)
	case objects.INTF_FSM_STATE_P2P:
		state = strings.ToUpper(objects.INTF_FSM_STATE_P2P_STR)
	}
	designatedRouter := convertUint32ToDotNotation(obj.DesignatedRouter)
	designatedRouterId := convertUint32ToDotNotation(obj.DesignatedRouterId)
	backupDesignatedRouter := convertUint32ToDotNotation(obj.BackupDesignatedRouter)
	backupDesignatedRouterId := convertUint32ToDotNotation(obj.BackupDesignatedRouterId)
	return &ospfv2d.Ospfv2IntfState{
		IpAddress:                ipAddr,
		AddressLessIfIdx:         int32(obj.AddressLessIfIdx),
		State:                    state,
		DesignatedRouter:         designatedRouter,
		DesignatedRouterId:       designatedRouterId,
		BackupDesignatedRouter:   backupDesignatedRouter,
		BackupDesignatedRouterId: backupDesignatedRouterId,
		NumOfRouterLSA:           int32(obj.NumOfRouterLSA),
		NumOfNetworkLSA:          int32(obj.NumOfNetworkLSA),
		NumOfSummary3LSA:         int32(obj.NumOfSummary3LSA),
		NumOfSummary4LSA:         int32(obj.NumOfSummary4LSA),
		NumOfASExternalLSA:       int32(obj.NumOfASExternalLSA),
		NumOfLSA:                 int32(obj.NumOfLSA),
		NumOfNbrs:                int32(obj.NumOfNbrs),
		NumOfRoutes:              int32(obj.NumOfRoutes),
		Mtu:                      int32(obj.Mtu),
		Cost:                     int32(obj.Cost),
		NumOfStateChange:         int32(obj.NumOfStateChange),
		TimeOfStateChange:        obj.TimeOfStateChange,
	}
}

func convertFromRPCFmtLSType(LSType string) (uint8, error) {
	var lsType uint8

	switch strings.ToLower(LSType) {
	case objects.ROUTER_LSA_STR:
		lsType = objects.ROUTER_LSA
	case objects.NETWORK_LSA_STR:
		lsType = objects.NETWORK_LSA
	case objects.SUMMARY3_LSA_STR:
		lsType = objects.SUMMARY3_LSA
	case objects.SUMMARY4_LSA_STR:
		lsType = objects.SUMMARY4_LSA
	case objects.ASExternal_LSA_STR:
		lsType = objects.ASExternal_LSA
	default:
		return 0, errors.New("Invalid LSA Type")
	}
	return lsType, nil
}

func convertToRPCFmtOspfv2LsdbState(obj *objects.Ospfv2LsdbState) *ospfv2d.Ospfv2LsdbState {
	var lsType string
	switch obj.LSType {
	case objects.ROUTER_LSA:
		lsType = strings.ToUpper(objects.ROUTER_LSA_STR)
	case objects.NETWORK_LSA:
		lsType = strings.ToUpper(objects.NETWORK_LSA_STR)
	case objects.SUMMARY3_LSA:
		lsType = strings.ToUpper(objects.SUMMARY3_LSA_STR)
	case objects.SUMMARY4_LSA:
		lsType = strings.ToUpper(objects.SUMMARY4_LSA_STR)
	case objects.ASExternal_LSA:
		lsType = strings.ToUpper(objects.ASExternal_LSA_STR)
	}
	lsId := convertUint32ToDotNotation(obj.LSId)
	areaId := convertUint32ToDotNotation(obj.AreaId)
	advRtrId := convertUint32ToDotNotation(obj.AdvRouterId)
	seqNum := fmt.Sprintf("0x%X", obj.SequenceNum)
	return &ospfv2d.Ospfv2LsdbState{
		LSType:        lsType,
		LSId:          lsId,
		AreaId:        areaId,
		AdvRouterId:   advRtrId,
		SequenceNum:   seqNum,
		Age:           int16(obj.Age),
		Checksum:      int16(obj.Checksum),
		Options:       int8(obj.Options),
		Length:        int16(obj.Length),
		Advertisement: obj.Advertisement,
	}
}

func convertToRPCFmtOspfv2NbrState(obj *objects.Ospfv2NbrState) *ospfv2d.Ospfv2NbrState {
	ipAddr := convertUint32ToDotNotation(obj.IpAddr)
	rtrId := convertUint32ToDotNotation(obj.RtrId)
	var state string
	switch obj.State {
	case objects.NBR_STATE_ONE_WAY:
		state = strings.ToUpper(objects.NBR_STATE_ONE_WAY_STR)
	case objects.NBR_STATE_TWO_WAY:
		state = strings.ToUpper(objects.NBR_STATE_TWO_WAY_STR)
	case objects.NBR_STATE_INIT:
		state = strings.ToUpper(objects.NBR_STATE_INIT_STR)
	case objects.NBR_STATE_EXSTART:
		state = strings.ToUpper(objects.NBR_STATE_EXSTART_STR)
	case objects.NBR_STATE_EXCHANGE:
		state = strings.ToUpper(objects.NBR_STATE_EXCHANGE_STR)
	case objects.NBR_STATE_LOADING:
		state = strings.ToUpper(objects.NBR_STATE_LOADING_STR)
	case objects.NBR_STATE_ATTEMPT:
		state = strings.ToUpper(objects.NBR_STATE_ATTEMPT_STR)
	case objects.NBR_STATE_DOWN:
		state = strings.ToUpper(objects.NBR_STATE_DOWN_STR)
	case objects.NBR_STATE_FULL:
		state = strings.ToUpper(objects.NBR_STATE_FULL_STR)
	}
	return &ospfv2d.Ospfv2NbrState{
		IpAddr:           ipAddr,
		AddressLessIfIdx: int32(obj.AddressLessIfIdx),
		RtrId:            rtrId,
		Options:          obj.Options,
		State:            state,
	}
}

/*
func convertToRPCFmtOspfv2RouteState(obj *objects.Ospfv2RouteState) *ospfv2d.Ospfv2RouteState {
	return &ospfv2d.Ospfv2RouteState{}
}
*/
