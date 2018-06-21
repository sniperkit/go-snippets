package n2xsdk

type PortSelector struct {
 Handler string
}

func(np *PortSelector) ListModules ()(string, error) {
 //parameters: 
 //AgtPortSelector ListModules
 return "", nil
}

func(np *PortSelector) GetLastModule ()(string, error) {
 //parameters: 
 //AgtPortSelector GetLastModule
 return "", nil
}

func(np *PortSelector) GetModuleDescription ()(string, error) {
 //parameters: ModuleNumber
 //AgtPortSelector GetModuleDescription
 return "", nil
}

func(np *PortSelector) IsChassisBlade () error {
 //parameters: ModuleNumber
 //AgtPortSelector IsChassisBlade
 return nil
}

func(np *PortSelector) GetModuleExtendedType ()(string, error) {
 //parameters: ModuleNumber
 //AgtPortSelector GetModuleExtendedType
 return "", nil
}

func(np *PortSelector) GetChassisNumber ()(string, error) {
 //parameters: ModuleNumber
 //AgtPortSelector GetChassisNumber
 return "", nil
}

func(np *PortSelector) GetChassisSlotNumber ()(string, error) {
 //parameters: ModuleNumber
 //AgtPortSelector GetChassisSlotNumber
 return "", nil
}

func(np *PortSelector) GetPortType ()(string, error) {
 //parameters: ModuleNumber portNumber
 //AgtPortSelector GetPortType
 return "", nil
}

func(np *PortSelector) GetPortLabel ()(string, error) {
 //parameters: ModuleNumber portNumber
 //AgtPortSelector GetPortLabel
 return "", nil
}

func(np *PortSelector) GetModuleClass ()(string, error) {
 //parameters: ModuleType
 //AgtPortSelector GetModuleClass
 return "", nil
}

func(np *PortSelector) ListModuleTypes ()(string, error) {
 //parameters: ModuleNumber
 //AgtPortSelector ListModuleTypes
 return "", nil
}

func(np *PortSelector) GetModuleType ()(string, error) {
 //parameters: ModuleNumber
 //AgtPortSelector GetModuleType
 return "", nil
}

func(np *PortSelector) SetModuleType () error {
 //parameters: ModuleNumber ModuleType
 //AgtPortSelector SetModuleType
 return nil
}

func(np *PortSelector) GetModuleLimit ()(string, error) {
 //parameters: 
 //AgtPortSelector GetModuleLimit
 return "", nil
}

func(np *PortSelector) IsModuleSupported () error {
 //parameters: ModuleNumber
 //AgtPortSelector IsModuleSupported
 return nil
}

func(np *PortSelector) ListSessionModuleTypes ()(string, error) {
 //parameters: ModuleNumber
 //AgtPortSelector ListSessionModuleTypes
 return "", nil
}

func(np *PortSelector) GetModuleState ()(string, error) {
 //parameters: ModuleNumber
 //AgtPortSelector GetModuleState
 return "", nil
}

func(np *PortSelector) GetModuleLock ()(string, error) {
 //parameters: ModuleNumber
 //AgtPortSelector GetModuleLock
 return "", nil
}

func(np *PortSelector) GetChassisUpstreamLock ()(string, error) {
 //parameters: ChassisNumber
 //AgtPortSelector GetChassisUpstreamLock
 return "", nil
}

func(np *PortSelector) GetChassisDownstreamLock ()(string, error) {
 //parameters: ChassisNumber
 //AgtPortSelector GetChassisDownstreamLock
 return "", nil
}

func(np *PortSelector) GetLockedModuleList ()(string, error) {
 //parameters: 
 //AgtPortSelector GetLockedModuleList
 return "", nil
}

func(np *PortSelector) GetLockedModules ()(string, error) {
 //parameters: 
 //AgtPortSelector GetLockedModules
 return "", nil
}

func(np *PortSelector) ListRequiredModules ()(string, error) {
 //parameters: SelectedCount SelectedModules
 //AgtPortSelector ListRequiredModules
 return "", nil
}

func(np *PortSelector) ListUnavailableModules ()(string, error) {
 //parameters: SelectedCount SelectedModules
 //AgtPortSelector ListUnavailableModules
 return "", nil
}

func(np *PortSelector) ListUnavailableModulesInList ()(string, error) {
 //parameters: SelectedCount SelectedModules
 //AgtPortSelector ListUnavailableModulesInList
 return "", nil
}

func(np *PortSelector) ListLockedModules ()(string, error) {
 //parameters: 
 //AgtPortSelector ListLockedModules
 return "", nil
}

func(np *PortSelector) ListUnavailableModulesAfterSelection ()(string, error) {
 //parameters: InputCount ProposedSelectionModules
 //AgtPortSelector ListUnavailableModulesAfterSelection
 return "", nil
}

func(np *PortSelector) ListAvailableModules ()(string, error) {
 //parameters: 
 //AgtPortSelector ListAvailableModules
 return "", nil
}

func(np *PortSelector) GetAvailableModuleSets ()(string, error) {
 //parameters: 
 //AgtPortSelector GetAvailableModuleSets
 return "", nil
}

func(np *PortSelector) ListAvailableModuleSet ()(string, error) {
 //parameters: SetNumber
 //AgtPortSelector ListAvailableModuleSet
 return "", nil
}

func(np *PortSelector) AddPort () error {
 //parameters: ModuleNumber portNumber
 //AgtPortSelector AddPort
 return nil
}

func(np *PortSelector) AddPorts () error {
 //parameters: LabelCount PortLabels
 //AgtPortSelector AddPorts
 return nil
}

func(np *PortSelector) AddPortsWithLock () error {
 //parameters: LabelCount PortLabels ModuleCount ModuleNumbers
 //AgtPortSelector AddPortsWithLock
 return nil
}

func(np *PortSelector) RemovePort () error {
 //parameters: PortHandle
 //AgtPortSelector RemovePort
 return nil
}

func(np *PortSelector) RemovePorts () error {
 //parameters: PortCount PortHandles
 //AgtPortSelector RemovePorts
 return nil
}

func(np *PortSelector) ListPorts ()(string, error) {
 //parameters: 
 //AgtPortSelector ListPorts
 return "", nil
}

func(np *PortSelector) FindPortHandle () error {
 //parameters: ModuleNumber portNumber
 //AgtPortSelector FindPortHandle
 return nil
}

func(np *PortSelector) FindPortHandleFromLabel () error {
 //parameters: portLabel
 //AgtPortSelector FindPortHandleFromLabel
 return nil
}

func(np *PortSelector) GetPortDetails ()(string, error) {
 //parameters: PortHandle
 //AgtPortSelector GetPortDetails
 return "", nil
}

func(np *PortSelector) GetPortLabelFromHandle ()(string, error) {
 //parameters: PortHandle
 //AgtPortSelector GetPortLabelFromHandle
 return "", nil
}

func(np *PortSelector) IsDummyPort () error {
 //parameters: PortHandle
 //AgtPortSelector IsDummyPort
 return nil
}

func(np *PortSelector) SetPortComment () error {
 //parameters: PortHandle PortComment
 //AgtPortSelector SetPortComment
 return nil
}

func(np *PortSelector) GetPortComment ()(string, error) {
 //parameters: PortHandle
 //AgtPortSelector GetPortComment
 return "", nil
}

func(np *PortSelector) AddGroup () error {
 //parameters: Count PortHandles
 //AgtPortSelector AddGroup
 return nil
}

func(np *PortSelector) RemoveGroup () error {
 //parameters: GroupHandle
 //AgtPortSelector RemoveGroup
 return nil
}

func(np *PortSelector) ListGroups ()(string, error) {
 //parameters: 
 //AgtPortSelector ListGroups
 return "", nil
}

func(np *PortSelector) ListPortsInGroup ()(string, error) {
 //parameters: GroupHandle
 //AgtPortSelector ListPortsInGroup
 return "", nil
}

func(np *PortSelector) AddModule () error {
 //parameters: ModuleType
 //AgtPortSelector AddModule
 return nil
}

func(np *PortSelector) RemoveModule () error {
 //parameters: ModuleNumber
 //AgtPortSelector RemoveModule
 return nil
}

func(np *PortSelector) ListPortsInModule ()(string, error) {
 //parameters: ModuleNumber
 //AgtPortSelector ListPortsInModule
 return "", nil
}

func(np *PortSelector) GetPortClass ()(string, error) {
 //parameters: PortType
 //AgtPortSelector GetPortClass
 return "", nil
}

func(np *PortSelector) IsPerPortPersonalitySupported () error {
 //parameters: ModuleNumber
 //AgtPortSelector IsPerPortPersonalitySupported
 return nil
}

func(np *PortSelector) ListPortPersonalities ()(string, error) {
 //parameters: ModuleNumber portNumber
 //AgtPortSelector ListPortPersonalities
 return "", nil
}

func(np *PortSelector) GetPortPersonality ()(string, error) {
 //parameters: ModuleNumber portNumber
 //AgtPortSelector GetPortPersonality
 return "", nil
}

func(np *PortSelector) SetPortPersonality () error {
 //parameters: ModuleNumber portNumber Personality
 //AgtPortSelector SetPortPersonality
 return nil
}

func(np *PortSelector) ListModulePersonalities ()(string, error) {
 //parameters: ModuleNumber
 //AgtPortSelector ListModulePersonalities
 return "", nil
}

func(np *PortSelector) GetModulePersonality ()(string, error) {
 //parameters: ModuleNumber
 //AgtPortSelector GetModulePersonality
 return "", nil
}

func(np *PortSelector) SetModulePersonality () error {
 //parameters: ModuleNumber Personality
 //AgtPortSelector SetModulePersonality
 return nil
}

func(np *PortSelector) GetPortState ()(string, error) {
 //parameters: ModuleNumber portNumber
 //AgtPortSelector GetPortState
 return "", nil
}

func(np *PortSelector) GetPortLock ()(string, error) {
 //parameters: ModuleNumber portNumber
 //AgtPortSelector GetPortLock
 return "", nil
}

func(np *PortSelector) ListLockedPorts ()(string, error) {
 //parameters: 
 //AgtPortSelector ListLockedPorts
 return "", nil
}

func(np *PortSelector) ListAvailablePorts ()(string, error) {
 //parameters: 
 //AgtPortSelector ListAvailablePorts
 return "", nil
}

func(np *PortSelector) ListAvailablePortsOnModule ()(string, error) {
 //parameters: ModuleNumber
 //AgtPortSelector ListAvailablePortsOnModule
 return "", nil
}

func(np *PortSelector) IsMultiUserChassisSystem () error {
 //parameters: 
 //AgtPortSelector IsMultiUserChassisSystem
 return nil
}

func(np *PortSelector) SetNextPortHandle () error {
 //parameters: NextAvailablePort
 //AgtPortSelector SetNextPortHandle
 return nil
}

func(np *PortSelector) GetNextPortHandle ()(string, error) {
 //parameters: 
 //AgtPortSelector GetNextPortHandle
 return "", nil
}

func(np *PortSelector) GetMasterPort ()(string, error) {
 //parameters: ModuleNumber portNumber
 //AgtPortSelector GetMasterPort
 return "", nil
}

func(np *PortSelector) ListSlavePorts ()(string, error) {
 //parameters: ModuleNumber MasterPortNumber
 //AgtPortSelector ListSlavePorts
 return "", nil
}

func(np *PortSelector) ListMasterPorts ()(string, error) {
 //parameters: ModuleNumber
 //AgtPortSelector ListMasterPorts
 return "", nil
}

func(np *PortSelector) SetPersonality () error {
 //parameters: ModuleNumber MasterPortNumber Personality
 //AgtPortSelector SetPersonality
 return nil
}

func(np *PortSelector) GetPersonality ()(string, error) {
 //parameters: ModuleNumber MasterPortNumber
 //AgtPortSelector GetPersonality
 return "", nil
}

func(np *PortSelector) ListPersonalities ()(string, error) {
 //parameters: ModuleNumber portNumber
 //AgtPortSelector ListPersonalities
 return "", nil
}

