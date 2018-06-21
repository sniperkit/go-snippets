package n2xsdk

type Ospfv3ExternalRoutePool struct {
 Handler string
}

func(np *Ospfv3ExternalRoutePool) SetRoutes () error {
 //parameters: RoutePoolHandle FirstRoute PrefixLength NumRoutes Modifier
 //AgtOspfv3ExternalRoutePool SetRoutes
 return nil
}

func(np *Ospfv3ExternalRoutePool) GetRoutes ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtOspfv3ExternalRoutePool GetRoutes
 return "", nil
}

func(np *Ospfv3ExternalRoutePool) EnableTrafficDestinations () error {
 //parameters: RoutePoolHandle
 //AgtOspfv3ExternalRoutePool EnableTrafficDestinations
 return nil
}

func(np *Ospfv3ExternalRoutePool) DisableTrafficDestinations () error {
 //parameters: RoutePoolHandle
 //AgtOspfv3ExternalRoutePool DisableTrafficDestinations
 return nil
}

func(np *Ospfv3ExternalRoutePool) IsTrafficDestinationEnabled () error {
 //parameters: RoutePoolHandle
 //AgtOspfv3ExternalRoutePool IsTrafficDestinationEnabled
 return nil
}

func(np *Ospfv3ExternalRoutePool) SetRoutePoolFlag () error {
 //parameters: RoutePoolHandle Flag
 //AgtOspfv3ExternalRoutePool SetRoutePoolFlag
 return nil
}

func(np *Ospfv3ExternalRoutePool) UnsetRoutePoolFlag () error {
 //parameters: RoutePoolHandle Flag
 //AgtOspfv3ExternalRoutePool UnsetRoutePoolFlag
 return nil
}

func(np *Ospfv3ExternalRoutePool) IsRoutePoolFlagSet () error {
 //parameters: RoutePoolHandle Flag
 //AgtOspfv3ExternalRoutePool IsRoutePoolFlagSet
 return nil
}

func(np *Ospfv3ExternalRoutePool) GetRoutePoolStatus ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtOspfv3ExternalRoutePool GetRoutePoolStatus
 return "", nil
}

func(np *Ospfv3ExternalRoutePool) GetNthRoute ()(string, error) {
 //parameters: RoutePoolHandle Index
 //AgtOspfv3ExternalRoutePool GetNthRoute
 return "", nil
}

func(np *Ospfv3ExternalRoutePool) SetResourceClasses () error {
 //parameters: RoutePoolHandle Count psaClasses
 //AgtOspfv3ExternalRoutePool SetResourceClasses
 return nil
}

func(np *Ospfv3ExternalRoutePool) GetResourceClasses ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtOspfv3ExternalRoutePool GetResourceClasses
 return "", nil
}

func(np *Ospfv3ExternalRoutePool) GetConnectedRouter ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtOspfv3ExternalRoutePool GetConnectedRouter
 return "", nil
}

func(np *Ospfv3ExternalRoutePool) SetNssaFlag () error {
 //parameters: RoutePoolHandle NssaFlag
 //AgtOspfv3ExternalRoutePool SetNssaFlag
 return nil
}

func(np *Ospfv3ExternalRoutePool) GetNssaFlag ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtOspfv3ExternalRoutePool GetNssaFlag
 return "", nil
}

func(np *Ospfv3ExternalRoutePool) SetPropagateLsasFlag () error {
 //parameters: RoutePoolHandle PropagateFlag
 //AgtOspfv3ExternalRoutePool SetPropagateLsasFlag
 return nil
}

func(np *Ospfv3ExternalRoutePool) GetPropagateLsasFlag ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtOspfv3ExternalRoutePool GetPropagateLsasFlag
 return "", nil
}

func(np *Ospfv3ExternalRoutePool) SetAutoGenerateInterAreaRouterLsaFlag () error {
 //parameters: RoutePoolHandle AsbrFlag
 //AgtOspfv3ExternalRoutePool SetAutoGenerateInterAreaRouterLsaFlag
 return nil
}

func(np *Ospfv3ExternalRoutePool) GetAutoGenerateInterAreaRouterLsaFlag ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtOspfv3ExternalRoutePool GetAutoGenerateInterAreaRouterLsaFlag
 return "", nil
}

func(np *Ospfv3ExternalRoutePool) SetAdvertisingRouter () error {
 //parameters: RoutePoolHandle RouterId
 //AgtOspfv3ExternalRoutePool SetAdvertisingRouter
 return nil
}

func(np *Ospfv3ExternalRoutePool) GetAdvertisingRouter ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtOspfv3ExternalRoutePool GetAdvertisingRouter
 return "", nil
}

func(np *Ospfv3ExternalRoutePool) SetMetricType () error {
 //parameters: RoutePoolHandle Type
 //AgtOspfv3ExternalRoutePool SetMetricType
 return nil
}

func(np *Ospfv3ExternalRoutePool) GetMetricType ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtOspfv3ExternalRoutePool GetMetricType
 return "", nil
}

func(np *Ospfv3ExternalRoutePool) SetMetric () error {
 //parameters: RoutePoolHandle Metric
 //AgtOspfv3ExternalRoutePool SetMetric
 return nil
}

func(np *Ospfv3ExternalRoutePool) GetMetric ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtOspfv3ExternalRoutePool GetMetric
 return "", nil
}

func(np *Ospfv3ExternalRoutePool) SetPrefixOptions () error {
 //parameters: RoutePoolHandle PrefixOptions
 //AgtOspfv3ExternalRoutePool SetPrefixOptions
 return nil
}

func(np *Ospfv3ExternalRoutePool) GetPrefixOptions ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtOspfv3ExternalRoutePool GetPrefixOptions
 return "", nil
}

func(np *Ospfv3ExternalRoutePool) SetReferencedLsType () error {
 //parameters: LsaHandle LsType
 //AgtOspfv3ExternalRoutePool SetReferencedLsType
 return nil
}

func(np *Ospfv3ExternalRoutePool) GetReferencedLsType ()(string, error) {
 //parameters: LsaHandle
 //AgtOspfv3ExternalRoutePool GetReferencedLsType
 return "", nil
}

func(np *Ospfv3ExternalRoutePool) SetFbit () error {
 //parameters: RoutePoolHandle FbitFlag
 //AgtOspfv3ExternalRoutePool SetFbit
 return nil
}

func(np *Ospfv3ExternalRoutePool) GetFbit ()(string, error) {
 //parameters: LsaHandle
 //AgtOspfv3ExternalRoutePool GetFbit
 return "", nil
}

func(np *Ospfv3ExternalRoutePool) SetForwardingAddress () error {
 //parameters: LsaHandle ForwardingAddress
 //AgtOspfv3ExternalRoutePool SetForwardingAddress
 return nil
}

func(np *Ospfv3ExternalRoutePool) GetForwardingAddress ()(string, error) {
 //parameters: LsaHandle
 //AgtOspfv3ExternalRoutePool GetForwardingAddress
 return "", nil
}

func(np *Ospfv3ExternalRoutePool) SetTbit () error {
 //parameters: LsaHandle TbitFlag
 //AgtOspfv3ExternalRoutePool SetTbit
 return nil
}

func(np *Ospfv3ExternalRoutePool) GetTbit ()(string, error) {
 //parameters: LsaHandle
 //AgtOspfv3ExternalRoutePool GetTbit
 return "", nil
}

func(np *Ospfv3ExternalRoutePool) SetExternalRouteTag () error {
 //parameters: LsaHandle RouteTag
 //AgtOspfv3ExternalRoutePool SetExternalRouteTag
 return nil
}

func(np *Ospfv3ExternalRoutePool) GetExternalRouteTag ()(string, error) {
 //parameters: LsaHandle
 //AgtOspfv3ExternalRoutePool GetExternalRouteTag
 return "", nil
}

func(np *Ospfv3ExternalRoutePool) SetReferencedLinkStateId () error {
 //parameters: RoutePoolHandle RefLsId
 //AgtOspfv3ExternalRoutePool SetReferencedLinkStateId
 return nil
}

func(np *Ospfv3ExternalRoutePool) GetReferencedLinkStateId ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtOspfv3ExternalRoutePool GetReferencedLinkStateId
 return "", nil
}

func(np *Ospfv3ExternalRoutePool) GetOspfv3ExternalLsaPool ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtOspfv3ExternalRoutePool GetOspfv3ExternalLsaPool
 return "", nil
}

func(np *Ospfv3ExternalRoutePool) GetInterAreaRouterLsa ()(string, error) {
 //parameters: RoutePoolHandle
 //AgtOspfv3ExternalRoutePool GetInterAreaRouterLsa
 return "", nil
}

