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
	_ "fmt"
	"github.com/google/gopacket"
	"l2/lldp/utils"
	"time"
)

/* Go routine to recieve lldp frames. This go routine is created for all the
 * ports which are in up state.
 */
func (intf *LLDPGlobalInfo) ReceiveFrames(lldpRxPktCh chan InPktChannel) {
	pktSrc := gopacket.NewPacketSource(intf.PcapHandle, intf.PcapHandle.LinkType())
	in := pktSrc.Packets()
	// process packets
	for {
		select {
		case pkt, ok := <-in:
			if ok {
				lldpRxPktCh <- InPktChannel{pkt, intf.Port.IfIndex}
			} else {
				debug.Logger.Info("quit for ifIndex", intf.Port.Name, "rx exiting go routine")
				return
			}
		case <-intf.RxKill:
			debug.Logger.Info("quit for ifIndex", intf.Port.Name, "rx exiting go routine")
			intf.RxKill <- true
			return
		}
	}
}

/*  lldp server go routine to handle tx timer... once the timer fires we will
 *  send the ifindex on the channel to handle send info
 *
 *  For fast learning we will send out 5 frames in 5 seconds and then every 30 seconds an update frame will
 *  be send out
 */
func (intf *LLDPGlobalInfo) StartTxTimer(lldpTxPktCh chan SendPktChannel) {
	if intf.TxInfo.TxTimer != nil {
		if intf.counter.Send > LLDP_FAST_LEARN_MAX_FRAMES_SEND {
			intf.TxInfo.TxTimer.Reset(time.Duration(intf.TxInfo.MessageTxInterval) * time.Second)
		} else {
			intf.TxInfo.TxTimer.Reset(time.Duration(LLDP_FAST_LEARN_TIMER) * time.Second)
		}
	} else {
		var TxTimerHandler_func func()
		TxTimerHandler_func = func() {
			lldpTxPktCh <- SendPktChannel{intf.Port.IfIndex}
		}

		// Create an After Func and go routine for it, so that on timer stop TX is stopped automatically
		intf.TxInfo.TxTimer = time.AfterFunc(time.Duration(LLDP_FAST_LEARN_TIMER)*time.Second,
			TxTimerHandler_func)
	}
}

/*  Write packet is helper function to send packet on wire.
 *  It will inform caller that packet was send successfully and you can go ahead
 *  and cache the pkt or else do not cache the packet as it is corrupted or there
 *  was some error
 */
func (intf *LLDPGlobalInfo) WritePacket(pkt []byte) bool {
	var err error
	if len(pkt) == 0 {
		return false
	}
	if intf.PcapHandle != nil {
		err = intf.PcapHandle.WritePacketData(pkt)
	} else {
		err = errors.New("Pcap Handle is invalid for " + intf.Port.Name)
	}
	if err != nil {
		debug.Logger.Err("Sending packet failed Error:", err, "for Port:", intf.Port.Name)
		return false
	}
	return true
}
