package n2xsdk

type OspfExternalLsaPool struct {
 Handler string
}

func(np *OspfExternalLsaPool) SetAdvertisingRouter () error {
 //parameters: LsaHandle RouterId
 //AgtOspfExternalLsaPool SetAdvertisingRouter
 return nil
}

func(np *OspfExternalLsaPool) GetAdvertisingRouter ()(string, error) {
 //parameters: LsaHandle
 //AgtOspfExternalLsaPool GetAdvertisingRouter
 return "", nil
}

func(np *OspfExternalLsaPool) Advertise () error {
 //parameters: LsaHandle
 //AgtOspfExternalLsaPool Advertise
 return nil
}

func(np *OspfExternalLsaPool) Withdraw () error {
 //parameters: LsaHandle
 //AgtOspfExternalLsaPool Withdraw
 return nil
}

func(np *OspfExternalLsaPool) GetAdvertiseFlag ()(string, error) {
 //parameters: LsaHandle
 //AgtOspfExternalLsaPool GetAdvertiseFlag
 return "", nil
}

func(np *OspfExternalLsaPool) SetRoutes () error {
 //parameters: LsaHandle FirstRoute PrefixLength NumRoutes Modifier
 //AgtOspfExternalLsaPool SetRoutes
 return nil
}

func(np *OspfExternalLsaPool) GetRoutes ()(string, error) {
 //parameters: LsaHandle
 //AgtOspfExternalLsaPool GetRoutes
 return "", nil
}

func(np *OspfExternalLsaPool) SetMetric () error {
 //parameters: LsaHandle Metric
 //AgtOspfExternalLsaPool SetMetric
 return nil
}

func(np *OspfExternalLsaPool) GetMetric ()(string, error) {
 //parameters: LsaHandle
 //AgtOspfExternalLsaPool GetMetric
 return "", nil
}

func(np *OspfExternalLsaPool) SetEbit () error {
 //parameters: LsaHandle Flag
 //AgtOspfExternalLsaPool SetEbit
 return nil
}

func(np *OspfExternalLsaPool) GetEbit ()(string, error) {
 //parameters: LsaHandle
 //AgtOspfExternalLsaPool GetEbit
 return "", nil
}

func(np *OspfExternalLsaPool) SetForwardingAddress () error {
 //parameters: LsaHandle IpAddress
 //AgtOspfExternalLsaPool SetForwardingAddress
 return nil
}

func(np *OspfExternalLsaPool) GetForwardingAddress ()(string, error) {
 //parameters: LsaHandle
 //AgtOspfExternalLsaPool GetForwardingAddress
 return "", nil
}

func(np *OspfExternalLsaPool) SetExternalRouteTag () error {
 //parameters: LsaHandle RouteTag
 //AgtOspfExternalLsaPool SetExternalRouteTag
 return nil
}

func(np *OspfExternalLsaPool) GetExternalRouteTag ()(string, error) {
 //parameters: LsaHandle
 //AgtOspfExternalLsaPool GetExternalRouteTag
 return "", nil
}

