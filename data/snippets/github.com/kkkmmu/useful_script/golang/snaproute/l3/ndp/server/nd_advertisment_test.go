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
	_ "l3/ndp/config"
	"l3/ndp/packet"
	"net"
	_ "reflect"
	"testing"
)

const (
	//Neighbor Advertisement Const
	testNaSrcMac = "f6:6d:e4:22:75:9e"
	testNaDstMac = "00:1f:16:25:3e:71"
	testNaSrcIp  = "2149::61:123:1"
	testNaDstIp  = "2149::61:123:2"
)

// eth1_icmpv6.pcap
var naBaseTestPkt = []byte{
	0x00, 0x1f, 0x16, 0x25, 0x3e, 0x71, 0xf6, 0x6d, 0xe4, 0x22, 0x75, 0x9e, 0x86, 0xdd, 0x60, 0x00,
	0x00, 0x00, 0x00, 0x20, 0x3a, 0xff, 0x21, 0x49, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x61, 0x01, 0x23, 0x00, 0x01, 0x21, 0x49, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x61, 0x01, 0x23, 0x00, 0x02, 0x88, 0x00, 0xdd, 0x08, 0xe0, 0x00, 0x00, 0x00, 0x21, 0x49,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x61, 0x01, 0x23, 0x00, 0x01, 0x02, 0x01,
	0xf6, 0x6d, 0xe4, 0x22, 0x75, 0x9e,
}

func testInitNANdInfo() *packet.NDInfo {
	wantNDinfo := &packet.NDInfo{
		SrcMac:        testNaSrcMac, // Update SRC MAC From ethernet
		DstMac:        testNaDstMac, // Update DST MAC from ethernet
		SrcIp:         testNaSrcIp,  // copy sender ip address to this
		DstIp:         testNaDstIp,  // copy destination ip
		TargetAddress: net.ParseIP(testNaSrcIp),
		PktType:       layers.ICMPv6TypeNeighborAdvertisement,
	}
	ndOpt := &packet.NDOption{
		Type:   packet.NDOptionTypeTargetLinkLayerAddress,
		Length: 1,
		Value:  []byte{0xf6, 0x6d, 0xe4, 0x22, 0x75, 0x9e},
	}
	wantNDinfo.Options = append(wantNDinfo.Options, ndOpt)
	return wantNDinfo
}

func TestProcessNAPkt(t *testing.T) {
	// create ipv6 interface
	TestIPv6IntfCreate(t)
	ndInfo := testInitNANdInfo()
	l3Port, exists := testNdpServer.L3Port[testIfIndex]
	if !exists {
		t.Error("Failed to get L3 Port for ifIndex:", testIfIndex)
		return
	}
	_, oper := l3Port.processNA(ndInfo)
	if oper == CREATE {
		t.Error("if entry doesn't exists then NA should not have create operation")
		return
	}
	/*
		@TODO: add test case for R,S,O flags
		wantNbrInfo := &config.NeighborConfig{
			MacAddr: testNaSrcMac,
			IpAddr:  testNaSrcIp,
			Intf:    testIntfRef,
			IfIndex: testIfIndex,
		}
		if !reflect.DeepEqual(wantNbrInfo, nbrInfo) {
			t.Error("Want Neigbor Info:", *wantNbrInfo, "but received nbrInfo:", *nbrInfo)
			return
		}
		nbrInfo, oper = l3Port.processNA(ndInfo)
		if oper != UPDATE {
			t.Error("Failed to update existing neighbor entry for ndp")
			return
		}

		if !reflect.DeepEqual(wantNbrInfo, nbrInfo) {
			t.Error("During Update Want Neigbor Info:", *wantNbrInfo, "but received nbrInfo:", *nbrInfo)
			return
		}
	*/
}

func constructInvalidNdInfoNA() *packet.NDInfo {
	wantNDinfo := &packet.NDInfo{
		SrcMac:        testNaSrcMac, // Update SRC MAC From ethernet
		DstMac:        testNaDstMac, // Update DST MAC from ethernet
		SrcIp:         testRADstIp,
		DstIp:         testRADstIp,
		TargetAddress: net.ParseIP(testNaSrcIp),
		PktType:       layers.ICMPv6TypeNeighborAdvertisement,
	}
	ndOpt := &packet.NDOption{
		Type:   packet.NDOptionTypeTargetLinkLayerAddress,
		Length: 1,
		Value:  []byte{0xf6, 0x6d, 0xe4, 0x22, 0x75, 0x9e},
	}
	wantNDinfo.Options = append(wantNDinfo.Options, ndOpt)
	return wantNDinfo
}

func TestInvalidProcessNA(t *testing.T) {
	// create ipv6 interface
	TestIPv6IntfCreate(t)

	ndInfo := constructInvalidNdInfoNA()
	l3Port, exists := testNdpServer.L3Port[testIfIndex]
	if !exists {
		t.Error("Failed to get L3 Port for ifIndex:", testIfIndex)
		return
	}
	nbrInfo, oper := l3Port.processNA(ndInfo)
	if oper == CREATE || oper == DELETE {
		t.Error("Failed to ignore NA with ipv6 multicast prefix for ndInfo:", *ndInfo)
		return
	}
	if nbrInfo != nil {
		t.Error("Failed to ignore NA with ipv6 multicast prefix for ndInfo:", *ndInfo)
		return
	}
}
