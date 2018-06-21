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

// conversationId.go
package drcp

import (
	"fmt"
	"l2/lacp/protocol/utils"
)

// holds the current conversation map values for the system
// Each DR should be updated based on the information contained
// in this map
var ConversationIdMap [MAX_CONVERSATION_IDS]ConvIdTypeValue

type ConvIdTypeValue struct {
	Valid      bool
	Refcnt     int
	Idtype     GatewayAlgorithm
	Isid       uint32
	Cvlan      uint16
	Svlan      uint16
	Bvid       uint16
	Psuedowire uint32
	PortList   []int32
}

// GetAllCVIDConversations: Fill in the mapping of vlan -> conversation id which is 1:1
func GetAllCVIDConversations() {
	curMark := 0
	count := 100
	more := true
	for more {
		for _, client := range utils.GetAsicDPluginList() {

			bulkVlanInfo, _ := client.GetBulkVlan(curMark, count)
			if bulkVlanInfo != nil {
				objCnt := int(bulkVlanInfo.Count)
				more = bool(bulkVlanInfo.More)
				curMark = int(bulkVlanInfo.EndIdx)
				for i := 0; i < objCnt; i++ {
					vlan := bulkVlanInfo.VlanList[i].VlanId
					ent := ConversationIdMap[uint16(vlan)]
					ent.Valid = true
					ent.Refcnt = 1
					ent.Idtype = GATEWAY_ALGORITHM_CVID
					ent.Cvlan = uint16(vlan)
					if ent.PortList == nil {
						ent.PortList = make([]int32, 0)
					}
					// combine untagged and tagged ports
					ent.PortList = bulkVlanInfo.VlanList[i].UntagIfIndexList
					for _, ifindex := range bulkVlanInfo.VlanList[i].IfIndexList {
						ent.PortList = append(ent.PortList, ifindex)
					}
					//fmt.Println("Creating Conversation Id", ent)
					ConversationIdMap[uint16(vlan)] = ent
				}
			} else {
				more = false
			}
		}
	}
}

// CreateConversationId is a config api to handle conversationId updates
func CreateConversationId(cfg *DRConversationConfig) {

	// only supported converstation at this time
	if cfg.Idtype == GATEWAY_ALGORITHM_CVID {
		if cfg.Cvlan < MAX_CONVERSATION_IDS {
			if ConversationIdMap[cfg.Cvlan].Valid {
				ent := ConversationIdMap[cfg.Cvlan]
				ent.Refcnt++
				// add any new ports into the ConversationIdMap
				for _, p := range cfg.PortList {
					foundEntry := false
					for _, p2 := range ent.PortList {
						if p == p2 {
							foundEntry = true
						}
					}
					if !foundEntry {
						ent.PortList = append(ent.PortList, p)
					}
				}
				ConversationIdMap[cfg.Cvlan] = ent
			} else {
				ent := ConversationIdMap[cfg.Cvlan]
				ent.Valid = true
				ent.Refcnt = 1
				ent.Idtype = GATEWAY_ALGORITHM_CVID
				ent.Cvlan = uint16(cfg.Cvlan)
				ent.PortList = nil
				if cfg.PortList != nil {
					ent.PortList = make([]int32, 0)
				}

				for _, p := range cfg.PortList {
					ent.PortList = append(ent.PortList, p)
				}

				ConversationIdMap[uint16(cfg.Cvlan)] = ent

			}
			// update the local digests and converstaion lists
			for _, dr := range DistributedRelayDBList {
				if dr.DrniName == cfg.DrniName {
					dr.LaDrLog(fmt.Sprintf("Creating Converstaion %d", cfg.Cvlan))
					dr.SetTimeSharingPortAndGatwewayDigest()
				}
			}
		}
	}
}

// CreateConversationId is a config api to handle conversationId updates
func DeleteConversationId(cfg *DRConversationConfig, force bool) {

	// only supported converstation at this time
	if cfg.Idtype == GATEWAY_ALGORITHM_CVID {
		if cfg.Cvlan < MAX_CONVERSATION_IDS && (ConversationIdMap[cfg.Cvlan].Valid || force) {
			ent := ConversationIdMap[cfg.Cvlan]
			if ent.Refcnt > 1 {
				// TODO FUTURE when you can map multiple conversations types
				// to the same conversation id
				ent.Refcnt--
			} else {
				ent.Valid = false
				ent.Refcnt = 0
				ent.Idtype = GATEWAY_ALGORITHM_NULL
				ent.Cvlan = uint16(cfg.Cvlan)
				ent.PortList = nil
				ConversationIdMap[uint16(cfg.Cvlan)] = ent

				// update the local digests and converstaion lists
				for _, dr := range DistributedRelayDBList {
					if dr.DrniName == cfg.DrniName {
						dr.LaDrLog(fmt.Sprintf("Deleting Converstaion %d", cfg.Cvlan))
						dr.SetTimeSharingPortAndGatwewayDigest()
					}
				}
			}
		}
	}
}

// UpdateConversationId is a config api to handle conversationId updates
// likely port id updates
// NOTE: portList should always contain the complete valid port list
func UpdateConversationId(cfg *DRConversationConfig) {

	// only supported converstation at this time
	if cfg.Idtype == GATEWAY_ALGORITHM_CVID {
		if cfg.Cvlan < MAX_CONVERSATION_IDS && ConversationIdMap[cfg.Cvlan].Valid {
			ent := ConversationIdMap[cfg.Cvlan]
			ent.Valid = true
			ent.Idtype = GATEWAY_ALGORITHM_NULL
			ent.Cvlan = uint16(cfg.Cvlan)
			ent.PortList = nil
			ConversationIdMap[uint16(cfg.Cvlan)] = ent
			ent.PortList = nil
			if cfg.PortList != nil {
				ent.PortList = make([]int32, 0)
			}

			// cfg.PortList contains the final valid list so lets just
			// overwrite the port list
			for _, p := range cfg.PortList {
				ent.PortList = append(ent.PortList, p)
			}
			ConversationIdMap[uint16(cfg.Cvlan)] = ent

			// update the local digests and converstaion lists
			for _, dr := range DistributedRelayDBList {
				dr.SetTimeSharingPortAndGatwewayDigest()
			}
		}
	}
}
