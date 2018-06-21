package n2xsdk

type Bgp4VpnIpv6RoutePool struct {
 Handler string
}

func(np *Bgp4VpnIpv6RoutePool) SetRoutes () error {
 //parameters: RoutePoolHandle FirstRoute PrefixLength NumRoutes Modifier
 //AgtBgp4VpnIpv6RoutePool SetRoutes
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) SetNextHop () error {
 //parameters: RoutePoolHandle NextHop
 //AgtBgp4VpnIpv6RoutePool SetNextHop
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) GetNextHop ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool GetNextHop
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) SetNextHops () error {
 //parameters: RoutePoolHandle NextHops
 //AgtBgp4VpnIpv6RoutePool SetNextHops
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) GetNextHops ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool GetNextHops
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) SetOrigin () error {
 //parameters: RoutePoolHandle Origin
 //AgtBgp4VpnIpv6RoutePool SetOrigin
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) GetOrigin ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool GetOrigin
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) SetCommunities () error {
 //parameters: RoutePoolHandle CommunitiesList
 //AgtBgp4VpnIpv6RoutePool SetCommunities
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) GetCommunities ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool GetCommunities
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) SetPathAttributeString () error {
 //parameters: RoutePoolHandle PathAttribute PathAttributeString
 //AgtBgp4VpnIpv6RoutePool SetPathAttributeString
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) GetPathAttributeString ()(string, error) {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4VpnIpv6RoutePool GetPathAttributeString
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) SetPathAttributeLong () error {
 //parameters: RoutePoolHandle PathAttribute PathAttributeLong
 //AgtBgp4VpnIpv6RoutePool SetPathAttributeLong
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) GetPathAttributeLong ()(string, error) {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4VpnIpv6RoutePool GetPathAttributeLong
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) IsPathAttributeSelected () error {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4VpnIpv6RoutePool IsPathAttributeSelected
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) IsPathAttributeSet () error {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4VpnIpv6RoutePool IsPathAttributeSet
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) SetAsPath () error {
 //parameters: RoutePoolHandle AsPathAttributeType AsNumber
 //AgtBgp4VpnIpv6RoutePool SetAsPath
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) GetAsPath ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool GetAsPath
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) GetAsPathByType ()(string, error) {
 //parameters: RoutePoolHandle AsPathAttributeType
 //AgtBgp4VpnIpv6RoutePool GetAsPathByType
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) GetMaximumAsPathLength ()(string, error) {
 //parameters: 
 //AgtBgp4VpnIpv6RoutePool GetMaximumAsPathLength
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) SetAggregator () error {
 //parameters: RoutePoolHandle AsNumber IpAddress
 //AgtBgp4VpnIpv6RoutePool SetAggregator
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) GetAggregator ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool GetAggregator
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) SetOriginatorId () error {
 //parameters: RoutePoolHandle OriginatorId
 //AgtBgp4VpnIpv6RoutePool SetOriginatorId
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) GetOriginatorId ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool GetOriginatorId
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) SetClusterList ()(string, error) {
 //parameters: RoutePoolHandle ClusterList
 //AgtBgp4VpnIpv6RoutePool SetClusterList
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) GetClusterList ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool GetClusterList
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) Enable () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool Enable
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) Disable () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool Disable
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) IsEnabled () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool IsEnabled
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) Advertise () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool Advertise
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) AdvertisePools () error {
 //parameters: RoutePoolHandles RoutesPerUpdate InterUpdateDelay
 //AgtBgp4VpnIpv6RoutePool AdvertisePools
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) Withdraw () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool Withdraw
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) WithdrawPools () error {
 //parameters: RoutePoolHandles RoutesPerUpdate InterUpdateDelay
 //AgtBgp4VpnIpv6RoutePool WithdrawPools
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) GetState ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool GetState
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) SelectPathAttribute () error {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4VpnIpv6RoutePool SelectPathAttribute
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) DeselectPathAttribute () error {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4VpnIpv6RoutePool DeselectPathAttribute
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) EnableTrafficDestinations () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool EnableTrafficDestinations
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) DisableTrafficDestinations () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool DisableTrafficDestinations
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) EnableTraffic () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool EnableTraffic
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) DisableTraffic () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool DisableTraffic
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) SetRoutePoolFlag () error {
 //parameters: RoutePoolHandle RoutePoolFlag
 //AgtBgp4VpnIpv6RoutePool SetRoutePoolFlag
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) UnsetRoutePoolFlag () error {
 //parameters: RoutePoolHandle RoutePoolFlag
 //AgtBgp4VpnIpv6RoutePool UnsetRoutePoolFlag
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) SetResourceClasses () error {
 //parameters: RoutePoolHandle ResourceClasses
 //AgtBgp4VpnIpv6RoutePool SetResourceClasses
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) GetNumAdvertisedRoutes ()(string, error) {
 //parameters: RoutePoolHandle NumRoutes
 //AgtBgp4VpnIpv6RoutePool GetNumAdvertisedRoutes
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) GetNumWithdrawnRoutes ()(string, error) {
 //parameters: RoutePoolHandle NumRoutes
 //AgtBgp4VpnIpv6RoutePool GetNumWithdrawnRoutes
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) SetExportTarget () error {
 //parameters: RoutePoolHandle RouteTargetType Administrator AssignedNumber
 //AgtBgp4VpnIpv6RoutePool SetExportTarget
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) GetExportTarget ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool GetExportTarget
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) GetRouteDistinguisher ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool GetRouteDistinguisher
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) GetSiteHandle ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool GetSiteHandle
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) GetRoutes ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool GetRoutes
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) IsTrafficDestinationEnabled () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool IsTrafficDestinationEnabled
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) IsRoutePoolFlagSet () error {
 //parameters: RoutePoolHandle Flag
 //AgtBgp4VpnIpv6RoutePool IsRoutePoolFlagSet
 return nil
}

func(np *Bgp4VpnIpv6RoutePool) GetRoutePoolStatus ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool GetRoutePoolStatus
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) GetNthRoute ()(string, error) {
 //parameters: RoutePoolHandle Index
 //AgtBgp4VpnIpv6RoutePool GetNthRoute
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) GetResourceClasses ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv6RoutePool GetResourceClasses
 return "", nil
}

func(np *Bgp4VpnIpv6RoutePool) GetSessionHandle ()(string, error) {
 //parameters: hRoutePool
 //AgtBgp4VpnIpv6RoutePool GetSessionHandle
 return "", nil
}

