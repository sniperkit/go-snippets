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

// tx
package drcp

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"l2/lacp/protocol/utils"
	"net"
)

// bridge will simulate communication between two channels
type SimulationNeighborBridge struct {
	Port1      uint32
	Port2      uint32
	RxIppPort1 chan gopacket.Packet
	RxIppPort2 chan gopacket.Packet
}

func (bridge *SimulationNeighborBridge) TxViaGoChannel(key IppDbKey, dmac net.HardwareAddr, pdu interface{}) {

	var p *DRCPIpp
	if DRFindPortByKey(key, &p) {

		// Set up all the layers' fields we can.
		eth := layers.Ethernet{
			SrcMAC:       net.HardwareAddr{0x00, uint8(p.Id & 0xff), 0x00, 0x01, 0x01, 0x01},
			DstMAC:       dmac,
			EthernetType: layers.EthernetTypeDRCP,
		}
		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		}

		switch pdu.(type) {
		case *layers.DRCP:
			drcp := pdu.(*layers.DRCP)
			gopacket.SerializeLayers(buf, opts, &eth, drcp)
		}

		pkt := gopacket.NewPacket(buf.Bytes(), layers.LinkTypeEthernet, gopacket.Default)
		if p.Id != bridge.Port1 && bridge.RxIppPort1 != nil {
			//fmt.Println("TX channel: Tx From port", p.Id, "bridge Port Rx", bridge.Port1)
			//fmt.Printf("TX: %+v", pkt)
			bridge.RxIppPort1 <- pkt
		} else if bridge.RxIppPort2 != nil {
			//fmt.Println("TX channel: Tx From port", p.Id, "bridge Port Rx", bridge.Port2)
			//fmt.Println("TX: %+v", pkt)
			bridge.RxIppPort2 <- pkt
		}
	} else {
		utils.GlobalLogger.Err(fmt.Sprintf("Unable to find port %d (%s) in tx", p.Id, p.Name))
	}
}

func TxViaLinuxIf(key IppDbKey, dmac net.HardwareAddr, pdu interface{}) {
	var p *DRCPIpp
	if DRFindPortByKey(key, &p) {

		txIface, err := net.InterfaceByName(p.Name)
		if err == nil {
			// conver the packet to a go packet
			// Set up all the layers' fields we can.
			eth := layers.Ethernet{
				SrcMAC:       txIface.HardwareAddr,
				DstMAC:       dmac,
				EthernetType: layers.EthernetTypeDRCP,
			}

			// Set up buffer and options for serialization.
			buf := gopacket.NewSerializeBuffer()
			opts := gopacket.SerializeOptions{
				FixLengths:       true,
				ComputeChecksums: true,
			}

			switch pdu.(type) {
			case *layers.DRCP:
				drcp := pdu.(*layers.DRCP)
				gopacket.SerializeLayers(buf, opts, &eth, drcp)
			}

			// Send one packet for every address.
			if err := p.handle.WritePacketData(buf.Bytes()); err != nil {
				utils.GlobalLogger.Err(fmt.Sprintf("%s\n", err))
			}
		} else {
			utils.GlobalLogger.Err(fmt.Sprintln("ERROR could not find IPP interface", p.Name, err))
		}
	} else {
		utils.GlobalLogger.Err(fmt.Sprintf("Unable to find port %d (%s) in tx", p.Id, p.Name))
	}
}
