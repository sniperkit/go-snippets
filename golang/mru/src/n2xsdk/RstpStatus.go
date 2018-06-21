package n2xsdk

type RstpStatus struct {
 Handler string
}

func(np *RstpStatus) GetNeighborStatus ()(string, error) {
 //parameters: SessionPoolHandle SessionInstance NeighborStatusParameter
 //AgtRstpStatus GetNeighborStatus
 return "", nil
}

func(np *RstpStatus) GetOperatingMode ()(string, error) {
 //parameters: SessionPoolHandle SessionInstanceHandle
 //AgtRstpStatus GetOperatingMode
 return "", nil
}

func(np *RstpStatus) GetPortState ()(string, error) {
 //parameters: SessionPoolHandle SessionInstanceHandle
 //AgtRstpStatus GetPortState
 return "", nil
}

func(np *RstpStatus) GetPortRole ()(string, error) {
 //parameters: SessionPoolHandle SessionInstanceHandle
 //AgtRstpStatus GetPortRole
 return "", nil
}

func(np *RstpStatus) GetAutomaticPathCost ()(string, error) {
 //parameters: SessionPoolHandle SessionInstanceHandle
 //AgtRstpStatus GetAutomaticPathCost
 return "", nil
}

func(np *RstpStatus) GetRootIdentifier ()(string, error) {
 //parameters: SessionPoolHandle SessionInstanceHandle
 //AgtRstpStatus GetRootIdentifier
 return "", nil
}

func(np *RstpStatus) GetRootPathCost ()(string, error) {
 //parameters: SessionPoolHandle SessionInstanceHandle
 //AgtRstpStatus GetRootPathCost
 return "", nil
}

func(np *RstpStatus) IsEdge () error {
 //parameters: SessionPoolHandle SessionInstanceHandle
 //AgtRstpStatus IsEdge
 return nil
}

func(np *RstpStatus) GetTime ()(string, error) {
 //parameters: SessionPoolHandle SessionInstanceHandle TimeParameter
 //AgtRstpStatus GetTime
 return "", nil
}

func(np *RstpStatus) GetNeighborDetails ()(string, error) {
 //parameters: SessionPoolHandle SessionInstanceHandle NeighborParameter
 //AgtRstpStatus GetNeighborDetails
 return "", nil
}

func(np *RstpStatus) GetSourcePortVlanId ()(string, error) {
 //parameters: SessionPoolHandle SessionInstanceHandle
 //AgtRstpStatus GetSourcePortVlanId
 return "", nil
}

