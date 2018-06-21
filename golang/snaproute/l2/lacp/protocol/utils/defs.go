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

// defs
package utils

import (
	"net"
	"strconv"
	"strings"
	"utils/fsm"
)

var LaSwitchMac [6]uint8

// MachineEvent machine events will be sent
// with this struct and will provide extra data
// in order to provide async communication between
// sender and receiver
type MachineEvent struct {
	E            fsm.Event
	Src          string
	ResponseChan chan string
}

func SendResponse(msg string, responseChan chan string) {
	responseChan <- msg
}

type StateEvent struct {
	// current State
	s fsm.State
	// previous State
	ps fsm.State
	// current event
	e fsm.Event
	// previous event
	pe fsm.Event

	// event src
	esrc        string
	Owner       string
	StrStateMap map[fsm.State]string
	LogEna      bool
	Logger      func(string)
}

func (se *StateEvent) LoggerSet(log func(string))                 { se.Logger = log }
func (se *StateEvent) EnableLogging(ena bool)                     { se.LogEna = ena }
func (se *StateEvent) IsLoggerEna() bool                          { return se.LogEna }
func (se *StateEvent) StateStrMapSet(strMap map[fsm.State]string) { se.StrStateMap = strMap }
func (se *StateEvent) PreviousState() fsm.State                   { return se.ps }
func (se *StateEvent) CurrentState() fsm.State                    { return se.s }
func (se *StateEvent) PreviousEvent() fsm.Event                   { return se.pe }
func (se *StateEvent) CurrentEvent() fsm.Event                    { return se.e }
func (se *StateEvent) SetEvent(es string, e fsm.Event) {
	se.esrc = es
	se.pe = se.e
	se.e = e
}
func (se *StateEvent) SetState(s fsm.State) {
	se.ps = se.s
	se.s = s
	if se.IsLoggerEna() && se.ps != se.s {
		se.Logger((strings.Join([]string{"Src", se.esrc, "OldState", se.StrStateMap[se.ps], "Evt", strconv.Itoa(int(se.e)), "NewState", se.StrStateMap[s]}, ":")))
	}
}

func SaveSwitchMac(switchMac string) {
	netAddr, _ := net.ParseMAC(switchMac)
	LaSwitchMac = [6]uint8{netAddr[0], netAddr[1], netAddr[2], netAddr[3], netAddr[4], netAddr[5]}
}

func GetSwitchMac() [6]uint8 {
	return LaSwitchMac
}
