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
// This class module should represet a mocked interface class for users to use.  It should be used as a 'BASE'
// class for testing purposes only.  It will help in saving time in having to create each of the methods
// when testing.  The testing method should overwrite any methods it sees fit for the test.

package mockasicdclientplugin

import (
	"fmt"
	"utils/commonDefs"
)

type MockAsicdClientMgr struct {
	Val int
}

func (asicdClientMgr *MockAsicdClientMgr) CreateIPv4Neighbor(ipAddr, macAddr string, vlanId, ifIdx int32) (int32, error) {
	fmt.Println(ipAddr, macAddr, vlanId, ifIdx, asicdClientMgr.Val)
	return 0, nil
}

func (asicdClientMgr *MockAsicdClientMgr) UpdateIPv4Neighbor(ipAddr, macAddr string, vlanId, ifIdx int32) (int32, error) {
	fmt.Println(ipAddr, macAddr, vlanId, ifIdx, asicdClientMgr.Val)
	return 0, nil
}

func (asicdClientMgr *MockAsicdClientMgr) CreateIPv6Neighbor(ipAddr, macAddr string, vlanId, ifIdx int32) (int32, error) {
	fmt.Println(ipAddr, macAddr, vlanId, ifIdx, asicdClientMgr.Val)
	return 0, nil
}

func (asicdClientMgr *MockAsicdClientMgr) UpdateIPv6Neighbor(ipAddr, macAddr string, vlanId, ifIdx int32) (int32, error) {
	fmt.Println(ipAddr, macAddr, vlanId, ifIdx, asicdClientMgr.Val)
	return 0, nil
}

func (asicdClientMgr *MockAsicdClientMgr) DeleteIPv4Neighbor(ipAddr string) (int32, error) {
	fmt.Println(ipAddr, asicdClientMgr.Val)
	return 0, nil
}

func (asicdClientMgr *MockAsicdClientMgr) DeleteIPv6Neighbor(ipAddr string) (int32, error) {
	fmt.Println(ipAddr, asicdClientMgr.Val)
	return 0, nil
}

func (asicdClientMgr *MockAsicdClientMgr) GetBulkIPv4IntfState(curMark, count int) (*commonDefs.IPv4IntfStateGetInfo, error) {
	fmt.Println("IPv4 Intf State", curMark, count, asicdClientMgr.Val)
	return nil, nil
}

func (asicdClientMgr *MockAsicdClientMgr) GetBulkPort(curMark, count int) (*commonDefs.PortGetInfo, error) {
	fmt.Println("Port Get info", curMark, count, asicdClientMgr.Val)
	return nil, nil
}

func (asicdClientMgr *MockAsicdClientMgr) GetBulkPortState(curMark, count int) (*commonDefs.PortStateGetInfo, error) {
	fmt.Println("Port State Get info", curMark, count, asicdClientMgr.Val)
	return nil, nil
}

func (asicdClientMgr *MockAsicdClientMgr) GetBulkVlanState(curMark, count int) (*commonDefs.VlanStateGetInfo, error) {
	fmt.Println("Vlan State Get info", curMark, count, asicdClientMgr.Val)
	return nil, nil
}

func (asicdClientMgr *MockAsicdClientMgr) GetBulkLag(curMark, count int) (*commonDefs.LagGetInfo, error) {
	fmt.Println("Lag Get info", curMark, count, asicdClientMgr.Val)
	return nil, nil
}

func (asicdClientMgr *MockAsicdClientMgr) GetBulkVlan(curMark, count int) (*commonDefs.VlanGetInfo, error) {
	fmt.Println("Vlan Get info", curMark, count, asicdClientMgr.Val)
	return nil, nil
}

func (asicdClientMgr *MockAsicdClientMgr) GetAllIPv4IntfState() ([]*commonDefs.IPv4IntfState, error) {
	fmt.Println("Get all IPv4 Intf State called")
	return nil, nil
}

func (asicdClientMgr *MockAsicdClientMgr) GetAllIPv6IntfState() ([]*commonDefs.IPv6IntfState, error) {
	fmt.Println("Get all IPv6 Intf State called")
	return nil, nil
}

func (asicdClientMgr *MockAsicdClientMgr) GetAllPortState() ([]*commonDefs.PortState, error) {
	fmt.Println("Get all Port Intf State called")
	return nil, nil
}

func (asicdClientMgr *MockAsicdClientMgr) GetAllVlanState() ([]*commonDefs.VlanState, error) {
	fmt.Println("Get all Vlan State called")
	return nil, nil
}

func (asicdClientMgr *MockAsicdClientMgr) GetAllVlan() ([]*commonDefs.Vlan, error) {
	fmt.Println("Get all Vlan called")
	return nil, nil
}

func (asicdClientMgr *MockAsicdClientMgr) GetPort(name string) (*commonDefs.Port, error) {
	fmt.Println("Get Port Intf State called")
	return nil, nil
}

func (asicdClientMgr *MockAsicdClientMgr) DetermineRouterId() string {
	return "0.0.0.0"
}

func (asicdClientMgr *MockAsicdClientMgr) GetPortLinkStatus(pId int32) bool {
	return true
}

func (asicdClientMgr *MockAsicdClientMgr) CreateStgBridge(vlanList []uint16) int32 {
	return -1
}

func (asicdClientMgr *MockAsicdClientMgr) DeleteStgBridge(stgid int32, vlanList []uint16) error {
	return nil
}

func (asicdClientMgr *MockAsicdClientMgr) SetStgPortState(stgid int32, ifindex int32, state int) error {
	return nil
}

func (asicdClientMgr *MockAsicdClientMgr) FlushStgFdb(stgid, ifindex int32) error {
	return nil
}

func (asicdClientMgr *MockAsicdClientMgr) BPDUGuardDetected(ifindex int32, enable bool) error {
	return nil
}

func (asicdClientMgr *MockAsicdClientMgr) GetSwitchMAC(paramsPath string) string {
	return "00:00:00:00:00:00"
}

func (asicdClientMgr *MockAsicdClientMgr) CreateLag(ifname string, hashType int32, ports string) (hwAggId int32, err error) {
	return -1, err
}

func (asicdClientMgr *MockAsicdClientMgr) DeleteLag(hwAggId int32) (err error) {
	return err
}

func (asicdClientMgr *MockAsicdClientMgr) UpdateLag(ifIndex, hashType int32, ports string) (err error) {
	return err
}

func (asicdClientMgr *MockAsicdClientMgr) UpdateLagCfgIntfList(ifName string, ifIndexList []int32) bool {
	return true
}

func (asicdClientMgr *MockAsicdClientMgr) EnablePacketReception(mac string, vlan int, ifindex int32) (err error) {
	return err
}

func (asicdClientMgr *MockAsicdClientMgr) DisablePacketReception(mac string, vlan int, ifindex int32) (err error) {
	return err
}

func (asicdClientMgr *MockAsicdClientMgr) IppIngressEgressDrop(srcIfIndex, dstIfIndex string) (err error) {
	return err
}

func (asicdClientMgr *MockAsicdClientMgr) IppIngressEgressPass(srcIfIndex, dstIfIndex string) (err error) {
	return err
}

func (asicdClientMgr *MockAsicdClientMgr) IppVlanConversationSet(vlan uint16, ifindex int32) (err error) {
	return err
}
func (asicdClientMgr *MockAsicdClientMgr) IppVlanConversationClear(vlan uint16, ifindex int32) (err error) {
	return err
}
func (asicdClientMgr *MockAsicdClientMgr) IsLoopbackType(ifIndex int32) bool {
	return true
}

func (asicdClientMgr *MockAsicdClientMgr) CreateVirtualIPv4Intf(intRef, ipAddr, macAddr string, enable bool) (err error) {
	return err
}

func (asicdClientMgr *MockAsicdClientMgr) CreateVirtualIPv6Intf(intRef, ipAddr, macAddr string, enable bool) (err error) {
	return err
}

func (asicdClientMgr *MockAsicdClientMgr) UpdateVirtualIPv4Intf(intRef, ipAddr, macAddr string, enable bool) (err error) {
	return err
}

func (asicdClientMgr *MockAsicdClientMgr) UpdateVirtualIPv6Intf(intRef, ipAddr, macAddr string, enable bool) (err error) {
	return err
}

func (asicdClientMgr *MockAsicdClientMgr) DeleteVirtualIPv4Intf(intRef, ipAddr, macAddr string, enable bool) (err error) {
	return err
}

func (asicdClientMgr *MockAsicdClientMgr) DeleteVirtualIPv6Intf(intRef, ipAddr, macAddr string, enable bool) (err error) {
	return err
}

func (asicdClientMgr *MockAsicdClientMgr) GetAllSubIPv4IntfState() ([]*commonDefs.SubIPv4IntfState, error) {
	return nil, nil
}

func (asicdClientMgr *MockAsicdClientMgr) GetAllSubIPv6IntfState() ([]*commonDefs.SubIPv6IntfState, error) {
	return nil, nil
}

func (asicdClientMgr *MockAsicdClientMgr) IsLinuxOnlyPlugin() (bool, error) {
	return false, nil
}
