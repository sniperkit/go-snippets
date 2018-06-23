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

package asicdClntDefs

import ()

const (
	//Notification msgs copied from asicd always add notification to the bottom of the list
	NOTIFY_L2INTF_STATE_CHANGE           = iota // 0
	NOTIFY_IPV4_L3INTF_STATE_CHANGE             // 1
	NOTIFY_IPV6_L3INTF_STATE_CHANGE             // 2
	NOTIFY_VLAN_CREATE                          // 3
	NOTIFY_VLAN_DELETE                          // 4
	NOTIFY_VLAN_UPDATE                          // 5
	NOTIFY_LOGICAL_INTF_CREATE                  // 6
	NOTIFY_LOGICAL_INTF_DELETE                  // 7
	NOTIFY_LOGICAL_INTF_UPDATE                  // 8
	NOTIFY_IPV4INTF_CREATE                      // 9
	NOTIFY_IPV4INTF_DELETE                      // 10
	NOTIFY_IPV6INTF_CREATE                      // 11
	NOTIFY_IPV6INTF_DELETE                      // 12
	NOTIFY_LAG_CREATE                           // 13
	NOTIFY_LAG_DELETE                           // 14
	NOTIFY_LAG_UPDATE                           // 15
	NOTIFY_IPV4NBR_MAC_MOVE                     // 16
	NOTIFY_IPV6NBR_MAC_MOVE                     // 17
	NOTIFY_IPV4_ROUTE_CREATE_FAILURE            // 18
	NOTIFY_IPV4_ROUTE_DELETE_FAILURE            // 19
	NOTIFY_IPV6_ROUTE_CREATE_FAILURE            // 20
	NOTIFY_IPV6_ROUTE_DELETE_FAILURE            // 21
	NOTIFY_VTEP_CREATE                          // 22
	NOTIFY_VTEP_DELETE                          // 23
	NOTIFY_MPLSINTF_STATE_CHANGE                // 24
	NOTIFY_MPLSINTF_CREATE                      // 25
	NOTIFY_MPLSINTF_DELETE                      // 26
	NOTIFY_PORT_CONFIG_MODE_CHANGE              // 27
	NOTIFY_PORT_ATTR_CHANGE                     // 28
	NOTIFY_IPV4VIRTUAL_INTF_CREATE              // 29
	NOTIFY_IPV4VIRTUAL_INTF_DELETE              // 30
	NOTIFY_IPV6VIRTUAL_INTF_CREATE              // 31
	NOTIFY_IPV6VIRTUAL_INTF_DELETE              // 32
	NOTIFY_IPV4_VIRTUALINTF_STATE_CHANGE        // 33
	NOTIFY_IPV6_VIRTUALINTF_STATE_CHANGE        // 34
)

type L2IntfStateNotifyMsg struct {
	MsgType uint8
	IfIndex int32
	IfState uint8
}

type IPv4L3IntfStateNotifyMsg struct {
	MsgType uint8
	IpAddr  string
	IfIndex int32
	IfState uint8
}

type IPv6L3IntfStateNotifyMsg struct {
	MsgType uint8
	IpAddr  string
	IfIndex int32
	IfState uint8
}

type VlanNotifyMsg struct {
	MsgType     uint8
	VlanId      uint16
	VlanIfIndex int32
	VlanName    string
	TagPorts    []int32
	UntagPorts  []int32
}

type LogicalIntfNotifyMsg struct {
	MsgType         uint8
	IfIndex         int32
	LogicalIntfName string
}

type LagNotifyMsg struct {
	MsgType     uint8
	LagName     string
	IfIndex     int32
	IfIndexList []int32
}

type IPv4IntfNotifyMsg struct {
	MsgType uint8
	IpAddr  string
	IfIndex int32
	IntfRef string
}

type IPv4NbrMacMoveNotifyMsg struct {
	MsgType uint8
	IpAddr  string
	IfIndex int32
	VlanId  int32
}

type IPv6NbrMacMoveNotifyMsg struct {
	MsgType uint8
	IpAddr  string
	IfIndex int32
	VlanId  int32
}

type IPv4RouteAddDelFailNotifyMsg struct {
	MsgType uint8
	//TODO: Discuss this with Madhavi
}

type IPv6RouteAddDelFailNotifyMsg struct {
	MsgType uint8
	//TODO: Discuss this with Madhavi
}

type VtepNotifyMsg struct {
	MsgType    uint8
	VtepName   string
	IfIndex    int32
	Vni        int32
	SrcIfIndex int32
	SrcIfName  string
}

type MplsIntfStateNotifyMsg struct {
	//TODO: Need to be done along with MPLS Changes
	MsgType uint8
	IfIndex int32
	IfState uint8
}

type MplsIntfNotifyMsg struct {
	//TODO: Need to be done along with MPLS Changes
	MsgType uint8
	IfIndex int32
}

type IPv6IntfNotifyMsg struct {
	MsgType uint8
	IpAddr  string
	IfIndex int32
	IntfRef string
}

type PortConfigModeChgNotifyMsg struct {
	MsgType uint8
	IfIndex int32
	OldMode string
	NewMode string
}

type PortAttrChangeNotifyMsg struct {
	MsgType     uint8
	IfIndex     int32
	Mtu         int32
	Description string
	Pvid        int32
	AttrMask    int32
}

type IPv4VirtualIntfNotifyMsg struct {
	MsgType       uint8
	IfIndex       int32
	ParentIfIndex int32
	IpAddr        string
	MacAddr       string
	IfName        string
}

type IPv6VirtualIntfNotifyMsg struct {
	MsgType       uint8
	IfIndex       int32
	ParentIfIndex int32
	IpAddr        string
	MacAddr       string
	IfName        string
}

type IPv4VirtualIntfStateNotifyMsg struct {
	MsgType uint8
	IfIndex int32
	IpAddr  string
	IfState uint8
}

type IPv6VirtualIntfStateNotifyMsg struct {
	MsgType uint8
	IfIndex int32
	IpAddr  string
	IfState uint8
}
