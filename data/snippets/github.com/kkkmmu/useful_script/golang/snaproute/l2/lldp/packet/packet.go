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

package packet

import (
	"github.com/google/gopacket/layers"
	"net"
	"time"
)

const (
	LLDP_MAX_TTL             = 65535
	LLDP_PROTO_DST_MAC       = "01:80:c2:00:00:0e"
	LLDP_TOTAL_TLV_SUPPORTED = 8
)

type RX struct {
	RxRunning bool
	// ethernet frame Info (used for rx/tx)
	SrcMAC net.HardwareAddr // NOTE: Please be informed this is Peer Mac Addr
	DstMAC net.HardwareAddr

	// lldp rx information
	RxFrame         *layers.LinkLayerDiscovery
	RxLinkInfo      *layers.LinkLayerDiscoveryInfo
	OrgTLV802_1     layers.LLDPInfo8021
	ClearCacheTimer *time.Timer

	// cache last packet and see if we need to update current information or not
	LastPkt []byte
}

type TX struct {
	// tx information
	ttl                     int
	DstMAC                  net.HardwareAddr
	MessageTxInterval       int
	MessageTxHoldMultiplier int
	useCacheFrame           bool
	cacheFrame              []byte
	TxTimer                 *time.Timer
}
