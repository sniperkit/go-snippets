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
	"l3/ndp/config"
	"l3/ndp/debug"
	"l3/ndp/packet"
	"net"
)

const (
	ROUTER_SET    = 0x80
	SOLICITED_SET = 0x40
	OVERRIDE_SET  = 0x20
)

func (intf *Interface) createNANbrKey(ndInfo *packet.NDInfo) (tgtMac, nbrKey string) {
	for _, option := range ndInfo.Options {
		if option.Type == packet.NDOptionTypeTargetLinkLayerAddress {
			tgtMac = net.HardwareAddr(option.Value).String()
			nbrKey = tgtMac + "_" + ndInfo.SrcIp + "_" + ndInfo.LearnedIntfRef
			debug.Logger.Debug("NA nbrKey created:", nbrKey)
			break
		}
	}
	return tgtMac, nbrKey
}

/*
 * When we get advertisement packet we need to update the mac address of peer and move the state to
 * REACHABLE
 *
 * If srcIP is my own IP then linux is responding for earlier solicitation message and hence we need to update
 * our cache entry with reachable
 * If srcIP is peer ip then we need to use dst ip to get link information and then update cache entry to be
 * reachable and also update peer mac address into the cache
 * @TODO: handle un-solicited Neighbor Advertisemtn
 */
func (intf *Interface) processNA(ndInfo *packet.NDInfo) (nbrInfo *config.NeighborConfig, oper NDP_OPERATION) {
	if ndInfo.SrcIp == intf.linkScope || ndInfo.SrcIp == intf.globalScope {
		// NA was generated locally or it is multicast-solicitation message
		return nil, IGNORE
	}
	var nbrKey string
	var tgtMac string
	if (ndInfo.ReservedFlags & OVERRIDE_SET) != 0x00 {
		debug.Logger.Debug("override is set and hence using create na nbr key")
		tgtMac, nbrKey = intf.createNANbrKey(ndInfo)
	} else {
		nbrKey = intf.createNbrKey(ndInfo)
	}
	if !intf.validNbrKey(nbrKey) {
		return nil, IGNORE
	}
	nbr, exists := intf.Neighbor[nbrKey]
	if !exists {
		// not exists lets check incomplete Nbr list
		// e.g: NA nbrKey created: d4:ae:52:e7:bb:b0_4000::105:1:1:1_fpPort48]
		_, exists := intf.incompleteNbr[ndInfo.SrcIp]
		if !exists {
			return nil, IGNORE
		}
		// create new neighbor as the entry has moved from INCOMPLETE to reachable with the NA
		nbr.InitCache(intf.reachableTime, intf.retransTime, nbrKey, intf.PktDataCh, intf.IfIndex)
		intf.deleteIncompleteEntry(ndInfo.SrcIp)
		oper = CREATE
		nbr.State = REACHABLE
		goto early_exit
	}

	// override flag is clear
	if ndInfo.ReservedFlags&OVERRIDE_SET == 0x00 && nbr.State != INCOMPLETE {
		if nbr.State == REACHABLE {
			nbr.State = STALE
		} else {
			return nil, IGNORE
		}
	} else {
		debug.Logger.Debug("override is set check SOLICITED_SET")
		// if override flag is set or supplied link-layer address is the same as in cache (this is validated by nbrKey)
		// @TODO: or no Target Link Layer Address was supplied
		/*
			The link-layer address in the Target Link-Layer Address option
			MUST be inserted in the cache (if one is supplied and differs
			from the already recorded address).
		*/
		if ndInfo.ReservedFlags&SOLICITED_SET != 0x0 {
			debug.Logger.Debug("SOLICITED_SET is set mark nbr as reachable and move on")
			nbr.State = REACHABLE
			nbr.UpdateProbe()
			nbr.RchTimer()
		} else if ndInfo.ReservedFlags&SOLICITED_SET == 0x0 && tgtMac != ndInfo.SrcMac {
			nbr.State = STALE
		}
	}
	oper = UPDATE
early_exit:
	nbrInfo = nbr.populateNbrInfo(intf.IfIndex, intf.IntfRef)
	nbr.updatePktRxStateInfo()
	intf.Neighbor[nbrKey] = nbr
	return nbrInfo, oper
}
