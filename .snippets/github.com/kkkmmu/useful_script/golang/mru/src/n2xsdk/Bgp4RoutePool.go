package n2xsdk

type Bgp4RoutePool struct {
 Handler string
}

func(np *Bgp4RoutePool) SetRoutes () error {
 //parameters: RoutePoolHandle FirstRoute PrefixLength NumRoutes Modifier
 //AgtBgp4RoutePool SetRoutes
 return nil
}

func(np *Bgp4RoutePool) SetNextHop () error {
 //parameters: RoutePoolHandle NextHop
 //AgtBgp4RoutePool SetNextHop
 return nil
}

func(np *Bgp4RoutePool) GetNextHop ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool GetNextHop
 return "", nil
}

func(np *Bgp4RoutePool) SetNextHops () error {
 //parameters: RoutePoolHandle NextHops
 //AgtBgp4RoutePool SetNextHops
 return nil
}

func(np *Bgp4RoutePool) SetOrigin () error {
 //parameters: RoutePoolHandle Origin
 //AgtBgp4RoutePool SetOrigin
 return nil
}

func(np *Bgp4RoutePool) GetOrigin ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool GetOrigin
 return "", nil
}

func(np *Bgp4RoutePool) SetCommunities () error {
 //parameters: RoutePoolHandle CommunitiesList
 //AgtBgp4RoutePool SetCommunities
 return nil
}

func(np *Bgp4RoutePool) GetCommunities ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool GetCommunities
 return "", nil
}

func(np *Bgp4RoutePool) SetPathAttributeString () error {
 //parameters: RoutePoolHandle PathAttribute PathAttributeString
 //AgtBgp4RoutePool SetPathAttributeString
 return nil
}

func(np *Bgp4RoutePool) GetPathAttributeString ()(string, error) {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4RoutePool GetPathAttributeString
 return "", nil
}

func(np *Bgp4RoutePool) SetPathAttributeLong () error {
 //parameters: RoutePoolHandle PathAttribute PathAttributeLong
 //AgtBgp4RoutePool SetPathAttributeLong
 return nil
}

func(np *Bgp4RoutePool) GetPathAttributeLong ()(string, error) {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4RoutePool GetPathAttributeLong
 return "", nil
}

func(np *Bgp4RoutePool) IsPathAttributeSelected () error {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4RoutePool IsPathAttributeSelected
 return nil
}

func(np *Bgp4RoutePool) IsPathAttributeSet () error {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4RoutePool IsPathAttributeSet
 return nil
}

func(np *Bgp4RoutePool) SetAsPath () error {
 //parameters: RoutePoolHandle AsPathAttributeType AsNumber
 //AgtBgp4RoutePool SetAsPath
 return nil
}

func(np *Bgp4RoutePool) GetAsPath ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool GetAsPath
 return "", nil
}

func(np *Bgp4RoutePool) GetAsPathByType ()(string, error) {
 //parameters: RoutePoolHandle AsPathAttributeType
 //AgtBgp4RoutePool GetAsPathByType
 return "", nil
}

func(np *Bgp4RoutePool) GetMaximumAsPathLength ()(string, error) {
 //parameters: 
 //AgtBgp4RoutePool GetMaximumAsPathLength
 return "", nil
}

func(np *Bgp4RoutePool) SetAggregator () error {
 //parameters: RoutePoolHandle AsNumber IpAddress
 //AgtBgp4RoutePool SetAggregator
 return nil
}

func(np *Bgp4RoutePool) GetAggregator ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool GetAggregator
 return "", nil
}

func(np *Bgp4RoutePool) SetOriginatorId () error {
 //parameters: RoutePoolHandle OriginatorId
 //AgtBgp4RoutePool SetOriginatorId
 return nil
}

func(np *Bgp4RoutePool) GetOriginatorId ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool GetOriginatorId
 return "", nil
}

func(np *Bgp4RoutePool) SetClusterList ()(string, error) {
 //parameters: RoutePoolHandle ClusterList
 //AgtBgp4RoutePool SetClusterList
 return "", nil
}

func(np *Bgp4RoutePool) GetClusterList ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool GetClusterList
 return "", nil
}

func(np *Bgp4RoutePool) Enable () error {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool Enable
 return nil
}

func(np *Bgp4RoutePool) Disable () error {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool Disable
 return nil
}

func(np *Bgp4RoutePool) IsEnabled () error {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool IsEnabled
 return nil
}

func(np *Bgp4RoutePool) Advertise () error {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool Advertise
 return nil
}

func(np *Bgp4RoutePool) AdvertisePools () error {
 //parameters: RoutePoolHandles RoutesPerUpdate InterUpdateDelay
 //AgtBgp4RoutePool AdvertisePools
 return nil
}

func(np *Bgp4RoutePool) Withdraw () error {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool Withdraw
 return nil
}

func(np *Bgp4RoutePool) WithdrawPools () error {
 //parameters: RoutePoolHandles RoutesPerUpdate InterUpdateDelay
 //AgtBgp4RoutePool WithdrawPools
 return nil
}

func(np *Bgp4RoutePool) GetState ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool GetState
 return "", nil
}

func(np *Bgp4RoutePool) SelectPathAttribute () error {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4RoutePool SelectPathAttribute
 return nil
}

func(np *Bgp4RoutePool) DeselectPathAttribute () error {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4RoutePool DeselectPathAttribute
 return nil
}

func(np *Bgp4RoutePool) EnableTrafficDestinations () error {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool EnableTrafficDestinations
 return nil
}

func(np *Bgp4RoutePool) DisableTrafficDestinations () error {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool DisableTrafficDestinations
 return nil
}

func(np *Bgp4RoutePool) EnableTraffic () error {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool EnableTraffic
 return nil
}

func(np *Bgp4RoutePool) DisableTraffic () error {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool DisableTraffic
 return nil
}

func(np *Bgp4RoutePool) SetRoutePoolFlag () error {
 //parameters: RoutePoolHandle RoutePoolFlag
 //AgtBgp4RoutePool SetRoutePoolFlag
 return nil
}

func(np *Bgp4RoutePool) UnsetRoutePoolFlag () error {
 //parameters: RoutePoolHandle RoutePoolFlag
 //AgtBgp4RoutePool UnsetRoutePoolFlag
 return nil
}

func(np *Bgp4RoutePool) SetResourceClasses () error {
 //parameters: RoutePoolHandle ResourceClasses
 //AgtBgp4RoutePool SetResourceClasses
 return nil
}

func(np *Bgp4RoutePool) GetNumAdvertisedRoutes ()(string, error) {
 //parameters: RoutePoolHandle NumRoutes
 //AgtBgp4RoutePool GetNumAdvertisedRoutes
 return "", nil
}

func(np *Bgp4RoutePool) GetNumWithdrawnRoutes ()(string, error) {
 //parameters: RoutePoolHandle NumRoutes
 //AgtBgp4RoutePool GetNumWithdrawnRoutes
 return "", nil
}

func(np *Bgp4RoutePool) GetRoutes ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool GetRoutes
 return "", nil
}

func(np *Bgp4RoutePool) IsTrafficDestinationEnabled () error {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool IsTrafficDestinationEnabled
 return nil
}

func(np *Bgp4RoutePool) IsRoutePoolFlagSet () error {
 //parameters: RoutePoolHandle Flag
 //AgtBgp4RoutePool IsRoutePoolFlagSet
 return nil
}

func(np *Bgp4RoutePool) GetRoutePoolStatus ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool GetRoutePoolStatus
 return "", nil
}

func(np *Bgp4RoutePool) GetNthRoute ()(string, error) {
 //parameters: RoutePoolHandle Index
 //AgtBgp4RoutePool GetNthRoute
 return "", nil
}

func(np *Bgp4RoutePool) GetResourceClasses ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4RoutePool GetResourceClasses
 return "", nil
}

func(np *Bgp4RoutePool) GetSessionHandle ()(string, error) {
 //parameters: hRoutePool
 //AgtBgp4RoutePool GetSessionHandle
 return "", nil
}

