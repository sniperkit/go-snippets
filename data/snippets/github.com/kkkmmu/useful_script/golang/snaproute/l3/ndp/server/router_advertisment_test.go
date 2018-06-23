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
	"github.com/google/gopacket/layers"
	"l3/ndp/config"
	"l3/ndp/packet"
	"reflect"
	"testing"
)

const (
	//Router Advertisement Const
	testRASrcMac      = "88:1d:fc:cf:15:fc"
	testRADstMac      = "33:33:00:00:00:01"
	testRALinkScopeIp = "fe80::8a1d:fcff:fecf:15fc"
	testRADstIp       = "ff02::1"

	testLinkScopeIp = "fe80::8a1d:fcff:fecf:15fc"
)

var raBaseTestPkt = []byte{
	0x33, 0x33, 0x00, 0x00, 0x00, 0x01, 0x88, 0x1d, 0xfc, 0xcf, 0x15, 0xfc, 0x86, 0xdd, 0x60, 0x00,
	0x00, 0x00, 0x00, 0x20, 0x3a, 0xff, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x1d,
	0xfc, 0xff, 0xfe, 0xcf, 0x15, 0xfc, 0xff, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x86, 0x00, 0xf2, 0x66, 0x40, 0x00, 0x07, 0x08, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0x88, 0x1d, 0xfc, 0xcf, 0x15, 0xfc, 0x05, 0x01,
	0x00, 0x00, 0x00, 0x00, 0x05, 0xdc,
}

func constructBaseWantNDInfoForRA() *packet.NDInfo {
	wantBaseNDInfo := &packet.NDInfo{
		CurHopLimit:    64,
		ReservedFlags:  0,
		RouterLifetime: 1800,
		ReachableTime:  0,
		RetransTime:    0,
		PktType:        layers.ICMPv6TypeRouterAdvertisement,
		SrcMac:         testRASrcMac,
		DstMac:         testRADstMac,
		SrcIp:          testRALinkScopeIp,
		DstIp:          testRADstIp,
	}

	sourcendOpt := &packet.NDOption{
		Type:   packet.NDOptionTypeSourceLinkLayerAddress,
		Length: 1,
	}
	sourcendOpt.Value = make([]byte, 6)
	copy(sourcendOpt.Value, raBaseTestPkt[72:78])
	mtuOpt := &packet.NDOption{
		Type:   packet.NDOptionTypeMTU,
		Length: 1,
	}
	for i := 0; i < 4; i++ {
		mtuOpt.Value = append(mtuOpt.Value, 0)
	}
	mtuOpt.Value = append(mtuOpt.Value, 0x05)
	mtuOpt.Value = append(mtuOpt.Value, 0xdc)
	wantBaseNDInfo.Options = append(wantBaseNDInfo.Options, sourcendOpt)
	wantBaseNDInfo.Options = append(wantBaseNDInfo.Options, mtuOpt)
	return wantBaseNDInfo
}

func TestSendRA(t *testing.T) {
	initServerBasic()
	intf := Interface{
		linkScope: testLinkScopeIp,
	}
	intf.SendRA(testSrcMac)
}

func TestProcessRA(t *testing.T) {
	// create ipv6 interface
	TestIPv6IntfCreate(t)
	ndInfo := constructBaseWantNDInfoForRA()
	wantNbrInfo := &config.NeighborConfig{
		MacAddr: testRASrcMac,
		IpAddr:  testRALinkScopeIp,
		Intf:    testIntfRef,
		IfIndex: testIfIndex,
	}
	l3Port, exists := testNdpServer.L3Port[testIfIndex]
	if !exists {
		t.Error("Failed to get L3 Port for ifIndex:", testIfIndex)
		return
	}
	nbrInfo, oper := l3Port.processRA(ndInfo)
	if oper != CREATE {
		t.Error("Failed to create a new neighbor entry on RA packet")
		return
	}
	if !reflect.DeepEqual(wantNbrInfo, nbrInfo) {
		t.Error("Want Neigbor Info:", *wantNbrInfo, "but received nbrInfo:", *nbrInfo)
		return
	}
	testNdpServer.CreateNeighborInfo(nbrInfo, testIfIndex, testIntfRef)
	if len(testNdpServer.NeighborInfo) == 0 {
		t.Error("Failed to add new learned neighbor information after processing router advertisement")
		return
	}

	if len(testNdpServer.neighborKey) == 0 {
		t.Error("Neighbor insert into NeighborInfo map but key is not added into neighborKey slice")
		return
	}
	nbrInfo, oper = l3Port.processRA(ndInfo)
	if oper != UPDATE {
		t.Error("Failed to create a new neighbor entry on RA packet")
		return
	}
	if !reflect.DeepEqual(wantNbrInfo, nbrInfo) {
		t.Error("During Update Want Neigbor Info:", *wantNbrInfo, "but received nbrInfo:", *nbrInfo)
		return
	}

	// test router advertisement delete packet
	ndInfo.RouterLifetime = 0
	nbrInfo, oper = l3Port.processRA(ndInfo)
	if oper != DELETE {
		t.Error("Failed to delete neighbor entry when Router Advertisement packet is received with routerLifeTime = 0")
		return
	}
	if nbrInfo == nil {
		t.Error("During router advertisement delete nbr Info should not be nil")
		return
	}
	// delete neighbor entry that just got created using router advertisement
	testNdpServer.deleteNeighbor(createNeighborKey(nbrInfo.MacAddr, nbrInfo.IpAddr, testIntfRef), testIfIndex)
	if len(testNdpServer.NeighborInfo) != 0 {
		t.Error("Failed to delete learned neighbor information after processing router advertisement")
		return
	}

	if len(testNdpServer.neighborKey) != 0 {
		t.Error("Neighbor delete into NeighborInfo map but key is not delete from into neighborKey slice")
		return
	}
}

func constructNDInfo() *packet.NDInfo {
	wantBaseNDInfo := &packet.NDInfo{
		CurHopLimit:    64,
		ReservedFlags:  0,
		RouterLifetime: 1800,
		ReachableTime:  0,
		RetransTime:    0,
		PktType:        layers.ICMPv6TypeRouterAdvertisement,
		SrcMac:         testRASrcMac,
		DstMac:         testRADstMac,
		DstIp:          testRALinkScopeIp,
		SrcIp:          testRADstIp,
	}

	sourcendOpt := &packet.NDOption{
		Type:   packet.NDOptionTypeSourceLinkLayerAddress,
		Length: 1,
	}
	sourcendOpt.Value = make([]byte, 6)
	copy(sourcendOpt.Value, raBaseTestPkt[72:78])
	mtuOpt := &packet.NDOption{
		Type:   packet.NDOptionTypeMTU,
		Length: 1,
	}
	for i := 0; i < 4; i++ {
		mtuOpt.Value = append(mtuOpt.Value, 0)
	}
	mtuOpt.Value = append(mtuOpt.Value, 0x05)
	mtuOpt.Value = append(mtuOpt.Value, 0xdc)
	wantBaseNDInfo.Options = append(wantBaseNDInfo.Options, sourcendOpt)
	wantBaseNDInfo.Options = append(wantBaseNDInfo.Options, mtuOpt)
	return wantBaseNDInfo
}

func TestInvalidProcessRA(t *testing.T) {
	// create ipv6 interface
	TestIPv6IntfCreate(t)
	l3Port, exists := testNdpServer.L3Port[testIfIndex]
	if !exists {
		t.Error("Failed to get L3 Port for ifIndex:", testIfIndex)
		return
	}

	t.Log(l3Port)
	ndInfo := constructNDInfo()
	nbrInfo, oper := l3Port.processRA(ndInfo)
	if nbrInfo != nil {
		t.Error("For ndInfo", *ndInfo, "which has ipv6 muticast prefix we should not create nbr Entry", nbrInfo)
		return
	}
	if oper == CREATE || oper == DELETE {
		t.Error("For ndInfo", *ndInfo, "which has ipv6 muticast prefix we should not have oper as ", oper)
		return
	}
}
