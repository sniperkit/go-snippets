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
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	_ "github.com/google/gopacket/pcap"
	"l2/lldp/config"
	"l2/lldp/utils"
	"net"
	"time"
)

func RxInit() *RX {
	var err error
	rxInfo := &RX{}
	rxInfo.DstMAC, err = net.ParseMAC(LLDP_PROTO_DST_MAC)
	if err != nil {
		debug.Logger.Err(fmt.Sprintln("parsing lldp protocol Mac failed",
			err))
	}

	return rxInfo
}

/*  Upon receiving incoming packet check whether all the madatory layer info is
 *  correct or not.. If not then treat the packet as corrupted and move on
 */
func (p *RX) VerifyFrame(lldpInfo *layers.LinkLayerDiscovery) error {

	if lldpInfo.ChassisID.Subtype > layers.LLDPChassisIDSubTypeLocal {
		return errors.New("Invalid chassis id subtype")
	}

	if lldpInfo.PortID.Subtype > layers.LLDPPortIDSubtypeLocal {
		return errors.New("Invalid port id subtype")
	}

	if lldpInfo.TTL > uint16(LLDP_MAX_TTL) {
		return errors.New("Invalid TTL value")
	}
	return nil
}

func (p *RX) Process(rxInfo *RX, pkt gopacket.Packet) (int, error) {
	event := config.NoOp
	ethernetLayer := pkt.Layer(layers.LayerTypeEthernet)
	if ethernetLayer == nil {
		return event, errors.New("Invalid eth layer")
	}
	eth := ethernetLayer.(*layers.Ethernet)
	// copy src mac and dst mac
	rxInfo.SrcMAC = eth.SrcMAC
	if rxInfo.DstMAC.String() != eth.DstMAC.String() {
		return event, errors.New("Invalid DST MAC in rx frame")
	}
	// Get lldp manadatory layer and optional info
	lldpLayer := pkt.Layer(layers.LayerTypeLinkLayerDiscovery)
	lldpLayerInfo := pkt.Layer(layers.LayerTypeLinkLayerDiscoveryInfo)
	// Verify that the information is not nil
	if lldpLayer == nil || lldpLayerInfo == nil {
		return event, errors.New("Invalid Frame")
	}

	// Verify that the mandatory layer info is indeed correct
	err := p.VerifyFrame(lldpLayer.(*layers.LinkLayerDiscovery))
	if err != nil {
		return event, err
	}
	// Update last packet byte for cacheing...
	if len(rxInfo.LastPkt) == 0 || rxInfo.LastPkt == nil {
		//this is new cache set event state to be learned
		rxInfo.LastPkt = pkt.Data()
		event = config.Learned
	} else {
		// if incoming packet has difference then it means that we need to publish event
		if bytes.Compare(rxInfo.LastPkt, pkt.Data()) != 0 {
			event = config.Updated
		}
	}

	if rxInfo.RxFrame == nil {
		rxInfo.RxFrame = new(layers.LinkLayerDiscovery)
	}
	// Store lldp frame information received from direct connection
	*rxInfo.RxFrame = *lldpLayer.(*layers.LinkLayerDiscovery)

	if rxInfo.RxLinkInfo == nil {
		rxInfo.RxLinkInfo = new(layers.LinkLayerDiscoveryInfo)
	}
	// Store lldp link layer optional tlv information
	*rxInfo.RxLinkInfo = *lldpLayerInfo.(*layers.LinkLayerDiscoveryInfo)
	// LLDPInfo8021 represents the information carried in 802.1 Org-specific TLVs
	rxInfo.OrgTLV802_1, err = rxInfo.RxLinkInfo.Decode8021()
	if err != nil {
		debug.Logger.Warning("failed to decode 802.11 TLV")
	}

	return event, nil
}

/*
 *  Handle TTL timer. Once the timer expires, we will delete the remote entry
 *  if timer is running then reset the value
 */
func (rxInfo *RX) CheckPeerEntry(port string, eCh chan config.EventInfo, ifIndex int32) {
	if rxInfo.ClearCacheTimer != nil {
		// timer is running reset the time so that it doesn't expire
		rxInfo.ClearCacheTimer.Reset(time.Duration(rxInfo.RxFrame.TTL) * time.Second)
	} else {
		var clearPeerInfo_func func()
		// On timer expiration we will delete peer info and set it to nil
		clearPeerInfo_func = func() {
			debug.Logger.Info("Recipient info delete timer expired for " + "peer connected to port " +
				port + " and hence deleting peer information from runtime")
			rxInfo.RxFrame = nil
			rxInfo.RxLinkInfo = nil
			rxInfo.LastPkt = nil
			eCh <- config.EventInfo{
				EventType: config.Removed,
				IfIndex:   ifIndex,
			}
		}
		// First time start function
		rxInfo.ClearCacheTimer = time.AfterFunc(time.Duration(rxInfo.RxFrame.TTL)*time.Second,
			clearPeerInfo_func)
	}
}
