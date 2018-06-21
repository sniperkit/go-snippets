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

package objects

const (
	AUTH_TYPE_NONE_STR            string = "none"
	AUTH_TYPE_SIMPLE_PASSWORD_STR string = "simplepassword"
	AUTH_TYPE_MD5_STR             string = "md5"
)

const (
	AUTH_TYPE_NONE            uint8 = 0
	AUTH_TYPE_SIMPLE_PASSWORD uint8 = 1
	AUTH_TYPE_MD5             uint8 = 2
)

const (
	AREA_ADMIN_STATE_UP   bool = true
	AREA_ADMIN_STATE_DOWN bool = false
)

const (
	AREA_ADMIN_STATE_UP_STR   string = "up"
	AREA_ADMIN_STATE_DOWN_STR string = "down"
)

const (
	OSPFV2_AREA_UPDATE_ADMIN_STATE      = 0x1
	OSPFV2_AREA_UPDATE_AUTH_TYPE        = 0x2
	OSPFV2_AREA_UPDATE_IMPORT_AS_EXTERN = 0x3
)

type Ospfv2Area struct {
	AreaId         uint32
	AdminState     bool
	AuthType       uint8
	ImportASExtern bool
}

type Ospfv2AreaState struct {
	AreaId uint32
	//NumSpfRuns       uint32
	//NumBdrRtr        uint32
	//NumAsBdrRtr      uint32
	NumOfRouterLSA     uint32
	NumOfNetworkLSA    uint32
	NumOfSummary3LSA   uint32
	NumOfSummary4LSA   uint32
	NumOfASExternalLSA uint32
	NumOfIntfs         uint32
	NumOfLSA           uint32
	NumOfNbrs          uint32
	NumOfRoutes        uint32
}

type Ospfv2AreaStateGetInfo struct {
	EndIdx int
	Count  int
	More   bool
	List   []*Ospfv2AreaState
}

const (
	GLOBAL_ADMIN_STATE_UP   bool = true
	GLOBAL_ADMIN_STATE_DOWN bool = false
)

const (
	GLOBAL_ADMIN_STATE_UP_STR   string = "up"
	GLOBAL_ADMIN_STATE_DOWN_STR string = "down"
)

const (
	OSPFV2_GLOBAL_UPDATE_ROUTER_ID           = 0x1
	OSPFV2_GLOBAL_UPDATE_ADMIN_STATE         = 0x2
	OSPFV2_GLOBAL_UPDATE_AS_BDR_RTR_STATUS   = 0x4
	OSPFV2_GLOBAL_UPDATE_REFERENCE_BANDWIDTH = 0x8
)

type Ospfv2Global struct {
	Vrf                string
	RouterId           uint32
	AdminState         bool
	ASBdrRtrStatus     bool
	ReferenceBandwidth uint32
}

type Ospfv2GlobalState struct {
	Vrf                string
	AreaBdrRtrStatus   bool
	NumOfAreas         uint32
	NumOfIntfs         uint32
	NumOfNbrs          uint32
	NumOfLSA           uint32
	NumOfRouterLSA     uint32
	NumOfNetworkLSA    uint32
	NumOfSummary3LSA   uint32
	NumOfSummary4LSA   uint32
	NumOfASExternalLSA uint32
	NumOfRoutes        uint32
}

type Ospfv2GlobalStateGetInfo struct {
	EndIdx int
	Count  int
	More   bool
	List   []*Ospfv2GlobalState
}

const (
	INTF_ADMIN_STATE_DOWN bool = false
	INTF_ADMIN_STATE_UP   bool = true
)

const (
	INTF_ADMIN_STATE_DOWN_STR string = "down"
	INTF_ADMIN_STATE_UP_STR   string = "up"
)

const (
	INTF_TYPE_POINT2POINT_STR string = "pointtopoint"
	INTF_TYPE_BROADCAST_STR   string = "broadcast"
)

const (
	INTF_TYPE_POINT2POINT uint8 = 0
	INTF_TYPE_BROADCAST   uint8 = 1
)

const (
	INTF_FSM_STATE_UNKNOWN  uint8 = 0
	INTF_FSM_STATE_DOWN     uint8 = 1
	INTF_FSM_STATE_WAITING  uint8 = 2
	INTF_FSM_STATE_LOOPBACK uint8 = 3
	INTF_FSM_STATE_P2P      uint8 = 4
	INTF_FSM_STATE_OTHER_DR uint8 = 5
	INTF_FSM_STATE_DR       uint8 = 6
	INTF_FSM_STATE_BDR      uint8 = 7
)

const (
	INTF_FSM_STATE_UNKNOWN_STR  string = "unknown"
	INTF_FSM_STATE_DOWN_STR     string = "down"
	INTF_FSM_STATE_WAITING_STR  string = "waiting"
	INTF_FSM_STATE_LOOPBACK_STR string = "loopback"
	INTF_FSM_STATE_P2P_STR      string = "point-to-point"
	INTF_FSM_STATE_OTHER_DR_STR string = "other-dr"
	INTF_FSM_STATE_DR_STR       string = "dr"
	INTF_FSM_STATE_BDR_STR      string = "bdr"
)

const (
	OSPFV2_INTF_UPDATE_ADMIN_STATE       = 0x2
	OSPFV2_INTF_UPDATE_AREA_ID           = 0x4
	OSPFV2_INTF_UPDATE_TYPE              = 0x8
	OSPFV2_INTF_UPDATE_RTR_PRIORITY      = 0x10
	OSPFV2_INTF_UPDATE_TRANSIT_DELAY     = 0x20
	OSPFV2_INTF_UPDATE_RETRANS_INTERVAL  = 0x40
	OSPFV2_INTF_UPDATE_HELLO_INTERVAL    = 0x80
	OSPFV2_INTF_UPDATE_RTR_DEAD_INTERVAL = 0x100
	OSPFV2_INTF_UPDATE_METRIC_VALUE      = 0x200
)

type Ospfv2Intf struct {
	IpAddress        uint32
	AddressLessIfIdx uint32
	AdminState       bool
	AreaId           uint32
	Type             uint8
	RtrPriority      uint8
	TransitDelay     uint16
	RetransInterval  uint16
	HelloInterval    uint16
	RtrDeadInterval  uint32
	MetricValue      uint16
}

type Ospfv2IntfState struct {
	IpAddress                uint32
	AddressLessIfIdx         uint32
	State                    uint8
	DesignatedRouter         uint32
	DesignatedRouterId       uint32
	BackupDesignatedRouter   uint32
	BackupDesignatedRouterId uint32
	NumOfRouterLSA           uint32
	NumOfNetworkLSA          uint32
	NumOfSummary3LSA         uint32
	NumOfSummary4LSA         uint32
	NumOfASExternalLSA       uint32
	NumOfLSA                 uint32
	NumOfNbrs                uint32
	NumOfRoutes              uint32
	Mtu                      uint32
	Cost                     uint32
	NumOfStateChange         uint32
	TimeOfStateChange        string
}
type Ospfv2IntfStateGetInfo struct {
	EndIdx int
	Count  int
	More   bool
	List   []*Ospfv2IntfState
}

const (
	ROUTER_LSA     uint8 = 1
	NETWORK_LSA    uint8 = 2
	SUMMARY3_LSA   uint8 = 3
	SUMMARY4_LSA   uint8 = 4
	ASExternal_LSA uint8 = 5
)

const (
	ROUTER_LSA_STR     string = "router"
	NETWORK_LSA_STR    string = "network"
	SUMMARY3_LSA_STR   string = "summary3"
	SUMMARY4_LSA_STR   string = "summary4"
	ASExternal_LSA_STR string = "asexternal"
)

type Ospfv2LsdbState struct {
	LSType        uint8
	LSId          uint32
	AreaId        uint32
	AdvRouterId   uint32
	SequenceNum   uint32
	Age           uint16
	Checksum      uint16
	Options       uint8
	Length        uint16
	Advertisement string
}

type Ospfv2LsdbStateGetInfo struct {
	EndIdx int
	Count  int
	More   bool
	List   []*Ospfv2LsdbState
}

const (
	NBR_STATE_ONE_WAY_STR  string = "oneway"
	NBR_STATE_TWO_WAY_STR  string = "twoway"
	NBR_STATE_INIT_STR     string = "init"
	NBR_STATE_EXSTART_STR  string = "exstart"
	NBR_STATE_EXCHANGE_STR string = "exchange"
	NBR_STATE_LOADING_STR  string = "loading"
	NBR_STATE_ATTEMPT_STR  string = "attempt"
	NBR_STATE_DOWN_STR     string = "down"
	NBR_STATE_FULL_STR     string = "full"
)

const (
	NBR_STATE_ONE_WAY  uint8 = 0
	NBR_STATE_TWO_WAY  uint8 = 1
	NBR_STATE_INIT     uint8 = 2
	NBR_STATE_EXSTART  uint8 = 3
	NBR_STATE_EXCHANGE uint8 = 4
	NBR_STATE_LOADING  uint8 = 5
	NBR_STATE_ATTEMPT  uint8 = 6
	NBR_STATE_DOWN     uint8 = 7
	NBR_STATE_FULL     uint8 = 8
)

type Ospfv2NbrState struct {
	IpAddr           uint32
	AddressLessIfIdx uint32
	RtrId            uint32
	Options          int32
	State            uint8
}

type Ospfv2NbrStateGetInfo struct {
	EndIdx int
	Count  int
	More   bool
	List   []*Ospfv2NbrState
}

/*
type Ospfv2NextHop struct {
	IntfIPAddr    uint32
	IntfIfIdx     uint32
	NextHopIPAddr uint32
	AdvRtrId      uint32
}

type Ospfv2RouteState struct {
	DestId            uint32
	AddrMask          uint32
	DestType          uint8
	OptCapabilities   int32
	AreaId            uint32
	PathType          uint8
	Cost              uint32
	Type2Cost         uint32
	NumOfPaths        uint16
	LSOriginLSType    uint8
	LSOriginLSId      uint32
	LSOriginAdvRouter uint32
	NextHops          []Ospfv2NextHop
}

type Ospfv2RouteStateGetInfo struct {
	EndIdx int
	Count  int
	More   bool
	List   []*Ospfv2RouteState
}
*/
