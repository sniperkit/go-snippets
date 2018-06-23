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
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"l3/ndp/debug"
)

/*
 * internal api for creating pcap handler for l2 untagged/tagged physical port for RX
 */
func (l2Port *PhyPort) createPortPcap(pktCh chan *RxPktInfo, name string) (err error) {
	if l2Port.RX == nil {
		debug.Logger.Debug("creating l2 rx pcap for", name, l2Port.Info.IfIndex)
		l2Port.RX, err = pcap.OpenLive(name, NDP_PCAP_SNAPSHOTlEN, NDP_PCAP_PROMISCUOUS, NDP_PCAP_TIMEOUT)
		if err != nil {
			debug.Logger.Err("Creating Pcap Handler failed for l2 interface:", name, "Error:", err)
			return err
		}
		err = l2Port.RX.SetBPFFilter(NDP_PCAP_FILTER)
		if err != nil {
			debug.Logger.Err("Creating BPF Filter failed Error", err)
			l2Port.RX = nil
			return err
		}
		debug.Logger.Info("Created l2 Pcap handler for port:", name, "now start receiving NdpPkts")
		go l2Port.L2ReceiveNdpPkts(pktCh)
	}
	return nil
}

/*
 * internal api for creating pcap handler for l2 physical port for RX
 */
func (l2Port *PhyPort) deletePcap() {
	if l2Port.RX != nil {
		l2Port.RX.Close()
		l2Port.RX = nil
	}
}

/*
 * Receive Ndp Packet and push it on the pktCh
 */
func (l2Port *PhyPort) L2ReceiveNdpPkts(pktCh chan *RxPktInfo) error {
	if l2Port.RX == nil {
		debug.Logger.Err("pcap handler for port:", l2Port.Info.Name, "is not valid. ABORT!!!!")
		return errors.New(fmt.Sprintln("pcap handler for port:", l2Port.Info.Name, "is not valid. ABORT!!!!"))
	}
	src := gopacket.NewPacketSource(l2Port.RX, layers.LayerTypeEthernet)
	in := src.Packets()
	for {
		select {
		case pkt, ok := <-in:
			if ok {
				pktCh <- &RxPktInfo{pkt, l2Port.Info.IfIndex}
			} else {
				debug.Logger.Debug("L2 Pcap closed as in is invalid exiting go routine for port:", l2Port.Info.Name)
				return nil
			}
		}
	}
	return nil
}

func (l2Port *PhyPort) updateFilter(macAddr string) {
	if l2Port.RX != nil {
		err := l2Port.RX.SetBPFFilter(getNewFilter(macAddr))
		if err != nil {
			debug.Logger.Err("Updating BPF Filter failed for interface:", l2Port.Info.Name, "Error", err)
			return
		}
		debug.Logger.Debug("Updating BPF Filter for interface:", l2Port.Info.Name)

	}
}

func (l2Port *PhyPort) resetFilter() {
	if l2Port.RX != nil {
		err := l2Port.RX.SetBPFFilter(NDP_PCAP_FILTER)
		if err != nil {
			debug.Logger.Err("Reseting BPF Filter failed for interface:", l2Port.Info.Name, "Error", err)
		}
	}
}
