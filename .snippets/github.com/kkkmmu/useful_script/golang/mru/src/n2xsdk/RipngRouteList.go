package n2xsdk

type RipngRouteList struct {
 Handler string
}

func(np *RipngRouteLis) AddRoutePools () error {
 //parameters: SessionHandle PfxCount saPrefixLengths RouteCount saRoutes
 //AgtRipngRouteList AddRoutePools
 return nil
}

func(np *RipngRouteLis) ListRoutePools ()(string, error) {
 //parameters: SessionHandle
 //AgtRipngRouteList ListRoutePools
 return "", nil
}

func(np *RipngRouteLis) RemoveRoutePools () error {
 //parameters: SessionHandle Count saRoutePools
 //AgtRipngRouteList RemoveRoutePools
 return nil
}

func(np *RipngRouteLis) ClearRoutePools () error {
 //parameters: SessionHandle
 //AgtRipngRouteList ClearRoutePools
 return nil
}

