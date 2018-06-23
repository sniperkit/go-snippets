package n2xsdk

type PppMesh struct {
 Handler string
}

func(np *PppMesh) GetType ()(string, error) {
 //parameters: Handle
 //AgtPppMesh GetType
 return "", nil
}

func(np *PppMesh) GetName ()(string, error) {
 //parameters: Handle
 //AgtPppMesh GetName
 return "", nil
}

func(np *PppMesh) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtPppMesh GetLockCount
 return "", nil
}

func(np *PppMesh) Enable () error {
 //parameters: MeshHandle
 //AgtPppMesh Enable
 return nil
}

func(np *PppMesh) Disable () error {
 //parameters: MeshHandle
 //AgtPppMesh Disable
 return nil
}

func(np *PppMesh) IsEnabled () error {
 //parameters: MeshHandle
 //AgtPppMesh IsEnabled
 return nil
}

func(np *PppMesh) SetTrafficOrientation () error {
 //parameters: MeshHandle Orientation
 //AgtPppMesh SetTrafficOrientation
 return nil
}

func(np *PppMesh) GetTrafficOrientation ()(string, error) {
 //parameters: MeshHandle
 //AgtPppMesh GetTrafficOrientation
 return "", nil
}

func(np *PppMesh) SetLengthMode () error {
 //parameters: MeshHandle LengthMode
 //AgtPppMesh SetLengthMode
 return nil
}

func(np *PppMesh) GetLengthMode ()(string, error) {
 //parameters: MeshHandle
 //AgtPppMesh GetLengthMode
 return "", nil
}

func(np *PppMesh) SetLength () error {
 //parameters: MeshHandle LengthType Count psaLengthParameterList
 //AgtPppMesh SetLength
 return nil
}

func(np *PppMesh) GetLength ()(string, error) {
 //parameters: MeshHandle
 //AgtPppMesh GetLength
 return "", nil
}

func(np *PppMesh) ListSources ()(string, error) {
 //parameters: MeshHandle
 //AgtPppMesh ListSources
 return "", nil
}

func(np *PppMesh) ListDestinations ()(string, error) {
 //parameters: MeshHandle
 //AgtPppMesh ListDestinations
 return "", nil
}

func(np *PppMesh) ListStreamGroups ()(string, error) {
 //parameters: MeshHandle
 //AgtPppMesh ListStreamGroups
 return "", nil
}

func(np *PppMesh) ListStreamGroupsBySource ()(string, error) {
 //parameters: MeshHandle Source
 //AgtPppMesh ListStreamGroupsBySource
 return "", nil
}

func(np *PppMesh) ListStreamGroupProfiles ()(string, error) {
 //parameters: MeshHandle StreamGroupCount psaStreamGroups
 //AgtPppMesh ListStreamGroupProfiles
 return "", nil
}

func(np *PppMesh) EnableStreamGeneration () error {
 //parameters: MeshHandle
 //AgtPppMesh EnableStreamGeneration
 return nil
}

func(np *PppMesh) DisableStreamGeneration () error {
 //parameters: MeshHandle
 //AgtPppMesh DisableStreamGeneration
 return nil
}

func(np *PppMesh) IsStreamGenerationEnabled () error {
 //parameters: MeshHandle
 //AgtPppMesh IsStreamGenerationEnabled
 return nil
}

func(np *PppMesh) SetAverageLoad () error {
 //parameters: MeshHandle AverageLoad Units
 //AgtPppMesh SetAverageLoad
 return nil
}

func(np *PppMesh) GetAverageLoad ()(string, error) {
 //parameters: MeshHandle
 //AgtPppMesh GetAverageLoad
 return "", nil
}

func(np *PppMesh) IsMeshAddedOnNewProfiles () error {
 //parameters: MeshHandle
 //AgtPppMesh IsMeshAddedOnNewProfiles
 return nil
}

func(np *PppMesh) AddSessionPools () error {
 //parameters: MeshHandle SourceCount SessionPools DestinationCount DestinationPorts
 //AgtPppMesh AddSessionPools
 return nil
}

func(np *PppMesh) RemoveSessionPools () error {
 //parameters: MeshHandle Count SessionPools
 //AgtPppMesh RemoveSessionPools
 return nil
}

func(np *PppMesh) ListSessionPools ()(string, error) {
 //parameters: MeshHandle
 //AgtPppMesh ListSessionPools
 return "", nil
}

func(np *PppMesh) GetStreamGroup ()(string, error) {
 //parameters: MeshHandle SessionPool
 //AgtPppMesh GetStreamGroup
 return "", nil
}

func(np *PppMesh) GetSessionPoolForStreamGroup ()(string, error) {
 //parameters: MeshHandle StreamGroup
 //AgtPppMesh GetSessionPoolForStreamGroup
 return "", nil
}

func(np *PppMesh) SetEndpointTypes () error {
 //parameters: MeshHandle SourceType DestinationType
 //AgtPppMesh SetEndpointTypes
 return nil
}

func(np *PppMesh) GetEndpointTypes ()(string, error) {
 //parameters: MeshHandle
 //AgtPppMesh GetEndpointTypes
 return "", nil
}

func(np *PppMesh) AddSourcesAndDestinations () error {
 //parameters: MeshHandle SourceCount Sources DestinationCount Destinations
 //AgtPppMesh AddSourcesAndDestinations
 return nil
}

func(np *PppMesh) AddDestinationsBySource () error {
 //parameters: MeshHandle Source DestinationCount Destinations
 //AgtPppMesh AddDestinationsBySource
 return nil
}

func(np *PppMesh) ListAllSources ()(string, error) {
 //parameters: MeshHandle
 //AgtPppMesh ListAllSources
 return "", nil
}

func(np *PppMesh) ListAllDestinations ()(string, error) {
 //parameters: MeshHandle
 //AgtPppMesh ListAllDestinations
 return "", nil
}

func(np *PppMesh) ListDestinationsBySource ()(string, error) {
 //parameters: MeshHandle Source
 //AgtPppMesh ListDestinationsBySource
 return "", nil
}

func(np *PppMesh) GetStreamGroupForSourceAndDestination ()(string, error) {
 //parameters: MeshHandle Source Destination
 //AgtPppMesh GetStreamGroupForSourceAndDestination
 return "", nil
}

func(np *PppMesh) GetSourceAndDestinationForStreamGroup ()(string, error) {
 //parameters: MeshHandle StreamGroup
 //AgtPppMesh GetSourceAndDestinationForStreamGroup
 return "", nil
}

func(np *PppMesh) RemoveSources () error {
 //parameters: MeshHandle Count Sources
 //AgtPppMesh RemoveSources
 return nil
}

func(np *PppMesh) RemoveDestinationsBySource () error {
 //parameters: MeshHandle Source Count Destinations
 //AgtPppMesh RemoveDestinationsBySource
 return nil
}

func(np *PppMesh) Clear () error {
 //parameters: MeshHandle
 //AgtPppMesh Clear
 return nil
}

