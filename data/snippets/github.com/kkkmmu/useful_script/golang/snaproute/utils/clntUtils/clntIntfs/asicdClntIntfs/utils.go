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

package asicdClntIntfs

import (
	"models/objects"
	"utils/clntUtils/clntDefs/asicdClntDefs"
)

func ConvertFromClntDefsToObjectPort(obj *asicdClntDefs.Port, retObj *objects.Port) {
	retObj.IntfRef = string(obj.IntfRef)
	retObj.IfIndex = int32(obj.IfIndex)
	retObj.Description = string(obj.Description)
	retObj.PhyIntfType = string(obj.PhyIntfType)
	retObj.AdminState = string(obj.AdminState)
	retObj.MacAddr = string(obj.MacAddr)
	retObj.Speed = int32(obj.Speed)
	retObj.Duplex = string(obj.Duplex)
	retObj.Autoneg = string(obj.Autoneg)
	retObj.MediaType = string(obj.MediaType)
	retObj.Mtu = int32(obj.Mtu)
	retObj.BreakOutMode = string(obj.BreakOutMode)
	retObj.LoopbackMode = string(obj.LoopbackMode)
	retObj.EnableFEC = bool(obj.EnableFEC)
	retObj.PRBSTxEnable = bool(obj.PRBSTxEnable)
	retObj.PRBSRxEnable = bool(obj.PRBSRxEnable)
	retObj.PRBSPolynomial = string(obj.PRBSPolynomial)
}
