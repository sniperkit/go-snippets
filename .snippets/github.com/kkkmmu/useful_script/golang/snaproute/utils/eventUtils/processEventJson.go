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

package eventUtils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type FaultDetail struct {
	RaiseFault       bool
	ClearingEventId  int
	ClearingDaemonId int
	AlarmSeverity    string
}

type EventStruct struct {
	EventId     int
	EventName   string
	Description string
	SrcObjName  string
	EventEnable bool
	IsFault     bool
	Fault       FaultDetail
}

type DaemonEvent struct {
	DaemonId          int
	DaemonName        string
	DaemonEventEnable bool
	EventBufferSize   int
	EventList         []EventStruct
}

type EventJson struct {
	DaemonEvents []DaemonEvent
}

const (
	EventDir string = "/etc/flexswitch/"
)

func ParseEventsJson() (evtJson EventJson, err error) {
	eventsFile := EventDir + "events.json"
	bytes, err := ioutil.ReadFile(eventsFile)
	if err != nil {
		err := errors.New(fmt.Sprintln("Error in reading ", eventsFile, " file."))
		return evtJson, err
	}

	err = json.Unmarshal(bytes, &evtJson)
	if err != nil {
		err := errors.New(fmt.Sprintln("Errors in unmarshalling json file: ", eventsFile))
		return evtJson, err
	}
	return evtJson, err
}
