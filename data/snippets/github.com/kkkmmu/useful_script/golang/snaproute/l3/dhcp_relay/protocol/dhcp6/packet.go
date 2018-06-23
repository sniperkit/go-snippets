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

package dhcp6

import (
	"encoding/binary"
	"net"
)

type DhcpPacket []byte
type DhcpRelayPacket []byte
type MsgType uint8
type OptionType uint16

type DhcpOption []byte
type DhcpOptionMap map[OptionType]DhcpOption

type OptionIANA []byte
type OptionIATA []byte
type OptionIAAddr []byte
type Duid []byte

// CLIENT/SERVER Packet
func (p DhcpPacket) GetMsgType() uint8 {
	return uint8(p[0])
}

func (p DhcpPacket) GetOptionsField() []byte {
	return p[4:]
}

func (p DhcpPacket) Validate() bool {
	if len(p) < 5 {
		return false
	}
	if MsgType(p.GetMsgType()) < SOLICIT ||
		MsgType(p.GetMsgType()) > RELAY_REPL {

		return false
	}
	return true
}

func filterOption(filter []OptionType, opCode OptionType) bool {
	if filter == nil || len(filter) == 0 {
		return true
	}
	for _, curOpCode := range filter {
		if curOpCode == opCode {
			return true
		}
	}
	return false
}

func (p DhcpPacket) ParseOptions(filter []OptionType) DhcpOptionMap {
	dOpt := DhcpOption(p.GetOptionsField())
	optionMap := make(DhcpOptionMap, 0)
	for {
		if !dOpt.Validate() {
			break
		}
		opCode := OptionType(dOpt.GetCode())
		if filterOption(filter, opCode) {
			optionMap[opCode] = dOpt
		}
		dOpt = DhcpOption(dOpt[4+int(dOpt.GetLen()):])
	}
	return optionMap
}

func (p DhcpPacket) GetDuidField() []byte {
	i := 4
	var dOpt DhcpOption
	for {
		if i >= len(p) {
			break
		}
		dOpt = DhcpOption(p[i:])
		if !dOpt.Validate() {
			break
		}
		opCode := dOpt.GetCode()
		if OptionType(opCode) == OPTION_CLIENTID {
			return dOpt.GetPayload()
		}
		i = i + 4 + int(dOpt.GetLen())
	}
	return nil
}

func (p DhcpPacket) GetClientHwAddr() net.HardwareAddr {
	var duidField []byte
	switch MsgType(p.GetMsgType()) {
	case SOLICIT, REQUEST, CONFIRM, RENEW, REBIND, DECLINE, RELEASE,
		INFO_REQ, ADVERTISE, REPLY, RECONFIGURE:

		duidField = p.GetDuidField()
	default:
		duidField = nil
	}
	if duidField == nil {
		return nil
	}
	macAddrField := Duid(duidField).GetLinkLayerAddr()
	if macAddrField == nil {
		return nil
	}
	return net.HardwareAddr(macAddrField)
}

// RELAY-FORW/RELAY-REPL Packet
func (p DhcpRelayPacket) GetMsgType() uint8 {
	return uint8(p[0])
}

func (p DhcpRelayPacket) GetHopCount() uint8 {
	return uint8(p[1])
}

func (p DhcpRelayPacket) GetLinkAddrField() []byte {
	return p[2:18]
}

func (p DhcpRelayPacket) GetPeerAddrField() []byte {
	return p[18:34]
}

func (p DhcpRelayPacket) GetOptionsField() []byte {
	return p[34:]
}

func (p DhcpRelayPacket) GetDRAOptionField() DhcpOption {
	i := 34
	var dOpt DhcpOption
	for { // Iterate over Options
		if i >= len(p) {
			break
		}
		dOpt = DhcpOption(p[i:])
		if !dOpt.Validate() {
			break
		}
		opCode := dOpt.GetCode()
		if opCode == uint16(OPTION_RELAY_MSG) {
			return dOpt
		}
		i = i + 4 + int(dOpt.GetLen())
	}
	return nil
}

func (p DhcpRelayPacket) AddRelayMsgOption(payload []byte) DhcpRelayPacket {
	newPkt := p[:len(p)+4+len(payload)]
	dhcpOp := DhcpOption(newPkt[len(p):])
	dhcpOp.SetCodeField(uint16(OPTION_RELAY_MSG))
	dhcpOp.SetPayload(payload)
	return newPkt
}

//func (p DhcpRelayPacket) GetDuidField() []byte {
//	var draOpt DhcpOption
//	var tmpBuf []byte = p
//	for {
//		draOpt = DhcpRelayPacket(tmpBuf).GetDRAOptionField()
//		if draOpt == nil {
//			break
//		}
//		if !draOpt.Validate() {
//			break
//		}
//		tmpBuf = draOpt.GetPayload()
//		if MsgType(tmpBuf[0]) >= SOLICIT &&
//			MsgType(tmpBuf[0]) <= INFO_REQ {

//			return DhcpPacket(tmpBuf).GetDuidField()
//		} else if MsgType(tmpBuf[0]) == RELAY_FORW ||
//			MsgType(tmpBuf[0]) == RELAY_REPL {

//			continue
//		} else {
//			break
//		}
//	}
//	return nil
//}

// (NOTE): Decapsulates DRA Packet recursively to get DHCP Packet
//func (p DhcpRelayPacket) GetDhcpPacket() []byte {
//	var draOpt DhcpOption
//	var tmpBuf []byte = p
//	for {
//		draOpt = DhcpRelayPacket(tmpBuf).GetDRAOptionField()
//		if draOpt == nil {
//			break
//		}
//		if !draOpt.Validate() {
//			break
//		}
//		tmpBuf = draOpt.GetPayload()
//		if MsgType(tmpBuf[0]) >= SOLICIT &&
//			MsgType(tmpBuf[0]) <= INFO_REQ {
//
//			return tmpBuf
//		} else if MsgType(tmpBuf[0]) == RELAY_FORW ||
//			MsgType(tmpBuf[0]) == RELAY_REPL {
//
//			continue
//		} else {
//			break
//		}
//	}
//	return nil
//}

func (p DhcpRelayPacket) Validate() bool {
	if len(p) < 34 {
		return false
	}
	if MsgType(p.GetMsgType()) == RELAY_FORW ||
		MsgType(p.GetMsgType()) == RELAY_REPL {

		return true
	} else {
		return false
	}
}

func (p DhcpRelayPacket) SetMsgType(mType int8) {
	p[0] = byte(mType)
}

func (p DhcpRelayPacket) SetHopCount(cnt uint8) {
	p[1] = byte(cnt)
}

func (p DhcpRelayPacket) SetLinkAddrField(addr []byte) {
	copy(p.GetLinkAddrField(), addr)
}

func (p DhcpRelayPacket) SetPeerAddrField(addr []byte) {
	copy(p.GetPeerAddrField(), addr)
}

func (p DhcpRelayPacket) AddDRAOption(dOpt DhcpOption) {
	copy(p[34:], dOpt[:4+dOpt.GetLen()])
}

// DHCP OPTION
func (p DhcpOption) GetCodeField() []byte {
	return p[0:2]
}

func (p DhcpOption) GetLenField() []byte {
	return p[2:4]
}

func (p DhcpOption) GetCode() uint16 {
	return binary.BigEndian.Uint16(p[0:2])
}

func (p DhcpOption) GetLen() uint16 {
	return binary.BigEndian.Uint16(p[2:4])
}

func (p DhcpOption) GetPayload() []byte {
	return p[4 : 4+p.GetLen()]
}

func (p DhcpOption) Validate() bool {
	if len(p) < 4 {
		return false
	}
	if uint16(len(p)-4) < p.GetLen() {
		return false
	}
	if p.GetCode() == 0 {
		return false
	}
	return true
}

func (p DhcpOption) SetCodeField(code uint16) {
	binary.BigEndian.PutUint16(p.GetCodeField(), code)
}

func (p DhcpOption) SetPayload(payload []byte) {
	payloadLen := uint16(len(payload))
	binary.BigEndian.PutUint16(p.GetLenField(), payloadLen)
	copy(p[4:], payload)
}

// DUID
// (TODO): Add bounds validation
func (p Duid) GetLinkLayerAddr() []byte {
	const DUID_LLT uint16 = 1
	const DUID_EN uint16 = 2
	const DUID_LL uint16 = 3
	duidType := binary.BigEndian.Uint16(p[0:2])
	switch duidType {
	case DUID_LLT:
		hwType := binary.BigEndian.Uint16(p[2:4])
		if hwType == 1 { // Ethernet
			return p[8:len(p)]
		} else {
			return nil
		}
	case DUID_EN:
		return nil
	case DUID_LL:
		hwType := binary.BigEndian.Uint16(p[2:4])
		if hwType == 1 { // Ethernet
			return p[4:len(p)]
		} else {
			return nil
		}
	default:
		return nil
	}
}

func getOptionIAAddr(p []byte) OptionIAAddr {
	dOpt := DhcpOption(p)
	for {
		if !dOpt.Validate() {
			break
		}
		if OptionType(dOpt.GetCode()) == OPTION_IAADDR {
			return OptionIAAddr(dOpt)
		}
		dOpt = DhcpOption(dOpt[4+int(dOpt.GetLen()):])
	}
	return nil
}

func (p OptionIANA) GetOptions() []byte {
	return p[16:]
}

func (p OptionIANA) GetIpAddr() []byte {
	if !DhcpOption(p).Validate() {
		return nil
	}
	iaAddrOpt := getOptionIAAddr(p.GetOptions())
	if iaAddrOpt == nil {
		return nil
	}
	if !DhcpOption(iaAddrOpt).Validate() {
		return nil
	}
	return iaAddrOpt.GetIpAddr()
}

func (p OptionIATA) GetOptions() []byte {
	return p[8:]
}

func (p OptionIATA) GetIpAddr() []byte {
	if !DhcpOption(p).Validate() {
		return nil
	}
	iaAddrOpt := getOptionIAAddr(p.GetOptions())
	if iaAddrOpt == nil {
		return nil
	}
	if !DhcpOption(iaAddrOpt).Validate() {
		return nil
	}
	return iaAddrOpt.GetIpAddr()
}

func (p OptionIAAddr) GetIpAddr() []byte {
	return p[4:20]
}
