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

package api

import (
	bgppolicy "l3/bgp/policy"
	"sync"
	utilspolicy "utils/policy"
)

type PolicyApiLayer struct {
	policyManager *bgppolicy.BGPPolicyManager
}

var bgppolicyapi *PolicyApiLayer = nil
var policyOnce sync.Once

/*  Singleton instance should be accesible only within api
 */
func getPolicyInstance() *PolicyApiLayer {
	policyOnce.Do(func() {
		bgppolicyapi = &PolicyApiLayer{}
	})
	return bgppolicyapi
}

/*  Initialize bgp api layer with the channels that will be used for communicating
 *  with the policy engine server
 */
func InitPolicy(policyEngine *bgppolicy.BGPPolicyManager) {
	bgppolicyapi = getPolicyInstance()
	bgppolicyapi.policyManager = policyEngine
}

func AddPolicyCondition(condition utilspolicy.PolicyConditionConfig) {
	bgppolicyapi.policyManager.ConditionCfgCh <- condition
}

func RemovePolicyCondition(conditionName string) {
	bgppolicyapi.policyManager.ConditionDelCh <- conditionName
}

func UpdatePolicyCondition(condition utilspolicy.PolicyConditionConfig) {
	return
}

func AddPolicyStmt(stmt utilspolicy.PolicyStmtConfig) {
	bgppolicyapi.policyManager.StmtCfgCh <- stmt
}

func RemovePolicyStmt(stmtName string) {
	bgppolicyapi.policyManager.StmtDelCh <- stmtName
}

func UpdatePolicyStmt(stmt utilspolicy.PolicyStmtConfig) {
	return
}

func AddPolicyDefinition(definition utilspolicy.PolicyDefinitionConfig) {
	bgppolicyapi.policyManager.DefinitionCfgCh <- definition
}

func RemovePolicyDefinition(definitionName string) {
	bgppolicyapi.policyManager.DefinitionDelCh <- definitionName
}

func UpdatePolicyDefinition(definition utilspolicy.PolicyDefinitionConfig) {
	return
}
