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
package packet

import (
	"bytes"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"l3/vrrp/common"
	"l3/vrrp/debug"
	"log/syslog"
	"net"
	"reflect"
	"syscall"
	"testing"
	"utils/logging"
)

var testPktInfo *PacketInfo

var testEncodePkt = []byte{
	0x01, 0x00, 0x5e, 0x00, 0x00, 0x12, 0x00, 0x00, 0x5e, 0x00, 0x01, 0x01, 0x08, 0x00, 0x45, 0x00,
	0x00, 0x28, 0x00, 0x00, 0x00, 0x00, 0xff, 0x70, 0x1a, 0x8d, 0xc0, 0xa8, 0x00, 0x1e, 0xe0, 0x00,
	0x00, 0x12, 0x21, 0x01, 0x64, 0x01, 0x00, 0x01, 0xba, 0x52, 0xc0, 0xa8, 0x00, 0x01, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

var version3Packet = []byte{
	0x01, 0x00, 0x5e, 0x00, 0x00, 0x12, 0x00, 0x00, 0x5e, 0x00, 0x01, 0x01, 0x08, 0x00, 0x45, 0x00,
	0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0xff, 0x70, 0x2f, 0x47, 0xac, 0x12, 0x00, 0x02, 0xe0, 0x00,
	0x00, 0x12, 0x31, 0x01, 0x64, 0x01, 0x00, 0x64, 0xbe, 0x85, 0xac, 0x12, 0x00, 0x01, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

var ipv6Packet = []byte{
	0x33, 0x33, 0x00, 0x00, 0x00, 0x12, 0x00, 0x00, 0x5e, 0x00, 0x02, 0x01, 0x86, 0xdd, 0x60, 0x00,
	0x00, 0x00, 0x00, 0x18, 0x70, 0xff, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x1d,
	0xfc, 0xff, 0xfe, 0xcf, 0x15, 0xfc, 0xff, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x12, 0x31, 0x01, 0x64, 0x01, 0x00, 0x64, 0xd1, 0x9d, 0xfe, 0x80,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x70, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01,
}
var testVrid = uint8(1)
var testPriority = uint8(100)
var testAdvInt = uint16(1)
var testVMac = "00:00:5e:00:01:01"
var testSrcIp = "192.168.0.30"
var testSrcIpV3 = "172.18.0.2"
var testVip = "192.168.0.1"
var testVipV3 = "172.18.0.1"
var testVMacIPv6 = "00:00:5e:00:02:01"
var testSrcIPv6 = "fe80::8a1d:fcff:fecf:15fc"
var testVipIPv6 = "fe80::70:1:1:1"

func TestInit(t *testing.T) {
	testPktInfo = Init()
	if testPktInfo == nil {
		t.Error("failed to initialize packet info")
		return
	}
	var err error
	logger := new(logging.Writer)
	logger.MyComponentName = "VRRPD"
	logger.SysLogger, err = syslog.New(syslog.LOG_INFO|syslog.LOG_DAEMON, "VRRPTEST")
	if err != nil {
		t.Error("failed to initialize syslog:", err)
	} else {
		logger.MyLogLevel = 9 // trace level
		debug.SetLogger(logger)
	}
}

func TestEncodeV2(t *testing.T) {
	TestInit(t)
	pktInfo := &PacketInfo{
		Version:      common.VERSION2,
		Vrid:         testVrid,
		Priority:     testPriority,
		AdvertiseInt: testAdvInt,
		VirutalMac:   testVMac,
		IpAddr:       testSrcIp,
		Vip:          testVip,
		IpType:       syscall.AF_INET,
	}
	encodedPkt := testPktInfo.Encode(pktInfo)
	if len(encodedPkt) != len(testEncodePkt) {
		t.Error("mis-match in length:", len(encodedPkt), len(testEncodePkt))
		return
	}
	if !bytes.Equal(encodedPkt, testEncodePkt) {
		t.Error("Failed to encode packet for pktInfo:", *pktInfo)
		t.Error("	testEncodePkt:", testEncodePkt)
		t.Error("	encoded Pkt:", encodedPkt)
		for idx, _ := range encodedPkt {
			if encodedPkt[idx] != testEncodePkt[idx] {
				t.Error("byte:", idx+1, "is not equal")
				t.Error(fmt.Sprintf("encoded Byte is:0x%x but wanted byte is:0x%x", encodedPkt[idx], testEncodePkt[idx]))
			}
		}
		return
	}
}

func TestEncodeV3(t *testing.T) {
	TestInit(t)
	pktInfo := &PacketInfo{
		Version:      common.VERSION3,
		Vrid:         testVrid,
		Priority:     testPriority,
		AdvertiseInt: testAdvInt,
		VirutalMac:   testVMac,
		IpAddr:       testSrcIpV3,
		Vip:          testVipV3,
		IpType:       syscall.AF_INET,
	}
	encodedPkt := testPktInfo.Encode(pktInfo)
	if len(encodedPkt) != len(version3Packet) {
		t.Error("mis-match in length:", len(encodedPkt), len(version3Packet))
		return
	}
	if !bytes.Equal(encodedPkt, version3Packet) {
		t.Error("Failed to encode packet for pktInfo:", *pktInfo)
		t.Error("	testEncodePkt:", version3Packet)
		t.Error("	encoded Pkt:", encodedPkt)
		for idx, _ := range encodedPkt {
			if encodedPkt[idx] != version3Packet[idx] {
				t.Error("byte:", idx+1, "is not equal")
				t.Error(fmt.Sprintf("encoded Byte is:0x%x but wanted byte is:0x%x", encodedPkt[idx], version3Packet[idx]))
			}
		}
		return
	}
}

func TestEncodeIPv6(t *testing.T) {
	TestInit(t)
	pktInfo := &PacketInfo{
		Version:      common.VERSION3,
		Vrid:         testVrid,
		Priority:     testPriority,
		AdvertiseInt: testAdvInt,
		VirutalMac:   testVMacIPv6,
		IpAddr:       testSrcIPv6,
		Vip:          testVipIPv6,
		IpType:       syscall.AF_INET6,
	}
	encodedPkt := testPktInfo.Encode(pktInfo)
	if len(encodedPkt) != len(ipv6Packet) {
		t.Error("mis-match in length:", len(encodedPkt), len(ipv6Packet))
		return
	}
	if !bytes.Equal(encodedPkt, ipv6Packet) {
		t.Error("Failed to encode packet for pktInfo:", *pktInfo)
		t.Error("	testEncodePkt:", ipv6Packet)
		t.Error("	encoded Pkt:", encodedPkt)
		for idx, _ := range encodedPkt {
			if encodedPkt[idx] != ipv6Packet[idx] {
				t.Error("byte:", idx+1, "is not equal")
				t.Error(fmt.Sprintf("encoded Byte is:0x%x but wanted byte is:0x%x", encodedPkt[idx], ipv6Packet[idx]))
			}
		}
		return
	}
}

func TestDecodeV2(t *testing.T) {
	TestInit(t)
	p := gopacket.NewPacket(testEncodePkt, layers.LinkTypeEthernet, gopacket.Default)
	decodePkt := testPktInfo.Decode(p)
	if decodePkt == nil {
		t.Error("failed to decode packet")
		return
	}
	wantPktInfo := &PacketInfo{
		DstMac: VRRP_PROTOCOL_MAC,
		SrcMac: testVMac,
		IpAddr: testSrcIp,
		DstIp:  VRRP_V4_GROUP_IP,
		Hdr: &Header{
			Version:      common.VERSION2,
			Type:         VRRP_PKT_TYPE_ADVERTISEMENT,
			VirtualRtrId: testVrid,
			Priority:     testPriority,
			CountIPAddr:  1,
			Rsvd:         0,
			MaxAdverInt:  testAdvInt,
			CheckSum:     uint16(47698),
		},
	}
	wantPktInfo.Hdr.IpAddr = append(wantPktInfo.Hdr.IpAddr, net.ParseIP(testVip).To4())
	if !reflect.DeepEqual(wantPktInfo, decodePkt) {
		t.Error("failed to decode packet")
		t.Error("wantPktInfo header is:", *wantPktInfo.Hdr, "entire packet info:", wantPktInfo)
		t.Error("decodePktInfo header is:", *decodePkt.Hdr, "entire packet info:", decodePkt)
		return
	}
}

func TestDecodeV3(t *testing.T) {
	TestInit(t)
	p := gopacket.NewPacket(version3Packet, layers.LinkTypeEthernet, gopacket.Default)
	decodePkt := testPktInfo.Decode(p)
	if decodePkt == nil {
		t.Error("failed to decode packet")
		return
	}

	wantPktInfo := &PacketInfo{
		DstMac: VRRP_PROTOCOL_MAC,
		SrcMac: testVMac,
		IpAddr: testSrcIpV3,
		DstIp:  VRRP_V4_GROUP_IP,
		Hdr: &Header{
			Version:      common.VERSION3,
			Type:         VRRP_PKT_TYPE_ADVERTISEMENT,
			VirtualRtrId: testVrid,
			Priority:     testPriority,
			CountIPAddr:  1,
			Rsvd:         0,
			MaxAdverInt:  testAdvInt,
			CheckSum:     uint16(48773),
		},
	}
	wantPktInfo.Hdr.IpAddr = append(wantPktInfo.Hdr.IpAddr, net.ParseIP(testVipV3).To4())
	if !reflect.DeepEqual(wantPktInfo, decodePkt) {
		t.Error("failed to decode packet")
		t.Error("wantPktInfo header is:", *wantPktInfo.Hdr, "entire packet info:", wantPktInfo)
		t.Error("decodePktInfo header is:", *decodePkt.Hdr, "entire packet info:", decodePkt)
		return
	}
}

func TestDecodeIPv6(t *testing.T) {
	TestInit(t)
	p := gopacket.NewPacket(ipv6Packet, layers.LinkTypeEthernet, gopacket.Default)
	decodePkt := testPktInfo.Decode(p)
	if decodePkt == nil {
		t.Error("failed to decode packet")
		return
	}

	wantPktInfo := &PacketInfo{
		DstMac: VRRP_V6_PROTOCOL_MAC,
		SrcMac: testVMacIPv6,
		IpAddr: testSrcIPv6,
		DstIp:  VRRP_V6_GROUP_IP,
		Hdr: &Header{
			Version:      common.VERSION3,
			Type:         VRRP_PKT_TYPE_ADVERTISEMENT,
			VirtualRtrId: testVrid,
			Priority:     testPriority,
			CountIPAddr:  1,
			Rsvd:         0,
			MaxAdverInt:  testAdvInt,
			CheckSum:     uint16(53661),
		},
	}
	wantPktInfo.Hdr.IpAddr = append(wantPktInfo.Hdr.IpAddr, net.ParseIP(testVipIPv6).To16())
	if !reflect.DeepEqual(wantPktInfo, decodePkt) {
		t.Error("failed to decode packet")
		t.Error("wantPktInfo header is:", *wantPktInfo.Hdr, "entire packet info:", wantPktInfo)
		t.Error("decodePktInfo header is:", *decodePkt.Hdr, "entire packet info:", decodePkt)
		return
	}
}
