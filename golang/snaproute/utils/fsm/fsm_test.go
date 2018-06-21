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

// fsm_test
package fsm_test

import (
	"fmt"
	"testing"
	"utils/fsm"
)

type MyFSM struct {
	FSM *fsm.Machine
}

type MyStateEvent struct {
	PrevState fsm.State
	PrevEvent fsm.Event
	State     fsm.State
	Event     fsm.Event
}

func (se *MyStateEvent) CurrentState() fsm.State         { return se.State }
func (se *MyStateEvent) CurrentEvent() fsm.Event         { return se.Event }
func (se *MyStateEvent) PreviousState() fsm.State        { return se.State }
func (se *MyStateEvent) PreviousEvent() fsm.Event        { return se.Event }
func (se *MyStateEvent) SetState(s fsm.State)            { se.State = s }
func (se *MyStateEvent) SetEvent(es string, e fsm.Event) { se.Event = e }

func (se *MyStateEvent) LoggerSet(log func(string))                 {}
func (se *MyStateEvent) EnableLogging(ena bool)                     {}
func (se *MyStateEvent) IsLoggerEna() bool                          { return false }
func (se *MyStateEvent) StateStrMapSet(strMap map[fsm.State]string) {}

const (
	exampleState1 = iota
	exampleState2
	exampleState3
)
const (
	exampleEvent1 = iota + 1
	exampleEvent2
)

func TestProcessEventNoStartCalled(t *testing.T) {

	rules := fsm.Ruleset{}

	// example rules
	rules.AddRule(exampleState1, exampleEvent1, func(m fsm.Machine, data interface{}) fsm.State { return exampleState2 })
	rules.AddRule(exampleState2, exampleEvent2, func(m fsm.Machine, data interface{}) fsm.State { return exampleState3 })

	myFsm := &MyFSM{FSM: &fsm.Machine{Curr: &MyStateEvent{},
		Rules: &rules}}

	rv := myFsm.FSM.ProcessEvent("FSM", exampleEvent1, nil)
	if rv != fsm.ErrorMachineNotStarted {
		t.Error("Expected Error", fsm.ErrorMachineNotStarted)
	}

	if 0 != myFsm.FSM.Curr.CurrentState() {
		t.Error("Expected state", nil, "\nActual state", myFsm.FSM.Curr.CurrentState())
	}

	if 0 != myFsm.FSM.Curr.CurrentEvent() {
		t.Error("Expected no valid event stored\nActual", myFsm.FSM.Curr.CurrentEvent())
	}

}

func TestAddRuleDuplicateAdd(t *testing.T) {

	rules := fsm.Ruleset{}

	// example rules
	rv := rules.AddRule(exampleState1, exampleEvent1, func(m fsm.Machine, data interface{}) fsm.State { return exampleState2 })
	rv2 := rules.AddRule(exampleState1, exampleEvent2, func(m fsm.Machine, data interface{}) fsm.State { return exampleState3 })
	rv3 := rules.AddRule(exampleState1, exampleEvent2, func(m fsm.Machine, data interface{}) fsm.State { return exampleState3 })

	if rv != nil {
		t.Error("Expected no error")
	}
	if rv2 != nil {
		t.Error("Expected no error")
	}
	if rv3 != fsm.ErrorMachineStateEventExists {
		t.Error("Expected Error", fsm.ErrorMachineStateEventExists)
	}
}

func TestProcessEventBadEventForGivenState(t *testing.T) {

	rules := fsm.Ruleset{}

	// example rules
	rules.AddRule(exampleState1, exampleEvent1, func(m fsm.Machine, data interface{}) fsm.State { return exampleState2 })
	rules.AddRule(exampleState2, exampleEvent2, func(m fsm.Machine, data interface{}) fsm.State { return exampleState3 })

	myFsm := &MyFSM{FSM: &fsm.Machine{Curr: &MyStateEvent{},
		Rules: &rules,
		Begin: false}}

	// start state
	begin := myFsm.FSM.Start(exampleState1)
	fmt.Println("Begin", begin)

	rv := myFsm.FSM.ProcessEvent("FSM", exampleEvent2, nil)
	if rv != fsm.InvalidStateEvent {
		t.Error("Expected Error", fsm.InvalidStateEvent, "\nActual", rv)
	}
	if exampleState1 != myFsm.FSM.Curr.CurrentState() {
		t.Error("Expected state", exampleState1, "\nActual state", myFsm.FSM.Curr.CurrentState())
	}

	if 0 != myFsm.FSM.Curr.CurrentEvent() {
		t.Error("Expected no valid event stored\nActual", myFsm.FSM.Curr.CurrentEvent())
	}

}

func TestProcessEventGoodStateTransition(t *testing.T) {

	rules := fsm.Ruleset{}

	// example rules
	rules.AddRule(exampleState1, exampleEvent1, func(m fsm.Machine, data interface{}) fsm.State { return exampleState2 })
	rules.AddRule(exampleState2, exampleEvent2, func(m fsm.Machine, data interface{}) fsm.State { return exampleState3 })

	myFsm := &MyFSM{FSM: &fsm.Machine{Curr: &MyStateEvent{},
		Rules: &rules}}

	// start state
	myFsm.FSM.Start(exampleState1)

	// First transition
	rv := myFsm.FSM.ProcessEvent("FSM", exampleEvent1, nil)
	if rv != nil {
		t.Error("Expected no error")
	}
	if exampleState2 != myFsm.FSM.Curr.CurrentState() {
		t.Error("Expected state", exampleState2, "\nActual state", myFsm.FSM.Curr.CurrentState())
	}
	if exampleEvent1 != myFsm.FSM.Curr.CurrentEvent() {
		t.Error("Expected event", exampleEvent1, "\nActual event", myFsm.FSM.Curr.CurrentEvent())
	}

	// Second transition
	rv2 := myFsm.FSM.ProcessEvent("FSM", exampleEvent2, nil)
	if rv2 != nil {
		t.Error("Expected no error")
	}
	if exampleState3 != myFsm.FSM.Curr.CurrentState() {
		t.Error("Expected state", exampleState2, "\nActual state", myFsm.FSM.Curr.CurrentState())
	}
	if exampleEvent2 != myFsm.FSM.Curr.CurrentEvent() {
		t.Error("Expected event", exampleEvent2, "\nActual event", myFsm.FSM.Curr.CurrentEvent())
	}

}
