package n2xsdk

type MulticastMesh struct {
 Handler string
}

func(np *MulticastMesh) GetType ()(string, error) {
 //parameters: Handle
 //AgtMulticastMesh GetType
 return "", nil
}

func(np *MulticastMesh) GetName ()(string, error) {
 //parameters: Handle
 //AgtMulticastMesh GetName
 return "", nil
}

func(np *MulticastMesh) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtMulticastMesh GetLockCount
 return "", nil
}

func(np *MulticastMesh) Enable () error {
 //parameters: MeshHandle
 //AgtMulticastMesh Enable
 return nil
}

func(np *MulticastMesh) Disable () error {
 //parameters: MeshHandle
 //AgtMulticastMesh Disable
 return nil
}

func(np *MulticastMesh) IsEnabled () error {
 //parameters: MeshHandle
 //AgtMulticastMesh IsEnabled
 return nil
}

func(np *MulticastMesh) SetTrafficOrientation () error {
 //parameters: MeshHandle Orientation
 //AgtMulticastMesh SetTrafficOrientation
 return nil
}

func(np *MulticastMesh) GetTrafficOrientation ()(string, error) {
 //parameters: MeshHandle
 //AgtMulticastMesh GetTrafficOrientation
 return "", nil
}

func(np *MulticastMesh) SetLengthMode () error {
 //parameters: MeshHandle LengthMode
 //AgtMulticastMesh SetLengthMode
 return nil
}

func(np *MulticastMesh) GetLengthMode ()(string, error) {
 //parameters: MeshHandle
 //AgtMulticastMesh GetLengthMode
 return "", nil
}

func(np *MulticastMesh) SetLength () error {
 //parameters: MeshHandle LengthType Count psaLengthParameterList
 //AgtMulticastMesh SetLength
 return nil
}

func(np *MulticastMesh) GetLength ()(string, error) {
 //parameters: MeshHandle
 //AgtMulticastMesh GetLength
 return "", nil
}

func(np *MulticastMesh) ListSources ()(string, error) {
 //parameters: MeshHandle
 //AgtMulticastMesh ListSources
 return "", nil
}

func(np *MulticastMesh) ListDestinations ()(string, error) {
 //parameters: MeshHandle
 //AgtMulticastMesh ListDestinations
 return "", nil
}

func(np *MulticastMesh) ListStreamGroups ()(string, error) {
 //parameters: MeshHandle
 //AgtMulticastMesh ListStreamGroups
 return "", nil
}

func(np *MulticastMesh) ListStreamGroupsBySource ()(string, error) {
 //parameters: MeshHandle Source
 //AgtMulticastMesh ListStreamGroupsBySource
 return "", nil
}

func(np *MulticastMesh) ListStreamGroupProfiles ()(string, error) {
 //parameters: MeshHandle StreamGroupCount psaStreamGroups
 //AgtMulticastMesh ListStreamGroupProfiles
 return "", nil
}

func(np *MulticastMesh) EnableStreamGeneration () error {
 //parameters: MeshHandle
 //AgtMulticastMesh EnableStreamGeneration
 return nil
}

func(np *MulticastMesh) DisableStreamGeneration () error {
 //parameters: MeshHandle
 //AgtMulticastMesh DisableStreamGeneration
 return nil
}

func(np *MulticastMesh) IsStreamGenerationEnabled () error {
 //parameters: MeshHandle
 //AgtMulticastMesh IsStreamGenerationEnabled
 return nil
}

func(np *MulticastMesh) SetAverageLoad () error {
 //parameters: MeshHandle AverageLoad Units
 //AgtMulticastMesh SetAverageLoad
 return nil
}

func(np *MulticastMesh) GetAverageLoad ()(string, error) {
 //parameters: MeshHandle
 //AgtMulticastMesh GetAverageLoad
 return "", nil
}

func(np *MulticastMesh) IsMeshAddedOnNewProfiles () error {
 //parameters: MeshHandle
 //AgtMulticastMesh IsMeshAddedOnNewProfiles
 return nil
}

func(np *MulticastMesh) SetSourceEndpointType () error {
 //parameters: MeshHandle SourceType
 //AgtMulticastMesh SetSourceEndpointType
 return nil
}

func(np *MulticastMesh) GetSourceEndpointType ()(string, error) {
 //parameters: MeshHandle
 //AgtMulticastMesh GetSourceEndpointType
 return "", nil
}

func(np *MulticastMesh) SetDestinationEndpointType () error {
 //parameters: MeshHandle DestinationType
 //AgtMulticastMesh SetDestinationEndpointType
 return nil
}

func(np *MulticastMesh) GetDestinationEndpointType ()(string, error) {
 //parameters: MeshHandle
 //AgtMulticastMesh GetDestinationEndpointType
 return "", nil
}

func(np *MulticastMesh) AddSources () error {
 //parameters: MeshHandle Count psaSources
 //AgtMulticastMesh AddSources
 return nil
}

func(np *MulticastMesh) RemoveSources () error {
 //parameters: MeshHandle Count psaSources
 //AgtMulticastMesh RemoveSources
 return nil
}

func(np *MulticastMesh) AddDestinations () error {
 //parameters: MeshHandle Count psaDestinations
 //AgtMulticastMesh AddDestinations
 return nil
}

func(np *MulticastMesh) RemoveDestinations () error {
 //parameters: MeshHandle Count psaDestinations
 //AgtMulticastMesh RemoveDestinations
 return nil
}

func(np *MulticastMesh) UpdateStreamGroups () error {
 //parameters: MeshHandle
 //AgtMulticastMesh UpdateStreamGroups
 return nil
}

func(np *MulticastMesh) GetStreamGroupsByPort ()(string, error) {
 //parameters: MeshHandle Source Destination
 //AgtMulticastMesh GetStreamGroupsByPort
 return "", nil
}

func(np *MulticastMesh) GetStreamGroupsByPool ()(string, error) {
 //parameters: MeshHandle SourcePool GroupPool
 //AgtMulticastMesh GetStreamGroupsByPool
 return "", nil
}

func(np *MulticastMesh) GetPortsForStreamGroup ()(string, error) {
 //parameters: MeshHandle StreamGroup
 //AgtMulticastMesh GetPortsForStreamGroup
 return "", nil
}

func(np *MulticastMesh) GetMulticastAddressPoolsForStreamGroup ()(string, error) {
 //parameters: MeshHandle StreamGroup
 //AgtMulticastMesh GetMulticastAddressPoolsForStreamGroup
 return "", nil
}

func(np *MulticastMesh) GetTotalSGPairsInUse ()(string, error) {
 //parameters: MeshHandle Source
 //AgtMulticastMesh GetTotalSGPairsInUse
 return "", nil
}

