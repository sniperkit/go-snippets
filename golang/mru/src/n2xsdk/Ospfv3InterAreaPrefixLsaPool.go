package n2xsdk

type Ospfv3InterAreaPrefixLsaPool struct {
 Handler string
}

func(np *Ospfv3InterAreaPrefixLsaPool) SetAdvertisingRouter () error {
 //parameters: LsaHandle RouterId
 //AgtOspfv3InterAreaPrefixLsaPool SetAdvertisingRouter
 return nil
}

func(np *Ospfv3InterAreaPrefixLsaPool) GetAdvertisingRouter ()(string, error) {
 //parameters: LsaHandle
 //AgtOspfv3InterAreaPrefixLsaPool GetAdvertisingRouter
 return "", nil
}

func(np *Ospfv3InterAreaPrefixLsaPool) Advertise () error {
 //parameters: LsaHandle
 //AgtOspfv3InterAreaPrefixLsaPool Advertise
 return nil
}

func(np *Ospfv3InterAreaPrefixLsaPool) Withdraw () error {
 //parameters: LsaHandle
 //AgtOspfv3InterAreaPrefixLsaPool Withdraw
 return nil
}

func(np *Ospfv3InterAreaPrefixLsaPool) GetAdvertiseFlag ()(string, error) {
 //parameters: LsaHandle
 //AgtOspfv3InterAreaPrefixLsaPool GetAdvertiseFlag
 return "", nil
}

func(np *Ospfv3InterAreaPrefixLsaPool) SetRoutes () error {
 //parameters: LsaHandle FirstRoute PrefixLength NumRoutes Modifier
 //AgtOspfv3InterAreaPrefixLsaPool SetRoutes
 return nil
}

func(np *Ospfv3InterAreaPrefixLsaPool) GetRoutes ()(string, error) {
 //parameters: LsaHandle
 //AgtOspfv3InterAreaPrefixLsaPool GetRoutes
 return "", nil
}

func(np *Ospfv3InterAreaPrefixLsaPool) SetMetric () error {
 //parameters: LsaHandle Metric
 //AgtOspfv3InterAreaPrefixLsaPool SetMetric
 return nil
}

func(np *Ospfv3InterAreaPrefixLsaPool) GetMetric ()(string, error) {
 //parameters: LsaHandle
 //AgtOspfv3InterAreaPrefixLsaPool GetMetric
 return "", nil
}

func(np *Ospfv3InterAreaPrefixLsaPool) GetPrefixOptions ()(string, error) {
 //parameters: LsaHandle
 //AgtOspfv3InterAreaPrefixLsaPool GetPrefixOptions
 return "", nil
}

func(np *Ospfv3InterAreaPrefixLsaPool) SetPrefixOptions () error {
 //parameters: LsaHandle PrefixOptions
 //AgtOspfv3InterAreaPrefixLsaPool SetPrefixOptions
 return nil
}

