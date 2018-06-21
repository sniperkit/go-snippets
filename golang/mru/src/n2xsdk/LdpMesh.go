package n2xsdk

type LdpMesh struct {
 Handler string
}

func(np *LdpMesh) GetType ()(string, error) {
 //parameters: Handle
 //AgtLdpMesh GetType
 return "", nil
}

func(np *LdpMesh) GetName ()(string, error) {
 //parameters: Handle
 //AgtLdpMesh GetName
 return "", nil
}

func(np *LdpMesh) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtLdpMesh GetLockCount
 return "", nil
}

func(np *LdpMesh) Enable () error {
 //parameters: MeshHandle
 //AgtLdpMesh Enable
 return nil
}

func(np *LdpMesh) Disable () error {
 //parameters: MeshHandle
 //AgtLdpMesh Disable
 return nil
}

func(np *LdpMesh) IsEnabled () error {
 //parameters: MeshHandle
 //AgtLdpMesh IsEnabled
 return nil
}

func(np *LdpMesh) SetTrafficOrientation () error {
 //parameters: MeshHandle Orientation
 //AgtLdpMesh SetTrafficOrientation
 return nil
}

func(np *LdpMesh) GetTrafficOrientation ()(string, error) {
 //parameters: MeshHandle
 //AgtLdpMesh GetTrafficOrientation
 return "", nil
}

func(np *LdpMesh) SetLengthMode () error {
 //parameters: MeshHandle LengthMode
 //AgtLdpMesh SetLengthMode
 return nil
}

func(np *LdpMesh) GetLengthMode ()(string, error) {
 //parameters: MeshHandle
 //AgtLdpMesh GetLengthMode
 return "", nil
}

func(np *LdpMesh) SetLength () error {
 //parameters: MeshHandle LengthType Count psaLengthParameterList
 //AgtLdpMesh SetLength
 return nil
}

func(np *LdpMesh) GetLength ()(string, error) {
 //parameters: MeshHandle
 //AgtLdpMesh GetLength
 return "", nil
}

func(np *LdpMesh) ListSources ()(string, error) {
 //parameters: MeshHandle
 //AgtLdpMesh ListSources
 return "", nil
}

func(np *LdpMesh) ListDestinations ()(string, error) {
 //parameters: MeshHandle
 //AgtLdpMesh ListDestinations
 return "", nil
}

func(np *LdpMesh) ListStreamGroups ()(string, error) {
 //parameters: MeshHandle
 //AgtLdpMesh ListStreamGroups
 return "", nil
}

func(np *LdpMesh) ListStreamGroupsBySource ()(string, error) {
 //parameters: MeshHandle Source
 //AgtLdpMesh ListStreamGroupsBySource
 return "", nil
}

func(np *LdpMesh) ListStreamGroupProfiles ()(string, error) {
 //parameters: MeshHandle StreamGroupCount psaStreamGroups
 //AgtLdpMesh ListStreamGroupProfiles
 return "", nil
}

func(np *LdpMesh) EnableStreamGeneration () error {
 //parameters: MeshHandle
 //AgtLdpMesh EnableStreamGeneration
 return nil
}

func(np *LdpMesh) DisableStreamGeneration () error {
 //parameters: MeshHandle
 //AgtLdpMesh DisableStreamGeneration
 return nil
}

func(np *LdpMesh) IsStreamGenerationEnabled () error {
 //parameters: MeshHandle
 //AgtLdpMesh IsStreamGenerationEnabled
 return nil
}

func(np *LdpMesh) SetAverageLoad () error {
 //parameters: MeshHandle AverageLoad Units
 //AgtLdpMesh SetAverageLoad
 return nil
}

func(np *LdpMesh) GetAverageLoad ()(string, error) {
 //parameters: MeshHandle
 //AgtLdpMesh GetAverageLoad
 return "", nil
}

func(np *LdpMesh) IsMeshAddedOnNewProfiles () error {
 //parameters: MeshHandle
 //AgtLdpMesh IsMeshAddedOnNewProfiles
 return nil
}

func(np *LdpMesh) AddLspAddressPool () error {
 //parameters: MeshHandle SourceLdpSessionHandle DestinationPort DestinationFirstIpAddress DestinationPrefixLength DestinationNumAddresses DestinationAddressModifier
 //AgtLdpMesh AddLspAddressPool
 return nil
}

func(np *LdpMesh) RemoveLspAddressPool () error {
 //parameters: MeshHandle AddressPoolHandle
 //AgtLdpMesh RemoveLspAddressPool
 return nil
}

func(np *LdpMesh) ListLspAddressPools ()(string, error) {
 //parameters: MeshHandle
 //AgtLdpMesh ListLspAddressPools
 return "", nil
}

func(np *LdpMesh) GetLspAddressPool ()(string, error) {
 //parameters: MeshHandle AddressPoolHandle
 //AgtLdpMesh GetLspAddressPool
 return "", nil
}

func(np *LdpMesh) GetTotalDestinationAddresses ()(string, error) {
 //parameters: MeshHandle
 //AgtLdpMesh GetTotalDestinationAddresses
 return "", nil
}

func(np *LdpMesh) GetStreamGroup ()(string, error) {
 //parameters: MeshHandle AddressPoolHandle
 //AgtLdpMesh GetStreamGroup
 return "", nil
}

