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

// bgp_test.go
package packet

import (
	"encoding/hex"
	"fmt"
	"net"
	"testing"
)

func TestMPNextHopIP(t *testing.T) {
	//var mpNH MPNextHop
	mpNHIP := NewMPNextHopIP()
	var ip net.IP
	err := mpNHIP.SetNextHop(ip)
	if err == nil {
		t.Fatal("MPNextHopIP.SetNextHop for ip nil, expected failure, got NO error")
	} else {
		t.Log("MPNextHopIP.SetNextHop for ip nil, expected failure, got error:", err)
	}

	ips := []string{"10.1.10.1", "2002::1"}
	expectedLen := []uint8{5, 17}
	for i, ipStr := range ips {
		ip := net.ParseIP(ipStr)
		err = mpNHIP.SetNextHop(ip)
		if err != nil {
			t.Fatal("MPNextHopIP.SetNextHop failed for ip:", ipStr, "with error:", err)
		} else {
			t.Log("MPNextHopIP.SetNextHop successful for ip:", ipStr)
		}

		mpNHIPLen := mpNHIP.Len()
		if mpNHIPLen != expectedLen[i] {
			t.Fatal("MPNextHopIP.Len failed for ip:", ipStr, "expected:", expectedLen[i], "got:", mpNHIPLen)
		}
	}
}

func TestMPNextHopIP6(t *testing.T) {
	//var mpNH MPNextHop
	mpNHIP := NewMPNextHopIP6()
	var ip net.IP
	err := mpNHIP.SetGlobalNextHop(ip)
	if err == nil {
		t.Fatal("MPNextHopIP6.SetGlobalNextHop for ip nil, expected failure, got NO error")
	} else {
		t.Log("MPNextHopIP6.SetGlobalNextHop for ip nil, expected failure, got error:", err)
	}

	globalIPs := []string{"2002::1", "2002::1", "4001::1"}
	localIPs := []string{"", "fe80:2001::1", "fe80:2001::1"}
	expectedLen := []uint8{17, 33, 33}
	for i, ipStr := range globalIPs {
		gIP := net.ParseIP(ipStr)
		err = mpNHIP.SetGlobalNextHop(gIP)
		if err != nil {
			t.Fatal("MPNextHopIP6.SetGlobalNextHop failed for ip:", ipStr, "with error:", err)
		} else {
			t.Log("MPNextHopIP6.SetGlobalNextHop successful for ip:", ipStr)
		}

		if localIPs[i] != "" {
			lIP := net.ParseIP(localIPs[i])
			err = mpNHIP.SetLinkLocalNextHop(lIP)
			if err != nil {
				t.Fatal("MPNextHopIP6.SetLinkLocalNextHop failed for ip:", localIPs[i], "with error:", err)
			} else {
				t.Log("MPNextHopIP6.SetLinkLocalNextHop successful for ip:", localIPs[i])
			}
		}

		mpNHIPLen := mpNHIP.Len()
		if mpNHIPLen != expectedLen[i] {
			t.Fatal("MPNextHopIP6.Len failed for global ip:", ipStr, "local ip:", localIPs[i], "expected:",
				expectedLen[i], "got:", mpNHIPLen)
		}
	}
}

func TestMPReachNLRIDecode(t *testing.T) {
	packets := make([]string, 0)
	packets = append(packets, "800E10000102040A010A0101001814010A0A0A01")
	packets = append(packets, "800E11000102040A010A01001814010A0A0A01")
	packets = append(packets, "900E002F000202102001000000000000000000000A010A01001814010A0A0A014020013001400150014B90028002700260025002")
	packets = append(packets, "900E0031000202102001000000000000000000000A010A01001814010A0A0A014020013001400150014B90028002700260025002")

	for _, strPkt := range packets {
		hexPkt, err := hex.DecodeString(strPkt)
		fmt.Printf("packet = %x, len = %d\n", hexPkt, len(hexPkt))
		if err != nil {
			t.Fatal("Failed to decode the string to hex, string =", strPkt)
		}

		mpReach := NewBGPPathAttrMPReachNLRI()
		peerAttrs := BGPPeerAttrs{
			ASSize:           4,
			AddPathsRxActual: false,
		}
		err = mpReach.Decode(hexPkt, peerAttrs)
		if err == nil {
			t.Fatal("BGP MPReachNLRI decode called... expected failure, got NO errors")
		} else {
			t.Log("BGP MPReachNLRI decode called... expected failure, got error:", err)
		}
	}

	packets = make([]string, 0)
	packets = append(packets, "800E10000102040A010A01001814010A0A0A01")
	packets = append(packets, "900E0030000202102001000000000000000000000A010A01001814010A0A0A014020013001400150014B90028002700260025002")
	packets = append(packets, "900E0040000202202001000000000000000000000A010A01FE80000000000000000000000A010A01001814010A0A0A014020013001400150014B90028002700260025002")
	packets = append(packets, "800E10ABCD02040A010A01001814010A0A0A01")
	packets = append(packets, "900E0040DEAD02202001000000000000000000000A010A01FE80000000000000000000000A010A01001814010A0A0A014020013001400150014B90028002700260025002")

	for _, strPkt := range packets {
		hexPkt, err := hex.DecodeString(strPkt)
		fmt.Printf("packet = %x, len = %d\n", hexPkt, len(hexPkt))
		if err != nil {
			t.Fatal("Failed to decode the string to hex, string =", strPkt)
		}

		mpReach := NewBGPPathAttrMPReachNLRI()
		peerAttrs := BGPPeerAttrs{
			ASSize:           4,
			AddPathsRxActual: false,
		}
		err = mpReach.Decode(hexPkt, peerAttrs)
		if err != nil {
			t.Fatal("BGP MPReachNLRI decode failed with error:", err)
		} else {
			t.Log("BGP MPReachNLRI decode succeeded")
		}
	}
}

func TestMPUnreachNLRIDecode(t *testing.T) {
	packets := make([]string, 0)
	packets = append(packets, "800F13000102000000011814010A000000020A0A01")
	packets = append(packets, "900F002D000202010203041814010A010203050A0A0101020307402001300140015001010203084B90028002700260025002")
	packets = append(packets, "800F11ABCD02010203041814010A010203050A0A01")
	packets = append(packets, "900F002FDEAD02010203041814010A010203050A0A0101020307402001300140015001010203084B90028002700260025002")

	for _, strPkt := range packets {
		hexPkt, err := hex.DecodeString(strPkt)
		fmt.Printf("packet = %x, len = %d\n", hexPkt, len(hexPkt))
		if err != nil {
			t.Fatal("Failed to decode the string to hex, string =", strPkt)
		}

		mpReach := NewBGPPathAttrMPReachNLRI()
		peerAttrs := BGPPeerAttrs{
			ASSize:           4,
			AddPathsRxActual: false,
		}
		err = mpReach.Decode(hexPkt, peerAttrs)
		if err == nil {
			t.Fatal("BGP MPUnreachNLRI decode called... expected failure, got NO errors")
		} else {
			t.Log("BGP MPUnreachNLRI decode called... expected failure, got error:", err)
		}
	}

	packets = make([]string, 0)
	packets = append(packets, "800F12000102000000011814010A000000020A0A01")
	packets = append(packets, "900F002E000202010203041814010A010203050A0A0101020307402001300140015001010203084B90028002700260025002")
	packets = append(packets, "800F12ABCD02010203041814010A010203050A0A01")
	packets = append(packets, "900F002EDEAD02010203041814010A010203050A0A0101020307402001300140015001010203084B90028002700260025002")

	for _, strPkt := range packets {
		hexPkt, err := hex.DecodeString(strPkt)
		fmt.Printf("packet = %x, len = %d\n", hexPkt, len(hexPkt))
		if err != nil {
			t.Fatal("Failed to decode the string to hex, string =", strPkt)
		}

		mpReach := NewBGPPathAttrMPUnreachNLRI()
		peerAttrs := BGPPeerAttrs{
			ASSize:           4,
			AddPathsRxActual: true,
		}
		err = mpReach.Decode(hexPkt, peerAttrs)
		if err != nil {
			t.Fatal("BGP MPUnreachNLRI decode failed with error:", err)
		} else {
			t.Log("BGP MPUnreachNLRI decode succeeded")
		}
	}
}
