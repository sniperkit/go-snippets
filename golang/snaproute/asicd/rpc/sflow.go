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

// This file defines all Sflow interfaces provided Asicd Services

package rpc

import (
	"asicdInt"
)

func (svcHdlr AsicDaemonServiceHandler) EnableSFlowSampling(PortIfIdx int32, Enable bool) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) SetSFlowSamplingRate(PortIfIdx, Rate int32) (bool, error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkSFlowIntfInfo() (*asicdInt.SFlowIntfBulkInfo, error) {
	bulkObj := new(asicdInt.SFlowIntfBulkInfo)
	return bulkObj, nil
}
