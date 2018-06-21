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

// packetUtils.go
package policy

import (
	"l3/bgp/packet"
	bgprib "l3/bgp/rib"
	"l3/bgp/utils"
	"strings"
	utilspolicy "utils/policy"
	"utils/policy/policyCommonDefs"
)

func ApplyActionsToPacket(pa []packet.BGPPathAttr, stmt utilspolicy.PolicyStmt) ([]packet.BGPPathAttr, bool) {
	medUpdated := false
	for _, action := range stmt.SetActionsState {
		utils.Logger.Infof("ApplyActionsToPacket - action:%+v", action)
		switch action.Attr {
		case policyCommonDefs.PolicyActionTypeSetLocalPref:
			pa = packet.SetLocalPrefToPathAttrs(pa, action.LocalPref, true)

		case policyCommonDefs.PolicyActionTypeSetCommunity:
			pa = packet.AddCommunityToPathAttrs(pa, action.Community)

		case policyCommonDefs.PolicyActionTypeSetExtendedCommunity:
			pa = packet.AddExtCommunityToPathAttrs(pa, action.ExtendedCommunity)

		case policyCommonDefs.PolicyActionTypeSetMED:
			pa = packet.SetMEDToPathAttrs(pa, action.MED, true)
			medUpdated = true

		case policyCommonDefs.PolicyActionTypeSetPrependASPath:
			asPathStrList := strings.Split(action.PrependASPath, " ")
			utils.Logger.Info("PrependASPath list =", asPathStrList)
			asList := make([]uint32, 0)
			for _, asStr := range asPathStrList {
				as := strings.TrimSpace(asStr)
				if len(as) > 0 {
					if asNum, err := utils.GetAsNum(as); err == nil {
						asList = append(asList, uint32(asNum))
					} else {
						utils.Logger.Err("Apply action for AS prepend failed, as path:", action.PrependASPath,
							"failed convert as num", as, "to integer with err:", err)
						asList = nil
						break
					}
				}
			}

			if asList != nil {
				pa = packet.PrependASList(pa, asList, 4)
			}
		}
	}
	return pa, medUpdated
}

func GetPolicyEngineFilterEntity(path *bgprib.Path) *utilspolicy.PolicyEngineFilterEntityParams {
	params := &utilspolicy.PolicyEngineFilterEntityParams{}

	community := packet.GetCommunityValues(path.PathAttrs)
	if community != nil {
		params.Community = community
	}

	extComm := packet.GetExtCommunityValues(path.PathAttrs)
	if extComm != nil {
		params.ExtendedCommunity = extComm
	}

	if asPathStr, err := packet.GetASPathAsString(path.PathAttrs); err == nil {
		params.ASPath = asPathStr
	}

	if localPref, err := packet.GetLocalPref(path.PathAttrs); err == nil {
		params.LocalPref = localPref
	}

	if med, ok := packet.GetMED(path.PathAttrs); ok == true {
		params.MED = med
	}

	return params
}
