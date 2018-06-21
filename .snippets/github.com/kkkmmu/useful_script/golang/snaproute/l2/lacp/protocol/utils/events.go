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
package utils

import (
	"fmt"
	"models/events"
	"utils/eventUtils"
)

type ifindex_event struct {
	ifindex int32
	event   events.EventId
}

// events will hold negitive events and will be cleared
//
var EventMap map[ifindex_event]bool

func CreateEventMap(ifindex int32) {

	evt := ifindex_event{
		ifindex: ifindex,
		event:   events.LacpdEventPortOperStateDown,
	}

	EventMap[evt] = false
	evt.event = events.LacpdEventGroupOperStateDown
	EventMap[evt] = false
	evt.event = events.LacpdEventPortPartnerInfoMismatch
	EventMap[evt] = false

}

func DeleteEventMap(ifindex int32) {
	evt := ifindex_event{
		ifindex: ifindex,
		event:   events.LacpdEventPortOperStateDown,
	}

	delete(EventMap, evt)
	evt.event = events.LacpdEventGroupOperStateDown
	delete(EventMap, evt)
	evt.event = events.LacpdEventPortPartnerInfoMismatch
	delete(EventMap, evt)
}

func ProcessLacpGroupOperStateDown(ifindex int32) {
	intfref := GetAggNameFromIfIndex(ifindex)

	if intfref != "" {
		evt := ifindex_event{
			ifindex: ifindex,
			event:   events.LacpdEventGroupOperStateDown,
		}

		if isset, ok := EventMap[evt]; ok {
			if !isset {
				EventMap[evt] = true
				evtKey := events.LacpEntryKey{
					IntfRef: intfref,
				}
				txEvent := eventUtils.TxEvent{
					EventId: events.LacpdEventGroupOperStateDown,
					Key:     evtKey,
				}
				err := eventUtils.PublishEvents(&txEvent)
				if err != nil {
					GlobalLogger.Err("Error in publishing LacpdEventPortOperStateDown Event")
				}
			}
		}
	} else {
		GlobalLogger.Err(fmt.Sprintf("Error in publishing LacpdEventPortOperStateDown Event, ifindex %d not found", ifindex))
	}
}

func ProcessLacpGroupOperStateUp(ifindex int32) {
	intfref := GetAggNameFromIfIndex(ifindex)

	if intfref != "" {
		evt := ifindex_event{
			ifindex: ifindex,
			event:   events.LacpdEventGroupOperStateDown,
		}
		if isset, ok := EventMap[evt]; ok {
			if isset {
				EventMap[evt] = false
				evtKey := events.LacpEntryKey{
					IntfRef: intfref,
				}
				txEvent := eventUtils.TxEvent{
					EventId: events.LacpdEventGroupOperStateUp,
					Key:     evtKey,
				}
				err := eventUtils.PublishEvents(&txEvent)
				if err != nil {
					GlobalLogger.Err("Error in publishing LacpdEventGroupOperStateUp Event")
				}
			}
		}
	} else {
		GlobalLogger.Err(fmt.Sprintf("Error in publishing LacpdEventGroupOperStateUp Event, ifindex %d not found", ifindex))
	}
}

func ProcessLacpPortOperStateDown(ifindex int32) {
	intfref := GetNameFromIfIndex(ifindex)

	if intfref != "" {
		evt := ifindex_event{
			ifindex: ifindex,
			event:   events.LacpdEventPortOperStateDown,
		}

		if isset, ok := EventMap[evt]; ok {
			if !isset {
				EventMap[evt] = true

				evtKey := events.LacpPortEntryKey{
					IntfRef: intfref,
				}
				txEvent := eventUtils.TxEvent{
					EventId: events.LacpdEventPortOperStateDown,
					Key:     evtKey,
				}
				err := eventUtils.PublishEvents(&txEvent)
				if err != nil {
					GlobalLogger.Err("Error in publishing LacpdEventPortOperStateDown Event")
				}
			}
		}
	} else {
		GlobalLogger.Err(fmt.Sprintf("Error in publishing LacpdEventPortOperStateDown Event, ifindex %d not found", ifindex))
	}
}

func ProcessLacpPortOperStateUp(ifindex int32) {
	intfref := GetNameFromIfIndex(ifindex)

	if intfref != "" {
		evt := ifindex_event{
			ifindex: ifindex,
			event:   events.LacpdEventPortOperStateDown,
		}
		if isset, ok := EventMap[evt]; ok {
			if isset {
				EventMap[evt] = false
				evtKey := events.LacpPortEntryKey{
					IntfRef: intfref,
				}
				txEvent := eventUtils.TxEvent{
					EventId: events.LacpdEventPortOperStateUp,
					Key:     evtKey,
				}
				err := eventUtils.PublishEvents(&txEvent)
				if err != nil {
					GlobalLogger.Err("Error in publishing LacpdEventPortOperStateUp Event")
				}
			}
		}
	} else {
		GlobalLogger.Err(fmt.Sprintf("Error in publishing LacpdEventPortOperStateUp Event, ifindex %d not found", ifindex))
	}
}

func ProcessLacpPortPartnerInfoMismatch(ifindex int32) {
	intfref := GetNameFromIfIndex(ifindex)

	if intfref != "" {
		evt := ifindex_event{
			ifindex: ifindex,
			event:   events.LacpdEventPortPartnerInfoMismatch,
		}

		if isset, ok := EventMap[evt]; ok {
			if !isset {
				EventMap[evt] = true
				evtKey := events.LacpPortEntryKey{
					IntfRef: intfref,
				}
				txEvent := eventUtils.TxEvent{
					EventId: events.LacpdEventPortPartnerInfoMismatch,
					Key:     evtKey,
				}
				err := eventUtils.PublishEvents(&txEvent)
				if err != nil {
					GlobalLogger.Err("Error in publishing LacpdEventPortPartnerInfoMismatch Event")
				}
			}
		}
	} else {
		GlobalLogger.Err(fmt.Sprintf("Error in publishing LacpdEventPortPartnerInfoMismatch Event, ifindex %d not found", ifindex))
	}
}

func ProcessLacpPortPartnerInfoSync(ifindex int32) {
	intfref := GetNameFromIfIndex(ifindex)

	if intfref != "" {
		evt := ifindex_event{
			ifindex: ifindex,
			event:   events.LacpdEventPortPartnerInfoMismatch,
		}

		if isset, ok := EventMap[evt]; ok {
			if isset {
				EventMap[evt] = false
				evtKey := events.LacpPortEntryKey{
					IntfRef: intfref,
				}
				txEvent := eventUtils.TxEvent{
					EventId: events.LacpdEventPortPartnerInfoSync,
					Key:     evtKey,
				}
				err := eventUtils.PublishEvents(&txEvent)
				if err != nil {
					GlobalLogger.Err("Error in publishing LacpdEventPortPartnerInfoSync Event")
				}
			}
		}
	} else {
		GlobalLogger.Err(fmt.Sprintf("Error in publishing LacpdEventPortPartnerInfoSync Event, ifindex %d not found", ifindex))
	}
}
