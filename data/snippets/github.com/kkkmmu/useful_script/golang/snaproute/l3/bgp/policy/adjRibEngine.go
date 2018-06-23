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

// adjRibEngine.go
package policy

import (
	bgprib "l3/bgp/rib"
	"utils/logging"
	utilspolicy "utils/policy"
)

type AdjRibPolicyExtensions struct {
	HitCounter    int
	RouteList     []string
	RouteInfoList []*bgprib.AdjRIBPathIdRoute
}

type AdjRibPPolicyEngine struct {
	BasePolicyEngine
}

func NewAdjRibPolicyEngine(logger *logging.Writer) *AdjRibPPolicyEngine {
	policyEngine := &AdjRibPPolicyEngine{
		BasePolicyEngine: NewBasePolicyEngine(logger, utilspolicy.NewPolicyEngineDB(logger)),
	}
	policyEngine.SetGetPolicyEntityMapIndexFunc(getPolicyEnityKey)
	return policyEngine
}

func (eng *AdjRibPPolicyEngine) CreatePolicyDefinition(defCfg utilspolicy.PolicyDefinitionConfig) error {
	defCfg.Extensions = AdjRibPolicyExtensions{}
	return eng.PolicyEngine.CreatePolicyDefinition(defCfg)
}
