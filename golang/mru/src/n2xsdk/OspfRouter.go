package n2xsdk

type OspfRouter struct {
 Handler string
}

func(np *OspfRouter) SetRouterId () error {
 //parameters: RouterHandle RouterId
 //AgtOspfRouter SetRouterId
 return nil
}

func(np *OspfRouter) GetRouterId ()(string, error) {
 //parameters: RouterHandle
 //AgtOspfRouter GetRouterId
 return "", nil
}

func(np *OspfRouter) SetRouterType () error {
 //parameters: RouterHandle RouterTypeValue
 //AgtOspfRouter SetRouterType
 return nil
}

func(np *OspfRouter) GetRouterType ()(string, error) {
 //parameters: RouterHandle
 //AgtOspfRouter GetRouterType
 return "", nil
}

func(np *OspfRouter) EnableTe () error {
 //parameters: RouterHandle
 //AgtOspfRouter EnableTe
 return nil
}

func(np *OspfRouter) DisableTe () error {
 //parameters: RouterHandle
 //AgtOspfRouter DisableTe
 return nil
}

func(np *OspfRouter) IsTeEnabled () error {
 //parameters: RouterHandle
 //AgtOspfRouter IsTeEnabled
 return nil
}

func(np *OspfRouter) EnableLinkTe () error {
 //parameters: RouterHandle OtherObjectHandle
 //AgtOspfRouter EnableLinkTe
 return nil
}

func(np *OspfRouter) DisableLinkTe () error {
 //parameters: RouterHandle OtherObjectHandle
 //AgtOspfRouter DisableLinkTe
 return nil
}

func(np *OspfRouter) IsLinkTeEnabled () error {
 //parameters: RouterHandle OtherObjectHandle
 //AgtOspfRouter IsLinkTeEnabled
 return nil
}

func(np *OspfRouter) GetRouterLsa ()(string, error) {
 //parameters: RouterHandle
 //AgtOspfRouter GetRouterLsa
 return "", nil
}

func(np *OspfRouter) GetLinkLsa ()(string, error) {
 //parameters: RouterHandle
 //AgtOspfRouter GetLinkLsa
 return "", nil
}

func(np *OspfRouter) GetIntraAreaPrefixLsa ()(string, error) {
 //parameters: RouterHandle
 //AgtOspfRouter GetIntraAreaPrefixLsa
 return "", nil
}

func(np *OspfRouter) GetTeRouterLsa ()(string, error) {
 //parameters: RouterHandle
 //AgtOspfRouter GetTeRouterLsa
 return "", nil
}

func(np *OspfRouter) ListTeLinkLsas ()(string, error) {
 //parameters: RouterHandle
 //AgtOspfRouter ListTeLinkLsas
 return "", nil
}

func(np *OspfRouter) ConnectObjects () error {
 //parameters: RouterHandle OtherObjectHandle
 //AgtOspfRouter ConnectObjects
 return nil
}

func(np *OspfRouter) DisconnectObjects () error {
 //parameters: RouterHandle OtherObjectHandle
 //AgtOspfRouter DisconnectObjects
 return nil
}

func(np *OspfRouter) ListObjectConnections ()(string, error) {
 //parameters: RouterHandle
 //AgtOspfRouter ListObjectConnections
 return "", nil
}

func(np *OspfRouter) GetLinkMetric ()(string, error) {
 //parameters: RouterHandle OtherObjectHandle
 //AgtOspfRouter GetLinkMetric
 return "", nil
}

func(np *OspfRouter) SetLinkMetric () error {
 //parameters: RouterHandle OtherObjectHandle Metric
 //AgtOspfRouter SetLinkMetric
 return nil
}

func(np *OspfRouter) EnableLink () error {
 //parameters: RouterHandle OtherObjectHandle
 //AgtOspfRouter EnableLink
 return nil
}

func(np *OspfRouter) DisableLink () error {
 //parameters: RouterHandle OtherObjectHandle
 //AgtOspfRouter DisableLink
 return nil
}

func(np *OspfRouter) IsLinkEnabled () error {
 //parameters: RouterHandle OtherObjectHandle
 //AgtOspfRouter IsLinkEnabled
 return nil
}

func(np *OspfRouter) Advertise () error {
 //parameters: RouterHandle
 //AgtOspfRouter Advertise
 return nil
}

func(np *OspfRouter) Withdraw () error {
 //parameters: RouterHandle
 //AgtOspfRouter Withdraw
 return nil
}

func(np *OspfRouter) GetLinkType ()(string, error) {
 //parameters: RouterHandle OtherObjectHandle
 //AgtOspfRouter GetLinkType
 return "", nil
}

func(np *OspfRouter) SetLinkInterfaceAddress () error {
 //parameters: RouterHandle OtherObjectHandle IpAddress
 //AgtOspfRouter SetLinkInterfaceAddress
 return nil
}

func(np *OspfRouter) GetLinkInterfaceAddress ()(string, error) {
 //parameters: RouterHandle OtherObjectHandle
 //AgtOspfRouter GetLinkInterfaceAddress
 return "", nil
}

func(np *OspfRouter) GetTeLinkLsa ()(string, error) {
 //parameters: RouterHandle OtherObjectHandle
 //AgtOspfRouter GetTeLinkLsa
 return "", nil
}

func(np *OspfRouter) GetLinkTeInstance ()(string, error) {
 //parameters: RouterHandle OtherObjectHandle
 //AgtOspfRouter GetLinkTeInstance
 return "", nil
}

func(np *OspfRouter) SetLinkTeInstance () error {
 //parameters: RouterHandle OtherObjectHandle Instance
 //AgtOspfRouter SetLinkTeInstance
 return nil
}

func(np *OspfRouter) GetAdvertiseFlag ()(string, error) {
 //parameters: RouterHandle
 //AgtOspfRouter GetAdvertiseFlag
 return "", nil
}

func(np *OspfRouter) GetVersion ()(string, error) {
 //parameters: RouterHandle
 //AgtOspfRouter GetVersion
 return "", nil
}

func(np *OspfRouter) GetRouterDetails ()(string, error) {
 //parameters: RouterHandle
 //AgtOspfRouter GetRouterDetails
 return "", nil
}

func(np *OspfRouter) GetLinkDetails ()(string, error) {
 //parameters: RouterHandle OtherObjectHandle
 //AgtOspfRouter GetLinkDetails
 return "", nil
}

