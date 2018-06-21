package n2xsdk

type PimSession struct {
 Handler string
}

func(np *PimSession) SetPimMode () error {
 //parameters: SessionHandle PimMode
 //AgtPimSession SetPimMode
 return nil
}

func(np *PimSession) GetPimMode ()(string, error) {
 //parameters: SessionHandle
 //AgtPimSession GetPimMode
 return "", nil
}

func(np *PimSession) SetRouterId () error {
 //parameters: SessionHandle RouterId
 //AgtPimSession SetRouterId
 return nil
}

func(np *PimSession) RemoveAllNeighborIpAddresses () error {
 //parameters: SessionHandle
 //AgtPimSession RemoveAllNeighborIpAddresses
 return nil
}

func(np *PimSession) ListNeighborIpAddresses ()(string, error) {
 //parameters: SessionHandle
 //AgtPimSession ListNeighborIpAddresses
 return "", nil
}

func(np *PimSession) GetNeighborIpAddress ()(string, error) {
 //parameters: SessionHandle SutAddressHandle
 //AgtPimSession GetNeighborIpAddress
 return "", nil
}

func(np *PimSession) AddMembershipPools () error {
 //parameters: SessionHandle PoolMode GroupPoolHandles
 //AgtPimSession AddMembershipPools
 return nil
}

func(np *PimSession) SetMembershipPoolMode () error {
 //parameters: SessionHandle MembershipPoolHandle PoolMode
 //AgtPimSession SetMembershipPoolMode
 return nil
}

func(np *PimSession) GetMembershipPoolMode ()(string, error) {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimSession GetMembershipPoolMode
 return "", nil
}

func(np *PimSession) ListMembershipPoolsByMode ()(string, error) {
 //parameters: SessionHandle PoolMode
 //AgtPimSession ListMembershipPoolsByMode
 return "", nil
}

func(np *PimSession) SetPoolRpAddress () error {
 //parameters: SessionHandle MembershipPoolHandle RpIpAddress
 //AgtPimSession SetPoolRpAddress
 return nil
}

func(np *PimSession) GetPoolRpAddress ()(string, error) {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimSession GetPoolRpAddress
 return "", nil
}

func(np *PimSession) SetPimTimer () error {
 //parameters: SessionHandle TimerType Period
 //AgtPimSession SetPimTimer
 return nil
}

func(np *PimSession) GetPimTimer ()(string, error) {
 //parameters: SessionHandle TimerType
 //AgtPimSession GetPimTimer
 return "", nil
}

func(np *PimSession) AddHelloTlv () error {
 //parameters: SessionHandle Type Length Value
 //AgtPimSession AddHelloTlv
 return nil
}

func(np *PimSession) RemoveHelloTlv () error {
 //parameters: SessionHandle TlvHandle
 //AgtPimSession RemoveHelloTlv
 return nil
}

func(np *PimSession) RemoveAllHelloTlvs () error {
 //parameters: SessionHandle
 //AgtPimSession RemoveAllHelloTlvs
 return nil
}

func(np *PimSession) ListHelloTlvs ()(string, error) {
 //parameters: SessionHandle
 //AgtPimSession ListHelloTlvs
 return "", nil
}

func(np *PimSession) GetHelloTlvDetails ()(string, error) {
 //parameters: SessionHandle TlvHandle
 //AgtPimSession GetHelloTlvDetails
 return "", nil
}

func(np *PimSession) GetLastError ()(string, error) {
 //parameters: SessionHandle
 //AgtPimSession GetLastError
 return "", nil
}

func(np *PimSession) GetPimv4PortStateSummary ()(string, error) {
 //parameters: PortHandle
 //AgtPimSession GetPimv4PortStateSummary
 return "", nil
}

func(np *PimSession) GetPimv4StateSummary ()(string, error) {
 //parameters: 
 //AgtPimSession GetPimv4StateSummary
 return "", nil
}

func(np *PimSession) Enable () error {
 //parameters: SessionHandle
 //AgtPimSession Enable
 return nil
}

func(np *PimSession) Disable () error {
 //parameters: SessionHandle
 //AgtPimSession Disable
 return nil
}

func(np *PimSession) EnableGreTunnel () error {
 //parameters: SessionHandle
 //AgtPimSession EnableGreTunnel
 return nil
}

func(np *PimSession) DisableGreTunnel () error {
 //parameters: SessionHandle
 //AgtPimSession DisableGreTunnel
 return nil
}

func(np *PimSession) IsGreTunnelEnabled () error {
 //parameters: SessionHandle
 //AgtPimSession IsGreTunnelEnabled
 return nil
}

func(np *PimSession) EnableGreTunnelChecksumField () error {
 //parameters: SessionHandle
 //AgtPimSession EnableGreTunnelChecksumField
 return nil
}

func(np *PimSession) DisableGreTunnelChecksumField () error {
 //parameters: SessionHandle
 //AgtPimSession DisableGreTunnelChecksumField
 return nil
}

func(np *PimSession) IsGreTunnelChecksumFieldEnabled () error {
 //parameters: SessionHandle
 //AgtPimSession IsGreTunnelChecksumFieldEnabled
 return nil
}

func(np *PimSession) SetGreTunnelLocalAddress () error {
 //parameters: SessionHandle IpAddress
 //AgtPimSession SetGreTunnelLocalAddress
 return nil
}

func(np *PimSession) GetGreTunnelLocalAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtPimSession GetGreTunnelLocalAddress
 return "", nil
}

func(np *PimSession) SetGreTunnelRemoteAddress () error {
 //parameters: SessionHandle IpAddress
 //AgtPimSession SetGreTunnelRemoteAddress
 return nil
}

func(np *PimSession) GetGreTunnelRemoteAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtPimSession GetGreTunnelRemoteAddress
 return "", nil
}

func(np *PimSession) SetPimTos () error {
 //parameters: SessionHandle Tos
 //AgtPimSession SetPimTos
 return nil
}

func(np *PimSession) GetPimTos ()(string, error) {
 //parameters: SessionHandle
 //AgtPimSession GetPimTos
 return "", nil
}

func(np *PimSession) SetGreTunnelTos () error {
 //parameters: SessionHandle Tos
 //AgtPimSession SetGreTunnelTos
 return nil
}

func(np *PimSession) GetGreTunnelTos ()(string, error) {
 //parameters: SessionHandle
 //AgtPimSession GetGreTunnelTos
 return "", nil
}

func(np *PimSession) EnablePimSmSimulatedRendezvousPoint () error {
 //parameters: SessionHandle
 //AgtPimSession EnablePimSmSimulatedRendezvousPoint
 return nil
}

func(np *PimSession) DisablePimSmSimulatedRendezvousPoint () error {
 //parameters: SessionHandle
 //AgtPimSession DisablePimSmSimulatedRendezvousPoint
 return nil
}

func(np *PimSession) EnableAll () error {
 //parameters: 
 //AgtPimSession EnableAll
 return nil
}

func(np *PimSession) DisableAll () error {
 //parameters: 
 //AgtPimSession DisableAll
 return nil
}

func(np *PimSession) SetInterfaceIpAddress () error {
 //parameters: 
 //AgtPimSession SetInterfaceIpAddress
 return nil
}

func(np *PimSession) GetInterfaceIpAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtPimSession GetInterfaceIpAddress
 return "", nil
}

func(np *PimSession) GetRouterId ()(string, error) {
 //parameters: SessionHandle
 //AgtPimSession GetRouterId
 return "", nil
}

func(np *PimSession) GetLinkLocalInterfaceAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtPimSession GetLinkLocalInterfaceAddress
 return "", nil
}

func(np *PimSession) GetNumberOfGroupAddresses ()(string, error) {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimSession GetNumberOfGroupAddresses
 return "", nil
}

func(np *PimSession) GetNumberOfSourceAddresses ()(string, error) {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimSession GetNumberOfSourceAddresses
 return "", nil
}

func(np *PimSession) RemoveMembershipPools () error {
 //parameters: SessionHandle Count psaMembershipPoolHandles
 //AgtPimSession RemoveMembershipPools
 return nil
}

func(np *PimSession) RemoveAllMembershipPools () error {
 //parameters: SessionHandle
 //AgtPimSession RemoveAllMembershipPools
 return nil
}

func(np *PimSession) ListMembershipPools ()(string, error) {
 //parameters: SessionHandle
 //AgtPimSession ListMembershipPools
 return "", nil
}

func(np *PimSession) SetGroupPoolHandle () error {
 //parameters: SessionHandle MembershipPoolHandle GroupPoolHandle
 //AgtPimSession SetGroupPoolHandle
 return nil
}

func(np *PimSession) GetGroupPoolHandle ()(string, error) {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimSession GetGroupPoolHandle
 return "", nil
}

func(np *PimSession) AddSourcePools () error {
 //parameters: SessionHandle MembershipPoolHandle NumSources psaSourcePoolHandles
 //AgtPimSession AddSourcePools
 return nil
}

func(np *PimSession) RemoveSourcePools () error {
 //parameters: SessionHandle MembershipPoolHandle Count psaSourcePoolHandles
 //AgtPimSession RemoveSourcePools
 return nil
}

func(np *PimSession) RemoveAllSourcePools () error {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimSession RemoveAllSourcePools
 return nil
}

func(np *PimSession) ListSourcePools ()(string, error) {
 //parameters: SessionHandle MembershipPoolHandle
 //AgtPimSession ListSourcePools
 return "", nil
}

func(np *PimSession) GetEnableFlag ()(string, error) {
 //parameters: SessionHandle
 //AgtPimSession GetEnableFlag
 return "", nil
}

func(np *PimSession) GetState ()(string, error) {
 //parameters: SessionHandle
 //AgtPimSession GetState
 return "", nil
}

func(np *PimSession) GetPortStateSummary ()(string, error) {
 //parameters: PortHandle
 //AgtPimSession GetPortStateSummary
 return "", nil
}

func(np *PimSession) GetStateSummary ()(string, error) {
 //parameters: 
 //AgtPimSession GetStateSummary
 return "", nil
}

func(np *PimSession) GetMulticastRoutingProtocol ()(string, error) {
 //parameters: SessionHandle
 //AgtPimSession GetMulticastRoutingProtocol
 return "", nil
}

func(np *PimSession) AddNeighborIpAddresses () error {
 //parameters: 
 //AgtPimSession AddNeighborIpAddresses
 return nil
}

func(np *PimSession) RemoveNeighborIpAddresses () error {
 //parameters: 
 //AgtPimSession RemoveNeighborIpAddresses
 return nil
}

