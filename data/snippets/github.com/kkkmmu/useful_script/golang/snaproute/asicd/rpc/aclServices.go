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

package rpc

import (
	"asicdServices"
)

func (svcHdlr AsicDaemonServiceHandler) CreateAclGlobal(dropObj *asicdServices.AclGlobal) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) UpdateAclGlobal(dropObjOld *asicdServices.AclGlobal, dropObjNew *asicdServices.AclGlobal, attr []bool, patch []*asicdServices.PatchOpInfo) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) DeleteAclGlobal(dropObj *asicdServices.AclGlobal) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) CreateAcl(aclObj *asicdServices.Acl) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) DeleteAcl(aclObj *asicdServices.Acl) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) UpdateAcl(oldAclObj, newAclObj *asicdServices.Acl, attrset []bool, op []*asicdServices.PatchOpInfo) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) CreateAclIpv4Filter(Ipv4Filter *asicdServices.AclIpv4Filter) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) UpdateAclIpv4Filter(oldIpv4Filter, newIpv4Filter *asicdServices.AclIpv4Filter, attrset []bool, op []*asicdServices.PatchOpInfo) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) DeleteAclIpv4Filter(Ipv4Filter *asicdServices.AclIpv4Filter) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) CreateAclMacFilter(macFilter *asicdServices.AclMacFilter) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) UpdateAclMacFilter(oldMacFilter, newMacFilter *asicdServices.AclMacFilter, attrset []bool, op []*asicdServices.PatchOpInfo) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) DeleteAclMacFilter(macFilter *asicdServices.AclMacFilter) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) CreateAclIpv6Filter(Ipv6Filter *asicdServices.AclIpv6Filter) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) UpdateAclIpv6Filter(oldIpv6Filter, newIpv6Filter *asicdServices.AclIpv6Filter, attrset []bool, op []*asicdServices.PatchOpInfo) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) DeleteAclIpv6Filter(Ipv6Filter *asicdServices.AclIpv6Filter) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) CreateCopp(copp *asicdServices.Copp) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) UpdateCopp(oldCopp *asicdServices.Copp, newCopp *asicdServices.Copp, attrset []bool, op []*asicdServices.PatchOpInfo) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) DeleteCopp(copp *asicdServices.Copp) (bool, error) {
	return true, nil
}

/*** get bulk apis **********/
func (svcHdlr AsicDaemonServiceHandler) GetCopp(proto string) (*asicdServices.Copp, error) {
	coppCfg := asicdServices.NewCopp()
	return coppCfg, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkCopp(currMarker, count asicdServices.Int) (*asicdServices.CoppGetInfo, error) {
	bulkObj := asicdServices.NewCoppGetInfo()
	return bulkObj, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetAclState(aclName string) (*asicdServices.AclState, error) {
	aclState := asicdServices.NewAclState()
	return aclState, nil

}

func (svcHdlr AsicDaemonServiceHandler) GetBulkAclState(currMarker, count asicdServices.Int) (*asicdServices.AclStateGetInfo, error) {
	bulkObj := asicdServices.NewAclStateGetInfo()
	return bulkObj, nil

}

func (svcHdlr AsicDaemonServiceHandler) GetBulkCoppStatState(currMarker, count asicdServices.Int) (*asicdServices.CoppStatStateGetInfo, error) {
	bulkObj := asicdServices.NewCoppStatStateGetInfo()
	return bulkObj, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetCoppStatState(proto string) (*asicdServices.CoppStatState, error) {
	coppState := asicdServices.NewCoppStatState()
	return coppState, nil
}
