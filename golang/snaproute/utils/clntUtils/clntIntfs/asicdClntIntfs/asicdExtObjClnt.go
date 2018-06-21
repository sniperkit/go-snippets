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

package asicdClntIntfs

import (
	"models/objects"
	"utils/clntUtils/clntDefs/asicdClntDefs"
)

type AsicdExtObjClntIntf interface {
	GetBulkAsicGlobalState(fromIndex int, count int) (*asicdClntDefs.AsicGlobalStateGetInfo, error)
	GetAsicGlobalState(ModuleId uint8) (*objects.AsicGlobalState, error)
	CreateAsicGlobalPM(cfg *objects.AsicGlobalPM) (bool, error)
	UpdateAsicGlobalPM(origCfg, newCfg *objects.AsicGlobalPM, attrset []bool, op []*objects.PatchOpInfo) (bool, error)
	DeleteAsicGlobalPM(cfg *objects.AsicGlobalPM) (bool, error)
	GetBulkAsicGlobalPM(fromIndex int, count int) (*asicdClntDefs.AsicGlobalPMGetInfo, error)
	GetAsicGlobalPM(ModuleId uint8, Resource string) (*objects.AsicGlobalPM, error)
	GetBulkAsicGlobalPMState(fromIndex int, count int) (*asicdClntDefs.AsicGlobalPMStateGetInfo, error)
	GetAsicGlobalPMState(ModuleId uint8, Resource string) (*objects.AsicGlobalPMState, error)
	CreateEthernetPM(cfg *objects.EthernetPM) (bool, error)
	UpdateEthernetPM(origCfg, newCfg *objects.EthernetPM, attrset []bool, op []*objects.PatchOpInfo) (bool, error)
	DeleteEthernetPM(cfg *objects.EthernetPM) (bool, error)
	GetBulkEthernetPM(fromIndex int, count int) (*asicdClntDefs.EthernetPMGetInfo, error)
	GetEthernetPM(IntfRef string, Resource string) (*objects.EthernetPM, error)
	GetBulkEthernetPMState(fromIndex int, count int) (*asicdClntDefs.EthernetPMStateGetInfo, error)
	GetEthernetPMState(IntfRef string, Resource string) (*objects.EthernetPMState, error)
	GetBulkAsicSummaryState(fromIndex int, count int) (*asicdClntDefs.AsicSummaryStateGetInfo, error)
	GetAsicSummaryState(ModuleId uint8) (*objects.AsicSummaryState, error)
	CreateVlan(cfg *objects.Vlan) (bool, error)
	UpdateVlan(origCfg, newCfg *objects.Vlan, attrset []bool, op []*objects.PatchOpInfo) (bool, error)
	DeleteVlan(cfg *objects.Vlan) (bool, error)
	GetBulkVlanState(fromIndex int, count int) (*asicdClntDefs.VlanStateGetInfo, error)
	GetVlanState(VlanId int32) (*objects.VlanState, error)
	CreateIPv4Intf(cfg *objects.IPv4Intf) (bool, error)
	UpdateIPv4Intf(origCfg, newCfg *objects.IPv4Intf, attrset []bool, op []*objects.PatchOpInfo) (bool, error)
	DeleteIPv4Intf(cfg *objects.IPv4Intf) (bool, error)
	GetBulkIPv4IntfState(fromIndex int, count int) (*asicdClntDefs.IPv4IntfStateGetInfo, error)
	GetIPv4IntfState(IntfRef string) (*objects.IPv4IntfState, error)
	CreatePort(cfg *objects.Port) (bool, error)
	UpdatePort(origCfg, newCfg *objects.Port, attrset []bool, op []*objects.PatchOpInfo) (bool, error)
	DeletePort(cfg *objects.Port) (bool, error)
	GetBulkPort(fromIndex int, count int) (*asicdClntDefs.PortGetInfo, error)
	GetPort(IntfRef string) (*objects.Port, error)
	GetBulkPortState(fromIndex int, count int) (*asicdClntDefs.PortStateGetInfo, error)
	GetPortState(IntfRef string) (*objects.PortState, error)
	GetBulkMacTableEntryState(fromIndex int, count int) (*asicdClntDefs.MacTableEntryStateGetInfo, error)
	GetMacTableEntryState(MacAddr string) (*objects.MacTableEntryState, error)
	GetBulkIPv4RouteHwState(fromIndex int, count int) (*asicdClntDefs.IPv4RouteHwStateGetInfo, error)
	GetIPv4RouteHwState(DestinationNw string) (*objects.IPv4RouteHwState, error)
	GetBulkIPv6RouteHwState(fromIndex int, count int) (*asicdClntDefs.IPv6RouteHwStateGetInfo, error)
	GetIPv6RouteHwState(DestinationNw string) (*objects.IPv6RouteHwState, error)
	GetBulkArpEntryHwState(fromIndex int, count int) (*asicdClntDefs.ArpEntryHwStateGetInfo, error)
	GetArpEntryHwState(IpAddr string) (*objects.ArpEntryHwState, error)
	GetBulkNdpEntryHwState(fromIndex int, count int) (*asicdClntDefs.NdpEntryHwStateGetInfo, error)
	GetNdpEntryHwState(IpAddr string) (*objects.NdpEntryHwState, error)
	CreateLogicalIntf(cfg *objects.LogicalIntf) (bool, error)
	UpdateLogicalIntf(origCfg, newCfg *objects.LogicalIntf, attrset []bool, op []*objects.PatchOpInfo) (bool, error)
	DeleteLogicalIntf(cfg *objects.LogicalIntf) (bool, error)
	GetBulkLogicalIntfState(fromIndex int, count int) (*asicdClntDefs.LogicalIntfStateGetInfo, error)
	GetLogicalIntfState(Name string) (*objects.LogicalIntfState, error)
	CreateSubIPv4Intf(cfg *objects.SubIPv4Intf) (bool, error)
	UpdateSubIPv4Intf(origCfg, newCfg *objects.SubIPv4Intf, attrset []bool, op []*objects.PatchOpInfo) (bool, error)
	DeleteSubIPv4Intf(cfg *objects.SubIPv4Intf) (bool, error)
	GetBulkSubIPv4IntfState(fromIndex int, count int) (*asicdClntDefs.SubIPv4IntfStateGetInfo, error)
	GetSubIPv4IntfState(IntfRef string, Type string) (*objects.SubIPv4IntfState, error)
	CreateIPv6Intf(cfg *objects.IPv6Intf) (bool, error)
	UpdateIPv6Intf(origCfg, newCfg *objects.IPv6Intf, attrset []bool, op []*objects.PatchOpInfo) (bool, error)
	DeleteIPv6Intf(cfg *objects.IPv6Intf) (bool, error)
	GetBulkIPv6IntfState(fromIndex int, count int) (*asicdClntDefs.IPv6IntfStateGetInfo, error)
	GetIPv6IntfState(IntfRef string) (*objects.IPv6IntfState, error)
	CreateSubIPv6Intf(cfg *objects.SubIPv6Intf) (bool, error)
	UpdateSubIPv6Intf(origCfg, newCfg *objects.SubIPv6Intf, attrset []bool, op []*objects.PatchOpInfo) (bool, error)
	DeleteSubIPv6Intf(cfg *objects.SubIPv6Intf) (bool, error)
	GetBulkSubIPv6IntfState(fromIndex int, count int) (*asicdClntDefs.SubIPv6IntfStateGetInfo, error)
	GetSubIPv6IntfState(IntfRef string, Type string) (*objects.SubIPv6IntfState, error)
	GetBulkBufferPortStatState(fromIndex int, count int) (*asicdClntDefs.BufferPortStatStateGetInfo, error)
	GetBufferPortStatState(IntfRef string) (*objects.BufferPortStatState, error)
	GetBulkBufferGlobalStatState(fromIndex int, count int) (*asicdClntDefs.BufferGlobalStatStateGetInfo, error)
	GetBufferGlobalStatState(DeviceId uint32) (*objects.BufferGlobalStatState, error)
	CreateAclGlobal(cfg *objects.AclGlobal) (bool, error)
	UpdateAclGlobal(origCfg, newCfg *objects.AclGlobal, attrset []bool, op []*objects.PatchOpInfo) (bool, error)
	DeleteAclGlobal(cfg *objects.AclGlobal) (bool, error)
	CreateAcl(cfg *objects.Acl) (bool, error)
	UpdateAcl(origCfg, newCfg *objects.Acl, attrset []bool, op []*objects.PatchOpInfo) (bool, error)
	DeleteAcl(cfg *objects.Acl) (bool, error)
	CreateAclIpv4Filter(cfg *objects.AclIpv4Filter) (bool, error)
	UpdateAclIpv4Filter(origCfg, newCfg *objects.AclIpv4Filter, attrset []bool, op []*objects.PatchOpInfo) (bool, error)
	DeleteAclIpv4Filter(cfg *objects.AclIpv4Filter) (bool, error)
	CreateAclMacFilter(cfg *objects.AclMacFilter) (bool, error)
	UpdateAclMacFilter(origCfg, newCfg *objects.AclMacFilter, attrset []bool, op []*objects.PatchOpInfo) (bool, error)
	DeleteAclMacFilter(cfg *objects.AclMacFilter) (bool, error)
	CreateAclIpv6Filter(cfg *objects.AclIpv6Filter) (bool, error)
	UpdateAclIpv6Filter(origCfg, newCfg *objects.AclIpv6Filter, attrset []bool, op []*objects.PatchOpInfo) (bool, error)
	DeleteAclIpv6Filter(cfg *objects.AclIpv6Filter) (bool, error)
	GetBulkAclState(fromIndex int, count int) (*asicdClntDefs.AclStateGetInfo, error)
	GetAclState(AclName string) (*objects.AclState, error)
	GetBulkLinkScopeIpState(fromIndex int, count int) (*asicdClntDefs.LinkScopeIpStateGetInfo, error)
	GetLinkScopeIpState(LinkScopeIp string) (*objects.LinkScopeIpState, error)
	GetBulkCoppStatState(fromIndex int, count int) (*asicdClntDefs.CoppStatStateGetInfo, error)
	GetCoppStatState(Protocol string) (*objects.CoppStatState, error)
}
