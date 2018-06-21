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

package api

import (
	"errors"
	"l3/ospfv2/objects"
	"l3/ospfv2/server"
)

var svr *server.OSPFV2Server

//Initialize server handle
func InitApiLayer(server *server.OSPFV2Server) {
	svr = server
}

func CreateOspfv2Area(cfg *objects.Ospfv2Area) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.CREATE_OSPFV2_AREA,
		Data: interface{}(&server.CreateOspfv2AreaInArgs{
			Cfg: cfg,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.CreateConfigOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response received from server during CreateArea")
}

func UpdateOspfv2Area(oldCfg, newCfg *objects.Ospfv2Area, attrset []bool) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.UPDATE_OSPFV2_AREA,
		Data: interface{}(&server.UpdateOspfv2AreaInArgs{
			OldCfg:  oldCfg,
			NewCfg:  newCfg,
			AttrSet: attrset,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.UpdateConfigOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response received from server during UpdateArea")
}

func DeleteOspfv2Area(cfg *objects.Ospfv2Area) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.DELETE_OSPFV2_AREA,
		Data: interface{}(&server.DeleteOspfv2AreaInArgs{
			Cfg: cfg,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.DeleteConfigOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response received from server during DeleteArea")
}

func GetOspfv2AreaState(areaId uint32) (*objects.Ospfv2AreaState, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_OSPFV2_AREA_STATE,
		Data: interface{}(&server.GetOspfv2AreaStateInArgs{
			AreaId: areaId,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetOspfv2AreaStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetOspfv2AreaState")
	}

}

func GetBulkOspfv2AreaState(fromIdx, count int) (*objects.Ospfv2AreaStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_OSPFV2_AREA_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkOspfv2AreaStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkOspfv2AreaState")
	}
}

func CreateOspfv2Global(cfg *objects.Ospfv2Global) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.CREATE_OSPFV2_GLOBAL,
		Data: interface{}(&server.CreateOspfv2GlobalInArgs{
			Cfg: cfg,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.CreateConfigOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response received from server during CreateGlobal")
}

func UpdateOspfv2Global(oldCfg, newCfg *objects.Ospfv2Global, attrset []bool) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.UPDATE_OSPFV2_GLOBAL,
		Data: interface{}(&server.UpdateOspfv2GlobalInArgs{
			OldCfg:  oldCfg,
			NewCfg:  newCfg,
			AttrSet: attrset,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.UpdateConfigOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response received from server during UpdateGlobal")
}

func DeleteOspfv2Global(cfg *objects.Ospfv2Global) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.DELETE_OSPFV2_GLOBAL,
		Data: interface{}(&server.DeleteOspfv2GlobalInArgs{
			Cfg: cfg,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.DeleteConfigOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response received from server during DeleteGlobal")
}

func GetOspfv2GlobalState(vrf string) (*objects.Ospfv2GlobalState, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_OSPFV2_GLOBAL_STATE,
		Data: interface{}(&server.GetOspfv2GlobalStateInArgs{
			Vrf: vrf,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetOspfv2GlobalStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetOspfv2GlobalState")
	}
}

func GetBulkOspfv2GlobalState(fromIdx, count int) (*objects.Ospfv2GlobalStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_OSPFV2_GLOBAL_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkOspfv2GlobalStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkOspfv2GlobalState")
	}
}

func CreateOspfv2Intf(cfg *objects.Ospfv2Intf) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.CREATE_OSPFV2_INTF,
		Data: interface{}(&server.CreateOspfv2IntfInArgs{
			Cfg: cfg,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.CreateConfigOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response received from server during CreateIntf")
}

func UpdateOspfv2Intf(oldCfg, newCfg *objects.Ospfv2Intf, attrset []bool) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.UPDATE_OSPFV2_INTF,
		Data: interface{}(&server.UpdateOspfv2IntfInArgs{
			OldCfg:  oldCfg,
			NewCfg:  newCfg,
			AttrSet: attrset,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.UpdateConfigOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response received from server during UpdateIntf")
}

func DeleteOspfv2Intf(cfg *objects.Ospfv2Intf) (bool, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.DELETE_OSPFV2_INTF,
		Data: interface{}(&server.DeleteOspfv2IntfInArgs{
			Cfg: cfg,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.DeleteConfigOutArgs); ok {
		return retObj.RetVal, retObj.Err
	}
	return false, errors.New("Error: Invalid response received from server during DeleteIntf")
	return true, nil
}

func GetOspfv2IntfState(ipAddr, addrLessIfIdx uint32) (*objects.Ospfv2IntfState, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_OSPFV2_INTF_STATE,
		Data: interface{}(&server.GetOspfv2IntfStateInArgs{
			IpAddr:           ipAddr,
			AddressLessIfIdx: addrLessIfIdx,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetOspfv2IntfStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetOspfv2IntfState")
	}
}

func GetBulkOspfv2IntfState(fromIdx, count int) (*objects.Ospfv2IntfStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_OSPFV2_INTF_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkOspfv2IntfStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkOspfv2IntfState")
	}
}

func GetOspfv2LsdbState(lsType uint8, lsId, areaId, advRtrId uint32) (*objects.Ospfv2LsdbState, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_OSPFV2_LSDB_STATE,
		Data: interface{}(&server.GetOspfv2LsdbStateInArgs{
			LSType:   lsType,
			LSId:     lsId,
			AreaId:   areaId,
			AdvRtrId: advRtrId,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetOspfv2LsdbStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetOspfv2LsdbState")
	}
}

func GetBulkOspfv2LsdbState(fromIdx, count int) (*objects.Ospfv2LsdbStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_OSPFV2_LSDB_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkOspfv2LsdbStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkOspfv2LsdbState")
	}
}

func GetOspfv2NbrState(ipAddr, addrLessIfIdx uint32) (*objects.Ospfv2NbrState, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_OSPFV2_NBR_STATE,
		Data: interface{}(&server.GetOspfv2NbrStateInArgs{
			IpAddr:           ipAddr,
			AddressLessIfIdx: addrLessIfIdx,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetOspfv2NbrStateOutArgs); ok {
		return retObj.Obj, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetOspfv2NbrState")
	}
}

func GetBulkOspfv2NbrState(fromIdx, count int) (*objects.Ospfv2NbrStateGetInfo, error) {
	svr.ReqChan <- &server.ServerRequest{
		Op: server.GET_BULK_OSPFV2_NBR_STATE,
		Data: interface{}(&server.GetBulkInArgs{
			FromIdx: fromIdx,
			Count:   count,
		}),
	}
	ret := <-svr.ReplyChan
	if retObj, ok := ret.(*server.GetBulkOspfv2NbrStateOutArgs); ok {
		return retObj.BulkInfo, retObj.Err
	} else {
		return nil, errors.New("Error: Invalid response received from server during GetBulkOspfv2NbrState")
	}
}

/*
func GetOspfv2RouteState(destId, addrMask, destType uint32) (*objects.Ospfv2RouteState, error) {
	return nil, nil
}

func GetBulkOspfv2RouteState(fromIdx, count int) (*objects.Ospfv2RouteStateGetInfo, error) {
	return nil, nil
}
*/
