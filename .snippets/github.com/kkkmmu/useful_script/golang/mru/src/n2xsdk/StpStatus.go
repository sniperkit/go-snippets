package n2xsdk

type StpStatus struct {
 Handler string
}

func(np *StpStatus) GetOperatingMode ()(string, error) {
 //parameters: SessionPoolHandle SessionInstanceHandle
 //AgtStpStatus GetOperatingMode
 return "", nil
}

func(np *StpStatus) GetPortState ()(string, error) {
 //parameters: SessionPoolHandle SessionInstanceHandle
 //AgtStpStatus GetPortState
 return "", nil
}

func(np *StpStatus) GetPortRole ()(string, error) {
 //parameters: SessionPoolHandle SessionInstanceHandle
 //AgtStpStatus GetPortRole
 return "", nil
}

func(np *StpStatus) GetAutomaticPathCost ()(string, error) {
 //parameters: SessionPoolHandle SessionInstanceHandle
 //AgtStpStatus GetAutomaticPathCost
 return "", nil
}

func(np *StpStatus) GetRootIdentifier ()(string, error) {
 //parameters: SessionPoolHandle SessionInstanceHandle
 //AgtStpStatus GetRootIdentifier
 return "", nil
}

func(np *StpStatus) GetRootPathCost ()(string, error) {
 //parameters: SessionPoolHandle SessionInstanceHandle
 //AgtStpStatus GetRootPathCost
 return "", nil
}

func(np *StpStatus) IsEdge () error {
 //parameters: SessionPoolHandle SessionInstanceHandle
 //AgtStpStatus IsEdge
 return nil
}

func(np *StpStatus) GetTime ()(string, error) {
 //parameters: SessionPoolHandle SessionInstanceHandle TimeParameter
 //AgtStpStatus GetTime
 return "", nil
}

func(np *StpStatus) GetNeighborDetails ()(string, error) {
 //parameters: SessionPoolHandle SessionInstanceHandle NeighborParameter
 //AgtStpStatus GetNeighborDetails
 return "", nil
}

func(np *StpStatus) GetSourcePortVlanId ()(string, error) {
 //parameters: SessionPoolHandle SessionInstanceHandle
 //AgtStpStatus GetSourcePortVlanId
 return "", nil
}

