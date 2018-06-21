package n2xsdk

type Bgp4VpnIpv4RouteList struct {
 Handler string
}

func(np *Bgp4VpnIpv4RouteLis) AddRoutePools () error {
 //parameters: SiteHandle PfxCount saPrefixLengths RouteCount saRoutes
 //AgtBgp4VpnIpv4RouteList AddRoutePools
 return nil
}

func(np *Bgp4VpnIpv4RouteLis) ListRoutePools ()(string, error) {
 //parameters: SiteHandle
 //AgtBgp4VpnIpv4RouteList ListRoutePools
 return "", nil
}

func(np *Bgp4VpnIpv4RouteLis) RemoveRoutePools () error {
 //parameters: SiteHandle Count saRoutePools
 //AgtBgp4VpnIpv4RouteList RemoveRoutePools
 return nil
}

func(np *Bgp4VpnIpv4RouteLis) ClearRoutePools () error {
 //parameters: SiteHandle
 //AgtBgp4VpnIpv4RouteList ClearRoutePools
 return nil
}

