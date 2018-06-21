// vtepFSM.go
// File contains the FSM related to setting up a VTEP endpoint
package vxlan

import (
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
	"utils/fsm"
)

// VxlanVtepMachineModuleStr FSM name
const VxlanVtepMachineModuleStr = "VTEPM"

// TODO, what should this time be ??
var retrytime = time.Second * 3

// VxlanVtepState are the states being handled to
// configure a VTEP.
// Init - Initial configuration received
// Disabled - Vtep interface or port is disabled
// Detached - Vtep has been provisioned without a proper vxlan vni reference
// Interface - Src Inteface has been resolved this means that the src ip/ src mac
//             of the VTEP has been discovered
// NextHopInfo - RIB has resolved the next hop info so that the packet can reach
//               the VTEP's destination IP
// ResolveNextHopMac - Need to resolve the dstmac of the VTEP, this will come from
//                     ARP
// Hwconfig - Configure the VTEP information into the HW
// Start - Start listening for traffic on the linux vtep interface for packets that
//         need to have the vxlan header added by VXLAND.  Similarly listen
//         for VXLAN packets on the next hop interface.
//
//
// NOTE: Order here matters, because it will allow for shorter state checking when
//       wanting to check if a state is less than or greater than a certain state
const (
	VxlanVtepStateNone = iota + 1
	VxlanVtepStateDisabled
	VxlanVtepStateInit
	VxlanVtepStateDetached
	VxlanVtepStateInterface
	VxlanVtepStateNextHopInfo
	VxlanVtepStateResolveNextHopMac
	VxlanVtepStateHwConfig
	VxlanVtepStateStart
)

// VxlanVtepStateStrMap is a map that holds the integer state number to string map
var VxlanVtepStateStrMap map[fsm.State]string

// VxlanVtepState map converts state to string
func VxlanVtepMachineStrStateMapInit() {
	VxlanVtepStateStrMap = make(map[fsm.State]string)
	VxlanVtepStateStrMap[VxlanVtepStateNone] = "UNINITIALIZED"
	VxlanVtepStateStrMap[VxlanVtepStateDisabled] = "DISABLED"
	VxlanVtepStateStrMap[VxlanVtepStateInit] = "INIT"
	VxlanVtepStateStrMap[VxlanVtepStateDetached] = "DETACHED"
	VxlanVtepStateStrMap[VxlanVtepStateInterface] = "INTERFACE"
	VxlanVtepStateStrMap[VxlanVtepStateNextHopInfo] = "NEXT HOP INFO"
	VxlanVtepStateStrMap[VxlanVtepStateResolveNextHopMac] = "RESOLVE NEXT HOP INFO"
	VxlanVtepStateStrMap[VxlanVtepStateHwConfig] = "HW CONFIG"
	VxlanVtepStateStrMap[VxlanVtepStateStart] = "LISTENER"
}

// VxlanVtepEvent is used to transition VTEP FSM to various
// configuration steps
const (
	VxlanVtepEventBegin = iota + 1
	VxlanVtepEventUnconditionalFallThrough
	VxlanVtepEventSrcInterfaceResolved
	VxlanVtepEventNextHopInfoResolved
	VxlanVtepEventNextHopInfoMacResolved
	VxlanVtepEventHwConfigComplete
	VxlanVtepEventStartListener
	VxlanVtepEventDisable
	VxlanVtepEventDetached
	VxlanVtepEventEnable
	VxlanVtepEventRetryTimerExpired
)

// VxlanVtepStateEvent is the interface struct for the fsm
type VxlanVtepStateEvent struct {
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
	owner       string
	logEna      bool
	strStateMap map[fsm.State]string
	logger      func(string)
}

func (se *VxlanVtepStateEvent) LoggerSet(log func(string))                 { se.logger = log }
func (se *VxlanVtepStateEvent) EnableLogging(ena bool)                     { se.logEna = ena }
func (se *VxlanVtepStateEvent) IsLoggerEna() bool                          { return se.logEna }
func (se *VxlanVtepStateEvent) StateStrMapSet(strMap map[fsm.State]string) { se.strStateMap = strMap }
func (se *VxlanVtepStateEvent) PreviousState() fsm.State                   { return se.ps }
func (se *VxlanVtepStateEvent) CurrentState() fsm.State                    { return se.s }
func (se *VxlanVtepStateEvent) PreviousEvent() fsm.Event                   { return se.pe }
func (se *VxlanVtepStateEvent) CurrentEvent() fsm.Event                    { return se.e }
func (se *VxlanVtepStateEvent) SetEvent(es string, e fsm.Event) {
	se.esrc = es
	se.pe = se.e
	se.e = e
}
func (se *VxlanVtepStateEvent) SetState(s fsm.State) {
	se.ps = se.s
	se.s = s
	if se.IsLoggerEna() && se.ps != se.s {
		//if se.IsLoggerEna() {
		se.logger((strings.Join([]string{"Src", se.esrc, "OldState", se.strStateMap[se.ps], "Evt", strconv.Itoa(int(se.e)), "NewState", se.strStateMap[s]}, ":")))
	}
}

// MachineEvent is a generic struct to pass data to the FSM
// contains the event which is being passed.  Who sourced the
// event.  The data associated with the event.  And a response
// channel in case the caller needs to wait for completion of the
// event to occur
type MachineEvent struct {
	E            fsm.Event
	Src          string
	Data         interface{}
	ResponseChan chan string
}

// VxlanVtepMachine holds FSM and current State
// and event channels for State transitions
type VxlanVtepMachine struct {
	Machine *fsm.Machine

	// State transition log
	log chan string

	// Reference to StpPort
	vtep *VtepDbEntry

	// machine specific events
	VxlanVtepEvents chan MachineEvent
}

func SendResponse(msg string, responseChan chan string) {
	responseChan <- msg
}

func (m *VxlanVtepMachine) GetCurrStateStr() string {
	return VxlanVtepStateStrMap[m.Machine.Curr.CurrentState()]
}

func (m *VxlanVtepMachine) GetPrevStateStr() string {
	return VxlanVtepStateStrMap[m.Machine.Curr.PreviousState()]
}

// NewVxlanVtepFSMMachine will create a new instance of the VxlanVtepMachine
func NewVxlanVtepFSMMachine(vtep *VtepDbEntry) *VxlanVtepMachine {
	m := &VxlanVtepMachine{
		vtep:            vtep,
		VxlanVtepEvents: make(chan MachineEvent, 50)}

	vtep.VxlanVtepMachineFsm = m

	return m
}

// VxlanVtepFsmLogger is a wrapper above the vxlan logger to be passed into the FSM
func VxlanVtepFsmLogger(str string) {
	logger.Info(str)
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (m *VxlanVtepMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if m.Machine == nil {
		m.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	m.Machine.Rules = r
	m.Machine.Curr = &VxlanVtepStateEvent{
		strStateMap: VxlanVtepStateStrMap,
		logEna:      true,
		logger:      VxlanVtepFsmLogger,
		owner:       VxlanVtepMachineModuleStr,
		ps:          VxlanVtepStateNone,
		s:           VxlanVtepStateNone,
	}

	return m.Machine
}

func (vm *VxlanVtepMachine) BEGIN() {

	vtep := vm.vtep
	// restart the timer on deprovisioning as we will retry each of the
	// state transitions again
	vtep.retrytimer.Reset(retrytime)

	vm.VxlanVtepEvents <- MachineEvent{
		E:   VxlanVtepEventBegin,
		Src: VxlanVtepMachineModuleStr,
	}
}

// Stop should clean up all resources
func (vm *VxlanVtepMachine) Stop() {

	vtep := vm.vtep
	responsechan := make(chan string)
	vm.VxlanVtepEvents <- MachineEvent{
		E:            VxlanVtepEventDetached,
		Src:          VxlanVtepMachineModuleStr,
		ResponseChan: responsechan,
	}

	<-responsechan

	logger.Info("Close VTEP MACHINE")
	close(vm.VxlanVtepEvents)

	vtep.wg.Wait()

	if vtep.retrytimer != nil {
		vtep.retrytimer.Stop()
		vtep.retrytimer = nil
	}

}

// VxlanVtepInit is the state at which we start the VTEP FSM timer
// Then try and resolve the Src Interface info which will be used
// as the Src Info of the VXLAN header (IP/MAC)
func (vm *VxlanVtepMachine) VxlanVtepInit(m fsm.Machine, data interface{}) fsm.State {

	vtep := vm.vtep

	//logger.Info(fmt.Sprintln("vxlandb", GetVxlanDB()))
	//logger.Debug("Vtep enabled %t", vtep.Enable)
	if vxlan, ok := GetVxlanDB()[vtep.Vni]; ok {
		//logger.Debug("Vxlan enabled %t", vxlan.Enable)
		if vxlan.Enable {
			if vtep.Enable {
				// src interface was supplied
				// lets lookup the appropriate info
				if vtep.SrcIfName != "" {
					for _, client := range ClientIntf {
						client.GetIntfInfo(vtep.SrcIfName, vm.VxlanVtepEvents)
					}
				} else if vtep.SrcIp.String() != "0.0.0.0" &&
					vtep.SrcIp != nil {
					// TODO need to handle case where src ip and mac is supplied
					// in typical cases the src mac is the switch mac, src ip should
					// be applied to the vtep itself
				}
			} else {
				// lets move from Init to Detached state since vni not attached
				vm.VxlanVtepEvents <- MachineEvent{
					E:   VxlanVtepEventDisable,
					Src: VxlanVtepMachineModuleStr,
				}
			}
		} else {
			// lets move from Init to Detached state since vni not attached
			vm.VxlanVtepEvents <- MachineEvent{
				E:   VxlanVtepEventDetached,
				Src: VxlanVtepMachineModuleStr,
			}
		}
	} else {
		// lets move from Init to Detached state since vni not attached
		vm.VxlanVtepEvents <- MachineEvent{
			E:   VxlanVtepEventDetached,
			Src: VxlanVtepMachineModuleStr,
		}
	}

	return VxlanVtepStateInit
}

// VxlanVtepInterface is the state which the VTEP has not been attached to a proper
// VXLAN Vni domain
func (vm *VxlanVtepMachine) VxlanVtepDetached(m fsm.Machine, data interface{}) fsm.State {

	vtep := vm.vtep

	if vm.Machine.Curr.PreviousState() != VxlanVtepStateInit {
		DeProvisionVtep(vtep, true)
	}

	return VxlanVtepStateDetached
}

// VxlanVtepInterface is the state at which the source interface has been
// resolved, thus we save the info found, then we try and resolve the next hop info
//  for the destination MAC. This info should come from RIB
func (vm *VxlanVtepMachine) VxlanVtepInterface(m fsm.Machine, data interface{}) fsm.State {

	vtep := vm.vtep
	logger.Info(fmt.Sprintln("VxlanVtepInterface", data, reflect.TypeOf(data)))
	switch data.(type) {
	case VxlanIntfInfo:
		intfinfo := data.(VxlanIntfInfo)

		// save off info related to the source interface
		vtep.SrcIfName = intfinfo.IntfName
		vtep.SrcIp = intfinfo.Ip
		vtep.SrcMac = intfinfo.Mac
		vtep.SrcIfIndex = intfinfo.IfIndex
		logger.Info(fmt.Sprintf("%s: resolved srcip %s src mac %s from intf %s", strings.TrimRight(vtep.VtepName, "Int"), vtep.SrcIp, vtep.SrcMac, vtep.SrcIfName))
	}
	// lets resolve the next hop ip and intf
	for _, client := range ClientIntf {
		if client.GetNextHopInfo(vtep.DstIp, vm.VxlanVtepEvents) {
			// found now register for events
			client.RegisterReachability(vtep.DstIp)
		}
	}

	return VxlanVtepStateInterface
}

// VxlanVtepNextHopInfo is the state at which the next hop info is found and the data
// is stored on the vtep.  It will then try to resolve the mac address of the next hop ip
func (vm *VxlanVtepMachine) VxlanVtepNextHopInfo(m fsm.Machine, data interface{}) fsm.State {
	vtep := vm.vtep

	logger.Info(fmt.Sprintln("VxlanVtepNextHopInfo", data))
	switch data.(type) {
	case VtepNextHopInfo:
		info := data.(VtepNextHopInfo)
		vtep.NextHop.Ip = info.Ip
		vtep.NextHop.IfIndex = info.IfIndex
		vtep.NextHop.IfName = info.IfName

		// TODO need create a port listener per next hop interface
		// lets start listening on this port for VXLAN frames
		// Call will protect against multiple calls to same port
		//CreatePort(VxlanNextHopIp.Intf, vtep.UDP)
		logger.Info(fmt.Sprintf("%s: found next hop ip %s and interface %s for dstip %s ", strings.TrimRight(vtep.VtepName, "Int"), vtep.NextHop.Ip, vtep.NextHop.IfName, vtep.DstIp))
		// next state
	}
	// lets resolve the next hop mac
	for _, client := range ClientIntf {
		client.ResolveNextHopMac(vtep.NextHop.Ip, vtep.NextHop.IfName, vm.VxlanVtepEvents)
	}

	return VxlanVtepStateNextHopInfo

}

// VxlanVtepResolveNextHopInfoMac is the state at which the next hop ip mac address has
// been resolved via ARP.  The info is then stored against the vtep.  No action as this
// state will immediately tranisition to VxlanVtepStateHwConfig
func (vm *VxlanVtepMachine) VxlanVtepResolveNextHopInfoMac(m fsm.Machine, data interface{}) fsm.State {

	vtep := vm.vtep
	logger.Info(fmt.Sprintln("VxlanVtepResolveNextHopInfoMac", data))

	switch data.(type) {
	case net.HardwareAddr:
		mac := data.(net.HardwareAddr)

		vtep.DstMac = mac
		logger.Info(fmt.Sprintf("%s: resolved mac %s for next hop ip %s ", strings.TrimRight(vtep.VtepName, "Int"), vtep.DstMac, vtep.NextHop.Ip))
	}

	return VxlanVtepStateResolveNextHopMac
}

// VxlanVtepConfigHw is the state which all VTEP info is known and now we are ready
// to configure this info into the hw.  Stop the running timer as well
func (vm *VxlanVtepMachine) VxlanVtepConfigHw(m fsm.Machine, data interface{}) fsm.State {
	vtep := vm.vtep

	for _, client := range ClientIntf {
		client.CreateVtep(vtep, vm.VxlanVtepEvents)
	}

	vtep.retrytimer.Stop()

	return VxlanVtepStateHwConfig
}

// VxlanVtepStartListener is the state which the VTEP will start to listen on the linux
// interfaces for VXLAN packets as well as listen for packets which are transmitted out
// of the vtep interface
func (vm *VxlanVtepMachine) VxlanVtepStartListener(m fsm.Machine, data interface{}) fsm.State {

	vtep := vm.vtep

	hwdata := data.(VtepCreateCfgData)
	vtep.SetIfIndex(hwdata.IfIndex)
	logger.Info(fmt.Sprintf("%s: Starting listening for packets on vtep intf %s and intf %s ifindex %d", vtep.VtepName, vtep.VtepHandleName, vtep.NextHop.IfName, vtep.VtepIfIndex))
	VxlanVtepRxTx(vtep)
	VxlanCreatePortRxTx(vtep.NextHop.IfName, vtep.UDP)
	return VxlanVtepStateStart
}

//VxlanVtepDisabled is the state which holds says the vtep interface is disabled and
// should not process any packets.  It will be deprovisioned in the hw as a result
func (vm *VxlanVtepMachine) VxlanVtepDisabled(m fsm.Machine, data interface{}) fsm.State {

	vtep := vm.vtep
	// we only want to deprovision if we are coming from a state other than init
	if vm.Machine.Curr.PreviousState() != VxlanVtepStateInit {
		DeProvisionVtep(vtep, false)
	}
	return VxlanVtepStateDisabled
}

func (vm *VxlanVtepMachine) ProcessPostStateProcessing(data interface{}) {

	if vm.Machine.Curr.CurrentState() == VxlanVtepStateResolveNextHopMac {
		vm.VxlanVtepEvents <- MachineEvent{
			E:   VxlanVtepEventUnconditionalFallThrough,
			Src: VxlanVtepMachineModuleStr,
		}
	}
}

func VxlanVtepMachineFSMBuild(vtep *VtepDbEntry) *VxlanVtepMachine {

	rules := fsm.Ruleset{}

	// Instantiate a new NewVxlanVtepFSMMachine
	vm := NewVxlanVtepFSMMachine(vtep)

	// BEGIN -> INIT
	rules.AddRule(VxlanVtepStateNone, VxlanVtepEventBegin, vm.VxlanVtepInit)
	rules.AddRule(VxlanVtepStateDisabled, VxlanVtepEventBegin, vm.VxlanVtepInit)
	rules.AddRule(VxlanVtepStateInterface, VxlanVtepEventBegin, vm.VxlanVtepInit)
	rules.AddRule(VxlanVtepStateNextHopInfo, VxlanVtepEventBegin, vm.VxlanVtepInit)
	rules.AddRule(VxlanVtepStateResolveNextHopMac, VxlanVtepEventBegin, vm.VxlanVtepInit)
	rules.AddRule(VxlanVtepStateHwConfig, VxlanVtepEventBegin, vm.VxlanVtepInit)
	rules.AddRule(VxlanVtepStateStart, VxlanVtepEventBegin, vm.VxlanVtepInit)
	rules.AddRule(VxlanVtepStateDetached, VxlanVtepEventBegin, vm.VxlanVtepInit)
	rules.AddRule(VxlanVtepStateInit, VxlanVtepEventBegin, vm.VxlanVtepInit)

	// DETACHED -> Init -> DETACHED
	rules.AddRule(VxlanVtepStateInit, VxlanVtepEventDetached, vm.VxlanVtepDetached)
	rules.AddRule(VxlanVtepStateInterface, VxlanVtepEventDetached, vm.VxlanVtepDetached)
	rules.AddRule(VxlanVtepStateNextHopInfo, VxlanVtepEventDetached, vm.VxlanVtepDetached)
	rules.AddRule(VxlanVtepStateResolveNextHopMac, VxlanVtepEventDetached, vm.VxlanVtepDetached)
	rules.AddRule(VxlanVtepStateHwConfig, VxlanVtepEventDetached, vm.VxlanVtepDetached)
	rules.AddRule(VxlanVtepStateStart, VxlanVtepEventDetached, vm.VxlanVtepDetached)
	rules.AddRule(VxlanVtepStateInit, VxlanVtepEventDetached, vm.VxlanVtepDetached)

	// SRC INTERFACE RESOLVED -> INTERFACE
	rules.AddRule(VxlanVtepStateInit, VxlanVtepEventSrcInterfaceResolved, vm.VxlanVtepInterface)

	// NEXT HOP INFO RESOLVED -> NEXT HOP INFO
	rules.AddRule(VxlanVtepStateInterface, VxlanVtepEventNextHopInfoResolved, vm.VxlanVtepNextHopInfo)

	// NEXT HOP INFO MAC RESOLVED -> NEXT HOP INFO MAC
	rules.AddRule(VxlanVtepStateNextHopInfo, VxlanVtepEventNextHopInfoMacResolved, vm.VxlanVtepResolveNextHopInfoMac)

	// UNCONDITIONAL FALL THROUGH -> HW CONFIG
	rules.AddRule(VxlanVtepStateResolveNextHopMac, VxlanVtepEventUnconditionalFallThrough, vm.VxlanVtepConfigHw)

	// HW CONFIG COMPLETE -> START LISTENER
	rules.AddRule(VxlanVtepStateHwConfig, VxlanVtepEventHwConfigComplete, vm.VxlanVtepStartListener)

	// DISABLED -> DISABLED
	rules.AddRule(VxlanVtepStateInterface, VxlanVtepEventDisable, vm.VxlanVtepDisabled)
	rules.AddRule(VxlanVtepStateNextHopInfo, VxlanVtepEventDisable, vm.VxlanVtepDisabled)
	rules.AddRule(VxlanVtepStateResolveNextHopMac, VxlanVtepEventDisable, vm.VxlanVtepDisabled)
	rules.AddRule(VxlanVtepStateHwConfig, VxlanVtepEventDisable, vm.VxlanVtepDisabled)
	rules.AddRule(VxlanVtepStateStart, VxlanVtepEventDisable, vm.VxlanVtepDisabled)
	rules.AddRule(VxlanVtepStateDetached, VxlanVtepEventDisable, vm.VxlanVtepDisabled)
	rules.AddRule(VxlanVtepStateInit, VxlanVtepEventDisable, vm.VxlanVtepDisabled)

	// Certain clients will need to have information polled if an event is not generated
	for _, client := range ClientIntf {
		//logger.Info("Adding State to", client, client.IsClientIntfType(client, VXLANSnapClientStr))
		if client.IsClientIntfType(client, VXLANSnapClientStr) {
			// user has not configured src interface
			rules.AddRule(VxlanVtepStateInit, VxlanVtepEventRetryTimerExpired, vm.VxlanVtepInit)
			// arpd is not sending an event for the resolved mac thus must poll till it is resolved
			rules.AddRule(VxlanVtepStateNextHopInfo, VxlanVtepEventRetryTimerExpired, vm.VxlanVtepNextHopInfo)
		} else if !client.IsClientIntfType(client, VXLANSnapClientStr) {
			// in mock environmnet no asicd
			rules.AddRule(VxlanVtepStateInit, VxlanVtepEventRetryTimerExpired, vm.VxlanVtepInit)
			// in mock environment no arp
			rules.AddRule(VxlanVtepStateNextHopInfo, VxlanVtepEventRetryTimerExpired, vm.VxlanVtepNextHopInfo)
			// in mock environment no rib
			rules.AddRule(VxlanVtepStateInterface, VxlanVtepEventRetryTimerExpired, vm.VxlanVtepInterface)

		}
	}

	// Create a new FSM and apply the rules
	vm.Apply(&rules)

	return vm
}

// VxlanVtepMachineMain:
func (vtep *VtepDbEntry) VxlanVtepMachineMain() {

	// Build the State machine
	vm := VxlanVtepMachineFSMBuild(vtep)
	vtep.wg.Add(1)

	// set the inital State
	vm.Machine.Start(vm.Machine.Curr.PreviousState())

	// start a timer to retry any state configuration
	vtep.retrytimer = time.NewTimer(retrytime)

	// lets create a go routing which will wait for the specific events
	// that the Port Timer State Machine should handle
	go func() {
		logger.Info(fmt.Sprintln("Vtep MACHINE Start", vtep.VtepName))
		defer vtep.wg.Done()
		for {
			select {
			case _, ok := <-vtep.retrytimer.C:

				//logger.Info("Timer Expired")
				if ok {

					// in the case that the interface call needs to be polled then add state
					vm.Machine.ProcessEvent(VxlanVtepMachineModuleStr, VxlanVtepEventRetryTimerExpired, nil)
					vtep.ticksTillConfig++
					vtep.retrytimer.Reset(retrytime)
				}

			case event, ok := <-vm.VxlanVtepEvents:
				//fmt.Println("VXLAN event", event)
				if ok {
					rv := vm.Machine.ProcessEvent(event.Src, event.E, event.Data)
					if rv != nil {
						logger.Err(fmt.Sprintf("%s: %s event[%d] currState[%s]\n", vtep.VtepName, rv, event.E, VxlanVtepStateStrMap[vm.Machine.Curr.CurrentState()]))
					} else {
						// POST events
						vm.ProcessPostStateProcessing(event.Data)
					}
					if event.ResponseChan != nil {
						SendResponse(VxlanVtepMachineModuleStr, event.ResponseChan)
					}
				} else {
					//logger.Info(fmt.Sprintln("Vtep MACHINE Stop", vtep.VtepName))
					return
				}
			}
		}
	}()
}
