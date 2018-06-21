package n2xsdk

type Ospfv2RouterPool struct {
 Handler string
}

func(np *Ospfv2RouterPool) GetState ()(string, error) {
 //parameters: RouterPoolHandle RouterPoolInstance
 //AgtOspfv2RouterPool GetState
 return "", nil
}

func(np *Ospfv2RouterPool) GetNetworkType ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtOspfv2RouterPool GetNetworkType
 return "", nil
}

func(np *Ospfv2RouterPool) SetRouterId () error {
 //parameters: RouterPoolHandle FirstRouterId Increment Repeat
 //AgtOspfv2RouterPool SetRouterId
 return nil
}

func(np *Ospfv2RouterPool) GetRouterId ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtOspfv2RouterPool GetRouterId
 return "", nil
}

func(np *Ospfv2RouterPool) SetAreaId () error {
 //parameters: RouterPoolHandle FirstAreaId Increment Repeat
 //AgtOspfv2RouterPool SetAreaId
 return nil
}

func(np *Ospfv2RouterPool) GetAreaId ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtOspfv2RouterPool GetAreaId
 return "", nil
}

func(np *Ospfv2RouterPool) EnableOption () error {
 //parameters: RouterPoolHandle OptionsFieldBitType
 //AgtOspfv2RouterPool EnableOption
 return nil
}

func(np *Ospfv2RouterPool) DisableOption () error {
 //parameters: RouterPoolHandle OptionsFieldBitType
 //AgtOspfv2RouterPool DisableOption
 return nil
}

func(np *Ospfv2RouterPool) IsOptionEnabled () error {
 //parameters: RouterPoolHandle OptionsFieldBitType
 //AgtOspfv2RouterPool IsOptionEnabled
 return nil
}

func(np *Ospfv2RouterPool) SetInterval () error {
 //parameters: RouterPoolHandle IntervalType Interval
 //AgtOspfv2RouterPool SetInterval
 return nil
}

func(np *Ospfv2RouterPool) GetInterval ()(string, error) {
 //parameters: RouterPoolHandle IntervalType
 //AgtOspfv2RouterPool GetInterval
 return "", nil
}

func(np *Ospfv2RouterPool) EnableChecksumValidation () error {
 //parameters: RouterPoolHandle ChecksumType
 //AgtOspfv2RouterPool EnableChecksumValidation
 return nil
}

func(np *Ospfv2RouterPool) DisableChecksumValidation () error {
 //parameters: RouterPoolHandle ChecksumType
 //AgtOspfv2RouterPool DisableChecksumValidation
 return nil
}

func(np *Ospfv2RouterPool) IsChecksumValidationEnabled () error {
 //parameters: RouterPoolHandle ChecksumType
 //AgtOspfv2RouterPool IsChecksumValidationEnabled
 return nil
}

func(np *Ospfv2RouterPool) SetInterfaceCost () error {
 //parameters: RouterPoolHandle InterfaceCost
 //AgtOspfv2RouterPool SetInterfaceCost
 return nil
}

func(np *Ospfv2RouterPool) GetInterfaceCost ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtOspfv2RouterPool GetInterfaceCost
 return "", nil
}

func(np *Ospfv2RouterPool) EnableUpdateRateControl () error {
 //parameters: RouterPoolHandle
 //AgtOspfv2RouterPool EnableUpdateRateControl
 return nil
}

func(np *Ospfv2RouterPool) DisableUpdateRateControl () error {
 //parameters: RouterPoolHandle
 //AgtOspfv2RouterPool DisableUpdateRateControl
 return nil
}

func(np *Ospfv2RouterPool) IsUpdateRateControlEnabled () error {
 //parameters: RouterPoolHandle
 //AgtOspfv2RouterPool IsUpdateRateControlEnabled
 return nil
}

func(np *Ospfv2RouterPool) SetInterUpdateDelay () error {
 //parameters: RouterPoolHandle InterUpdateDelay
 //AgtOspfv2RouterPool SetInterUpdateDelay
 return nil
}

func(np *Ospfv2RouterPool) GetInterUpdateDelay ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtOspfv2RouterPool GetInterUpdateDelay
 return "", nil
}

func(np *Ospfv2RouterPool) EnableLsaDiscard () error {
 //parameters: RouterPoolHandle
 //AgtOspfv2RouterPool EnableLsaDiscard
 return nil
}

func(np *Ospfv2RouterPool) DisableLsaDiscard () error {
 //parameters: RouterPoolHandle
 //AgtOspfv2RouterPool DisableLsaDiscard
 return nil
}

func(np *Ospfv2RouterPool) IsLsaDiscardEnabled () error {
 //parameters: RouterPoolHandle
 //AgtOspfv2RouterPool IsLsaDiscardEnabled
 return nil
}

func(np *Ospfv2RouterPool) AddLsaProfile () error {
 //parameters: RouterPoolHandle LsaProfileType
 //AgtOspfv2RouterPool AddLsaProfile
 return nil
}

func(np *Ospfv2RouterPool) RemoveLsaProfile () error {
 //parameters: RouterPoolHandle LsaProfileHandle
 //AgtOspfv2RouterPool RemoveLsaProfile
 return nil
}

func(np *Ospfv2RouterPool) RemoveAllLsaProfiles () error {
 //parameters: RouterPoolHandle
 //AgtOspfv2RouterPool RemoveAllLsaProfiles
 return nil
}

func(np *Ospfv2RouterPool) ListLsaProfiles ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtOspfv2RouterPool ListLsaProfiles
 return "", nil
}

func(np *Ospfv2RouterPool) ListLsaProfilesByType ()(string, error) {
 //parameters: RouterPoolHandle LsaProfileType
 //AgtOspfv2RouterPool ListLsaProfilesByType
 return "", nil
}

