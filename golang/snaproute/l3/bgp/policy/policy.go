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
	_ "fmt"
	"utils/logging"
	utilspolicy "utils/policy"
)

type PolicyActionFunc struct {
	ApplyFunc utilspolicy.Policyfunc
	UndoFunc  utilspolicy.UndoActionfunc
}

type BGPPolicyEngine interface {
	CreatePolicyCondition(utilspolicy.PolicyConditionConfig) (bool, error)
	CreatePolicyStmt(utilspolicy.PolicyStmtConfig) error
	CreatePolicyDefinition(utilspolicy.PolicyDefinitionConfig) error
	CreatePolicyAction(utilspolicy.PolicyActionConfig) (bool, error)
	DeletePolicyCondition(string) (bool, error)
	DeletePolicyStmt(string) error
	DeletePolicyDefinition(string) error
	DeletePolicyAction(string) (bool, error)
	UpdateApplyPolicy(utilspolicy.ApplyPolicyInfo, bool)
	UpdateUndoApplyPolicy(utilspolicy.ApplyPolicyInfo, bool)
	SetTraverseFuncs(utilspolicy.EntityTraverseAndApplyPolicyfunc, utilspolicy.EntityTraverseAndReversePolicyfunc)
	SetActionFuncs(map[int]PolicyActionFunc)
	SetEntityUpdateFunc(utilspolicy.EntityUpdatefunc)
	SetIsEntityPresentFunc(utilspolicy.PolicyCheckfunc)
	SetGetPolicyEntityMapIndexFunc(utilspolicy.GetPolicyEnityMapIndexFunc)
	GetPolicyEngine() *utilspolicy.PolicyEngineDB
}

type BasePolicyEngine struct {
	logger       *logging.Writer
	PolicyEngine *utilspolicy.PolicyEngineDB
}

func NewBasePolicyEngine(logger *logging.Writer, policyEngine *utilspolicy.PolicyEngineDB) BasePolicyEngine {
	return BasePolicyEngine{
		logger:       logger,
		PolicyEngine: policyEngine,
	}
}

func (eng *BasePolicyEngine) SetTraverseFuncs(traverseApplyFunc utilspolicy.EntityTraverseAndApplyPolicyfunc,
	traverseReverseFunc utilspolicy.EntityTraverseAndReversePolicyfunc) {
	eng.logger.Info("BasePolicyEngine:SetTraverseFunc traverse apply func %v", traverseApplyFunc)
	if traverseApplyFunc != nil {
		eng.PolicyEngine.SetTraverseAndApplyPolicyFunc(traverseApplyFunc)
	}
	eng.logger.Info("BasePolicyEngine:SetTraverseFunc traverse reverse func %v", traverseReverseFunc)
	if traverseReverseFunc != nil {
		eng.PolicyEngine.SetTraverseAndReversePolicyFunc(traverseReverseFunc)
	}
}

func (eng *BasePolicyEngine) SetActionFuncs(actionFuncMap map[int]PolicyActionFunc) {
	eng.logger.Infof("BasePolicyEngine:SetApplyActionFunc actionFuncMap %v", actionFuncMap)
	for actionType, actionFuncs := range actionFuncMap {
		eng.logger.Info("BasePolicyEngine:SetApplyActionFunc set apply/undo callbacks for action", actionType)
		if actionFuncs.ApplyFunc != nil {
			eng.PolicyEngine.SetActionFunc(actionType, actionFuncs.ApplyFunc)
		}
		if actionFuncs.UndoFunc != nil {
			eng.PolicyEngine.SetUndoActionFunc(actionType, actionFuncs.UndoFunc)
		}
	}
}

func (eng *BasePolicyEngine) SetEntityUpdateFunc(entityUpdateFunc utilspolicy.EntityUpdatefunc) {
	eng.logger.Info("BasePolicyEngine:SetEntityUpdateFunc func %v", entityUpdateFunc)
	if entityUpdateFunc != nil {
		eng.PolicyEngine.SetEntityUpdateFunc(entityUpdateFunc)
	}
}

func (eng *BasePolicyEngine) SetIsEntityPresentFunc(entityPresentFunc utilspolicy.PolicyCheckfunc) {
	eng.logger.Info("BasePolicyEngine:SetIsEntityPresentFunc func %v", entityPresentFunc)
	if entityPresentFunc != nil {
		eng.PolicyEngine.SetIsEntityPresentFunc(entityPresentFunc)
	}
}

func (eng *BasePolicyEngine) SetGetPolicyEntityMapIndexFunc(policyEntityKeyFunc utilspolicy.GetPolicyEnityMapIndexFunc) {
	eng.logger.Info("BasePolicyEngine:SetGetPolicyEntityMapIndexFunc func %v", policyEntityKeyFunc)
	if policyEntityKeyFunc != nil {
		eng.PolicyEngine.SetGetPolicyEntityMapIndexFunc(policyEntityKeyFunc)
	}
}

func (eng *BasePolicyEngine) CreatePolicyCondition(condCfg utilspolicy.PolicyConditionConfig) (bool, error) {
	return eng.PolicyEngine.CreatePolicyCondition(condCfg)
}

func (eng *BasePolicyEngine) CreatePolicyStmt(stmtCfg utilspolicy.PolicyStmtConfig) error {
	return eng.PolicyEngine.CreatePolicyStatement(stmtCfg)
}

func (eng *BasePolicyEngine) CreatePolicyAction(actionCfg utilspolicy.PolicyActionConfig) (bool, error) {
	return eng.PolicyEngine.CreatePolicyAggregateAction(actionCfg)
}

func (eng *BasePolicyEngine) DeletePolicyCondition(conditionName string) (bool, error) {
	conditionCfg := utilspolicy.PolicyConditionConfig{Name: conditionName}
	return eng.PolicyEngine.DeletePolicyCondition(conditionCfg)
}

func (eng *BasePolicyEngine) DeletePolicyStmt(stmtName string) error {
	stmtCfg := utilspolicy.PolicyStmtConfig{Name: stmtName}
	return eng.PolicyEngine.DeletePolicyStatement(stmtCfg)
}

func (eng *BasePolicyEngine) DeletePolicyDefinition(policyName string) error {
	policyCfg := utilspolicy.PolicyDefinitionConfig{Name: policyName}
	return eng.PolicyEngine.DeletePolicyDefinition(policyCfg)
}

func (eng *BasePolicyEngine) DeletePolicyAction(actionName string) (bool, error) {
	actionCfg := utilspolicy.PolicyActionConfig{Name: actionName}
	return eng.PolicyEngine.DeletePolicyAction(actionCfg)
}

func (eng *BasePolicyEngine) UpdateApplyPolicy(info utilspolicy.ApplyPolicyInfo, apply bool) {
	eng.PolicyEngine.UpdateApplyPolicy(info, apply)
}

func (eng *BasePolicyEngine) UpdateUndoApplyPolicy(info utilspolicy.ApplyPolicyInfo, apply bool) {
	eng.PolicyEngine.UpdateUndoApplyPolicy(info, apply)
}

func (eng *BasePolicyEngine) GetPolicyEngine() *utilspolicy.PolicyEngineDB {
	return eng.PolicyEngine
}
