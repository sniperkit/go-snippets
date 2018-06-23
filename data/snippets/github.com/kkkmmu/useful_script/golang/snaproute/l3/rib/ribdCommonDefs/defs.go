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

package ribdCommonDefs

import (
	"ribdInt"
	"utils/commonDefs"
)

type IPType int

const (
	IPv4 IPType = iota
	IPv6
)

type OpType int

const (
	NoOp OpType = iota
	Add
	AddBulk
	Del
	Update
	Addv6
	Delv6
	Updatev6
	Get
	AddFIBOnly
	DelFIBOnly
	Addv6FIBOnly
	Delv6FIBOnly
	AsicdFetchv4
	AsicdFetchv6
	DBFetch
	AddPolicyCondition
	DelPolicyCondition
	UpdatePolicyCondition
	AddPolicyAction
	DelPolicyAction
	UpdatePolicyAction
	AddPolicyStmt
	DelPolicyStmt
	UpdatePolicyStmt
	AddPolicyPrefixSet
	DelPolicyPrefixSet
	UpdatePolicyPrefixSet
	AddPolicyCommunitySet
	DelPolicyCommunitySet
	UpdatePolicyCommunitySet
	AddPolicyExtendedCommunitySet
	DelPolicyExtendedCommunitySet
	UpdatePolicyExtendedCommunitySet
	AddPolicyASPathSet
	DelPolicyASPathSet
	UpdatePolicyASPathSet
	AddPolicyDefinition
	DelPolicyDefinition
	UpdatePolicyDefinition
	ApplyPolicy
)
const (
	CONNECTED                                    = 0
	STATIC                                       = 1
	OSPF                                         = 89
	EBGP                                         = 8
	IBGP                                         = 9
	BGP                                          = 17
	PUB_SOCKET_ADDR                              = "ipc:///tmp/ribd.ipc"
	PUB_SOCKET_BGPD_ADDR                         = "ipc:///tmp/ribd_bgpd.ipc"
	PUB_SOCKET_OSPFD_ADDR                        = "ipc:///tmp/ribd_ospfd.ipc"
	PUB_SOCKET_BFDD_ADDR                         = "ipc:///tmp/ribd_bfdd.ipc"
	PUB_SOCKET_VXLAND_ADDR                       = "ipc:///tmp/ribd_vxland.ipc"
	PUB_SOCKET_POLICY_ADDR                       = "ipc:///tmp/ribd_policyd.ipc"
	NOTIFY_ROUTE_CREATED                         = 1
	NOTIFY_ROUTE_DELETED                         = 2
	NOTIFY_ROUTE_INVALIDATED                     = 3
	NOTIFY_ROUTE_REACHABILITY_STATUS_UPDATE      = 4
	NOTIFY_POLICY_CONDITION_CREATED              = 5
	NOTIFY_POLICY_CONDITION_DELETED              = 6
	NOTIFY_POLICY_CONDITION_UPDATED              = 7
	NOTIFY_POLICY_STMT_CREATED                   = 8
	NOTIFY_POLICY_STMT_DELETED                   = 9
	NOTIFY_POLICY_STMT_UPDATED                   = 10
	NOTIFY_POLICY_DEFINITION_CREATED             = 11
	NOTIFY_POLICY_DEFINITION_DELETED             = 12
	NOTIFY_POLICY_DEFINITION_UPDATED             = 13
	NOTIFY_POLICY_PREFIX_SET_CREATED             = 14
	NOTIFY_POLICY_PREFIX_SET_DELETED             = 15
	NOTIFY_POLICY_PREFIX_SET_UPDATED             = 16
	NOTIFY_POLICY_COMMUNITY_SET_CREATED          = 17
	NOTIFY_POLICY_COMMUNITY_SET_DELETED          = 18
	NOTIFY_POLICY_COMMUNITY_SET_UPDATED          = 19
	NOTIFY_POLICY_EXTENDED_COMMUNITY_SET_CREATED = 20
	NOTIFY_POLICY_EXTENDED_COMMUNITY_SET_DELETED = 21
	NOTIFY_POLICY_EXTENDED_COMMUNITY_SET_UPDATED = 22
	NOTIFY_POLICY_ASPATH_SET_CREATED             = 23
	NOTIFY_POLICY_ASPATH_SET_DELETED             = 24
	NOTIFY_POLICY_ASPATH_SET_UPDATED             = 25
	DEFAULT_NOTIFICATION_SIZE                    = 128
	RoutePolicyStateChangetoValid                = 1
	RoutePolicyStateChangetoInValid              = 2
	RoutePolicyStateChangeNoChange               = 3
)

type RibdNotifyMsg struct {
	MsgType uint16
	MsgBuf  []byte
}

type RoutelistInfo struct {
	RouteInfo ribdInt.Routes
}
type RouteReachabilityStatusMsgInfo struct {
	Network     string
	IsReachable bool
	NextHopIntf ribdInt.NextHopInfo
}

func GetNextHopIfTypeStr(nextHopIfType ribdInt.Int) (nextHopIfTypeStr string, err error) {
	nextHopIfTypeStr = ""
	switch nextHopIfType {
	case commonDefs.IfTypePort:
		nextHopIfTypeStr = "PHY"
		break
	case commonDefs.IfTypeVlan:
		nextHopIfTypeStr = "VLAN"
		break
	case commonDefs.IfTypeNull:
		nextHopIfTypeStr = "NULL"
		break
	case commonDefs.IfTypeLoopback:
		nextHopIfTypeStr = "Loopback"
		break
	}
	return nextHopIfTypeStr, err
}
