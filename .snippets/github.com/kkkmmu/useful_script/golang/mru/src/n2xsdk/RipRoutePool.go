package n2xsdk

type RipRoutePool struct {
 Handler string
}

func(np *RipRoutePool) SetRoutes () error {
 //parameters: RoutePoolHandle FirstRoute PrefixLength NumRoutes Modifier
 //AgtRipRoutePool SetRoutes
 return nil
}

func(np *RipRoutePool) GetRoutes ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtRipRoutePool GetRoutes
 return "", nil
}

func(np *RipRoutePool) EnableTrafficDestinations () error {
 //parameters: RoutePoolHandle
 //AgtRipRoutePool EnableTrafficDestinations
 return nil
}

func(np *RipRoutePool) DisableTrafficDestinations () error {
 //parameters: RoutePoolHandle
 //AgtRipRoutePool DisableTrafficDestinations
 return nil
}

func(np *RipRoutePool) IsTrafficDestinationEnabled () error {
 //parameters: RoutePoolHandle
 //AgtRipRoutePool IsTrafficDestinationEnabled
 return nil
}

func(np *RipRoutePool) SetRoutePoolFlag () error {
 //parameters: RoutePoolHandle Flag
 //AgtRipRoutePool SetRoutePoolFlag
 return nil
}

func(np *RipRoutePool) UnsetRoutePoolFlag () error {
 //parameters: RoutePoolHandle Flag
 //AgtRipRoutePool UnsetRoutePoolFlag
 return nil
}

func(np *RipRoutePool) IsRoutePoolFlagSet () error {
 //parameters: RoutePoolHandle Flag
 //AgtRipRoutePool IsRoutePoolFlagSet
 return nil
}

func(np *RipRoutePool) GetRoutePoolStatus ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtRipRoutePool GetRoutePoolStatus
 return "", nil
}

func(np *RipRoutePool) GetNthRoute ()(string, error) {
 //parameters: RoutePoolHandle Index
 //AgtRipRoutePool GetNthRoute
 return "", nil
}

func(np *RipRoutePool) SetResourceClasses () error {
 //parameters: RoutePoolHandle Count psaClasses
 //AgtRipRoutePool SetResourceClasses
 return nil
}

func(np *RipRoutePool) GetResourceClasses ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtRipRoutePool GetResourceClasses
 return "", nil
}

func(np *RipRoutePool) SetMetric () error {
 //parameters: hRoutePool Metric
 //AgtRipRoutePool SetMetric
 return nil
}

func(np *RipRoutePool) GetMetric ()(string, error) {
 //parameters: hRoutePool
 //AgtRipRoutePool GetMetric
 return "", nil
}

func(np *RipRoutePool) SetRouteTag () error {
 //parameters: hRoutePool RouteTag
 //AgtRipRoutePool SetRouteTag
 return nil
}

func(np *RipRoutePool) GetRouteTag ()(string, error) {
 //parameters: hRoutePool
 //AgtRipRoutePool GetRouteTag
 return "", nil
}

func(np *RipRoutePool) SetAddressFamily () error {
 //parameters: hRoutePool AddressFamily
 //AgtRipRoutePool SetAddressFamily
 return nil
}

func(np *RipRoutePool) GetAddressFamily ()(string, error) {
 //parameters: hRoutePool
 //AgtRipRoutePool GetAddressFamily
 return "", nil
}

func(np *RipRoutePool) SetNextHop () error {
 //parameters: hRoutePool NextHop
 //AgtRipRoutePool SetNextHop
 return nil
}

func(np *RipRoutePool) GetNextHop ()(string, error) {
 //parameters: hRoutePool
 //AgtRipRoutePool GetNextHop
 return "", nil
}

func(np *RipRoutePool) Advertise () error {
 //parameters: hRoutePool
 //AgtRipRoutePool Advertise
 return nil
}

func(np *RipRoutePool) Withdraw () error {
 //parameters: hRoutePool
 //AgtRipRoutePool Withdraw
 return nil
}

func(np *RipRoutePool) AdvertisePools () error {
 //parameters: Count saRoutePools
 //AgtRipRoutePool AdvertisePools
 return nil
}

func(np *RipRoutePool) WithdrawPools () error {
 //parameters: Count saRoutePools
 //AgtRipRoutePool WithdrawPools
 return nil
}

func(np *RipRoutePool) GetSessionHandle ()(string, error) {
 //parameters: hRoutePool
 //AgtRipRoutePool GetSessionHandle
 return "", nil
}

func(np *RipRoutePool) Enable () error {
 //parameters: hRoutePool
 //AgtRipRoutePool Enable
 return nil
}

func(np *RipRoutePool) Disable () error {
 //parameters: hRoutePool
 //AgtRipRoutePool Disable
 return nil
}

func(np *RipRoutePool) IsEnabled () error {
 //parameters: hRoutePool
 //AgtRipRoutePool IsEnabled
 return nil
}

