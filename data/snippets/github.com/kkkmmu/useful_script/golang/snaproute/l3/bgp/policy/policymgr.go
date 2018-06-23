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
	"l3/bgp/config"
	"models/objects"
	"utils/dbutils"
	"utils/logging"
	utilspolicy "utils/policy"
)

var PolicyManager *BGPPolicyManager

type BGPPolicyManager struct {
	logger          *logging.Writer
	policyEngines   []BGPPolicyEngine
	ConditionCfgCh  chan utilspolicy.PolicyConditionConfig
	ActionCfgCh     chan utilspolicy.PolicyActionConfig
	StmtCfgCh       chan utilspolicy.PolicyStmtConfig
	DefinitionCfgCh chan utilspolicy.PolicyDefinitionConfig
	ConditionDelCh  chan string
	ActionDelCh     chan string
	StmtDelCh       chan string
	DefinitionDelCh chan string
	policyPlugin    config.PolicyMgrIntf
}

func NewPolicyManager(logger *logging.Writer, pMgr config.PolicyMgrIntf) *BGPPolicyManager {
	if PolicyManager == nil {
		policyManager := &BGPPolicyManager{}
		policyManager.logger = logger
		policyManager.policyEngines = make([]BGPPolicyEngine, 0)
		policyManager.ConditionCfgCh = make(chan utilspolicy.PolicyConditionConfig)
		policyManager.ActionCfgCh = make(chan utilspolicy.PolicyActionConfig)
		policyManager.StmtCfgCh = make(chan utilspolicy.PolicyStmtConfig)
		policyManager.DefinitionCfgCh = make(chan utilspolicy.PolicyDefinitionConfig)
		policyManager.ConditionDelCh = make(chan string)
		policyManager.ActionDelCh = make(chan string)
		policyManager.StmtDelCh = make(chan string)
		policyManager.DefinitionDelCh = make(chan string)
		policyManager.policyPlugin = pMgr
		PolicyManager = policyManager
	}

	return PolicyManager
}

func (eng *BGPPolicyManager) AddPolicyEngine(bgpPE BGPPolicyEngine) {
	eng.policyEngines = append(eng.policyEngines, bgpPE)
}

func convertModelsToPolicyCondition(cfg objects.PolicyCondition) *utilspolicy.PolicyConditionConfig {
	destIPMatch := utilspolicy.PolicyDstIpMatchPrefixSetCondition{
		Prefix: utilspolicy.PolicyPrefix{
			IpPrefix:        cfg.IpPrefix,
			MasklengthRange: cfg.MaskLengthRange,
		},
	}

	return &utilspolicy.PolicyConditionConfig{
		Name:                          cfg.Name,
		ConditionType:                 cfg.ConditionType,
		MatchDstIpPrefixConditionInfo: destIPMatch,
	}
}

func (eng *BGPPolicyManager) readPolicyConditions(dbUtil *dbutils.DBUtil) error {
	eng.logger.Info("readPolicyConditions")
	var conditionObj objects.PolicyCondition
	conditionList, err := dbUtil.GetAllObjFromDb(conditionObj)
	if err != nil {
		eng.logger.Err("readPolicyConditions - GetAllObjFromDb for policy condition failed with error", err)
		return err
	}

	for idx := 0; idx < len(conditionList); idx++ {
		policyCondCfg := convertModelsToPolicyCondition(conditionList[idx].(objects.PolicyCondition))
		eng.logger.Info("readPolicyConditions - create policy condition", policyCondCfg.Name)
		for _, pe := range eng.policyEngines {
			pe.CreatePolicyCondition(*policyCondCfg)
		}
	}
	return nil
}

func convertModelsToPolicyStmt(cfg objects.PolicyStmt) *utilspolicy.PolicyStmtConfig {
	actions := make([]string, 1)
	actions[0] = cfg.Action

	return &utilspolicy.PolicyStmtConfig{
		Name:            cfg.Name,
		MatchConditions: cfg.MatchConditions,
		Conditions:      cfg.Conditions,
		Actions:         actions,
	}
}

func (eng *BGPPolicyManager) readPolicyStmts(dbUtil *dbutils.DBUtil) error {
	eng.logger.Info("readPolicyStmts")
	var stmtObj objects.PolicyStmt
	stmtList, err := dbUtil.GetAllObjFromDb(stmtObj)
	if err != nil {
		eng.logger.Err("readPolicyStmts - GetAllObjFromDb for policy statement failed with error", err)
		return err
	}

	for idx := 0; idx < len(stmtList); idx++ {
		policyStmtCfg := convertModelsToPolicyStmt(stmtList[idx].(objects.PolicyStmt))
		eng.logger.Info("readPolicyStmts - create policy statement", policyStmtCfg.Name)
		for _, pe := range eng.policyEngines {
			pe.CreatePolicyStmt(*policyStmtCfg)
		}
	}
	return nil
}

func convertModelsToPolicyDefinition(cfg objects.PolicyDefinition) *utilspolicy.PolicyDefinitionConfig {
	stmtPrecedenceList := make([]utilspolicy.PolicyDefinitionStmtPrecedence, 0)
	for i := 0; i < len(cfg.StatementList); i++ {
		stmtPrecedence := utilspolicy.PolicyDefinitionStmtPrecedence{
			Precedence: int(cfg.StatementList[i].Priority),
			Statement:  cfg.StatementList[i].Statement,
		}
		stmtPrecedenceList = append(stmtPrecedenceList, stmtPrecedence)
	}

	return &utilspolicy.PolicyDefinitionConfig{
		Name:                       cfg.Name,
		Precedence:                 int(cfg.Priority),
		MatchType:                  cfg.MatchType,
		PolicyDefinitionStatements: stmtPrecedenceList,
		PolicyType:                 cfg.PolicyType,
	}
}

func (eng *BGPPolicyManager) readPolicyDefinitions(dbUtil *dbutils.DBUtil) error {
	eng.logger.Info("readPolicyDefinitions")
	var defObj objects.PolicyDefinition
	definitionList, err := dbUtil.GetAllObjFromDb(defObj)
	if err != nil {
		eng.logger.Err("readPolicyDefinitions - GetAllObjFromDb for policy definition failed with error", err)
		return err
	}

	for idx := 0; idx < len(definitionList); idx++ {
		policyDefCfg := convertModelsToPolicyDefinition(definitionList[idx].(objects.PolicyDefinition))
		eng.logger.Info("readPolicyDefinitions - create policy definition", policyDefCfg.Name)
		for _, pe := range eng.policyEngines {
			pe.CreatePolicyDefinition(*policyDefCfg)
		}
	}
	return nil
}

func (eng *BGPPolicyManager) StartPolicyEngine(dbUtil *dbutils.DBUtil, doneCh chan bool) {
	eng.policyPlugin.Start()
	eng.readPolicyConditions(dbUtil)
	eng.readPolicyStmts(dbUtil)
	eng.readPolicyDefinitions(dbUtil)
	doneCh <- true
	for {
		select {
		case condCfg := <-eng.ConditionCfgCh:
			eng.logger.Info("BGPPolicyEngine - create policy condition", condCfg.Name)
			for _, pe := range eng.policyEngines {
				pe.CreatePolicyCondition(condCfg)
			}

		case actionCfg := <-eng.ActionCfgCh:
			eng.logger.Info("BGPPolicyEngine - create policy action", actionCfg.Name)
			for _, pe := range eng.policyEngines {
				pe.CreatePolicyAction(actionCfg)
			}

		case stmtCfg := <-eng.StmtCfgCh:
			eng.logger.Info("BGPPolicyEngine - create policy statement", stmtCfg.Name)
			for _, pe := range eng.policyEngines {
				pe.CreatePolicyStmt(stmtCfg)
			}

		case defCfg := <-eng.DefinitionCfgCh:
			eng.logger.Info("BGPPolicyEngine - create policy definition", defCfg.Name)
			for _, pe := range eng.policyEngines {
				pe.CreatePolicyDefinition(defCfg)
			}

		case conditionName := <-eng.ConditionDelCh:
			eng.logger.Info("BGPPolicyEngine - delete policy condition", conditionName)
			for _, pe := range eng.policyEngines {
				pe.DeletePolicyCondition(conditionName)
			}

		case actionName := <-eng.ActionDelCh:
			eng.logger.Info("BGPPolicyEngine - delete policy action", actionName)
			for _, pe := range eng.policyEngines {
				pe.DeletePolicyAction(actionName)
			}

		case stmtName := <-eng.StmtDelCh:
			eng.logger.Info("BGPPolicyEngine - delete policy statment", stmtName)
			for _, pe := range eng.policyEngines {
				pe.DeletePolicyStmt(stmtName)
			}

		case policyName := <-eng.DefinitionDelCh:
			eng.logger.Info("BGPPolicyEngine - delete policy definition", policyName)
			for _, pe := range eng.policyEngines {
				pe.DeletePolicyDefinition(policyName)
			}
		}
	}
}
