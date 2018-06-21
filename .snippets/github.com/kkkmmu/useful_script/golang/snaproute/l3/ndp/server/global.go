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
package server

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"l3/ndp/config"
	"l3/ndp/packet"
	"sync"
	"time"
	"utils/clntUtils/clntIntfs/asicdClntIntfs"
	"utils/dmnBase"
)

type RxPktInfo struct {
	pkt     gopacket.Packet
	ifIndex int32
}

type NdpConfig struct {
	Vrf               string
	ReachableTime     uint32
	RetransTime       uint32
	RaRestransmitTime uint8
}

type L3Info struct {
	Name     string
	IfIndex  int32
	PortType string // tag or untag
	Updated  bool   // if set then do not delete during update
}

type PhyPort struct {
	RX   *pcap.Handle
	Info config.PortInfo
}

type NDPServer struct {
	NdpConfig                                              // base config
	dmnBase               *dmnBase.FSBaseDmn               // base Daemon
	SwitchPlugin          asicdClntIntfs.AsicdClntIntf     // asicd plugin
	L2Port                map[int32]PhyPort                // key is l2 ifIndex
	L3Port                map[int32]Interface              // key is l3 ifIndex
	VlanInfo              map[int32]config.VlanInfo        // key is vlanIfIndex
	virtualIp             map[int32]*config.VirtualIpInfo  // key is virtual ip ifIndex
	VlanIfIdxVlanIdMap    map[string]int32                 // reverse map for vlanName ----> vlanId, used during ipv6 neig create
	SwitchMacMapEntries   map[string]struct{}              // cache entry for all mac addresses on a switch
	NeighborInfo          map[string]config.NeighborConfig // neighbor created by NDP used for STATE
	neighborKey           []string                         // keys for all neighbor entries is stored here for GET calls
	PhyPortToL3PortMap    map[int32]L3Info                 // reverse map for l2IfIndex ----> l3IfIndex, used during vlan RX Pcap
	Dot1QToVlanIfIndex    map[int32]int32                  // reverse map of vlanId (aka dot1q tag) to l3IfIndex for vlan
	GlobalCfg             chan NdpConfig                   // Configuration Channels
	NeigborEntryLock      *sync.RWMutex                    // Lock for reading/writing NeighorInfo
	IpIntfCh              chan *config.IPIntfNotification  // IPV6 Create/Delete State Up/Down Notification Channel
	VlanCh                chan *config.VlanNotification    // Vlan Create/Delete/Update Notification Channel
	MacMoveCh             chan *config.MacMoveNotification // Mac Move Notification Channel
	RxPktCh               chan *RxPktInfo                  // Received Pkt Channel
	PktDataCh             chan config.PacketData           // Package packet informs server over PktDataCh saying that send this packet..
	ActionCh              chan *config.ActionData          // Action Channel for NDP
	VirtualIpCh           chan *config.VirtualIpInfo       // Virtual Ip Channel used for virtual interface notification from asicd
	Packet                *packet.Packet                   // Neighbor Cache Information
	SwitchMac             string                           // @HACK: Need to find better way of getting Switch Mac Address
	notifyChan            chan<- []byte                    // Notification Channel for Publisher
	counter               PktCounter                       // counter for packets send and received
	ndpL3IntfStateSlice   []int32
	ndpUpL3IntfStateSlice []int32
	L3IfIntfRefToIfIndex  map[string]int32
	SnapShotLen           int32
	Promiscuous           bool
	Timeout               time.Duration
}

const (
	NDP_CPU_PROFILE_FILE                  = "/var/log/ndpd.prof"
	NDP_SERVER_MAP_INITIAL_CAP            = 30
	NDP_SERVER_ASICD_NOTIFICATION_CH_SIZE = 1
	NDP_SERVER_INITIAL_CHANNEL_SIZE       = 1
	INTF_REF_NOT_FOUND                    = "Not Found"
)
