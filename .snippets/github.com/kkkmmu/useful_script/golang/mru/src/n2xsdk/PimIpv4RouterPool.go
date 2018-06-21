package n2xsdk

type PimIpv4RouterPool struct {
 Handler string
}

func(np *PimIpv4RouterPool) AddGroup () error {
 //parameters: RouterPoolHandle GroupType
 //AgtPimIpv4RouterPool AddGroup
 return nil
}

func(np *PimIpv4RouterPool) RemoveGroup () error {
 //parameters: RouterPoolHandle GroupHandle
 //AgtPimIpv4RouterPool RemoveGroup
 return nil
}

func(np *PimIpv4RouterPool) RemoveAllGroups () error {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool RemoveAllGroups
 return nil
}

func(np *PimIpv4RouterPool) GetGroupType ()(string, error) {
 //parameters: RouterPoolHandle GroupHandle
 //AgtPimIpv4RouterPool GetGroupType
 return "", nil
}

func(np *PimIpv4RouterPool) GetGroupName ()(string, error) {
 //parameters: RouterPoolHandle GroupHandle
 //AgtPimIpv4RouterPool GetGroupName
 return "", nil
}

func(np *PimIpv4RouterPool) SetGroupName () error {
 //parameters: RouterPoolHandle GroupHandle GroupName
 //AgtPimIpv4RouterPool SetGroupName
 return nil
}

func(np *PimIpv4RouterPool) ListGroups ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool ListGroups
 return "", nil
}

func(np *PimIpv4RouterPool) ListGroupsByType ()(string, error) {
 //parameters: RouterPoolHandle GroupType
 //AgtPimIpv4RouterPool ListGroupsByType
 return "", nil
}

func(np *PimIpv4RouterPool) SetNeighborRouterAddressIncrementingRange () error {
 //parameters: RouterPoolHandle FirstIpAddress PrefixLength Increment Repeat Count
 //AgtPimIpv4RouterPool SetNeighborRouterAddressIncrementingRange
 return nil
}

func(np *PimIpv4RouterPool) GetNeighborRouterAddressIncrementingRange ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool GetNeighborRouterAddressIncrementingRange
 return "", nil
}

func(np *PimIpv4RouterPool) SetNeighborRouterAddressList ()(string, error) {
 //parameters: RouterPoolHandle IpAddressList
 //AgtPimIpv4RouterPool SetNeighborRouterAddressList
 return "", nil
}

func(np *PimIpv4RouterPool) GetNeighborRouterAddressList ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool GetNeighborRouterAddressList
 return "", nil
}

func(np *PimIpv4RouterPool) SetNeighborRouterAddress () error {
 //parameters: RouterPoolHandle IpAddress
 //AgtPimIpv4RouterPool SetNeighborRouterAddress
 return nil
}

func(np *PimIpv4RouterPool) GetNeighborRouterAddress ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool GetNeighborRouterAddress
 return "", nil
}

func(np *PimIpv4RouterPool) EnableDataMdt () error {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool EnableDataMdt
 return nil
}

func(np *PimIpv4RouterPool) DisableDataMdt () error {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool DisableDataMdt
 return nil
}

func(np *PimIpv4RouterPool) IsDataMdtEnabled () error {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool IsDataMdtEnabled
 return nil
}

func(np *PimIpv4RouterPool) SetDataMdtTimeout () error {
 //parameters: RouterPoolHandle DataMdtTimeout
 //AgtPimIpv4RouterPool SetDataMdtTimeout
 return nil
}

func(np *PimIpv4RouterPool) GetDataMdtTimeout ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool GetDataMdtTimeout
 return "", nil
}

func(np *PimIpv4RouterPool) RefreshDataMdts () error {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool RefreshDataMdts
 return nil
}

func(np *PimIpv4RouterPool) GetDefaultMdts ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool GetDefaultMdts
 return "", nil
}

func(np *PimIpv4RouterPool) GetDataMdts ()(string, error) {
 //parameters: RouterPoolHandle DefaultMdtGroupAddress
 //AgtPimIpv4RouterPool GetDataMdts
 return "", nil
}

func(np *PimIpv4RouterPool) SetTypeOfService () error {
 //parameters: RouterPoolHandle TypeOfService
 //AgtPimIpv4RouterPool SetTypeOfService
 return nil
}

func(np *PimIpv4RouterPool) GetTypeOfService ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool GetTypeOfService
 return "", nil
}

func(np *PimIpv4RouterPool) SetHoldTime () error {
 //parameters: RouterPoolHandle HoldTime
 //AgtPimIpv4RouterPool SetHoldTime
 return nil
}

func(np *PimIpv4RouterPool) GetHoldTime ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool GetHoldTime
 return "", nil
}

func(np *PimIpv4RouterPool) SetDrPriority () error {
 //parameters: RouterPoolHandle DrPriority
 //AgtPimIpv4RouterPool SetDrPriority
 return nil
}

func(np *PimIpv4RouterPool) GetDrPriority ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool GetDrPriority
 return "", nil
}

func(np *PimIpv4RouterPool) EnableGreTunnel () error {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool EnableGreTunnel
 return nil
}

func(np *PimIpv4RouterPool) DisableGreTunnel () error {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool DisableGreTunnel
 return nil
}

func(np *PimIpv4RouterPool) IsGreTunnelEnabled () error {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool IsGreTunnelEnabled
 return nil
}

func(np *PimIpv4RouterPool) SetGreTunnelLocalAddressIncrementingRange () error {
 //parameters: RouterPoolHandle FirstIpAddress PrefixLength Increment Repeat Count
 //AgtPimIpv4RouterPool SetGreTunnelLocalAddressIncrementingRange
 return nil
}

func(np *PimIpv4RouterPool) GetGreTunnelLocalAddressIncrementingRange ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool GetGreTunnelLocalAddressIncrementingRange
 return "", nil
}

func(np *PimIpv4RouterPool) SetGreTunnelLocalAddressList ()(string, error) {
 //parameters: RouterPoolHandle IpAddressList
 //AgtPimIpv4RouterPool SetGreTunnelLocalAddressList
 return "", nil
}

func(np *PimIpv4RouterPool) GetGreTunnelLocalAddressList ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool GetGreTunnelLocalAddressList
 return "", nil
}

func(np *PimIpv4RouterPool) SetGreTunnelLocalAddress () error {
 //parameters: RouterPoolHandle IpAddress
 //AgtPimIpv4RouterPool SetGreTunnelLocalAddress
 return nil
}

func(np *PimIpv4RouterPool) GetGreTunnelLocalAddress ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool GetGreTunnelLocalAddress
 return "", nil
}

func(np *PimIpv4RouterPool) SetGreTunnelRemoteAddressIncrementingRange () error {
 //parameters: RouterPoolHandle FirstIpAddress PrefixLength Increment Repeat Count
 //AgtPimIpv4RouterPool SetGreTunnelRemoteAddressIncrementingRange
 return nil
}

func(np *PimIpv4RouterPool) GetGreTunnelRemoteAddressIncrementingRange ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool GetGreTunnelRemoteAddressIncrementingRange
 return "", nil
}

func(np *PimIpv4RouterPool) SetGreTunnelRemoteAddressList ()(string, error) {
 //parameters: RouterPoolHandle IpAddressList
 //AgtPimIpv4RouterPool SetGreTunnelRemoteAddressList
 return "", nil
}

func(np *PimIpv4RouterPool) GetGreTunnelRemoteAddressList ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool GetGreTunnelRemoteAddressList
 return "", nil
}

func(np *PimIpv4RouterPool) SetGreTunnelRemoteAddress () error {
 //parameters: RouterPoolHandle IpAddress
 //AgtPimIpv4RouterPool SetGreTunnelRemoteAddress
 return nil
}

func(np *PimIpv4RouterPool) GetGreTunnelRemoteAddress ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool GetGreTunnelRemoteAddress
 return "", nil
}

func(np *PimIpv4RouterPool) EnableGreTunnelChecksum () error {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool EnableGreTunnelChecksum
 return nil
}

func(np *PimIpv4RouterPool) DisableGreTunnelChecksum () error {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool DisableGreTunnelChecksum
 return nil
}

func(np *PimIpv4RouterPool) IsGreTunnelChecksumEnabled () error {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool IsGreTunnelChecksumEnabled
 return nil
}

func(np *PimIpv4RouterPool) EnableNeighborRouter () error {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool EnableNeighborRouter
 return nil
}

func(np *PimIpv4RouterPool) DisableNeighborRouter () error {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool DisableNeighborRouter
 return nil
}

func(np *PimIpv4RouterPool) IsNeighborRouterEnabled () error {
 //parameters: RouterPoolHandle
 //AgtPimIpv4RouterPool IsNeighborRouterEnabled
 return nil
}

func(np *PimIpv4RouterPool) SetTimer () error {
 //parameters: RouterPoolHandle TimerType TimerValue
 //AgtPimIpv4RouterPool SetTimer
 return nil
}

func(np *PimIpv4RouterPool) GetTimer ()(string, error) {
 //parameters: RouterPoolHandle TimerType
 //AgtPimIpv4RouterPool GetTimer
 return "", nil
}

