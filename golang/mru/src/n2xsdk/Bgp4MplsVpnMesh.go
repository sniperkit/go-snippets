package n2xsdk

type Bgp4MplsVpnMesh struct {
 Handler string
}

func(np *Bgp4MplsVpnMesh) GetType ()(string, error) {
 //parameters: Handle
 //AgtBgp4MplsVpnMesh GetType
 return "", nil
}

func(np *Bgp4MplsVpnMesh) GetName ()(string, error) {
 //parameters: Handle
 //AgtBgp4MplsVpnMesh GetName
 return "", nil
}

func(np *Bgp4MplsVpnMesh) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtBgp4MplsVpnMesh GetLockCount
 return "", nil
}

func(np *Bgp4MplsVpnMesh) Enable () error {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh Enable
 return nil
}

func(np *Bgp4MplsVpnMesh) Disable () error {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh Disable
 return nil
}

func(np *Bgp4MplsVpnMesh) IsEnabled () error {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh IsEnabled
 return nil
}

func(np *Bgp4MplsVpnMesh) SetTrafficOrientation () error {
 //parameters: MeshHandle Orientation
 //AgtBgp4MplsVpnMesh SetTrafficOrientation
 return nil
}

func(np *Bgp4MplsVpnMesh) GetTrafficOrientation ()(string, error) {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh GetTrafficOrientation
 return "", nil
}

func(np *Bgp4MplsVpnMesh) SetLengthMode () error {
 //parameters: MeshHandle LengthMode
 //AgtBgp4MplsVpnMesh SetLengthMode
 return nil
}

func(np *Bgp4MplsVpnMesh) GetLengthMode ()(string, error) {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh GetLengthMode
 return "", nil
}

func(np *Bgp4MplsVpnMesh) SetLength () error {
 //parameters: MeshHandle LengthType Count psaLengthParameterList
 //AgtBgp4MplsVpnMesh SetLength
 return nil
}

func(np *Bgp4MplsVpnMesh) GetLength ()(string, error) {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh GetLength
 return "", nil
}

func(np *Bgp4MplsVpnMesh) ListSources ()(string, error) {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh ListSources
 return "", nil
}

func(np *Bgp4MplsVpnMesh) ListDestinations ()(string, error) {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh ListDestinations
 return "", nil
}

func(np *Bgp4MplsVpnMesh) ListStreamGroups ()(string, error) {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh ListStreamGroups
 return "", nil
}

func(np *Bgp4MplsVpnMesh) ListStreamGroupsBySource ()(string, error) {
 //parameters: MeshHandle Source
 //AgtBgp4MplsVpnMesh ListStreamGroupsBySource
 return "", nil
}

func(np *Bgp4MplsVpnMesh) ListStreamGroupProfiles ()(string, error) {
 //parameters: MeshHandle StreamGroupCount psaStreamGroups
 //AgtBgp4MplsVpnMesh ListStreamGroupProfiles
 return "", nil
}

func(np *Bgp4MplsVpnMesh) EnableStreamGeneration () error {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh EnableStreamGeneration
 return nil
}

func(np *Bgp4MplsVpnMesh) DisableStreamGeneration () error {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh DisableStreamGeneration
 return nil
}

func(np *Bgp4MplsVpnMesh) IsStreamGenerationEnabled () error {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh IsStreamGenerationEnabled
 return nil
}

func(np *Bgp4MplsVpnMesh) SetAverageLoad () error {
 //parameters: MeshHandle AverageLoad Units
 //AgtBgp4MplsVpnMesh SetAverageLoad
 return nil
}

func(np *Bgp4MplsVpnMesh) GetAverageLoad ()(string, error) {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh GetAverageLoad
 return "", nil
}

func(np *Bgp4MplsVpnMesh) IsMeshAddedOnNewProfiles () error {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh IsMeshAddedOnNewProfiles
 return nil
}

func(np *Bgp4MplsVpnMesh) SetSourceSite () error {
 //parameters: MeshHandle Count psaSiteHandle
 //AgtBgp4MplsVpnMesh SetSourceSite
 return nil
}

func(np *Bgp4MplsVpnMesh) GetSourceSite ()(string, error) {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh GetSourceSite
 return "", nil
}

func(np *Bgp4MplsVpnMesh) AddDestinationAddressPool () error {
 //parameters: MeshHandle DestinationPort FirstIpAddress PrefixLength NumAddresses Modifier
 //AgtBgp4MplsVpnMesh AddDestinationAddressPool
 return nil
}

func(np *Bgp4MplsVpnMesh) RemoveDestinationAddressPool () error {
 //parameters: MeshHandle AddressPoolHandle
 //AgtBgp4MplsVpnMesh RemoveDestinationAddressPool
 return nil
}

func(np *Bgp4MplsVpnMesh) ListDestinationAddressPools ()(string, error) {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh ListDestinationAddressPools
 return "", nil
}

func(np *Bgp4MplsVpnMesh) GetDestinationAddressPool ()(string, error) {
 //parameters: MeshHandle AddressPoolHandle
 //AgtBgp4MplsVpnMesh GetDestinationAddressPool
 return "", nil
}

func(np *Bgp4MplsVpnMesh) GetTotalDestinationAddresses ()(string, error) {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh GetTotalDestinationAddresses
 return "", nil
}

func(np *Bgp4MplsVpnMesh) GetStreamGroup ()(string, error) {
 //parameters: MeshHandle AddressPoolHandle
 //AgtBgp4MplsVpnMesh GetStreamGroup
 return "", nil
}

func(np *Bgp4MplsVpnMesh) SetIngressTrafficTunnelLdp () error {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh SetIngressTrafficTunnelLdp
 return nil
}

func(np *Bgp4MplsVpnMesh) SetIngressTrafficTunnelRsvp () error {
 //parameters: MeshHandle IngressTunnel
 //AgtBgp4MplsVpnMesh SetIngressTrafficTunnelRsvp
 return nil
}

func(np *Bgp4MplsVpnMesh) GetIngressTrafficTunnelRsvp ()(string, error) {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh GetIngressTrafficTunnelRsvp
 return "", nil
}

func(np *Bgp4MplsVpnMesh) SetIngressTrafficTunnelStatic () error {
 //parameters: MeshHandle StaticLabel
 //AgtBgp4MplsVpnMesh SetIngressTrafficTunnelStatic
 return nil
}

func(np *Bgp4MplsVpnMesh) GetIngressTrafficTunnelStatic ()(string, error) {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh GetIngressTrafficTunnelStatic
 return "", nil
}

func(np *Bgp4MplsVpnMesh) SetIngressTrafficTunnelProtocols () error {
 //parameters: MeshHandle Count psaProtocols
 //AgtBgp4MplsVpnMesh SetIngressTrafficTunnelProtocols
 return nil
}

func(np *Bgp4MplsVpnMesh) GetIngressTrafficTunnelProtocols ()(string, error) {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh GetIngressTrafficTunnelProtocols
 return "", nil
}

func(np *Bgp4MplsVpnMesh) GetIngressTrafficTunnelType ()(string, error) {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh GetIngressTrafficTunnelType
 return "", nil
}

func(np *Bgp4MplsVpnMesh) GetIngressMplsTrafficTunnelLabel ()(string, error) {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh GetIngressMplsTrafficTunnelLabel
 return "", nil
}

func(np *Bgp4MplsVpnMesh) GetIngressMplsVpnLabels ()(string, error) {
 //parameters: MeshHandle AddressPoolHandle
 //AgtBgp4MplsVpnMesh GetIngressMplsVpnLabels
 return "", nil
}

func(np *Bgp4MplsVpnMesh) GetIngressMplsVpnLabelState ()(string, error) {
 //parameters: MeshHandle AddressPoolHandle
 //AgtBgp4MplsVpnMesh GetIngressMplsVpnLabelState
 return "", nil
}

func(np *Bgp4MplsVpnMesh) SetIngressIntermediateMplsLabels () error {
 //parameters: MeshHandle Count psaIntermediateLabels
 //AgtBgp4MplsVpnMesh SetIngressIntermediateMplsLabels
 return nil
}

func(np *Bgp4MplsVpnMesh) GetIngressIntermediateMplsLabels ()(string, error) {
 //parameters: MeshHandle
 //AgtBgp4MplsVpnMesh GetIngressIntermediateMplsLabels
 return "", nil
}

