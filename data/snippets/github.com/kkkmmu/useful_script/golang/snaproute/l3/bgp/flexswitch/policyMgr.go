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

package FSMgr

import (
	"bytes"
	"encoding/json"
	_ "fmt"
	"l3/bgp/api"
	"l3/rib/ribdCommonDefs"
	"models/objects"
	"utils/logging"
	utilspolicy "utils/policy"

	nanomsg "github.com/op/go-nanomsg"
)

/*  Init policy manager with specific needs
 */
func NewFSPolicyMgr(logger *logging.Writer, fileName string) *FSPolicyMgr {

	mgr := &FSPolicyMgr{
		plugin: "ovsdb",
		logger: logger,
	}

	return mgr
}

/*  Start nano msg socket with ribd
 */
func (mgr *FSPolicyMgr) Start() {
	mgr.logger.Info("Starting policyMgr")
	mgr.policySubSocket, _ = mgr.setupSubSocket(ribdCommonDefs.PUB_SOCKET_POLICY_ADDR)
	go mgr.listenForPolicyUpdates(mgr.policySubSocket)
}

func (mgr *FSPolicyMgr) setupSubSocket(address string) (*nanomsg.SubSocket, error) {
	var err error
	var socket *nanomsg.SubSocket
	if socket, err = nanomsg.NewSubSocket(); err != nil {
		mgr.logger.Errf("Failed to create subscribe socket %s error:%s", address, err)
		return nil, err
	}

	if err = socket.Subscribe(""); err != nil {
		mgr.logger.Errf("Failed to subscribe to \"\" on subscribe socket %s, error:%s", address, err)
		return nil, err
	}

	if _, err = socket.Connect(address); err != nil {
		mgr.logger.Errf("Failed to connect to publisher socket %s, error:%s", address, err)
		return nil, err
	}

	mgr.logger.Infof("Connected to publisher socker %s", address)
	if err = socket.SetRecvBuffer(1024 * 1024); err != nil {
		mgr.logger.Err("Failed to set the buffer size for subscriber socket", address, "error:", err)
		return nil, err
	}
	return socket, nil
}

func convertModelsToPolicyConditionConfig(cfg *objects.PolicyCondition) *utilspolicy.PolicyConditionConfig {
	if cfg == nil {
		return nil
	}

	destIPMatch := utilspolicy.PolicyDstIpMatchPrefixSetCondition{
		Prefix: utilspolicy.PolicyPrefix{
			IpPrefix:        cfg.IpPrefix,
			MasklengthRange: cfg.MaskLengthRange,
		},
	}
	matchExtendedCommunityInfo := utilspolicy.PolicyExtendedCommunityInfo{cfg.ExtendedCommunityType,
		cfg.ExtendedCommunityValue}

	return &utilspolicy.PolicyConditionConfig{
		Name:                          cfg.Name,
		ConditionType:                 cfg.ConditionType,
		MatchDstIpPrefixConditionInfo: destIPMatch,
		MatchCommunityConditionInfo:   utilspolicy.PolicyMatchCommunitySetCondition{cfg.Community, cfg.CommunitySet},
		MatchASPathConditionInfo:      utilspolicy.PolicyMatchASPathSetCondition{cfg.ASPath, cfg.ASPathSet},
		MatchLocalPrefConditionInfo:   cfg.LocalPref,
		MatchMEDConditionInfo:         cfg.MED,
		MatchExtendedCommunityConditionInfo: utilspolicy.PolicyMatchExtendedCommunitySetCondition{
			ExtendedCommunity: matchExtendedCommunityInfo, ExtendedCommunitySet: cfg.ExtendedCommunitySet},
	}
}

func (mgr *FSPolicyMgr) handlePolicyConditionUpdates(msg ribdCommonDefs.RibdNotifyMsg) {
	policyCondition := objects.PolicyCondition{}
	var updateMsg string
	switch msg.MsgType {
	case ribdCommonDefs.NOTIFY_POLICY_CONDITION_CREATED:
		updateMsg = "Add"
	case ribdCommonDefs.NOTIFY_POLICY_CONDITION_DELETED:
		updateMsg = "Remove"
	case ribdCommonDefs.NOTIFY_POLICY_CONDITION_UPDATED:
		updateMsg = "Update"
	default:
		mgr.logger.Errf("Unknown policy condition notification type %d", msg.MsgType)
		return
	}

	err := json.Unmarshal(msg.MsgBuf, &policyCondition)
	if err != nil {
		mgr.logger.Errf("Unmarshal RIB policy condition update failed with err %s", err)
		return
	}

	mgr.logger.Info(updateMsg, "Policy Condition", policyCondition.Name, "type:", policyCondition.ConditionType)
	condition := convertModelsToPolicyConditionConfig(&policyCondition)
	if condition == nil {
		mgr.logger.Err(updateMsg, "Policy Condition", policyCondition.Name, "conversion failed")
		return
	}

	if msg.MsgType == ribdCommonDefs.NOTIFY_POLICY_CONDITION_CREATED {
		api.AddPolicyCondition(*condition)
	} else if msg.MsgType == ribdCommonDefs.NOTIFY_POLICY_CONDITION_DELETED {
		api.RemovePolicyCondition(condition.Name)
	} else if msg.MsgType == ribdCommonDefs.NOTIFY_POLICY_CONDITION_UPDATED {
		api.UpdatePolicyCondition(*condition)
	}
}

func convertModelsToPolicyStmtConfig(cfg *objects.PolicyStmt) *utilspolicy.PolicyStmtConfig {
	if cfg == nil {
		return nil
	}

	actions := make([]string, 1)
	actions[0] = cfg.Action

	setActions := make([]utilspolicy.PolicyActionCfg, 0)
	for _, setAction := range cfg.SetActions {
		setActions = append(setActions, utilspolicy.PolicyActionCfg{
			Attr:              setAction.Attr,
			Community:         setAction.Community,
			ExtendedCommunity: utilspolicy.PolicyExtendedCommunityInfo{setAction.ExtendedCommunityType, setAction.ExtendedCommunityValue},
			LocalPref:         setAction.LocalPref,
			MED:               setAction.MED,
			PrependASPath:     setAction.PrependASPath,
		})
	}
	return &utilspolicy.PolicyStmtConfig{
		Name:            cfg.Name,
		MatchConditions: cfg.MatchConditions,
		Conditions:      cfg.Conditions,
		Actions:         actions,
		SetActions:      setActions,
	}
}

func (mgr *FSPolicyMgr) handlePolicyStmtUpdates(msg ribdCommonDefs.RibdNotifyMsg) {
	policyStmt := objects.PolicyStmt{}
	var updateMsg string
	switch msg.MsgType {
	case ribdCommonDefs.NOTIFY_POLICY_STMT_CREATED:
		updateMsg = "Add"
	case ribdCommonDefs.NOTIFY_POLICY_STMT_DELETED:
		updateMsg = "Remove"
	case ribdCommonDefs.NOTIFY_POLICY_STMT_UPDATED:
		updateMsg = "Update"
	default:
		mgr.logger.Errf("Unknown policy statement notification type %d", msg.MsgType)
		return
	}

	err := json.Unmarshal(msg.MsgBuf, &policyStmt)
	if err != nil {
		mgr.logger.Errf("Unmarshal RIB policy condition update failed with err %s", err)
		return
	}

	mgr.logger.Info(updateMsg, "Policy statement", policyStmt.Name)
	stmt := convertModelsToPolicyStmtConfig(&policyStmt)
	if stmt == nil {
		mgr.logger.Err(updateMsg, "Policy statement", policyStmt.Name, "conversion failed")
		return
	}

	if msg.MsgType == ribdCommonDefs.NOTIFY_POLICY_STMT_CREATED {
		api.AddPolicyStmt(*stmt)
	} else if msg.MsgType == ribdCommonDefs.NOTIFY_POLICY_STMT_DELETED {
		api.RemovePolicyStmt(stmt.Name)
	} else if msg.MsgType == ribdCommonDefs.NOTIFY_POLICY_STMT_UPDATED {
		api.UpdatePolicyStmt(*stmt)
	}
}

func convertModelsToPolicyDefintionConfig(cfg *objects.PolicyDefinition) *utilspolicy.PolicyDefinitionConfig {
	if cfg == nil {
		return nil
	}
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

func (mgr *FSPolicyMgr) handlePolicyDefinitionUpdates(msg ribdCommonDefs.RibdNotifyMsg) {
	policyDefinition := objects.PolicyDefinition{}
	var updateMsg string
	switch msg.MsgType {
	case ribdCommonDefs.NOTIFY_POLICY_DEFINITION_CREATED:
		updateMsg = "Add"
	case ribdCommonDefs.NOTIFY_POLICY_DEFINITION_DELETED:
		updateMsg = "Remove"
	case ribdCommonDefs.NOTIFY_POLICY_DEFINITION_UPDATED:
		updateMsg = "Update"
	default:
		mgr.logger.Errf("Unknown policy definition notification type %d", msg.MsgType)
		return
	}

	err := json.Unmarshal(msg.MsgBuf, &policyDefinition)
	if err != nil {
		mgr.logger.Errf("Unmarshal RIB policy definition update failed with err %s", err)
		return
	}

	mgr.logger.Info(updateMsg, "Policy definition", policyDefinition.Name, " policy type:", policyDefinition.PolicyType)
	condition := convertModelsToPolicyDefintionConfig(&policyDefinition)
	if condition == nil {
		mgr.logger.Err(updateMsg, "Policy definition", policyDefinition.Name, "conversion failed")
		return
	}

	if msg.MsgType == ribdCommonDefs.NOTIFY_POLICY_DEFINITION_CREATED {
		api.AddPolicyDefinition(*condition)
	} else if msg.MsgType == ribdCommonDefs.NOTIFY_POLICY_DEFINITION_DELETED {
		api.RemovePolicyDefinition(condition.Name)
	} else if msg.MsgType == ribdCommonDefs.NOTIFY_POLICY_DEFINITION_UPDATED {
		api.UpdatePolicyDefinition(*condition)
	}
}

func (mgr *FSPolicyMgr) handlePolicyUpdates(rxBuf []byte) {
	reader := bytes.NewReader(rxBuf)
	decoder := json.NewDecoder(reader)
	msg := ribdCommonDefs.RibdNotifyMsg{}
	err := decoder.Decode(&msg)
	if err != nil {
		mgr.logger.Err("Error while decoding msg")
		return
	}
	switch msg.MsgType {
	case ribdCommonDefs.NOTIFY_POLICY_CONDITION_CREATED, ribdCommonDefs.NOTIFY_POLICY_CONDITION_DELETED,
		ribdCommonDefs.NOTIFY_POLICY_CONDITION_UPDATED:
		mgr.handlePolicyConditionUpdates(msg)
	case ribdCommonDefs.NOTIFY_POLICY_STMT_CREATED, ribdCommonDefs.NOTIFY_POLICY_STMT_DELETED,
		ribdCommonDefs.NOTIFY_POLICY_STMT_UPDATED:
		mgr.handlePolicyStmtUpdates(msg)
	case ribdCommonDefs.NOTIFY_POLICY_DEFINITION_CREATED, ribdCommonDefs.NOTIFY_POLICY_DEFINITION_DELETED,
		ribdCommonDefs.NOTIFY_POLICY_DEFINITION_UPDATED:
		mgr.handlePolicyDefinitionUpdates(msg)
	default:
		mgr.logger.Errf("**** Received Policy update with unknown type %d ****", msg.MsgType)
	}
}

func (mgr *FSPolicyMgr) listenForPolicyUpdates(socket *nanomsg.SubSocket) {
	for {
		mgr.logger.Info("Read on Policy subscriber socket...")
		rxBuf, err := socket.Recv(0)
		if err != nil {
			mgr.logger.Err("Recv on Policy subscriber socket failed with error:", err)
			continue
		}
		mgr.logger.Info("Policy subscriber recv returned:", rxBuf)
		mgr.handlePolicyUpdates(rxBuf)
	}
}
