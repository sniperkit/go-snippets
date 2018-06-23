package n2xsdk

type MstpSession struct {
 Handler string
}

func(np *MstpSession) SignalTopologyChange () error {
 //parameters: SessionInstanceHandle SpanningTreeIdentifier
 //AgtMstpSession SignalTopologyChange
 return nil
}

func(np *MstpSession) GetCistRegionalRootIdentifier ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtMstpSession GetCistRegionalRootIdentifier
 return "", nil
}

func(np *MstpSession) SetCistRegionalRootIdentifier () error {
 //parameters: SessionPoolHandle Priority MacAddress
 //AgtMstpSession SetCistRegionalRootIdentifier
 return nil
}

func(np *MstpSession) GetCistRootPathCostExternal ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtMstpSession GetCistRootPathCostExternal
 return "", nil
}

func(np *MstpSession) SetCistRootPathCostExternal () error {
 //parameters: SessionPoolHandle RootPathCost
 //AgtMstpSession SetCistRootPathCostExternal
 return nil
}

func(np *MstpSession) GetCistRootPathCostInternal ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtMstpSession GetCistRootPathCostInternal
 return "", nil
}

func(np *MstpSession) SetCistRootPathCostInternal () error {
 //parameters: SessionPoolHandle RootPathCost
 //AgtMstpSession SetCistRootPathCostInternal
 return nil
}

func(np *MstpSession) GetMstConfig ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtMstpSession GetMstConfig
 return "", nil
}

func(np *MstpSession) SetMstConfig () error {
 //parameters: SessionPoolHandle MstConfigName MstConfigRevision
 //AgtMstpSession SetMstConfig
 return nil
}

func(np *MstpSession) GetCistMaxHops ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtMstpSession GetCistMaxHops
 return "", nil
}

func(np *MstpSession) SetCistMaxHops () error {
 //parameters: SessionPoolHandle MaxHops
 //AgtMstpSession SetCistMaxHops
 return nil
}

func(np *MstpSession) SetCistIncrementingParameter () error {
 //parameters: SessionPoolHandle IncrementingCistParameter Value Repeat Increment
 //AgtMstpSession SetCistIncrementingParameter
 return nil
}

func(np *MstpSession) GetCistIncrementingParameter ()(string, error) {
 //parameters: SessionPoolHandle IncrementingCistParameter
 //AgtMstpSession GetCistIncrementingParameter
 return "", nil
}

func(np *MstpSession) GetMstiRegionalRootIdentifier ()(string, error) {
 //parameters: SessionPoolHandle MstiIdentifier
 //AgtMstpSession GetMstiRegionalRootIdentifier
 return "", nil
}

func(np *MstpSession) SetMstiRegionalRootIdentifier () error {
 //parameters: SessionPoolHandle MstiIdentifier Priority MacAddress
 //AgtMstpSession SetMstiRegionalRootIdentifier
 return nil
}

func(np *MstpSession) GetMstiRootPathCostInternal ()(string, error) {
 //parameters: SessionPoolHandle MstiIdentifier
 //AgtMstpSession GetMstiRootPathCostInternal
 return "", nil
}

func(np *MstpSession) SetMstiRootPathCostInternal () error {
 //parameters: SessionPoolHandle MstiIdentifier RootPathCost
 //AgtMstpSession SetMstiRootPathCostInternal
 return nil
}

func(np *MstpSession) GetMstiBridgePriority ()(string, error) {
 //parameters: SessionPoolHandle MstiIdentifier
 //AgtMstpSession GetMstiBridgePriority
 return "", nil
}

func(np *MstpSession) SetMstiBridgePriority () error {
 //parameters: SessionPoolHandle MstiIdentifier Priority
 //AgtMstpSession SetMstiBridgePriority
 return nil
}

func(np *MstpSession) GetMstiPortPriority ()(string, error) {
 //parameters: SessionPoolHandle MstiIdentifier
 //AgtMstpSession GetMstiPortPriority
 return "", nil
}

func(np *MstpSession) SetMstiPortPriority () error {
 //parameters: SessionPoolHandle MstiIdentifier Priority
 //AgtMstpSession SetMstiPortPriority
 return nil
}

func(np *MstpSession) GetMstiPathCost ()(string, error) {
 //parameters: SessionPoolHandle MstiIdentifier
 //AgtMstpSession GetMstiPathCost
 return "", nil
}

func(np *MstpSession) SetMstiPathCost () error {
 //parameters: SessionPoolHandle MstiIdentifier MstiPathCost
 //AgtMstpSession SetMstiPathCost
 return nil
}

func(np *MstpSession) GetMstiMaxHops ()(string, error) {
 //parameters: SessionPoolHandle MstiIdentifier
 //AgtMstpSession GetMstiMaxHops
 return "", nil
}

func(np *MstpSession) SetMstiMaxHops () error {
 //parameters: SessionPoolHandle MstiIdentifier MaxHops
 //AgtMstpSession SetMstiMaxHops
 return nil
}

func(np *MstpSession) SetMstiIncrementingParameter () error {
 //parameters: SessionPoolHandle MstiIdentifier IncrementingMstiParameter Value Repeat Increment
 //AgtMstpSession SetMstiIncrementingParameter
 return nil
}

func(np *MstpSession) GetMstiIncrementingParameter ()(string, error) {
 //parameters: SessionPoolHandle MstiIdentifier IncrementingMstiParameter
 //AgtMstpSession GetMstiIncrementingParameter
 return "", nil
}

func(np *MstpSession) AddMsti () error {
 //parameters: SessionPoolHandle MstiIdentifier
 //AgtMstpSession AddMsti
 return nil
}

func(np *MstpSession) RemoveMsti () error {
 //parameters: SessionPoolHandle MstiIdentifier
 //AgtMstpSession RemoveMsti
 return nil
}

func(np *MstpSession) ListMsti ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtMstpSession ListMsti
 return "", nil
}

func(np *MstpSession) RemoveAllMstis () error {
 //parameters: SessionPoolHandle
 //AgtMstpSession RemoveAllMstis
 return nil
}

func(np *MstpSession) MapVlanToMsti () error {
 //parameters: SessionPoolHandle MstiIdentifier VlanIdStart Repeat Increment Count
 //AgtMstpSession MapVlanToMsti
 return nil
}

func(np *MstpSession) ListMstiVlanMapEntryHandles ()(string, error) {
 //parameters: SessionPoolHandle MstiIdentifier
 //AgtMstpSession ListMstiVlanMapEntryHandles
 return "", nil
}

func(np *MstpSession) GetMstiVlanMapEntry ()(string, error) {
 //parameters: SessionPoolHandle MstiIdentifier VlanMapEntryHandle
 //AgtMstpSession GetMstiVlanMapEntry
 return "", nil
}

func(np *MstpSession) RemoveMstiVlanMapEntry () error {
 //parameters: SessionPoolHandle MstiIdentifier VlanMapEntryHandle
 //AgtMstpSession RemoveMstiVlanMapEntry
 return nil
}

func(np *MstpSession) ClearVlanMap () error {
 //parameters: SessionPoolHandle MstiIdentifier
 //AgtMstpSession ClearVlanMap
 return nil
}

func(np *MstpSession) ForceMigrationCheck () error {
 //parameters: SessionInstanceHandle
 //AgtMstpSession ForceMigrationCheck
 return nil
}

func(np *MstpSession) OpenActive () error {
 //parameters: SessionInstanceHandle
 //AgtMstpSession OpenActive
 return nil
}

func(np *MstpSession) OpenPassive () error {
 //parameters: SessionInstanceHandle
 //AgtMstpSession OpenPassive
 return nil
}

func(np *MstpSession) Close () error {
 //parameters: SessionInstanceHandle
 //AgtMstpSession Close
 return nil
}

func(np *MstpSession) Reset () error {
 //parameters: SessionInstanceHandle
 //AgtMstpSession Reset
 return nil
}

func(np *MstpSession) ResetMeasurements () error {
 //parameters: SessionInstanceHandle
 //AgtMstpSession ResetMeasurements
 return nil
}

func(np *MstpSession) GetIdentifierEncoding ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtMstpSession GetIdentifierEncoding
 return "", nil
}

func(np *MstpSession) SetIdentifierEncoding () error {
 //parameters: SessionPoolHandle IdentifierEncoding
 //AgtMstpSession SetIdentifierEncoding
 return nil
}

func(np *MstpSession) GetBridgeIdentifier ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtMstpSession GetBridgeIdentifier
 return "", nil
}

func(np *MstpSession) SetBridgeIdentifier () error {
 //parameters: SessionPoolHandle Priority MacAddress
 //AgtMstpSession SetBridgeIdentifier
 return nil
}

func(np *MstpSession) GetPortIdentifier ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtMstpSession GetPortIdentifier
 return "", nil
}

func(np *MstpSession) SetPortIdentifier () error {
 //parameters: SessionPoolHandle Priority Number
 //AgtMstpSession SetPortIdentifier
 return nil
}

func(np *MstpSession) GetPathCostMode ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtMstpSession GetPathCostMode
 return "", nil
}

func(np *MstpSession) SetPathCostMode () error {
 //parameters: SessionPoolHandle PathCostMode
 //AgtMstpSession SetPathCostMode
 return nil
}

func(np *MstpSession) GetPathCostValue ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtMstpSession GetPathCostValue
 return "", nil
}

func(np *MstpSession) SetPathCostValue () error {
 //parameters: SessionPoolHandle PathCost
 //AgtMstpSession SetPathCostValue
 return nil
}

func(np *MstpSession) GetRootIdentifier ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtMstpSession GetRootIdentifier
 return "", nil
}

func(np *MstpSession) SetRootIdentifier () error {
 //parameters: SessionPoolHandle Priority MacAddress
 //AgtMstpSession SetRootIdentifier
 return nil
}

func(np *MstpSession) GetRootPathCost ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtMstpSession GetRootPathCost
 return "", nil
}

func(np *MstpSession) SetRootPathCost () error {
 //parameters: SessionPoolHandle RootPathCost
 //AgtMstpSession SetRootPathCost
 return nil
}

func(np *MstpSession) GetTime ()(string, error) {
 //parameters: SessionPoolHandle TimeParameter
 //AgtMstpSession GetTime
 return "", nil
}

func(np *MstpSession) SetTime () error {
 //parameters: SessionPoolHandle TimeParameter Time
 //AgtMstpSession SetTime
 return nil
}

func(np *MstpSession) GetOperationalModifier ()(string, error) {
 //parameters: SessionPoolHandle OperationalParameter
 //AgtMstpSession GetOperationalModifier
 return "", nil
}

func(np *MstpSession) SetOperationalModifier () error {
 //parameters: SessionPoolHandle OperationalParameter Value
 //AgtMstpSession SetOperationalModifier
 return nil
}

func(np *MstpSession) SetIncrementingParameter () error {
 //parameters: SessionPoolHandle IncrementingParameter Value Repeat Increment
 //AgtMstpSession SetIncrementingParameter
 return nil
}

func(np *MstpSession) GetIncrementingParameter ()(string, error) {
 //parameters: SessionPoolHandle IncrementingParameter
 //AgtMstpSession GetIncrementingParameter
 return "", nil
}

func(np *MstpSession) GetBridgeGroupAddressType ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtMstpSession GetBridgeGroupAddressType
 return "", nil
}

func(np *MstpSession) SetBridgeGroupAddressType () error {
 //parameters: SessionPoolHandle BridgeGroupAddressType
 //AgtMstpSession SetBridgeGroupAddressType
 return nil
}

func(np *MstpSession) EnablePvstPlus () error {
 //parameters: SessionPoolHandle
 //AgtMstpSession EnablePvstPlus
 return nil
}

func(np *MstpSession) DisablePvstPlus () error {
 //parameters: SessionPoolHandle
 //AgtMstpSession DisablePvstPlus
 return nil
}

func(np *MstpSession) IsPvstPlusEnabled () error {
 //parameters: SessionPoolHandle
 //AgtMstpSession IsPvstPlusEnabled
 return nil
}

