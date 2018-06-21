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

// This file defines all interfaces provided for the L3 service
package rpc

import (
	"asicdInt"
	"asicdServices"
	"errors"
)

//Logical Intf related services
func (svcHdlr AsicDaemonServiceHandler) CreateLogicalIntf(confObj *asicdServices.LogicalIntf) (rv bool, err error) {
	return true, nil
}
func (svcHdlr AsicDaemonServiceHandler) UpdateLogicalIntf(oldConfIntfObj, newConfIntfObj *asicdServices.LogicalIntf, attrset []bool, op []*asicdServices.PatchOpInfo) (bool, error) {
	return false, errors.New("Update operation is not supported on LogicalIntf object")
}
func (svcHdlr AsicDaemonServiceHandler) DeleteLogicalIntf(confObj *asicdServices.LogicalIntf) (rv bool, err error) {
	return true, nil
}
func (svcHdlr AsicDaemonServiceHandler) GetBulkLogicalIntfState(currMarker, count asicdServices.Int) (*asicdServices.LogicalIntfStateGetInfo, error) {
	bulkObj := asicdServices.NewLogicalIntfStateGetInfo()
	return bulkObj, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetLogicalIntfState(name string) (*asicdServices.LogicalIntfState, error) {
	logicalIntfState := asicdServices.NewLogicalIntfState()
	return logicalIntfState, nil
}

//IPv4 interface related services
func (svcHdlr AsicDaemonServiceHandler) CreateIPv4Intf(ipv4IntfObj *asicdServices.IPv4Intf) (rv bool, err error) {
	return true, nil
}
func (svcHdlr AsicDaemonServiceHandler) UpdateIPv4Intf(oldIPv4IntfObj, newIPv4IntfObj *asicdServices.IPv4Intf,
	attrset []bool, op []*asicdServices.PatchOpInfo) (bool, error) {
	return true, nil
}
func (svcHdlr AsicDaemonServiceHandler) DeleteIPv4Intf(ipv4IntfObj *asicdServices.IPv4Intf) (rv bool, err error) {
	return true, nil
}
func (svcHdlr AsicDaemonServiceHandler) GetIPv4IntfState(intfRef string) (*asicdServices.IPv4IntfState, error) {
	ipv4IntfState := asicdServices.NewIPv4IntfState()
	return ipv4IntfState, nil
}
func (svcHdlr AsicDaemonServiceHandler) GetBulkIPv4IntfState(currMarker, count asicdServices.Int) (*asicdServices.IPv4IntfStateGetInfo, error) {
	bulkObj := asicdServices.NewIPv4IntfStateGetInfo()
	return bulkObj, nil
}

//IPv4 Neighbor related services
func (svcHdlr AsicDaemonServiceHandler) CreateIPv4Neighbor(ipAddr string, macAddr string, vlanId, ifIndex int32) (rval int32, err error) {
	return (int32)(0), nil
}
func (svcHdlr AsicDaemonServiceHandler) UpdateIPv4Neighbor(ipAddr string, macAddr string, vlanId, ifIndex int32) (rval int32, err error) {
	return (int32)(0), nil
}
func (svcHdlr AsicDaemonServiceHandler) DeleteIPv4Neighbor(ipAddr string, macAddr string, vlanId, ifIndex int32) (rval int32, err error) {
	return (int32)(0), nil
}
func (svcHdlr AsicDaemonServiceHandler) GetBulkArpEntryHwState(currMarker, count asicdServices.Int) (*asicdServices.ArpEntryHwStateGetInfo, error) {
	bulkObj := asicdServices.NewArpEntryHwStateGetInfo()
	return bulkObj, nil
}
func (svcHdlr AsicDaemonServiceHandler) GetArpEntryHwState(ipAddr string) (*asicdServices.ArpEntryHwState, error) {
	arpEntryHwState := asicdServices.NewArpEntryHwState()
	return arpEntryHwState, nil
}

//IPv4 Route related services
func (svcHdlr AsicDaemonServiceHandler) OnewayCreateIPv4Route(ipv4RouteList []*asicdInt.IPv4Route) error {
	return nil
}
func (svcHdlr AsicDaemonServiceHandler) OnewayDeleteIPv4Route(ipv4RouteList []*asicdInt.IPv4Route) error {
	return nil
}

//IPv6 Route related services
func (svcHdlr AsicDaemonServiceHandler) OnewayCreateIPv6Route(ipv6RouteList []*asicdInt.IPv6Route) error {
	return nil
}
func (svcHdlr AsicDaemonServiceHandler) OnewayDeleteIPv6Route(ipv6RouteList []*asicdInt.IPv6Route) error {
	return nil
}
func (svcHdlr AsicDaemonServiceHandler) GetBulkIPRouteHwState(currMarker, count asicdInt.Int) (*asicdInt.IPRouteHwStateGetInfo, error) {
	bulkObj := asicdInt.NewIPRouteHwStateGetInfo()
	return bulkObj, nil
}
func (svcHdlr AsicDaemonServiceHandler) GetIPRouteHwState(destinationNw string) (*asicdInt.IPRouteHwState, error) {
	ipRouteHwState := asicdInt.NewIPRouteHwState()
	return ipRouteHwState, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkIPv4RouteHwState(currMarker, count asicdServices.Int) (*asicdServices.IPv4RouteHwStateGetInfo, error) {
	bulkObj := asicdServices.NewIPv4RouteHwStateGetInfo()
	return bulkObj, nil
}
func (svcHdlr AsicDaemonServiceHandler) GetIPv4RouteHwState(destinationNw string) (*asicdServices.IPv4RouteHwState, error) {
	ipv4RouteHwState := asicdServices.NewIPv4RouteHwState()
	return ipv4RouteHwState, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkIPv6RouteHwState(currMarker, count asicdServices.Int) (*asicdServices.IPv6RouteHwStateGetInfo, error) {
	bulkObj := asicdServices.NewIPv6RouteHwStateGetInfo()
	return bulkObj, nil
}
func (svcHdlr AsicDaemonServiceHandler) GetIPv6RouteHwState(destinationNw string) (*asicdServices.IPv6RouteHwState, error) {
	ipv6RouteHwState := asicdServices.NewIPv6RouteHwState()
	return ipv6RouteHwState, nil
}

func (svcHdlr AsicDaemonServiceHandler) CreateSubIPv4Intf(obj *asicdServices.SubIPv4Intf) (rv bool, err error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) UpdateSubIPv4Intf(oldObj, newObj *asicdServices.SubIPv4Intf, attrset []bool, op []*asicdServices.PatchOpInfo) (rv bool, err error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) DeleteSubIPv4Intf(obj *asicdServices.SubIPv4Intf) (rv bool, err error) {
	return true, nil
}

// IPv6 Interface Create/Update/Delete API's
func (svcHdlr AsicDaemonServiceHandler) CreateIPv6Intf(ipv6IntfObj *asicdServices.IPv6Intf) (rv bool, err error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) UpdateIPv6Intf(oldIPv6IntfObj, newIPv6IntfObj *asicdServices.IPv6Intf, attrset []bool,
	op []*asicdServices.PatchOpInfo) (rv bool, err error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) DeleteIPv6Intf(ipv6IntfObj *asicdServices.IPv6Intf) (rv bool, err error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetIPv6IntfState(intfRef string) (ipv6IntfState *asicdServices.IPv6IntfState, err error) {
	ipv6IntfState = asicdServices.NewIPv6IntfState()
	return ipv6IntfState, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkIPv6IntfState(currMarker, count asicdServices.Int) (*asicdServices.IPv6IntfStateGetInfo, error) {
	bulkObj := asicdServices.NewIPv6IntfStateGetInfo()
	return bulkObj, nil
}

//IPv6 Neighbor related services
func (svcHdlr AsicDaemonServiceHandler) CreateIPv6Neighbor(ipAddr string, macAddr string, vlanId, ifIndex int32) (rval int32, err error) {
	return (int32)(0), nil
}
func (svcHdlr AsicDaemonServiceHandler) UpdateIPv6Neighbor(ipAddr string, macAddr string, vlanId, ifIndex int32) (rval int32, err error) {
	return (int32)(0), nil
}
func (svcHdlr AsicDaemonServiceHandler) DeleteIPv6Neighbor(ipAddr string, macAddr string, vlanId,
	ifIndex int32) (rval int32, err error) {
	return (int32)(0), nil
}

// Sub IPv6 Interface Create/Delete/Update API's
func (svcHdlr AsicDaemonServiceHandler) CreateSubIPv6Intf(obj *asicdServices.SubIPv6Intf) (rv bool, err error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) UpdateSubIPv6Intf(oldObj, newObj *asicdServices.SubIPv6Intf, attrset []bool, op []*asicdServices.PatchOpInfo) (rv bool, err error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) DeleteSubIPv6Intf(obj *asicdServices.SubIPv6Intf) (rv bool, err error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkNdpEntryHwState(currMarker, count asicdServices.Int) (*asicdServices.NdpEntryHwStateGetInfo, error) {
	bulkObj := asicdServices.NewNdpEntryHwStateGetInfo()
	return bulkObj, nil
}
func (svcHdlr AsicDaemonServiceHandler) GetNdpEntryHwState(ipAddr string) (*asicdServices.NdpEntryHwState, error) {
	ndpEntryHwState := asicdServices.NewNdpEntryHwState()
	return ndpEntryHwState, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkLinkScopeIpState(currMarker, count asicdServices.Int) (*asicdServices.LinkScopeIpStateGetInfo, error) {
	bulkObj := asicdServices.NewLinkScopeIpStateGetInfo()
	return bulkObj, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetLinkScopeIpState(lsIpAddr string) (*asicdServices.LinkScopeIpState, error) {
	linkScopeIpState := asicdServices.NewLinkScopeIpState()
	return linkScopeIpState, nil
}

func (svcHdlr AsicDaemonServiceHandler) UpdateVirtualIPv4Intf(intfRef, subType, ipAddr, macAddr string, enable bool) (rv bool, err error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) UpdateVirtualIPv6Intf(intfRef, subType, ipAddr, macAddr string, enable bool) (rv bool, err error) {
	return true, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkSubIPv4IntfState(currMarker, count asicdServices.Int) (*asicdServices.SubIPv4IntfStateGetInfo, error) {
	bulkObj := asicdServices.NewSubIPv4IntfStateGetInfo()
	return bulkObj, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetSubIPv4IntfState(intfRef string, subType string) (*asicdServices.SubIPv4IntfState, error) {
	ipv4IntfState := asicdServices.NewSubIPv4IntfState()
	return ipv4IntfState, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetBulkSubIPv6IntfState(currMarker, count asicdServices.Int) (*asicdServices.SubIPv6IntfStateGetInfo, error) {
	bulkObj := asicdServices.NewSubIPv6IntfStateGetInfo()
	return bulkObj, nil
}

func (svcHdlr AsicDaemonServiceHandler) GetSubIPv6IntfState(intfRef string, subType string) (*asicdServices.SubIPv6IntfState, error) {
	ipv6IntfState := asicdServices.NewSubIPv6IntfState()
	return ipv6IntfState, nil
}
