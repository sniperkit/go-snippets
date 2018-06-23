package n2xsdk

type OspfLsaDatabase struct {
 Handler string
}

func(np *OspfLsaDatabase) AddLsa () error {
 //parameters: DatabaseHandle LsaType
 //AgtOspfLsaDatabase AddLsa
 return nil
}

func(np *OspfLsaDatabase) RemoveLsa () error {
 //parameters: LsaHandle
 //AgtOspfLsaDatabase RemoveLsa
 return nil
}

func(np *OspfLsaDatabase) Clear () error {
 //parameters: DatabaseHandle
 //AgtOspfLsaDatabase Clear
 return nil
}

func(np *OspfLsaDatabase) AdvertiseLsa () error {
 //parameters: LsaHandle
 //AgtOspfLsaDatabase AdvertiseLsa
 return nil
}

func(np *OspfLsaDatabase) WithdrawLsa () error {
 //parameters: LsaHandle
 //AgtOspfLsaDatabase WithdrawLsa
 return nil
}

func(np *OspfLsaDatabase) AdvertiseAll () error {
 //parameters: DatabaseHandle
 //AgtOspfLsaDatabase AdvertiseAll
 return nil
}

func(np *OspfLsaDatabase) WithdrawAll () error {
 //parameters: DatabaseHandle
 //AgtOspfLsaDatabase WithdrawAll
 return nil
}

func(np *OspfLsaDatabase) GetLsaType ()(string, error) {
 //parameters: LsaHandle
 //AgtOspfLsaDatabase GetLsaType
 return "", nil
}

func(np *OspfLsaDatabase) GetLsaCount ()(string, error) {
 //parameters: DatabaseHandle
 //AgtOspfLsaDatabase GetLsaCount
 return "", nil
}

func(np *OspfLsaDatabase) GetLsaCountByType ()(string, error) {
 //parameters: DatabaseHandle LsaType
 //AgtOspfLsaDatabase GetLsaCountByType
 return "", nil
}

func(np *OspfLsaDatabase) ListLsasByType ()(string, error) {
 //parameters: DatabaseHandle LsaType
 //AgtOspfLsaDatabase ListLsasByType
 return "", nil
}

func(np *OspfLsaDatabase) ListLsas ()(string, error) {
 //parameters: DatabaseHandle
 //AgtOspfLsaDatabase ListLsas
 return "", nil
}

func(np *OspfLsaDatabase) GetLsaAdvertiseFlag ()(string, error) {
 //parameters: LsaHandle
 //AgtOspfLsaDatabase GetLsaAdvertiseFlag
 return "", nil
}

