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
	"l3/vrrp/common"
	"l3/vrrp/debug"
	"l3/vrrp/fsm"
)

type IPIntf interface {
	Init(*common.BaseIpInfo)
	Update(*common.BaseIpInfo)
	DeInit(*common.BaseIpInfo)
	GetObjFromDb(*common.BaseIpInfo)
	SetVrrpIntfKey(KeyInfo)
	GetVrrpIntfKey() *KeyInfo
	GetIntfRef() string
}

type VrrpInterface struct {
	L3     *common.BaseIpInfo // Vrrp Port Information Collected From System
	Config *common.IntfCfg    // Vrrp config for interface
	Fsm    *fsm.FSM           // Vrrp fsm information
}

func (intf *VrrpInterface) InitVrrpIntf(cfg *common.IntfCfg, l3Info *common.BaseIpInfo, vipCh chan *common.VirtualIpInfo, rxCh chan struct{}, txCh chan struct{}) {
	debug.Logger.Info("Initializing interface with config:", *cfg, "base ip interface:", *l3Info)
	intf.Config = cfg
	intf.L3 = l3Info
	// Init fsm
	intf.Fsm = fsm.InitFsm(intf.Config, l3Info, vipCh, rxCh, txCh)
}

func (intf *VrrpInterface) DeInitVrrpIntf() {
	intf.Fsm.DeInitFsm()
	intf.L3 = nil
	intf.Config = nil
	intf.Fsm = nil
}

func (intf *VrrpInterface) UpdateOperState(state string) {
	intf.L3.OperState = state
}

func (intf *VrrpInterface) StartFsm() {
	debug.Logger.Info("Starting Fsm for interface:", intf.L3.IntfRef, "vrid:", intf.Config.VRID)
	go intf.Fsm.StartFsm()
}

// should only be called if vrrp is disabled globally do not call this during de-init
func (intf *VrrpInterface) StopFsm() {
	intf.Fsm.IntfEventCh <- &fsm.IntfEvent{
		Event: fsm.TEAR_DOWN,
	}
}

func (intf *VrrpInterface) GetVMac() string {
	return intf.Fsm.VirtualMACAddress
}

func (intf *VrrpInterface) GetVirtualIpUpdateInfo() (string, string, string) {
	return intf.L3.IntfRef, intf.Config.VirtualIPAddr, intf.Fsm.VirtualMACAddress
}

func (intf *VrrpInterface) UpdateIpState() {
	if intf.Config.AdminState {
		if intf.Fsm.IsRunning() {
			// send out state up event
			debug.Logger.Info("fsm for interface:", intf.L3.IntfRef, "vrid:", intf.Config.VRID, "is running and hence sending state change")
			intf.Fsm.IntfEventCh <- &fsm.IntfEvent{
				Event:     fsm.STATE_CHANGE,
				OperState: intf.L3.OperState,
			}
		} else {
			intf.StartFsm()
		}
	}
}

func (intf *VrrpInterface) UpdateConfig(cfg *common.IntfCfg) (oper int) {
	debug.Logger.Info("Updating interface configuration from old Config:", *intf.Config, "to new Config:", *cfg)
	intf.Fsm.UpdateConfig(cfg)
	// Special case for handling interface admin state
	if intf.Config.AdminState != cfg.AdminState {
		if cfg.AdminState == false {
			// tear down fsm as admin state is down
			intf.StopFsm()
			oper = common.DELETE
		} else {
			// start fsm
			intf.StartFsm()
			oper = common.CREATE
		}
	}
	intf.Config = cfg
	return oper
}
