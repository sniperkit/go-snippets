package n2xsdk

type PimRouterPool struct {
 Handler string
}

func(np *PimRouterPool) SetTypeOfService () error {
 //parameters: RouterPoolHandle TypeOfService
 //AgtPimRouterPool SetTypeOfService
 return nil
}

func(np *PimRouterPool) GetTypeOfService ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool GetTypeOfService
 return "", nil
}

func(np *PimRouterPool) SetHoldTime () error {
 //parameters: RouterPoolHandle HoldTime
 //AgtPimRouterPool SetHoldTime
 return nil
}

func(np *PimRouterPool) GetHoldTime ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool GetHoldTime
 return "", nil
}

func(np *PimRouterPool) SetDrPriority () error {
 //parameters: RouterPoolHandle DrPriority
 //AgtPimRouterPool SetDrPriority
 return nil
}

func(np *PimRouterPool) GetDrPriority ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool GetDrPriority
 return "", nil
}

func(np *PimRouterPool) AddGroup () error {
 //parameters: RouterPoolHandle GroupType
 //AgtPimRouterPool AddGroup
 return nil
}

func(np *PimRouterPool) RemoveGroup () error {
 //parameters: RouterPoolHandle GroupHandle
 //AgtPimRouterPool RemoveGroup
 return nil
}

func(np *PimRouterPool) RemoveAllGroups () error {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool RemoveAllGroups
 return nil
}

func(np *PimRouterPool) GetGroupType ()(string, error) {
 //parameters: RouterPoolHandle GroupHandle
 //AgtPimRouterPool GetGroupType
 return "", nil
}

func(np *PimRouterPool) GetGroupName ()(string, error) {
 //parameters: RouterPoolHandle GroupHandle
 //AgtPimRouterPool GetGroupName
 return "", nil
}

func(np *PimRouterPool) SetGroupName () error {
 //parameters: RouterPoolHandle GroupHandle GroupName
 //AgtPimRouterPool SetGroupName
 return nil
}

func(np *PimRouterPool) ListGroups ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool ListGroups
 return "", nil
}

func(np *PimRouterPool) ListGroupsByType ()(string, error) {
 //parameters: RouterPoolHandle GroupType
 //AgtPimRouterPool ListGroupsByType
 return "", nil
}

func(np *PimRouterPool) EnableGreTunnel () error {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool EnableGreTunnel
 return nil
}

func(np *PimRouterPool) DisableGreTunnel () error {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool DisableGreTunnel
 return nil
}

func(np *PimRouterPool) IsGreTunnelEnabled () error {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool IsGreTunnelEnabled
 return nil
}

func(np *PimRouterPool) SetGreTunnelLocalAddressIncrementingRange () error {
 //parameters: RouterPoolHandle FirstIpAddress PrefixLength Increment Repeat Count
 //AgtPimRouterPool SetGreTunnelLocalAddressIncrementingRange
 return nil
}

func(np *PimRouterPool) GetGreTunnelLocalAddressIncrementingRange ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool GetGreTunnelLocalAddressIncrementingRange
 return "", nil
}

func(np *PimRouterPool) SetGreTunnelLocalAddressList ()(string, error) {
 //parameters: RouterPoolHandle IpAddressList
 //AgtPimRouterPool SetGreTunnelLocalAddressList
 return "", nil
}

func(np *PimRouterPool) GetGreTunnelLocalAddressList ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool GetGreTunnelLocalAddressList
 return "", nil
}

func(np *PimRouterPool) SetGreTunnelLocalAddress () error {
 //parameters: RouterPoolHandle IpAddress
 //AgtPimRouterPool SetGreTunnelLocalAddress
 return nil
}

func(np *PimRouterPool) GetGreTunnelLocalAddress ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool GetGreTunnelLocalAddress
 return "", nil
}

func(np *PimRouterPool) SetGreTunnelRemoteAddressIncrementingRange () error {
 //parameters: RouterPoolHandle FirstIpAddress PrefixLength Increment Repeat Count
 //AgtPimRouterPool SetGreTunnelRemoteAddressIncrementingRange
 return nil
}

func(np *PimRouterPool) GetGreTunnelRemoteAddressIncrementingRange ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool GetGreTunnelRemoteAddressIncrementingRange
 return "", nil
}

func(np *PimRouterPool) SetGreTunnelRemoteAddressList ()(string, error) {
 //parameters: RouterPoolHandle IpAddressList
 //AgtPimRouterPool SetGreTunnelRemoteAddressList
 return "", nil
}

func(np *PimRouterPool) GetGreTunnelRemoteAddressList ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool GetGreTunnelRemoteAddressList
 return "", nil
}

func(np *PimRouterPool) SetGreTunnelRemoteAddress () error {
 //parameters: RouterPoolHandle IpAddress
 //AgtPimRouterPool SetGreTunnelRemoteAddress
 return nil
}

func(np *PimRouterPool) GetGreTunnelRemoteAddress ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool GetGreTunnelRemoteAddress
 return "", nil
}

func(np *PimRouterPool) EnableGreTunnelChecksum () error {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool EnableGreTunnelChecksum
 return nil
}

func(np *PimRouterPool) DisableGreTunnelChecksum () error {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool DisableGreTunnelChecksum
 return nil
}

func(np *PimRouterPool) IsGreTunnelChecksumEnabled () error {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool IsGreTunnelChecksumEnabled
 return nil
}

func(np *PimRouterPool) EnableNeighborRouter () error {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool EnableNeighborRouter
 return nil
}

func(np *PimRouterPool) DisableNeighborRouter () error {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool DisableNeighborRouter
 return nil
}

func(np *PimRouterPool) IsNeighborRouterEnabled () error {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool IsNeighborRouterEnabled
 return nil
}

func(np *PimRouterPool) SetNeighborRouterAddressIncrementingRange () error {
 //parameters: RouterPoolHandle FirstIpAddress PrefixLength Increment Repeat Count
 //AgtPimRouterPool SetNeighborRouterAddressIncrementingRange
 return nil
}

func(np *PimRouterPool) GetNeighborRouterAddressIncrementingRange ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool GetNeighborRouterAddressIncrementingRange
 return "", nil
}

func(np *PimRouterPool) SetNeighborRouterAddressList ()(string, error) {
 //parameters: RouterPoolHandle IpAddressList
 //AgtPimRouterPool SetNeighborRouterAddressList
 return "", nil
}

func(np *PimRouterPool) GetNeighborRouterAddressList ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool GetNeighborRouterAddressList
 return "", nil
}

func(np *PimRouterPool) SetNeighborRouterAddress () error {
 //parameters: RouterPoolHandle IpAddress
 //AgtPimRouterPool SetNeighborRouterAddress
 return nil
}

func(np *PimRouterPool) GetNeighborRouterAddress ()(string, error) {
 //parameters: RouterPoolHandle
 //AgtPimRouterPool GetNeighborRouterAddress
 return "", nil
}

func(np *PimRouterPool) SetTimer () error {
 //parameters: RouterPoolHandle TimerType TimerValue
 //AgtPimRouterPool SetTimer
 return nil
}

func(np *PimRouterPool) GetTimer ()(string, error) {
 //parameters: RouterPoolHandle TimerType
 //AgtPimRouterPool GetTimer
 return "", nil
}

