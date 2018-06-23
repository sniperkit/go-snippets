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
	"encoding/binary"
	_ "encoding/json"
	"errors"
	_ "fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"l2/lldp/config"
	"l2/lldp/utils"
	"net"
)

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func TxInit(interval, hold int) *TX {
	var err error
	/*  Set tx interval during init or update
	 *  default value is 30
	 *  Set tx hold multiplier during init or update
	 *  default value is 4
	 */
	txInfo := &TX{
		MessageTxInterval:       interval,
		MessageTxHoldMultiplier: hold,
		useCacheFrame:           false,
	}
	/*  Set TTL Value at the time of init or update of lldp config
	 *  default value comes out to be 120
	 */
	txInfo.ttl = Min(LLDP_MAX_TTL, txInfo.MessageTxInterval*
		txInfo.MessageTxHoldMultiplier)
	txInfo.DstMAC, err = net.ParseMAC(LLDP_PROTO_DST_MAC)
	if err != nil {
		debug.Logger.Err("parsing lldp protocol Mac failed", err)
	}

	return txInfo
}

/*  Function to send out lldp frame to peer on timer expiry.
 *  if a cache entry is present then use that otherwise create a new lldp frame
 *  A new frame will be constructed:
 *		1) if it is first time send
 *		2) if there is config object update
 */
func (t *TX) Frame(port config.PortInfo, sysInfo *config.SystemInfo) []byte {
	temp := make([]byte, 0)
	// if cached then directly send the packet
	if t.useCacheFrame {
		return t.cacheFrame
	} else {
		srcmac, _ := net.ParseMAC(port.MacAddr)
		// we need to construct new lldp frame based of the information that we
		// have collected locally
		// Chassis ID: Mac Address of Port
		// Port ID: Port Name
		// TTL: calculated during port init default is 30 * 4 = 120
		payload := t.createPayload(srcmac, port, sysInfo)
		if payload == nil {
			debug.Logger.Err("Creating payload failed for port", port)
			t.useCacheFrame = false
			return temp
		}
		// Additional TLV's... @TODO: get it done later on
		// System information... like "show version" command at Cisco
		// System Capabilites...

		// Construct ethernet information
		eth := &layers.Ethernet{
			SrcMAC:       srcmac,
			DstMAC:       t.DstMAC,
			EthernetType: layers.EthernetTypeLinkLayerDiscovery,
		}

		// construct new buffer
		buffer := gopacket.NewSerializeBuffer()
		options := gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		}
		gopacket.SerializeLayers(buffer, options, eth, gopacket.Payload(payload))
		pkt := buffer.Bytes()
		t.cacheFrame = make([]byte, len(pkt))
		copied := copy(t.cacheFrame, pkt)
		if copied < len(pkt) {
			debug.Logger.Err("Cache cannot be created")
			t.cacheFrame = nil
			t.useCacheFrame = false
			// return should never happen
			return temp
		}
		debug.Logger.Info("Send Frame is cached for Port:", port.Name)
		t.useCacheFrame = true
		return pkt
	}
}

/*  helper function to create payload from lldp frame struct
 */
func (t *TX) createPayload(srcmac []byte, port config.PortInfo, sysInfo *config.SystemInfo) []byte {
	var payload []byte
	var err error
	tlvType := layers.LLDPTLVChassisID // start with chassis id always
	for {
		if tlvType > LLDP_TOTAL_TLV_SUPPORTED { // right now only minimal lldp tlv
			break
		} else if tlvType > layers.LLDPTLVTTL && sysInfo == nil {
			debug.Logger.Debug("Reading System Information from DB failed and hence sending out only " +
				"Mandatory TLV's")
			break
		}
		tlv := &layers.LinkLayerDiscoveryValue{}
		switch tlvType {
		case layers.LLDPTLVChassisID: // Chassis ID
			tlv.Type = layers.LLDPTLVChassisID
			tlv.Value = EncodeMandatoryTLV(byte(layers.LLDPChassisIDSubTypeMACAddr), srcmac)
			debug.Logger.Debug("Chassis id tlv", *tlv)

		case layers.LLDPTLVPortID: // Port ID
			tlv.Type = layers.LLDPTLVPortID
			tlv.Value = EncodeMandatoryTLV(byte(layers.LLDPPortIDSubtypeIfaceName), []byte(port.Name))
			debug.Logger.Debug("Port id tlv", *tlv)

		case layers.LLDPTLVTTL: // TTL
			tlv.Type = layers.LLDPTLVTTL
			tb := []byte{0, 0}
			binary.BigEndian.PutUint16(tb, uint16(t.ttl))
			tlv.Value = append(tlv.Value, tb...)
			debug.Logger.Debug("TTL tlv", *tlv)

		case layers.LLDPTLVPortDescription:
			tlv.Type = layers.LLDPTLVPortDescription
			tlv.Value = []byte(port.Description)
			debug.Logger.Debug("Port Description", *tlv)

		case layers.LLDPTLVSysDescription:
			tlv.Type = layers.LLDPTLVSysDescription
			tlv.Value = []byte(sysInfo.Description)
			debug.Logger.Debug("System Description", *tlv)

		case layers.LLDPTLVSysName:
			tlv.Type = layers.LLDPTLVSysName
			tlv.Value = []byte(sysInfo.Hostname)
			debug.Logger.Debug("System Name", *tlv)

		case layers.LLDPTLVSysCapabilities:
			err = errors.New("Tlv not supported")

		case layers.LLDPTLVMgmtAddress:
			/*
			 *  Value: N bytes
			 *     Subtype is 1 byte
			 *     Address is []byte
			 *     IntefaceSubtype is 1 byte
			 *     IntefaceNumber uint32 <<< this is system interface number which is IfIndex in out case
			 *     OID string

			 */
			tlv.Type = layers.LLDPTLVMgmtAddress
			mgmtInfo := &layers.LLDPMgmtAddress{
				Subtype:          layers.IANAAddressFamilyIPV4,
				Address:          net.ParseIP(sysInfo.MgmtIp).To4(),
				InterfaceSubtype: layers.LLDPInterfaceSubtypeifIndex,
				InterfaceNumber:  uint32(port.IfIndex),
			}
			tlv.Value = EncodeMgmtTLV(mgmtInfo)
		}
		if err == nil {
			tlv.Length = uint16(len(tlv.Value))
			payload = append(payload, EncodeTLV(tlv)...)
		}
		err = nil
		tlvType++
	}

	// add organization specific tlv for pvid
	tlv := &layers.LinkLayerDiscoveryValue{}
	tlv.Type = layers.LLDPTLVOrgSpecific
	temp := make([]byte, 2)
	binary.BigEndian.PutUint16(temp[0:2], uint16(port.Pvid))
	ieeeOui8021 := &layers.LLDPOrgSpecificTLV{
		OUI:     layers.IEEEOUI8021,
		SubType: layers.LLDP8021SubtypePortVLANID,
	}
	ieeeOui8021.Info = append(ieeeOui8021.Info, temp...)
	tlv.Value = EncodePvidTLV(ieeeOui8021)
	if len(tlv.Value) > 0 {
		tlv.Length = uint16(len(tlv.Value))
		payload = append(payload, EncodeTLV(tlv)...)
	}

	// After all TLV's are added we need to go ahead and Add LLDPTLVEnd
	tlv = &layers.LinkLayerDiscoveryValue{}
	tlv.Type = layers.LLDPTLVEnd
	tlv.Length = 0
	payload = append(payload, EncodeTLV(tlv)...)
	return payload
}

/*  Encode Mandatory tlv, chassis id and port id
 */
func EncodeMandatoryTLV(Subtype byte, ID []byte) []byte {
	// 1 byte: subtype
	// N bytes: ID
	b := make([]byte, 1+len(ID))
	b[0] = byte(Subtype)
	copy(b[1:], ID)

	return b
}

// Marshall tlv information into binary form
// 1) Check type value
// 2) Check Length
func EncodeTLV(tlv *layers.LinkLayerDiscoveryValue) []byte {

	// copy value into b
	// type : 7 bits
	// leng : 9 bits
	// value: N bytes
	typeLen := uint16(tlv.Type)<<9 | tlv.Length
	temp := make([]byte, 2+len(tlv.Value))
	binary.BigEndian.PutUint16(temp[0:2], typeLen)
	copy(temp[2:], tlv.Value)
	return temp
}

/*  TLV Type = 8, 7 bits             ----
					|--> 2 bytes
 *  TLV Length = 9 bits....	     ----
 *  Value: N bytes
 *     Subtype is 1 byte
 *     Address is []byte
 *     IntefaceSubtype is 1 byte
 *     IntefaceNumber uint32 <<< this is system interface number which is IfIndex in our case
 *     OID string
*/
func EncodeMgmtTLV(tlv *layers.LLDPMgmtAddress) []byte {
	var b []byte
	b = append(b, byte(len(tlv.Address)+1 /*IEEE Guys why you want redundant information?*/))
	b = append(b, byte(tlv.Subtype))
	b = append(b, tlv.Address...)
	b = append(b, byte(tlv.InterfaceSubtype))
	temp := make([]byte, 4 /*uint32*/ +1 /*Length of OID String*/)
	binary.BigEndian.PutUint32(temp[0:4], tlv.InterfaceNumber)
	temp[4] = 0
	b = append(b, temp...)
	debug.Logger.Debug("byte returned", b)
	return b
}

func EncodePvidTLV(tlv *layers.LLDPOrgSpecificTLV) []byte {
	b := make([]byte, 4 /*uint32 oui*/)
	ouiInfo := uint32(tlv.OUI)<<8 | uint32(tlv.SubType)
	binary.BigEndian.PutUint32(b[0:4], ouiInfo)
	b = append(b, tlv.Info...) // [] byte
	debug.Logger.Debug("byte returned for pvid", b)
	return b
}

func (t *TX) UseCache() bool {
	return t.useCacheFrame
}

func (t *TX) SetCache(use bool) {
	t.useCacheFrame = use
}

/*  We have deleted the pcap handler and hence we will invalid the cache buffer
 */
func (t *TX) DeleteCacheFrame() {
	t.useCacheFrame = false
	t.cacheFrame = nil
}

/*  Stop Send Tx timer... as we have already delete the pcap handle
 */
func (t *TX) StopTxTimer() {
	if t.TxTimer != nil {
		debug.Logger.Debug("tx timer stopped")
		t.TxTimer.Stop()
		t.TxTimer = nil
	}
}
