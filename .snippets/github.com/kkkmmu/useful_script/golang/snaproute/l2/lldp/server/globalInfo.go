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
	"l2/lldp/config"
	"l2/lldp/packet"
	"l2/lldp/plugin"
	"sync"
	"time"
	"utils/dbutils"
)

type InPktChannel struct {
	pkt     gopacket.Packet
	ifIndex int32
}

type SendPktChannel struct {
	ifIndex int32
}

type LLDPGlobalInfo struct {
	// Port information
	Port config.PortInfo
	// Pcap Handler for Each Port
	PcapHandle *pcap.Handle
	// rx information
	RxInfo *packet.RX
	// tx information
	TxInfo *packet.TX
	// State info
	enable bool
	// Reading received info & updating received info lock
	RxLock *sync.RWMutex
	// Go Routine Killer Channels
	RxKill chan bool
	TxDone chan bool
	// counter for total frames rx/tx
	counter Frame
	// last received packet time
	pktRcvdTime time.Time
	rxtxMode    uint8
}

type Frame struct {
	Send int32
	Rcvd int32
}

type LLDPServer struct {
	// Basic server start fields
	lldpDbHdl *dbutils.DBUtil
	paramsDir string

	asicPlugin plugin.AsicIntf
	CfgPlugin  plugin.ConfigIntf
	SysPlugin  plugin.SystemIntf

	//System Information
	SysInfo *config.SystemInfo

	// Global LLDP Information
	Global *config.Global

	// lldp per port global info
	lldpGblInfo            map[int32]LLDPGlobalInfo
	lldpIntfStateSlice     []int32
	lldpUpIntfStateSlice   []int32
	lldpIntfRef2IfIndexMap map[string]int32

	// lldp pcap handler default config values
	lldpSnapshotLen int32
	lldpPromiscuous bool
	lldpTimeout     time.Duration

	// lldp packet rx channel
	lldpRxPktCh chan InPktChannel
	// lldp send packet channel
	lldpTxPktCh chan SendPktChannel
	// lldp global config channel
	GblCfgCh chan *config.Global
	// lldp per port config
	IntfCfgCh chan *config.IntfConfig
	// lldp asic notification channel
	IfStateCh chan *config.PortState
	// Update Cache notification channel
	UpdateCacheCh chan *config.SystemInfo
	// Event Publish channel for server
	EventCh chan config.EventInfo
	// Port Attribute change channel
	PortAttrCh chan *config.PortAttrMsg

	// Frames Counter
	counter Frame
}

const (
	// LLDP profiling
	LLDP_CPU_PROFILE_FILE = "/var/log/lldp.prof"

	// Consts Init Size/Capacity
	LLDP_INITIAL_GLOBAL_INFO_CAPACITY   = 100
	LLDP_RX_PKT_CHANNEL_SIZE            = 30
	LLDP_TX_PKT_CHANNEL_SIZE            = 30
	LLDP_PORT_STATE_CHANGE_CHANNEL_SIZE = 200
	LLDP_PORT_CONFIG_CHANNEL_SIZE       = 5

	// Port Operation State
	LLDP_PORT_STATE_DOWN = "DOWN"
	LLDP_PORT_STATE_UP   = "UP"
	LLDP_PORT_BROKEN_OUT = "Port broken out"

	LLDP_BPF_FILTER                 = "ether proto 0x88cc"
	LLDP_DEFAULT_TX_INTERVAL        = 30
	LLDP_DEFAULT_TX_HOLD_MULTIPLIER = 4
	LLDP_MIN_FRAME_LENGTH           = 12 // this is 12 bytes
	LLDP_FAST_LEARN_TIMER           = 1  // in seconds
	LLDP_FAST_LEARN_MAX_FRAMES_SEND = 5
)
