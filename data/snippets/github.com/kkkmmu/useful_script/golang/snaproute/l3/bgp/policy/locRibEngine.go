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

// policy.go
package policy

import (
	bgprib "l3/bgp/rib"
	"utils/logging"
	utilspolicy "utils/policy"
)

type PolicyExtensions struct {
	HitCounter    int
	RouteList     []string
	RouteInfoList []*bgprib.Route
}

type LocRibPolicyEngine struct {
	BasePolicyEngine
	maxId uint32
	ids   []uint32
}

func NewLocRibPolicyEngine(logger *logging.Writer) *LocRibPolicyEngine {
	policyEngine := &LocRibPolicyEngine{
		BasePolicyEngine: NewBasePolicyEngine(logger, utilspolicy.NewPolicyEngineDB(logger)),
	}
	policyEngine.SetGetPolicyEntityMapIndexFunc(getPolicyEnityKey)
	return policyEngine
}

func (l *LocRibPolicyEngine) CreatePolicyDefinition(defCfg utilspolicy.PolicyDefinitionConfig) error {
	defCfg.Extensions = PolicyExtensions{}
	return l.PolicyEngine.CreatePolicyDefinition(defCfg)
}

func (l *LocRibPolicyEngine) GetNextId() uint32 {
	var id uint32
	if len(l.ids) > 0 {
		id = l.ids[len(l.ids)-1]
		l.ids = l.ids[:len(l.ids)-1]
		return id
	}

	id = l.maxId
	l.maxId++
	return id
}

func (l *LocRibPolicyEngine) ReleaseId(id uint32) {
	if id+1 == l.maxId {
		l.maxId--
		return
	}

	l.ids = append(l.ids, id)
}
