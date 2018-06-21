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

package server

import (
	"time"
)

const (
	snapshotLen             int32         = 65549 //packet capture length
	promiscuous             bool          = false //mode
	pcapTimeout             time.Duration = 5 * time.Second
	LOGICAL_INTF_MTU        int32         = 1512
	ALLSPFROUTER            uint32        = 0xE0000005
	ALLDROUTER              uint32        = 0xE0000006
	ALLSPFROUTERMAC         string        = "01:00:5e:00:00:05"
	ALLDROUTERMAC           string        = "01:00:5e:00:00:06"
	MASKMAC                 string        = "ff:ff:ff:ff:ff:ff"
	InitialSequenceNum      uint32        = 0x80000001
	MaxSequenceNumber       uint32        = 0x7fffffff
	LsaAgingTimeGranularity time.Duration = 1 * time.Second
	MAX_AGE                 uint16        = 3600 // 1 hour
	LS_REFRESH_TIME         uint16        = 1800 // 30 min
	MIN_LS_INTERVAL         uint16        = 5    // 5 second
	MIN_LS_ARRIVAL          uint16        = 1    // 1 second
	CHECK_AGE               uint16        = 300  // 5 mins
	MAX_AGE_DIFF            uint16        = 900  // 15 mins
	//LSSequenceNumber      int           = InitialSequenceNumber
	LSInfinity                 uint32 = 0x00ffffff
	FLETCHER_CHECKSUM_VALIDATE uint16 = 0xffff
)

const (
	OSPF_HELLO_MIN_SIZE        = 20
	OSPF_DBD_MIN_SIZE          = 8
	OSPF_LSA_HEADER_SIZE       = 20
	OSPF_LSA_REQ_SIZE          = 12
	OSPF_LSA_ACK_SIZE          = 20
	OSPF_HEADER_SIZE           = 24
	IP_HEADER_MIN_LEN          = 20
	OSPF_PROTO_ID        uint8 = 89
	OSPF_VERSION_2       uint8 = 2
	OSPF_NO_OF_LSA_FIELD       = 4
)

const (
	HelloType         uint8 = 1
	DBDescriptionType uint8 = 2
	LSRequestType     uint8 = 3
	LSUpdateType      uint8 = 4
	LSAckType         uint8 = 5
)

const (
	P2P_LINK     uint8 = 1
	TRANSIT_LINK uint8 = 2
	STUB_LINK    uint8 = 3
	VIRTUAL_LINK uint8 = 4
)

type DstIPType uint8

const (
	NormalType       DstIPType = 1
	AllSPFRouterType DstIPType = 2
	AllDRouterType   DstIPType = 3
)

const (
	EOption  = 0x02
	MCOption = 0x04
	NPOption = 0x08
	EAOption = 0x20
	DCOption = 0x40
)

const (
	LsdbEntryFound    bool = true
	LsdbEntryNotFound bool = false
)

const (
	LsdbAdd      uint8 = 0
	LsdbDel      uint8 = 1
	LsdbUpdate   uint8 = 2
	LsdbNoAction uint8 = 3
)

const (
	LSA_MAX_AGE      uint16 = 0x7fff
	LSA_MAX_AGE_DIFF uint16 = 0x7fff
	LSASELFLOOD             = 1 // flood for received LSA
	LSAINTF                 = 2 // Send LSA on the interface in reply to LSAREQ
	LSAAGE                  = 3 // flood aged LSAs.
	LSASUMMARYFLOOD         = 4 //flood summary LSAs in different areas.
	LSAEXTFLOOD             = 5 //flood AS External summary LSA
	LSAROUTERFLOOD          = 6 //flood only router LSA
)

const (
	AllSPFRouters = "224.0.0.5"
	AllDRouters   = "224.0.0.6"
	McastMAC      = "01:00:5e:00:00:05"
)
