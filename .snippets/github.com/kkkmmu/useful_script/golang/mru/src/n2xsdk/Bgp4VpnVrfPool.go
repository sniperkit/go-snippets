package n2xsdk

type Bgp4VpnVrfPool struct {
 Handler string
}

func(np *Bgp4VpnVrfPool) SetVrfCreationMode () error {
 //parameters: VrfPoolHandle VrfCreationMode
 //AgtBgp4VpnVrfPool SetVrfCreationMode
 return nil
}

func(np *Bgp4VpnVrfPool) GetVrfCreationMode ()(string, error) {
 //parameters: VrfPoolHandle
 //AgtBgp4VpnVrfPool GetVrfCreationMode
 return "", nil
}

func(np *Bgp4VpnVrfPool) GetVrfTable ()(string, error) {
 //parameters: VrfPoolHandle VrfIndex
 //AgtBgp4VpnVrfPool GetVrfTable
 return "", nil
}

func(np *Bgp4VpnVrfPool) QueryVrfEntry () error {
 //parameters: VrfPoolHandle VrfIndex Prefix PrefixLength
 //AgtBgp4VpnVrfPool QueryVrfEntry
 return nil
}

func(np *Bgp4VpnVrfPool) QueryLspLabel () error {
 //parameters: VrfPoolHandle VrfIndex Destination
 //AgtBgp4VpnVrfPool QueryLspLabel
 return nil
}

func(np *Bgp4VpnVrfPool) SetImportRouteTargetList ()(string, error) {
 //parameters: VrfPoolHandle RouteTargetType RouteTargetList
 //AgtBgp4VpnVrfPool SetImportRouteTargetList
 return "", nil
}

func(np *Bgp4VpnVrfPool) GetImportRouteTargetList ()(string, error) {
 //parameters: VrfPoolHandle RouteTargetType
 //AgtBgp4VpnVrfPool GetImportRouteTargetList
 return "", nil
}

func(np *Bgp4VpnVrfPool) SetImportRouteTarget () error {
 //parameters: VrfPoolHandle RouteTargetType RouteTarget
 //AgtBgp4VpnVrfPool SetImportRouteTarget
 return nil
}

func(np *Bgp4VpnVrfPool) GetImportRouteTarget ()(string, error) {
 //parameters: VrfPoolHandle RouteTargetType
 //AgtBgp4VpnVrfPool GetImportRouteTarget
 return "", nil
}

func(np *Bgp4VpnVrfPool) Enable () error {
 //parameters: VrfPoolHandle
 //AgtBgp4VpnVrfPool Enable
 return nil
}

func(np *Bgp4VpnVrfPool) Disable () error {
 //parameters: VrfPoolHandle
 //AgtBgp4VpnVrfPool Disable
 return nil
}

func(np *Bgp4VpnVrfPool) IsEnabled () error {
 //parameters: VrfPoolHandle
 //AgtBgp4VpnVrfPool IsEnabled
 return nil
}

func(np *Bgp4VpnVrfPool) SetImportRouteTargetIncrementingRange () error {
 //parameters: VrfPoolHandle RouteTargetType FirstRouteTarget RouteTargetIncrement RouteTargetRepeat PercentOverlap
 //AgtBgp4VpnVrfPool SetImportRouteTargetIncrementingRange
 return nil
}

func(np *Bgp4VpnVrfPool) GetImportRouteTargetIncrementingRange ()(string, error) {
 //parameters: VrfPoolHandle RouteTargetType
 //AgtBgp4VpnVrfPool GetImportRouteTargetIncrementingRange
 return "", nil
}

func(np *Bgp4VpnVrfPool) SetImportRouteTargetType () error {
 //parameters: VrfPoolHandle RouteTargetTypeValue
 //AgtBgp4VpnVrfPool SetImportRouteTargetType
 return nil
}

func(np *Bgp4VpnVrfPool) GetImportRouteTargetType ()(string, error) {
 //parameters: VrfPoolHandle
 //AgtBgp4VpnVrfPool GetImportRouteTargetType
 return "", nil
}

func(np *Bgp4VpnVrfPool) SetVpnsPerPeer () error {
 //parameters: VrfPoolHandle NumVpnsPerPeer
 //AgtBgp4VpnVrfPool SetVpnsPerPeer
 return nil
}

func(np *Bgp4VpnVrfPool) GetVpnsPerPeer ()(string, error) {
 //parameters: VrfPoolHandle
 //AgtBgp4VpnVrfPool GetVpnsPerPeer
 return "", nil
}

func(np *Bgp4VpnVrfPool) GetTotalNumberOfVpns ()(string, error) {
 //parameters: VrfPoolHandle
 //AgtBgp4VpnVrfPool GetTotalNumberOfVpns
 return "", nil
}

