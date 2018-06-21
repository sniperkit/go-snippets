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
	"github.com/google/gopacket/layers"
	"l3/ndp/config"
	"l3/ndp/debug"
	"l3/ndp/packet"
)

/*
 *  Generic API to send Neighbor Solicitation Packet based on inputs..
 */
func (intf *Interface) sendUnicastNS(srcMac, nbrMac, nbrIp string, isFastProbe bool) NDP_OPERATION {
	nbrKey := nbrIp + "_" + nbrMac
	nbr, exists := intf.Neighbor[nbrKey]
	if !exists {
		return IGNORE
	}

	if nbr.ProbesSent > MAX_UNICAST_SOLICIT {
		intf.FlushNeighborPerIp(nbrKey, nbrIp)
		return DELETE
	}

	pkt := &packet.Packet{
		SrcMac: srcMac,
		DstMac: nbrMac,
		DstIp:  nbrIp,
		PType:  layers.ICMPv6TypeNeighborSolicitation,
	}
	if isLinkLocal(nbrIp) {
		pkt.SrcIp = intf.linkScope
	} else {
		pkt.SrcIp = intf.globalScope
	}

	pktToSend := pkt.Encode()
	err := intf.writePkt(pktToSend)
	if err != nil {
		return IGNORE
	}

	// when sending unicast packet re-start retransmit/delay probe timer.. rest all will be taken care of when
	// NA packet is received..
	if isFastProbe == false {
		debug.Logger.Debug("Reachable Timer is probing Neighbor:", nbrIp, nbrMac)
		if nbr.State == REACHABLE {
			// This means that Reachable Timer has expierd and hence we are sending Unicast Message..
			// Lets set the time for delay first probe
			nbr.DelayProbe()
			nbr.State = DELAY
			nbr.ProbesSent = 0
		} else if nbr.State == DELAY || nbr.State == PROBE {
			// Probes Sent can still be zero but the state has changed to Delay..
			// Start Timer for Probe and move the state from delay to Probe
			nbr.Timer()
			nbr.State = PROBE
			nbr.ProbesSent += 1
		}
	} else {
		// reset the fast probe to reachable timer / 2 * FastProbleMultipliers
		//if nbr.StopFastProbe == false {
		debug.Logger.Debug("Fast Probe Timer is probing Neighbor:", nbrIp, nbrMac)
		//	nbr.FastProbe()
		//}
	}
	intf.counter.Send++
	nbr.counter.Send++
	intf.Neighbor[nbrKey] = nbr
	return IGNORE
}

func (intf *Interface) SendNS(myMac, nbrMac, nbrIp string, isFastProbe bool) NDP_OPERATION {
	if nbrIp == "" {
		//send multicast solicitation when we take over linux called when port comes up
	} else {
		return intf.sendUnicastNS(myMac, nbrMac, nbrIp, isFastProbe)
	}

	return IGNORE
}

func (intf *Interface) deleteIncompleteEntry(ipAddr string) {
	delete(intf.incompleteNbr, ipAddr)
	debug.Logger.Debug("incompleteNbr list is:", intf.incompleteNbr)
}

func (intf *Interface) insertIncompleteEntry(ndInfo *packet.NDInfo) {
	var empty struct{}
	intf.incompleteNbr[ndInfo.TargetAddress.String()] = empty
	debug.Logger.Debug("incompleteNbr list is:", intf.incompleteNbr)
}

/*
 *  helper function to handle incoming Neighbor solicitation messages...
 *  Case 1) SrcIP == "::"
 *		This is a message which is locally generated. In this case Target Address will be our own
 *		IP Address, which is not a Neighbor and hence we should not create a entry in NbrCache
 *  Case 2) SrcIP != "::"
 *		This is a message coming from our Neighbor. Ok now what do we need to do?
 *		If no cache entry:
 *		    Then create a cache entry and mark that entry as incomplete
 *		    nd info: {4000::105:1:1:1 0(135) -1 cc:37:ab:a9:23:f5 4000::105:1:1:2 ff02::1:ff01:1 33:33:ff:01:00:01 fpPort48 0 0 0 0 0 0 [0xc8211290c0]}
 *		If cache entry exists:
 *		    Then update the state to STALE
 */
func (intf *Interface) processNS(ndInfo *packet.NDInfo) (nbrInfo *config.NeighborConfig, oper NDP_OPERATION) {
	if ndInfo.SrcIp == "" || ndInfo.SrcIp == "::" {
		// NS was generated locally or it is multicast-solicitation message
		return nil, IGNORE
	}
	if (ndInfo.SrcIp == intf.linkScope || ndInfo.SrcIp == intf.globalScope) && len(ndInfo.TargetAddress) != 0 {
		intf.insertIncompleteEntry(ndInfo)
		return nil, IGNORE
	}
	nbrKey := intf.createNbrKey(ndInfo)
	if !intf.validNbrKey(nbrKey) {
		return nil, IGNORE
	}
	nbr, exists := intf.Neighbor[nbrKey]
	if exists {
		// update the neighbor ??? what to do in this case moving to stale
		//nbr.State = STALE
		oper = UPDATE
	} else {
		// create new neighbor
		nbr.InitCache(intf.reachableTime, intf.retransTime, nbrKey, intf.PktDataCh, intf.IfIndex)
		oper = CREATE
	}
	if len(ndInfo.Options) > 0 {
		for _, option := range ndInfo.Options {
			if option.Type == packet.NDOptionTypeSourceLinkLayerAddress {
				nbr.State = REACHABLE
			}
		}
	}
	nbrInfo = nbr.populateNbrInfo(intf.IfIndex, intf.IntfRef)
	nbr.updatePktRxStateInfo()
	intf.Neighbor[nbrKey] = nbr
	return nbrInfo, oper
}
