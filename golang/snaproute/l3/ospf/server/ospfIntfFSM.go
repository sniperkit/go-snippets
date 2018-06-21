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
	"encoding/binary"
	"fmt"
	"l3/ospf/config"
	"time"
)

func (server *OSPFServer) StartOspfIntfFSM(key IntfConfKey) {
	ent, ok := server.IntfConfMap[key]
	if !ok {
		server.logger.Err(fmt.Sprintln("IntfFSM: IntfConfMap doesnt exist for key ", key))
		return
	}
	areaId := convertIPv4ToUint32(ent.IfAreaId)
	msg := NetworkLSAChangeMsg{
		areaId:  areaId,
		intfKey: key,
	}

	server.logger.Info("Sending msg for router LSA generation")
	server.IntfStateChangeCh <- msg

	if ent.IfType == config.NumberedP2P || ent.IfType == config.UnnumberedP2P {
		server.StartOspfP2PIntfFSM(key)
	} else if ent.IfType == config.Broadcast {
		server.StartOspfBroadcastIntfFSM(key)
	}
}

func (server *OSPFServer) StartOspfP2PIntfFSM(key IntfConfKey) {
	server.StartSendHelloPkt(key)
	for {
		ent, _ := server.IntfConfMap[key]
		select {
		case <-ent.HelloIntervalTicker.C:
			server.StartSendHelloPkt(key)
		case createMsg := <-ent.NeighCreateCh:
			if bytesEqual(createMsg.DRtr, []byte{0, 0, 0, 0}) == false ||
				bytesEqual(createMsg.BDRtr, []byte{0, 0, 0, 0}) == false {
				server.logger.Err("DR or BDR is non zero")
				continue
			}
			ipad := convertUint32ToIPv4(createMsg.NbrIP)
			ospfNbrConfKey := NeighborConfKey{
				IPAddr:  config.IpAddress(ipad),
				IntfIdx: key.IntfIdx,
			}

			neighborEntry, exist := ent.NeighborMap[ospfNbrConfKey]
			if !exist {
				neighborEntry.NbrIP = createMsg.NbrIP
				neighborEntry.TwoWayStatus = createMsg.TwoWayStatus
				neighborEntry.RtrPrio = createMsg.RtrPrio
				neighborEntry.FullState = false
				ent.NeighborMap[ospfNbrConfKey] = neighborEntry
				server.IntfConfMap[key] = ent
				server.logger.Info(fmt.Sprintln("1 IntfConf neighbor entry", server.IntfConfMap[key].NeighborMap, "neighborKey:", ospfNbrConfKey))
			}
		case changeMsg := <-ent.NeighChangeCh:
			if bytesEqual(changeMsg.DRtr, []byte{0, 0, 0, 0}) == false ||
				bytesEqual(changeMsg.BDRtr, []byte{0, 0, 0, 0}) == false {
				server.logger.Err("DR or BDR is non zero")
				continue
			}
			ipad := convertUint32ToIPv4(changeMsg.NbrIP)
			ospfNbrConfKey := NeighborConfKey{
				IPAddr:  config.IpAddress(ipad),
				IntfIdx: key.IntfIdx,
			}
			neighborEntry, exist := ent.NeighborMap[ospfNbrConfKey]
			if exist {
				server.logger.Info(fmt.Sprintln("Change msg: ", changeMsg, "neighbor entry:", neighborEntry, "neighbor key:", ospfNbrConfKey))
				neighborEntry.NbrIP = changeMsg.NbrIP
				neighborEntry.TwoWayStatus = changeMsg.TwoWayStatus
				neighborEntry.RtrPrio = changeMsg.RtrPrio
				neighborEntry.DRtr = changeMsg.DRtr
				neighborEntry.BDRtr = changeMsg.BDRtr
				ent.NeighborMap[ospfNbrConfKey] = neighborEntry
				server.IntfConfMap[key] = ent
				server.logger.Info(fmt.Sprintln("2 IntfConf neighbor entry", server.IntfConfMap[key].NeighborMap))
			} else {
				server.logger.Err(fmt.Sprintln("Neighbor entry does not exists", ospfNbrConfKey.IPAddr))
			}
		case nbrStateChangeMsg := <-ent.NbrStateChangeCh:
			// Only when Neighbor Went Down from TwoWayStatus
			server.logger.Info(fmt.Sprintf("Recev Neighbor State Change message", nbrStateChangeMsg))
			server.processNbrDownEvent(nbrStateChangeMsg, key, true)
		case state := <-ent.FSMCtrlCh:
			if state == false {
				server.StopSendHelloPkt(key)
				ent.FSMCtrlStatusCh <- false
				return
			}
		}
	}

}

func (server *OSPFServer) StartOspfBroadcastIntfFSM(key IntfConfKey) {
	server.StartSendHelloPkt(key)
	for {
		ent, _ := server.IntfConfMap[key]
		select {
		case <-ent.HelloIntervalTicker.C:
			server.StartSendHelloPkt(key)
		case <-ent.WaitTimer.C:
			server.logger.Info("Wait timer expired")
			eventInfo := "Wait time expired for "
			server.AddOspfEventState(config.INTF, eventInfo)
			//server.IntfConfMap[key] = ent
			// Elect BDR And DR
			server.ElectBDRAndDR(key)
		case msg := <-ent.BackupSeenCh:
			server.logger.Info(fmt.Sprintf("Transit to action state because of backup seen", msg))
			server.AddOspfEventState(config.INTF, "Backup seen")
			server.ElectBDRAndDR(key)
		case createMsg := <-ent.NeighCreateCh:

			ipad := convertUint32ToIPv4(createMsg.NbrIP)
			ospfNbrConfKey := NeighborConfKey{
				IPAddr:  config.IpAddress(ipad),
				IntfIdx: key.IntfIdx,
			}

			neighborEntry, exist := ent.NeighborMap[ospfNbrConfKey]
			if !exist {
				neighborEntry.NbrIP = createMsg.NbrIP
				neighborEntry.TwoWayStatus = createMsg.TwoWayStatus
				neighborEntry.RtrPrio = createMsg.RtrPrio
				neighborEntry.DRtr = createMsg.DRtr
				neighborEntry.BDRtr = createMsg.BDRtr
				neighborEntry.FullState = false
				ent.NeighborMap[ospfNbrConfKey] = neighborEntry
				server.IntfConfMap[key] = ent
				server.logger.Info(fmt.Sprintln("1 IntfConf neighbor entry", server.IntfConfMap[key].NeighborMap, "neighborKey:", ospfNbrConfKey))
				if createMsg.TwoWayStatus == true &&
					ent.IfFSMState > config.Waiting {
					server.ElectBDRAndDR(key)
				}
			}
		case changeMsg := <-ent.NeighChangeCh:
			ipad := convertUint32ToIPv4(changeMsg.NbrIP)
			ospfNbrConfKey := NeighborConfKey{
				IPAddr:  config.IpAddress(ipad),
				IntfIdx: key.IntfIdx,
			}
			neighborEntry, exist := ent.NeighborMap[ospfNbrConfKey]
			if exist {
				server.logger.Info(fmt.Sprintln("Change msg: ", changeMsg, "neighbor entry:", neighborEntry, "neighbor key:", ospfNbrConfKey))
				//rtrId := changeMsg.RouterId
				NbrIP := changeMsg.NbrIP
				oldRtrPrio := neighborEntry.RtrPrio
				oldDRtr := binary.BigEndian.Uint32(neighborEntry.DRtr)
				oldBDRtr := binary.BigEndian.Uint32(neighborEntry.BDRtr)
				newDRtr := binary.BigEndian.Uint32(changeMsg.DRtr)
				newBDRtr := binary.BigEndian.Uint32(changeMsg.BDRtr)
				oldTwoWayStatus := neighborEntry.TwoWayStatus
				neighborEntry.NbrIP = changeMsg.NbrIP
				neighborEntry.TwoWayStatus = changeMsg.TwoWayStatus
				neighborEntry.RtrPrio = changeMsg.RtrPrio
				neighborEntry.DRtr = changeMsg.DRtr
				neighborEntry.BDRtr = changeMsg.BDRtr
				ent.NeighborMap[ospfNbrConfKey] = neighborEntry
				server.IntfConfMap[key] = ent
				server.logger.Info(fmt.Sprintln("2 IntfConf neighbor entry", server.IntfConfMap[key].NeighborMap))
				if ent.IfFSMState > config.Waiting {
					// RFC2328 Section 9.2 (Neighbor Change Event)
					if (oldDRtr == NbrIP && newDRtr != NbrIP && oldTwoWayStatus == true) ||
						(oldDRtr != NbrIP && newDRtr == NbrIP && oldTwoWayStatus == true) ||
						(oldBDRtr == NbrIP && newBDRtr != NbrIP && oldTwoWayStatus == true) ||
						(oldBDRtr != NbrIP && newBDRtr == NbrIP && oldTwoWayStatus == true) ||
						(oldTwoWayStatus != changeMsg.TwoWayStatus) ||
						(oldRtrPrio != changeMsg.RtrPrio && oldTwoWayStatus == true) {

						// Update Neighbor and Re-elect BDR And DR
						server.ElectBDRAndDR(key)
					}
				}
			}
		case nbrStateChangeMsg := <-ent.NbrStateChangeCh:
			// Only when Neighbor Went Down from TwoWayStatus
			// Todo: Handle NbrIP: Ashutosh
			server.logger.Info(fmt.Sprintf("Recev Neighbor State Change message", nbrStateChangeMsg))
			server.processNbrDownEvent(nbrStateChangeMsg, key, false)
		case state := <-ent.FSMCtrlCh:
			if state == false {
				server.StopSendHelloPkt(key)
				ent.FSMCtrlStatusCh <- false
				return
			}
		case msg := <-ent.NbrFullStateCh:
			// Note : NBR State Machine should only send message if
			// NBR State changes to/from Full (but not to Down)
			server.processNbrFullStateMsg(msg, key)
		}
	}
}

func (server *OSPFServer) processNbrDownEvent(msg NbrStateChangeMsg,
	key IntfConfKey, p2p bool) {
	ent, _ := server.IntfConfMap[key]

	neighborEntry, exist := ent.NeighborMap[msg.nbrKey]
	if exist {
		oldTwoWayStatus := neighborEntry.TwoWayStatus
		delete(ent.NeighborMap, msg.nbrKey)
		server.logger.Info(fmt.Sprintln("Deleting", msg.nbrKey))
		server.IntfConfMap[key] = ent
		if p2p == false {
			if ent.IfFSMState > config.Waiting {
				// RFC2328 Section 9.2 (Neighbor Change Event)
				/* Investigate - if neighbor goes to dead from full
				   to dead status oldTwoWayStatus is not true */
				oldTwoWayStatus = true // temp fix
				if oldTwoWayStatus == true {
					server.logger.Info(fmt.Sprintln("deleting nbr, call dr/bdr election."))
					server.ElectBDRAndDR(key)
				} else {
					server.logger.Info("Dont call elect DR/BDR as neighbr was not in 2 way")
				}
			}
		}
	}
}

// Nbr State machine has to send FullState change msg and then
// Send interface state down event
func (server *OSPFServer) processNbrFullStateMsg(msg NbrFullStateMsg,
	key IntfConfKey) {
	ent, _ := server.IntfConfMap[key]
	//areaId := convertIPv4ToUint32(ent.IfAreaId)
	if msg.FullState == true {
		server.logger.Info("Neighbor State changed to full state")
	} else {
		server.logger.Info("Neighbor State changed from full state")
	}

	nbrEntry, exist := ent.NeighborMap[msg.nbrKey]
	if exist {
		if msg.FullState != nbrEntry.FullState &&
			ent.IfFSMState == config.DesignatedRouter {
			nbrEntry.FullState = msg.FullState
			ent.NeighborMap[msg.nbrKey] = nbrEntry
			server.IntfConfMap[key] = ent
			/*lsaMsg := NetworkLSAChangeMsg{
				areaId:  areaId,
				intfKey: key,
			}*/
			//	server.CreateNetworkLSACh <- lsaMsg
		}
		if msg.FullState {
			/*
				msg := nbrStateChangeMsg{
					key: nbrKey.RouterId,
					areaId: areaId,
				} */
			//	server.neighborStateChangeCh <- msg
		}
	}
}

func (server *OSPFServer) ElectBDR(key IntfConfKey) ([]byte, uint32) {
	ent, _ := server.IntfConfMap[key]
	electedBDR := []byte{0, 0, 0, 0}
	var electedRtrPrio uint8
	var electedRtrId uint32
	var MaxRtrPrio uint8
	var RtrIdWithMaxPrio uint32
	var NbrIPWithMaxPrio uint32

	for nbrkey, nbrEntry := range ent.NeighborMap {
		nbrConf := server.NeighborConfigMap[nbrkey]
		if nbrEntry.TwoWayStatus == true &&
			nbrEntry.RtrPrio > 0 &&
			nbrEntry.NbrIP != 0 {
			tempDR := binary.BigEndian.Uint32(nbrEntry.DRtr)
			if tempDR == nbrEntry.NbrIP {
				continue
			}
			tempBDR := binary.BigEndian.Uint32(nbrEntry.BDRtr)
			if tempBDR == nbrEntry.NbrIP {
				if nbrEntry.RtrPrio > electedRtrPrio {
					electedRtrPrio = nbrEntry.RtrPrio
					electedRtrId = nbrConf.OspfNbrRtrId
					electedBDR = nbrEntry.BDRtr
				} else if nbrEntry.RtrPrio == electedRtrPrio {
					if electedRtrId < nbrConf.OspfNbrRtrId {
						electedRtrPrio = nbrEntry.RtrPrio
						electedRtrId = nbrConf.OspfNbrRtrId
						electedBDR = nbrEntry.BDRtr
					}
				}
			}
			if MaxRtrPrio < nbrEntry.RtrPrio {
				MaxRtrPrio = nbrEntry.RtrPrio
				RtrIdWithMaxPrio = nbrConf.OspfNbrRtrId
				NbrIPWithMaxPrio = nbrEntry.NbrIP
			} else if MaxRtrPrio == nbrEntry.RtrPrio {
				if RtrIdWithMaxPrio < nbrConf.OspfNbrRtrId {
					MaxRtrPrio = nbrEntry.RtrPrio
					RtrIdWithMaxPrio = nbrConf.OspfNbrRtrId
					NbrIPWithMaxPrio = nbrEntry.NbrIP
				}
			}
		}
	}

	if ent.IfRtrPriority != 0 &&
		bytesEqual(ent.IfIpAddr.To4(), []byte{0, 0, 0, 0}) == false {
		if bytesEqual(ent.IfIpAddr.To4(), ent.IfDRIp) == false {
			if bytesEqual(ent.IfIpAddr.To4(), ent.IfBDRIp) == true {
				rtrId := binary.BigEndian.Uint32(server.ospfGlobalConf.RouterId)
				if ent.IfRtrPriority > electedRtrPrio {
					electedRtrPrio = ent.IfRtrPriority
					electedRtrId = rtrId
					electedBDR = ent.IfIpAddr.To4()
				} else if ent.IfRtrPriority == electedRtrPrio {
					if electedRtrId < rtrId {
						electedRtrPrio = ent.IfRtrPriority
						electedRtrId = rtrId
						electedBDR = ent.IfIpAddr.To4()
					}
				}
			}

			tempRtrId := binary.BigEndian.Uint32(server.ospfGlobalConf.RouterId)
			if MaxRtrPrio < ent.IfRtrPriority {
				MaxRtrPrio = ent.IfRtrPriority
				NbrIPWithMaxPrio = binary.BigEndian.Uint32(ent.IfIpAddr.To4())
				RtrIdWithMaxPrio = tempRtrId
			} else if MaxRtrPrio == ent.IfRtrPriority {
				if RtrIdWithMaxPrio < tempRtrId {
					MaxRtrPrio = ent.IfRtrPriority
					NbrIPWithMaxPrio = binary.BigEndian.Uint32(ent.IfIpAddr.To4())
					RtrIdWithMaxPrio = tempRtrId
				}
			}

		}
	}
	if bytesEqual(electedBDR, []byte{0, 0, 0, 0}) == true {
		binary.BigEndian.PutUint32(electedBDR, NbrIPWithMaxPrio)
		electedRtrId = RtrIdWithMaxPrio
	}

	return electedBDR, electedRtrId
}

func (server *OSPFServer) ElectDR(key IntfConfKey, electedBDR []byte, electedBDRtrId uint32) ([]byte, uint32) {
	ent, _ := server.IntfConfMap[key]
	electedDR := []byte{0, 0, 0, 0}
	var electedRtrPrio uint8
	var electedDRtrId uint32

	for key, nbrEntry := range ent.NeighborMap {
		nbrConf := server.NeighborConfigMap[key]
		if nbrEntry.TwoWayStatus == true &&
			nbrEntry.RtrPrio > 0 &&
			nbrEntry.NbrIP != 0 {
			tempDR := binary.BigEndian.Uint32(nbrEntry.DRtr)
			if tempDR == nbrEntry.NbrIP {
				if nbrEntry.RtrPrio > electedRtrPrio {
					electedRtrPrio = nbrEntry.RtrPrio
					electedDRtrId = nbrConf.OspfNbrRtrId
					electedDR = nbrEntry.DRtr
				} else if nbrEntry.RtrPrio == electedRtrPrio {
					if electedDRtrId < nbrConf.OspfNbrRtrId {
						electedRtrPrio = nbrEntry.RtrPrio
						electedDRtrId = nbrConf.OspfNbrRtrId
						electedDR = nbrEntry.DRtr
					}
				}
			}
		}
	}

	if ent.IfRtrPriority > 0 &&
		bytesEqual(ent.IfIpAddr.To4(), []byte{0, 0, 0, 0}) == false {
		if bytesEqual(ent.IfIpAddr.To4(), ent.IfDRIp) == true {
			rtrId := binary.BigEndian.Uint32(server.ospfGlobalConf.RouterId)
			if ent.IfRtrPriority > electedRtrPrio {
				electedRtrPrio = ent.IfRtrPriority
				electedDRtrId = rtrId
				electedDR = ent.IfIpAddr.To4()
			} else if ent.IfRtrPriority == electedRtrPrio {
				if electedDRtrId < rtrId {
					electedRtrPrio = ent.IfRtrPriority
					electedDRtrId = rtrId
					electedDR = ent.IfIpAddr.To4()
				}
			}
		}
	}

	if bytesEqual(electedDR, []byte{0, 0, 0, 0}) == true {
		electedDR = electedBDR
		electedDRtrId = electedBDRtrId
	}
	return electedDR, electedDRtrId
}

func (server *OSPFServer) ElectBDRAndDR(key IntfConfKey) {
	ent, _ := server.IntfConfMap[key]
	server.logger.Info(fmt.Sprintln("Election of BDR andDR", ent.IfFSMState))

	oldDRtrId := ent.IfDRtrId
	oldBDRtrId := ent.IfBDRtrId
	//oldBDR := ent.IfBDRIp
	oldState := ent.IfFSMState
	var newState config.IfState

	electedBDR, electedBDRtrId := server.ElectBDR(key)
	ent.IfBDRIp = electedBDR
	ent.IfBDRtrId = electedBDRtrId
	electedDR, electedDRtrId := server.ElectDR(key, electedBDR, electedBDRtrId)
	ent.IfDRIp = electedDR
	ent.IfDRtrId = electedDRtrId
	if bytesEqual(ent.IfDRIp, ent.IfIpAddr.To4()) == true {
		newState = config.DesignatedRouter
	} else if bytesEqual(ent.IfBDRIp, ent.IfIpAddr.To4()) == true {
		newState = config.BackupDesignatedRouter
	} else {
		newState = config.OtherDesignatedRouter
	}

	server.logger.Info(fmt.Sprintln("1. Election of BDR:", ent.IfBDRIp, " and DR:", ent.IfDRIp, "new State:", newState, "DR Id:", ent.IfDRtrId, "BDR Id:", ent.IfBDRtrId))
	server.IntfConfMap[key] = ent

	if newState != oldState &&
		!(newState == config.OtherDesignatedRouter &&
			oldState < config.OtherDesignatedRouter) {
		ent, _ = server.IntfConfMap[key]
		electedBDR, electedBDRtrId = server.ElectBDR(key)
		ent.IfBDRIp = electedBDR
		ent.IfBDRtrId = electedBDRtrId
		electedDR, electedDRtrId = server.ElectDR(key, electedBDR, electedBDRtrId)
		ent.IfDRIp = electedDR
		ent.IfDRtrId = electedDRtrId
		if bytesEqual(ent.IfDRIp, ent.IfIpAddr.To4()) == true {
			newState = config.DesignatedRouter
		} else if bytesEqual(ent.IfBDRIp, ent.IfIpAddr.To4()) == true {
			newState = config.BackupDesignatedRouter
		} else {
			newState = config.OtherDesignatedRouter
		}
		server.logger.Info(fmt.Sprintln("2. Election of BDR:", ent.IfBDRIp, " and DR:", ent.IfDRIp, "new State:", newState, "DR Id:", ent.IfDRtrId, "BDR Id:", ent.IfBDRtrId))
		server.IntfConfMap[key] = ent
	}

	server.createAndSendEventsIntfFSM(key, oldState, newState, oldDRtrId, oldBDRtrId)
}

func (server *OSPFServer) createAndSendEventsIntfFSM(key IntfConfKey,
	oldState config.IfState, newState config.IfState, oldDRtrId uint32,
	oldBDRtrId uint32) {
	ent, _ := server.IntfConfMap[key]
	ent.IfFSMState = newState
	// Need to Check: do we need to add events even when we
	// come back to same state after DR or BDR Election
	ent.IfEvents = ent.IfEvents + 1
	server.IntfConfMap[key] = ent
	server.logger.Info(fmt.Sprintln("Final Election of BDR:", ent.IfBDRIp, " and DR:", ent.IfDRIp, "new State:", newState))

	areaId := convertIPv4ToUint32(ent.IfAreaId)

	msg1 := DrChangeMsg{
		areaId:   areaId,
		intfKey:  key,
		oldstate: oldState,
		newstate: newState,
	}

	server.logger.Info("DRBDR changed. Sending message for router/network LSA generation")
	server.NetworkDRChangeCh <- msg1
	server.logger.Info(fmt.Sprintln("oldState", oldState, " newState", newState))

}

func (server *OSPFServer) StopOspfIntfFSM(key IntfConfKey) {
	ent, _ := server.IntfConfMap[key]
	ent.FSMCtrlCh <- false
	cnt := 0
	for {
		select {
		case status := <-ent.FSMCtrlStatusCh:
			if status == false { // False Means Trans Pkt Thread Stopped
				server.logger.Info("Stopped Sending Hello Pkt")
				return
			}
		default:
			time.Sleep(time.Duration(10) * time.Millisecond)
			cnt = cnt + 1
			if cnt == 100 {
				server.logger.Err("Unable to stop the Tx thread")
				return
			}
		}
	}
}
