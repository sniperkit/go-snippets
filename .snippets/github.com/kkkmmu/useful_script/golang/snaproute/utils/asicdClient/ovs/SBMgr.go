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

package ovs

import (
	"fmt"
	"utils/commonDefs"
)

type OvsAsicdClientMgr struct {
	Val int
}

func (asicdClientMgr *OvsAsicdClientMgr) CreateIPv4Neighbor(ipAddr, macAddr string, vlanId, ifIdx int32) (int32, error) {
	fmt.Println(ipAddr, macAddr, vlanId, ifIdx, asicdClientMgr.Val)
	return 0, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) UpdateIPv4Neighbor(ipAddr, macAddr string, vlanId, ifIdx int32) (int32, error) {
	fmt.Println(ipAddr, macAddr, vlanId, ifIdx, asicdClientMgr.Val)
	return 0, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) CreateIPv6Neighbor(ipAddr, macAddr string, vlanId, ifIdx int32) (int32, error) {
	fmt.Println(ipAddr, macAddr, vlanId, ifIdx, asicdClientMgr.Val)
	return 0, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) UpdateIPv6Neighbor(ipAddr, macAddr string, vlanId, ifIdx int32) (int32, error) {
	fmt.Println(ipAddr, macAddr, vlanId, ifIdx, asicdClientMgr.Val)
	return 0, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) DeleteIPv4Neighbor(ipAddr string) (int32, error) {
	fmt.Println(ipAddr, asicdClientMgr.Val)
	return 0, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) DeleteIPv6Neighbor(ipAddr string) (int32, error) {
	fmt.Println(ipAddr, asicdClientMgr.Val)
	return 0, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) GetBulkIPv4IntfState(curMark, count int) (*commonDefs.IPv4IntfStateGetInfo, error) {
	fmt.Println("IPv4 Intf State", curMark, count, asicdClientMgr.Val)
	return nil, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) GetBulkPort(curMark, count int) (*commonDefs.PortGetInfo, error) {
	fmt.Println("Port Get info", curMark, count, asicdClientMgr.Val)
	return nil, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) GetBulkPortState(curMark, count int) (*commonDefs.PortStateGetInfo, error) {
	fmt.Println("Port State Get info", curMark, count, asicdClientMgr.Val)
	return nil, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) GetBulkVlanState(curMark, count int) (*commonDefs.VlanStateGetInfo, error) {
	fmt.Println("Vlan State Get info", curMark, count, asicdClientMgr.Val)
	return nil, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) GetBulkLag(curMark, count int) (*commonDefs.LagGetInfo, error) {
	fmt.Println("Lag Get info", curMark, count, asicdClientMgr.Val)
	return nil, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) GetBulkVlan(curMark, count int) (*commonDefs.VlanGetInfo, error) {
	fmt.Println("Vlan Get info", curMark, count, asicdClientMgr.Val)
	return nil, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) GetAllIPv4IntfState() ([]*commonDefs.IPv4IntfState, error) {
	fmt.Println("Get all IPv4 Intf State called")
	return nil, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) GetAllIPv6IntfState() ([]*commonDefs.IPv6IntfState, error) {
	fmt.Println("Get all IPv6 Intf State called")
	return nil, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) GetAllPortState() ([]*commonDefs.PortState, error) {
	fmt.Println("Get all Port Intf State called")
	return nil, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) GetAllVlanState() ([]*commonDefs.VlanState, error) {
	fmt.Println("Get all Vlan State called")
	return nil, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) GetAllVlan() ([]*commonDefs.Vlan, error) {
	fmt.Println("Get all Vlan called")
	return nil, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) GetPort(name string) (*commonDefs.Port, error) {
	fmt.Println("Get Port Intf State called")
	return nil, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) DetermineRouterId() string {
	return "0.0.0.0"
}

func (asicdClientMgr *OvsAsicdClientMgr) GetPortLinkStatus(pId int32) bool {
	return true
}

func (asicdClientMgr *OvsAsicdClientMgr) CreateStgBridge(vlanList []uint16) int32 {
	return -1
}

func (asicdClientMgr *OvsAsicdClientMgr) DeleteStgBridge(stgid int32, vlanList []uint16) error {
	return nil
}

func (asicdClientMgr *OvsAsicdClientMgr) SetStgPortState(stgid int32, ifindex int32, state int) error {
	return nil
}

func (asicdClientMgr *OvsAsicdClientMgr) FlushStgFdb(stgid, ifindex int32) error {
	return nil
}

func (asicdClientMgr *OvsAsicdClientMgr) BPDUGuardDetected(ifindex int32, enable bool) error {
	return nil
}

func (asicdClientMgr *OvsAsicdClientMgr) GetSwitchMAC(paramsPath string) string {
	return "00:00:00:00:00:00"
}

func (asicdClientMgr *OvsAsicdClientMgr) CreateLag(ifname string, hashType int32, ports string) (hwAggId int32, err error) {
	return -1, err
}

func (asicdClientMgr *OvsAsicdClientMgr) DeleteLag(hwAggId int32) (err error) {
	return err
}

func (asicdClientMgr *OvsAsicdClientMgr) UpdateLag(ifIndex, hashType int32, ports string) (err error) {
	return err
}

func (asicdClientMgr *OvsAsicdClientMgr) UpdateLagCfgIntfList(ifName string, ifIndexList []int32) bool {
	return true
}

func (asicdClientMgr *OvsAsicdClientMgr) EnablePacketReception(mac string, vlan int, ifindex int32) (err error) {
	return err
}

func (asicdClientMgr *OvsAsicdClientMgr) DisablePacketReception(mac string, vlan int, ifindex int32) (err error) {
	return err
}

func (asicdClientMgr *OvsAsicdClientMgr) IppIngressEgressDrop(srcIfIndex, dstIfIndex string) (err error) {
	return err
}

func (asicdClientMgr *OvsAsicdClientMgr) IppIngressEgressPass(srcIfIndex, dstIfIndex string) (err error) {
	return err
}

func (asicdClientMgr *OvsAsicdClientMgr) IppVlanConversationSet(vlan uint16, ifindex int32) (err error) {
	return err
}
func (asicdClientMgr *OvsAsicdClientMgr) IppVlanConversationClear(vlan uint16, ifindex int32) (err error) {
	return err
}

func (asicdClientMgr *OvsAsicdClientMgr) IsLoopbackType(ifindex int32) bool {
	return true
}

func (asicdClientMgr *OvsAsicdClientMgr) CreateVirtualIPv4Intf(intRef, ipAddr, macAddr string, enable bool) (err error) {
	return err
}

func (asicdClientMgr *OvsAsicdClientMgr) CreateVirtualIPv6Intf(intRef, ipAddr, macAddr string, enable bool) (err error) {
	return err
}

func (asicdClientMgr *OvsAsicdClientMgr) UpdateVirtualIPv4Intf(intRef, ipAddr, macAddr string, enable bool) (err error) {
	return err
}

func (asicdClientMgr *OvsAsicdClientMgr) UpdateVirtualIPv6Intf(intRef, ipAddr, macAddr string, enable bool) (err error) {
	return err
}

func (asicdClientMgr *OvsAsicdClientMgr) DeleteVirtualIPv4Intf(intRef, ipAddr, macAddr string, enable bool) (err error) {
	return err
}

func (asicdClientMgr *OvsAsicdClientMgr) DeleteVirtualIPv6Intf(intRef, ipAddr, macAddr string, enable bool) (err error) {
	return err
}
func (asicdClientMgr *OvsAsicdClientMgr) GetAllSubIPv4IntfState() ([]*commonDefs.SubIPv4IntfState, error) {
	return nil, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) GetAllSubIPv6IntfState() ([]*commonDefs.SubIPv6IntfState, error) {
	return nil, nil
}

func (asicdClientMgr *OvsAsicdClientMgr) IsLinuxOnlyPlugin() (bool, error) {
	return false, nil
}
