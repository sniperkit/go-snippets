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
package dhcp4

import (
	"net"
)

// type is similar to typedef in c
type DhcpOptionCode byte
type OpCode byte
type MessageType byte // Option 53

// Map of DHCP options
type DhcpRelayAgentOptions map[DhcpOptionCode][]byte

// A DHCP packet
type DhcpRelayAgentPacket []byte

var dhcprelayPadder [DHCP_PACKET_MIN_SIZE]byte

type Option struct {
	Code  DhcpOptionCode
	Value []byte
}

/*
 * API to create a new Dhcp packet with Relay Agent information in it
 */
func DhcpRelayAgentCreateNewPacket(opCode OpCode, inReq DhcpRelayAgentPacket) DhcpRelayAgentPacket {
	p := make(DhcpRelayAgentPacket, DHCP_PACKET_MIN_BYTES+1) //241
	p.SetHeaderType(inReq.GetHeaderType())                   // Ethernet
	p.SetCookie(inReq.GetCookie())                           // copy cookie from original pkt
	p.SetOpCode(opCode)                                      // opcode can be request or reply
	p.SetXId(inReq.GetXId())                                 // copy from org pkt
	p.SetFlags(inReq.GetFlags())                             // copy from org pkt
	p.SetYIAddr(inReq.GetYIAddr())                           // copy from org pkt
	p.SetCHAddr(inReq.GetCHAddr())                           // copy from org pkt
	p.SetSecs(inReq.GetSecs())                               // copy from org pkt
	p.SetSName(inReq.GetSName())                             // copy from org pkt
	p.SetFile(inReq.GetFile())                               // copy from org pkt
	p[DHCP_PACKET_MIN_BYTES] = byte(End)                     // set opcode END at the very last
	return p
}

func (p DhcpRelayAgentPacket) GetHeaderLen() byte {
	return p[2]
}

func (p DhcpRelayAgentPacket) GetOpCode() OpCode {
	return OpCode(p[0])
}
func (p DhcpRelayAgentPacket) GetHeaderType() byte {
	return p[1]
}
func (p DhcpRelayAgentPacket) GetHops() byte {
	return p[3]
}
func (p DhcpRelayAgentPacket) GetXId() []byte {
	return p[4:8]
}
func (p DhcpRelayAgentPacket) GetSecs() []byte {
	return p[8:10]
}
func (p DhcpRelayAgentPacket) GetFlags() []byte {
	return p[10:12]
}
func (p DhcpRelayAgentPacket) GetCIAddr() net.IP {
	return net.IP(p[12:16])
}
func (p DhcpRelayAgentPacket) GetYIAddr() net.IP {
	return net.IP(p[16:20])
}
func (p DhcpRelayAgentPacket) GetSIAddr() net.IP {
	return net.IP(p[20:24])
}
func (p DhcpRelayAgentPacket) GetGIAddr() net.IP {
	return net.IP(p[24:28])
}
func (p DhcpRelayAgentPacket) GetCHAddr() net.HardwareAddr {
	hLen := p.GetHeaderLen()
	if hLen > DHCP_PACKET_HEADER_SIZE { // Prevent chaddr exceeding p boundary
		hLen = DHCP_PACKET_HEADER_SIZE
	}
	return net.HardwareAddr(p[28 : 28+hLen]) // max endPos 44
}

func UtiltrimNull(d []byte) []byte {
	for i, v := range d {
		if v == 0 {
			return d[:i]
		}
	}
	return d
}
func (p DhcpRelayAgentPacket) GetCookie() []byte {
	return p[236:240]
}

// BOOTP legacy
func (p DhcpRelayAgentPacket) GetSName() []byte {
	return UtiltrimNull(p[44:108])
}

// BOOTP legacy
func (p DhcpRelayAgentPacket) GetFile() []byte {
	return UtiltrimNull(p[108:236])
}

func ParseMessageTypeToString(mtype MessageType) string {
	switch mtype {
	case 1:
		Logger.Debug("DRA: Message Type: DhcpDiscover")
		return "DHCPDISCOVER"
	case 2:
		Logger.Debug("DRA: Message Type: DhcpOffer")
		return "DHCPOFFER"
	case 3:
		Logger.Debug("DRA: Message Type: DhcpRequest")
		return "DHCPREQUEST"
	case 4:
		Logger.Debug("DRA: Message Type: DhcpDecline")
		return "DHCPDECLINE"
	case 5:
		Logger.Debug("DRA: Message Type: DhcpACK")
		return "DHCPACK"
	case 6:
		Logger.Debug("DRA: Message Type: DhcpNAK")
		return "DHCPNAK"
	case 7:
		Logger.Debug("DRA: Message Type: DhcpRelease")
		return "DHCPRELEASE"
	case 8:
		Logger.Debug("DRA: Message Type: DhcpInform")
		return "DHCPINFORM"
	default:
		Logger.Debug("DRA: Message Type: UnKnown...Discard the Packet")
		return "UNKNOWN REQUEST TYPE"
	}
}

/*
 * ========================SET API's FOR ABOVE MESSAGE FORMAT==================
 */
func (p DhcpRelayAgentPacket) SetOpCode(c OpCode) {
	p[0] = byte(c)
}

func (p DhcpRelayAgentPacket) SetCHAddr(a net.HardwareAddr) {
	copy(p[28:44], a)
	p[2] = byte(len(a))
}

func (p DhcpRelayAgentPacket) SetHeaderType(hType byte) {
	p[1] = hType
}

func (p DhcpRelayAgentPacket) SetCookie(cookie []byte) {
	copy(p.GetCookie(), cookie)
}

func (p DhcpRelayAgentPacket) SetHops(hops byte) {
	p[3] = hops
}

func (p DhcpRelayAgentPacket) SetXId(xId []byte) {
	copy(p.GetXId(), xId)
}

func (p DhcpRelayAgentPacket) SetSecs(secs []byte) {
	copy(p.GetSecs(), secs)
}

func (p DhcpRelayAgentPacket) SetFlags(flags []byte) {
	copy(p.GetFlags(), flags)
}

func (p DhcpRelayAgentPacket) SetCIAddr(ip net.IP) {
	copy(p.GetCIAddr(), ip.To4())
}

func (p DhcpRelayAgentPacket) SetYIAddr(ip net.IP) {
	copy(p.GetYIAddr(), ip.To4())
}

func (p DhcpRelayAgentPacket) SetSIAddr(ip net.IP) {
	copy(p.GetSIAddr(), ip.To4())
}

func (p DhcpRelayAgentPacket) SetGIAddr(ip net.IP) {
	copy(p.GetGIAddr(), ip.To4())
}

// BOOTP legacy
func (p DhcpRelayAgentPacket) SetSName(sName []byte) {
	copy(p[44:108], sName)
	if len(sName) < 64 {
		p[44+len(sName)] = 0
	}
}

// BOOTP legacy
func (p DhcpRelayAgentPacket) SetFile(file []byte) {
	copy(p[108:236], file)
	if len(file) < 128 {
		p[108+len(file)] = 0
	}
}

func (p DhcpRelayAgentPacket) AllocateOptions() []byte {
	if len(p) > DHCP_PACKET_MIN_BYTES {
		return p[DHCP_PACKET_MIN_BYTES:]
	}
	return nil
}

func (p *DhcpRelayAgentPacket) PadToMinSize() {
	sizeofPacket := len(*p)
	if sizeofPacket < DHCP_PACKET_MIN_SIZE {
		// adding whatever is left out to the padder
		*p = append(*p, dhcprelayPadder[:DHCP_PACKET_MIN_SIZE-sizeofPacket]...)
	}
}

// Parses the packet's options into an Options map
func (p DhcpRelayAgentPacket) ParseDhcpOptions() DhcpRelayAgentOptions {
	opts := p.AllocateOptions()
	// create basic dhcp options...
	doptions := make(DhcpRelayAgentOptions, 15)
	for len(opts) >= 2 && DhcpOptionCode(opts[0]) != End {
		if DhcpOptionCode(opts[0]) == Pad {
			opts = opts[1:]
			continue
		}
		size := int(opts[1])
		if len(opts) < 2+size {
			break
		}
		doptions[DhcpOptionCode(opts[0])] = opts[2 : 2+size]
		opts = opts[2+size:]
	}
	return doptions
}

// Appends a DHCP option to the end of a packet
func (p *DhcpRelayAgentPacket) AddDhcpOptions(op DhcpOptionCode, value []byte) {
	// Strip off End, Add OptionCode and Length
	*p = append((*p)[:len(*p)-1], []byte{byte(op), byte(len(value))}...)
	*p = append(*p, value...)  // Add Option Value
	*p = append(*p, byte(End)) // Add on new End
}

// SelectOrder returns a slice of options ordered and selected by a byte array
// usually defined by OptionParameterRequestList.  This result is expected to be
// used in ReplyPacket()'s []Option parameter.
func (o DhcpRelayAgentOptions) SelectOrder(order []byte) []Option {
	opts := make([]Option, 0, len(order))
	for _, v := range order {
		if data, ok := o[DhcpOptionCode(v)]; ok {
			opts = append(opts, Option{Code: DhcpOptionCode(v),
				Value: data})
		}
	}
	return opts
}

// SelectOrderOrAll has same functionality as SelectOrder, except if the order
// param is nil, whereby all options are added (in arbitary order).
func (o DhcpRelayAgentOptions) SelectOrderOrAll(order []byte) []Option {
	if order == nil {
		opts := make([]Option, 0, len(o))
		for i, v := range o {
			opts = append(opts, Option{Code: i, Value: v})
		}
		return opts
	}
	return o.SelectOrder(order)
}

/*========================= END OF HELPER FUNCTION ===========================*/
/*
 * APT to decode incoming Packet by converting the byte into DHCP packet format
 */
func DhcpRelayAgentDecodeInPkt(data []byte, bytesRead int) (DhcpRelayAgentPacket,
	DhcpRelayAgentOptions, MessageType) {
	inRequest := DhcpRelayAgentPacket(data[:bytesRead])
	if inRequest.GetHeaderLen() > DHCP_PACKET_HEADER_SIZE {
		Logger.Warning("Header Lenght is invalid... don't do anything")
		return nil, nil, 0
	}
	reqOptions := inRequest.ParseDhcpOptions()
	/*
		logger.Debug("DRA: CIAddr is " + inRequest.GetCIAddr().String())
		logger.Debug("DRA: CHaddr is " + inRequest.GetCHAddr().String())
		logger.Debug("DRA: YIAddr is " + inRequest.GetYIAddr().String())
		logger.Debug("DRA: GIAddr is " + inRequest.GetGIAddr().String())

		logger.Debug("DRA: Cookie is ", inRequest.GetCookie())
	*/
	mType := reqOptions[OptionDHCPMessageType]
	return inRequest, reqOptions, MessageType(mType[0])
}
