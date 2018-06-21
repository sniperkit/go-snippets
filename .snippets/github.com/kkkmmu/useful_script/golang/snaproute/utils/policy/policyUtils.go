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

// policyUtils.go
package policy

import (
	"bytes"
	"errors"
	//	"log"
	//	"log/syslog"
	//	"os"
	//"fmt"
	//"strconv"
	//"strings"
	//commondefs "utils/commonDefs"
	bgpUtils "utils/bgpUtils"
	"utils/logging"
	"utils/patriciaDB"
	"utils/policy/policyCommonDefs"
)

const (
	add = iota
	del
	delAll
	invalidate
)
const (
	Invalid = -1
	Valid   = 0
)

type ConditionsAndActionsList struct {
	ConditionList []PolicyCondition
	ActionList    []PolicyAction
}

type PolicyStmtMap struct {
	PolicyStmtMap map[string]ConditionsAndActionsList
}
type PolicyEngineFilterEntityParams struct {
	DestNetIp         string //CIDR format
	NextHopIp         string
	RouteProtocol     string
	Neighbor          string
	Community         []uint32
	ExtendedCommunity []uint64
	ASPath            string
	LocalPref         uint32
	MED               uint32
	CreatePath        bool
	DeletePath        bool
	PolicyList        []string
	PolicyHitCounter  int
}
type PolicyEngineApplyInfo struct {
	ApplyPolicy    ApplyPolicyInfo
	StmtList       []string
	ConditionsList []string
}

//struct sent to the application for updating its local maps/DBs
type PolicyDetails struct {
	Policy        string
	PolicyStmt    string
	ConditionList []PolicyCondition
	ActionList    []PolicyAction
	EntityDeleted bool //whether this policy/stmt resulted in deleting the entity
}
type ApplyPolicyInfo struct {
	ApplyPolicy Policy
	Action      PolicyAction
	Conditions  []string //extra condition names
}
type ApplyPolicyMapInfo struct {
	Count    int
	InfoList []ApplyPolicyInfo
}
type LocalDB struct {
	Prefix     patriciaDB.Prefix
	IsValid    bool
	Precedence int
}
type LocalDBSlice []LocalDB

func (slice *LocalDBSlice) updateLocalDB(prefix patriciaDB.Prefix, op int) {
	if slice == nil {
		return
	}
	tempSlice := *slice
	if op == add {
		localDBRecord := LocalDB{Prefix: prefix, IsValid: true}
		tempSlice = append(tempSlice, localDBRecord)
	} else if op == del {
		found := false
		var i int
		for i = 0; i < len(tempSlice); i++ {
			if bytes.Equal(tempSlice[i].Prefix, prefix) {
				found = true
				break
			}
		}
		if found == true {
			if len(tempSlice) <= i+1 {
				tempSlice = tempSlice[:i]
			} else {
				tempSlice = append(tempSlice[:i], tempSlice[i+1:]...)
			}
		}
	}
	*slice = tempSlice
}

type Policyfunc func(actionInfo interface{}, conditionInfo []interface{}, policy Policy, params interface{}, policyStmt PolicyStmt)
type PolicyConditionCheckfunc func(entity PolicyEngineFilterEntityParams, condition PolicyCondition) bool
type UndoActionfunc func(actionInfo interface{}, conditionList []interface{}, policy Policy, params interface{}, policyStmt PolicyStmt)
type PolicyCheckfunc func(params interface{}) bool
type EntityUpdatefunc func(details PolicyDetails, params interface{})
type PolicyApplyfunc func(entity PolicyEngineFilterEntityParams, policyData interface{}, params interface{})
type EntityTraverseAndApplyPolicyfunc func(data interface{}, updatefunc PolicyApplyfunc)
type EntityTraverseAndReversePolicyfunc func(data interface{})
type PolicyEntityMapIndex interface{}
type GetPolicyEnityMapIndexFunc func(entity PolicyEngineFilterEntityParams, policy string) PolicyEntityMapIndex

type PolicyEngineDB struct {
	Logger                            *logging.Writer //*log.Logger
	PolicyPrefixSetDB                 *patriciaDB.Trie
	LocalPolicyPrefixSetDB            *LocalDBSlice
	PolicyASPathSetDB                 *patriciaDB.Trie
	LocalPolicyASPathSetDB            *LocalDBSlice
	PolicyCommunitySetDB              *patriciaDB.Trie
	LocalPolicyCommunitySetDB         *LocalDBSlice
	PolicyExtendedCommunitySetDB      *patriciaDB.Trie
	LocalPolicyExtendedCommunitySetDB *LocalDBSlice
	PolicyConditionsDB                *patriciaDB.Trie
	LocalPolicyConditionsDB           *LocalDBSlice
	PolicyActionsDB                   *patriciaDB.Trie
	LocalPolicyActionsDB              *LocalDBSlice
	PolicyStmtDB                      *patriciaDB.Trie
	LocalPolicyStmtDB                 *LocalDBSlice
	PolicyDB                          *patriciaDB.Trie
	LocalPolicyDB                     *LocalDBSlice
	PolicyStmtPolicyMapDB             map[string][]string //policies using this statement
	PrefixPolicyListDB                *patriciaDB.Trie
	ProtocolPolicyListDB              map[string][]string //policystmt names assoociated with every protocol type
	ImportPolicyPrecedenceMap         map[int]string
	ExportPolicyPrecedenceMap         map[int]string
	ApplyPolicyMap                    map[string]ApplyPolicyMapInfo
	PolicyEntityMap                   map[PolicyEntityMapIndex]PolicyStmtMap
	DefaultImportPolicyActionFunc     Policyfunc
	DefaultExportPolicyActionFunc     Policyfunc
	IsEntityPresentFunc               PolicyCheckfunc
	GetPolicyEntityMapIndex           GetPolicyEnityMapIndexFunc
	UpdateEntityDB                    EntityUpdatefunc
	ConditionCheckfuncMap             map[int]PolicyConditionCheckfunc
	ActionfuncMap                     map[int]Policyfunc
	UndoActionfuncMap                 map[int]UndoActionfunc
	TraverseAndApplyPolicyFunc        EntityTraverseAndApplyPolicyfunc
	TraverseAndReversePolicyFunc      EntityTraverseAndReversePolicyfunc
	ValidConditionsForPolicyTypeMap   map[string][]int //map of policyType to list of valid conditions
	ValidActionsForPolicyTypeMap      map[string][]int //map of policyType to list of valid actions
	Global                            bool             //this variable is to say whether this engine is for storing the policies only (true)) or the actual engine : default is false, meaning it is an application engine
}

func (db *PolicyEngineDB) buildPolicyConditionCheckfuncMap() {
	db.Logger.Info("buildPolicyConditionCheckfuncMap")
	db.ConditionCheckfuncMap[policyCommonDefs.PolicyConditionTypeDstIpPrefixMatch] = db.DstIpPrefixMatchConditionfunc
	db.ConditionCheckfuncMap[policyCommonDefs.PolicyConditionTypeProtocolMatch] = db.ProtocolMatchConditionfunc
	db.ConditionCheckfuncMap[policyCommonDefs.PolicyConditionTypeNeighborMatch] = db.NeighborMatchConditionfunc
	db.ConditionCheckfuncMap[policyCommonDefs.PolicyConditionTypeCommunityMatch] = db.CommunityMatchConditionfunc
	db.ConditionCheckfuncMap[policyCommonDefs.PolicyConditionTypeExtendedCommunityMatch] = db.ExtendedCommunityMatchConditionfunc
	db.ConditionCheckfuncMap[policyCommonDefs.PolicyConditionTypeASPathMatch] = db.ASPathMatchConditionfunc
	db.ConditionCheckfuncMap[policyCommonDefs.PolicyConditionTypeLocalPrefMatch] = db.LocalPrefMatchConditionfunc
	db.ConditionCheckfuncMap[policyCommonDefs.PolicyConditionTypeMEDMatch] = db.MEDMatchConditionfunc
}
func (db *PolicyEngineDB) buildPolicyValidConditionsForPolicyTypeMap() {
	db.Logger.Info("buildPolicyValidConditionsForPolicyTypeMap")
	db.ValidConditionsForPolicyTypeMap["ALL"] = []int{policyCommonDefs.PolicyConditionTypeDstIpPrefixMatch,
		policyCommonDefs.PolicyConditionTypeProtocolMatch,
		policyCommonDefs.PolicyConditionTypeNeighborMatch,
		policyCommonDefs.PolicyConditionTypeCommunityMatch,
		policyCommonDefs.PolicyConditionTypeExtendedCommunityMatch,
		policyCommonDefs.PolicyConditionTypeLocalPrefMatch,
		policyCommonDefs.PolicyConditionTypeASPathMatch,
		policyCommonDefs.PolicyConditionTypeMEDMatch}
	db.ValidConditionsForPolicyTypeMap["BGP"] = []int{policyCommonDefs.PolicyConditionTypeDstIpPrefixMatch,
		policyCommonDefs.PolicyConditionTypeNeighborMatch,
		policyCommonDefs.PolicyConditionTypeCommunityMatch,
		policyCommonDefs.PolicyConditionTypeExtendedCommunityMatch,
		policyCommonDefs.PolicyConditionTypeLocalPrefMatch,
		policyCommonDefs.PolicyConditionTypeASPathMatch,
		policyCommonDefs.PolicyConditionTypeMEDMatch}
	db.ValidConditionsForPolicyTypeMap["OSPF"] = []int{policyCommonDefs.PolicyConditionTypeDstIpPrefixMatch}
}
func (db *PolicyEngineDB) buildPolicyValidActionsForPolicyTypeMap() {
	db.Logger.Info("buildPolicyValidActionsForPolicyTypeMap")
	db.ValidActionsForPolicyTypeMap["ALL"] = []int{policyCommonDefs.PolicyActionTypeRouteDisposition,
		policyCommonDefs.PolicyActionTypeRouteRedistribute,
		policyCommonDefs.PolicyActionTypeRIBIn,
		policyCommonDefs.PolicyActionTypeRIBOut,
		policyCommonDefs.PolicyActionTypeSetCommunity,
		policyCommonDefs.PolicyActionTypeSetExtendedCommunity,
		policyCommonDefs.PolicyActionTypeSetLocalPref,
		policyCommonDefs.PolicyActionTypeSetPrependASPath,
		policyCommonDefs.PolicyActionTypeSetMED}
	db.ValidActionsForPolicyTypeMap["BGP"] = []int{policyCommonDefs.PolicyActionTypeAggregate,
		policyCommonDefs.PolicyActionTypeRIBIn,
		policyCommonDefs.PolicyActionTypeRIBOut,
		policyCommonDefs.PolicyActionTypeSetCommunity,
		policyCommonDefs.PolicyActionTypeSetExtendedCommunity,
		policyCommonDefs.PolicyActionTypeSetLocalPref,
		policyCommonDefs.PolicyActionTypeSetPrependASPath,
		policyCommonDefs.PolicyActionTypeSetMED}
}
func NewPolicyEngineDB(logger *logging.Writer) (policyEngineDB *PolicyEngineDB) {
	policyEngineDB = &PolicyEngineDB{}
	/*	if policyEngineDB.Logger == nil {
		policyEngineDB.Logger = log.New(os.Stdout, "PolicyEngine :", log.Ldate|log.Ltime|log.Lshortfile)

		syslogger, err := syslog.New(syslog.LOG_NOTICE|syslog.LOG_INFO|syslog.LOG_DAEMON, "PolicyEngine")
		if err == nil {
			syslogger.Info("### PolicyEngineDB initailized")
			policyEngineDB.Logger.SetOutput(syslogger)
		}
	}*/
	policyEngineDB.Logger = logger
	policyEngineDB.PolicyActionsDB = patriciaDB.NewTrie()
	LocalPolicyActionsDB := make([]LocalDB, 0)
	localActionSlice := LocalDBSlice(LocalPolicyActionsDB)
	policyEngineDB.LocalPolicyActionsDB = &localActionSlice

	policyEngineDB.PolicyPrefixSetDB = patriciaDB.NewTrie()
	LocalPolicyPrefixSetDB := make([]LocalDB, 0)
	localPrefixSetSlice := LocalDBSlice(LocalPolicyPrefixSetDB)
	policyEngineDB.LocalPolicyPrefixSetDB = &localPrefixSetSlice

	policyEngineDB.PolicyCommunitySetDB = patriciaDB.NewTrie()
	LocalPolicyCommunitySetDB := make([]LocalDB, 0)
	localCommunitySetSlice := LocalDBSlice(LocalPolicyCommunitySetDB)
	policyEngineDB.LocalPolicyCommunitySetDB = &localCommunitySetSlice

	policyEngineDB.PolicyExtendedCommunitySetDB = patriciaDB.NewTrie()
	LocalPolicyExtendedCommunitySetDB := make([]LocalDB, 0)
	localExtendedCommunitySetSlice := LocalDBSlice(LocalPolicyExtendedCommunitySetDB)
	policyEngineDB.LocalPolicyExtendedCommunitySetDB = &localExtendedCommunitySetSlice

	policyEngineDB.PolicyASPathSetDB = patriciaDB.NewTrie()
	LocalPolicyASPathSetDB := make([]LocalDB, 0)
	localASPathSetSlice := LocalDBSlice(LocalPolicyASPathSetDB)
	policyEngineDB.LocalPolicyASPathSetDB = &localASPathSetSlice

	policyEngineDB.PolicyConditionsDB = patriciaDB.NewTrie()
	LocalPolicyConditionsDB := make([]LocalDB, 0)
	localConditionSlice := LocalDBSlice(LocalPolicyConditionsDB)
	policyEngineDB.LocalPolicyConditionsDB = &localConditionSlice

	policyEngineDB.PolicyStmtDB = patriciaDB.NewTrie()
	LocalPolicyStmtDB := make([]LocalDB, 0)
	localStmtSlice := LocalDBSlice(LocalPolicyStmtDB)
	policyEngineDB.LocalPolicyStmtDB = &localStmtSlice

	policyEngineDB.PolicyDB = patriciaDB.NewTrie()
	LocalPolicyDB := make([]LocalDB, 0)
	localPolicySlice := LocalDBSlice(LocalPolicyDB)
	policyEngineDB.LocalPolicyDB = &localPolicySlice

	policyEngineDB.PolicyStmtPolicyMapDB = make(map[string][]string)
	policyEngineDB.PolicyEntityMap = make(map[PolicyEntityMapIndex]PolicyStmtMap)
	policyEngineDB.PrefixPolicyListDB = patriciaDB.NewTrie()
	policyEngineDB.ProtocolPolicyListDB = make(map[string][]string)
	policyEngineDB.ImportPolicyPrecedenceMap = make(map[int]string)
	policyEngineDB.ExportPolicyPrecedenceMap = make(map[int]string)
	policyEngineDB.ApplyPolicyMap = make(map[string]ApplyPolicyMapInfo)
	policyEngineDB.ConditionCheckfuncMap = make(map[int]PolicyConditionCheckfunc)
	policyEngineDB.ValidConditionsForPolicyTypeMap = make(map[string][]int)
	policyEngineDB.ValidActionsForPolicyTypeMap = make(map[string][]int)
	policyEngineDB.buildPolicyConditionCheckfuncMap()
	policyEngineDB.buildPolicyValidConditionsForPolicyTypeMap()
	policyEngineDB.buildPolicyValidActionsForPolicyTypeMap()
	policyEngineDB.ActionfuncMap = make(map[int]Policyfunc)
	policyEngineDB.UndoActionfuncMap = make(map[int]UndoActionfunc)
	policyEngineDB.Global = false
	return policyEngineDB
}

func (db *PolicyEngineDB) SetDefaultImportPolicyActionFunc(defaultfunc Policyfunc) {
	db.DefaultImportPolicyActionFunc = defaultfunc
}
func (db *PolicyEngineDB) SetDefaultExportPolicyActionFunc(defaultfunc Policyfunc) {
	db.DefaultExportPolicyActionFunc = defaultfunc
}
func (db *PolicyEngineDB) SetIsEntityPresentFunc(IsPresent PolicyCheckfunc) {
	db.IsEntityPresentFunc = IsPresent
}
func (db *PolicyEngineDB) SetEntityUpdateFunc(updatefunc EntityUpdatefunc) {
	db.UpdateEntityDB = updatefunc
}
func (db *PolicyEngineDB) SetActionFunc(action int, setfunc Policyfunc) {
	db.ActionfuncMap[action] = setfunc
}
func (db *PolicyEngineDB) SetUndoActionFunc(action int, setfunc UndoActionfunc) {
	db.UndoActionfuncMap[action] = setfunc
}
func (db *PolicyEngineDB) SetTraverseAndApplyPolicyFunc(updatefunc EntityTraverseAndApplyPolicyfunc) {
	db.TraverseAndApplyPolicyFunc = updatefunc
}
func (db *PolicyEngineDB) SetTraverseAndReversePolicyFunc(updatefunc func(policyItem interface{})) {
	db.TraverseAndReversePolicyFunc = updatefunc
}
func (db *PolicyEngineDB) SetGetPolicyEntityMapIndexFunc(getfunc func(entity PolicyEngineFilterEntityParams, policy string) PolicyEntityMapIndex) {
	db.GetPolicyEntityMapIndex = getfunc
}
func isPolicyTypeSame(oldPolicy Policy, policy Policy) (same bool) {
	if oldPolicy.ExportPolicy == policy.ExportPolicy && oldPolicy.ImportPolicy == policy.ImportPolicy {
		same = true
	}
	return same
}
func (db *PolicyEngineDB) AddPolicyEntityMapEntry(entity PolicyEngineFilterEntityParams, policy string,
	policyStmt string, conditionList []PolicyCondition, actionList []PolicyAction) {
	db.Logger.Info("AddPolicyEntityMapEntry")
	var policyStmtMap PolicyStmtMap
	var conditionsAndActionsList ConditionsAndActionsList
	if db.PolicyEntityMap == nil {
		db.PolicyEntityMap = make(map[PolicyEntityMapIndex]PolicyStmtMap)
	}
	if db.GetPolicyEntityMapIndex == nil {
		return
	}
	policyEntityMapIndex := db.GetPolicyEntityMapIndex(entity, policy)
	if policyEntityMapIndex == nil {
		db.Logger.Err("policyEntityMapKey nil")
		return
	}
	policyStmtMap, ok := db.PolicyEntityMap[policyEntityMapIndex]
	if !ok {
		policyStmtMap.PolicyStmtMap = make(map[string]ConditionsAndActionsList)
	}
	_, ok = policyStmtMap.PolicyStmtMap[policyStmt]
	if ok {
		db.Logger.Err("policy statement map for statement ", policyStmt, " already in place for policy ", policy, " : ", policyStmtMap.PolicyStmtMap[policyStmt])
		//	conditionsAndActionsList.ConditionList = policyStmtMap.PolicyStmtMap[policyStmt].ConditionList
		//	conditionsAndActionsList.ActionList = policyStmtMap.PolicyStmtMap[policyStmt].ActionList
		return
	} //else {
	conditionsAndActionsList.ConditionList = make([]PolicyCondition, 0)
	conditionsAndActionsList.ActionList = make([]PolicyAction, 0)
	//}
	for i := 0; conditionList != nil && i < len(conditionList); i++ {
		conditionsAndActionsList.ConditionList = append(conditionsAndActionsList.ConditionList, conditionList[i])
	}
	for i := 0; actionList != nil && i < len(actionList); i++ {
		conditionsAndActionsList.ActionList = append(conditionsAndActionsList.ActionList, actionList[i])
	}
	policyStmtMap.PolicyStmtMap[policyStmt] = conditionsAndActionsList
	db.PolicyEntityMap[policyEntityMapIndex] = policyStmtMap
	db.Logger.Info("Adding entry for index ", policyEntityMapIndex)
}
func (db *PolicyEngineDB) DeletePolicyEntityMapEntry(entity PolicyEngineFilterEntityParams, policy string) {
	db.Logger.Info("DeletePolicyEntityMapEntry for policy ", policy)
	if db.PolicyEntityMap == nil {
		db.Logger.Err("PolicyEntityMap empty")
		return
	}
	if db.GetPolicyEntityMapIndex == nil {
		return
	}
	policyEntityMapIndex := db.GetPolicyEntityMapIndex(entity, policy)
	if policyEntityMapIndex == nil {
		db.Logger.Err("policyEntityMapIndex nil")
		return
	}
	//PolicyRouteMap[policyRouteIndex].policyStmtMap=nil
	delete(db.PolicyEntityMap, policyEntityMapIndex)
}
func (db *PolicyEngineDB) PolicyActionType(actionType int) (exportTypeAction bool, importTypeAction bool, globalTypeAction bool) {
	db.Logger.Info("PolicyActionType for type ", actionType)
	switch actionType {
	case policyCommonDefs.PoilcyActionTypeSetAdminDistance:
		globalTypeAction = true
		db.Logger.Info("PoilcyActionTypeSetAdminDistance, setting globalTypeAction true")
		break
	case policyCommonDefs.PolicyActionTypeAggregate:
		exportTypeAction = true
		db.Logger.Info("PolicyActionTypeAggregate: setting exportTypeAction true")
		break
	case policyCommonDefs.PolicyActionTypeRouteRedistribute:
		exportTypeAction = true
		db.Logger.Info("PolicyActionTypeRouteRedistribute: setting exportTypeAction true")
		break
	case policyCommonDefs.PolicyActionTypeNetworkStatementAdvertise:
		exportTypeAction = true
		db.Logger.Info("PolicyActionTypeNetworkStatementAdvertise: setting exportTypeAction true")
		break
	case policyCommonDefs.PolicyActionTypeRouteDisposition:
		importTypeAction = true
		db.Logger.Info("setting importTypeAction true")
		break
	case policyCommonDefs.PolicyActionTypeRIBIn:
		importTypeAction = true
		db.Logger.Info("setting importTypeAction true")
		break
	case policyCommonDefs.PolicyActionTypeRIBOut:
		exportTypeAction = true
		db.Logger.Info("setting exportTypeAction true")
		break
	case policyCommonDefs.PolicyActionTypeSetCommunity:
		exportTypeAction = true
		db.Logger.Info("setting exportTypeAction true")
		break
	case policyCommonDefs.PolicyActionTypeSetLocalPref:
		exportTypeAction = true
		db.Logger.Info("setting exportTypeAction true")
		break
	case policyCommonDefs.PolicyActionTypeSetMED:
		exportTypeAction = true
		db.Logger.Info("setting exportTypeAction true")
		break
	case policyCommonDefs.PolicyActionTypeSetExtendedCommunity:
		exportTypeAction = true
		db.Logger.Info("setting exportTypeAction true")
		break
	case policyCommonDefs.PolicyActionTypeSetPrependASPath:
		exportTypeAction = true
		db.Logger.Info("setting exportTypeAction true")
		break
	default:
		db.Logger.Err("Unknown action type")
		break
	}
	return exportTypeAction, importTypeAction, globalTypeAction
}
func PolicyActionStrToIntType(action string) (actionType int, err error) {
	switch action {
	case "RouteDisposition":
		actionType = policyCommonDefs.PolicyActionTypeRouteDisposition
		break
	case "Redistribution":
		actionType = policyCommonDefs.PolicyActionTypeRouteRedistribute
		break
	case "SetAdminDistance":
		actionType = policyCommonDefs.PoilcyActionTypeSetAdminDistance
		break
	case "NetworkStatementAdvertise":
		actionType = policyCommonDefs.PolicyActionTypeNetworkStatementAdvertise
		break
	case "Aggregate":
		actionType = policyCommonDefs.PolicyActionTypeAggregate
		break
	case "RIBIn":
		actionType = policyCommonDefs.PolicyActionTypeRIBIn
		break
	case "RIBOut":
		actionType = policyCommonDefs.PolicyActionTypeRIBOut
		break
	case "LocalPref":
		actionType = policyCommonDefs.PolicyActionTypeSetLocalPref
		break
	case "Community":
		actionType = policyCommonDefs.PolicyActionTypeSetCommunity
		break
	case "ExtendedCommunity":
		actionType = policyCommonDefs.PolicyActionTypeSetExtendedCommunity
		break
	case "MED":
		actionType = policyCommonDefs.PolicyActionTypeSetMED
		break
	case "PrependASPath":
		actionType = policyCommonDefs.PolicyActionTypeSetPrependASPath
		break
	default:
		return -1, errors.New("Unknown ActionType")
	}
	return actionType, nil
}
func HasActionInfo(infoLIst []ApplyPolicyInfo, action PolicyAction) bool {
	for i := 0; i < len(infoLIst); i++ {
		if infoLIst[i].Action.ActionType == action.ActionType && infoLIst[i].Action.ActionInfo == action.ActionInfo {
			return true
		}
	}
	return false
}

/*func (db *PolicyEngineDB) SetAndValidatePolicyType(policy *Policy, stmt PolicyStmt) (err error) {
	db.Logger.Info(fmt.Sprintln("SetPolicyTypeFromPolicyStmt"))
	if policy.ExportPolicy == false && policy.ImportPolicy == false && policy.GlobalPolicy == false {
		db.Logger.Info(fmt.Sprintln("Policy is still not associated with a type, set it from stmt"))
		policy.ExportPolicy = stmt.ExportStmt
		policy.ImportPolicy = stmt.ImportStmt
		policy.GlobalPolicy = stmt.GlobalStmt

		if policy.ImportPolicy && db.ImportPolicyPrecedenceMap != nil {
			_, ok := db.ImportPolicyPrecedenceMap[int(policy.Precedence)]
			if ok {
				db.Logger.Err(fmt.Sprintln("There is already a import policy with this precedence."))
				err = errors.New("There is already a import policy with this precedence.")
				return err
			}
		} else if policy.ExportPolicy && db.ExportPolicyPrecedenceMap != nil {
			_, ok := db.ExportPolicyPrecedenceMap[int(policy.Precedence)]
			if ok {
				db.Logger.Err(fmt.Sprintln("There is already a export policy with this precedence."))
				err = errors.New("There is already a export policy with this precedence.")
				return err
			}
		} else if policy.GlobalPolicy {
			db.Logger.Info(fmt.Sprintln("This is a global policy"))
		}
		return err
	}
	if policy.ExportPolicy != stmt.ExportStmt ||
		policy.ImportPolicy != stmt.ImportStmt ||
		policy.GlobalPolicy != stmt.GlobalStmt {
		db.Logger.Err(fmt.Sprintln("Policy type settings, export/import/global :", policy.ExportPolicy, "/", policy.ImportPolicy, "/", policy.GlobalPolicy, " does not match the export/import/global settings on the stmt: ", stmt.ExportStmt, "/", stmt.ImportStmt, "/", stmt.GlobalStmt))
		err = errors.New("Mismatch on policy type")
		return err
	}
	return err
}
*/
func (db *PolicyEngineDB) ConditionCheckForPolicyType(conditionName string, policyType string) bool {
	validList := db.ValidConditionsForPolicyTypeMap[policyType]
	if validList == nil || len(validList) == 0 {
		db.Logger.Err("Valid Conditions not defined for policyType: ", policyType)
		return false
	}
	item := db.PolicyConditionsDB.Get(patriciaDB.Prefix(conditionName))
	if item == nil {
		db.Logger.Err("Condtition with conditionName ", conditionName, " not defined")
		return false
	}
	condition := item.(PolicyCondition)
	for j := 0; j < len(validList); j++ {
		if validList[j] == condition.ConditionType {
			db.Logger.Info("Condition ", conditionName, " valid for policyType: ", policyType)
			return true
		}
	}
	db.Logger.Info("Condition ", conditionName, " not valid for policyType: ", policyType)
	return false
} /*
func (db *PolicyEngineDB) GetCommunityValue(inp string) (uint32, error) {
	var val uint32
	info, ok := commondefs.BGPWellKnownCommunitiesMap[inp]
	if ok {
		val = info
	} else if strings.HasPrefix(inp, "0x") {
		info, err := strconv.ParseInt(inp, 0, 64)
		if err != nil {
			return val, err
		} else {
			val = uint32(info)
		}
	} else {
		//split with :
		a := strings.Split(inp, ":")
		if len(a) > 2 {
			db.Logger.Err("Incorrect format for community ", inp)
			return val, errors.New(fmt.Sprintln("Incorrect format for community ", inp))
		}
		if len(a) == 2 {
			a0num, _ := strconv.Atoi(a[0])
			as := fmt.Sprintf("%4x", a0num)
			as = strings.Replace(as, " ", "0", -1)
			a1num, _ := strconv.Atoi(a[1])
			num := fmt.Sprintf("%4x", a1num)
			num = strings.Replace(num, " ", "0", -1)

			comm := "0x" + as + num
			valint, err := strconv.ParseInt(comm, 0, 64)
			if err != nil {
				return val, err
			}
			val = uint32(valint)
		} else if len(a) == 1 {
			//just a integer
			info, err := strconv.Atoi(inp)
			db.Logger.Info("err:", err, " while caling strconv for ", inp)
			if err == nil {
				val = uint32(info)
			} else {
				db.Logger.Err("Incorrect community input:", inp)
				return val, err
			}
		} else {
			db.Logger.Err("Incorrect community input:", inp)
			return val, errors.New(fmt.Sprintln("Incorrect community input:", inp))
		}
	}
	return val, nil
}*/
func (db *PolicyEngineDB) GetPolicySetAction(in PolicyActionCfg) (out PolicyActionState, err error) {
	if _, ok := policyCommonDefs.SetActionMap[in.Attr]; !ok {
		db.Logger.Err("GetPolicySetAction, inValid attr:", in.Attr)
		return out, errors.New("Invalid attr")
	}
	out.Attr = policyCommonDefs.SetActionMap[in.Attr]
	switch out.Attr {
	case policyCommonDefs.PolicyActionTypeSetLocalPref:
		out.LocalPref = uint32(in.LocalPref)
		break
	case policyCommonDefs.PolicyActionTypeSetMED:
		out.MED = uint32(in.MED)
		break
	case policyCommonDefs.PolicyActionTypeSetPrependASPath:
		out.PrependASPath = in.PrependASPath
		break
	case policyCommonDefs.PolicyActionTypeSetCommunity:
		out.Community, err = bgpUtils.GetCommunityValue(in.Community)
		if err != nil {
			db.Logger.Err("SetAction returned error:", err, " while setting community policy action")
			return out, err
		}
		break
	case policyCommonDefs.PolicyActionTypeSetExtendedCommunity:
		out.ExtendedCommunity, err = bgpUtils.EncodeExtCommunity(bgpUtils.ExtCommunity{in.ExtendedCommunity.Type, in.ExtendedCommunity.Value})
		if err != nil {
			db.Logger.Err("SetAction returned error:", err, " while setting extended community policy action")
			return out, err
		}
		break
	}
	return out, err
}
