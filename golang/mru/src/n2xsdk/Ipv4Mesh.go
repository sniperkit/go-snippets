package n2xsdk

type Ipv4Mesh struct {
 Handler string
}

func(np *Ipv4Mesh) GetType ()(string, error) {
 //parameters: Handle
 //AgtIpv4Mesh GetType
 return "", nil
}

func(np *Ipv4Mesh) GetName ()(string, error) {
 //parameters: Handle
 //AgtIpv4Mesh GetName
 return "", nil
}

func(np *Ipv4Mesh) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtIpv4Mesh GetLockCount
 return "", nil
}

func(np *Ipv4Mesh) Enable () error {
 //parameters: MeshHandle
 //AgtIpv4Mesh Enable
 return nil
}

func(np *Ipv4Mesh) Disable () error {
 //parameters: MeshHandle
 //AgtIpv4Mesh Disable
 return nil
}

func(np *Ipv4Mesh) IsEnabled () error {
 //parameters: MeshHandle
 //AgtIpv4Mesh IsEnabled
 return nil
}

func(np *Ipv4Mesh) SetTrafficOrientation () error {
 //parameters: MeshHandle Orientation
 //AgtIpv4Mesh SetTrafficOrientation
 return nil
}

func(np *Ipv4Mesh) GetTrafficOrientation ()(string, error) {
 //parameters: MeshHandle
 //AgtIpv4Mesh GetTrafficOrientation
 return "", nil
}

func(np *Ipv4Mesh) SetLengthMode () error {
 //parameters: MeshHandle LengthMode
 //AgtIpv4Mesh SetLengthMode
 return nil
}

func(np *Ipv4Mesh) GetLengthMode ()(string, error) {
 //parameters: MeshHandle
 //AgtIpv4Mesh GetLengthMode
 return "", nil
}

func(np *Ipv4Mesh) SetLength () error {
 //parameters: MeshHandle LengthType Count psaLengthParameterList
 //AgtIpv4Mesh SetLength
 return nil
}

func(np *Ipv4Mesh) GetLength ()(string, error) {
 //parameters: MeshHandle
 //AgtIpv4Mesh GetLength
 return "", nil
}

func(np *Ipv4Mesh) ListSources ()(string, error) {
 //parameters: MeshHandle
 //AgtIpv4Mesh ListSources
 return "", nil
}

func(np *Ipv4Mesh) ListDestinations ()(string, error) {
 //parameters: MeshHandle
 //AgtIpv4Mesh ListDestinations
 return "", nil
}

func(np *Ipv4Mesh) ListStreamGroups ()(string, error) {
 //parameters: MeshHandle
 //AgtIpv4Mesh ListStreamGroups
 return "", nil
}

func(np *Ipv4Mesh) ListStreamGroupsBySource ()(string, error) {
 //parameters: MeshHandle Source
 //AgtIpv4Mesh ListStreamGroupsBySource
 return "", nil
}

func(np *Ipv4Mesh) ListStreamGroupProfiles ()(string, error) {
 //parameters: MeshHandle StreamGroupCount psaStreamGroups
 //AgtIpv4Mesh ListStreamGroupProfiles
 return "", nil
}

func(np *Ipv4Mesh) EnableStreamGeneration () error {
 //parameters: MeshHandle
 //AgtIpv4Mesh EnableStreamGeneration
 return nil
}

func(np *Ipv4Mesh) DisableStreamGeneration () error {
 //parameters: MeshHandle
 //AgtIpv4Mesh DisableStreamGeneration
 return nil
}

func(np *Ipv4Mesh) IsStreamGenerationEnabled () error {
 //parameters: MeshHandle
 //AgtIpv4Mesh IsStreamGenerationEnabled
 return nil
}

func(np *Ipv4Mesh) SetAverageLoad () error {
 //parameters: MeshHandle AverageLoad Units
 //AgtIpv4Mesh SetAverageLoad
 return nil
}

func(np *Ipv4Mesh) GetAverageLoad ()(string, error) {
 //parameters: MeshHandle
 //AgtIpv4Mesh GetAverageLoad
 return "", nil
}

func(np *Ipv4Mesh) IsMeshAddedOnNewProfiles () error {
 //parameters: MeshHandle
 //AgtIpv4Mesh IsMeshAddedOnNewProfiles
 return nil
}

func(np *Ipv4Mesh) AddSources () error {
 //parameters: MeshHandle Count psaSourcePorts
 //AgtIpv4Mesh AddSources
 return nil
}

func(np *Ipv4Mesh) RemoveSources () error {
 //parameters: MeshHandle Count psaSourcePorts
 //AgtIpv4Mesh RemoveSources
 return nil
}

func(np *Ipv4Mesh) AddDestinations () error {
 //parameters: MeshHandle Count psaDestinationPorts
 //AgtIpv4Mesh AddDestinations
 return nil
}

func(np *Ipv4Mesh) RemoveDestinations () error {
 //parameters: MeshHandle Count psaDestinationPorts
 //AgtIpv4Mesh RemoveDestinations
 return nil
}

func(np *Ipv4Mesh) SetTrafficDistribution () error {
 //parameters: MeshHandle Distribution
 //AgtIpv4Mesh SetTrafficDistribution
 return nil
}

func(np *Ipv4Mesh) GetTrafficDistribution ()(string, error) {
 //parameters: MeshHandle
 //AgtIpv4Mesh GetTrafficDistribution
 return "", nil
}

func(np *Ipv4Mesh) GetStreamGroups ()(string, error) {
 //parameters: MeshHandle Source Destination
 //AgtIpv4Mesh GetStreamGroups
 return "", nil
}

func(np *Ipv4Mesh) GetStreamGroup ()(string, error) {
 //parameters: MeshHandle Source Destination RoutePool
 //AgtIpv4Mesh GetStreamGroup
 return "", nil
}

func(np *Ipv4Mesh) GetPortsForStreamGroup ()(string, error) {
 //parameters: MeshHandle StreamGroup
 //AgtIpv4Mesh GetPortsForStreamGroup
 return "", nil
}

func(np *Ipv4Mesh) GetRoutePoolForStreamGroup ()(string, error) {
 //parameters: MeshHandle StreamGroup
 //AgtIpv4Mesh GetRoutePoolForStreamGroup
 return "", nil
}

func(np *Ipv4Mesh) SetSourceVirtualInterface () error {
 //parameters: MeshHandle Count psaInterfaceHandles
 //AgtIpv4Mesh SetSourceVirtualInterface
 return nil
}

func(np *Ipv4Mesh) GetSourceVirtualInterface ()(string, error) {
 //parameters: MeshHandle Source
 //AgtIpv4Mesh GetSourceVirtualInterface
 return "", nil
}

func(np *Ipv4Mesh) SetVirtualInterfaces () error {
 //parameters: MeshHandle Port Count psaVirtualInterfaces
 //AgtIpv4Mesh SetVirtualInterfaces
 return nil
}

func(np *Ipv4Mesh) GetVirtualInterfaces ()(string, error) {
 //parameters: MeshHandle Port
 //AgtIpv4Mesh GetVirtualInterfaces
 return "", nil
}

func(np *Ipv4Mesh) ListAvailableDestinationRoutePools ()(string, error) {
 //parameters: MeshHandle Destination
 //AgtIpv4Mesh ListAvailableDestinationRoutePools
 return "", nil
}

func(np *Ipv4Mesh) ListSelectedDestinationRoutePools ()(string, error) {
 //parameters: MeshHandle Destination
 //AgtIpv4Mesh ListSelectedDestinationRoutePools
 return "", nil
}

func(np *Ipv4Mesh) GetTotalDestinationRoutesInUse ()(string, error) {
 //parameters: MeshHandle Source Destination
 //AgtIpv4Mesh GetTotalDestinationRoutesInUse
 return "", nil
}

func(np *Ipv4Mesh) UpdateStreamGroups () error {
 //parameters: MeshHandle
 //AgtIpv4Mesh UpdateStreamGroups
 return nil
}

