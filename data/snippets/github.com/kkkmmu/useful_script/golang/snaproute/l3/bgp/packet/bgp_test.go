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
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"l3/bgp/utils"
	"math"
	"net"
	"testing"
	"utils/logging"
)

func TestBGPUpdatePacketsSliceBoundOutOfRange(t *testing.T) {
	strPkts := make([]string, 0)
	strPkts = append(strPkts, "000000204001010140020602011908b10a4003040a0a00c280040400000000c01c000100080a")
	strPkts = append(strPkts, "0000002a400101015002feff02011908b10a4003040a0a00c280040400000000800e0b000101040a0a00c200080a")
	strPkts = append(strPkts, "00000401c0fcfe020080000008800000108000001880000020800000288000003080000038800000408000004880000050800000588000006080000068800000708000"+
		"007880000080800000888000009080000098800000a0800000a8800000b0800000b8800000c0800000c8800000d0800000d8800000e0800000e8800000f0800000f8800001008000010880000110800001188000012080000128800001308000013880000140800001488000015"+
		"0800001588000016080000168800001708000017880000180800001888000019080000198800001a0800001a8800001b0800001b8800001c0800001c8800001d0800001d8800001e0800001e8800001f0800001f880000200800002088000021080000218800002208000022880"+
		"00023080000238800002408000024880000250800002588000026080000268800002708000027880000280800002888000029080000298800002a0800002a8800002b0800002b8800002c0800002c8800002d0800002d8800002e0800002e8800002f0800002f88000030080000"+
		"308800003108000031880000320800003288000033080000338800003408000034880000350800003588000036080000368800003708000037880000380800003888000039080000398800003a0800003a8800003b0800003b8800003c0800003c8800003d0800003d8800003e0"+
		"800003e8800003f0800003f88000040080000408800004108000041880000420800004288000043080000438800004408000044880000450800004588000046080000468800004708000047880000480800004888000049080000498800004a0800004a8800004b0800004b8800"+
		"004c0800004c8800004d0800004d8800004e0800004e8800004f0800004f880000500800005088000051080000518800005208000052880000530800005388000054080000548800005508000055880000560800005688000057080000578800005808000058880000590800005"+
		"98800005a0800005a8800005b0800005b8800005c0800005c8800005d0800005d8800005e0800005e8800005f0800005f88000060080000608800006108000061880000620800006288000063080000638800006408000064880000650800006588000066080000668800006708"+
		"000067880000680800006888000069080000698800006a0800006a8800006b0800006b8800006c0800006c8800006d0800006d8800006e0800006e8800006f0800006f8800007008000070880000710800007188000072080000728800007308000073880000740800007488000"+
		"0750800007588000076080000768800007708000077880000780800007888000079080000798800007a0800007a8800007b0800007b8800007c0800007c8800007d0800007d8800007e0800007e8800007f0800007f8")
	strPkts = append(strPkts, "0000000940022a02011908b10a")
	strPkts = append(strPkts, "0000042b4001010140020602011908b10ad011000802008000000880000010800000188000002080000028800000308000003880000040800000488000005080000058"+
		"8000006080000068800000708000007880000080800000888000009080000098800000a0800000a8800000b0800000b8800000c0800000c8800000d0800000d8800000e0800000e8800000f0800000f880000100800001088000011080000118800001208000012880000130800"+
		"00138800001408000014880000150800001588000016080000168800001708000017880000180800001888000019080000198800001a0800001a8800001b0800001b8800001c0800001c8800001d0800001d8800001e0800001e8800001f0800001f88000020080000208800002"+
		"108000021880000220800002288000023080000238800002408000024880000250800002588000026080000268800002708000027880000280800002888000029080000298800002a0800002a8800002b0800002b8800002c0800002c8800002d0800002d8800002e0800002e88"+
		"00002f0800002f88000030080000308800003108000031880000320800003288000033080000338800003408000034880000350800003588000036080000368800003708000037880000380800003888000039080000398800003a0800003a8800003b0800003b8800003c08000"+
		"03c8800003d0800003d8800003e0800003e8800003f0800003f88000040080000408800004108000041880000420800004288000043080000438800004408000044880000450800004588000046080000468800004708000047880000480800004888000049080000498800004a"+
		"0800004a8800004b0800004b8800004c0800004c8800004d0800004d8800004e0800004e8800004f0800004f88000050080000508800005108000051880000520800005288000053080000538800005408000054880000550800005588000056080000568800005708000057880"+
		"000580800005888000059080000598800005a0800005a8800005b0800005b8800005c0800005c8800005d0800005d8800005e0800005e8800005f0800005f8800006008000060880000610800006188000062080000628800006308000063880000640800006488000065080000"+
		"6588000066080000668800006708000067880000680800006888000069080000698800006a0800006a8800006b0800006b8800006c0800006c8800006d0800006d8800006e0800006e8800006f0800006f880000700800007088000071080000718800007208000072880000730"+
		"80000738800007408000074880000750800007588000076080000768800007708000077880000780800007888000079080000798800007a0800007a8800007b0800007b8800007c0800007c8800007d0800007d8800007e0800007e8800007f0800007f84003040a0a00c280040"+
		"400000000800e0b000101040a0a00c200080a")
	strPkts = append(strPkts, "000000204001010140020602011908b10a4003040a0a00c280040400000000c01c000100080a")
	for _, strPkt := range strPkts {
		hexPkt, err := hex.DecodeString(strPkt)
		fmt.Printf("packet = %x, len = %d\n", hexPkt, len(hexPkt))
		if err != nil {
			t.Fatal("Failed to decode the string to hex, string =", strPkt)
		}

		if len(hexPkt) > math.MaxUint16 {
			t.Fatal("Length of packet exceeded MAX size, packet len =", len(hexPkt))
		}

		pktLen := make([]byte, 2)
		binary.BigEndian.PutUint16(pktLen, uint16(len(hexPkt)+19))
		header := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x02}
		copy(header[16:18], pktLen)
		fmt.Printf("packet header = %x, len = %d\n", header, len(header))

		bgpHeader := NewBGPHeader()
		err = bgpHeader.Decode(header)
		if err != nil {
			t.Fatal("BGP packet header decode failed with error", err)
		}

		peerAttrs := BGPPeerAttrs{
			ASSize:           4,
			AddPathsRxActual: true,
		}
		bgpMessage := NewBGPMessage()
		err = bgpMessage.Decode(bgpHeader, hexPkt, peerAttrs)
		if err == nil {
			t.Fatal("BGP update message decode called... expected failure, got NO error")
		} else {
			t.Log("BGP update message decode called... expected failure, error:", err)
		}
	}
}

func TestBGPOpenPacketsIndexOutOfRange(t *testing.T) {
	strPkts := make([]string, 0)
	strPkts = append(strPkts, "045ba0000a0a0a00c21d020682070001010101f003020601040001000102020200020440020000")
	for _, strPkt := range strPkts {
		hexPkt, err := hex.DecodeString(strPkt)
		fmt.Printf("packet = %x, len = %d\n", hexPkt, len(hexPkt))
		if err != nil {
			t.Fatal("Failed to decode the string to hex, string =", strPkt)
		}

		if len(hexPkt) > math.MaxUint16 {
			t.Fatal("Length of packet exceeded MAX size, packet len =", len(hexPkt))
		}

		pktLen := make([]byte, 2)
		binary.BigEndian.PutUint16(pktLen, uint16(len(hexPkt)+19))
		header := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x01}
		copy(header[16:18], pktLen)
		fmt.Printf("packet header = %x, len = %d\n", header, len(header))

		bgpHeader := NewBGPHeader()
		err = bgpHeader.Decode(header)
		if err != nil {
			t.Fatal("BGP packet header decode failed with error", err)
		}

		peerAttrs := BGPPeerAttrs{
			ASSize:           2,
			AddPathsRxActual: false,
		}
		bgpMessage := NewBGPMessage()
		err = bgpMessage.Decode(bgpHeader, hexPkt, peerAttrs)
		if err == nil {
			t.Fatal("BGP open message decode called... expected failure, got NO error")
		} else {
			t.Log("BGP open message decode called... expected failure, error:", err)
		}
	}
}

func TestBGPUpdatePathAttrsBadFlags(t *testing.T) {
	strPkts := make([]string, 0)
	strPkts = append(strPkts, "0000001c40010100100200060201000002584003045a01010280040400000000183c010118500101184701011846010218460101183c0102")
	strPkts = append(strPkts, "0000001c40010100500200060201000002582003045a01010280040400000000183c010118500101184701011846010218460101183c0102")
	strPkts = append(strPkts, "0000001c40010100500200060201000002584003045a010102A0040400000000183c010118500101184701011846010218460101183c0102")

	pktPathAttrs := "0000002040010100500200060201000002584003045a01010280040400000000"
	nlri := "183c010118500101184701011846010218460101183c0102"
	pathAttrs := []string{"00000100", "20000100", "60000100", "A0000100"}
	for _, pa := range pathAttrs {
		pa = pa[:2] + fmt.Sprintf("%02x", BGPPathAttrTypeUnknown) + pa[4:]
		strPkts = append(strPkts, pktPathAttrs+pa+nlri)
	}

	for _, strPkt := range strPkts {
		hexPkt, err := hex.DecodeString(strPkt)
		fmt.Printf("packet = %x, len = %d\n", hexPkt, len(hexPkt))
		if err != nil {
			t.Fatal("Failed to decode the string to hex, string =", strPkt)
		}

		if len(hexPkt) > math.MaxUint16 {
			t.Fatal("Length of packet exceeded MAX size, packet len =", len(hexPkt))
		}

		pktLen := make([]byte, 2)
		binary.BigEndian.PutUint16(pktLen, uint16(len(hexPkt)+19))
		header := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x02}
		copy(header[16:18], pktLen)
		fmt.Printf("packet header = %x, len = %d\n", header, len(header))

		bgpHeader := NewBGPHeader()
		err = bgpHeader.Decode(header)
		if err != nil {
			t.Fatal("BGP packet header decode failed with error", err)
		}

		peerAttrs := BGPPeerAttrs{
			ASSize:           4,
			AddPathsRxActual: true,
		}
		bgpMessage := NewBGPMessage()
		err = bgpMessage.Decode(bgpHeader, hexPkt, peerAttrs)
		if err == nil {
			t.Error("BGP update message decode called... expected failure, got NO error")
		} else {
			t.Log("BGP update message decode called... expected failure, error:", err)
		}
	}
}

func TestBGPUpdatePathAttrsBadLength(t *testing.T) {
	strPkts := make([]string, 0)

	pktPathAttrs := "0000001c40010100500200060201000002584003045a01010280040400000000"
	nlri := "183c010118500101184701011846010218460101183c0102"
	pathAttrs := []string{"80000100"}
	for _, pa := range pathAttrs {
		pa = pa[:2] + fmt.Sprintf("%02x", BGPPathAttrTypeUnknown) + pa[4:]
		strPkts = append(strPkts, pktPathAttrs+pa+nlri)
	}

	for _, strPkt := range strPkts {
		hexPkt, err := hex.DecodeString(strPkt)
		fmt.Printf("packet = %x, len = %d\n", hexPkt, len(hexPkt))
		if err != nil {
			t.Fatal("Failed to decode the string to hex, string =", strPkt)
		}

		if len(hexPkt) > math.MaxUint16 {
			t.Fatal("Length of packet exceeded MAX size, packet len =", len(hexPkt))
		}

		pktLen := make([]byte, 2)
		binary.BigEndian.PutUint16(pktLen, uint16(len(hexPkt)+19))
		header := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x02}
		copy(header[16:18], pktLen)
		fmt.Printf("packet header = %x, len = %d\n", header, len(header))

		bgpHeader := NewBGPHeader()
		err = bgpHeader.Decode(header)
		if err != nil {
			t.Fatal("BGP packet header decode failed with error", err)
		}

		peerAttrs := BGPPeerAttrs{
			ASSize:           4,
			AddPathsRxActual: false,
		}
		bgpMessage := NewBGPMessage()
		err = bgpMessage.Decode(bgpHeader, hexPkt, peerAttrs)
		if err == nil {
			t.Error("BGP update message decode called... expected failure, got NO error")
		} else {
			t.Log("BGP update message decode called... expected failure, error:", err)
		}
	}
}

func TestBGPUpdatePacketDecode(t *testing.T) {
	strPkts := make([]string, 0)
	// With base Path attrs - ORIGIN, AS_PATH (4 byte), NEXT_HOP, MULTI_EXIT_DISC
	strPkts = append(strPkts, "0000001b4001010140020602011908b10a4003040a0a00c28004040000000000000001080a")
	// Added path attrs - LOCAL_PREF, ATOMIC_AGGREGATE
	strPkts = append(strPkts, "000000254001010140020602011908b10a4003040a0a00c2800404000000004005040102030440060000000001080a")

	// Added path attrs - AGGREGATOR (4 byte AS)
	strPkts = append(strPkts, "000000304001010140020602011908b10a4003040a0a00c28004040000000040050401020304400600C007081908b10b0a010a1c00000001080a")

	// Added path attrs - ORIGINATOR_ID, CLUSTER_LIST
	strPkts = append(strPkts, "000000334001010140020602011908b10a4003040a0a00c280040400000000400504010203044006008009040a010a32800A040102030400000001080a")

	// Added path attrs - MP_REACH_NLRI, MP_UNREACH_NLRI
	strPkts = append(strPkts, "000000474001010140020602011908b10a4003040a0a00c280040400000000800E1C000201100102030405060708091011121314151600000000020A0A80800F0A000201000000030A0BC000000001080a")

	for _, strPkt := range strPkts {
		hexPkt, err := hex.DecodeString(strPkt)
		fmt.Printf("packet = %x, len = %d\n", hexPkt, len(hexPkt))
		if err != nil {
			t.Fatal("Failed to decode the string to hex, string =", strPkt)
		}

		if len(hexPkt) > math.MaxUint16 {
			t.Fatal("Length of packet exceeded MAX size, packet len =", len(hexPkt))
		}

		pktLen := make([]byte, 2)
		binary.BigEndian.PutUint16(pktLen, uint16(len(hexPkt)+19))
		header := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x02}
		copy(header[16:18], pktLen)
		fmt.Printf("packet header = %x, len = %d\n", header, len(header))

		bgpHeader := NewBGPHeader()
		err = bgpHeader.Decode(header)
		if err != nil {
			t.Fatal("BGP packet header decode failed with error", err)
		}

		peerAttrs := BGPPeerAttrs{
			ASSize:           4,
			AddPathsRxActual: true,
		}
		bgpMessage := NewBGPMessage()
		err = bgpMessage.Decode(bgpHeader, hexPkt, peerAttrs)
		if err != nil {
			t.Fatal("BGP update message decode failed with error:", err)
		} else {
			t.Log("BGP update message decode succeeded")
		}
	}
}

func TestBGPUpdatePacketDecode2ByteAS(t *testing.T) {
	logger, err := logging.NewLogger("bgpd", "BGP", true)
	if err != nil {
		t.Fatal("Failed to start the logger. Exiting!!")
	}
	utils.SetLogger(logger)

	strPkts := make([]string, 0)
	// With Path attrs - ORIGIN, AS_PATH (2 byte), NEXT_HOP, AS4_PATH
	//strPkts = append(strPkts, "000000214001010140020602025BA0010a4003040a0a00c2C0110A02021908b10a0000010a00000001080a")
	strPkts = append(strPkts, "000000214001010140020602025BA0010a4003040a0a00c2C0110A02021908b10a0000010a080a")
	// Added path attrs - AGGREGATOR (2 byte AS), AS4_AGGREGATOR
	//strPkts = append(strPkts, "000000354001010140020602025BA0010a4003040a0a00c2C0110A02021908b10a0000010aC007065BA00a010a1cC012081908b10b0a010a1c00000001080a")
	strPkts = append(strPkts, "000000354001010140020602025BA0010a4003040a0a00c2C0110A02021908b10a0000010aC007065BA00a010a1cC012081908b10b0a010a1c080a")
	for _, strPkt := range strPkts {
		hexPkt, err := hex.DecodeString(strPkt)
		fmt.Printf("packet = %x, len = %d\n", hexPkt, len(hexPkt))
		if err != nil {
			t.Fatal("Failed to decode the string to hex, string =", strPkt)
		}

		if len(hexPkt) > math.MaxUint16 {
			t.Fatal("Length of packet exceeded MAX size, packet len =", len(hexPkt))
		}

		pktLen := make([]byte, 2)
		binary.BigEndian.PutUint16(pktLen, uint16(len(hexPkt)+19))
		header := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x02}
		copy(header[16:18], pktLen)
		fmt.Printf("packet header = %x, len = %d\n", header, len(header))

		bgpHeader := NewBGPHeader()
		err = bgpHeader.Decode(header)
		if err != nil {
			t.Fatal("BGP packet header decode failed with error", err)
		}

		peerAttrs := BGPPeerAttrs{
			ASSize:           2,
			AddPathsRxActual: false,
		}
		bgpMessage := NewBGPMessage()
		err = bgpMessage.Decode(bgpHeader, hexPkt, peerAttrs)
		if err != nil {
			t.Fatal("BGP update message decode failed with error:", err)
		} else {
			t.Log("BGP update message decode succeeded")
		}
	}
}

func TestBGPUpdateEncode(t *testing.T) {
	ip := net.ParseIP("10.1.0.0")
	prefix := NewIPPrefix(ip, 16)
	extNLRI := NewExtNLRI(1, prefix)
	nlri := make([]NLRI, 0)
	nlri = append(nlri, extNLRI)

	pa := make([]BGPPathAttr, 0)

	origin := NewBGPPathAttrOrigin(BGPPathAttrOriginIncomplete)
	pa = append(pa, origin)

	asPathSeq := NewBGPAS4PathSegmentSeq()
	asPathSeq.AppendAS(1)
	asPathSeq.PrependAS(2)
	asPathSet := NewBGPAS4PathSegmentSet()
	asPathSet.AppendAS(1)
	asPathSet.PrependAS(2)
	asPath := NewBGPPathAttrASPath()
	asPath.PrependASPathSegment(asPathSeq)
	asPath.AppendASPathSegment(asPathSet)
	pa = append(pa, asPath)

	nextHop := NewBGPPathAttrNextHop()
	nextHop.Value = net.ParseIP("10.1.10.1")
	pa = append(pa, nextHop)

	med := NewBGPPathAttrMultiExitDisc()
	med.Value = 1
	pa = append(pa, med)

	localPref := NewBGPPathAttrLocalPref()
	localPref.Value = 100
	pa = append(pa, localPref)

	atomicAgg := NewBGPPathAttrAtomicAggregate()
	pa = append(pa, atomicAgg)

	aggregator := NewBGPPathAttrAggregator()
	aggAS := NewBGPAggregator4ByteAS()
	aggAS.AS = 200
	aggregator.SetBGPAggregatorAS(aggAS)
	aggregator.IP = net.ParseIP("20.1.20.1")
	pa = append(pa, aggregator)

	originatorId := NewBGPPathAttrOriginatorId(net.ParseIP("30.1.30.1"))
	pa = append(pa, originatorId)

	clusterList := NewBGPPathAttrClusterList()
	clusterList.PrependId(1234)
	pa = append(pa, clusterList)

	mpReachNLRI := NewBGPPathAttrMPReachNLRI()
	mpReachNLRI.AFI = AfiIP6
	mpReachNLRI.SAFI = SafiUnicast
	mpNextHop := NewMPNextHopIP()
	mpNextHop.SetNextHop(net.ParseIP("2001::1"))
	mpReachNLRI.SetNextHop(mpNextHop)
	mpIP := NewIPPrefix(net.ParseIP("1.2.0.0"), 16)
	mpReachNLRI.AddNLRI(mpIP)
	pa = append(pa, mpReachNLRI)

	mpUnreachNLRI := NewBGPPathAttrMPUnreachNLRI()
	mpUnreachNLRI.AFI = AfiIP6
	mpUnreachNLRI.SAFI = SafiUnicast
	mpUnreachIP := NewIPPrefix(net.ParseIP("3.4.0.0"), 16)
	mpUnreachNLRI.AddNLRI(mpUnreachIP)
	pa = append(pa, mpUnreachNLRI)

	updateMsg := NewBGPUpdateMessage(make([]NLRI, 0), pa, nlri)
	pkt, err := updateMsg.Encode()
	if err != nil {
		t.Fatal("BGP update message encode failed with error:", err)
	}
	t.Log("BGP update message:", pkt)

	newUpdateMsg := updateMsg.Clone()
	newPkt, err := newUpdateMsg.Encode()
	if err != nil {
		t.Fatal("Cloned BGP update message encode failed with error:", err)
	}
	t.Log("Cloned BGP update message:", newPkt)

	if !bytes.Equal(pkt, newPkt) {
		t.Fatal("Cloned update message is not the same as the original message")
	}
}
