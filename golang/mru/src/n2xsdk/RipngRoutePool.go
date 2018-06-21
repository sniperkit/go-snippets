package n2xsdk

type RipngRoutePool struct {
 Handler string
}

func(np *RipngRoutePool) SetRoutes () error {
 //parameters: RoutePoolHandle FirstRoute PrefixLength NumRoutes Modifier
 //AgtRipngRoutePool SetRoutes
 return nil
}

func(np *RipngRoutePool) GetRoutes ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtRipngRoutePool GetRoutes
 return "", nil
}

func(np *RipngRoutePool) EnableTrafficDestinations () error {
 //parameters: RoutePoolHandle
 //AgtRipngRoutePool EnableTrafficDestinations
 return nil
}

func(np *RipngRoutePool) DisableTrafficDestinations () error {
 //parameters: RoutePoolHandle
 //AgtRipngRoutePool DisableTrafficDestinations
 return nil
}

func(np *RipngRoutePool) IsTrafficDestinationEnabled () error {
 //parameters: RoutePoolHandle
 //AgtRipngRoutePool IsTrafficDestinationEnabled
 return nil
}

func(np *RipngRoutePool) SetRoutePoolFlag () error {
 //parameters: RoutePoolHandle Flag
 //AgtRipngRoutePool SetRoutePoolFlag
 return nil
}

func(np *RipngRoutePool) UnsetRoutePoolFlag () error {
 //parameters: RoutePoolHandle Flag
 //AgtRipngRoutePool UnsetRoutePoolFlag
 return nil
}

func(np *RipngRoutePool) IsRoutePoolFlagSet () error {
 //parameters: RoutePoolHandle Flag
 //AgtRipngRoutePool IsRoutePoolFlagSet
 return nil
}

func(np *RipngRoutePool) GetRoutePoolStatus ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtRipngRoutePool GetRoutePoolStatus
 return "", nil
}

func(np *RipngRoutePool) GetNthRoute ()(string, error) {
 //parameters: RoutePoolHandle Index
 //AgtRipngRoutePool GetNthRoute
 return "", nil
}

func(np *RipngRoutePool) SetResourceClasses () error {
 //parameters: RoutePoolHandle Count psaClasses
 //AgtRipngRoutePool SetResourceClasses
 return nil
}

func(np *RipngRoutePool) GetResourceClasses ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtRipngRoutePool GetResourceClasses
 return "", nil
}

func(np *RipngRoutePool) SetMetric () error {
 //parameters: hRoutePool Metric
 //AgtRipngRoutePool SetMetric
 return nil
}

func(np *RipngRoutePool) GetMetric ()(string, error) {
 //parameters: hRoutePool
 //AgtRipngRoutePool GetMetric
 return "", nil
}

func(np *RipngRoutePool) SetRouteTag () error {
 //parameters: hRoutePool RouteTag
 //AgtRipngRoutePool SetRouteTag
 return nil
}

func(np *RipngRoutePool) GetRouteTag ()(string, error) {
 //parameters: hRoutePool
 //AgtRipngRoutePool GetRouteTag
 return "", nil
}

func(np *RipngRoutePool) SetAddressFamily () error {
 //parameters: hRoutePool AddressFamily
 //AgtRipngRoutePool SetAddressFamily
 return nil
}

func(np *RipngRoutePool) GetAddressFamily ()(string, error) {
 //parameters: hRoutePool
 //AgtRipngRoutePool GetAddressFamily
 return "", nil
}

func(np *RipngRoutePool) SetNextHop () error {
 //parameters: hRoutePool NextHop
 //AgtRipngRoutePool SetNextHop
 return nil
}

func(np *RipngRoutePool) GetNextHop ()(string, error) {
 //parameters: hRoutePool
 //AgtRipngRoutePool GetNextHop
 return "", nil
}

func(np *RipngRoutePool) Advertise () error {
 //parameters: hRoutePool
 //AgtRipngRoutePool Advertise
 return nil
}

func(np *RipngRoutePool) Withdraw () error {
 //parameters: hRoutePool
 //AgtRipngRoutePool Withdraw
 return nil
}

func(np *RipngRoutePool) AdvertisePools () error {
 //parameters: Count saRoutePools
 //AgtRipngRoutePool AdvertisePools
 return nil
}

func(np *RipngRoutePool) WithdrawPools () error {
 //parameters: Count saRoutePools
 //AgtRipngRoutePool WithdrawPools
 return nil
}

func(np *RipngRoutePool) GetSessionHandle ()(string, error) {
 //parameters: hRoutePool
 //AgtRipngRoutePool GetSessionHandle
 return "", nil
}

func(np *RipngRoutePool) Enable () error {
 //parameters: hRoutePool
 //AgtRipngRoutePool Enable
 return nil
}

func(np *RipngRoutePool) Disable () error {
 //parameters: hRoutePool
 //AgtRipngRoutePool Disable
 return nil
}

func(np *RipngRoutePool) IsEnabled () error {
 //parameters: hRoutePool
 //AgtRipngRoutePool IsEnabled
 return nil
}

func(np *RipngRoutePool) SelectNextHop () error {
 //parameters: hRoutePool
 //AgtRipngRoutePool SelectNextHop
 return nil
}

func(np *RipngRoutePool) DeselectNextHop () error {
 //parameters: hRoutePool
 //AgtRipngRoutePool DeselectNextHop
 return nil
}

func(np *RipngRoutePool) IsNextHopSelected () error {
 //parameters: hRoutePool
 //AgtRipngRoutePool IsNextHopSelected
 return nil
}

