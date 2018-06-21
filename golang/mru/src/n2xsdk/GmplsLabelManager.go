package n2xsdk

type GmplsLabelManager struct {
 Handler string
}

func(np *GmplsLabelManager) GetLabelSessionType ()(string, error) {
 //parameters: PortHandle
 //AgtGmplsLabelManager GetLabelSessionType
 return "", nil
}

func(np *GmplsLabelManager) SetSwitchingCapability () error {
 //parameters: PortHandle SwitchingCapability
 //AgtGmplsLabelManager SetSwitchingCapability
 return nil
}

func(np *GmplsLabelManager) GetSwitchingCapability ()(string, error) {
 //parameters: PortHandle
 //AgtGmplsLabelManager GetSwitchingCapability
 return "", nil
}

func(np *GmplsLabelManager) SetInterfaceNumbering () error {
 //parameters: PortHandle InterfaceNumbering
 //AgtGmplsLabelManager SetInterfaceNumbering
 return nil
}

func(np *GmplsLabelManager) GetInterfaceNumbering ()(string, error) {
 //parameters: PortHandle
 //AgtGmplsLabelManager GetInterfaceNumbering
 return "", nil
}

func(np *GmplsLabelManager) AddUnnumberedInterfaceGroup () error {
 //parameters: PortHandle InterfaceIdStart InterfaceIdIncr NumberOfInterfaces
 //AgtGmplsLabelManager AddUnnumberedInterfaceGroup
 return nil
}

func(np *GmplsLabelManager) AddNumberedInterfaceGroup () error {
 //parameters: PortHandle InterfaceIdStart InterfaceIdIncr NumberOfInterfaces RemoteIdStart RemoteIdIncr
 //AgtGmplsLabelManager AddNumberedInterfaceGroup
 return nil
}

func(np *GmplsLabelManager) SetUnnumberedInterfaceGroup () error {
 //parameters: PortHandle GroupHandle InterfaceIdStart InterfaceIdIncr NumberOfInterfaces
 //AgtGmplsLabelManager SetUnnumberedInterfaceGroup
 return nil
}

func(np *GmplsLabelManager) SetNumberedInterfaceGroup () error {
 //parameters: PortHandle GroupHandle InterfaceIdStart InterfaceIdIncr NumberOfInterfaces RemoteIdStart RemoteIdIncr
 //AgtGmplsLabelManager SetNumberedInterfaceGroup
 return nil
}

func(np *GmplsLabelManager) GetUnnumberedInterfaceGroup ()(string, error) {
 //parameters: PortHandle GroupHandle
 //AgtGmplsLabelManager GetUnnumberedInterfaceGroup
 return "", nil
}

func(np *GmplsLabelManager) GetNumberedInterfaceGroup ()(string, error) {
 //parameters: PortHandle GroupHandle
 //AgtGmplsLabelManager GetNumberedInterfaceGroup
 return "", nil
}

func(np *GmplsLabelManager) SetInterfaceGroupTeLink () error {
 //parameters: PortHandle GroupHandle UseTeLink TeLinkAddr
 //AgtGmplsLabelManager SetInterfaceGroupTeLink
 return nil
}

func(np *GmplsLabelManager) GetInterfaceGroupTeLink ()(string, error) {
 //parameters: PortHandle GroupHandle
 //AgtGmplsLabelManager GetInterfaceGroupTeLink
 return "", nil
}

func(np *GmplsLabelManager) SetUnnumberedInterfaceGroupRemoteId () error {
 //parameters: PortHandle GroupHandle RemoteIdStart RemoteIdIncr
 //AgtGmplsLabelManager SetUnnumberedInterfaceGroupRemoteId
 return nil
}

func(np *GmplsLabelManager) GetUnnumberedInterfaceGroupRemoteId ()(string, error) {
 //parameters: PortHandle GroupHandle
 //AgtGmplsLabelManager GetUnnumberedInterfaceGroupRemoteId
 return "", nil
}

func(np *GmplsLabelManager) EnableChannelization () error {
 //parameters: PortHandle GroupHandle
 //AgtGmplsLabelManager EnableChannelization
 return nil
}

func(np *GmplsLabelManager) DisableChannelization () error {
 //parameters: PortHandle GroupHandle
 //AgtGmplsLabelManager DisableChannelization
 return nil
}

func(np *GmplsLabelManager) IsChannelizationEnabled () error {
 //parameters: PortHandle GroupHandle
 //AgtGmplsLabelManager IsChannelizationEnabled
 return nil
}

func(np *GmplsLabelManager) SetTdmInterfaceGroupRate () error {
 //parameters: PortHandle GroupHandle InterfaceRate
 //AgtGmplsLabelManager SetTdmInterfaceGroupRate
 return nil
}

func(np *GmplsLabelManager) GetTdmInterfaceGroupRate ()(string, error) {
 //parameters: PortHandle GroupHandle
 //AgtGmplsLabelManager GetTdmInterfaceGroupRate
 return "", nil
}

func(np *GmplsLabelManager) SetLscInterfaceGroupRate () error {
 //parameters: PortHandle GroupHandle LambdaIdStart LambdaIdIncr NumberOfLambdas InterfaceRate
 //AgtGmplsLabelManager SetLscInterfaceGroupRate
 return nil
}

func(np *GmplsLabelManager) GetLscInterfaceGroupRate ()(string, error) {
 //parameters: PortHandle GroupHandle
 //AgtGmplsLabelManager GetLscInterfaceGroupRate
 return "", nil
}

func(np *GmplsLabelManager) SetFscInterfaceGroupRate () error {
 //parameters: PortHandle GroupHandle FscInterfaceStart FscInterfaceIncr InterfaceRate
 //AgtGmplsLabelManager SetFscInterfaceGroupRate
 return nil
}

func(np *GmplsLabelManager) GetFscInterfaceGroupRate ()(string, error) {
 //parameters: PortHandle GroupHandle
 //AgtGmplsLabelManager GetFscInterfaceGroupRate
 return "", nil
}

func(np *GmplsLabelManager) RemoveInterfaceGroup () error {
 //parameters: PortHandle GroupHandle
 //AgtGmplsLabelManager RemoveInterfaceGroup
 return nil
}

func(np *GmplsLabelManager) ListInterfaceGroups ()(string, error) {
 //parameters: PortHandle
 //AgtGmplsLabelManager ListInterfaceGroups
 return "", nil
}

