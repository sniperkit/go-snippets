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
	"asicdInt"
)

// vxlan related apis
func (svcHdlr AsicDaemonServiceHandler) CreateVxlan(config *asicdInt.Vxlan) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) DeleteVxlan(config *asicdInt.Vxlan) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) UpdateVxlan(vni int32, addvlanlist, delvlanlist, addUntaggedVlanlist, delUntaggedVlanlist []int16) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) CreateVxlanVtep(config *asicdInt.Vtep) (int32, error) {
	return 0, nil
}

func (svcHdlr AsicDaemonServiceHandler) DeleteVxlanVtep(config *asicdInt.Vtep) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) UpdateVxlanVtepAttr(vtepName string, vni int32, tos, ttl int16, mtu int32) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) LearnFdbVtep(mac string, name string, ifIndex int32) (rval bool, err error) {
	return true, nil
}
