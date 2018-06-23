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

package commonDefs

import (
	"utils/logging"
)

type IPv4IntfState struct {
	IntfRef           string
	IfIndex           int32
	IpAddr            string
	OperState         string
	NumUpEvents       int32
	LastUpEventTime   string
	NumDownEvents     int32
	LastDownEventTime string
	L2IntfType        string
	L2IntfId          int32
}

type IPv4IntfStateGetInfo struct {
	StartIdx          int32
	EndIdx            int32
	Count             int32
	More              bool
	IPv4IntfStateList []IPv4IntfState
}

type Port struct {
	IntfRef     string
	IfIndex     int32
	Description string
	PhyIntfType string
	AdminState  string
	MacAddr     string
	Speed       int32
	Duplex      string
	Autoneg     string
	MediaType   string
	Mtu         int32
}

type PortGetInfo struct {
	StartIdx int32
	EndIdx   int32
	Count    int32
	More     bool
	PortList []Port
}

type PortState struct {
	IntfRef           string
	IfIndex           int32
	Name              string
	OperState         string
	NumUpEvents       int32
	LastUpEventTime   string
	NumDownEvents     int32
	LastDownEventTime string
	Pvid              int32
	IfInOctets        int64
	IfInUcastPkts     int64
	IfInDiscards      int64
	IfInErrors        int64
	IfInUnknownProtos int64
	IfOutOctets       int64
	IfOutUcastPkts    int64
	IfOutDiscards     int64
	IfOutErrors       int64
	ErrDisableReason  string
}

type PortStateGetInfo struct {
	StartIdx      int32
	EndIdx        int32
	Count         int32
	More          bool
	PortStateList []PortState
}

type Lag struct {
	LagIfIndex  int32
	HashType    int32
	IfIndexList []int32
	LagName     string
}

type LagGetInfo struct {
	StartIdx int32
	EndIdx   int32
	Count    int32
	More     bool
	LagList  []Lag
}

type Vlan struct {
	VlanId           int32
	IfIndexList      []int32
	UntagIfIndexList []int32
}

type VlanGetInfo struct {
	StartIdx int32
	EndIdx   int32
	Count    int32
	More     bool
	VlanList []Vlan
}

type VlanState struct {
	VlanId    int32
	VlanName  string
	OperState string
	IfIndex   int32
}

type VlanStateGetInfo struct {
	StartIdx      int32
	EndIdx        int32
	Count         int32
	More          bool
	VlanStateList []VlanState
}

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

const (
	// state values copied from asicd
	INTF_STATE_DOWN = 0
	INTF_STATE_UP   = 1
)

const (
	// this needs to match asicd server
	PORT_ATTR_PHY_INTF_TYPE = 0x00000001
	PORT_ATTR_ADMIN_STATE   = 0x00000002
	PORT_ATTR_MAC_ADDR      = 0x00000004
	PORT_ATTR_SPEED         = 0x00000008
	PORT_ATTR_DUPLEX        = 0x00000010
	PORT_ATTR_AUTONEG       = 0x00000020
	PORT_ATTR_MEDIA_TYPE    = 0x00000040
	PORT_ATTR_MTU           = 0x00000080
	PORT_ATTR_BREAKOUT_MODE = 0x00000100
	PORT_ATTR_LOOPBACK_MODE = 0x00000200
	PORT_ATTR_ENABLE_FEC    = 0x00000400
	PORT_ATTR_TX_PRBS_EN    = 0x00000800
	PORT_ATTR_RX_PRBS_EN    = 0x00001000
	PORT_ATTR_PRBS_POLY     = 0x00002000
	PORT_ATTR_DESCRIPTION   = 0x00004000
	PORT_ATTR_PVID          = 0x00008000
)

type AsicdNotification map[uint8]bool

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

type IPv6IntfNotifyMsg struct {
	MsgType uint8
	IpAddr  string
	IfIndex int32
	IntfRef string
}

type PortConfigModeChgNotifyMsg struct {
	IfIndex int32
	OldMode string
	NewMode string
}

type PortAttrChangeNotifyMsg struct {
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

type AsicdNotificationHdl interface {
	ProcessNotification(msg AsicdNotifyMsg)
}

// Empty Interface
type AsicdNotifyMsg interface {
}

type AsicdClientStruct struct {
	Logger *logging.Writer
	NHdl   AsicdNotificationHdl
	NMap   AsicdNotification
}

type IPv6IntfState struct {
	IntfRef           string
	IfIndex           int32
	IpAddr            string
	OperState         string
	NumUpEvents       int32
	LastUpEventTime   string
	NumDownEvents     int32
	LastDownEventTime string
	L2IntfType        string
	L2IntfId          int32
}

type SubIPv4IntfState struct {
	IntfRef       string
	Type          string
	IfIndex       int32
	IfName        string
	ParentIfIndex int32
	IpAddr        string
	MacAddr       string
	OperState     string
}

type SubIPv6IntfState struct {
	IntfRef       string
	Type          string
	IfIndex       int32
	IfName        string
	ParentIfIndex int32
	IpAddr        string
	MacAddr       string
	OperState     string
}
