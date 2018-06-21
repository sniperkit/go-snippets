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

import (
	"models/objects"
)

type AsicGlobalStateGetInfo struct {
	StartIdx            int32
	EndIdx              int32
	Count               int32
	More                bool
	AsicGlobalStateList []*objects.AsicGlobalState
}

type AsicGlobalPMGetInfo struct {
	StartIdx         int32
	EndIdx           int32
	Count            int32
	More             bool
	AsicGlobalPMList []*objects.AsicGlobalPM
}

type AsicGlobalPMStateGetInfo struct {
	StartIdx              int32
	EndIdx                int32
	Count                 int32
	More                  bool
	AsicGlobalPMStateList []*objects.AsicGlobalPMState
}

type EthernetPMGetInfo struct {
	StartIdx       int32
	EndIdx         int32
	Count          int32
	More           bool
	EthernetPMList []*objects.EthernetPM
}

type EthernetPMStateGetInfo struct {
	StartIdx            int32
	EndIdx              int32
	Count               int32
	More                bool
	EthernetPMStateList []*objects.EthernetPMState
}

type AsicSummaryStateGetInfo struct {
	StartIdx             int32
	EndIdx               int32
	Count                int32
	More                 bool
	AsicSummaryStateList []*objects.AsicSummaryState
}

type VlanStateGetInfo struct {
	StartIdx      int32
	EndIdx        int32
	Count         int32
	More          bool
	VlanStateList []*objects.VlanState
}

type IPv4IntfStateGetInfo struct {
	StartIdx          int32
	EndIdx            int32
	Count             int32
	More              bool
	IPv4IntfStateList []*objects.IPv4IntfState
}

type PortGetInfo struct {
	StartIdx int32
	EndIdx   int32
	Count    int32
	More     bool
	PortList []*objects.Port
}

type PortStateGetInfo struct {
	StartIdx      int32
	EndIdx        int32
	Count         int32
	More          bool
	PortStateList []*objects.PortState
}

type MacTableEntryStateGetInfo struct {
	StartIdx               int32
	EndIdx                 int32
	Count                  int32
	More                   bool
	MacTableEntryStateList []*objects.MacTableEntryState
}

type IPv4RouteHwStateGetInfo struct {
	StartIdx             int32
	EndIdx               int32
	Count                int32
	More                 bool
	IPv4RouteHwStateList []*objects.IPv4RouteHwState
}

type IPv6RouteHwStateGetInfo struct {
	StartIdx             int32
	EndIdx               int32
	Count                int32
	More                 bool
	IPv6RouteHwStateList []*objects.IPv6RouteHwState
}

type ArpEntryHwStateGetInfo struct {
	StartIdx            int32
	EndIdx              int32
	Count               int32
	More                bool
	ArpEntryHwStateList []*objects.ArpEntryHwState
}

type NdpEntryHwStateGetInfo struct {
	StartIdx            int32
	EndIdx              int32
	Count               int32
	More                bool
	NdpEntryHwStateList []*objects.NdpEntryHwState
}

type LogicalIntfStateGetInfo struct {
	StartIdx             int32
	EndIdx               int32
	Count                int32
	More                 bool
	LogicalIntfStateList []*objects.LogicalIntfState
}

type SubIPv4IntfStateGetInfo struct {
	StartIdx             int32
	EndIdx               int32
	Count                int32
	More                 bool
	SubIPv4IntfStateList []*objects.SubIPv4IntfState
}

type IPv6IntfStateGetInfo struct {
	StartIdx          int32
	EndIdx            int32
	Count             int32
	More              bool
	IPv6IntfStateList []*objects.IPv6IntfState
}

type SubIPv6IntfStateGetInfo struct {
	StartIdx             int32
	EndIdx               int32
	Count                int32
	More                 bool
	SubIPv6IntfStateList []*objects.SubIPv6IntfState
}

type BufferPortStatStateGetInfo struct {
	StartIdx                int32
	EndIdx                  int32
	Count                   int32
	More                    bool
	BufferPortStatStateList []*objects.BufferPortStatState
}

type BufferGlobalStatStateGetInfo struct {
	StartIdx                  int32
	EndIdx                    int32
	Count                     int32
	More                      bool
	BufferGlobalStatStateList []*objects.BufferGlobalStatState
}

type AclStateGetInfo struct {
	StartIdx     int32
	EndIdx       int32
	Count        int32
	More         bool
	AclStateList []*objects.AclState
}

type LinkScopeIpStateGetInfo struct {
	StartIdx             int32
	EndIdx               int32
	Count                int32
	More                 bool
	LinkScopeIpStateList []*objects.LinkScopeIpState
}

type CoppStatStateGetInfo struct {
	StartIdx          int32
	EndIdx            int32
	Count             int32
	More              bool
	CoppStatStateList []*objects.CoppStatState
}

/* Internal Vlan Object */
type Vlan struct {
	VlanId           int32
	IfIndexList      []int32
	UntagIfIndexList []int32
	VlanName         string
}

type VlanGetInfo struct {
	StartIdx int32
	EndIdx   int32
	Count    int32
	More     bool
	VlanList []*Vlan
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
	LagList  []*Lag
}

type Port struct {
	IntfRef        string
	IfIndex        int32
	Description    string
	PhyIntfType    string
	AdminState     string
	MacAddr        string
	Speed          int32
	Duplex         string
	Autoneg        string
	MediaType      string
	Mtu            int32
	BreakOutMode   string
	LoopbackMode   string
	EnableFEC      bool
	PRBSTxEnable   bool
	PRBSRxEnable   bool
	PRBSPolynomial string
}

type IPv4NextHop struct {
	NextHopIp     string
	Weight        int32
	NextHopIfType int32
}

type IPv6NextHop struct {
	NextHopIp     string
	Weight        int32
	NextHopIfType int32
}

type IPv4Route struct {
	DestinationNw string
	NetworkMask   string
	NextHopList   []*IPv4NextHop
}

type IPv6Route struct {
	DestinationNw string
	NetworkMask   string
	NextHopList   []*IPv6NextHop
}
