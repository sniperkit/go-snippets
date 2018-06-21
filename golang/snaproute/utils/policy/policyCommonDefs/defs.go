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

package policyCommonDefs

const (
	PolicyConditionTypeDstIpPrefixMatch       = 0
	PolicyConditionTypeProtocolMatch          = 1
	PolicyConditionTypeNeighborMatch          = 2
	PolicyConditionTypeCommunityMatch         = 3
	PolicyConditionTypeExtendedCommunityMatch = 4
	PolicyConditionTypeLocalPrefMatch         = 5
	PolicyConditionTypeASPathMatch            = 6
	PolicyConditionTypeMEDMatch               = 7
	PolicyActionTypeRouteDisposition          = 0
	PolicyActionTypeRouteRedistribute         = 1
	PoilcyActionTypeSetAdminDistance          = 2
	PolicyActionTypeNetworkStatementAdvertise = 3
	PolicyActionTypeAggregate                 = 4
	PolicyActionTypeRIBIn                     = 5
	PolicyActionTypeRIBOut                    = 6
	PolicyActionTypeSetCommunity              = 7
	PolicyActionTypeSetExtendedCommunity      = 8
	PolicyActionTypeSetLocalPref              = 9
	PolicyActionTypeSetPrependASPath          = 10
	PolicyActionTypeSetMED                    = 11
	PolicyPath_Import                         = 1
	PolicyPath_Export                         = 2
	PolicyPath_All                            = 3
)

var SetActionMap = map[string]int{
	"Community":         PolicyActionTypeSetCommunity,
	"ExtendedCommunity": PolicyActionTypeSetExtendedCommunity,
	"LocalPref":         PolicyActionTypeSetLocalPref,
	"MED":               PolicyActionTypeSetMED,
	"PrependASPath":     PolicyActionTypeSetPrependASPath,
}
