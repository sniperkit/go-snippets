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

// tx.go
package stp

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func ConvertBoolToUint8(v bool) (rv uint8) {
	if v {
		rv = 1
	}
	return rv
}

func ConvertRoleToPktRole(r PortRole) (rv uint8) {

	switch r {
	case PortRoleAlternatePort, PortRoleBackupPort:
		rv = layers.RoleAlternateBackupPort
	case PortRoleRootPort:
		rv = layers.RoleRootPort
	case PortRoleDesignatedPort:
		rv = layers.RoleDesignatedPort
	default:
		rv = layers.RoleMasterPort
	}
	return rv
}

func (p *StpPort) BuildRSTPEthernetLlcHeaders() (eth layers.Ethernet, llc layers.LLC) {
	pIntf, _ := PortConfigMap[p.IfIndex]

	eth = layers.Ethernet{
		SrcMAC: pIntf.HardwareAddr,
		DstMAC: layers.BpduDMAC,
		// length
		EthernetType: layers.EthernetTypeLLC,
		Length:       uint16(layers.STPProtocolLength + 3), // found from PCAP from packetlife.net
	}

	llc = layers.LLC{
		DSAP:    0x42,
		IG:      false,
		SSAP:    0x42,
		CR:      false,
		Control: 0x03,
	}
	return eth, llc
}

func (p *StpPort) TxPVST() {
	if p.handle != nil {
		pIntf, _ := PortConfigMap[p.IfIndex]

		eth := layers.Ethernet{
			SrcMAC: pIntf.HardwareAddr,
			DstMAC: layers.BpduPVSTDMAC,
			// length
			EthernetType: layers.EthernetTypeDot1Q,
		}

		vlan := layers.Dot1Q{
			Priority:       PVST_VLAN_PRIORITY,
			DropEligible:   false,
			VLANIdentifier: p.b.Vlan,
			Type:           layers.EthernetType(layers.PVSTProtocolLength + 3 + 5), // length
		}

		llc := layers.LLC{
			DSAP:    0xAA,
			IG:      false,
			SSAP:    0xAA,
			CR:      false,
			Control: 0x03,
		}

		snap := layers.SNAP{
			OrganizationalCode: []byte{0x00, 0x00, 0x0C},
			Type:               0x010b,
		}

		pvst := layers.PVST{
			ProtocolId:        layers.RSTPProtocolIdentifier,
			ProtocolVersionId: p.BridgeProtocolVersionGet(),
			BPDUType:          layers.BPDUTypeRSTP,
			Flags:             0,
			RootId:            p.PortPriority.RootBridgeId,
			RootPathCost:      uint32(p.b.BridgePriority.RootPathCost),
			BridgeId:          p.b.BridgePriority.DesignatedBridgeId,
			PortId:            uint16(p.PortId | p.Priority<<8),
			MsgAge:            uint16(p.b.RootTimes.MessageAge << 8),
			MaxAge:            uint16(p.b.RootTimes.MaxAge << 8),
			HelloTime:         uint16(p.b.RootTimes.HelloTime << 8),
			FwdDelay:          uint16(p.b.RootTimes.ForwardingDelay << 8),
			Version1Length:    0,
			OriginatingVlan: layers.STPOriginatingVlanTlv{
				Type:     0,
				Length:   2,
				OrigVlan: p.b.Vlan,
			},
		}

		var flags uint8
		StpSetBpduFlags(ConvertBoolToUint8(p.TcAck),
			ConvertBoolToUint8(p.Agree),
			ConvertBoolToUint8(p.Forwarding),
			ConvertBoolToUint8(p.Learning),
			ConvertRoleToPktRole(p.Role),
			ConvertBoolToUint8(p.Proposed),
			ConvertBoolToUint8(p.TcWhileTimer.count != 0),
			&flags)

		pvst.Flags = layers.StpFlags(flags)

		/* NOT VALID within PVST STP frames should have been detected outside of this logic
		if !p.SendRSTP {
			pvst.ProtocolId = layers.RSTPProtocolIdentifier
			pvst.ProtocolVersionId = layers.STPProtocolVersion
			pvst.BPDUType = layers.BPDUTypeSTP
			// only tc and tc ack are valid for stp
			StpSetBpduFlags(ConvertBoolToUint8(p.TcAck),
				0,
				0,
				0,
				0,
				0,
				ConvertBoolToUint8(p.TcWhileTimer.count != 0),
				&flags)

			pvst.Flags = layers.StpFlags(flags)
		}
		*/

		// Set up buffer and options for serialization.
		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		}
		// Send one packet for every address.
		gopacket.SerializeLayers(buf, opts, &eth, &vlan, &llc, &snap, &pvst)
		if err := p.handle.WritePacketData(buf.Bytes()); err != nil {
			StpLogger("ERROR", fmt.Sprintf("Error writing packet to interface %s\n", err))
			return
		}

		p.SetTxPortCounters(BPDURxTypePVST)
		if p.TcWhileTimer.count != 0 {
			StpMachineLogger("DEBUG", "TX", p.IfIndex, p.BrgIfIndex, fmt.Sprintf("Sent TC packet on interface %s\n", pIntf.Name))
			p.SetTxPortCounters(BPDURxTypeTopo)
		}
		if p.TcAck {
			StpMachineLogger("DEBUG", "TX", p.IfIndex, p.BrgIfIndex, fmt.Sprintf("Sent TC Ack packet on interface %s\n", pIntf.Name))
			p.SetTxPortCounters(BPDURxTypeTopoAck)
		}

		//StpLogger("DEBUG", fmt.Sprintf("Sent PVST packet on interface %s %#v\n", pIntf.Name, pvst))
	}
}

func (p *StpPort) TxRSTP() {

	if p.handle != nil {
		if p.b.Vlan != DEFAULT_STP_BRIDGE_VLAN {
			p.TxPVST()
			return
		}

		eth, llc := p.BuildRSTPEthernetLlcHeaders()

		rstp := layers.RSTP{
			ProtocolId:        layers.RSTPProtocolIdentifier,
			ProtocolVersionId: p.BridgeProtocolVersionGet(),
			BPDUType:          layers.BPDUTypeRSTP,
			Flags:             0,
			RootId:            p.PortPriority.RootBridgeId,
			RootPathCost:      uint32(p.b.BridgePriority.RootPathCost),
			BridgeId:          p.b.BridgePriority.DesignatedBridgeId,
			PortId:            uint16(p.PortId | p.Priority<<8),
			MsgAge:            uint16(p.b.RootTimes.MessageAge << 8),
			MaxAge:            uint16(p.b.RootTimes.MaxAge << 8),
			HelloTime:         uint16(p.b.RootTimes.HelloTime << 8),
			FwdDelay:          uint16(p.b.RootTimes.ForwardingDelay << 8),
			Version1Length:    0,
		}

		var flags uint8
		StpSetBpduFlags(ConvertBoolToUint8(p.TcAck),
			ConvertBoolToUint8(p.Agree),
			ConvertBoolToUint8(p.Forwarding),
			ConvertBoolToUint8(p.Learning),
			ConvertRoleToPktRole(p.Role),
			ConvertBoolToUint8(p.Proposed),
			ConvertBoolToUint8(p.TcWhileTimer.count != 0),
			&flags)

		rstp.Flags = layers.StpFlags(flags)

		// Set up buffer and options for serialization.
		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		}
		// Send one packet for every address.
		gopacket.SerializeLayers(buf, opts, &eth, &llc, &rstp)
		if err := p.handle.WritePacketData(buf.Bytes()); err != nil {
			StpLogger("ERROR", fmt.Sprintf("Error writing packet to interface %s\n", err))
			return
		}
		pIntf, _ := PortConfigMap[p.IfIndex]
		p.SetTxPortCounters(BPDURxTypeRSTP)
		if p.TcWhileTimer.count != 0 {
			StpMachineLogger("DEBUG", "TX", p.IfIndex, p.BrgIfIndex, fmt.Sprintf("Sent TC packet on interface %s\n", pIntf.Name))
			p.SetTxPortCounters(BPDURxTypeTopo)
		}
		if p.TcAck {
			StpMachineLogger("DEBUG", "TX", p.IfIndex, p.BrgIfIndex, fmt.Sprintf("Sent TC Ack packet on interface %s\n", pIntf.Name))
			p.SetTxPortCounters(BPDURxTypeTopoAck)
		}
		//StpLogger("DEBUG", fmt.Sprintf("Sent RSTP packet on interface %s %#v\n", pIntf.Name, rstp))
	}
}

func (p *StpPort) TxTCN() {
	if p.handle != nil {
		eth, llc := p.BuildRSTPEthernetLlcHeaders()

		if !p.SendRSTP {

			topo := layers.BPDUTopology{
				ProtocolId:        layers.RSTPProtocolIdentifier,
				ProtocolVersionId: layers.STPProtocolVersion,
				BPDUType:          layers.BPDUTypeTopoChange,
			}

			// Set up buffer and options for serialization.
			buf := gopacket.NewSerializeBuffer()
			opts := gopacket.SerializeOptions{
				FixLengths:       true,
				ComputeChecksums: true,
			}
			// Send one packet for every address.
			gopacket.SerializeLayers(buf, opts, &eth, &llc, &topo)
			if err := p.handle.WritePacketData(buf.Bytes()); err != nil {
				StpLogger("ERROR", fmt.Sprintf("Error writing packet to interface %s\n", err))
				return
			}
		} else {
			p.TxRSTP()
		}

		p.SetTxPortCounters(BPDURxTypeSTP)
		p.SetTxPortCounters(BPDURxTypeTopo)
		pIntf, _ := PortConfigMap[p.IfIndex]
		StpMachineLogger("DEBUG", "TX", p.IfIndex, p.BrgIfIndex, fmt.Sprintf("Sent TCN packet on interface %s\n", pIntf.Name))
	}
}

func (p *StpPort) TxConfig() {
	if p.handle != nil {
		eth, llc := p.BuildRSTPEthernetLlcHeaders()

		if p.b.Vlan != DEFAULT_STP_BRIDGE_VLAN {
			p.TxPVST()
			return
		}

		stp := layers.STP{
			ProtocolId:        layers.RSTPProtocolIdentifier,
			ProtocolVersionId: layers.STPProtocolVersion,
			BPDUType:          layers.BPDUTypeSTP,
			Flags:             0,
			RootId:            p.PortPriority.RootBridgeId,
			RootPathCost:      uint32(p.b.BridgePriority.RootPathCost),
			BridgeId:          p.b.BridgePriority.DesignatedBridgeId,
			PortId:            uint16(p.PortId | p.Priority<<8),
			MsgAge:            uint16(p.b.RootTimes.MessageAge << 8),
			MaxAge:            uint16(p.b.RootTimes.MaxAge << 8),
			HelloTime:         uint16(p.b.RootTimes.HelloTime << 8),
			FwdDelay:          uint16(p.b.RootTimes.ForwardingDelay << 8),
		}
		var flags uint8
		// only tc and tc ack are valid for stp
		StpSetBpduFlags(ConvertBoolToUint8(p.TcAck),
			0,
			0,
			0,
			0,
			0,
			ConvertBoolToUint8(p.TcWhileTimer.count != 0),
			&flags)

		stp.Flags = layers.StpFlags(flags)

		// Set up buffer and options for serialization.
		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		}
		// Send one packet for every address.
		gopacket.SerializeLayers(buf, opts, &eth, &llc, &stp)
		if err := p.handle.WritePacketData(buf.Bytes()); err != nil {
			StpLogger("ERROR", fmt.Sprintf("Error writing packet to interface %s\n", err))
			return
		}

		p.SetTxPortCounters(BPDURxTypeSTP)
		if p.TcWhileTimer.count != 0 {
			p.SetTxPortCounters(BPDURxTypeTopo)
		}
		if p.TcAck {
			p.SetTxPortCounters(BPDURxTypeTopoAck)
		}
	}
	//pIntf, _ := PortConfigMap[p.IfIndex]
	//StpLogger("DEBUG", fmt.Sprintf("Sent Config packet on interface %s %#v\n", pIntf.Name, stp))
}
