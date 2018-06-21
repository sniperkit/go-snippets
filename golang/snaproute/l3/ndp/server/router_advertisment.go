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
	"net"
)

func (intf *Interface) createRANbrKey(ndInfo *packet.NDInfo) (tgtMac, nbrKey string) {
	for _, option := range ndInfo.Options {
		if option.Type == packet.NDOptionTypeSourceLinkLayerAddress {
			tgtMac = net.HardwareAddr(option.Value).String()
			nbrKey = tgtMac + "_" + ndInfo.SrcIp + "_" + ndInfo.LearnedIntfRef
			debug.Logger.Debug("RA nbrKey created:", nbrKey)
			break
		}
	}
	return tgtMac, nbrKey
}

/*
 * When we get router advertisement packet we need to update the mac address of peer and move the state to
 * REACHABLE
 *
 * Based on ifIndex we will get a prefixLink which contains all the prefixes for that link
 *
 * fill the NDInfo and then return it back to caller
 */
func (intf *Interface) processRA(ndInfo *packet.NDInfo) (nbrInfo *config.NeighborConfig, oper NDP_OPERATION) {
	// RA generated for my own ip address
	if ndInfo.SrcIp == intf.linkScope || ndInfo.SrcIp == intf.globalScope {
		return nil, IGNORE
	}
	var nbrKey string
	// we need to use Option field to determine who is the original sender of this advertisement... if no
	// options then the port itself is the sender
	_, nbrKey = intf.createRANbrKey(ndInfo)
	if nbrKey == "" {
		nbrKey = intf.createNbrKey(ndInfo)
	}
	if !intf.validNbrKey(nbrKey) {
		return nil, IGNORE
	}
	nbr, exists := intf.Neighbor[nbrKey]
	if exists {
		if ndInfo.RouterLifetime == 0 {
			// delete this neighbor
			nbrInfo = nbr.populateNbrInfo(intf.IfIndex, intf.IntfRef)
			nbr.DeInit()
			delete(intf.Neighbor, nbr.IpAddr)
			return nbrInfo, DELETE
		} else {
			// update existing neighbor timers
			nbr.State = REACHABLE
			// Router Lifetime/Invalidation Timer reset
			nbr.InValidTimer(ndInfo.RouterLifetime)
			// Recahable timer reset
			// Stop any probes
			nbr.RchTimer()
			oper = UPDATE
		}
	} else {
		// create new neighbor
		nbr.InitCache(intf.reachableTime, intf.retransTime, nbrKey, intf.PktDataCh, intf.IfIndex)
		nbr.InValidTimer(ndInfo.RouterLifetime)
		nbr.RchTimer()
		nbr.State = REACHABLE
		oper = CREATE
	}
	nbrInfo = nbr.populateNbrInfo(intf.IfIndex, intf.IntfRef)
	nbr.updatePktRxStateInfo()
	intf.Neighbor[nbrKey] = nbr
	return nbrInfo, oper
}

/*
 *  Router Advertisement Packet is send out for both link scope ip and global scope ip on timer expiry & port
 *  up notification
 */
func (intf *Interface) SendRA(srcMac string) {
	pkt := &packet.Packet{
		SrcMac: srcMac,
		DstMac: ALL_NODES_MULTICAST_LINK_LAYER_ADDRESS,
		DstIp:  ALL_NODES_MULTICAST_IPV6_ADDRESS,
		PType:  layers.ICMPv6TypeRouterAdvertisement,
	}
	if intf.linkScope != "" {
		pkt.SrcIp = intf.linkScope
		pktToSend := pkt.Encode()
		intf.writePkt(pktToSend)
		intf.counter.Send++
	}
	if intf.globalScope != "" {
		pkt.SrcIp = intf.globalScope
		pktToSend := pkt.Encode()
		intf.writePkt(pktToSend)
		intf.counter.Send++
	}

	intf.RAResTransmitTimer()
}
