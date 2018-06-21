package n2xsdk

type Bgp4Ipv6RoutePool struct {
 Handler string
}

func(np *Bgp4Ipv6RoutePool) SetRoutes () error {
 //parameters: RoutePoolHandle FirstRoute PrefixLength NumRoutes Modifier
 //AgtBgp4Ipv6RoutePool SetRoutes
 return nil
}

func(np *Bgp4Ipv6RoutePool) SetNextHop () error {
 //parameters: RoutePoolHandle NextHop
 //AgtBgp4Ipv6RoutePool SetNextHop
 return nil
}

func(np *Bgp4Ipv6RoutePool) GetNextHop ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool GetNextHop
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) SetNextHops () error {
 //parameters: RoutePoolHandle NextHops
 //AgtBgp4Ipv6RoutePool SetNextHops
 return nil
}

func(np *Bgp4Ipv6RoutePool) GetNextHops ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool GetNextHops
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) SetOrigin () error {
 //parameters: RoutePoolHandle Origin
 //AgtBgp4Ipv6RoutePool SetOrigin
 return nil
}

func(np *Bgp4Ipv6RoutePool) GetOrigin ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool GetOrigin
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) SetCommunities () error {
 //parameters: RoutePoolHandle CommunitiesList
 //AgtBgp4Ipv6RoutePool SetCommunities
 return nil
}

func(np *Bgp4Ipv6RoutePool) GetCommunities ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool GetCommunities
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) SetPathAttributeString () error {
 //parameters: RoutePoolHandle PathAttribute PathAttributeString
 //AgtBgp4Ipv6RoutePool SetPathAttributeString
 return nil
}

func(np *Bgp4Ipv6RoutePool) GetPathAttributeString ()(string, error) {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4Ipv6RoutePool GetPathAttributeString
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) SetPathAttributeLong () error {
 //parameters: RoutePoolHandle PathAttribute PathAttributeLong
 //AgtBgp4Ipv6RoutePool SetPathAttributeLong
 return nil
}

func(np *Bgp4Ipv6RoutePool) GetPathAttributeLong ()(string, error) {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4Ipv6RoutePool GetPathAttributeLong
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) IsPathAttributeSelected () error {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4Ipv6RoutePool IsPathAttributeSelected
 return nil
}

func(np *Bgp4Ipv6RoutePool) IsPathAttributeSet () error {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4Ipv6RoutePool IsPathAttributeSet
 return nil
}

func(np *Bgp4Ipv6RoutePool) SetAsPath () error {
 //parameters: RoutePoolHandle AsPathAttributeType AsNumber
 //AgtBgp4Ipv6RoutePool SetAsPath
 return nil
}

func(np *Bgp4Ipv6RoutePool) GetAsPath ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool GetAsPath
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) GetAsPathByType ()(string, error) {
 //parameters: RoutePoolHandle AsPathAttributeType
 //AgtBgp4Ipv6RoutePool GetAsPathByType
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) GetMaximumAsPathLength ()(string, error) {
 //parameters: 
 //AgtBgp4Ipv6RoutePool GetMaximumAsPathLength
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) SetAggregator () error {
 //parameters: RoutePoolHandle AsNumber IpAddress
 //AgtBgp4Ipv6RoutePool SetAggregator
 return nil
}

func(np *Bgp4Ipv6RoutePool) GetAggregator ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool GetAggregator
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) SetOriginatorId () error {
 //parameters: RoutePoolHandle OriginatorId
 //AgtBgp4Ipv6RoutePool SetOriginatorId
 return nil
}

func(np *Bgp4Ipv6RoutePool) GetOriginatorId ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool GetOriginatorId
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) SetClusterList ()(string, error) {
 //parameters: RoutePoolHandle ClusterList
 //AgtBgp4Ipv6RoutePool SetClusterList
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) GetClusterList ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool GetClusterList
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) Enable () error {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool Enable
 return nil
}

func(np *Bgp4Ipv6RoutePool) Disable () error {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool Disable
 return nil
}

func(np *Bgp4Ipv6RoutePool) IsEnabled () error {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool IsEnabled
 return nil
}

func(np *Bgp4Ipv6RoutePool) Advertise () error {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool Advertise
 return nil
}

func(np *Bgp4Ipv6RoutePool) AdvertisePools () error {
 //parameters: RoutePoolHandles RoutesPerUpdate InterUpdateDelay
 //AgtBgp4Ipv6RoutePool AdvertisePools
 return nil
}

func(np *Bgp4Ipv6RoutePool) Withdraw () error {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool Withdraw
 return nil
}

func(np *Bgp4Ipv6RoutePool) WithdrawPools () error {
 //parameters: RoutePoolHandles RoutesPerUpdate InterUpdateDelay
 //AgtBgp4Ipv6RoutePool WithdrawPools
 return nil
}

func(np *Bgp4Ipv6RoutePool) GetState ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool GetState
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) SelectPathAttribute () error {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4Ipv6RoutePool SelectPathAttribute
 return nil
}

func(np *Bgp4Ipv6RoutePool) DeselectPathAttribute () error {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4Ipv6RoutePool DeselectPathAttribute
 return nil
}

func(np *Bgp4Ipv6RoutePool) EnableTrafficDestinations () error {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool EnableTrafficDestinations
 return nil
}

func(np *Bgp4Ipv6RoutePool) DisableTrafficDestinations () error {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool DisableTrafficDestinations
 return nil
}

func(np *Bgp4Ipv6RoutePool) EnableTraffic () error {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool EnableTraffic
 return nil
}

func(np *Bgp4Ipv6RoutePool) DisableTraffic () error {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool DisableTraffic
 return nil
}

func(np *Bgp4Ipv6RoutePool) SetRoutePoolFlag () error {
 //parameters: RoutePoolHandle RoutePoolFlag
 //AgtBgp4Ipv6RoutePool SetRoutePoolFlag
 return nil
}

func(np *Bgp4Ipv6RoutePool) UnsetRoutePoolFlag () error {
 //parameters: RoutePoolHandle RoutePoolFlag
 //AgtBgp4Ipv6RoutePool UnsetRoutePoolFlag
 return nil
}

func(np *Bgp4Ipv6RoutePool) SetResourceClasses () error {
 //parameters: RoutePoolHandle ResourceClasses
 //AgtBgp4Ipv6RoutePool SetResourceClasses
 return nil
}

func(np *Bgp4Ipv6RoutePool) GetNumAdvertisedRoutes ()(string, error) {
 //parameters: RoutePoolHandle NumRoutes
 //AgtBgp4Ipv6RoutePool GetNumAdvertisedRoutes
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) GetNumWithdrawnRoutes ()(string, error) {
 //parameters: RoutePoolHandle NumRoutes
 //AgtBgp4Ipv6RoutePool GetNumWithdrawnRoutes
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) GetRoutes ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool GetRoutes
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) IsTrafficDestinationEnabled () error {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool IsTrafficDestinationEnabled
 return nil
}

func(np *Bgp4Ipv6RoutePool) IsRoutePoolFlagSet () error {
 //parameters: RoutePoolHandle Flag
 //AgtBgp4Ipv6RoutePool IsRoutePoolFlagSet
 return nil
}

func(np *Bgp4Ipv6RoutePool) GetRoutePoolStatus ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool GetRoutePoolStatus
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) GetNthRoute ()(string, error) {
 //parameters: RoutePoolHandle Index
 //AgtBgp4Ipv6RoutePool GetNthRoute
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) GetResourceClasses ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4Ipv6RoutePool GetResourceClasses
 return "", nil
}

func(np *Bgp4Ipv6RoutePool) GetSessionHandle ()(string, error) {
 //parameters: hRoutePool
 //AgtBgp4Ipv6RoutePool GetSessionHandle
 return "", nil
}

