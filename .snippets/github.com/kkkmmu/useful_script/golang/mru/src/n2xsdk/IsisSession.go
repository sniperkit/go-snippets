package n2xsdk

type IsisSession struct {
 Handler string
}

func(np *IsisSession) SetMaximumPacketSize () error {
 //parameters: SessionHandle MaximumPacketSize
 //AgtIsisSession SetMaximumPacketSize
 return nil
}

func(np *IsisSession) GetMaximumPacketSize ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession GetMaximumPacketSize
 return "", nil
}

func(np *IsisSession) Restart () error {
 //parameters: SessionHandle
 //AgtIsisSession Restart
 return nil
}

func(np *IsisSession) RestartAllSessions () error {
 //parameters: 
 //AgtIsisSession RestartAllSessions
 return nil
}

func(np *IsisSession) GetNeighborState ()(string, error) {
 //parameters: SessionHandle NeighborHandle
 //AgtIsisSession GetNeighborState
 return "", nil
}

func(np *IsisSession) GetSutLinkState ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession GetSutLinkState
 return "", nil
}

func(np *IsisSession) GetSessionState ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession GetSessionState
 return "", nil
}

func(np *IsisSession) AddNeighbor () error {
 //parameters: SessionHandle
 //AgtIsisSession AddNeighbor
 return nil
}

func(np *IsisSession) RemoveNeighbor () error {
 //parameters: SessionHandle NeighborHandle
 //AgtIsisSession RemoveNeighbor
 return nil
}

func(np *IsisSession) RemoveAllNeighbors () error {
 //parameters: SessionHandle
 //AgtIsisSession RemoveAllNeighbors
 return nil
}

func(np *IsisSession) ListNeighbors ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession ListNeighbors
 return "", nil
}

func(np *IsisSession) SetSutLinkIpAddress () error {
 //parameters: SessionHandle Ipv4Address Ipv4PrefixLength
 //AgtIsisSession SetSutLinkIpAddress
 return nil
}

func(np *IsisSession) GetSutLinkIpAddress ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession GetSutLinkIpAddress
 return "", nil
}

func(np *IsisSession) SetSutLinkIpv4Address () error {
 //parameters: SessionHandle Ipv4Address Ipv4PrefixLength
 //AgtIsisSession SetSutLinkIpv4Address
 return nil
}

func(np *IsisSession) GetSutLinkIpv4Address ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession GetSutLinkIpv4Address
 return "", nil
}

func(np *IsisSession) SetSutLinkIpv6Address () error {
 //parameters: SessionHandle Ipv6Address Ipv6PrefixLength
 //AgtIsisSession SetSutLinkIpv6Address
 return nil
}

func(np *IsisSession) GetSutLinkIpv6Address ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession GetSutLinkIpv6Address
 return "", nil
}

func(np *IsisSession) GetSutLinkLocalIpv6Address ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession GetSutLinkLocalIpv6Address
 return "", nil
}

func(np *IsisSession) SetAreaId () error {
 //parameters: SessionHandle AreaId
 //AgtIsisSession SetAreaId
 return nil
}

func(np *IsisSession) GetAreaId ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession GetAreaId
 return "", nil
}

func(np *IsisSession) SetAreaList ()(string, error) {
 //parameters: SessionHandle AreaList
 //AgtIsisSession SetAreaList
 return "", nil
}

func(np *IsisSession) GetAreaList ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession GetAreaList
 return "", nil
}

func(np *IsisSession) SetSystemId () error {
 //parameters: SessionHandle SystemId
 //AgtIsisSession SetSystemId
 return nil
}

func(np *IsisSession) GetSystemId ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession GetSystemId
 return "", nil
}

func(np *IsisSession) IsEnabled () error {
 //parameters: SessionHandle
 //AgtIsisSession IsEnabled
 return nil
}

func(np *IsisSession) GetStateSummary ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession GetStateSummary
 return "", nil
}

func(np *IsisSession) GetLspDatabase ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession GetLspDatabase
 return "", nil
}

func(np *IsisSession) GetSessionRouter ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession GetSessionRouter
 return "", nil
}

func(np *IsisSession) SetTeRouterId () error {
 //parameters: SessionHandle TeRouterId
 //AgtIsisSession SetTeRouterId
 return nil
}

func(np *IsisSession) GetTeRouterId ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession GetTeRouterId
 return "", nil
}

func(np *IsisSession) SetPseudonodeNumber () error {
 //parameters: SessionHandle PseudonodeNumber
 //AgtIsisSession SetPseudonodeNumber
 return nil
}

func(np *IsisSession) GetPseudonodeNumber ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession GetPseudonodeNumber
 return "", nil
}

func(np *IsisSession) SetRoutingLevel () error {
 //parameters: SessionHandle RoutingLevel
 //AgtIsisSession SetRoutingLevel
 return nil
}

func(np *IsisSession) GetRoutingLevel ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession GetRoutingLevel
 return "", nil
}

func(np *IsisSession) SetNetworkType () error {
 //parameters: SessionHandle NetworkType
 //AgtIsisSession SetNetworkType
 return nil
}

func(np *IsisSession) GetNetworkType ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession GetNetworkType
 return "", nil
}

func(np *IsisSession) EnableWideMetrics () error {
 //parameters: SessionHandle
 //AgtIsisSession EnableWideMetrics
 return nil
}

func(np *IsisSession) DisableWideMetrics () error {
 //parameters: SessionHandle
 //AgtIsisSession DisableWideMetrics
 return nil
}

func(np *IsisSession) IsWideMetricsEnabled () error {
 //parameters: SessionHandle
 //AgtIsisSession IsWideMetricsEnabled
 return nil
}

func(np *IsisSession) SetIpMode () error {
 //parameters: SessionHandle IpMode
 //AgtIsisSession SetIpMode
 return nil
}

func(np *IsisSession) GetIpMode ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession GetIpMode
 return "", nil
}

func(np *IsisSession) EnableLspDiscardMode () error {
 //parameters: SessionHandle
 //AgtIsisSession EnableLspDiscardMode
 return nil
}

func(np *IsisSession) GetLspDiscardModeFlag ()(string, error) {
 //parameters: SessionHandle
 //AgtIsisSession GetLspDiscardModeFlag
 return "", nil
}

func(np *IsisSession) DisableLspDiscardMode () error {
 //parameters: SessionHandle
 //AgtIsisSession DisableLspDiscardMode
 return nil
}

func(np *IsisSession) EnableMultiTopologies () error {
 //parameters: SessionHandle
 //AgtIsisSession EnableMultiTopologies
 return nil
}

func(np *IsisSession) DisableMultiTopologies () error {
 //parameters: SessionHandle
 //AgtIsisSession DisableMultiTopologies
 return nil
}

func(np *IsisSession) IsMultiTopologiesEnabled () error {
 //parameters: SessionHandle
 //AgtIsisSession IsMultiTopologiesEnabled
 return nil
}

func(np *IsisSession) Enable3WayHandshake () error {
 //parameters: SessionHandle
 //AgtIsisSession Enable3WayHandshake
 return nil
}

func(np *IsisSession) Disable3WayHandshake () error {
 //parameters: SessionHandle
 //AgtIsisSession Disable3WayHandshake
 return nil
}

func(np *IsisSession) Is3WayHandshakeEnabled () error {
 //parameters: SessionHandle
 //AgtIsisSession Is3WayHandshakeEnabled
 return nil
}

func(np *IsisSession) EnableGracefulRestart () error {
 //parameters: SessionHandle
 //AgtIsisSession EnableGracefulRestart
 return nil
}

func(np *IsisSession) DisableGracefulRestart () error {
 //parameters: SessionHandle
 //AgtIsisSession DisableGracefulRestart
 return nil
}

func(np *IsisSession) IsGracefulRestartEnabled () error {
 //parameters: SessionHandle
 //AgtIsisSession IsGracefulRestartEnabled
 return nil
}

func(np *IsisSession) EnableIihPadding () error {
 //parameters: SessionHandle
 //AgtIsisSession EnableIihPadding
 return nil
}

func(np *IsisSession) DisableIihPadding () error {
 //parameters: SessionHandle
 //AgtIsisSession DisableIihPadding
 return nil
}

func(np *IsisSession) IsIihPaddingEnabled () error {
 //parameters: SessionHandle
 //AgtIsisSession IsIihPaddingEnabled
 return nil
}

func(np *IsisSession) SetSutLinkParameter () error {
 //parameters: SessionHandle SutLinkParameter SutLinkParameterValue
 //AgtIsisSession SetSutLinkParameter
 return nil
}

func(np *IsisSession) GetSutLinkParameter ()(string, error) {
 //parameters: SessionHandle SutLinkParameter
 //AgtIsisSession GetSutLinkParameter
 return "", nil
}

func(np *IsisSession) Enable () error {
 //parameters: SessionHandle
 //AgtIsisSession Enable
 return nil
}

func(np *IsisSession) Disable () error {
 //parameters: SessionHandle
 //AgtIsisSession Disable
 return nil
}

func(np *IsisSession) GetRoutersAndConnections ()(string, error) {
 //parameters: SessionHandle IsGetConnections
 //AgtIsisSession GetRoutersAndConnections
 return "", nil
}

func(np *IsisSession) SetAuthenticationMode () error {
 //parameters: SessionHandle KeyType AuthenticationMode
 //AgtIsisSession SetAuthenticationMode
 return nil
}

func(np *IsisSession) GetAuthenticationMode ()(string, error) {
 //parameters: SessionHandle KeyType
 //AgtIsisSession GetAuthenticationMode
 return "", nil
}

func(np *IsisSession) EnableAuthenticationKey () error {
 //parameters: SessionHandle KeyType
 //AgtIsisSession EnableAuthenticationKey
 return nil
}

func(np *IsisSession) DisableAuthenticationKey () error {
 //parameters: SessionHandle KeyType
 //AgtIsisSession DisableAuthenticationKey
 return nil
}

func(np *IsisSession) IsAuthenticationKeyEnabled () error {
 //parameters: SessionHandle KeyType
 //AgtIsisSession IsAuthenticationKeyEnabled
 return nil
}

func(np *IsisSession) EnableAuthenticationSendKey () error {
 //parameters: SessionHandle KeyType
 //AgtIsisSession EnableAuthenticationSendKey
 return nil
}

func(np *IsisSession) DisableAuthenticationSendKey () error {
 //parameters: SessionHandle KeyType
 //AgtIsisSession DisableAuthenticationSendKey
 return nil
}

func(np *IsisSession) IsAuthenticationSendKeyEnabled () error {
 //parameters: SessionHandle KeyType
 //AgtIsisSession IsAuthenticationSendKeyEnabled
 return nil
}

func(np *IsisSession) EnableAuthenticationReceiveKey () error {
 //parameters: SessionHandle KeyType
 //AgtIsisSession EnableAuthenticationReceiveKey
 return nil
}

func(np *IsisSession) DisableAuthenticationReceiveKey () error {
 //parameters: SessionHandle KeyType
 //AgtIsisSession DisableAuthenticationReceiveKey
 return nil
}

func(np *IsisSession) IsAuthenticationReceiveKeyEnabled () error {
 //parameters: SessionHandle KeyType
 //AgtIsisSession IsAuthenticationReceiveKeyEnabled
 return nil
}

func(np *IsisSession) SetAuthenticationSendKey () error {
 //parameters: SessionHandle KeyType Key
 //AgtIsisSession SetAuthenticationSendKey
 return nil
}

func(np *IsisSession) GetAuthenticationSendKey ()(string, error) {
 //parameters: SessionHandle KeyType
 //AgtIsisSession GetAuthenticationSendKey
 return "", nil
}

func(np *IsisSession) SetAuthenticationReceiveKey () error {
 //parameters: SessionHandle KeyType Key
 //AgtIsisSession SetAuthenticationReceiveKey
 return nil
}

func(np *IsisSession) GetAuthenticationReceiveKey ()(string, error) {
 //parameters: SessionHandle KeyType
 //AgtIsisSession GetAuthenticationReceiveKey
 return "", nil
}

func(np *IsisSession) EnableAllSessions () error {
 //parameters: 
 //AgtIsisSession EnableAllSessions
 return nil
}

func(np *IsisSession) DisableAllSessions () error {
 //parameters: 
 //AgtIsisSession DisableAllSessions
 return nil
}

func(np *IsisSession) SetMacAddressSource () error {
 //parameters: MacAddressSourceType
 //AgtIsisSession SetMacAddressSource
 return nil
}

func(np *IsisSession) GetMacAddressSource ()(string, error) {
 //parameters: 
 //AgtIsisSession GetMacAddressSource
 return "", nil
}

func(np *IsisSession) EnableAuthentication () error {
 //parameters: SessionHandle AuthenticationType
 //AgtIsisSession EnableAuthentication
 return nil
}

func(np *IsisSession) DisableAuthentication () error {
 //parameters: SessionHandle AuthenticationType
 //AgtIsisSession DisableAuthentication
 return nil
}

func(np *IsisSession) IsAuthenticationEnabled () error {
 //parameters: SessionHandle AuthenticationType
 //AgtIsisSession IsAuthenticationEnabled
 return nil
}

func(np *IsisSession) SetSendPassword () error {
 //parameters: SessionHandle AuthenticationType Password
 //AgtIsisSession SetSendPassword
 return nil
}

func(np *IsisSession) GetSendPassword ()(string, error) {
 //parameters: SessionHandle AuthenticationType
 //AgtIsisSession GetSendPassword
 return "", nil
}

func(np *IsisSession) SetReceivePassword () error {
 //parameters: SessionHandle AuthenticationType Count PasswordList
 //AgtIsisSession SetReceivePassword
 return nil
}

func(np *IsisSession) GetReceivePassword ()(string, error) {
 //parameters: SessionHandle AuthenticationType
 //AgtIsisSession GetReceivePassword
 return "", nil
}

