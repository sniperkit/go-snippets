package n2xsdk

type Bgp4RouteList struct {
 Handler string
}

func(np *Bgp4RouteLis) AddRoutePools () error {
 //parameters: SessionHandle PfxCount saPrefixLengths RouteCount saRoutes
 //AgtBgp4RouteList AddRoutePools
 return nil
}

func(np *Bgp4RouteLis) ListRoutePools ()(string, error) {
 //parameters: SessionHandle
 //AgtBgp4RouteList ListRoutePools
 return "", nil
}

func(np *Bgp4RouteLis) RemoveRoutePools () error {
 //parameters: SessionHandle Count saRoutePools
 //AgtBgp4RouteList RemoveRoutePools
 return nil
}

func(np *Bgp4RouteLis) ClearRoutePools () error {
 //parameters: SessionHandle
 //AgtBgp4RouteList ClearRoutePools
 return nil
}

