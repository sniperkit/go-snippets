package n2xsdk

type DeviceMesh struct {
 Handler string
}

func(np *DeviceMesh) GetType ()(string, error) {
 //parameters: Handle
 //AgtDeviceMesh GetType
 return "", nil
}

func(np *DeviceMesh) GetName ()(string, error) {
 //parameters: Handle
 //AgtDeviceMesh GetName
 return "", nil
}

func(np *DeviceMesh) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtDeviceMesh GetLockCount
 return "", nil
}

func(np *DeviceMesh) Enable () error {
 //parameters: MeshHandle
 //AgtDeviceMesh Enable
 return nil
}

func(np *DeviceMesh) Disable () error {
 //parameters: MeshHandle
 //AgtDeviceMesh Disable
 return nil
}

func(np *DeviceMesh) IsEnabled () error {
 //parameters: MeshHandle
 //AgtDeviceMesh IsEnabled
 return nil
}

func(np *DeviceMesh) SetTrafficOrientation () error {
 //parameters: MeshHandle Orientation
 //AgtDeviceMesh SetTrafficOrientation
 return nil
}

func(np *DeviceMesh) GetTrafficOrientation ()(string, error) {
 //parameters: MeshHandle
 //AgtDeviceMesh GetTrafficOrientation
 return "", nil
}

func(np *DeviceMesh) SetLengthMode () error {
 //parameters: MeshHandle LengthMode
 //AgtDeviceMesh SetLengthMode
 return nil
}

func(np *DeviceMesh) GetLengthMode ()(string, error) {
 //parameters: MeshHandle
 //AgtDeviceMesh GetLengthMode
 return "", nil
}

func(np *DeviceMesh) SetLength () error {
 //parameters: MeshHandle LengthType Count psaLengthParameterList
 //AgtDeviceMesh SetLength
 return nil
}

func(np *DeviceMesh) GetLength ()(string, error) {
 //parameters: MeshHandle
 //AgtDeviceMesh GetLength
 return "", nil
}

func(np *DeviceMesh) ListSources ()(string, error) {
 //parameters: MeshHandle
 //AgtDeviceMesh ListSources
 return "", nil
}

func(np *DeviceMesh) ListDestinations ()(string, error) {
 //parameters: MeshHandle
 //AgtDeviceMesh ListDestinations
 return "", nil
}

func(np *DeviceMesh) ListStreamGroups ()(string, error) {
 //parameters: MeshHandle
 //AgtDeviceMesh ListStreamGroups
 return "", nil
}

func(np *DeviceMesh) ListStreamGroupsBySource ()(string, error) {
 //parameters: MeshHandle Source
 //AgtDeviceMesh ListStreamGroupsBySource
 return "", nil
}

func(np *DeviceMesh) ListStreamGroupProfiles ()(string, error) {
 //parameters: MeshHandle StreamGroupCount psaStreamGroups
 //AgtDeviceMesh ListStreamGroupProfiles
 return "", nil
}

func(np *DeviceMesh) EnableStreamGeneration () error {
 //parameters: MeshHandle
 //AgtDeviceMesh EnableStreamGeneration
 return nil
}

func(np *DeviceMesh) DisableStreamGeneration () error {
 //parameters: MeshHandle
 //AgtDeviceMesh DisableStreamGeneration
 return nil
}

func(np *DeviceMesh) IsStreamGenerationEnabled () error {
 //parameters: MeshHandle
 //AgtDeviceMesh IsStreamGenerationEnabled
 return nil
}

func(np *DeviceMesh) SetAverageLoad () error {
 //parameters: MeshHandle AverageLoad Units
 //AgtDeviceMesh SetAverageLoad
 return nil
}

func(np *DeviceMesh) GetAverageLoad ()(string, error) {
 //parameters: MeshHandle
 //AgtDeviceMesh GetAverageLoad
 return "", nil
}

func(np *DeviceMesh) IsMeshAddedOnNewProfiles () error {
 //parameters: MeshHandle
 //AgtDeviceMesh IsMeshAddedOnNewProfiles
 return nil
}

func(np *DeviceMesh) SetEndpointTypes () error {
 //parameters: MeshHandle SourceType DestinationType
 //AgtDeviceMesh SetEndpointTypes
 return nil
}

func(np *DeviceMesh) GetEndpointTypes ()(string, error) {
 //parameters: MeshHandle
 //AgtDeviceMesh GetEndpointTypes
 return "", nil
}

func(np *DeviceMesh) EnableManualSelectionOfAddressTypes () error {
 //parameters: MeshHandle
 //AgtDeviceMesh EnableManualSelectionOfAddressTypes
 return nil
}

func(np *DeviceMesh) DisableManualSelectionOfAddressTypes () error {
 //parameters: MeshHandle
 //AgtDeviceMesh DisableManualSelectionOfAddressTypes
 return nil
}

func(np *DeviceMesh) IsManualSelectionOfAddressTypesEnabled () error {
 //parameters: MeshHandle
 //AgtDeviceMesh IsManualSelectionOfAddressTypesEnabled
 return nil
}

func(np *DeviceMesh) SetEndpointAddressTypes () error {
 //parameters: MeshHandle SourceAddressType DestinationAddressType
 //AgtDeviceMesh SetEndpointAddressTypes
 return nil
}

func(np *DeviceMesh) GetEndpointAddressTypes ()(string, error) {
 //parameters: MeshHandle
 //AgtDeviceMesh GetEndpointAddressTypes
 return "", nil
}

func(np *DeviceMesh) ListSupportedEndpointsByAddressType ()(string, error) {
 //parameters: EndpointType EndpointAddressType
 //AgtDeviceMesh ListSupportedEndpointsByAddressType
 return "", nil
}

func(np *DeviceMesh) ListSupportedEndpointAddressTypes ()(string, error) {
 //parameters: EndpointType Endpoint
 //AgtDeviceMesh ListSupportedEndpointAddressTypes
 return "", nil
}

func(np *DeviceMesh) SetTrafficDistribution () error {
 //parameters: MeshHandle Distribution
 //AgtDeviceMesh SetTrafficDistribution
 return nil
}

func(np *DeviceMesh) GetTrafficDistribution ()(string, error) {
 //parameters: MeshHandle
 //AgtDeviceMesh GetTrafficDistribution
 return "", nil
}

func(np *DeviceMesh) SetDeviceUnderTestType () error {
 //parameters: MeshHandle DutType
 //AgtDeviceMesh SetDeviceUnderTestType
 return nil
}

func(np *DeviceMesh) GetDeviceUnderTestType ()(string, error) {
 //parameters: MeshHandle
 //AgtDeviceMesh GetDeviceUnderTestType
 return "", nil
}

func(np *DeviceMesh) AddSources () error {
 //parameters: MeshHandle Count psaSources
 //AgtDeviceMesh AddSources
 return nil
}

func(np *DeviceMesh) RemoveSources () error {
 //parameters: MeshHandle Count psaSources
 //AgtDeviceMesh RemoveSources
 return nil
}

func(np *DeviceMesh) AddDestinations () error {
 //parameters: MeshHandle Count psaDestinations
 //AgtDeviceMesh AddDestinations
 return nil
}

func(np *DeviceMesh) RemoveDestinations () error {
 //parameters: MeshHandle Count psaDestinations
 //AgtDeviceMesh RemoveDestinations
 return nil
}

func(np *DeviceMesh) RemoveDestinationsForSources () error {
 //parameters: MeshHandle NumSources psaSources NumDestinations psaDestinations
 //AgtDeviceMesh RemoveDestinationsForSources
 return nil
}

func(np *DeviceMesh) SetRemoteHostIpv6AddressSuffix () error {
 //parameters: MeshHandle Endpoint HostIpv6AddressSuffix
 //AgtDeviceMesh SetRemoteHostIpv6AddressSuffix
 return nil
}

func(np *DeviceMesh) GetRemoteHostIpv6AddressSuffix ()(string, error) {
 //parameters: MeshHandle Endpoint
 //AgtDeviceMesh GetRemoteHostIpv6AddressSuffix
 return "", nil
}

func(np *DeviceMesh) SetRemoteHostCount () error {
 //parameters: MeshHandle Endpoint HostCountPerClient
 //AgtDeviceMesh SetRemoteHostCount
 return nil
}

func(np *DeviceMesh) GetRemoteHostCount ()(string, error) {
 //parameters: MeshHandle Endpoint
 //AgtDeviceMesh GetRemoteHostCount
 return "", nil
}

func(np *DeviceMesh) Clear () error {
 //parameters: MeshHandle
 //AgtDeviceMesh Clear
 return nil
}

func(np *DeviceMesh) GetForwardStreamGroup ()(string, error) {
 //parameters: MeshHandle Source Destination
 //AgtDeviceMesh GetForwardStreamGroup
 return "", nil
}

func(np *DeviceMesh) GetReverseStreamGroup ()(string, error) {
 //parameters: MeshHandle Source Destination
 //AgtDeviceMesh GetReverseStreamGroup
 return "", nil
}

func(np *DeviceMesh) GetEndpointsForStreamGroup ()(string, error) {
 //parameters: MeshHandle StreamGroup
 //AgtDeviceMesh GetEndpointsForStreamGroup
 return "", nil
}

func(np *DeviceMesh) ListSourcesForDestination ()(string, error) {
 //parameters: MeshHandle Destination
 //AgtDeviceMesh ListSourcesForDestination
 return "", nil
}

func(np *DeviceMesh) ListDestinationsForSource ()(string, error) {
 //parameters: MeshHandle Source
 //AgtDeviceMesh ListDestinationsForSource
 return "", nil
}

func(np *DeviceMesh) UpdateStreamGroups () error {
 //parameters: MeshHandle
 //AgtDeviceMesh UpdateStreamGroups
 return nil
}

func(np *DeviceMesh) SetStreamGenerationParameter () error {
 //parameters: MeshHandle Protocol ProtocolInstance Field
 //AgtDeviceMesh SetStreamGenerationParameter
 return nil
}

func(np *DeviceMesh) GetStreamGenerationParameter ()(string, error) {
 //parameters: MeshHandle
 //AgtDeviceMesh GetStreamGenerationParameter
 return "", nil
}

func(np *DeviceMesh) ResetPduHeaders () error {
 //parameters: MeshHandle
 //AgtDeviceMesh ResetPduHeaders
 return nil
}

