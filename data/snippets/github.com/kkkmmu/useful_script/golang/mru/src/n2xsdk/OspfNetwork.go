package n2xsdk

type OspfNetwork struct {
 Handler string
}

func(np *OspfNetwork) SetNetwork () error {
 //parameters: NetworkHandle IpAdress PrefixLength
 //AgtOspfNetwork SetNetwork
 return nil
}

func(np *OspfNetwork) GetNetwork ()(string, error) {
 //parameters: NetworkHandle
 //AgtOspfNetwork GetNetwork
 return "", nil
}

func(np *OspfNetwork) GetNetworkLsa ()(string, error) {
 //parameters: NetworkHandle
 //AgtOspfNetwork GetNetworkLsa
 return "", nil
}

func(np *OspfNetwork) GetIntraAreaPrefixLsa ()(string, error) {
 //parameters: NetworkHandle
 //AgtOspfNetwork GetIntraAreaPrefixLsa
 return "", nil
}

func(np *OspfNetwork) SetDesignatedRouter () error {
 //parameters: NetworkHandle RouterHandle
 //AgtOspfNetwork SetDesignatedRouter
 return nil
}

func(np *OspfNetwork) GetDesignatedRouter ()(string, error) {
 //parameters: NetworkHandle
 //AgtOspfNetwork GetDesignatedRouter
 return "", nil
}

func(np *OspfNetwork) ListConnectedRouters ()(string, error) {
 //parameters: NetworkHandle
 //AgtOspfNetwork ListConnectedRouters
 return "", nil
}

func(np *OspfNetwork) GetVersion ()(string, error) {
 //parameters: NetworkHandle
 //AgtOspfNetwork GetVersion
 return "", nil
}

