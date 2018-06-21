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

package asicdClient

import (
	"utils/asicdClient/flexswitch"
	"utils/asicdClient/ovs"
	"utils/commonDefs"
)

type AsicdClientIntf interface {
	CreateIPv4Neighbor(ipAddr string, macAddr string, vlanId int32, ifIdx int32) (rv int32, err error)
	UpdateIPv4Neighbor(ipAddr string, macAddr string, vlanId int32, ifIdx int32) (rv int32, err error)
	DeleteIPv4Neighbor(ipAddr string) (rv int32, err error)

	CreateIPv6Neighbor(ipAddr string, macAddr string, vlanId int32, ifIdx int32) (rv int32, err error)
	UpdateIPv6Neighbor(ipAddr string, macAddr string, vlanId int32, ifIdx int32) (rv int32, err error)
	DeleteIPv6Neighbor(ipAddr string) (rv int32, err error)

	GetBulkIPv4IntfState(curMark, count int) (*commonDefs.IPv4IntfStateGetInfo, error)
	GetAllIPv4IntfState() ([]*commonDefs.IPv4IntfState, error)
	GetAllIPv6IntfState() ([]*commonDefs.IPv6IntfState, error)
	GetAllPortState() ([]*commonDefs.PortState, error)
	GetBulkPort(curMark, count int) (*commonDefs.PortGetInfo, error)
	GetBulkPortState(curMark, count int) (*commonDefs.PortStateGetInfo, error)
	GetBulkLag(curMark, count int) (*commonDefs.LagGetInfo, error)
	GetBulkVlan(curMark, count int) (*commonDefs.VlanGetInfo, error)
	GetBulkVlanState(curMark, count int) (*commonDefs.VlanStateGetInfo, error)
	GetAllVlanState() ([]*commonDefs.VlanState, error)
	GetAllVlan() ([]*commonDefs.Vlan, error)
	DetermineRouterId() string
	GetPort(string) (*commonDefs.Port, error)

	// get the switch MAC given in string format
	GetSwitchMAC(paramsPath string) string

	// Get the current link status of a link
	GetPortLinkStatus(pId int32) bool
	// create stp bridge, map vlans to stg, stgid returned by caller
	CreateStgBridge(vlanList []uint16) int32
	DeleteStgBridge(stgid int32, vlanList []uint16) error
	// set forwarding/learning/blocked state
	SetStgPortState(stgid int32, ifindex int32, state int) error
	// Flush the macs associated with this stg
	FlushStgFdb(stgid, ifindex int32) error
	// BPDU Guard detection
	BPDUGuardDetected(ifindex int32, enable bool) error

	CreateLag(ifname string, hashType int32, ports string) (int32, error)
	DeleteLag(ifIndex int32) error
	UpdateLag(ifIndex, hashType int32, ports string) error
	UpdateLagCfgIntfList(ifName string, ifIndexList []int32) bool

	EnablePacketReception(mac string, vlan int, ifindex int32) error
	DisablePacketReception(mac string, vlan int, ifindex int32) error

	// Distributed Relay (MLAG) Ipp actions when in Time Sharing Mode
	IppIngressEgressDrop(srcIfIndex, dstIfIndex string) error
	IppIngressEgressPass(srcIfIndex, dstIfIndex string) error
	IppVlanConversationSet(vlan uint16, ifindex int32) error
	IppVlanConversationClear(vlan uint16, ifindex int32) error

	IsLoopbackType(ifIndex int32) bool

	// Create Virtual Interface
	CreateVirtualIPv4Intf(intRef, ipAddr, macAddr string, enable bool) error
	CreateVirtualIPv6Intf(intRef, ipAddr, macAddr string, enable bool) error

	UpdateVirtualIPv4Intf(intfRef, ipAddr, macAddr string, enable bool) error
	UpdateVirtualIPv6Intf(intfRef, ipAddr, macAddr string, enable bool) error

	DeleteVirtualIPv4Intf(intRef, ipAddr, macAddr string, enable bool) error
	DeleteVirtualIPv6Intf(intRef, ipAddr, macAddr string, enable bool) error

	GetAllSubIPv4IntfState() ([]*commonDefs.SubIPv4IntfState, error)
	GetAllSubIPv6IntfState() ([]*commonDefs.SubIPv6IntfState, error)

	IsLinuxOnlyPlugin() (bool, error)
}

func NewAsicdClientInit(plugin string, paramsFile string, asicdHdl commonDefs.AsicdClientStruct) AsicdClientIntf {
	if plugin == "Flexswitch" {
		clientHdl := flexswitch.GetAsicdThriftClientHdl(paramsFile, asicdHdl.Logger)
		if clientHdl == nil {
			asicdHdl.Logger.Err("Unable Initialize Asicd Client")
			return nil
		}
		flexswitch.InitFSAsicdSubscriber(asicdHdl)
		return &flexswitch.FSAsicdClientMgr{clientHdl}
	} else if plugin == "OvsDB" {
		ovs.InitOvsAsicdSubscriber(asicdHdl)
		return &ovs.OvsAsicdClientMgr{100}
	}
	return nil
}
