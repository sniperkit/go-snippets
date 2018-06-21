package n2xsdk

type IsisLspDatabase struct {
 Handler string
}

func(np *IsisLspDatabase) AdvertiseLsp () error {
 //parameters: LspHandle
 //AgtIsisLspDatabase AdvertiseLsp
 return nil
}

func(np *IsisLspDatabase) WithdrawLsp () error {
 //parameters: LspHandle
 //AgtIsisLspDatabase WithdrawLsp
 return nil
}

func(np *IsisLspDatabase) AdvertiseAll () error {
 //parameters: DatabaseHandle
 //AgtIsisLspDatabase AdvertiseAll
 return nil
}

func(np *IsisLspDatabase) WithdrawAll () error {
 //parameters: DatabaseHandle
 //AgtIsisLspDatabase WithdrawAll
 return nil
}

func(np *IsisLspDatabase) GetLspCount ()(string, error) {
 //parameters: DatabaseHandle Level
 //AgtIsisLspDatabase GetLspCount
 return "", nil
}

func(np *IsisLspDatabase) ListLsps ()(string, error) {
 //parameters: DatabaseHandle
 //AgtIsisLspDatabase ListLsps
 return "", nil
}

func(np *IsisLspDatabase) GetLspAdvertiseFlag ()(string, error) {
 //parameters: LspHandle
 //AgtIsisLspDatabase GetLspAdvertiseFlag
 return "", nil
}

