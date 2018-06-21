package n2xsdk

type Bgp4MVrfPool struct {
 Handler string
}

func(np *Bgp4MVrfPool) QueryMVrfType1Entry () error {
 //parameters: VrfPoolHandle VrfIndex RouteTargetType RouteTarget PmsiTunnelEndPoint
 //AgtBgp4MVrfPool QueryMVrfType1Entry
 return nil
}

func(np *Bgp4MVrfPool) QueryMVrfType4Entry () error {
 //parameters: VrfPoolHandle VrfIndex RouteTargetType RouteTarget PmsiTunnelEndPoint SourceAddress GroupAddress
 //AgtBgp4MVrfPool QueryMVrfType4Entry
 return nil
}

func(np *Bgp4MVrfPool) QueryMVrfType5Entry () error {
 //parameters: VrfPoolHandle VrfIndex RouteTargetType RouteTarget SourceAddress GroupAddress
 //AgtBgp4MVrfPool QueryMVrfType5Entry
 return nil
}

func(np *Bgp4MVrfPool) QueryMVrfType7Entry () error {
 //parameters: VrfPoolHandle VrfIndex RouteTargetType RouteTarget SourceAddress GroupAddress
 //AgtBgp4MVrfPool QueryMVrfType7Entry
 return nil
}

func(np *Bgp4MVrfPool) SetPmsiTunnelType () error {
 //parameters: VrfPoolHandle PmsiTunnelType
 //AgtBgp4MVrfPool SetPmsiTunnelType
 return nil
}

func(np *Bgp4MVrfPool) GetPmsiTunnelType ()(string, error) {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool GetPmsiTunnelType
 return "", nil
}

func(np *Bgp4MVrfPool) EnableSourceActiveFlag () error {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool EnableSourceActiveFlag
 return nil
}

func(np *Bgp4MVrfPool) DisableSourceActiveFlag () error {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool DisableSourceActiveFlag
 return nil
}

func(np *Bgp4MVrfPool) IsSourceActiveFlagEnabled () error {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool IsSourceActiveFlagEnabled
 return nil
}

func(np *Bgp4MVrfPool) SetImportCMcastRtPeAddressIncrementingRange () error {
 //parameters: VrfPoolHandle FirstRouteTarget PrefixLength RouteTargetIncrement RouteTargetRepeat PercentOverlap
 //AgtBgp4MVrfPool SetImportCMcastRtPeAddressIncrementingRange
 return nil
}

func(np *Bgp4MVrfPool) GetImportCMcastRtPeAddressIncrementingRange ()(string, error) {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool GetImportCMcastRtPeAddressIncrementingRange
 return "", nil
}

func(np *Bgp4MVrfPool) SetImportCMcastRtPeAddressList ()(string, error) {
 //parameters: VrfPoolHandle RouteTargetList
 //AgtBgp4MVrfPool SetImportCMcastRtPeAddressList
 return "", nil
}

func(np *Bgp4MVrfPool) GetImportCMcastRtPeAddressList ()(string, error) {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool GetImportCMcastRtPeAddressList
 return "", nil
}

func(np *Bgp4MVrfPool) SetImportCMcastRtPeAddress () error {
 //parameters: VrfPoolHandle RouteTarget
 //AgtBgp4MVrfPool SetImportCMcastRtPeAddress
 return nil
}

func(np *Bgp4MVrfPool) GetImportCMcastRtPeAddress ()(string, error) {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool GetImportCMcastRtPeAddress
 return "", nil
}

func(np *Bgp4MVrfPool) SetImportCMcastRtVrfIdIncrementingRange () error {
 //parameters: VrfPoolHandle FirstRouteTargetId RouteTargetIncrement PercentOverlap
 //AgtBgp4MVrfPool SetImportCMcastRtVrfIdIncrementingRange
 return nil
}

func(np *Bgp4MVrfPool) GetImportCMcastRtVrfIdIncrementingRange ()(string, error) {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool GetImportCMcastRtVrfIdIncrementingRange
 return "", nil
}

func(np *Bgp4MVrfPool) SetImportCMcastRtVrfIdList ()(string, error) {
 //parameters: VrfPoolHandle RouteTargetIdList
 //AgtBgp4MVrfPool SetImportCMcastRtVrfIdList
 return "", nil
}

func(np *Bgp4MVrfPool) GetImportCMcastRtVrfIdList ()(string, error) {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool GetImportCMcastRtVrfIdList
 return "", nil
}

func(np *Bgp4MVrfPool) SetImportCMcastRtVrfId () error {
 //parameters: VrfPoolHandle RouteTargetId
 //AgtBgp4MVrfPool SetImportCMcastRtVrfId
 return nil
}

func(np *Bgp4MVrfPool) GetImportCMcastRtVrfId ()(string, error) {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool GetImportCMcastRtVrfId
 return "", nil
}

func(np *Bgp4MVrfPool) SetImportLeafADRtIncrementingRange () error {
 //parameters: VrfPoolHandle FirstRouteTarget RouteTargetIncrement RouteTargetRepeat PercentOverlap
 //AgtBgp4MVrfPool SetImportLeafADRtIncrementingRange
 return nil
}

func(np *Bgp4MVrfPool) GetImportLeafADRtIncrementingRange ()(string, error) {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool GetImportLeafADRtIncrementingRange
 return "", nil
}

func(np *Bgp4MVrfPool) SetImportLeafADRtList ()(string, error) {
 //parameters: VrfPoolHandle RouteTargetList
 //AgtBgp4MVrfPool SetImportLeafADRtList
 return "", nil
}

func(np *Bgp4MVrfPool) GetImportLeafADRtList ()(string, error) {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool GetImportLeafADRtList
 return "", nil
}

func(np *Bgp4MVrfPool) SetImportLeafADRt () error {
 //parameters: VrfPoolHandle RouteTarget
 //AgtBgp4MVrfPool SetImportLeafADRt
 return nil
}

func(np *Bgp4MVrfPool) GetImportLeafADRt ()(string, error) {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool GetImportLeafADRt
 return "", nil
}

func(np *Bgp4MVrfPool) SetImportRouteTargetList ()(string, error) {
 //parameters: VrfPoolHandle RouteTargetType RouteTargetList
 //AgtBgp4MVrfPool SetImportRouteTargetList
 return "", nil
}

func(np *Bgp4MVrfPool) GetImportRouteTargetList ()(string, error) {
 //parameters: VrfPoolHandle RouteTargetType
 //AgtBgp4MVrfPool GetImportRouteTargetList
 return "", nil
}

func(np *Bgp4MVrfPool) SetImportRouteTarget () error {
 //parameters: VrfPoolHandle RouteTargetType RouteTarget
 //AgtBgp4MVrfPool SetImportRouteTarget
 return nil
}

func(np *Bgp4MVrfPool) GetImportRouteTarget ()(string, error) {
 //parameters: VrfPoolHandle RouteTargetType
 //AgtBgp4MVrfPool GetImportRouteTarget
 return "", nil
}

func(np *Bgp4MVrfPool) WithdrawAllType4Routes () error {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool WithdrawAllType4Routes
 return nil
}

func(np *Bgp4MVrfPool) WithdrawType4RoutesInRt () error {
 //parameters: VrfPoolHandle VrfIndex RouteTargetType RouteTarget PmsiTunnelEndPoint
 //AgtBgp4MVrfPool WithdrawType4RoutesInRt
 return nil
}

func(np *Bgp4MVrfPool) WithdrawType4Route () error {
 //parameters: VrfPoolHandle VrfIndex RouteTargetType RouteTarget PmsiTunnelEndPoint SourceAddress GroupAddress
 //AgtBgp4MVrfPool WithdrawType4Route
 return nil
}

func(np *Bgp4MVrfPool) WithdrawAllType5Routes () error {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool WithdrawAllType5Routes
 return nil
}

func(np *Bgp4MVrfPool) WithdrawType5RoutesInRt () error {
 //parameters: VrfPoolHandle VrfIndex RouteTargetType RouteTarget
 //AgtBgp4MVrfPool WithdrawType5RoutesInRt
 return nil
}

func(np *Bgp4MVrfPool) WithdrawType5Route () error {
 //parameters: VrfPoolHandle VrfIndex RouteTargetType RouteTarget SourceAddress GroupAddress
 //AgtBgp4MVrfPool WithdrawType5Route
 return nil
}

func(np *Bgp4MVrfPool) WithdrawAllType7Routes () error {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool WithdrawAllType7Routes
 return nil
}

func(np *Bgp4MVrfPool) WithdrawType7RoutesInRt () error {
 //parameters: VrfPoolHandle VrfIndex RouteTargetType RouteTarget
 //AgtBgp4MVrfPool WithdrawType7RoutesInRt
 return nil
}

func(np *Bgp4MVrfPool) WithdrawType7Route () error {
 //parameters: VrfPoolHandle VrfIndex RouteTargetType RouteTarget SourceAddress GroupAddress
 //AgtBgp4MVrfPool WithdrawType7Route
 return nil
}

func(np *Bgp4MVrfPool) Enable () error {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool Enable
 return nil
}

func(np *Bgp4MVrfPool) Disable () error {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool Disable
 return nil
}

func(np *Bgp4MVrfPool) IsEnabled () error {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool IsEnabled
 return nil
}

func(np *Bgp4MVrfPool) SetImportRouteTargetIncrementingRange () error {
 //parameters: VrfPoolHandle RouteTargetType FirstRouteTarget RouteTargetIncrement RouteTargetRepeat PercentOverlap
 //AgtBgp4MVrfPool SetImportRouteTargetIncrementingRange
 return nil
}

func(np *Bgp4MVrfPool) GetImportRouteTargetIncrementingRange ()(string, error) {
 //parameters: VrfPoolHandle RouteTargetType
 //AgtBgp4MVrfPool GetImportRouteTargetIncrementingRange
 return "", nil
}

func(np *Bgp4MVrfPool) SetImportRouteTargetType () error {
 //parameters: VrfPoolHandle RouteTargetTypeValue
 //AgtBgp4MVrfPool SetImportRouteTargetType
 return nil
}

func(np *Bgp4MVrfPool) GetImportRouteTargetType ()(string, error) {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool GetImportRouteTargetType
 return "", nil
}

func(np *Bgp4MVrfPool) SetVpnsPerPeer () error {
 //parameters: VrfPoolHandle NumVpnsPerPeer
 //AgtBgp4MVrfPool SetVpnsPerPeer
 return nil
}

func(np *Bgp4MVrfPool) GetVpnsPerPeer ()(string, error) {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool GetVpnsPerPeer
 return "", nil
}

func(np *Bgp4MVrfPool) GetTotalNumberOfVpns ()(string, error) {
 //parameters: VrfPoolHandle
 //AgtBgp4MVrfPool GetTotalNumberOfVpns
 return "", nil
}

