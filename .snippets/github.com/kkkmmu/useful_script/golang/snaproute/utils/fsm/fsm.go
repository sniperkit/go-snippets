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

package fsm

import "errors"

// State defines the users state type
type State int

// Event defines the users state type
type Event int

// Key identifies the key of for the RulesSet map
type Key struct {
	currState State
	event     Event
}

type Callback func(m Machine, data interface{}) State

var (
	ErrorMachineNotStarted       = errors.New("FSM: ERROR Start() not called")
	InvalidStateEvent            = errors.New("FSM: ERROR Invalid FSM State-Event")
	ErrorMachineStateEventExists = errors.New("FSM: ERROR FSM State-Event already exists. FSM only supports one callback")
)

// StateEvent is the key to callbacks
type StateEvent interface {
	Current() State
	Signal() Event
}

// FSMKey implements the StateEvent interface; it provides a default
// implementation of a Key.
type FSMKey struct {
	S State
	E Event
}

func (k FSMKey) Current() State { return k.S }
func (k FSMKey) Signal() Event  { return k.E }

// Ruleset stores the rules for the state machine.
type Ruleset map[StateEvent]Callback

// AddRule Adds the rules for the callbacks
func (r Ruleset) AddRule(s State, e Event, cb Callback) error {
	k := FSMKey{s, e}
	if _, ok := r[k]; ok {
		// not adding rule
		return ErrorMachineStateEventExists
	}
	r[k] = cb
	return nil
}

// Stater can be passed into the FSM. The Stater is reponsible for setting
// its own default state. Behavior of a Stater without a State is undefined.
type MachineState interface {
	CurrentState() State
	CurrentEvent() Event
	PreviousState() State
	PreviousEvent() Event
	SetState(State)
	SetEvent(string, Event)
	LoggerSet(func(string))
	IsLoggerEna() bool
	EnableLogging(bool)
	StateStrMapSet(map[State]string)
	// TODO History(State, Event)
}

// Machine is a pairing of Rules and a Subject.
// The subject or rules may be changed at any time within
// the machine's lifecycle.
type Machine struct {
	Begin bool
	Curr  MachineState
	Rules *Ruleset
}

// ProcessEvent will attemt to call a callback based on
// the current state of the machine and the event passed in
// dbdata will be called as an input to the callback func
func (m *Machine) ProcessEvent(es string, e Event, cbdata interface{}) error {

	if !m.Begin {
		return ErrorMachineNotStarted
	}
	k := FSMKey{m.Curr.CurrentState(), e}
	r := *m.Rules

	if f, ok := r[k]; ok {
		// save off current event
		m.Curr.SetEvent(es, e)
		// callbacks responsibility to return current state
		m.Curr.SetState(f(*m, cbdata))
		return nil
	}

	return InvalidStateEvent
}

// Start initializes the state machine with
// an initial state and allows for processing of events
// to occur
func (m *Machine) Start(s State) bool {
	m.Curr.SetState(s)
	m.Begin = true
	return m.Begin
}

// New initializes a machine
func New(opts ...func(*Machine)) Machine {
	var m Machine

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

// WithMachineState is intended to be passed to New to set the MachineState
func WithMachineState(ms MachineState) func(*Machine) {
	return func(m *Machine) {
		m.Curr = ms
	}
}

// WithRules is intended to be passed to New to set the Rules
func WithRules(r Ruleset) func(*Machine) {
	return func(m *Machine) {
		m.Rules = &r
	}
}
