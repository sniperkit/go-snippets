package n2xsdk

type IgmpQuerierPool struct {
 Handler string
}

func(np *IgmpQuerierPool) SetVersion () error {
 //parameters: DeviceHandle IgmpVersion
 //AgtIgmpQuerierPool SetVersion
 return nil
}

func(np *IgmpQuerierPool) GetVersion ()(string, error) {
 //parameters: DeviceHandle
 //AgtIgmpQuerierPool GetVersion
 return "", nil
}

func(np *IgmpQuerierPool) SetParameter () error {
 //parameters: DeviceHandle Parameter Value
 //AgtIgmpQuerierPool SetParameter
 return nil
}

func(np *IgmpQuerierPool) GetParameter ()(string, error) {
 //parameters: DeviceHandle Parameter
 //AgtIgmpQuerierPool GetParameter
 return "", nil
}

func(np *IgmpQuerierPool) GetInfo ()(string, error) {
 //parameters: DeviceHandle Info
 //AgtIgmpQuerierPool GetInfo
 return "", nil
}

func(np *IgmpQuerierPool) SetGeneralQueryInterval () error {
 //parameters: DeviceHandle QueryInterval
 //AgtIgmpQuerierPool SetGeneralQueryInterval
 return nil
}

func(np *IgmpQuerierPool) GetGeneralQueryInterval ()(string, error) {
 //parameters: DeviceHandle
 //AgtIgmpQuerierPool GetGeneralQueryInterval
 return "", nil
}

func(np *IgmpQuerierPool) AddGroupQueryMessage () error {
 //parameters: DeviceHandle GroupPoolHandles
 //AgtIgmpQuerierPool AddGroupQueryMessage
 return nil
}

func(np *IgmpQuerierPool) RemoveGroupQueryMessage () error {
 //parameters: DeviceHandle GroupPoolHandle
 //AgtIgmpQuerierPool RemoveGroupQueryMessage
 return nil
}

func(np *IgmpQuerierPool) RemoveAllQueryMessages () error {
 //parameters: DeviceHandle
 //AgtIgmpQuerierPool RemoveAllQueryMessages
 return nil
}

func(np *IgmpQuerierPool) ListGroupQueryMessages ()(string, error) {
 //parameters: DeviceHandle
 //AgtIgmpQuerierPool ListGroupQueryMessages
 return "", nil
}

func(np *IgmpQuerierPool) AddSourcePools () error {
 //parameters: DeviceHandle GroupPoolHandle SourcePoolHandles
 //AgtIgmpQuerierPool AddSourcePools
 return nil
}

func(np *IgmpQuerierPool) RemoveSourcePools () error {
 //parameters: DeviceHandle GroupPoolHandle SourcePoolHandles
 //AgtIgmpQuerierPool RemoveSourcePools
 return nil
}

func(np *IgmpQuerierPool) RemoveAllSourcePools () error {
 //parameters: DeviceHandle GroupPoolHandle
 //AgtIgmpQuerierPool RemoveAllSourcePools
 return nil
}

func(np *IgmpQuerierPool) ListSourcePools ()(string, error) {
 //parameters: DeviceHandle GroupPoolHandle
 //AgtIgmpQuerierPool ListSourcePools
 return "", nil
}

func(np *IgmpQuerierPool) GetNumberOfSourcesInGroupPool ()(string, error) {
 //parameters: DeviceHandle GroupPoolHandle
 //AgtIgmpQuerierPool GetNumberOfSourcesInGroupPool
 return "", nil
}

func(np *IgmpQuerierPool) SetGroupQueryParameters () error {
 //parameters: DeviceHandle GroupPoolHandle Lmqi Lmqc
 //AgtIgmpQuerierPool SetGroupQueryParameters
 return nil
}

func(np *IgmpQuerierPool) GetGroupQueryParameters ()(string, error) {
 //parameters: DeviceHandle GroupPoolHandle
 //AgtIgmpQuerierPool GetGroupQueryParameters
 return "", nil
}

func(np *IgmpQuerierPool) StartGeneralQueryMessage () error {
 //parameters: DeviceHandle
 //AgtIgmpQuerierPool StartGeneralQueryMessage
 return nil
}

func(np *IgmpQuerierPool) StopGeneralQueryMessage () error {
 //parameters: DeviceHandle
 //AgtIgmpQuerierPool StopGeneralQueryMessage
 return nil
}

func(np *IgmpQuerierPool) IsGeneralQueryRunning () error {
 //parameters: DeviceHandle
 //AgtIgmpQuerierPool IsGeneralQueryRunning
 return nil
}

func(np *IgmpQuerierPool) StartGroupQueryMessages () error {
 //parameters: DeviceHandle GroupPoolHandle
 //AgtIgmpQuerierPool StartGroupQueryMessages
 return nil
}

func(np *IgmpQuerierPool) StartAllQueryMessages () error {
 //parameters: DeviceHandle
 //AgtIgmpQuerierPool StartAllQueryMessages
 return nil
}

func(np *IgmpQuerierPool) StopGroupQueryMessages () error {
 //parameters: DeviceHandle GroupPoolHandle
 //AgtIgmpQuerierPool StopGroupQueryMessages
 return nil
}

func(np *IgmpQuerierPool) StopAllQueryMessages () error {
 //parameters: DeviceHandle
 //AgtIgmpQuerierPool StopAllQueryMessages
 return nil
}

func(np *IgmpQuerierPool) ListGroupQueryMessagesRunning ()(string, error) {
 //parameters: DeviceHandle
 //AgtIgmpQuerierPool ListGroupQueryMessagesRunning
 return "", nil
}

func(np *IgmpQuerierPool) IsGroupQueryRunning () error {
 //parameters: DeviceHandle GroupPoolHandle
 //AgtIgmpQuerierPool IsGroupQueryRunning
 return nil
}

func(np *IgmpQuerierPool) IsIpv4PriorityNoCodePointFieldSelected () error {
 //parameters: DeviceHandle
 //AgtIgmpQuerierPool IsIpv4PriorityNoCodePointFieldSelected
 return nil
}

func(np *IgmpQuerierPool) SelectIpv4PriorityNoCodePointField () error {
 //parameters: DeviceHandle
 //AgtIgmpQuerierPool SelectIpv4PriorityNoCodePointField
 return nil
}

func(np *IgmpQuerierPool) GetIpv4Priority ()(string, error) {
 //parameters: DeviceHandle
 //AgtIgmpQuerierPool GetIpv4Priority
 return "", nil
}

func(np *IgmpQuerierPool) SetIpv4Priority () error {
 //parameters: DeviceHandle Ipv4Priority
 //AgtIgmpQuerierPool SetIpv4Priority
 return nil
}

func(np *IgmpQuerierPool) IsIpv4PriorityTypeOfServiceFieldSelected () error {
 //parameters: DeviceHandle
 //AgtIgmpQuerierPool IsIpv4PriorityTypeOfServiceFieldSelected
 return nil
}

func(np *IgmpQuerierPool) SelectIpv4PriorityTypeOfServiceField () error {
 //parameters: DeviceHandle
 //AgtIgmpQuerierPool SelectIpv4PriorityTypeOfServiceField
 return nil
}

func(np *IgmpQuerierPool) GetIpv4PriorityTypeOfServiceValue ()(string, error) {
 //parameters: DeviceHandle TosCodePointField
 //AgtIgmpQuerierPool GetIpv4PriorityTypeOfServiceValue
 return "", nil
}

func(np *IgmpQuerierPool) SetIpv4PriorityTypeOfServiceValue () error {
 //parameters: DeviceHandle TosCodePointField Value
 //AgtIgmpQuerierPool SetIpv4PriorityTypeOfServiceValue
 return nil
}

func(np *IgmpQuerierPool) IsIpv4PriorityDiffServFieldSelected () error {
 //parameters: DeviceHandle DiffServPerHopBehaviorGroup
 //AgtIgmpQuerierPool IsIpv4PriorityDiffServFieldSelected
 return nil
}

func(np *IgmpQuerierPool) SelectIpv4PriorityDiffServField () error {
 //parameters: DeviceHandle DiffServPerHopBehaviorGroup
 //AgtIgmpQuerierPool SelectIpv4PriorityDiffServField
 return nil
}

func(np *IgmpQuerierPool) GetIpv4PriorityDiffServValue ()(string, error) {
 //parameters: DeviceHandle DiffServCodePointField
 //AgtIgmpQuerierPool GetIpv4PriorityDiffServValue
 return "", nil
}

func(np *IgmpQuerierPool) SetIpv4PriorityDiffServValue () error {
 //parameters: DeviceHandle DiffServCodePointConfigurableField Value
 //AgtIgmpQuerierPool SetIpv4PriorityDiffServValue
 return nil
}

func(np *IgmpQuerierPool) GetVlanPriority ()(string, error) {
 //parameters: DeviceHandle VlanTagIndex
 //AgtIgmpQuerierPool GetVlanPriority
 return "", nil
}

func(np *IgmpQuerierPool) SetVlanPriority () error {
 //parameters: DeviceHandle VlanTagIndex VlanPriority
 //AgtIgmpQuerierPool SetVlanPriority
 return nil
}

