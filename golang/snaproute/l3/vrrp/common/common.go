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
package common

const (
	_ = iota
	_ // skipping value 1
	VERSION2
	VERSION3
)

const (
	USE_CONFIG_ADVERTISEMENT = -100
	STATE_UP                 = "UP"
	STATE_DOWN               = "DOWN"
	IP_MSG_CREATE            = "ip_create"
	IP_MSG_DELETE            = "ip_delete"
	IP_MSG_STATE_CHANGE      = "state_change"
	NETMASK_DELIMITER        = "/"
	SLASH_32                 = "32"
	SLASH_64                 = "64"
	VERSION2_STR             = "version2"
	VERSION3_STR             = "version3"
)

const (
	_ = iota
	CREATE
	UPDATE
	DELETE
	DB_UPDATE
)

type GlobalConfig struct {
	Vrf       string
	Enable    bool
	Operation uint8
}

type PhyPort struct {
	IfIndex   int32
	IntfRef   string
	OperState string
	MacAddr   string
}

type VlanInfo struct {
	Id            int32
	IfIndex       int32
	Name          string
	UntagPortsMap map[int32]bool
	TagPortsMap   map[int32]bool
	OperState     string
}

type BaseIpInfo struct {
	IfIndex   int32
	IntfRef   string
	IpAddr    string
	OperState string
	MsgType   string
	IpType    int // useful information for server
}

type Ipv4Info struct {
	Info BaseIpInfo
}

type Ipv6Info struct {
	Info          BaseIpInfo
	GlobalScopeIp string
}

type IntfCfg struct {
	IntfRef               string
	IfIndex               int32
	VRID                  int32
	Priority              int32
	VirtualIPAddr         string
	AdvertisementInterval int32
	PreemptMode           bool
	AcceptMode            bool
	AdminState            bool
	Version               uint8 // Information that will be used by server.. as all configs will be passed onto one channel only
	Operation             uint8 // Information that will be used by server
	IpType                int   // Information that will be used by server
}

type State struct {
	IntfRef                 string // interface where vrrp is configured
	Vrid                    int32  // virtual router id
	OperState               string // Informs whether vrrp is up or down
	CurrentFsmState         string // current fsm state
	MasterIp                string // Remote Master Ip Address
	AdverRx                 uint32 // Total advertisement received
	AdverTx                 uint32 // Total advertisement send out
	LastAdverRx             string // Last advertisement received
	LastAdverTx             string // Last advertisment send out
	IpAddr                  string // l3 interface ip address
	VirtualIp               string // router ip address
	VirtualRouterMACAddress string // route mac address
	AdvertisementInterval   int32  // time when master will send out advertisements
	MasterDownTimer         int32  // time to declare master is down
}

type VirtualIpInfo struct {
	IntfRef string
	IpAddr  string
	MacAddr string
	Enable  bool
	Version uint8
	IpType  int
}

type GlobalState struct {
	Vrf           string
	Status        bool
	V4Intfs       int32
	V6Intfs       int32
	TotalRxFrames int32
	TotalTxFrames int32
}
