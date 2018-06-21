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

// ribdPolicyServer.go
package server

import (
	"encoding/json"
	"github.com/op/go-nanomsg"
	defs "l3/rib/ribdCommonDefs"
	"models/objects"
	"ribd"
	"ribdInt"
)

type ApplyPolicyInfo struct {
	Source     string //source application/protocol
	Policy     string
	Action     string
	Conditions []*ribdInt.ConditionInfo
}

type ApplyPolicyList struct {
	ApplyList []*ribdInt.ApplyPolicyInfo
	UndoList  []*ribdInt.ApplyPolicyInfo
}

/*
   Function to send PolicyPrefixSet Notification
*/
func (ribdServiceHandler *RIBDServer) PolicyPrefixSetNotificationSend(PUB *nanomsg.PubSocket, cfg ribd.PolicyPrefixSet, evt int) {
	logger.Info("PolicyPrefixSetNotificationSend")
	msgBuf := objects.PolicyPrefixSet{}
	objects.ConvertThriftToribdPolicyPrefixSetObj(&cfg, &msgBuf)
	msgbufbytes, err := json.Marshal(msgBuf)
	msg := defs.RibdNotifyMsg{MsgType: uint16(evt), MsgBuf: msgbufbytes}
	buf, err := json.Marshal(msg)
	if err != nil {
		logger.Err("Error in marshalling Json")
		return
	}
	var evtStr string
	if evt == defs.NOTIFY_POLICY_PREFIX_SET_CREATED {
		evtStr = " POLICY_PREFIX_SET_CREATED "
	} else if evt == defs.NOTIFY_POLICY_PREFIX_SET_DELETED {
		evtStr = " POLICY_PREFIX_SET_DELETED "
	} else if evt == defs.NOTIFY_POLICY_PREFIX_SET_UPDATED {
		evtStr = " POLICY_PREFIX_SET_UPDATED "
	}
	eventInfo := evtStr + " for prefix set " + cfg.Name
	logger.Debug("Adding ", evtStr, " to notification channel")
	ribdServiceHandler.NotificationChannel <- NotificationMsg{PUB, buf, eventInfo}
}

/*
   Function to send PolicyCommunitySet Notification
*/
func (ribdServiceHandler *RIBDServer) PolicyCommunitySetNotificationSend(PUB *nanomsg.PubSocket, cfg ribd.PolicyCommunitySet, evt int) {
	logger.Info("PolicyCommunitySetNotificationSend")
	msgBuf := objects.PolicyCommunitySet{}
	objects.ConvertThriftToribdPolicyCommunitySetObj(&cfg, &msgBuf)
	msgbufbytes, err := json.Marshal(msgBuf)
	msg := defs.RibdNotifyMsg{MsgType: uint16(evt), MsgBuf: msgbufbytes}
	buf, err := json.Marshal(msg)
	if err != nil {
		logger.Err("Error in marshalling Json")
		return
	}
	var evtStr string
	if evt == defs.NOTIFY_POLICY_COMMUNITY_SET_CREATED {
		evtStr = " POLICY_COMMUNITY_SET_CREATED "
	} else if evt == defs.NOTIFY_POLICY_COMMUNITY_SET_DELETED {
		evtStr = " POLICY_COMMUNITY_SET_DELETED "
	} else if evt == defs.NOTIFY_POLICY_COMMUNITY_SET_UPDATED {
		evtStr = " POLICY_COMMUNITY_SET_UPDATED "
	}
	eventInfo := evtStr + " for community set " + cfg.Name
	logger.Debug("Adding ", evtStr, " to notification channel")
	ribdServiceHandler.NotificationChannel <- NotificationMsg{PUB, buf, eventInfo}
}

/*
   Function to send PolicyExtendedCommunitySet Notification
*/
func (ribdServiceHandler *RIBDServer) PolicyExtendedCommunitySetNotificationSend(PUB *nanomsg.PubSocket, cfg ribd.PolicyExtendedCommunitySet, evt int) {
	logger.Info("PolicyExtendedCommunitySetNotificationSend")
	msgBuf := objects.PolicyExtendedCommunitySet{}
	objects.ConvertThriftToribdPolicyExtendedCommunitySetObj(&cfg, &msgBuf)
	msgbufbytes, err := json.Marshal(msgBuf)
	msg := defs.RibdNotifyMsg{MsgType: uint16(evt), MsgBuf: msgbufbytes}
	buf, err := json.Marshal(msg)
	if err != nil {
		logger.Err("Error in marshalling Json")
		return
	}
	var evtStr string
	if evt == defs.NOTIFY_POLICY_EXTENDED_COMMUNITY_SET_CREATED {
		evtStr = " POLICY_EXTENDED_COMMUNITY_SET_CREATED "
	} else if evt == defs.NOTIFY_POLICY_EXTENDED_COMMUNITY_SET_DELETED {
		evtStr = " POLICY_EXTENDED_COMMUNITY_SET_DELETED "
	} else if evt == defs.NOTIFY_POLICY_EXTENDED_COMMUNITY_SET_UPDATED {
		evtStr = " POLICY_EXTENDED_COMMUNITY_SET_UPDATED "
	}
	eventInfo := evtStr + " for extended community set " + cfg.Name
	logger.Debug("Adding ", evtStr, " to notification channel")
	ribdServiceHandler.NotificationChannel <- NotificationMsg{PUB, buf, eventInfo}
}

/*
   Function to send PolicyASPathSet Notification
*/
func (ribdServiceHandler *RIBDServer) PolicyASPathSetNotificationSend(PUB *nanomsg.PubSocket, cfg ribd.PolicyASPathSet, evt int) {
	logger.Info("PolicyASPathSetNotificationSend")
	msgBuf := objects.PolicyASPathSet{}
	objects.ConvertThriftToribdPolicyASPathSetObj(&cfg, &msgBuf)
	msgbufbytes, err := json.Marshal(msgBuf)
	msg := defs.RibdNotifyMsg{MsgType: uint16(evt), MsgBuf: msgbufbytes}
	buf, err := json.Marshal(msg)
	if err != nil {
		logger.Err("Error in marshalling Json")
		return
	}
	var evtStr string
	if evt == defs.NOTIFY_POLICY_ASPATH_SET_CREATED {
		evtStr = " POLICY_ASPATH_SET_CREATED "
	} else if evt == defs.NOTIFY_POLICY_ASPATH_SET_DELETED {
		evtStr = " POLICY_ASPATH_SET_DELETED "
	} else if evt == defs.NOTIFY_POLICY_ASPATH_SET_UPDATED {
		evtStr = " POLICY_ASPATH_SET_UPDATED "
	}
	eventInfo := evtStr + " for aspath set " + cfg.Name
	logger.Debug("Adding ", evtStr, " to notification channel")
	ribdServiceHandler.NotificationChannel <- NotificationMsg{PUB, buf, eventInfo}
}

/*
   Function to send PolicyCondition Notification
*/
func (ribdServiceHandler *RIBDServer) PolicyConditionNotificationSend(PUB *nanomsg.PubSocket, cfg ribd.PolicyCondition, evt int) {
	logger.Info("PolicyConditionNotificationSend")
	msgBuf := objects.PolicyCondition{}
	objects.ConvertThriftToribdPolicyConditionObj(&cfg, &msgBuf)
	/*	msgBuf := objects.PolicyConditionConfig{
				Name : cfg.Name,
				ConditionType   :cfg.ConditionType,
				Protocol        : cfg.Protocol,
				IpPrefix        : cfg.IpPrefix,
				MasklengthRange : cfg.MasklengthRange
	}*/
	msgbufbytes, err := json.Marshal(msgBuf)
	msg := defs.RibdNotifyMsg{MsgType: uint16(evt), MsgBuf: msgbufbytes}
	buf, err := json.Marshal(msg)
	if err != nil {
		logger.Err("Error in marshalling Json")
		return
	}
	var evtStr string
	if evt == defs.NOTIFY_POLICY_CONDITION_CREATED {
		evtStr = " POLICY_CONDITION_CREATED "
	} else if evt == defs.NOTIFY_POLICY_CONDITION_DELETED {
		evtStr = " POLICY_CONDITION_DELETED "
	} else if evt == defs.NOTIFY_POLICY_CONDITION_UPDATED {
		evtStr = " POLICY_CONDITION_UPDATED "
	}
	eventInfo := evtStr + " for condition " + cfg.Name + " " + " type " + cfg.ConditionType
	logger.Debug("Adding ", evtStr, " to notification channel")
	ribdServiceHandler.NotificationChannel <- NotificationMsg{PUB, buf, eventInfo}
}

/*
   Function to send PolicyStmt Notification
*/
func (ribdServiceHandler *RIBDServer) PolicyStmtNotificationSend(PUB *nanomsg.PubSocket, cfg ribd.PolicyStmt, evt int) {
	logger.Info("PolicyStmtNotificationSend")
	msgBuf := objects.PolicyStmt{}
	objects.ConvertThriftToribdPolicyStmtObj(&cfg, &msgBuf)
	/*	msgBuf := objects.PolicyStmtConfig{
					Name               : cfg.Name,
					MatchConditions    : cfg.MatchConditions,
					Action             : cfg.Action
		}
		msgBuf.Conditions = make([]string,0)
		for i := 0;i<len(cfg.Conditions);i++ {
			msgBuf.Conditions = append(msgBuf.Conditions,cfg.Conditions[i])
		}*/
	msgbufbytes, err := json.Marshal(msgBuf)
	msg := defs.RibdNotifyMsg{MsgType: uint16(evt), MsgBuf: msgbufbytes}
	buf, err := json.Marshal(msg)
	if err != nil {
		logger.Err("Error in marshalling Json")
		return
	}
	var evtStr string
	if evt == defs.NOTIFY_POLICY_STMT_CREATED {
		evtStr = " POLICY_STMT_CREATED "
	} else if evt == defs.NOTIFY_POLICY_STMT_DELETED {
		evtStr = " POLICY_STMT_DELETED "
	} else if evt == defs.NOTIFY_POLICY_STMT_UPDATED {
		evtStr = " POLICY_STMT_UPDATED "
	}
	eventInfo := evtStr + " for policy stmt " + cfg.Name
	logger.Debug("Adding ", evtStr, " to notification channel")
	ribdServiceHandler.NotificationChannel <- NotificationMsg{PUB, buf, eventInfo}
}

/*
   Function to send PolicyDefinition Notification
*/
func (ribdServiceHandler *RIBDServer) PolicyDefinitionNotificationSend(PUB *nanomsg.PubSocket, cfg ribd.PolicyDefinition, evt int) {
	logger.Info("PolicyDefinitionNotificationSend")
	msgBuf := objects.PolicyDefinition{}
	objects.ConvertThriftToribdPolicyDefinitionObj(&cfg, &msgBuf)
	/*	msgBuf := objects.PolicyDefinitionConfig{
					Name        : cfg.Name,
					Priority    : cfg.Priority,
					MatchType             : cfg.MatchType,
					PolicyType  : cfg.PolicyType
		}
		msgBuf.PolicyDefinitionStatements = make([]policy.PolicyDefinitionStmtPrecedence, 0)
		var policyDefinitionStatement policy.PolicyDefinitionStmtPrecedence
		for i := 0; i < len(cfg.StatementList); i++ {
			policyDefinitionStatement.Precedence = int(cfg.StatementList[i].Priority)
			policyDefinitionStatement.Statement = cfg.StatementList[i].Statement
			msgBuf.PolicyDefinitionStatements = append(msgBuf.PolicyDefinitionStatements, policyDefinitionStatement)
		}*/
	msgbufbytes, err := json.Marshal(msgBuf)
	msg := defs.RibdNotifyMsg{MsgType: uint16(evt), MsgBuf: msgbufbytes}
	buf, err := json.Marshal(msg)
	if err != nil {
		logger.Err("Error in marshalling Json")
		return
	}
	var evtStr string
	if evt == defs.NOTIFY_POLICY_DEFINITION_CREATED {
		evtStr = " POLICY_DEFINITION_CREATED "
	} else if evt == defs.NOTIFY_POLICY_DEFINITION_DELETED {
		evtStr = " POLICY_DEFINITION_DELETED "
	} else if evt == defs.NOTIFY_POLICY_DEFINITION_UPDATED {
		evtStr = "POLICY_DEFINITION_UPDATED"
	}
	eventInfo := evtStr + " for policy " + cfg.Name
	logger.Debug("Adding ", evtStr, " to notification channel")
	ribdServiceHandler.NotificationChannel <- NotificationMsg{PUB, buf, eventInfo}
}

/*
   Handles all policy object config based server updates. The flow is:
                                               policyChannels
   userConfig------rpc(policyHandler_functions)----------------policyServer
                                                                 |
                                                              policy objects updated in RIB's GlobalpolicyEngine
													      send events to applications about these object configs
														 policy objects updated in RIB's local policy engine which
														 functions as filter
*/
func (ribdServiceHandler *RIBDServer) StartPolicyServer() {
	logger.Debug("Starting the policy server loop")
	for {
		select {
		case conf := <-ribdServiceHandler.PolicyConfCh:
			logger.Debug("received message on PolicyConfCh channel, op: ", conf.Op, " policyconfdone:", ribdServiceHandler.PolicyConfDone)
			var err error
			if conf.Op == defs.AddPolicyCondition {
				_, err = ribdServiceHandler.ProcessPolicyConditionConfigCreate(conf.OrigConfigObject.(*ribd.PolicyCondition), ribdServiceHandler.GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyConditionNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyCondition)), defs.NOTIFY_POLICY_CONDITION_CREATED)
					_, err = ribdServiceHandler.ProcessPolicyConditionConfigCreate(conf.OrigConfigObject.(*ribd.PolicyCondition), ribdServiceHandler.PolicyEngineDB)
				}
			} else if conf.Op == defs.DelPolicyCondition {
				_, err = ribdServiceHandler.ProcessPolicyConditionConfigDelete(conf.OrigConfigObject.(*ribd.PolicyCondition), ribdServiceHandler.GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyConditionNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyCondition)), defs.NOTIFY_POLICY_CONDITION_DELETED)
					_, err = ribdServiceHandler.ProcessPolicyConditionConfigDelete(conf.OrigConfigObject.(*ribd.PolicyCondition), ribdServiceHandler.PolicyEngineDB)
				}
			} else if conf.Op == defs.UpdatePolicyCondition {
				logger.Debug("Received updatePolicyCondition on policy server channel")
				var err error
				if conf.PatchOp == nil || len(conf.PatchOp) == 0 {
					err = ribdServiceHandler.ProcessPolicyConditionConfigUpdate(conf.OrigConfigObject.(*ribd.PolicyCondition), conf.NewConfigObject.(*ribd.PolicyCondition), conf.AttrSet, GlobalPolicyEngineDB)
					if err == nil {
						ribdServiceHandler.PolicyConditionNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyCondition)), defs.NOTIFY_POLICY_CONDITION_UPDATED)
						ribdServiceHandler.ProcessPolicyConditionConfigUpdate(conf.OrigConfigObject.(*ribd.PolicyCondition), conf.NewConfigObject.(*ribd.PolicyCondition), conf.AttrSet, ribdServiceHandler.PolicyEngineDB)
					}
				}
			} else if conf.Op == defs.AddPolicyPrefixSet {
				_, err = ribdServiceHandler.ProcessPolicyPrefixSetConfigCreate(conf.OrigConfigObject.(*ribd.PolicyPrefixSet), GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyPrefixSetNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyPrefixSet)), defs.NOTIFY_POLICY_PREFIX_SET_CREATED)
					_, err = ribdServiceHandler.ProcessPolicyPrefixSetConfigCreate(conf.OrigConfigObject.(*ribd.PolicyPrefixSet), ribdServiceHandler.PolicyEngineDB)
				}
			} else if conf.Op == defs.DelPolicyPrefixSet {
				_, err = ribdServiceHandler.ProcessPolicyPrefixSetConfigDelete(conf.OrigConfigObject.(*ribd.PolicyPrefixSet), GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyPrefixSetNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyPrefixSet)), defs.NOTIFY_POLICY_PREFIX_SET_DELETED)
					_, err = ribdServiceHandler.ProcessPolicyPrefixSetConfigDelete(conf.OrigConfigObject.(*ribd.PolicyPrefixSet), ribdServiceHandler.PolicyEngineDB)
				}
			} else if conf.Op == defs.UpdatePolicyPrefixSet {
				logger.Debug("Received updatePolicyPrefixSet on policy server channel")
				var err error
				if conf.PatchOp == nil || len(conf.PatchOp) == 0 {
					err = ribdServiceHandler.ProcessPolicyPrefixSetConfigUpdate(conf.OrigConfigObject.(*ribd.PolicyPrefixSet), conf.NewConfigObject.(*ribd.PolicyPrefixSet), conf.AttrSet, GlobalPolicyEngineDB)
					if err == nil {
						ribdServiceHandler.PolicyPrefixSetNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyPrefixSet)), defs.NOTIFY_POLICY_PREFIX_SET_UPDATED)
						ribdServiceHandler.ProcessPolicyPrefixSetConfigUpdate(conf.OrigConfigObject.(*ribd.PolicyPrefixSet), conf.NewConfigObject.(*ribd.PolicyPrefixSet), conf.AttrSet, ribdServiceHandler.PolicyEngineDB)
					}
				} else {
					err = ribdServiceHandler.ProcessPolicyPrefixSetConfigPatchUpdate(conf.OrigConfigObject.(*ribd.PolicyPrefixSet), conf.NewConfigObject.(*ribd.PolicyPrefixSet), conf.PatchOp, GlobalPolicyEngineDB)
					if err == nil {
						ribdServiceHandler.PolicyPrefixSetNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyPrefixSet)), defs.NOTIFY_POLICY_PREFIX_SET_UPDATED)
						ribdServiceHandler.ProcessPolicyPrefixSetConfigPatchUpdate(conf.OrigConfigObject.(*ribd.PolicyPrefixSet), conf.NewConfigObject.(*ribd.PolicyPrefixSet), conf.PatchOp, ribdServiceHandler.PolicyEngineDB)
					}
				}
			} else if conf.Op == defs.AddPolicyCommunitySet {
				_, err = ribdServiceHandler.ProcessPolicyCommunitySetConfigCreate(conf.OrigConfigObject.(*ribd.PolicyCommunitySet), GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyCommunitySetNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyCommunitySet)), defs.NOTIFY_POLICY_COMMUNITY_SET_CREATED)
					_, err = ribdServiceHandler.ProcessPolicyCommunitySetConfigCreate(conf.OrigConfigObject.(*ribd.PolicyCommunitySet), ribdServiceHandler.PolicyEngineDB)
				}
			} else if conf.Op == defs.DelPolicyCommunitySet {
				_, err = ribdServiceHandler.ProcessPolicyCommunitySetConfigDelete(conf.OrigConfigObject.(*ribd.PolicyCommunitySet), GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyCommunitySetNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyCommunitySet)), defs.NOTIFY_POLICY_COMMUNITY_SET_DELETED)
					_, err = ribdServiceHandler.ProcessPolicyCommunitySetConfigDelete(conf.OrigConfigObject.(*ribd.PolicyCommunitySet), ribdServiceHandler.PolicyEngineDB)
				}
			} else if conf.Op == defs.UpdatePolicyCommunitySet {
				logger.Debug("Received updatePolicyCommunitySet on policy server channel")
				var err error
				if conf.PatchOp == nil || len(conf.PatchOp) == 0 {
					err = ribdServiceHandler.ProcessPolicyCommunitySetConfigUpdate(conf.OrigConfigObject.(*ribd.PolicyCommunitySet), conf.NewConfigObject.(*ribd.PolicyCommunitySet), conf.AttrSet, GlobalPolicyEngineDB)
					if err == nil {
						ribdServiceHandler.PolicyCommunitySetNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyCommunitySet)), defs.NOTIFY_POLICY_COMMUNITY_SET_UPDATED)
						ribdServiceHandler.ProcessPolicyCommunitySetConfigUpdate(conf.OrigConfigObject.(*ribd.PolicyCommunitySet), conf.NewConfigObject.(*ribd.PolicyCommunitySet), conf.AttrSet, ribdServiceHandler.PolicyEngineDB)
					}
				} else {
					err = ribdServiceHandler.ProcessPolicyCommunitySetConfigPatchUpdate(conf.OrigConfigObject.(*ribd.PolicyCommunitySet), conf.NewConfigObject.(*ribd.PolicyCommunitySet), conf.PatchOp, GlobalPolicyEngineDB)
					if err == nil {
						ribdServiceHandler.PolicyCommunitySetNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyCommunitySet)), defs.NOTIFY_POLICY_COMMUNITY_SET_UPDATED)
						ribdServiceHandler.ProcessPolicyCommunitySetConfigPatchUpdate(conf.OrigConfigObject.(*ribd.PolicyCommunitySet), conf.NewConfigObject.(*ribd.PolicyCommunitySet), conf.PatchOp, ribdServiceHandler.PolicyEngineDB)
					}
				}
			} else if conf.Op == defs.AddPolicyExtendedCommunitySet {
				_, err = ribdServiceHandler.ProcessPolicyExtendedCommunitySetConfigCreate(conf.OrigConfigObject.(*ribd.PolicyExtendedCommunitySet), GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyExtendedCommunitySetNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyExtendedCommunitySet)), defs.NOTIFY_POLICY_EXTENDED_COMMUNITY_SET_CREATED)
					_, err = ribdServiceHandler.ProcessPolicyExtendedCommunitySetConfigCreate(conf.OrigConfigObject.(*ribd.PolicyExtendedCommunitySet), ribdServiceHandler.PolicyEngineDB)
				}
			} else if conf.Op == defs.DelPolicyExtendedCommunitySet {
				_, err = ribdServiceHandler.ProcessPolicyExtendedCommunitySetConfigDelete(conf.OrigConfigObject.(*ribd.PolicyExtendedCommunitySet), GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyExtendedCommunitySetNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyExtendedCommunitySet)), defs.NOTIFY_POLICY_EXTENDED_COMMUNITY_SET_DELETED)
					_, err = ribdServiceHandler.ProcessPolicyExtendedCommunitySetConfigDelete(conf.OrigConfigObject.(*ribd.PolicyExtendedCommunitySet), ribdServiceHandler.PolicyEngineDB)
				}
			} else if conf.Op == defs.UpdatePolicyExtendedCommunitySet {
				logger.Debug("Received updatePolicyExtendedCommunitySet on policy server channel")
				var err error
				if conf.PatchOp == nil || len(conf.PatchOp) == 0 {
					err = ribdServiceHandler.ProcessPolicyExtendedCommunitySetConfigUpdate(conf.OrigConfigObject.(*ribd.PolicyExtendedCommunitySet), conf.NewConfigObject.(*ribd.PolicyExtendedCommunitySet), conf.AttrSet, GlobalPolicyEngineDB)
					if err == nil {
						ribdServiceHandler.PolicyExtendedCommunitySetNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyExtendedCommunitySet)), defs.NOTIFY_POLICY_EXTENDED_COMMUNITY_SET_UPDATED)
						ribdServiceHandler.ProcessPolicyExtendedCommunitySetConfigUpdate(conf.OrigConfigObject.(*ribd.PolicyExtendedCommunitySet), conf.NewConfigObject.(*ribd.PolicyExtendedCommunitySet), conf.AttrSet, ribdServiceHandler.PolicyEngineDB)
					}
				} else {
					err = ribdServiceHandler.ProcessPolicyExtendedCommunitySetConfigPatchUpdate(conf.OrigConfigObject.(*ribd.PolicyExtendedCommunitySet), conf.NewConfigObject.(*ribd.PolicyExtendedCommunitySet), conf.PatchOp, GlobalPolicyEngineDB)
					if err == nil {
						ribdServiceHandler.PolicyExtendedCommunitySetNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyExtendedCommunitySet)), defs.NOTIFY_POLICY_EXTENDED_COMMUNITY_SET_UPDATED)
						ribdServiceHandler.ProcessPolicyExtendedCommunitySetConfigPatchUpdate(conf.OrigConfigObject.(*ribd.PolicyExtendedCommunitySet), conf.NewConfigObject.(*ribd.PolicyExtendedCommunitySet), conf.PatchOp, ribdServiceHandler.PolicyEngineDB)
					}
				}
			} else if conf.Op == defs.AddPolicyASPathSet {
				_, err = ribdServiceHandler.ProcessPolicyASPathSetConfigCreate(conf.OrigConfigObject.(*ribd.PolicyASPathSet), GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyASPathSetNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyASPathSet)), defs.NOTIFY_POLICY_ASPATH_SET_CREATED)
					_, err = ribdServiceHandler.ProcessPolicyASPathSetConfigCreate(conf.OrigConfigObject.(*ribd.PolicyASPathSet), ribdServiceHandler.PolicyEngineDB)
				}
			} else if conf.Op == defs.DelPolicyASPathSet {
				_, err = ribdServiceHandler.ProcessPolicyASPathSetConfigDelete(conf.OrigConfigObject.(*ribd.PolicyASPathSet), GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyASPathSetNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyASPathSet)), defs.NOTIFY_POLICY_ASPATH_SET_DELETED)
					_, err = ribdServiceHandler.ProcessPolicyASPathSetConfigDelete(conf.OrigConfigObject.(*ribd.PolicyASPathSet), ribdServiceHandler.PolicyEngineDB)
				}
			} else if conf.Op == defs.UpdatePolicyASPathSet {
				logger.Debug("Received updatePolicyASPathSet on policy server channel")
				var err error
				if conf.PatchOp == nil || len(conf.PatchOp) == 0 {
					err = ribdServiceHandler.ProcessPolicyASPathSetConfigUpdate(conf.OrigConfigObject.(*ribd.PolicyASPathSet), conf.NewConfigObject.(*ribd.PolicyASPathSet), conf.AttrSet, GlobalPolicyEngineDB)
					if err == nil {
						ribdServiceHandler.PolicyASPathSetNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyASPathSet)), defs.NOTIFY_POLICY_ASPATH_SET_UPDATED)
						ribdServiceHandler.ProcessPolicyASPathSetConfigUpdate(conf.OrigConfigObject.(*ribd.PolicyASPathSet), conf.NewConfigObject.(*ribd.PolicyASPathSet), conf.AttrSet, ribdServiceHandler.PolicyEngineDB)
					}
				} else {
					err = ribdServiceHandler.ProcessPolicyASPathSetConfigPatchUpdate(conf.OrigConfigObject.(*ribd.PolicyASPathSet), conf.NewConfigObject.(*ribd.PolicyASPathSet), conf.PatchOp, GlobalPolicyEngineDB)
					if err == nil {
						ribdServiceHandler.PolicyASPathSetNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyASPathSet)), defs.NOTIFY_POLICY_ASPATH_SET_UPDATED)
						ribdServiceHandler.ProcessPolicyASPathSetConfigPatchUpdate(conf.OrigConfigObject.(*ribd.PolicyASPathSet), conf.NewConfigObject.(*ribd.PolicyASPathSet), conf.PatchOp, ribdServiceHandler.PolicyEngineDB)
					}
				}
			} else if conf.Op == defs.AddPolicyStmt {
				err = ribdServiceHandler.ProcessPolicyStmtConfigCreate(conf.OrigConfigObject.(*ribd.PolicyStmt), GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyStmtNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyStmt)), defs.NOTIFY_POLICY_STMT_CREATED)
					err = ribdServiceHandler.ProcessPolicyStmtConfigCreate(conf.OrigConfigObject.(*ribd.PolicyStmt), ribdServiceHandler.PolicyEngineDB)
				}
			} else if conf.Op == defs.DelPolicyStmt {
				err = ribdServiceHandler.ProcessPolicyStmtConfigDelete(conf.OrigConfigObject.(*ribd.PolicyStmt), GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyStmtNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyStmt)), defs.NOTIFY_POLICY_STMT_DELETED)
					err = ribdServiceHandler.ProcessPolicyStmtConfigDelete(conf.OrigConfigObject.(*ribd.PolicyStmt), ribdServiceHandler.PolicyEngineDB)
				}
			} else if conf.Op == defs.UpdatePolicyStmt {
				logger.Debug("Received updatePolicyStmt on policy server channel")
				var err error
				if conf.PatchOp == nil || len(conf.PatchOp) == 0 {
					err = ribdServiceHandler.ProcessPolicyStmtConfigUpdate(conf.OrigConfigObject.(*ribd.PolicyStmt), conf.NewConfigObject.(*ribd.PolicyStmt), conf.AttrSet, GlobalPolicyEngineDB)
					if err == nil {
						ribdServiceHandler.PolicyStmtNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyStmt)), defs.NOTIFY_POLICY_STMT_UPDATED)
						ribdServiceHandler.ProcessPolicyStmtConfigUpdate(conf.OrigConfigObject.(*ribd.PolicyStmt), conf.NewConfigObject.(*ribd.PolicyStmt), conf.AttrSet, ribdServiceHandler.PolicyEngineDB)
					}
				} else {
					err = ribdServiceHandler.ProcessPolicyStmtConfigPatchUpdate(conf.OrigConfigObject.(*ribd.PolicyStmt), conf.NewConfigObject.(*ribd.PolicyStmt), conf.PatchOp, GlobalPolicyEngineDB)
					if err == nil {
						ribdServiceHandler.PolicyStmtNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyStmt)), defs.NOTIFY_POLICY_STMT_UPDATED)
						ribdServiceHandler.ProcessPolicyStmtConfigPatchUpdate(conf.OrigConfigObject.(*ribd.PolicyStmt), conf.NewConfigObject.(*ribd.PolicyStmt), conf.PatchOp, ribdServiceHandler.PolicyEngineDB)
					}
				}
			} else if conf.Op == defs.AddPolicyDefinition {
				err = ribdServiceHandler.ProcessPolicyDefinitionConfigCreate(conf.OrigConfigObject.(*ribd.PolicyDefinition), GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyDefinitionNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyDefinition)), defs.NOTIFY_POLICY_DEFINITION_CREATED)
					err = ribdServiceHandler.ProcessPolicyDefinitionConfigCreate(conf.OrigConfigObject.(*ribd.PolicyDefinition), ribdServiceHandler.PolicyEngineDB)
				}
			} else if conf.Op == defs.DelPolicyDefinition {
				err = ribdServiceHandler.ProcessPolicyDefinitionConfigDelete(conf.OrigConfigObject.(*ribd.PolicyDefinition), GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyDefinitionNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyDefinition)), defs.NOTIFY_POLICY_DEFINITION_DELETED)
					err = ribdServiceHandler.ProcessPolicyDefinitionConfigDelete(conf.OrigConfigObject.(*ribd.PolicyDefinition), ribdServiceHandler.PolicyEngineDB)
				}
			} else if conf.Op == defs.UpdatePolicyDefinition {
				logger.Debug("Received updatePolicyDefinition on policy server channel")
				var err error
				if conf.PatchOp == nil || len(conf.PatchOp) == 0 {
					err = ribdServiceHandler.ProcessPolicyDefinitionConfigUpdate(conf.OrigConfigObject.(*ribd.PolicyDefinition), conf.NewConfigObject.(*ribd.PolicyDefinition), conf.AttrSet, GlobalPolicyEngineDB)
					if err == nil {
						ribdServiceHandler.PolicyDefinitionNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyDefinition)), defs.NOTIFY_POLICY_DEFINITION_UPDATED)
						ribdServiceHandler.ProcessPolicyDefinitionConfigUpdate(conf.OrigConfigObject.(*ribd.PolicyDefinition), conf.NewConfigObject.(*ribd.PolicyDefinition), conf.AttrSet, ribdServiceHandler.PolicyEngineDB)
					}
				} else {
					err = ribdServiceHandler.ProcessPolicyDefinitionConfigPatchUpdate(conf.OrigConfigObject.(*ribd.PolicyDefinition), conf.NewConfigObject.(*ribd.PolicyDefinition), conf.PatchOp, GlobalPolicyEngineDB)
					if err == nil {
						ribdServiceHandler.PolicyDefinitionNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyDefinition)), defs.NOTIFY_POLICY_DEFINITION_UPDATED)
						ribdServiceHandler.ProcessPolicyDefinitionConfigPatchUpdate(conf.OrigConfigObject.(*ribd.PolicyDefinition), conf.NewConfigObject.(*ribd.PolicyDefinition), conf.PatchOp, ribdServiceHandler.PolicyEngineDB)
					}
				}
			} else if conf.Op == defs.ApplyPolicy {
				ribdServiceHandler.UpdateApplyPolicyList(conf.PolicyList.ApplyList, conf.PolicyList.UndoList, true, ribdServiceHandler.PolicyEngineDB)
				ribdServiceHandler.UpdateApplyPolicyList(conf.PolicyList.ApplyList, conf.PolicyList.UndoList, false, GlobalPolicyEngineDB)
			}
			ribdServiceHandler.PolicyConfDone <- err
		case info := <-ribdServiceHandler.PolicyUpdateApplyCh:
			/*
			   This channel update/processing occurs when an application applies a policy.
			*/
			logger.Debug("received message on PolicyUpdateApplyCh channel")
			//update the global policyEngineDB
			ribdServiceHandler.UpdateApplyPolicyList(info.ApplyList, info.UndoList, false, GlobalPolicyEngineDB)
		}
	}
}

/*
func (ribdServiceHandler *RIBDServer) StartPolicyServer() {
	logger.Debug("Starting the policy server loop")
	for {
		select {
		case condConf := <-ribdServiceHandler.PolicyConditionConfCh:
			logger.Debug("received message on PolicyConditionConfCh channel, op: ", condConf.Op)
			if condConf.Op == "add" {
				_, err := ribdServiceHandler.ProcessPolicyConditionConfigCreate(condConf.OrigConfigObject.(*ribd.PolicyCondition), ribdServiceHandler.GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyConditionNotificationSend(RIBD_POLICY_PUB, *(condConf.OrigConfigObject.(*ribd.PolicyCondition)), defs.NOTIFY_POLICY_CONDITION_CREATED)
					ribdServiceHandler.ProcessPolicyConditionConfigCreate(condConf.OrigConfigObject.(*ribd.PolicyCondition), ribdServiceHandler.PolicyEngineDB)
				}
			} else if condConf.Op == "del" {
				_, err := ribdServiceHandler.ProcessPolicyConditionConfigDelete(condConf.OrigConfigObject.(*ribd.PolicyCondition), ribdServiceHandler.GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyConditionNotificationSend(RIBD_POLICY_PUB, *(condConf.OrigConfigObject.(*ribd.PolicyCondition)), defs.NOTIFY_POLICY_CONDITION_DELETED)
					ribdServiceHandler.ProcessPolicyConditionConfigDelete(condConf.OrigConfigObject.(*ribd.PolicyCondition), ribdServiceHandler.PolicyEngineDB)
				}
			}
		case stmtConf := <-ribdServiceHandler.PolicyStmtConfCh:
			logger.Debug("received message on PolicyStmtConfCh channel, op: ", stmtConf.Op)
			if stmtConf.Op == "add" {
				err := ribdServiceHandler.ProcessPolicyStmtConfigCreate(stmtConf.OrigConfigObject.(*ribd.PolicyStmt), GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyStmtNotificationSend(RIBD_POLICY_PUB, *(stmtConf.OrigConfigObject.(*ribd.PolicyStmt)), defs.NOTIFY_POLICY_STMT_CREATED)
					ribdServiceHandler.ProcessPolicyStmtConfigCreate(stmtConf.OrigConfigObject.(*ribd.PolicyStmt), ribdServiceHandler.PolicyEngineDB)
				}
			} else if stmtConf.Op == "del" {
				err := ribdServiceHandler.ProcessPolicyStmtConfigDelete(stmtConf.OrigConfigObject.(*ribd.PolicyStmt), GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyStmtNotificationSend(RIBD_POLICY_PUB, *(stmtConf.OrigConfigObject.(*ribd.PolicyStmt)), defs.NOTIFY_POLICY_STMT_DELETED)
					ribdServiceHandler.ProcessPolicyStmtConfigDelete(stmtConf.OrigConfigObject.(*ribd.PolicyStmt), ribdServiceHandler.PolicyEngineDB)
				}
			}
		case conf := <-ribdServiceHandler.PolicyDefinitionConfCh:
			logger.Debug("received message on PolicyDefinitionConfCh channel, op:", conf.Op)
			if conf.Op == "add" {
				err := ribdServiceHandler.ProcessPolicyDefinitionConfigCreate(conf.OrigConfigObject.(*ribd.PolicyDefinition), GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyDefinitionNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyDefinition)), defs.NOTIFY_POLICY_DEFINITION_CREATED)
					ribdServiceHandler.ProcessPolicyDefinitionConfigCreate(conf.OrigConfigObject.(*ribd.PolicyDefinition), ribdServiceHandler.PolicyEngineDB)
				}
			} else if conf.Op == "del" {
				err := ribdServiceHandler.ProcessPolicyDefinitionConfigDelete(conf.OrigConfigObject.(*ribd.PolicyDefinition), GlobalPolicyEngineDB)
				if err == nil {
					ribdServiceHandler.PolicyDefinitionNotificationSend(RIBD_POLICY_PUB, *(conf.OrigConfigObject.(*ribd.PolicyDefinition)), defs.NOTIFY_POLICY_DEFINITION_DELETED)
					ribdServiceHandler.ProcessPolicyDefinitionConfigDelete(conf.OrigConfigObject.(*ribd.PolicyDefinition), ribdServiceHandler.PolicyEngineDB)
				}
			}
		case info := <-ribdServiceHandler.PolicyUpdateApplyCh:
			/*
			   This channel update/processing occurs when an application applies a policy.
*/
/*			logger.Debug("received message on PolicyUpdateApplyCh channel")
			//update the global policyEngineDB
			ribdServiceHandler.UpdateApplyPolicyList(info.ApplyList, info.UndoList, false, GlobalPolicyEngineDB)
		}
	}
}
*/
