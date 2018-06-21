package n2xsdk

type StpSession struct {
 Handler string
}

func(np *StpSession) OpenActive () error {
 //parameters: SessionInstanceHandle
 //AgtStpSession OpenActive
 return nil
}

func(np *StpSession) OpenPassive () error {
 //parameters: SessionInstanceHandle
 //AgtStpSession OpenPassive
 return nil
}

func(np *StpSession) Close () error {
 //parameters: SessionInstanceHandle
 //AgtStpSession Close
 return nil
}

func(np *StpSession) Reset () error {
 //parameters: SessionInstanceHandle
 //AgtStpSession Reset
 return nil
}

func(np *StpSession) ResetMeasurements () error {
 //parameters: SessionInstanceHandle
 //AgtStpSession ResetMeasurements
 return nil
}

func(np *StpSession) SignalTopologyChange () error {
 //parameters: SessionInstanceHandle
 //AgtStpSession SignalTopologyChange
 return nil
}

func(np *StpSession) GetIdentifierEncoding ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtStpSession GetIdentifierEncoding
 return "", nil
}

func(np *StpSession) SetIdentifierEncoding () error {
 //parameters: SessionPoolHandle IdentifierEncoding
 //AgtStpSession SetIdentifierEncoding
 return nil
}

func(np *StpSession) GetBridgeIdentifier ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtStpSession GetBridgeIdentifier
 return "", nil
}

func(np *StpSession) SetBridgeIdentifier () error {
 //parameters: SessionPoolHandle Priority MacAddress
 //AgtStpSession SetBridgeIdentifier
 return nil
}

func(np *StpSession) GetPortIdentifier ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtStpSession GetPortIdentifier
 return "", nil
}

func(np *StpSession) SetPortIdentifier () error {
 //parameters: SessionPoolHandle Priority Number
 //AgtStpSession SetPortIdentifier
 return nil
}

func(np *StpSession) GetPathCostMode ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtStpSession GetPathCostMode
 return "", nil
}

func(np *StpSession) SetPathCostMode () error {
 //parameters: SessionPoolHandle PathCostMode
 //AgtStpSession SetPathCostMode
 return nil
}

func(np *StpSession) GetPathCostValue ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtStpSession GetPathCostValue
 return "", nil
}

func(np *StpSession) SetPathCostValue () error {
 //parameters: SessionPoolHandle PathCost
 //AgtStpSession SetPathCostValue
 return nil
}

func(np *StpSession) GetRootIdentifier ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtStpSession GetRootIdentifier
 return "", nil
}

func(np *StpSession) SetRootIdentifier () error {
 //parameters: SessionPoolHandle Priority MacAddress
 //AgtStpSession SetRootIdentifier
 return nil
}

func(np *StpSession) GetRootPathCost ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtStpSession GetRootPathCost
 return "", nil
}

func(np *StpSession) SetRootPathCost () error {
 //parameters: SessionPoolHandle RootPathCost
 //AgtStpSession SetRootPathCost
 return nil
}

func(np *StpSession) GetTime ()(string, error) {
 //parameters: SessionPoolHandle TimeParameter
 //AgtStpSession GetTime
 return "", nil
}

func(np *StpSession) SetTime () error {
 //parameters: SessionPoolHandle TimeParameter Time
 //AgtStpSession SetTime
 return nil
}

func(np *StpSession) GetOperationalModifier ()(string, error) {
 //parameters: SessionPoolHandle OperationalParameter
 //AgtStpSession GetOperationalModifier
 return "", nil
}

func(np *StpSession) SetOperationalModifier () error {
 //parameters: SessionPoolHandle OperationalParameter Value
 //AgtStpSession SetOperationalModifier
 return nil
}

func(np *StpSession) SetIncrementingParameter () error {
 //parameters: SessionPoolHandle IncrementingParameter Value Repeat Increment
 //AgtStpSession SetIncrementingParameter
 return nil
}

func(np *StpSession) GetIncrementingParameter ()(string, error) {
 //parameters: SessionPoolHandle IncrementingParameter
 //AgtStpSession GetIncrementingParameter
 return "", nil
}

func(np *StpSession) GetBridgeGroupAddressType ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtStpSession GetBridgeGroupAddressType
 return "", nil
}

func(np *StpSession) SetBridgeGroupAddressType () error {
 //parameters: SessionPoolHandle BridgeGroupAddressType
 //AgtStpSession SetBridgeGroupAddressType
 return nil
}

func(np *StpSession) EnablePvstPlus () error {
 //parameters: SessionPoolHandle
 //AgtStpSession EnablePvstPlus
 return nil
}

func(np *StpSession) DisablePvstPlus () error {
 //parameters: SessionPoolHandle
 //AgtStpSession DisablePvstPlus
 return nil
}

func(np *StpSession) IsPvstPlusEnabled () error {
 //parameters: SessionPoolHandle
 //AgtStpSession IsPvstPlusEnabled
 return nil
}

