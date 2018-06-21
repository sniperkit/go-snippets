//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
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
	"l3/ospfv2/objects"
)

type ServerOpId int

const (
	CREATE_OSPFV2_AREA ServerOpId = iota
	UPDATE_OSPFV2_AREA
	DELETE_OSPFV2_AREA
	GET_OSPFV2_AREA_STATE
	GET_BULK_OSPFV2_AREA_STATE
	CREATE_OSPFV2_GLOBAL
	UPDATE_OSPFV2_GLOBAL
	DELETE_OSPFV2_GLOBAL
	GET_OSPFV2_GLOBAL_STATE
	GET_BULK_OSPFV2_GLOBAL_STATE
	CREATE_OSPFV2_INTF
	UPDATE_OSPFV2_INTF
	DELETE_OSPFV2_INTF
	GET_OSPFV2_INTF_STATE
	GET_BULK_OSPFV2_INTF_STATE
	GET_OSPFV2_NBR_STATE
	GET_BULK_OSPFV2_NBR_STATE
	GET_OSPFV2_LSDB_STATE
	GET_BULK_OSPFV2_LSDB_STATE
)

type ServerRequest struct {
	Op   ServerOpId
	Data interface{}
}

type CreateConfigOutArgs struct {
	RetVal bool
	Err    error
}

type DeleteConfigOutArgs struct {
	RetVal bool
	Err    error
}

type UpdateConfigOutArgs struct {
	RetVal bool
	Err    error
}

type GetBulkInArgs struct {
	FromIdx int
	Count   int
}

type CreateOspfv2AreaInArgs struct {
	Cfg *objects.Ospfv2Area
}

type UpdateOspfv2AreaInArgs struct {
	OldCfg  *objects.Ospfv2Area
	NewCfg  *objects.Ospfv2Area
	AttrSet []bool
}

type DeleteOspfv2AreaInArgs struct {
	Cfg *objects.Ospfv2Area
}

type GetOspfv2AreaStateInArgs struct {
	AreaId uint32
}

type GetOspfv2AreaStateOutArgs struct {
	Obj *objects.Ospfv2AreaState
	Err error
}

type GetBulkOspfv2AreaStateOutArgs struct {
	BulkInfo *objects.Ospfv2AreaStateGetInfo
	Err      error
}

type CreateOspfv2GlobalInArgs struct {
	Cfg *objects.Ospfv2Global
}

type UpdateOspfv2GlobalInArgs struct {
	OldCfg  *objects.Ospfv2Global
	NewCfg  *objects.Ospfv2Global
	AttrSet []bool
}

type DeleteOspfv2GlobalInArgs struct {
	Cfg *objects.Ospfv2Global
}

type GetOspfv2GlobalStateInArgs struct {
	Vrf string
}

type GetOspfv2GlobalStateOutArgs struct {
	Obj *objects.Ospfv2GlobalState
	Err error
}

type GetBulkOspfv2GlobalStateOutArgs struct {
	BulkInfo *objects.Ospfv2GlobalStateGetInfo
	Err      error
}

type CreateOspfv2IntfInArgs struct {
	Cfg *objects.Ospfv2Intf
}

type UpdateOspfv2IntfInArgs struct {
	OldCfg  *objects.Ospfv2Intf
	NewCfg  *objects.Ospfv2Intf
	AttrSet []bool
}

type DeleteOspfv2IntfInArgs struct {
	Cfg *objects.Ospfv2Intf
}

type GetOspfv2IntfStateInArgs struct {
	IpAddr           uint32
	AddressLessIfIdx uint32
}

type GetOspfv2IntfStateOutArgs struct {
	Obj *objects.Ospfv2IntfState
	Err error
}

type GetBulkOspfv2IntfStateOutArgs struct {
	BulkInfo *objects.Ospfv2IntfStateGetInfo
	Err      error
}

type GetOspfv2NbrStateInArgs struct {
	IpAddr           uint32
	AddressLessIfIdx uint32
}

type GetOspfv2NbrStateOutArgs struct {
	Obj *objects.Ospfv2NbrState
	Err error
}

type GetBulkOspfv2NbrStateOutArgs struct {
	BulkInfo *objects.Ospfv2NbrStateGetInfo
	Err      error
}

type GetOspfv2LsdbStateInArgs struct {
	LSType   uint8
	LSId     uint32
	AreaId   uint32
	AdvRtrId uint32
}

type GetOspfv2LsdbStateOutArgs struct {
	Obj *objects.Ospfv2LsdbState
	Err error
}

type GetBulkOspfv2LsdbStateOutArgs struct {
	BulkInfo *objects.Ospfv2LsdbStateGetInfo
	Err      error
}
