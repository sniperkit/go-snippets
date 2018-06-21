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

package server

import (
	"fmt"
	defs "l3/rib/ribdCommonDefs"
	"ribd"
	"ribdInt"
	"testing"
	"time"
	"utils/policy"
)

var conditionsList []*ribd.PolicyCondition
var stmtsList []*ribd.PolicyStmt
var policyDefinitionsList []*ribd.PolicyDefinition
var applyPolicyList []*ribdInt.ApplyPolicyInfo

func InitConditionsList() {
	conditionsList = make([]*ribd.PolicyCondition, 0)
	conditionsList = append(conditionsList, &ribd.PolicyCondition{
		Name:          "MatchAll",
		ConditionType: "all",
	})
	conditionsList = append(conditionsList, &ribd.PolicyCondition{
		Name:          "MatchConnected",
		ConditionType: "MatchProtocol",
	})
	conditionsList = append(conditionsList, &ribd.PolicyCondition{
		Name:            "Match11.1.10Network",
		ConditionType:   "MatchDstIpPrefix",
		IpPrefix:        "11.1.10.0/24",
		MaskLengthRange: "exact",
	})
	conditionsList = append(conditionsList, &ribd.PolicyCondition{
		Name:            "Match40.1.10Network",
		ConditionType:   "MatchDstIpPrefix",
		IpPrefix:        "40.1.10.0/24",
		MaskLengthRange: "exact",
	})
	conditionsList = append(conditionsList, &ribd.PolicyCondition{
		Name:            "Match50.1Network",
		ConditionType:   "MatchDstIpPrefix",
		IpPrefix:        "50.1.0.0/16",
		MaskLengthRange: "16-24",
	})
	conditionsList = append(conditionsList, &ribd.PolicyCondition{
		Name:            "Match40.1.10Network",
		ConditionType:   "MatchDstIpPrefix",
		IpPrefix:        "40.1.10.0/24",
		MaskLengthRange: "exact",
	})
}
func InitStmtsList() {
	stmtsList = make([]*ribd.PolicyStmt, 0)
	stmtsList = append(stmtsList, &ribd.PolicyStmt{
		Name:            "redistConnectStmt",
		MatchConditions: "all",
		Conditions:      []string{"MatchConnected"},
	})
	stmtsList = append(stmtsList, &ribd.PolicyStmt{
		Name:            "redistStaticStmt",
		MatchConditions: "all",
		Conditions:      []string{"MatchStatic"},
	})
	stmtsList = append(stmtsList, &ribd.PolicyStmt{
		Name:            "redistConnected11.1.10NetworkStmt",
		MatchConditions: "all",
		Conditions:      []string{"MatchConnected", "Match11.1.10Network"},
	})
	stmtsList = append(stmtsList, &ribd.PolicyStmt{
		Name:            "redistNetworkStmt",
		MatchConditions: "any",
		Conditions:      []string{"Match50.1Network", "Match40.1.10Network"},
	})
}
func InitPolicyDefinitionsList() {
	policyDefinitionsList = make([]*ribd.PolicyDefinition, 0)
	policyDefinitionsList = append(policyDefinitionsList, &ribd.PolicyDefinition{
		Name:       "redistConnect",
		Priority:   1,
		MatchType:  "all",
		PolicyType: "BGP",
		StatementList: []*ribd.PolicyDefinitionStmtPriority{
			&ribd.PolicyDefinitionStmtPriority{Priority: 1, Statement: "redistConnectStmt"},
		},
	})
	policyDefinitionsList = append(policyDefinitionsList, &ribd.PolicyDefinition{
		Name:       "redistConnectAndStatic1",
		Priority:   1,
		MatchType:  "all",
		PolicyType: "BGP",
		StatementList: []*ribd.PolicyDefinitionStmtPriority{
			&ribd.PolicyDefinitionStmtPriority{Priority: 1, Statement: "redistConnectStmt"},
			&ribd.PolicyDefinitionStmtPriority{Priority: 1, Statement: "redistStaticStmt"},
		},
	})
	policyDefinitionsList = append(policyDefinitionsList, &ribd.PolicyDefinition{
		Name:       "redistConnectAndStatic2",
		Priority:   1,
		MatchType:  "all",
		PolicyType: "BGP",
		StatementList: []*ribd.PolicyDefinitionStmtPriority{
			&ribd.PolicyDefinitionStmtPriority{Priority: 1, Statement: "redistStaticStmt"},
			&ribd.PolicyDefinitionStmtPriority{Priority: 1, Statement: "redistConnectStmt"},
		},
	})
	policyDefinitionsList = append(policyDefinitionsList, &ribd.PolicyDefinition{
		Name:       "redistConnectAndStatic3",
		Priority:   1,
		MatchType:  "all",
		PolicyType: "BGP",
		StatementList: []*ribd.PolicyDefinitionStmtPriority{
			&ribd.PolicyDefinitionStmtPriority{Priority: 1, Statement: "redistStaticStmt"},
			&ribd.PolicyDefinitionStmtPriority{Priority: 2, Statement: "redistConnectStmt"},
		},
	})
	policyDefinitionsList = append(policyDefinitionsList, &ribd.PolicyDefinition{
		Name:       "redistConnectAndNetwork",
		Priority:   1,
		MatchType:  "all",
		PolicyType: "BGP",
		StatementList: []*ribd.PolicyDefinitionStmtPriority{
			&ribd.PolicyDefinitionStmtPriority{Priority: 1, Statement: "redistNetwork"},
			&ribd.PolicyDefinitionStmtPriority{Priority: 2, Statement: "redistConnectStmt"},
		},
	})
}
func InitApplyPolicyInfo() {
	applyPolicyList = make([]*ribdInt.ApplyPolicyInfo, 0)
	applyPolicyList = append(applyPolicyList, &ribdInt.ApplyPolicyInfo{
		Source:     "BGP",
		Policy:     "redistConnect",
		Action:     "Redistribution",
		Conditions: []*ribdInt.ConditionInfo{},
	})
	applyPolicyList = append(applyPolicyList, &ribdInt.ApplyPolicyInfo{
		Source: "BGP",
		Policy: "redistNetworkStmt",
		Action: "Redistribution",
		Conditions: []*ribdInt.ConditionInfo{
			&ribdInt.ConditionInfo{
				ConditionType: "MatchProtocol",
				Protocol:      "Connected",
			},
		},
	})

}
func PolicyConditionConfigCreate(config *ribd.PolicyCondition) {
	_, err := server.ProcessPolicyConditionConfigCreate(config, server.GlobalPolicyEngineDB)
	fmt.Println("err:", err, " for condition:", config)
}
func PolicyStmtConfigCreate(config *ribd.PolicyStmt) {
	err := server.ProcessPolicyStmtConfigCreate(config, server.GlobalPolicyEngineDB)
	fmt.Println("err:", err, " for stmt:", config)
}
func PolicyDefinitionConfigCreate(config *ribd.PolicyDefinition) {
	err := server.ProcessPolicyDefinitionConfigCreate(config, server.GlobalPolicyEngineDB)
	fmt.Println("err:", err, " for policy:", config)
}
func PolicyConditionConfigDelete(config *ribd.PolicyCondition) {
	_, err := server.ProcessPolicyConditionConfigDelete(config, server.GlobalPolicyEngineDB)
	fmt.Println("err:", err, " for condition:", config)
}
func PolicyStmtConfigDelete(config *ribd.PolicyStmt) {
	err := server.ProcessPolicyStmtConfigDelete(config, server.GlobalPolicyEngineDB)
	fmt.Println("err:", err, " for stmt:", config)
}
func PolicyDefinitionConfigDelete(config *ribd.PolicyDefinition) {
	err := server.ProcessPolicyDefinitionConfigDelete(config, server.GlobalPolicyEngineDB)
	fmt.Println("err:", err, " for policy:", config)
}
func CallUpdateApplyPolicy(config *ribdInt.ApplyPolicyInfo, apply bool, db *policy.PolicyEngineDB) {
	server.UpdateApplyPolicy(config, apply, db)
}
func TestInitPolicyProcessServer(t *testing.T) {
	fmt.Println("****Init Policy Process Server****")
	StartTestServer()
	InitConditionsList()
	InitStmtsList()
	InitPolicyDefinitionsList()
	InitApplyPolicyInfo()
	fmt.Println("****************")
}
func TestGetBulkPolicyConditionState(t *testing.T) {
	fmt.Println("Policy Conditions in global DB:")
	states, _ := server.GetBulkPolicyConditionState(0, 100, server.GlobalPolicyEngineDB)
	fmt.Println(states)
	fmt.Println("Policy Conditions in rib local DB:")
	states, _ = server.GetBulkPolicyConditionState(0, 100, server.PolicyEngineDB)
	fmt.Println(states)
}
func TestGetBulkPolicyStmtState(t *testing.T) {
	fmt.Println("Policy Statements in global DB:")
	states, _ := server.GetBulkPolicyStmtState(0, 100, server.GlobalPolicyEngineDB)
	fmt.Println(states)
	fmt.Println("Policy Statements in rib local DB:")
	states, _ = server.GetBulkPolicyStmtState(0, 100, server.PolicyEngineDB)
	fmt.Println(states)
}
func TestGetBulkPolicyDefinitionState(t *testing.T) {
	fmt.Println("Policy Definitions in global DB:")
	states, _ := server.GetBulkPolicyDefinitionState(0, 100, server.GlobalPolicyEngineDB)
	fmt.Println(states)
	fmt.Println("Policy Definitions in rib local DB:")
	states, _ = server.GetBulkPolicyDefinitionState(0, 100, server.PolicyEngineDB)
	fmt.Println(states)
}
func TestPolicyConditionConfigCreate(t *testing.T) {
	fmt.Println("**** TestPolicyConditionConfigCreate  ****")
	for _, condition := range conditionsList {
		PolicyConditionConfigCreate(condition)
	}
	TestGetBulkPolicyConditionState(t)
	fmt.Println("***************************************")
}
func TestPolicyStmtConfigCreate(t *testing.T) {
	fmt.Println("**** TestPolicyStmtConfigCreate  ****")
	for _, stmt := range stmtsList {
		PolicyStmtConfigCreate(stmt)
	}
	TestGetBulkPolicyStmtState(t)
	TestPolicyConditionConfigDelete(t)
	fmt.Println("***************************************")
}

func TestPolicyDefinitionConfigCreate(t *testing.T) {
	fmt.Println("**** TestPolicyDefinitionConfigCreate  ****")
	for _, policy := range policyDefinitionsList {
		PolicyDefinitionConfigCreate(policy)
	}
	TestGetBulkPolicyDefinitionState(t)
	TestPolicyStmtConfigDelete(t)
	fmt.Println("***************************************")
}
func TestUpdateApplyPolicy(t *testing.T) {
	fmt.Println("TestUpdateApplyPolicy")
	TestGetBulkPolicyDefinitionState(t)
	for _, applyPolicyInfo := range applyPolicyList {
		fmt.Println("Calling applyPolicyInfo:", applyPolicyInfo, " true, PolicyEngineDB")
		CallUpdateApplyPolicy(applyPolicyInfo, true, PolicyEngineDB)
	}
	for _, applyPolicyInfo := range applyPolicyList {
		fmt.Println("Calling applyPolicyInfo:", applyPolicyInfo, " false, GlobalPolicyEngineDB")
		CallUpdateApplyPolicy(applyPolicyInfo, false, GlobalPolicyEngineDB)
	}
	fmt.Println("*********************")
}
func TestPolicyDefinitionConfigDelete(t *testing.T) {
	fmt.Println("**** TestPolicyDefinitionConfigDelete  ****")
	for _, policy := range policyDefinitionsList {
		PolicyDefinitionConfigDelete(policy)
	}
	TestGetBulkPolicyDefinitionState(t)
	fmt.Println("***************************************")
}
func TestPolicyStmtConfigDelete(t *testing.T) {
	fmt.Println("**** TestPolicyStmtConfigDelete  ****")
	for _, stmt := range stmtsList {
		PolicyStmtConfigDelete(stmt)
	}
	TestGetBulkPolicyStmtState(t)
	fmt.Println("***************************************")
}
func TestPolicyConditionConfigDelete(t *testing.T) {
	fmt.Println("**** TestPolicyConditionConfigDelete  ****")
	for _, condition := range conditionsList {
		PolicyConditionConfigDelete(condition)
	}
	TestGetBulkPolicyConditionState(t)
	fmt.Println("***************************************")
}
func TestPolicyServer(t *testing.T) {
	time.Sleep(1)
	fmt.Println("**** TestPolicyServer ****")
	for _, condition := range conditionsList {
		server.PolicyConfCh <- RIBdServerConfig{
			OrigConfigObject: condition,
			Op:               defs.AddPolicyCondition,
		}
	}
	fmt.Println("Conditions Created")
	time.Sleep(1)
	TestGetBulkPolicyConditionState(t)
	for _, stmt := range stmtsList {
		server.PolicyConfCh <- RIBdServerConfig{
			OrigConfigObject: stmt,
			Op:               defs.AddPolicyStmt,
		}
	}
	fmt.Println("Stmts Created")
	time.Sleep(1)
	TestGetBulkPolicyStmtState(t)
	for _, policy := range policyDefinitionsList {
		server.PolicyConfCh <- RIBdServerConfig{
			OrigConfigObject: policy,
			Op:               defs.AddPolicyDefinition,
		}
	}
	fmt.Println("Definitions Created")
	time.Sleep(10)
	TestGetBulkPolicyConditionState(t)
	TestGetBulkPolicyStmtState(t)
	TestGetBulkPolicyDefinitionState(t)
	for _, applyPolicyInfo := range applyPolicyList {
		server.PolicyConfCh <- RIBdServerConfig{
			PolicyList: ApplyPolicyList{[]*ribdInt.ApplyPolicyInfo{applyPolicyInfo}, make([]*ribdInt.ApplyPolicyInfo, 0)},
			Op:         defs.ApplyPolicy,
		}
	}
	fmt.Println("Policies applied")
	time.Sleep(1)
	for _, policy := range policyDefinitionsList {
		server.PolicyConfCh <- RIBdServerConfig{
			OrigConfigObject: policy,
			Op:               defs.DelPolicyDefinition,
		}
	}
	fmt.Println("Definitions Deleted")
	time.Sleep(1)
	TestGetBulkPolicyDefinitionState(t)
	for _, stmt := range stmtsList {
		server.PolicyConfCh <- RIBdServerConfig{
			OrigConfigObject: stmt,
			Op:               defs.DelPolicyStmt,
		}
	}
	fmt.Println("Stmts Deleted")
	time.Sleep(1)
	TestGetBulkPolicyStmtState(t)
	for _, condition := range conditionsList {
		server.PolicyConfCh <- RIBdServerConfig{
			OrigConfigObject: condition,
			Op:               defs.DelPolicyCondition,
		}
	}
	fmt.Println("Conditions Deleted")
	time.Sleep(1)
	TestGetBulkPolicyConditionState(t)
	fmt.Println("***************************************")
}
