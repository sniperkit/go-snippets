package n2xsdk

type IsisIpv4RoutePool struct {
 Handler string
}

func(np *IsisIpv4RoutePool) SetRoutes () error {
 //parameters: RoutePoolHandle FirstRoute PrefixLength NumRoutes Modifier
 //AgtIsisIpv4RoutePool SetRoutes
 return nil
}

func(np *IsisIpv4RoutePool) GetRoutes ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtIsisIpv4RoutePool GetRoutes
 return "", nil
}

func(np *IsisIpv4RoutePool) EnableTrafficDestinations () error {
 //parameters: RoutePoolHandle
 //AgtIsisIpv4RoutePool EnableTrafficDestinations
 return nil
}

func(np *IsisIpv4RoutePool) DisableTrafficDestinations () error {
 //parameters: RoutePoolHandle
 //AgtIsisIpv4RoutePool DisableTrafficDestinations
 return nil
}

func(np *IsisIpv4RoutePool) IsTrafficDestinationEnabled () error {
 //parameters: RoutePoolHandle
 //AgtIsisIpv4RoutePool IsTrafficDestinationEnabled
 return nil
}

func(np *IsisIpv4RoutePool) SetRoutePoolFlag () error {
 //parameters: RoutePoolHandle Flag
 //AgtIsisIpv4RoutePool SetRoutePoolFlag
 return nil
}

func(np *IsisIpv4RoutePool) UnsetRoutePoolFlag () error {
 //parameters: RoutePoolHandle Flag
 //AgtIsisIpv4RoutePool UnsetRoutePoolFlag
 return nil
}

func(np *IsisIpv4RoutePool) IsRoutePoolFlagSet () error {
 //parameters: RoutePoolHandle Flag
 //AgtIsisIpv4RoutePool IsRoutePoolFlagSet
 return nil
}

func(np *IsisIpv4RoutePool) GetRoutePoolStatus ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtIsisIpv4RoutePool GetRoutePoolStatus
 return "", nil
}

func(np *IsisIpv4RoutePool) GetNthRoute ()(string, error) {
 //parameters: RoutePoolHandle Index
 //AgtIsisIpv4RoutePool GetNthRoute
 return "", nil
}

func(np *IsisIpv4RoutePool) SetResourceClasses () error {
 //parameters: RoutePoolHandle Count psaClasses
 //AgtIsisIpv4RoutePool SetResourceClasses
 return nil
}

func(np *IsisIpv4RoutePool) GetResourceClasses ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtIsisIpv4RoutePool GetResourceClasses
 return "", nil
}

func(np *IsisIpv4RoutePool) GetReachablePrefixCount ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtIsisIpv4RoutePool GetReachablePrefixCount
 return "", nil
}

func(np *IsisIpv4RoutePool) EnableRouteTagging () error {
 //parameters: RoutePoolHandle
 //AgtIsisIpv4RoutePool EnableRouteTagging
 return nil
}

func(np *IsisIpv4RoutePool) DisableRouteTagging () error {
 //parameters: RoutePoolHandle
 //AgtIsisIpv4RoutePool DisableRouteTagging
 return nil
}

func(np *IsisIpv4RoutePool) IsRouteTaggingEnabled () error {
 //parameters: RoutePoolHandle
 //AgtIsisIpv4RoutePool IsRouteTaggingEnabled
 return nil
}

func(np *IsisIpv4RoutePool) SetRouteTag () error {
 //parameters: RoutePoolHandle RouteTagValue
 //AgtIsisIpv4RoutePool SetRouteTag
 return nil
}

func(np *IsisIpv4RoutePool) GetRouteTag ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtIsisIpv4RoutePool GetRouteTag
 return "", nil
}

