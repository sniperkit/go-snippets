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

// conn.go
package config

import (
	"l3/rib/ribdCommonDefs"
)

type ReachabilityResult struct {
	Err         error
	NextHopInfo *NextHopInfo
}

type ReachabilityInfo struct {
	IP          string
	IfIndex     int32
	ReachableCh chan ReachabilityResult
}

type Operation int

const (
	NOTIFY_ROUTE_CREATED Operation = iota + 1
	NOTIFY_ROUTE_DELETED
	BFD_STATE_VALID
	BFD_STATE_INVALID
	INTF_CREATED
	INTF_DELETED
	INTFV6_CREATED
	INTFV6_DELETED
	IPV6_NEIGHBOR_CREATED
	IPV6_NEIGHBOR_DELETED
	INTF_STATE_DOWN
	INTF_STATE_UP
	NOTIFY_POLICY_CONDITION_CREATED
	NOTIFY_POLICY_CONDITION_DELETED
	NOTIFY_POLICY_CONDITION_UPDATED
	NOTIFY_POLICY_STMT_CREATED
	NOTIFY_POLICY_STMT_DELETED
	NOTIFY_POLICY_STMT_UPDATED
	NOTIFY_POLICY_DEFINITION_CREATED
	NOTIFY_POLICY_DEFINITION_DELETED
	NOTIFY_POLICY_DEFINITION_UPDATED
)

type BfdInfo struct {
	Oper   Operation
	DestIp string
	State  bool
}

type IntfStateInfo struct {
	Idx         int32
	IPAddr      string
	LinkLocalIP string
	State       Operation
}

type IntfMapInfo struct {
	Idx    int32
	IfName string
}

func NewIntfStateInfo(idx int32, ipAddr string, linklocalIp string, state Operation) *IntfStateInfo {
	return &IntfStateInfo{
		Idx:         idx,
		IPAddr:      ipAddr,
		LinkLocalIP: linklocalIp,
		State:       state,
	}
}

/*  This is mimic of ribd object...@TODO: need to change this to bgp server object
 */
type RouteInfo struct {
	IPAddr           string
	Mask             string
	NextHopIp        string
	Prototype        int
	NetworkStatement bool
	RouteOrigin      string
	AddressType      ribdCommonDefs.IPType
}

type RouteCh struct {
	Add    []*RouteInfo
	Remove []*RouteInfo
}

type NextHopInfo struct {
	IPAddr         string
	Mask           string
	Metric         int32
	NextHopIp      string
	IsReachable    bool
	NextHopIfType  int32
	NextHopIfIndex int32
}

type ApplyPolicyInfo struct {
	Protocol   string
	Policy     string
	Action     string
	Conditions []*ConditionInfo
}
