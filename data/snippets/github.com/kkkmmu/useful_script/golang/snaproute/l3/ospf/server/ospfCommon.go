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
	"github.com/google/gopacket/pcap"
	"net"
	"sync"
	"time"
)

var ALLSPFROUTER string = "224.0.0.5"
var ALLDROUTER string = "224.0.0.6"
var ALLSPFROUTERMAC string = "01:00:5e:00:00:05"
var ALLDROUTERMAC string = "01:00:5e:00:00:06"
var MASKMAC string = "ff:ff:ff:ff:ff:ff"

const (
	DEFAULT_VLAN_COST uint32 = 10
)

var LSInfinity uint32 = 0x00ffffff

type OspfHdrMetadata struct {
	pktType  OspfType
	pktlen   uint16
	backbone bool
	routerId []byte
	areaId   uint32
}

func NewOspfHdrMetadata() *OspfHdrMetadata {
	return &OspfHdrMetadata{}
}

type DstIPType uint8

const (
	Normal       DstIPType = 1
	AllSPFRouter DstIPType = 2
	AllDRouter   DstIPType = 3
)

type IpHdrMetadata struct {
	srcIP     []byte
	dstIP     []byte
	dstIPType DstIPType
}

func NewIpHdrMetadata() *IpHdrMetadata {
	return &IpHdrMetadata{}
}

type EthHdrMetadata struct {
	srcMAC net.HardwareAddr
}

func NewEthHdrMetadata() *EthHdrMetadata {
	return &EthHdrMetadata{}
}

var (
	snapshot_len int32         = 65549 //packet capture length
	promiscuous  bool          = false //mode
	timeout_pcap time.Duration = 5 * time.Second
)

const (
	OSPF_HELLO_MIN_SIZE  = 20
	OSPF_DBD_MIN_SIZE    = 8
	OSPF_LSA_HEADER_SIZE = 20
	OSPF_LSA_REQ_SIZE    = 12
	OSPF_LSA_ACK_SIZE    = 20
	OSPF_HEADER_SIZE     = 24
	IP_HEADER_MIN_LEN    = 20
	OSPF_PROTO_ID        = 89
	OSPF_VERSION_2       = 2
	OSPF_NO_OF_LSA_FIELD = 4
)

type OspfType uint8

const (
	HelloType         OspfType = 1
	DBDescriptionType OspfType = 2
	LSRequestType     OspfType = 3
	LSUpdateType      OspfType = 4
	LSAckType         OspfType = 5
)

type IntfToNeighMsg struct {
	IntfConfKey  IntfConfKey
	RouterId     uint32
	RtrPrio      uint8
	NeighborIP   net.IP
	nbrDeadTimer time.Duration
	TwoWayStatus bool
	nbrDR        []byte
	nbrBDR       []byte
	nbrMAC       net.HardwareAddr
}

type NbrStateChangeMsg struct {
	nbrKey NeighborConfKey
}

const (
	EOption  = 0x02
	MCOption = 0x04
	NPOption = 0x08
	EAOption = 0x20
	DCOption = 0x40
)

type IntfTxHandle struct {
	SendPcapHdl *pcap.Handle
	SendMutex   *sync.Mutex
}

type IntfRxHandle struct {
	RecvPcapHdl     *pcap.Handle
	PktRecvCh       chan bool
	PktRecvStatusCh chan bool
	//RecvMutex               *sync.Mutex
}

type AdjOKEvtMsg struct {
	NewDRtrId  uint32
	OldDRtrId  uint32
	NewBDRtrId uint32
	OldBDRtrId uint32
}

type NbrFullStateMsg struct {
	FullState bool
	NbrRtrId  uint32
	nbrKey    NeighborConfKey
}

const (
	LsdbEntryFound    = 0
	LsdbEntryNotFound = 1
)
