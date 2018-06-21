package n2xsdk

type MldQuerierPool struct {
 Handler string
}

func(np *MldQuerierPool) SetVersion () error {
 //parameters: DeviceHandle MldVersion
 //AgtMldQuerierPool SetVersion
 return nil
}

func(np *MldQuerierPool) GetVersion ()(string, error) {
 //parameters: DeviceHandle
 //AgtMldQuerierPool GetVersion
 return "", nil
}

func(np *MldQuerierPool) SetParameter () error {
 //parameters: DeviceHandle Parameter Value
 //AgtMldQuerierPool SetParameter
 return nil
}

func(np *MldQuerierPool) GetParameter ()(string, error) {
 //parameters: DeviceHandle Parameter
 //AgtMldQuerierPool GetParameter
 return "", nil
}

func(np *MldQuerierPool) GetInfo ()(string, error) {
 //parameters: DeviceHandle Info
 //AgtMldQuerierPool GetInfo
 return "", nil
}

func(np *MldQuerierPool) SetGeneralQueryInterval () error {
 //parameters: DeviceHandle QueryInterval
 //AgtMldQuerierPool SetGeneralQueryInterval
 return nil
}

func(np *MldQuerierPool) GetGeneralQueryInterval ()(string, error) {
 //parameters: DeviceHandle
 //AgtMldQuerierPool GetGeneralQueryInterval
 return "", nil
}

func(np *MldQuerierPool) AddGroupQueryMessage () error {
 //parameters: DeviceHandle GroupPoolHandles
 //AgtMldQuerierPool AddGroupQueryMessage
 return nil
}

func(np *MldQuerierPool) RemoveGroupQueryMessage () error {
 //parameters: DeviceHandle GroupPoolHandle
 //AgtMldQuerierPool RemoveGroupQueryMessage
 return nil
}

func(np *MldQuerierPool) RemoveAllQueryMessages () error {
 //parameters: DeviceHandle
 //AgtMldQuerierPool RemoveAllQueryMessages
 return nil
}

func(np *MldQuerierPool) ListGroupQueryMessages ()(string, error) {
 //parameters: DeviceHandle
 //AgtMldQuerierPool ListGroupQueryMessages
 return "", nil
}

func(np *MldQuerierPool) AddSourcePools () error {
 //parameters: DeviceHandle GroupPoolHandle SourcePoolHandles
 //AgtMldQuerierPool AddSourcePools
 return nil
}

func(np *MldQuerierPool) RemoveSourcePools () error {
 //parameters: DeviceHandle GroupPoolHandle SourcePoolHandles
 //AgtMldQuerierPool RemoveSourcePools
 return nil
}

func(np *MldQuerierPool) RemoveAllSourcePools () error {
 //parameters: DeviceHandle GroupPoolHandle
 //AgtMldQuerierPool RemoveAllSourcePools
 return nil
}

func(np *MldQuerierPool) ListSourcePools ()(string, error) {
 //parameters: DeviceHandle GroupPoolHandle
 //AgtMldQuerierPool ListSourcePools
 return "", nil
}

func(np *MldQuerierPool) GetNumberOfSourcesInGroupPool ()(string, error) {
 //parameters: DeviceHandle GroupPoolHandle
 //AgtMldQuerierPool GetNumberOfSourcesInGroupPool
 return "", nil
}

func(np *MldQuerierPool) SetGroupQueryParameters () error {
 //parameters: DeviceHandle GroupPoolHandle Lmqi Lmqc
 //AgtMldQuerierPool SetGroupQueryParameters
 return nil
}

func(np *MldQuerierPool) GetGroupQueryParameters ()(string, error) {
 //parameters: DeviceHandle GroupPoolHandle
 //AgtMldQuerierPool GetGroupQueryParameters
 return "", nil
}

func(np *MldQuerierPool) StartGeneralQueryMessage () error {
 //parameters: DeviceHandle
 //AgtMldQuerierPool StartGeneralQueryMessage
 return nil
}

func(np *MldQuerierPool) StopGeneralQueryMessage () error {
 //parameters: DeviceHandle
 //AgtMldQuerierPool StopGeneralQueryMessage
 return nil
}

func(np *MldQuerierPool) IsGeneralQueryRunning () error {
 //parameters: DeviceHandle
 //AgtMldQuerierPool IsGeneralQueryRunning
 return nil
}

func(np *MldQuerierPool) StartGroupQueryMessages () error {
 //parameters: DeviceHandle GroupPoolHandle
 //AgtMldQuerierPool StartGroupQueryMessages
 return nil
}

func(np *MldQuerierPool) StartAllQueryMessages () error {
 //parameters: DeviceHandle
 //AgtMldQuerierPool StartAllQueryMessages
 return nil
}

func(np *MldQuerierPool) StopGroupQueryMessages () error {
 //parameters: DeviceHandle GroupPoolHandle
 //AgtMldQuerierPool StopGroupQueryMessages
 return nil
}

func(np *MldQuerierPool) StopAllQueryMessages () error {
 //parameters: DeviceHandle
 //AgtMldQuerierPool StopAllQueryMessages
 return nil
}

func(np *MldQuerierPool) ListGroupQueryMessagesRunning ()(string, error) {
 //parameters: DeviceHandle
 //AgtMldQuerierPool ListGroupQueryMessagesRunning
 return "", nil
}

func(np *MldQuerierPool) IsGroupQueryRunning () error {
 //parameters: DeviceHandle GroupPoolHandle
 //AgtMldQuerierPool IsGroupQueryRunning
 return nil
}

func(np *MldQuerierPool) GetVlanPriority ()(string, error) {
 //parameters: DeviceHandle VlanTagIndex
 //AgtMldQuerierPool GetVlanPriority
 return "", nil
}

func(np *MldQuerierPool) SetVlanPriority () error {
 //parameters: DeviceHandle VlanTagIndex VlanPriority
 //AgtMldQuerierPool SetVlanPriority
 return nil
}

func(np *MldQuerierPool) IsIpv6PriorityNoCodePointFieldSelected () error {
 //parameters: DeviceHandle
 //AgtMldQuerierPool IsIpv6PriorityNoCodePointFieldSelected
 return nil
}

func(np *MldQuerierPool) SelectIpv6PriorityNoCodePointField () error {
 //parameters: DeviceHandle
 //AgtMldQuerierPool SelectIpv6PriorityNoCodePointField
 return nil
}

func(np *MldQuerierPool) GetIpv6Priority ()(string, error) {
 //parameters: DeviceHandle
 //AgtMldQuerierPool GetIpv6Priority
 return "", nil
}

func(np *MldQuerierPool) SetIpv6Priority () error {
 //parameters: DeviceHandle Ipv6Priority
 //AgtMldQuerierPool SetIpv6Priority
 return nil
}

func(np *MldQuerierPool) IsIpv6PriorityDiffServFieldSelected () error {
 //parameters: DeviceHandle DiffServPerHopBehaviorGroup
 //AgtMldQuerierPool IsIpv6PriorityDiffServFieldSelected
 return nil
}

func(np *MldQuerierPool) SelectIpv6PriorityDiffServField () error {
 //parameters: DeviceHandle DiffServPerHopBehaviorGroup
 //AgtMldQuerierPool SelectIpv6PriorityDiffServField
 return nil
}

func(np *MldQuerierPool) GetIpv6PriorityDiffServValue ()(string, error) {
 //parameters: DeviceHandle DiffServCodePointField
 //AgtMldQuerierPool GetIpv6PriorityDiffServValue
 return "", nil
}

func(np *MldQuerierPool) SetIpv6PriorityDiffServValue () error {
 //parameters: DeviceHandle DiffServCodePointConfigurableField Value
 //AgtMldQuerierPool SetIpv6PriorityDiffServValue
 return nil
}

