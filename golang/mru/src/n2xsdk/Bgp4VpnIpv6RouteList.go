package n2xsdk

type Bgp4VpnIpv6RouteList struct {
 Handler string
}

func(np *Bgp4VpnIpv6RouteLis) AddRoutePools () error {
 //parameters: SiteHandle PfxCount saPrefixLengths RouteCount saNumRoutes
 //AgtBgp4VpnIpv6RouteList AddRoutePools
 return nil
}

func(np *Bgp4VpnIpv6RouteLis) ListRoutePools ()(string, error) {
 //parameters: SiteHandle
 //AgtBgp4VpnIpv6RouteList ListRoutePools
 return "", nil
}

func(np *Bgp4VpnIpv6RouteLis) RemoveRoutePools () error {
 //parameters: SiteHandle Count saRoutePools
 //AgtBgp4VpnIpv6RouteList RemoveRoutePools
 return nil
}

func(np *Bgp4VpnIpv6RouteLis) ClearRoutePools () error {
 //parameters: SiteHandle
 //AgtBgp4VpnIpv6RouteList ClearRoutePools
 return nil
}

