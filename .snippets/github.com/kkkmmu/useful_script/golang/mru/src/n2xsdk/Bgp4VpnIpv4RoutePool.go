package n2xsdk

type Bgp4VpnIpv4RoutePool struct {
 Handler string
}

func(np *Bgp4VpnIpv4RoutePool) SetRoutes () error {
 //parameters: RoutePoolHandle FirstRoute PrefixLength NumRoutes Modifier
 //AgtBgp4VpnIpv4RoutePool SetRoutes
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) GetLabel ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool GetLabel
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) SetNextHop () error {
 //parameters: RoutePoolHandle NextHop
 //AgtBgp4VpnIpv4RoutePool SetNextHop
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) GetNextHop ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool GetNextHop
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) SetNextHops () error {
 //parameters: RoutePoolHandle NextHops
 //AgtBgp4VpnIpv4RoutePool SetNextHops
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) SetOrigin () error {
 //parameters: RoutePoolHandle Origin
 //AgtBgp4VpnIpv4RoutePool SetOrigin
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) GetOrigin ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool GetOrigin
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) SetCommunities () error {
 //parameters: RoutePoolHandle CommunitiesList
 //AgtBgp4VpnIpv4RoutePool SetCommunities
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) GetCommunities ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool GetCommunities
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) SetPathAttributeString () error {
 //parameters: RoutePoolHandle PathAttribute PathAttributeString
 //AgtBgp4VpnIpv4RoutePool SetPathAttributeString
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) GetPathAttributeString ()(string, error) {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4VpnIpv4RoutePool GetPathAttributeString
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) SetPathAttributeLong () error {
 //parameters: RoutePoolHandle PathAttribute PathAttributeLong
 //AgtBgp4VpnIpv4RoutePool SetPathAttributeLong
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) GetPathAttributeLong ()(string, error) {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4VpnIpv4RoutePool GetPathAttributeLong
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) IsPathAttributeSelected () error {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4VpnIpv4RoutePool IsPathAttributeSelected
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) IsPathAttributeSet () error {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4VpnIpv4RoutePool IsPathAttributeSet
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) SetAsPath () error {
 //parameters: RoutePoolHandle AsPathAttributeType AsNumber
 //AgtBgp4VpnIpv4RoutePool SetAsPath
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) GetAsPath ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool GetAsPath
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) GetAsPathByType ()(string, error) {
 //parameters: RoutePoolHandle AsPathAttributeType
 //AgtBgp4VpnIpv4RoutePool GetAsPathByType
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) GetMaximumAsPathLength ()(string, error) {
 //parameters: 
 //AgtBgp4VpnIpv4RoutePool GetMaximumAsPathLength
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) SetAggregator () error {
 //parameters: RoutePoolHandle AsNumber IpAddress
 //AgtBgp4VpnIpv4RoutePool SetAggregator
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) GetAggregator ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool GetAggregator
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) SetOriginatorId () error {
 //parameters: RoutePoolHandle OriginatorId
 //AgtBgp4VpnIpv4RoutePool SetOriginatorId
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) GetOriginatorId ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool GetOriginatorId
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) SetClusterList ()(string, error) {
 //parameters: RoutePoolHandle ClusterList
 //AgtBgp4VpnIpv4RoutePool SetClusterList
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) GetClusterList ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool GetClusterList
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) Enable () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool Enable
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) Disable () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool Disable
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) IsEnabled () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool IsEnabled
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) Advertise () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool Advertise
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) AdvertisePools () error {
 //parameters: RoutePoolHandles RoutesPerUpdate InterUpdateDelay
 //AgtBgp4VpnIpv4RoutePool AdvertisePools
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) Withdraw () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool Withdraw
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) WithdrawPools () error {
 //parameters: RoutePoolHandles RoutesPerUpdate InterUpdateDelay
 //AgtBgp4VpnIpv4RoutePool WithdrawPools
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) GetState ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool GetState
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) SelectPathAttribute () error {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4VpnIpv4RoutePool SelectPathAttribute
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) DeselectPathAttribute () error {
 //parameters: RoutePoolHandle PathAttribute
 //AgtBgp4VpnIpv4RoutePool DeselectPathAttribute
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) EnableTrafficDestinations () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool EnableTrafficDestinations
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) DisableTrafficDestinations () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool DisableTrafficDestinations
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) EnableTraffic () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool EnableTraffic
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) DisableTraffic () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool DisableTraffic
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) SetRoutePoolFlag () error {
 //parameters: RoutePoolHandle RoutePoolFlag
 //AgtBgp4VpnIpv4RoutePool SetRoutePoolFlag
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) UnsetRoutePoolFlag () error {
 //parameters: RoutePoolHandle RoutePoolFlag
 //AgtBgp4VpnIpv4RoutePool UnsetRoutePoolFlag
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) SetResourceClasses () error {
 //parameters: RoutePoolHandle ResourceClasses
 //AgtBgp4VpnIpv4RoutePool SetResourceClasses
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) GetNumAdvertisedRoutes ()(string, error) {
 //parameters: RoutePoolHandle NumRoutes
 //AgtBgp4VpnIpv4RoutePool GetNumAdvertisedRoutes
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) GetNumWithdrawnRoutes ()(string, error) {
 //parameters: RoutePoolHandle NumRoutes
 //AgtBgp4VpnIpv4RoutePool GetNumWithdrawnRoutes
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) SetExportTarget () error {
 //parameters: RoutePoolHandle RouteTargetType Administrator AssignedNumber
 //AgtBgp4VpnIpv4RoutePool SetExportTarget
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) GetExportTarget ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool GetExportTarget
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) GetRouteDistinguisher ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool GetRouteDistinguisher
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) GetSiteHandle ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool GetSiteHandle
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) GetRoutes ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool GetRoutes
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) IsTrafficDestinationEnabled () error {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool IsTrafficDestinationEnabled
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) IsRoutePoolFlagSet () error {
 //parameters: RoutePoolHandle Flag
 //AgtBgp4VpnIpv4RoutePool IsRoutePoolFlagSet
 return nil
}

func(np *Bgp4VpnIpv4RoutePool) GetRoutePoolStatus ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool GetRoutePoolStatus
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) GetNthRoute ()(string, error) {
 //parameters: RoutePoolHandle Index
 //AgtBgp4VpnIpv4RoutePool GetNthRoute
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) GetResourceClasses ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtBgp4VpnIpv4RoutePool GetResourceClasses
 return "", nil
}

func(np *Bgp4VpnIpv4RoutePool) GetSessionHandle ()(string, error) {
 //parameters: hRoutePool
 //AgtBgp4VpnIpv4RoutePool GetSessionHandle
 return "", nil
}

