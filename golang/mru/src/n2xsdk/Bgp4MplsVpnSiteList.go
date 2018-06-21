package n2xsdk

type Bgp4MplsVpnSiteList struct {
 Handler string
}

func(np *Bgp4MplsVpnSiteLis) AddSites () error {
 //parameters: SessionHandle TypeCount psaTypes AdminCount psaAdministrators NumCount psaAssignedNumbers
 //AgtBgp4MplsVpnSiteList AddSites
 return nil
}

func(np *Bgp4MplsVpnSiteLis) AddSiteRange () error {
 //parameters: SessionHandle NumSites Type FirstAdministrator AdministratorModifier FirstAssignedNumber AssignedNumberModifier
 //AgtBgp4MplsVpnSiteList AddSiteRange
 return nil
}

func(np *Bgp4MplsVpnSiteLis) RemoveSites () error {
 //parameters: SessionHandle Count psaSiteHandles
 //AgtBgp4MplsVpnSiteList RemoveSites
 return nil
}

func(np *Bgp4MplsVpnSiteLis) ListSites ()(string, error) {
 //parameters: SessionHandle
 //AgtBgp4MplsVpnSiteList ListSites
 return "", nil
}

func(np *Bgp4MplsVpnSiteLis) ClearSites () error {
 //parameters: SessionHandle
 //AgtBgp4MplsVpnSiteList ClearSites
 return nil
}

