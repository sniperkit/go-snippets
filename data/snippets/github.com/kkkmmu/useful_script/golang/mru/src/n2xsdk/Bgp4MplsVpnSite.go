package n2xsdk

type Bgp4MplsVpnSite struct {
 Handler string
}

func(np *Bgp4MplsVpnSite) SetRouteDistinguisher () error {
 //parameters: SiteHandle Type Administrator AssignedNumber
 //AgtBgp4MplsVpnSite SetRouteDistinguisher
 return nil
}

func(np *Bgp4MplsVpnSite) GetRouteDistinguisher ()(string, error) {
 //parameters: SiteHandle
 //AgtBgp4MplsVpnSite GetRouteDistinguisher
 return "", nil
}

func(np *Bgp4MplsVpnSite) SetName () error {
 //parameters: SiteHandle Name
 //AgtBgp4MplsVpnSite SetName
 return nil
}

func(np *Bgp4MplsVpnSite) GetName ()(string, error) {
 //parameters: SiteHandle
 //AgtBgp4MplsVpnSite GetName
 return "", nil
}

func(np *Bgp4MplsVpnSite) AddVpns () error {
 //parameters: SiteHandle Count psaVpnHandles
 //AgtBgp4MplsVpnSite AddVpns
 return nil
}

func(np *Bgp4MplsVpnSite) RemoveVpns () error {
 //parameters: SiteHandle Count psaVpnHandles
 //AgtBgp4MplsVpnSite RemoveVpns
 return nil
}

func(np *Bgp4MplsVpnSite) ListVpns ()(string, error) {
 //parameters: SiteHandle
 //AgtBgp4MplsVpnSite ListVpns
 return "", nil
}

func(np *Bgp4MplsVpnSite) ClearVpns () error {
 //parameters: SiteHandle
 //AgtBgp4MplsVpnSite ClearVpns
 return nil
}

func(np *Bgp4MplsVpnSite) ListImportTargets ()(string, error) {
 //parameters: SiteHandle
 //AgtBgp4MplsVpnSite ListImportTargets
 return "", nil
}

func(np *Bgp4MplsVpnSite) GetLabel ()(string, error) {
 //parameters: SiteHandle
 //AgtBgp4MplsVpnSite GetLabel
 return "", nil
}

func(np *Bgp4MplsVpnSite) GetVrfTable ()(string, error) {
 //parameters: SiteCount psaSiteHandle
 //AgtBgp4MplsVpnSite GetVrfTable
 return "", nil
}

func(np *Bgp4MplsVpnSite) QueryVrfEntry () error {
 //parameters: SiteCount psaSiteHandle Prefix PrefixLength
 //AgtBgp4MplsVpnSite QueryVrfEntry
 return nil
}

func(np *Bgp4MplsVpnSite) QueryLspLabel () error {
 //parameters: SiteCount psaSiteHandle Destination
 //AgtBgp4MplsVpnSite QueryLspLabel
 return nil
}

func(np *Bgp4MplsVpnSite) Enable () error {
 //parameters: SiteHandle
 //AgtBgp4MplsVpnSite Enable
 return nil
}

func(np *Bgp4MplsVpnSite) Disable () error {
 //parameters: SiteHandle
 //AgtBgp4MplsVpnSite Disable
 return nil
}

func(np *Bgp4MplsVpnSite) IsEnabled () error {
 //parameters: SiteHandle
 //AgtBgp4MplsVpnSite IsEnabled
 return nil
}

