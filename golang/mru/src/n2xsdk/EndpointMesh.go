package n2xsdk

type EndpointMesh struct {
 Handler string
}

func(np *EndpointMesh) GetType ()(string, error) {
 //parameters: Handle
 //AgtEndpointMesh GetType
 return "", nil
}

func(np *EndpointMesh) GetName ()(string, error) {
 //parameters: Handle
 //AgtEndpointMesh GetName
 return "", nil
}

func(np *EndpointMesh) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtEndpointMesh GetLockCount
 return "", nil
}

func(np *EndpointMesh) Enable () error {
 //parameters: MeshHandle
 //AgtEndpointMesh Enable
 return nil
}

func(np *EndpointMesh) Disable () error {
 //parameters: MeshHandle
 //AgtEndpointMesh Disable
 return nil
}

func(np *EndpointMesh) IsEnabled () error {
 //parameters: MeshHandle
 //AgtEndpointMesh IsEnabled
 return nil
}

func(np *EndpointMesh) SetTrafficOrientation () error {
 //parameters: MeshHandle Orientation
 //AgtEndpointMesh SetTrafficOrientation
 return nil
}

func(np *EndpointMesh) GetTrafficOrientation ()(string, error) {
 //parameters: MeshHandle
 //AgtEndpointMesh GetTrafficOrientation
 return "", nil
}

func(np *EndpointMesh) SetLengthMode () error {
 //parameters: MeshHandle LengthMode
 //AgtEndpointMesh SetLengthMode
 return nil
}

func(np *EndpointMesh) GetLengthMode ()(string, error) {
 //parameters: MeshHandle
 //AgtEndpointMesh GetLengthMode
 return "", nil
}

func(np *EndpointMesh) SetLength () error {
 //parameters: MeshHandle LengthType Count psaLengthParameterList
 //AgtEndpointMesh SetLength
 return nil
}

func(np *EndpointMesh) GetLength ()(string, error) {
 //parameters: MeshHandle
 //AgtEndpointMesh GetLength
 return "", nil
}

func(np *EndpointMesh) ListSources ()(string, error) {
 //parameters: MeshHandle
 //AgtEndpointMesh ListSources
 return "", nil
}

func(np *EndpointMesh) ListDestinations ()(string, error) {
 //parameters: MeshHandle
 //AgtEndpointMesh ListDestinations
 return "", nil
}

func(np *EndpointMesh) ListStreamGroups ()(string, error) {
 //parameters: MeshHandle
 //AgtEndpointMesh ListStreamGroups
 return "", nil
}

func(np *EndpointMesh) ListStreamGroupsBySource ()(string, error) {
 //parameters: MeshHandle Source
 //AgtEndpointMesh ListStreamGroupsBySource
 return "", nil
}

func(np *EndpointMesh) ListStreamGroupProfiles ()(string, error) {
 //parameters: MeshHandle StreamGroupCount psaStreamGroups
 //AgtEndpointMesh ListStreamGroupProfiles
 return "", nil
}

func(np *EndpointMesh) EnableStreamGeneration () error {
 //parameters: MeshHandle
 //AgtEndpointMesh EnableStreamGeneration
 return nil
}

func(np *EndpointMesh) DisableStreamGeneration () error {
 //parameters: MeshHandle
 //AgtEndpointMesh DisableStreamGeneration
 return nil
}

func(np *EndpointMesh) IsStreamGenerationEnabled () error {
 //parameters: MeshHandle
 //AgtEndpointMesh IsStreamGenerationEnabled
 return nil
}

func(np *EndpointMesh) SetAverageLoad () error {
 //parameters: MeshHandle AverageLoad Units
 //AgtEndpointMesh SetAverageLoad
 return nil
}

func(np *EndpointMesh) GetAverageLoad ()(string, error) {
 //parameters: MeshHandle
 //AgtEndpointMesh GetAverageLoad
 return "", nil
}

func(np *EndpointMesh) IsMeshAddedOnNewProfiles () error {
 //parameters: MeshHandle
 //AgtEndpointMesh IsMeshAddedOnNewProfiles
 return nil
}

func(np *EndpointMesh) AddSources () error {
 //parameters: MeshHandle NumSources psaSources
 //AgtEndpointMesh AddSources
 return nil
}

func(np *EndpointMesh) RemoveSources () error {
 //parameters: MeshHandle NumSources psaSources
 //AgtEndpointMesh RemoveSources
 return nil
}

func(np *EndpointMesh) AddDestinations () error {
 //parameters: MeshHandle NumDestinations psaDestinations
 //AgtEndpointMesh AddDestinations
 return nil
}

func(np *EndpointMesh) RemoveDestinations () error {
 //parameters: MeshHandle NumDestinations psaDestinations
 //AgtEndpointMesh RemoveDestinations
 return nil
}

func(np *EndpointMesh) SetEndpointTypes () error {
 //parameters: MeshHandle SourceType DestinationType
 //AgtEndpointMesh SetEndpointTypes
 return nil
}

func(np *EndpointMesh) GetEndpointTypes ()(string, error) {
 //parameters: MeshHandle
 //AgtEndpointMesh GetEndpointTypes
 return "", nil
}

func(np *EndpointMesh) SetTrafficDistribution () error {
 //parameters: MeshHandle Distribution
 //AgtEndpointMesh SetTrafficDistribution
 return nil
}

func(np *EndpointMesh) GetTrafficDistribution ()(string, error) {
 //parameters: MeshHandle
 //AgtEndpointMesh GetTrafficDistribution
 return "", nil
}

func(np *EndpointMesh) SetDeviceUnderTestType () error {
 //parameters: MeshHandle DutType
 //AgtEndpointMesh SetDeviceUnderTestType
 return nil
}

func(np *EndpointMesh) GetDeviceUnderTestType ()(string, error) {
 //parameters: MeshHandle
 //AgtEndpointMesh GetDeviceUnderTestType
 return "", nil
}

func(np *EndpointMesh) EnableFullMeshTransmitWithinPool () error {
 //parameters: MeshHandle
 //AgtEndpointMesh EnableFullMeshTransmitWithinPool
 return nil
}

func(np *EndpointMesh) DisableFullMeshTransmitWithinPool () error {
 //parameters: MeshHandle
 //AgtEndpointMesh DisableFullMeshTransmitWithinPool
 return nil
}

func(np *EndpointMesh) IsFullMeshTransmitWithinPoolEnabled () error {
 //parameters: MeshHandle
 //AgtEndpointMesh IsFullMeshTransmitWithinPoolEnabled
 return nil
}

func(np *EndpointMesh) ListStreamGroupsBySourceDestinationPair ()(string, error) {
 //parameters: MeshHandle Source Destination
 //AgtEndpointMesh ListStreamGroupsBySourceDestinationPair
 return "", nil
}

func(np *EndpointMesh) GetEndpointsForStreamGroup ()(string, error) {
 //parameters: MeshHandle StreamGroup
 //AgtEndpointMesh GetEndpointsForStreamGroup
 return "", nil
}

func(np *EndpointMesh) Refresh () error {
 //parameters: MeshHandle
 //AgtEndpointMesh Refresh
 return nil
}

func(np *EndpointMesh) SetIpMode () error {
 //parameters: MeshHandle IpMode
 //AgtEndpointMesh SetIpMode
 return nil
}

func(np *EndpointMesh) GetIpMode ()(string, error) {
 //parameters: MeshHandle
 //AgtEndpointMesh GetIpMode
 return "", nil
}

func(np *EndpointMesh) EnableEmulationEncapsulation () error {
 //parameters: MeshHandle
 //AgtEndpointMesh EnableEmulationEncapsulation
 return nil
}

func(np *EndpointMesh) DisableEmulationEncapsulation () error {
 //parameters: MeshHandle
 //AgtEndpointMesh DisableEmulationEncapsulation
 return nil
}

func(np *EndpointMesh) IsEmulationEncapsulationEnabled () error {
 //parameters: MeshHandle
 //AgtEndpointMesh IsEmulationEncapsulationEnabled
 return nil
}

