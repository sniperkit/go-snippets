//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
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
	"fmt"
	//"github.com/google/gopacket"
	//"github.com/google/gopacket/layers"
	//"l3/ospf/config"
	"errors"
	"l3/ospfv2/objects"
	"math"
	"net"
)

func (server *OSPFV2Server) selfGenLsaCheck(key LsaKey) bool {
	rtr_id := server.globalData.RouterId
	if key.AdvRouter == rtr_id {
		return true
	}
	return false
}
func (server *OSPFV2Server) lsaUpdDiscardCheck(nbrConf NbrConf, data []byte) bool {
	if nbrConf.State < NbrExchange {
		server.logger.Info(fmt.Sprintln("LSAUPD: Discard .. Nbrstate (expected less than exchange)", nbrConf.State))
		return true
	}

	return false
}
func (server *OSPFV2Server) lsAgeCheck(intf IntfConfKey, lsa_max_age bool, exist bool) bool {

	send_ack := true
	/*
	           if the LSA's LS age is equal to MaxAge, and there is
	       currently no instance of the LSA in the router's link state
	       database, and none of router's neighbors are in states Exchange
	       or Loading, then take the following actions: a) Acknowledge the
	       receipt of the LSA by sending a Link State Acknowledgment packet
	       back to the sending neighbor (see Section 13.5), and b) Discard
	       the LSA and examine the next LSA (if any) listed in the Link
	   State Update packet.
	*/
	data := server.NbrConfData.IntfToNbrMap[intf]
	for _, nbrKey := range data {
		nbr := server.NbrConfMap[nbrKey]
		if nbr.State == NbrExchange || nbr.State == NbrLoading {
			continue
		} else {
			send_ack = false
		}
	}
	if send_ack && !exist && lsa_max_age {
		return true
	}
	return false
}

func (server *OSPFV2Server) sanityCheckRouterLsa(rlsa RouterLsa, drlsa RouterLsa, nbr NbrConf, intf IntfConf, exist bool, lsa_max_age bool) (discard bool, op uint8) {
	discard = false
	op = LsdbAdd
	send_ack := server.lsAgeCheck(nbr.IntfKey, lsa_max_age, exist)
	if send_ack {
		op = LsdbNoAction
		discard = true
		server.logger.Info(fmt.Sprintln("LSAUPD: Router LSA Discard. link details", rlsa.LinkDetails, " nbr ", nbr))
		return discard, op
	} else {
		isNew := server.validateLsaIsNew(rlsa.LsaMd, drlsa.LsaMd)
		// TODO check if lsa is installed before MinLSArrival
		if isNew {
			op = FloodLsa
			discard = false
		} else {
			server.logger.Info(fmt.Sprintln("LSAUPD: Router LSA Discard.Already present in lsdb. link details", rlsa.LinkDetails, " nbr ", nbr))
			discard = true
			op = LsdbNoAction
		}
	}

	return discard, op
}
func (server *OSPFV2Server) sanityCheckNetworkLsa(lsaKey LsaKey, nlsa NetworkLsa, dnlsa NetworkLsa, nbr NbrConf, intf IntfConf, exist bool, lsa_max_age bool) (discard bool, op uint8) {
	discard = false
	op = LsdbAdd
	send_ack := server.lsAgeCheck(nbr.IntfKey, lsa_max_age, exist)
	if send_ack {
		op = LsdbNoAction
		discard = true
		server.logger.Info(fmt.Sprintln("LSAUPD: Network LSA Discard. ", " nbr ", nbr))
		return discard, op
	} else {
		isNew := server.validateLsaIsNew(nlsa.LsaMd, dnlsa.LsaMd)
		if isNew {
			op = FloodLsa
			discard = false
		} else {
			discard = true
			op = LsdbNoAction
		}
	}
	//if i am DR and receive nw LSA from neighbor discard it.
	rtr_id := server.globalData.RouterId
	if intf.DRtrId == rtr_id {
		nbrIp := nbr.NbrIP
		if lsaKey.LSId == nbrIp {
			server.logger.Info(fmt.Sprintln("DISCARD: I am dr. received nw LSA from nbr . LSA id ", nbr.NbrIP))
			discard = true
			op = LsdbNoAction
		}
	}
	return discard, op
}
func (server *OSPFV2Server) sanityCheckSummaryLsa(slsa SummaryLsa, dslsa SummaryLsa, nbr NbrConf, intf IntfConf, exist bool, lsa_max_age bool) (discard bool, op uint8) {
	discard = false
	op = LsdbAdd
	send_ack := server.lsAgeCheck(nbr.IntfKey, lsa_max_age, exist)
	if send_ack {
		op = LsdbNoAction
		discard = true
		server.logger.Info(fmt.Sprintln("LSAUPD: Summary LSA Discard. ", " nbr ", nbr))
		return discard, op
	} else {
		isNew := server.validateLsaIsNew(slsa.LsaMd, dslsa.LsaMd)
		if isNew {
			op = FloodLsa
			discard = false
		} else {
			server.logger.Info(fmt.Sprintln("LSAUPD: Discard Summary LSA slsa from nbr"))
			discard = true
			op = LsdbNoAction
		}
	}
	return discard, op
}

func (server *OSPFV2Server) sanityCheckASExternalLsa(alsa ASExternalLsa, dalsa ASExternalLsa, nbr NbrConf, intf IntfConf, exist bool, lsa_max_age bool) (discard bool, op uint8) {
	discard = false
	op = LsdbAdd
	// TODO Reject this lsa if area is configured as stub area.
	send_ack := server.lsAgeCheck(nbr.IntfKey, lsa_max_age, exist)
	if send_ack {
		op = LsdbNoAction
		discard = true
		server.logger.Info(fmt.Sprintln("LSAUPD: As external LSA Discard.", " nbr ", nbr))
		return discard, op
	} else {
		isNew := server.validateLsaIsNew(alsa.LsaMd, dalsa.LsaMd)
		if isNew {
			op = FloodLsa
			discard = false
		} else {
			discard = true
			op = LsdbNoAction
		}
	}
	return discard, op
}

func validateChecksum(data []byte) bool {
	csum := computeFletcherChecksum(data[2:], FLETCHER_CHECKSUM_VALIDATE)
	if csum != 0 {
		//server.logger.Err("LSAUPD: Invalid Router LSA Checksum")
		return false
	}
	return true
}

func (server *OSPFV2Server) validateLsaIsNew(rlsamd LsaMetadata, dlsamd LsaMetadata) bool {
	if rlsamd.LSSequenceNum > dlsamd.LSSequenceNum {
		server.logger.Debug("LSA: received lsseq num > db seq num. ")
		return true
	}
	if rlsamd.LSChecksum > dlsamd.LSChecksum {
		server.logger.Info(fmt.Sprintln("LSA: received lsa checksum > db chceksum "))
		return true
	}
	if rlsamd.LSAge == LSA_MAX_AGE {
		server.logger.Info(fmt.Sprintln("LSA: LSA is maxage "))
		return true
	}
	age_diff := math.Abs(float64(rlsamd.LSAge - dlsamd.LSAge))
	if age_diff > float64(LSA_MAX_AGE_DIFF) &&
		rlsamd.LSAge < rlsamd.LSAge {
		return true
	}
	/* Debug further - currently it doesnt return true for latest LSA */
	return true
}

func (server *OSPFV2Server) GetLsaByteFromLsaKey(areaId uint32, lsaKey LsaKey) []byte {
	lsaByte := []byte{}
	switch lsaKey.LSType {
	case RouterLSA:
		rlsa, valid := server.getRouterLsaFromLsdb(areaId, lsaKey)
		if valid == LsdbEntryNotFound {
			return nil
		}
		lsaByte = encodeRouterLsa(rlsa, lsaKey)
		break
	case NetworkLSA:
		nlsa, valid := server.getNetworkLsaFromLsdb(areaId, lsaKey)
		if valid == LsdbEntryNotFound {
			return nil
		}
		lsaByte = encodeNetworkLsa(nlsa, lsaKey)
		break
	case Summary3LSA, Summary4LSA:
		slsa, valid := server.getSummaryLsaFromLsdb(areaId, lsaKey)
		if valid == LsdbEntryNotFound {
			return nil
		}
		lsaByte = encodeSummaryLsa(slsa, lsaKey)
		break

	case ASExternalLSA:
		alsa, valid := server.getASExternalLsaFromLsdb(areaId, lsaKey)
		if valid == LsdbEntryNotFound {
			return nil
		}
		lsaByte = encodeASExternalLsa(alsa, lsaKey)
		break

	default:
		server.logger.Debug("Flood: Invalid lsa type . ", lsaKey)
		return nil
	}
	return lsaByte
}

func (server *OSPFV2Server) getLsaByteFromLsa(msg LsdbToFloodLSAMsg) []byte {
	lsaByte := []byte{}
	switch msg.LsaKey.LSType {
	case RouterLSA:
		if lsa, ok := msg.LsaData.(RouterLsa); ok {
			server.logger.Debug("Flood: Retrieved router lsa  from lsdb")
			lsaByte = encodeRouterLsa(lsa, msg.LsaKey)
		}
	case NetworkLSA:
		if lsa, ok := msg.LsaData.(NetworkLsa); ok {
			server.logger.Debug("Flood: Retrieved network lsa  from lsdb")
			lsaByte = encodeNetworkLsa(lsa, msg.LsaKey)
		}

	case Summary3LSA, Summary4LSA:
		if lsa, ok := msg.LsaData.(SummaryLsa); ok {
			server.logger.Debug("Flood: Retrieved network lsa  from lsdb")
			lsaByte = encodeSummaryLsa(lsa, msg.LsaKey)
		}

	case ASExternalLSA:
		if lsa, ok := msg.LsaData.(ASExternalLsa); ok {
			server.logger.Debug("Flood: Retrieved as external  lsa  from lsdb")
			lsaByte = encodeASExternalLsa(lsa, msg.LsaKey)
		}
	default:
		server.logger.Err("Flood: Invalid LSA type . Not able to decode message from lsdb ", msg.LsaKey)
	}
	return lsaByte
}

/*
On broadcast networks, the Link State Update packets are
            multicast.  The destination IP address specified for the
            Link State Update Packet depends on the state of the
            interface.  If the interface state is DR or Backup, the
            address AllSPFRouters should be used.  Otherwise, the
            address AllDRouters should be used.
            On non-broadcast networks, separate Link State Update
            packets must be sent, as unicasts, to each adjacent neighbor
            (i.e., those in state Exchange or greater).  The destination
            IP addresses for these packets are the neighbors' IP
            addresses.
*/
func (server *OSPFV2Server) GetDestIpForFlood(intfKey IntfConfKey, nbrIP uint32, nbrMac net.HardwareAddr) (net.IP, net.HardwareAddr, error) {
	var destIp net.IP
	var destMac net.HardwareAddr
	intf, exist := server.IntfConfMap[intfKey]
	if !exist {
		return nil, nil, errors.New("intf conf does not exist ")
	}
	if intf.Type == objects.INTF_TYPE_BROADCAST {
		if intf.DRtrId == server.globalData.RouterId || intf.BDRtrId == server.globalData.RouterId {
			destIp = net.ParseIP(AllSPFRouters)
		} else {
			destIp = net.ParseIP(AllDRouters)
		}
		destMac, _ = net.ParseMAC(McastMAC)
	} else if intf.Type == objects.INTF_TYPE_POINT2POINT {
		destIp = net.ParseIP(AllSPFRouters)
		destMac, _ = net.ParseMAC(ALLSPFROUTERMAC)
	}

	return destIp, destMac, nil
}
