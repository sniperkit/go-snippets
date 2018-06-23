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
	"bytes"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"net"
	"time"
)

func (session *BfdSession) StartPerLinkSessionServer(bfdServer *BFDServer) error {
	var err error
	var myMacAddr net.HardwareAddr
	ifName := session.state.Interface
	bfdServer.logger.Info("Starting perlink session ", session.state.SessionId, " on ", ifName)
	sessionId := session.state.SessionId
	myMacAddr, err = bfdServer.getMacAddrFromIntfName(ifName)
	if err != nil {
		bfdServer.logger.Info("Unable to get the MAC addr of ", ifName, err)
		return err
	}
	bfdServer.logger.Info("MAC is  ", myMacAddr, " on ", ifName)
	bfdPcapTimeout := time.Duration(session.state.RequiredMinRxInterval / 1000000)
	session.recvPcapHandle, err = pcap.OpenLive(ifName, bfdSnapshotLen, bfdPromiscuous, bfdPcapTimeout)
	if session.recvPcapHandle == nil {
		bfdServer.logger.Info("Failed to open recvPcapHandle for ", ifName, err)
		return err
	} else {
		err = session.recvPcapHandle.SetBPFFilter(bfdPcapFilterLag)
		if err != nil {
			bfdServer.logger.Info("Unable to set filter on", ifName, err)
			return err
		}
	}
	bfdPacketSrc := gopacket.NewPacketSource(session.recvPcapHandle, layers.LayerTypeEthernet)
	defer session.recvPcapHandle.Close()
	for receivedPacket := range bfdPacketSrc.Packets() {
		if bfdServer.bfdGlobal.Sessions[sessionId] == nil {
			return nil
		}
		bfdServer.logger.Info("Receive packet ", receivedPacket)

		ethLayer := receivedPacket.Layer(layers.LayerTypeEthernet)
		ethPacket, _ := ethLayer.(*layers.Ethernet)
		bfdServer.logger.Info("Ethernet ", ethPacket.SrcMAC, ethPacket.DstMAC)
		nwLayer := receivedPacket.Layer(layers.LayerTypeIPv4)
		ipPacket, _ := nwLayer.(*layers.IPv4)
		bfdServer.logger.Info("Network ", ipPacket.SrcIP, ipPacket.DstIP)
		transLayer := receivedPacket.Layer(layers.LayerTypeUDP)
		udpPacket, _ := transLayer.(*layers.UDP)
		bfdServer.logger.Info("Transport ", udpPacket.SrcPort, udpPacket.DstPort)
		appLayer := receivedPacket.ApplicationLayer()
		bfdServer.logger.Info("Application ", appLayer)

		if bytes.Equal(ethPacket.SrcMAC, myMacAddr) {
			bfdServer.logger.Info("My packet looped back")
			continue
		}

		buf := transLayer.LayerPayload()
		if len(buf) >= DEFAULT_CONTROL_PACKET_LEN {
			bfdPacket, err := DecodeBfdControlPacket(buf)
			if err == nil {
				sessionId := int32(bfdPacket.YourDiscriminator)
				if sessionId == 0 {
					bfdServer.logger.Info("Ignore bfd packet for session ", sessionId)
				} else {
					bfdSession := bfdServer.bfdGlobal.Sessions[sessionId]
					match := bytes.Equal(bfdSession.state.RemoteMacAddr, ethPacket.SrcMAC)
					if !match {
						bfdSession.state.RemoteMacAddr = ethPacket.DstMAC
					}
					bfdSession.state.NumRxPackets++
					bfdSession.ProcessBfdPacket(bfdPacket)
				}
			} else {
				bfdServer.logger.Info("Failed to decode packet - ", err)
			}
		}
	}
	return nil
}

func (session *BfdSession) StartPerLinkSessionClient(bfdServer *BFDServer) error {
	var err error
	var myMacAddr net.HardwareAddr
	ifName := session.state.Interface
	bfdServer.logger.Info("Starting perlink session ", session.state.SessionId, " on ", ifName)
	myMacAddr, err = bfdServer.getMacAddrFromIntfName(ifName)
	if err != nil {
		bfdServer.logger.Info("Unable to get the MAC addr of ", ifName, err)
		bfdServer.FailedSessionClientCh <- session.state.SessionId
		return err
	}
	bfdServer.logger.Info("MAC is  ", myMacAddr, " on ", ifName)
	bfdPcapTimeout := time.Duration(session.state.DesiredMinTxInterval / 1000000)
	session.sendPcapHandle, err = pcap.OpenLive(ifName, bfdSnapshotLen, bfdPromiscuous, bfdPcapTimeout)
	if session.sendPcapHandle == nil {
		bfdServer.logger.Info("Failed to open sendPcapHandle for ", ifName, err)
		bfdServer.FailedSessionClientCh <- session.state.SessionId
		return err
	}
	session.TxTimeoutCh = make(chan int32)
	session.SessionTimeoutCh = make(chan int32)
	sessionTimeoutMS := time.Duration(session.state.RequiredMinRxInterval * session.state.DetectionMultiplier / 1000)
	txTimerMS := time.Duration(session.state.DesiredMinTxInterval / 1000)
	session.sessionTimer = time.AfterFunc(time.Millisecond*sessionTimeoutMS, func() { session.SessionTimeoutCh <- session.state.SessionId })
	session.txTimer = time.AfterFunc(time.Millisecond*txTimerMS, func() { session.TxTimeoutCh <- session.state.SessionId })
	defer session.sendPcapHandle.Close()
	for {
		select {
		case sessionId := <-session.TxTimeoutCh:
			var destMac net.HardwareAddr
			bfdSession := bfdServer.bfdGlobal.Sessions[sessionId]
			if bfdSession.useDedicatedMac {
				destMac, _ = net.ParseMAC(bfdDedicatedMac)
			} else {
				destMac = bfdSession.state.RemoteMacAddr
			}
			ethLayer := &layers.Ethernet{
				SrcMAC:       bfdSession.state.LocalMacAddr,
				DstMAC:       destMac,
				EthernetType: layers.EthernetTypeIPv4,
			}
			ipLayer := &layers.IPv4{
				DstIP:    net.ParseIP(bfdSession.state.IpAddr),
				Protocol: layers.IPProtocolUDP,
			}
			udpLayer := &layers.UDP{
				SrcPort: layers.UDPPort(SRC_PORT_LAG),
				DstPort: layers.UDPPort(DEST_PORT_LAG),
			}
			options := gopacket.SerializeOptions{
				FixLengths:       true,
				ComputeChecksums: true,
			}
			bfdSession.UpdateBfdSessionControlPacket()
			bfdPacket, err := bfdSession.bfdPacket.CreateBfdControlPacket()
			if err != nil {
				bfdServer.logger.Info("Failed to create bfd control packet for session ", bfdSession.state.SessionId)
			}
			buffer := gopacket.NewSerializeBuffer()
			gopacket.SerializeLayers(buffer, options, ethLayer, ipLayer, udpLayer, gopacket.Payload(bfdPacket))
			outgoingPacket := buffer.Bytes()
			err = bfdSession.sendPcapHandle.WritePacketData(outgoingPacket)
			if err != nil {
				bfdServer.logger.Info("Failed to create complete packet for session ", bfdSession.state.SessionId)
			} else {
				if bfdSession.state.SessionState == STATE_UP {
					bfdSession.useDedicatedMac = false
				}
				bfdSession.state.NumTxPackets++
				txTimerMS = time.Duration(bfdSession.state.DesiredMinTxInterval / 1000)
				bfdSession.txTimer.Reset(time.Millisecond * txTimerMS)
			}
		case sessionId := <-session.SessionTimeoutCh:
			bfdSession := bfdServer.bfdGlobal.Sessions[sessionId]
			bfdSession.state.LocalDiagType = DIAG_TIME_EXPIRED
			bfdSession.EventHandler(TIMEOUT)
			sessionTimeoutMS = time.Duration(bfdSession.state.RequiredMinRxInterval * bfdSession.state.DetectionMultiplier / 1000)
			bfdSession.sessionTimer.Reset(time.Millisecond * sessionTimeoutMS)
		case <-session.SessionStopClientCh:
			return nil
		}
	}
	return nil
}
