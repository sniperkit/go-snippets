package n2xsdk

type Ipv6Mesh struct {
 Handler string
}

func(np *Ipv6Mesh) GetType ()(string, error) {
 //parameters: Handle
 //AgtIpv6Mesh GetType
 return "", nil
}

func(np *Ipv6Mesh) GetName ()(string, error) {
 //parameters: Handle
 //AgtIpv6Mesh GetName
 return "", nil
}

func(np *Ipv6Mesh) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtIpv6Mesh GetLockCount
 return "", nil
}

func(np *Ipv6Mesh) Enable () error {
 //parameters: MeshHandle
 //AgtIpv6Mesh Enable
 return nil
}

func(np *Ipv6Mesh) Disable () error {
 //parameters: MeshHandle
 //AgtIpv6Mesh Disable
 return nil
}

func(np *Ipv6Mesh) IsEnabled () error {
 //parameters: MeshHandle
 //AgtIpv6Mesh IsEnabled
 return nil
}

func(np *Ipv6Mesh) SetTrafficOrientation () error {
 //parameters: MeshHandle Orientation
 //AgtIpv6Mesh SetTrafficOrientation
 return nil
}

func(np *Ipv6Mesh) GetTrafficOrientation ()(string, error) {
 //parameters: MeshHandle
 //AgtIpv6Mesh GetTrafficOrientation
 return "", nil
}

func(np *Ipv6Mesh) SetLengthMode () error {
 //parameters: MeshHandle LengthMode
 //AgtIpv6Mesh SetLengthMode
 return nil
}

func(np *Ipv6Mesh) GetLengthMode ()(string, error) {
 //parameters: MeshHandle
 //AgtIpv6Mesh GetLengthMode
 return "", nil
}

func(np *Ipv6Mesh) SetLength () error {
 //parameters: MeshHandle LengthType Count psaLengthParameterList
 //AgtIpv6Mesh SetLength
 return nil
}

func(np *Ipv6Mesh) GetLength ()(string, error) {
 //parameters: MeshHandle
 //AgtIpv6Mesh GetLength
 return "", nil
}

func(np *Ipv6Mesh) ListSources ()(string, error) {
 //parameters: MeshHandle
 //AgtIpv6Mesh ListSources
 return "", nil
}

func(np *Ipv6Mesh) ListDestinations ()(string, error) {
 //parameters: MeshHandle
 //AgtIpv6Mesh ListDestinations
 return "", nil
}

func(np *Ipv6Mesh) ListStreamGroups ()(string, error) {
 //parameters: MeshHandle
 //AgtIpv6Mesh ListStreamGroups
 return "", nil
}

func(np *Ipv6Mesh) ListStreamGroupsBySource ()(string, error) {
 //parameters: MeshHandle Source
 //AgtIpv6Mesh ListStreamGroupsBySource
 return "", nil
}

func(np *Ipv6Mesh) ListStreamGroupProfiles ()(string, error) {
 //parameters: MeshHandle StreamGroupCount psaStreamGroups
 //AgtIpv6Mesh ListStreamGroupProfiles
 return "", nil
}

func(np *Ipv6Mesh) EnableStreamGeneration () error {
 //parameters: MeshHandle
 //AgtIpv6Mesh EnableStreamGeneration
 return nil
}

func(np *Ipv6Mesh) DisableStreamGeneration () error {
 //parameters: MeshHandle
 //AgtIpv6Mesh DisableStreamGeneration
 return nil
}

func(np *Ipv6Mesh) IsStreamGenerationEnabled () error {
 //parameters: MeshHandle
 //AgtIpv6Mesh IsStreamGenerationEnabled
 return nil
}

func(np *Ipv6Mesh) SetAverageLoad () error {
 //parameters: MeshHandle AverageLoad Units
 //AgtIpv6Mesh SetAverageLoad
 return nil
}

func(np *Ipv6Mesh) GetAverageLoad ()(string, error) {
 //parameters: MeshHandle
 //AgtIpv6Mesh GetAverageLoad
 return "", nil
}

func(np *Ipv6Mesh) IsMeshAddedOnNewProfiles () error {
 //parameters: MeshHandle
 //AgtIpv6Mesh IsMeshAddedOnNewProfiles
 return nil
}

func(np *Ipv6Mesh) AddSources () error {
 //parameters: MeshHandle Count psaSourcePorts
 //AgtIpv6Mesh AddSources
 return nil
}

func(np *Ipv6Mesh) RemoveSources () error {
 //parameters: MeshHandle Count psaSourcePorts
 //AgtIpv6Mesh RemoveSources
 return nil
}

func(np *Ipv6Mesh) AddDestinations () error {
 //parameters: MeshHandle Count psaDestinationPorts
 //AgtIpv6Mesh AddDestinations
 return nil
}

func(np *Ipv6Mesh) RemoveDestinations () error {
 //parameters: MeshHandle Count psaDestinationPorts
 //AgtIpv6Mesh RemoveDestinations
 return nil
}

func(np *Ipv6Mesh) SetTrafficDistribution () error {
 //parameters: MeshHandle Distribution
 //AgtIpv6Mesh SetTrafficDistribution
 return nil
}

func(np *Ipv6Mesh) GetTrafficDistribution ()(string, error) {
 //parameters: MeshHandle
 //AgtIpv6Mesh GetTrafficDistribution
 return "", nil
}

func(np *Ipv6Mesh) GetStreamGroups ()(string, error) {
 //parameters: MeshHandle Source Destination
 //AgtIpv6Mesh GetStreamGroups
 return "", nil
}

func(np *Ipv6Mesh) GetStreamGroup ()(string, error) {
 //parameters: MeshHandle Source Destination RoutePool
 //AgtIpv6Mesh GetStreamGroup
 return "", nil
}

func(np *Ipv6Mesh) GetPortsForStreamGroup ()(string, error) {
 //parameters: MeshHandle StreamGroup
 //AgtIpv6Mesh GetPortsForStreamGroup
 return "", nil
}

func(np *Ipv6Mesh) GetRoutePoolForStreamGroup ()(string, error) {
 //parameters: MeshHandle StreamGroup
 //AgtIpv6Mesh GetRoutePoolForStreamGroup
 return "", nil
}

func(np *Ipv6Mesh) SetSourceVirtualInterface () error {
 //parameters: MeshHandle Count psaInterfaceHandles
 //AgtIpv6Mesh SetSourceVirtualInterface
 return nil
}

func(np *Ipv6Mesh) GetSourceVirtualInterface ()(string, error) {
 //parameters: MeshHandle Source
 //AgtIpv6Mesh GetSourceVirtualInterface
 return "", nil
}

func(np *Ipv6Mesh) SetVirtualInterfaces () error {
 //parameters: MeshHandle Port Count psaVirtualInterfaces
 //AgtIpv6Mesh SetVirtualInterfaces
 return nil
}

func(np *Ipv6Mesh) GetVirtualInterfaces ()(string, error) {
 //parameters: MeshHandle Port
 //AgtIpv6Mesh GetVirtualInterfaces
 return "", nil
}

func(np *Ipv6Mesh) ListAvailableDestinationRoutePools ()(string, error) {
 //parameters: MeshHandle Destination
 //AgtIpv6Mesh ListAvailableDestinationRoutePools
 return "", nil
}

func(np *Ipv6Mesh) ListSelectedDestinationRoutePools ()(string, error) {
 //parameters: MeshHandle Destination
 //AgtIpv6Mesh ListSelectedDestinationRoutePools
 return "", nil
}

func(np *Ipv6Mesh) GetTotalDestinationRoutesInUse ()(string, error) {
 //parameters: MeshHandle Source Destination
 //AgtIpv6Mesh GetTotalDestinationRoutesInUse
 return "", nil
}

func(np *Ipv6Mesh) UpdateStreamGroups () error {
 //parameters: MeshHandle
 //AgtIpv6Mesh UpdateStreamGroups
 return nil
}

func(np *Ipv6Mesh) SetEncapsulationProtocols () error {
 //parameters: MeshHandle Source Count psaProtocols
 //AgtIpv6Mesh SetEncapsulationProtocols
 return nil
}

func(np *Ipv6Mesh) GetEncapsulationProtocols ()(string, error) {
 //parameters: MeshHandle Source
 //AgtIpv6Mesh GetEncapsulationProtocols
 return "", nil
}

