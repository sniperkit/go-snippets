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
	"l3/ndp/config"
	"l3/ndp/packet"
	"reflect"
	"testing"
)

const (
	testMulticastSolicitationAddr = "ff02::1:ff7c:ca9f"
	testUnspecifiecSrcIp          = "::"
	//Unicast Neighbor Solicitation Const
	testNsSrcMac    = "00:1f:16:25:33:ce"
	testNsDstMac    = "00:1f:16:25:34:31"
	testNsSrcIp     = "fe80::21f:16ff:fe25:33ce"
	testNsDstIp     = "2001:db8:0:f101::1"
	testIsFastProbe = false
)

var nsBaseTestPkt = []byte{
	0x00, 0x1f, 0x16, 0x25, 0x34, 0x31, 0x00, 0x1f, 0x16, 0x25, 0x33, 0xce, 0x86, 0xdd, 0x60, 0x00,
	0x00, 0x00, 0x00, 0x20, 0x3a, 0xff, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x1f,
	0x16, 0xff, 0xfe, 0x25, 0x33, 0xce, 0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0xf1, 0x01, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x87, 0x00, 0xa6, 0x86, 0x00, 0x00, 0x00, 0x00, 0x20, 0x01,
	0x0d, 0xb8, 0x00, 0x00, 0xf1, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0x01,
	0x00, 0x1f, 0x16, 0x25, 0x33, 0xce,
}

func TestProcessNS(t *testing.T) {
	intf := &Interface{}
	ndInfo := &packet.NDInfo{
		DstIp: testMulticastSolicitationAddr,
	}
	nbrInfo, operType := intf.processNS(ndInfo)
	if nbrInfo != nil {
		t.Error("for testMulticastSolicitationAddr nbrInfo should be nil")
		return
	}
	if operType != IGNORE {
		t.Error("for testMulticastSolicitationAddr operation should be IGNORE, but got:", operType)
		return
	}
	ndInfo.SrcIp = testUnspecifiecSrcIp
	ndInfo.DstIp = ""
	nbrInfo, operType = intf.processNS(ndInfo)
	if nbrInfo != nil {
		t.Error("for testUnspecifiecSrcIp nbrInfo should be nil")
		return
	}
	if operType != IGNORE {
		t.Error("for testUnspecifiecSrcIp  operation should be IGNORE, but got:", operType)
		return
	}
}

func testCreateNeighbor(t *testing.T) *packet.NDInfo {
	ndInfo := &packet.NDInfo{
		SrcMac: testNsSrcMac,
		DstMac: testNsDstMac,
		SrcIp:  testNsSrcIp,
		DstIp:  testNsDstIp,
	}
	ndOpt := &packet.NDOption{
		Type:   packet.NDOptionTypeSourceLinkLayerAddress,
		Length: 1,
		Value:  []byte{0x00, 0x1f, 0x16, 0x25, 0x33, 0xce},
	}
	ndInfo.Options = append(ndInfo.Options, ndOpt)
	return ndInfo
}

func TestUnicastNS(t *testing.T) {
	// create ipv6 interface
	TestIPv6IntfCreate(t)
	myMac := testNsDstMac
	nbrMac := testNsSrcMac
	nbrIp := testNsSrcIp
	wantNbrInfo := &config.NeighborConfig{
		MacAddr: testNsSrcMac,
		IpAddr:  testNsSrcIp,
		Intf:    testIntfRef,
		IfIndex: testIfIndex,
	}

	l3Port, exists := testNdpServer.L3Port[testIfIndex]
	if !exists {
		t.Error("Failed to get L3 Port for ifIndex:", testIfIndex)
		return
	}
	operation := l3Port.SendNS(myMac, nbrMac, nbrIp, testIsFastProbe)
	if operation != IGNORE {
		t.Error("When no neighbor entries are present operation should be Ignore")
		return
	}
	// now create a neighbor entry and use that for sending unicast packet
	ndInfo := testCreateNeighbor(t)
	nbrInfo, oper := l3Port.processNS(ndInfo)
	if oper != CREATE {
		t.Error("When processing valid new neighbor solicitation message operation should be create")
		return
	}

	if !reflect.DeepEqual(nbrInfo, wantNbrInfo) {
		t.Error("Want Neigbor Info:", *wantNbrInfo, "but received nbrInfo:", *nbrInfo)
		return
	}
	//t.Log(*nbrInfo, oper)
	//t.Log(l3Port.Neighbor)
	operation = l3Port.SendNS(myMac, nbrMac, nbrIp, testIsFastProbe)
	if operation != IGNORE {
		t.Error("When neighbor entries are present operation should be Ignore and unicast packet should be send out")
		return
	}
}

func constructInvalidNdInfoNS() *packet.NDInfo {
	ndInfo := &packet.NDInfo{
		SrcMac: testNsSrcMac,
		DstMac: testNsDstMac,
		SrcIp:  testMulticastSolicitationAddr,
		DstIp:  testMulticastSolicitationAddr,
	}
	ndOpt := &packet.NDOption{
		Type:   packet.NDOptionTypeSourceLinkLayerAddress,
		Length: 1,
		Value:  []byte{0x00, 0x1f, 0x16, 0x25, 0x33, 0xce},
	}
	ndInfo.Options = append(ndInfo.Options, ndOpt)
	return ndInfo
}

func TestInvalidProcessNS(t *testing.T) {
	// create ipv6 interface
	TestIPv6IntfCreate(t)
	l3Port, exists := testNdpServer.L3Port[testIfIndex]
	if !exists {
		t.Error("Failed to get L3 Port for ifIndex:", testIfIndex)
		return
	}
	ndInfo := constructInvalidNdInfoNS()
	nbrInfo, oper := l3Port.processNS(ndInfo)
	if oper == CREATE || oper == DELETE {
		t.Error("Failed to ignore NS with ipv6 multicast prefix for ndInfo:", *ndInfo)
		return
	}
	if nbrInfo != nil {
		t.Error("Failed to ignore NS with ipv6 multicast prefix for ndInfo:", *ndInfo)
		return
	}
	ndInfo.SrcIp = ""
	nbrInfo, oper = l3Port.processNS(ndInfo)
	if oper == CREATE || oper == DELETE {
		t.Error("Failed to ignore NS with ipv6 multicast prefix for ndInfo:", *ndInfo)
		return
	}
	if nbrInfo != nil {
		t.Error("Failed to ignore NS with ipv6 multicast prefix for ndInfo:", *ndInfo)
		return
	}
	ndInfo.SrcIp = "::"
	nbrInfo, oper = l3Port.processNS(ndInfo)
	if oper == CREATE || oper == DELETE {
		t.Error("Failed to ignore NS with ipv6 multicast prefix for ndInfo:", *ndInfo)
		return
	}
	if nbrInfo != nil {
		t.Error("Failed to ignore NS with ipv6 multicast prefix for ndInfo:", *ndInfo)
		return
	}
}
