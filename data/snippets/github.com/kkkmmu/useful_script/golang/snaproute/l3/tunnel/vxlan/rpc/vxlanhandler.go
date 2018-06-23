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

// lahandler
package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	vxlan "l3/tunnel/vxlan/protocol"
	"l3/tunnel/vxlan/server"
	"models/objects"
	"reflect"
	"utils/dbutils"
	"utils/logging"
	"vxland"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/garyburd/redigo/redis"
)

const DBName string = "UsrConfDb.db"

type ClientJson struct {
	Name string `json:Name`
	Port int    `json:Port`
}

type VXLANDServiceHandler struct {
	server       *server.VXLANServer
	logger       *logging.Writer
	Thriftserver *thrift.TSimpleServer
	running      bool
	readdb       bool
}

// look up the various other daemons based on c string
func GetClientPort(paramsFile string, c string) int {
	var clientsList []ClientJson

	bytes, err := ioutil.ReadFile(paramsFile)
	if err != nil {
		return 0
	}

	err = json.Unmarshal(bytes, &clientsList)
	if err != nil {
		return 0
	}

	for _, client := range clientsList {
		if client.Name == c {
			return client.Port
		}
	}
	return 0
}

func NewVXLANDServiceHandler(server *server.VXLANServer, logger *logging.Writer) *VXLANDServiceHandler {
	//lacp.LacpStartTime = time.Now()
	// link up/down events for now
	//startEvtHandler()
	handler := &VXLANDServiceHandler{
		server: server,
		logger: logger,
	}

	handler.CreateThriftServer()

	return handler
}

func (v *VXLANDServiceHandler) CreateThriftServer() {

	var transport thrift.TServerTransport
	var err error

	fileName := v.server.Paramspath + "clients.json"
	port := GetClientPort(fileName, "vxland")
	if port != 0 {
		addr := fmt.Sprintf("localhost:%d", port)
		transport, err = thrift.NewTServerSocket(addr)
		if err != nil {
			panic(fmt.Sprintf("Failed to create Socket with:", addr))
		}

		processor := vxland.NewVXLANDServicesProcessor(v)
		transportFactory := thrift.NewTBufferedTransportFactory(8192)
		protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
		v.Thriftserver = thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	}
}

// StopThriftServer for purposes of stopping vxlan config from coming from confd
func (v *VXLANDServiceHandler) StopCfgServerLoop() {
	if v.Thriftserver != nil {
		// Not going to do anything as we don't want to edit thrift code at this time.
		// See thrift code changes
		//err := v.Thriftserver.Stop()
		//v.running = False
		//v.logger.Info("Stopping Cfg Server loop", err)
		v.logger.Info("Stopping Cfg Server loop")
	}
}

// Enable listening of server config
func (v *VXLANDServiceHandler) StartCfgServerLoop() {
	if v.Thriftserver != nil && !v.running {
		v.logger.Info("Starting Cfg Server loop")
		// Not going to do anything as we don't want to edit thrift code at this time.
		// See thrift code changes

		// lets read the db on creation, necessary for restart
		if !v.readdb {
			prevState := vxlan.VxlanGlobalStateGet()
			v.ReadConfigFromDB(prevState)
			v.readdb = true
		}

		v.running = true
		err := v.Thriftserver.Serve()
		v.logger.Info("Cfg Server loop stopped", err)
	}
}

func (v *VXLANDServiceHandler) CreateVxlanGlobal(config *vxland.VxlanGlobal) (rv bool, err error) {
	rv = true
	v.logger.Info(fmt.Sprintf("CreateVxlanGlobal (server): %s", config.AdminState))

	if config.AdminState == "UP" {
		vxlan.VxlanGlobalStateSet(vxlan.VXLAN_GLOBAL_ENABLE)
	} else if config.AdminState == "DOWN" {
		vxlan.VxlanGlobalStateSet(vxlan.VXLAN_GLOBAL_DISABLE)
	} else {
		return rv, errors.New(fmt.Sprintln("Error VxlanGlobal unknown Admin State setting", config.AdminState))
	}

	return rv, err
}

func (v *VXLANDServiceHandler) DeleteVxlanGlobal(config *vxland.VxlanGlobal) (bool, error) {
	return true, nil
}

func (v *VXLANDServiceHandler) UpdateVxlanGlobal(origconfig *vxland.VxlanGlobal, updateconfig *vxland.VxlanGlobal, attrset []bool, op []*vxland.PatchOpInfo) (rv bool, err error) {
	v.logger.Info(fmt.Sprintf("UpdateVxlanGlobal (server): %s", updateconfig.AdminState))
	rv = true
	prevState := vxlan.VxlanGlobalStateGet()

	if updateconfig.AdminState == "UP" {
		vxlan.VxlanGlobalStateSet(vxlan.VXLAN_GLOBAL_ENABLE)
	} else if updateconfig.AdminState == "DOWN" {
		vxlan.VxlanGlobalStateSet(vxlan.VXLAN_GLOBAL_DISABLE_PENDING)
	} else {
		return rv, errors.New(fmt.Sprintln("Error Update VxlanGlobal unknown Admin State setting", updateconfig.AdminState))
	}
	if prevState != vxlan.VxlanGlobalStateGet() {
		v.ReadConfigFromDB(prevState)
	}
	return rv, err
}

func (v *VXLANDServiceHandler) GetBulkVxlanGlobalState(fromIndex vxland.Int, count vxland.Int) (obj *vxland.VxlanGlobalStateGetInfo, err error) {
	var returnVxlanGlobalStates []*vxland.VxlanGlobalState
	var returnVxlanGlobalStateGetInfo vxland.VxlanGlobalStateGetInfo
	toIndex := fromIndex
	obj = &returnVxlanGlobalStateGetInfo

	nextVxlanGlobalState, gserr := v.GetVxlanGlobalState("default")
	if gserr == nil {
		if len(returnVxlanGlobalStates) == 0 {
			returnVxlanGlobalStates = make([]*vxland.VxlanGlobalState, 0)
		}
		returnVxlanGlobalStates = append(returnVxlanGlobalStates, nextVxlanGlobalState)
	}
	obj.VxlanGlobalStateList = returnVxlanGlobalStates
	obj.StartIdx = fromIndex
	obj.EndIdx = toIndex + 1
	obj.More = false
	obj.Count = 1
	return obj, err
}
func (v *VXLANDServiceHandler) GetVxlanGlobalState(Vrf string) (*vxland.VxlanGlobalState, error) {

	vg := &vxland.VxlanGlobalState{}

	state := vxlan.VxlanGlobalStateGet()
	if state == vxlan.VXLAN_GLOBAL_ENABLE {
		vg.OperState = "UP"
	} else {
		vg.OperState = "DOWN"
	}
	vg.RxInvalidVtepCnt = 0
	vg.NumVteps = int64(len(vxlan.GetVtepDB()))
	return vg, nil
}

func (v *VXLANDServiceHandler) CreateVxlanInstance(config *vxland.VxlanInstance) (bool, error) {
	v.logger.Info(fmt.Sprintf("CreateVxlanConfigInstance %#v", config))

	c, err := vxlan.ConvertVxlanInstanceToVxlanConfig(config, true)
	if err == nil {
		err = vxlan.VxlanConfigCheck(c)
		if err == nil {
			v.server.Configchans.Vxlancreate <- *c
			return true, nil
		}
	}
	return false, err
}

func (v *VXLANDServiceHandler) DeleteVxlanInstance(config *vxland.VxlanInstance) (bool, error) {
	v.logger.Info(fmt.Sprintf("DeleteVxlanConfigInstance %#v", config))
	c, err := vxlan.ConvertVxlanInstanceToVxlanConfig(config, false)
	if err == nil {
		v.server.Configchans.Vxlandelete <- *c
		return true, nil
	}
	return false, err
}

func (v *VXLANDServiceHandler) UpdateVxlanInstance(origconfig *vxland.VxlanInstance, newconfig *vxland.VxlanInstance, attrset []bool, op []*vxland.PatchOpInfo) (bool, error) {
	v.logger.Info(fmt.Sprintf("UpdateVxlanConfigInstance orig[%#v] new[%#v] attrset[%#v]", origconfig, newconfig, attrset))
	oc, err1 := vxlan.ConvertVxlanInstanceToVxlanConfig(origconfig, false)
	nc, err2 := vxlan.ConvertVxlanInstanceToVxlanConfig(newconfig, false)
	if err1 == nil &&
		err2 == nil {
		err := vxlan.VxlanConfigUpdateCheck(oc, nc)
		if err == nil {

			strattr := make([]string, 0)
			objTyp := reflect.TypeOf(*origconfig)

			// important to note that the attrset starts at index 0 which is the BaseObj
			// which is not the first element on the thrift obj, thus we need to skip
			// this attribute
			for i := 0; i < objTyp.NumField(); i++ {
				objName := objTyp.Field(i).Name
				if attrset[i] {
					strattr = append(strattr, objName)
				}
			}
			update := vxlan.VxlanUpdate{
				Oldconfig: *oc,
				Newconfig: *nc,
				Attr:      strattr,
			}
			v.server.Configchans.Vxlanupdate <- update
			return true, nil
		} else {
			return false, err
		}
	}
	if err1 != nil {
		return false, err1
	}

	return false, err2
}

func (v *VXLANDServiceHandler) CreateVxlanVtepInstance(config *vxland.VxlanVtepInstance) (bool, error) {
	v.logger.Info(fmt.Sprintf("CreateVxlanVtepInstance %#v", config))
	cs, err := vxlan.ConvertVxlanVtepInstanceToVtepConfig(config)
	if err == nil {
		for _, c := range cs {
			err = vxlan.VtepConfigCheck(c, true)
			if err != nil {
	                	return false, err
			}
			v.server.Configchans.Vtepcreate <- *c
		}
	}
	return true, err
}

func (v *VXLANDServiceHandler) DeleteVxlanVtepInstance(config *vxland.VxlanVtepInstance) (bool, error) {
	v.logger.Info(fmt.Sprintf("DeleteVxlanVtepInstance %#v", config))
	cs, err := vxlan.ConvertVxlanVtepInstanceToVtepConfig(config)
	if err == nil {
		for _, c := range cs {
			v.server.Configchans.Vtepdelete <- *c
			return true, nil
		}
	}
	return false, err
}

func (v *VXLANDServiceHandler) UpdateVxlanVtepInstance(origconfig *vxland.VxlanVtepInstance, newconfig *vxland.VxlanVtepInstance, attrset []bool, op []*vxland.PatchOpInfo) (bool, error) {
	v.logger.Info(fmt.Sprintf("UpdateVxlanVtepInstances orig[%#v] new[%#v]", origconfig, newconfig))
	ocs, _ := vxlan.ConvertVxlanVtepInstanceToVtepConfig(origconfig)
	ncs, err := vxlan.ConvertVxlanVtepInstanceToVtepConfig(newconfig)
	if err == nil {
		for _, c := range ncs {
			err = vxlan.VtepConfigCheck(c, false)
			if err == nil {
				strattr := make([]string, 0)
				objTyp := reflect.TypeOf(*origconfig)

				// important to note that the attrset starts at index 0 which is the BaseObj
				// which is not the first element on the thrift obj, thus we need to skip
				// this attribute
				for i := 0; i < objTyp.NumField(); i++ {
					objName := objTyp.Field(i).Name
					if attrset[i] {
						if objName == "DstIp" {
							continue
						}
						strattr = append(strattr, objName)
					}
				}
				foundDstIp := false
				for _, oc := range ocs {

					if oc.TunnelDstIp.String() == c.TunnelDstIp.String() {
						foundDstIp = true
						update := vxlan.VtepUpdate{
							Oldconfig: *oc,
							Newconfig: *c,
							Attr:      strattr,
						}
						v.server.Configchans.Vtepupdate <- update
						break
					}
				}
				// we have a new vtep
				if !foundDstIp {
					v.server.Configchans.Vtepcreate <- *c
				}
			}
		}

		for _, c := range ocs {
			foundDstIp := false
			for _, nc := range ncs {
				if nc.TunnelDstIp.String() == c.TunnelDstIp.String() {
					foundDstIp = true
				}
			}
			// delete the vtep
			if !foundDstIp {
				v.server.Configchans.Vtepdelete <- *c
			}
		}

		return true, nil
	}

	return false, err
}

func (v *VXLANDServiceHandler) GetVxlanInstanceState(vni int32) (*vxland.VxlanInstanceState, error) {
	vis := &vxland.VxlanInstanceState{}
	if v, ok := vxlan.GetVxlanDB()[uint32(vni)]; ok {
		vis.Vni = int32(vni)
		if v.Enable {
			vis.OperState = "UP"
		} else {
			vis.OperState = "DOWN"
		}
		for _, vlan := range v.VlanId {
			vis.VlanId = append(vis.VlanId, int16(vlan))
		}

	} else {
		return nil, errors.New(fmt.Sprintf("Error could not find vni instance %d", vni))
	}
	return vis, nil
}

func (la *VXLANDServiceHandler) GetBulkVxlanInstanceState(fromIndex vxland.Int, count vxland.Int) (obj *vxland.VxlanInstanceStateGetInfo, err error) {

	var vxlanStateList []vxland.VxlanInstanceState = make([]vxland.VxlanInstanceState, count)
	var nextVxlanState *vxland.VxlanInstanceState
	var returnVxlanStates []*vxland.VxlanInstanceState
	var returnVxlanStateGetInfo vxland.VxlanInstanceStateGetInfo
	validCount := vxland.Int(0)
	toIndex := fromIndex
	moreRoutes := false
	obj = &returnVxlanStateGetInfo

	var v *vxlan.VxlanDbEntry
	currIndex := vxland.Int(0)

	for currIndex = vxland.Int(0); validCount != count && vxlan.GetVxlanDbListEntry(int32(currIndex), &v); currIndex++ {

		if currIndex < fromIndex {
			continue
		} else {

			nextVxlanState = &vxlanStateList[validCount]
			if v.Enable {
				nextVxlanState.OperState = "UP"
			} else {
				nextVxlanState.OperState = "DOWN"
			}
			nextVxlanState.Vni = int32(v.VNI)
			for _, vlan := range v.VlanId {
				nextVxlanState.VlanId = append(nextVxlanState.VlanId, int16(vlan))
			}
			if len(returnVxlanStates) == 0 {
				returnVxlanStates = make([]*vxland.VxlanInstanceState, 0)
			}
			returnVxlanStates = append(returnVxlanStates, nextVxlanState)
			validCount++
			toIndex++
		}
	}
	// lets try and get the next agg if one exists then there are more routes
	if v != nil {
		moreRoutes = vxlan.GetVxlanDbListEntry(int32(currIndex), &v)
	}
	obj.VxlanInstanceStateList = returnVxlanStates
	obj.StartIdx = fromIndex
	obj.EndIdx = toIndex + 1
	obj.More = moreRoutes
	obj.Count = validCount

	return obj, nil
}

func (v *VXLANDServiceHandler) GetVxlanVtepInstanceState(intf string, vni int32) (*vxland.VxlanVtepInstanceState, error) {
	vis := &vxland.VxlanVtepInstanceState{}
	//key := vxlan.VtepDbKey{
	//	Name: intf,
	//	Vni:  uint32(vni),
	//}

	for _, v := range vxlan.GetVtepDB() {
                
		if v.VtepConfigName == intf &&
			v.Vni == uint32(vni) {
			OperState := "UNKNOWN"
			if v.Enable && v.VxlanVtepMachineFsm.Machine.Curr.CurrentState() == vxlan.VxlanVtepStateStart {
				OperState = "UP"
			} else {
				OperState = "DOWN"
			}
			vis.Intf = v.VtepConfigName
			vis.IntfRef = v.SrcIfName
			vis.IfIndex = v.VtepIfIndex
			vis.Vni = int32(v.Vni)
			vis.DstUDP = int16(v.UDP)
			vis.TTL = int16(v.TTL)
			vis.TOS = int16(v.TOS)
			if v.FilterUnknownCustVlan {
				vis.InnerVlanHandlingMode = 0
			} else {
				vis.InnerVlanHandlingMode = 1
			}

			vis.SrcIp = v.SrcIp.String()
			vis.VlanId = int16(v.VlanId)
			vis.Mtu = int32(v.MTU)
			vis.PerDstIpState = append(vis.PerDstIpState, &vxland.VtepStateEntry{
				SubVtepName:        v.VtepName,
				DstIp:              v.DstIp.String(),
				RxSwPkts:           int64(v.GetStats().Rxpkts),
				RxSwBytes:          int64(v.GetStats().Rxbytes),
				RxSwDropPkts:       int64(v.GetStats().Rxdroppkts),
				RxSwDropBytes:      int64(v.GetStats().Rxdropbytes),
				RxSwFwdPkts:        int64(v.GetStats().Rxfwdpkts),
				RxSwFwdBytes:       int64(v.GetStats().Rxfwdbytes),
				TxSwPkts:           int64(v.GetStats().Txpkts),
				TxSwBytes:          int64(v.GetStats().Txbytes),
				TxSwDropPkts:       int64(v.GetStats().Txdroppkts),
				TxSwDropBytes:      int64(v.GetStats().Txdropbytes),
				TxSwFwdPkts:        int64(v.GetStats().Txfwdpkts),
				TxSwFwdBytes:       int64(v.GetStats().Txfwdbytes),
				LastSwRxDropReason: v.GetStats().Lastrxdropreason,
				LastSwTxDropReason: v.GetStats().Lasttxdropreason,
				//vis.RxFwdPkts             uint64 `DESCRIPTION: Rx Forwaded Packets`
				//vis.RxDropPkts            uint64 `DESCRIPTION: Rx Dropped Packets`
				//vis.RxUnknownVni          uint64 `DESCRIPTION: Rx Unknown Vni in frame`
				VtepFsmState:     vxlan.VxlanVtepStateStrMap[v.VxlanVtepMachineFsm.Machine.Curr.CurrentState()],
				VtepFsmPrevState: vxlan.VxlanVtepStateStrMap[v.VxlanVtepMachineFsm.Machine.Curr.PreviousState()],
				OperState:        OperState,
			})

		} else {
			return nil, errors.New(fmt.Sprintf("Error could not find vni instance %s", intf))
		}
	}
	return vis, nil
}

func (la *VXLANDServiceHandler) GetBulkVxlanVtepInstanceState(fromIndex vxland.Int, count vxland.Int) (obj *vxland.VxlanVtepInstanceStateGetInfo, err error) {

	var vxlanVtepStateList []vxland.VxlanVtepInstanceState = make([]vxland.VxlanVtepInstanceState, count)
	var nextVxlanVtepState *vxland.VxlanVtepInstanceState
	var returnVxlanVtepStates []*vxland.VxlanVtepInstanceState
	var returnVxlanVtepStateGetInfo vxland.VxlanVtepInstanceStateGetInfo
	validCount := vxland.Int(0)
	toIndex := fromIndex
	moreRoutes := false
	obj = &returnVxlanVtepStateGetInfo

	var v *vxlan.VtepDbEntry
	currIndex := vxland.Int(0)
	indexListAdded := make([]int, 0)
	for currIndex = vxland.Int(0); validCount != count && vxlan.GetVtepDbListEntry(int32(currIndex), &v); currIndex++ {

		if currIndex < fromIndex {
			continue
		} else {

			foundIndex := false
			for _, idx := range indexListAdded {
				if idx == int(currIndex) {
					foundIndex = true
				}
			}

			if foundIndex {
				continue
			}

			nextVxlanVtepState = &vxlanVtepStateList[validCount]
			nextVxlanVtepState.Intf = v.VtepName
			nextVxlanVtepState.IntfRef = v.SrcIfName
			nextVxlanVtepState.IfIndex = v.VtepIfIndex
			nextVxlanVtepState.Vni = int32(v.Vni)
			nextVxlanVtepState.DstUDP = int16(v.UDP)
			nextVxlanVtepState.TTL = int16(v.TTL)
			nextVxlanVtepState.TOS = int16(v.TOS)
			if v.FilterUnknownCustVlan {
				nextVxlanVtepState.InnerVlanHandlingMode = 0
			} else {
				nextVxlanVtepState.InnerVlanHandlingMode = 1
			}
			nextVxlanVtepState.SrcIp = v.SrcIp.String()
			nextVxlanVtepState.VlanId = int16(v.VlanId)
			nextVxlanVtepState.Mtu = int32(v.MTU)

			var v2 *vxlan.VtepDbEntry
			for currIndex2 := vxland.Int(0); vxlan.GetVtepDbListEntry(int32(currIndex2), &v2); currIndex2++ {

				if v2 != nil &&
					v2.VtepConfigName == v.VtepConfigName &&
					v2.Vni == v.Vni {
					foundIndex = false
					for _, idx := range indexListAdded {
						if idx == int(currIndex2) {
							foundIndex = true
						}
					}

					if !foundIndex {
						indexListAdded = append(indexListAdded, int(currIndex2))

					} else {
						continue
					}

					OperState := "UNKNOWN"
					if v2.Enable && v2.VxlanVtepMachineFsm.Machine.Curr.CurrentState() == vxlan.VxlanVtepStateStart {
						OperState = "UP"
					} else {
						OperState = "DOWN"
					}

					nextVxlanVtepState.PerDstIpState = append(nextVxlanVtepState.PerDstIpState, &vxland.VtepStateEntry{
						SubVtepName:        v2.VtepName,
						DstIp:              v2.DstIp.String(),
						RxSwPkts:           int64(v2.GetStats().Rxpkts),
						RxSwBytes:          int64(v2.GetStats().Rxbytes),
						RxSwDropPkts:       int64(v2.GetStats().Rxdroppkts),
						RxSwDropBytes:      int64(v2.GetStats().Rxdropbytes),
						RxSwFwdPkts:        int64(v2.GetStats().Rxfwdpkts),
						RxSwFwdBytes:       int64(v2.GetStats().Rxfwdbytes),
						TxSwPkts:           int64(v2.GetStats().Txpkts),
						TxSwBytes:          int64(v2.GetStats().Txbytes),
						TxSwDropPkts:       int64(v2.GetStats().Txdroppkts),
						TxSwDropBytes:      int64(v2.GetStats().Txdropbytes),
						TxSwFwdPkts:        int64(v2.GetStats().Txfwdpkts),
						TxSwFwdBytes:       int64(v2.GetStats().Txfwdbytes),
						LastSwRxDropReason: v2.GetStats().Lastrxdropreason,
						LastSwTxDropReason: v2.GetStats().Lasttxdropreason,
						//vis.RxFwdPkts             uint64 `DESCRIPTION: Rx Forwaded Packets`
						//vis.RxDropPkts            uint64 `DESCRIPTION: Rx Dropped Packets`
						//vis.RxUnknownVni          uint64 `DESCRIPTION: Rx Unknown Vni in frame`
						VtepFsmState:     vxlan.VxlanVtepStateStrMap[v2.VxlanVtepMachineFsm.Machine.Curr.CurrentState()],
						VtepFsmPrevState: vxlan.VxlanVtepStateStrMap[v2.VxlanVtepMachineFsm.Machine.Curr.PreviousState()],
						OperState:        OperState,
					})
				}
			}
			if len(returnVxlanVtepStates) == 0 {
				returnVxlanVtepStates = make([]*vxland.VxlanVtepInstanceState, 0)
			}
			returnVxlanVtepStates = append(returnVxlanVtepStates, nextVxlanVtepState)
			validCount++
			toIndex++
		}
	}
	// lets try and get the next agg if one exists then there are more routes
	if v != nil {
		moreRoutes = vxlan.GetVtepDbListEntry(int32(currIndex), &v)
	}
	obj.VxlanVtepInstanceStateList = returnVxlanVtepStates
	obj.StartIdx = fromIndex
	obj.EndIdx = toIndex + 1
	obj.More = moreRoutes
	obj.Count = validCount

	return obj, nil
}

func (v *VXLANDServiceHandler) HandleDbReadVxlanGlobal(dbHdl redis.Conn) error {
	if dbHdl != nil {
		var dbObj objects.VxlanGlobal
		objList, err := dbObj.GetAllObjFromDb(dbHdl)
		if err != nil {
			v.logger.Warning("DB Query failed when retrieving VxlanInstance objects")
			return err
		}
		for idx := 0; idx < len(objList); idx++ {
			obj := vxland.NewVxlanGlobal()
			dbObject := objList[idx].(objects.VxlanGlobal)
			objects.ConvertvxlandVxlanGlobalObjToThrift(&dbObject, obj)
			v.CreateVxlanGlobal(obj)
		}
	}
	return nil
}

func (v *VXLANDServiceHandler) HandleDbReadVxlanInstance(dbHdl redis.Conn, disable bool) error {
	if dbHdl != nil {
		var dbObj objects.VxlanInstance
		objList, err := dbObj.GetAllObjFromDb(dbHdl)
		if err != nil {
			v.logger.Warning("DB Query failed when retrieving VxlanInstance objects")
			return err
		}
		for idx := 0; idx < len(objList); idx++ {
			obj := vxland.NewVxlanInstance()
			dbObject := objList[idx].(objects.VxlanInstance)
			objects.ConvertvxlandVxlanInstanceObjToThrift(&dbObject, obj)

			if disable {
				_, err = v.DeleteVxlanInstance(obj)
			} else {
				_, err = v.CreateVxlanInstance(obj)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *VXLANDServiceHandler) HandleDbReadVxlanVtepInstance(dbHdl redis.Conn, disable bool) error {
	if dbHdl != nil {
		var dbObj objects.VxlanVtepInstance
		objList, err := dbObj.GetAllObjFromDb(dbHdl)
		if err != nil {
			v.logger.Warning("DB Query failed when retrieving VxlanVtepInstance objects")
			return err
		}
		for idx := 0; idx < len(objList); idx++ {
			obj := vxland.NewVxlanVtepInstance()
			dbObject := objList[idx].(objects.VxlanVtepInstance)
			objects.ConvertvxlandVxlanVtepInstanceObjToThrift(&dbObject, obj)
			if disable {
				_, err = v.DeleteVxlanVtepInstance(obj)
			} else {
				_, err = v.CreateVxlanVtepInstance(obj)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *VXLANDServiceHandler) ReadConfigFromDB(prevState int) error {

	dbHdl := dbutils.NewDBUtil(v.logger)
	err := dbHdl.Connect()
	if err != nil {
		v.logger.Err("Unable to connect to db")
		return err
	}
	defer dbHdl.Disconnect()

	if prevState == vxlan.VXLAN_GLOBAL_INIT {

		if err := v.HandleDbReadVxlanGlobal(dbHdl); err != nil {
			fmt.Println("Error getting All LacpGlobal objects")
			return err
		}
	}
	currState := vxlan.VxlanGlobalStateGet()

	v.logger.Info(fmt.Sprintf("Global State prev %d curr %d", prevState, currState))

	if currState == vxlan.VXLAN_GLOBAL_DISABLE_PENDING &&
		prevState == vxlan.VXLAN_GLOBAL_ENABLE {

		// lets delete the Aggregator first
		if err := v.HandleDbReadVxlanVtepInstance(dbHdl, true); err != nil {
			v.logger.Err("Error getting All VxlanVtep objects")
			return err
		}

		if err := v.HandleDbReadVxlanInstance(dbHdl, true); err != nil {
			v.logger.Err("Error getting All Vxlan objects")
			return err
		}
		vxlan.VxlanGlobalStateSet(vxlan.VXLAN_GLOBAL_DISABLE)
	} else if prevState != vxlan.VXLAN_GLOBAL_ENABLE &&
		currState == vxlan.VXLAN_GLOBAL_ENABLE {

		if err := v.HandleDbReadVxlanInstance(dbHdl, false); err != nil {
			fmt.Println("Error getting All Vxlan objects")
			return err
		}

		if err := v.HandleDbReadVxlanVtepInstance(dbHdl, false); err != nil {
			fmt.Println("Error getting All VxlanVtep objects")
			return err
		}
	}

	return nil
}
