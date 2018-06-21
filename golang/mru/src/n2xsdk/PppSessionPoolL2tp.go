package n2xsdk

type PppSessionPoolL2tp struct {
 Handler string
}

func(np *PppSessionPoolL2tp) SetTesterIpAddress () error {
 //parameters: SessionPoolHandle TesterIpAddress
 //AgtPppSessionPoolL2tp SetTesterIpAddress
 return nil
}

func(np *PppSessionPoolL2tp) SetSutIpAddress () error {
 //parameters: SessionPoolHandle SutIpAddress
 //AgtPppSessionPoolL2tp SetSutIpAddress
 return nil
}

func(np *PppSessionPoolL2tp) SetIdentity () error {
 //parameters: SessionPoolHandle Identity
 //AgtPppSessionPoolL2tp SetIdentity
 return nil
}

func(np *PppSessionPoolL2tp) SetNumberOfCallsPerTunnel () error {
 //parameters: SessionPoolHandle CallsPerTunnel
 //AgtPppSessionPoolL2tp SetNumberOfCallsPerTunnel
 return nil
}

func(np *PppSessionPoolL2tp) GetTunnelLifetime ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetTunnelLifetime
 return "", nil
}

func(np *PppSessionPoolL2tp) SetTunnelLifetime () error {
 //parameters: SessionPoolHandle TunnelLifetime
 //AgtPppSessionPoolL2tp SetTunnelLifetime
 return nil
}

func(np *PppSessionPoolL2tp) SetHelloInterval () error {
 //parameters: SessionPoolHandle Interval
 //AgtPppSessionPoolL2tp SetHelloInterval
 return nil
}

func(np *PppSessionPoolL2tp) SetInitialAckTimeout () error {
 //parameters: SessionPoolHandle Timeout
 //AgtPppSessionPoolL2tp SetInitialAckTimeout
 return nil
}

func(np *PppSessionPoolL2tp) SetMaxTxRetries () error {
 //parameters: SessionPoolHandle Retries
 //AgtPppSessionPoolL2tp SetMaxTxRetries
 return nil
}

func(np *PppSessionPoolL2tp) SetHostName () error {
 //parameters: SessionPoolHandle HostName
 //AgtPppSessionPoolL2tp SetHostName
 return nil
}

func(np *PppSessionPoolL2tp) SetSharedSecret () error {
 //parameters: SessionPoolHandle Secret
 //AgtPppSessionPoolL2tp SetSharedSecret
 return nil
}

func(np *PppSessionPoolL2tp) SetReceiveWindowSize () error {
 //parameters: SessionPoolHandle WindowSize
 //AgtPppSessionPoolL2tp SetReceiveWindowSize
 return nil
}

func(np *PppSessionPoolL2tp) SetSourceUdpPort () error {
 //parameters: SessionPoolHandle SourcePort
 //AgtPppSessionPoolL2tp SetSourceUdpPort
 return nil
}

func(np *PppSessionPoolL2tp) EnableAuthenticateSut () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp EnableAuthenticateSut
 return nil
}

func(np *PppSessionPoolL2tp) DisableAuthenticateSut () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp DisableAuthenticateSut
 return nil
}

func(np *PppSessionPoolL2tp) GetAuthenticateSutFlag ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetAuthenticateSutFlag
 return "", nil
}

func(np *PppSessionPoolL2tp) GetBearerCapabilities ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetBearerCapabilities
 return "", nil
}

func(np *PppSessionPoolL2tp) SetBearerCapabilities () error {
 //parameters: SessionPoolHandle Capabilities
 //AgtPppSessionPoolL2tp SetBearerCapabilities
 return nil
}

func(np *PppSessionPoolL2tp) GetBearerType ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetBearerType
 return "", nil
}

func(np *PppSessionPoolL2tp) SetBearerType () error {
 //parameters: SessionPoolHandle BearerType
 //AgtPppSessionPoolL2tp SetBearerType
 return nil
}

func(np *PppSessionPoolL2tp) GetCalledNumber ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetCalledNumber
 return "", nil
}

func(np *PppSessionPoolL2tp) SetCalledNumber () error {
 //parameters: SessionPoolHandle CalledNumber
 //AgtPppSessionPoolL2tp SetCalledNumber
 return nil
}

func(np *PppSessionPoolL2tp) GetCallingNumber ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetCallingNumber
 return "", nil
}

func(np *PppSessionPoolL2tp) SetCallingNumber () error {
 //parameters: SessionPoolHandle CallingNumber
 //AgtPppSessionPoolL2tp SetCallingNumber
 return nil
}

func(np *PppSessionPoolL2tp) GetHelloInterval ()(string, error) {
 //parameters: SessionPoolHandle Interval
 //AgtPppSessionPoolL2tp GetHelloInterval
 return "", nil
}

func(np *PppSessionPoolL2tp) GetHostName ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetHostName
 return "", nil
}

func(np *PppSessionPoolL2tp) GetIdentity ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetIdentity
 return "", nil
}

func(np *PppSessionPoolL2tp) GetInitialAckTimeout ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetInitialAckTimeout
 return "", nil
}

func(np *PppSessionPoolL2tp) GetMaxTxRetries ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetMaxTxRetries
 return "", nil
}

func(np *PppSessionPoolL2tp) GetNumberOfCallsPerTunnel ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetNumberOfCallsPerTunnel
 return "", nil
}

func(np *PppSessionPoolL2tp) GetReceiveWindowSize ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetReceiveWindowSize
 return "", nil
}

func(np *PppSessionPoolL2tp) GetRxConnectSpeed ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetRxConnectSpeed
 return "", nil
}

func(np *PppSessionPoolL2tp) SetRxConnectSpeed () error {
 //parameters: SessionPoolHandle RxSpeed
 //AgtPppSessionPoolL2tp SetRxConnectSpeed
 return nil
}

func(np *PppSessionPoolL2tp) GetSessionIds ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetSessionIds
 return "", nil
}

func(np *PppSessionPoolL2tp) GetSharedSecret ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetSharedSecret
 return "", nil
}

func(np *PppSessionPoolL2tp) GetSourceUdpPort ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetSourceUdpPort
 return "", nil
}

func(np *PppSessionPoolL2tp) GetSutIpAddress ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetSutIpAddress
 return "", nil
}

func(np *PppSessionPoolL2tp) GetTesterIpAddress ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetTesterIpAddress
 return "", nil
}

func(np *PppSessionPoolL2tp) GetTunnelDestinationPorts ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetTunnelDestinationPorts
 return "", nil
}

func(np *PppSessionPoolL2tp) GetTunnelIds ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetTunnelIds
 return "", nil
}

func(np *PppSessionPoolL2tp) GetTxConnectSpeed ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetTxConnectSpeed
 return "", nil
}

func(np *PppSessionPoolL2tp) SetTxConnectSpeed () error {
 //parameters: SessionPoolHandle TxSpeed
 //AgtPppSessionPoolL2tp SetTxConnectSpeed
 return nil
}

func(np *PppSessionPoolL2tp) Enable () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp Enable
 return nil
}

func(np *PppSessionPoolL2tp) Disable () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp Disable
 return nil
}

func(np *PppSessionPoolL2tp) SetNumberOfSessions () error {
 //parameters: SessionPoolHandle NumSessions
 //AgtPppSessionPoolL2tp SetNumberOfSessions
 return nil
}

func(np *PppSessionPoolL2tp) SetSessionLifetime () error {
 //parameters: SessionPoolHandle MinLifetime MaxLifetime
 //AgtPppSessionPoolL2tp SetSessionLifetime
 return nil
}

func(np *PppSessionPoolL2tp) SetSessionInitiationRate () error {
 //parameters: SessionPoolHandle NumSessions Period
 //AgtPppSessionPoolL2tp SetSessionInitiationRate
 return nil
}

func(np *PppSessionPoolL2tp) SetLcpOption () error {
 //parameters: SessionPoolHandle LcpOption Value
 //AgtPppSessionPoolL2tp SetLcpOption
 return nil
}

func(np *PppSessionPoolL2tp) EnableAccm () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp EnableAccm
 return nil
}

func(np *PppSessionPoolL2tp) DisableAccm () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp DisableAccm
 return nil
}

func(np *PppSessionPoolL2tp) EnableAddressControlFieldCompression () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp EnableAddressControlFieldCompression
 return nil
}

func(np *PppSessionPoolL2tp) DisableAddressControlFieldCompression () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp DisableAddressControlFieldCompression
 return nil
}

func(np *PppSessionPoolL2tp) SetMode () error {
 //parameters: SessionPoolHandle Mode
 //AgtPppSessionPoolL2tp SetMode
 return nil
}

func(np *PppSessionPoolL2tp) SetAuthenticationProtocol () error {
 //parameters: SessionPoolHandle Protocol
 //AgtPppSessionPoolL2tp SetAuthenticationProtocol
 return nil
}

func(np *PppSessionPoolL2tp) SetAuthenticationParameter () error {
 //parameters: SessionPoolHandle Parameter Value
 //AgtPppSessionPoolL2tp SetAuthenticationParameter
 return nil
}

func(np *PppSessionPoolL2tp) EnableUniqueAuthenticationParameter () error {
 //parameters: SessionPoolHandle Parameter
 //AgtPppSessionPoolL2tp EnableUniqueAuthenticationParameter
 return nil
}

func(np *PppSessionPoolL2tp) DisableUniqueAuthenticationParameter () error {
 //parameters: SessionPoolHandle Parameter
 //AgtPppSessionPoolL2tp DisableUniqueAuthenticationParameter
 return nil
}

func(np *PppSessionPoolL2tp) SetAuthenticationParameterStartValue () error {
 //parameters: SessionPoolHandle Parameter StartValue
 //AgtPppSessionPoolL2tp SetAuthenticationParameterStartValue
 return nil
}

func(np *PppSessionPoolL2tp) GetAuthenticationParameter ()(string, error) {
 //parameters: SessionPoolHandle Parameter
 //AgtPppSessionPoolL2tp GetAuthenticationParameter
 return "", nil
}

func(np *PppSessionPoolL2tp) SetPeerAuthenticationParameter () error {
 //parameters: SessionPoolHandle PeerParameter Value
 //AgtPppSessionPoolL2tp SetPeerAuthenticationParameter
 return nil
}

func(np *PppSessionPoolL2tp) GetPeerAuthenticationParameter ()(string, error) {
 //parameters: SessionPoolHandle PeerParameter
 //AgtPppSessionPoolL2tp GetPeerAuthenticationParameter
 return "", nil
}

func(np *PppSessionPoolL2tp) EnableAcceptAllAuthenticationValues () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp EnableAcceptAllAuthenticationValues
 return nil
}

func(np *PppSessionPoolL2tp) DisableAcceptAllAuthenticationValues () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp DisableAcceptAllAuthenticationValues
 return nil
}

func(np *PppSessionPoolL2tp) IsAcceptAllAuthenticationValuesEnabled () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp IsAcceptAllAuthenticationValuesEnabled
 return nil
}

func(np *PppSessionPoolL2tp) SetUserToDomainRatio () error {
 //parameters: SessionPoolHandle UserToDomainRatio
 //AgtPppSessionPoolL2tp SetUserToDomainRatio
 return nil
}

func(np *PppSessionPoolL2tp) GetUserToDomainRatio ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetUserToDomainRatio
 return "", nil
}

func(np *PppSessionPoolL2tp) SetPasswordIncrementMethod () error {
 //parameters: SessionPoolHandle PasswordIncrementMethod
 //AgtPppSessionPoolL2tp SetPasswordIncrementMethod
 return nil
}

func(np *PppSessionPoolL2tp) GetPasswordIncrementMethod ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetPasswordIncrementMethod
 return "", nil
}

func(np *PppSessionPoolL2tp) SetDomainCountLimit () error {
 //parameters: SessionPoolHandle DomainCountLimit
 //AgtPppSessionPoolL2tp SetDomainCountLimit
 return nil
}

func(np *PppSessionPoolL2tp) GetDomainCountLimit ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetDomainCountLimit
 return "", nil
}

func(np *PppSessionPoolL2tp) EnableDomainNameIncrementAfterExtension () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp EnableDomainNameIncrementAfterExtension
 return nil
}

func(np *PppSessionPoolL2tp) DisableDomainNameIncrementAfterExtension () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp DisableDomainNameIncrementAfterExtension
 return nil
}

func(np *PppSessionPoolL2tp) IsDomainNameIncrementAfterExtensionEnabled () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp IsDomainNameIncrementAfterExtensionEnabled
 return nil
}

func(np *PppSessionPoolL2tp) SetAcceptedAuthenticationParameter () error {
 //parameters: SessionPoolHandle AcceptedParameter Value
 //AgtPppSessionPoolL2tp SetAcceptedAuthenticationParameter
 return nil
}

func(np *PppSessionPoolL2tp) GetAcceptedAuthenticationParameter ()(string, error) {
 //parameters: SessionPoolHandle AcceptedParameter
 //AgtPppSessionPoolL2tp GetAcceptedAuthenticationParameter
 return "", nil
}

func(np *PppSessionPoolL2tp) SetStartingSourceIpAddress () error {
 //parameters: SessionPoolHandle IpAddress
 //AgtPppSessionPoolL2tp SetStartingSourceIpAddress
 return nil
}

func(np *PppSessionPoolL2tp) SetSourceIpAddressModifier () error {
 //parameters: SessionPoolHandle Modifier
 //AgtPppSessionPoolL2tp SetSourceIpAddressModifier
 return nil
}

func(np *PppSessionPoolL2tp) SetMaxConnectionAttempts () error {
 //parameters: SessionPoolHandle MaxConnectionAttempts
 //AgtPppSessionPoolL2tp SetMaxConnectionAttempts
 return nil
}

func(np *PppSessionPoolL2tp) GetMaxConnectionAttempts ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetMaxConnectionAttempts
 return "", nil
}

func(np *PppSessionPoolL2tp) RetryConnections () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp RetryConnections
 return nil
}

func(np *PppSessionPoolL2tp) GetConnectionAttempts ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetConnectionAttempts
 return "", nil
}

func(np *PppSessionPoolL2tp) Cancel () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp Cancel
 return nil
}

func(np *PppSessionPoolL2tp) GetAccmFlag ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetAccmFlag
 return "", nil
}

func(np *PppSessionPoolL2tp) GetAddressControlFieldCompressionFlag ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetAddressControlFieldCompressionFlag
 return "", nil
}

func(np *PppSessionPoolL2tp) GetAuthenticationParameterStartValue ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetAuthenticationParameterStartValue
 return "", nil
}

func(np *PppSessionPoolL2tp) GetAuthenticationProtocol ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetAuthenticationProtocol
 return "", nil
}

func(np *PppSessionPoolL2tp) GetEnableFlag ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetEnableFlag
 return "", nil
}

func(np *PppSessionPoolL2tp) GetLcpOption ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetLcpOption
 return "", nil
}

func(np *PppSessionPoolL2tp) GetMode ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetMode
 return "", nil
}

func(np *PppSessionPoolL2tp) GetName ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetName
 return "", nil
}

func(np *PppSessionPoolL2tp) SetName () error {
 //parameters: SessionPoolHandle Name
 //AgtPppSessionPoolL2tp SetName
 return nil
}

func(np *PppSessionPoolL2tp) GetNegotiatedIpAddresses ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetNegotiatedIpAddresses
 return "", nil
}

func(np *PppSessionPoolL2tp) GetNumberOfEstablishedSessions ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetNumberOfEstablishedSessions
 return "", nil
}

func(np *PppSessionPoolL2tp) GetNumberOfSessions ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetNumberOfSessions
 return "", nil
}

func(np *PppSessionPoolL2tp) GetRemoteIpAddresses ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetRemoteIpAddresses
 return "", nil
}

func(np *PppSessionPoolL2tp) GetSessionInitiationRate ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetSessionInitiationRate
 return "", nil
}

func(np *PppSessionPoolL2tp) GetSessionLifetime ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetSessionLifetime
 return "", nil
}

func(np *PppSessionPoolL2tp) GetSessionStateDescriptions ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetSessionStateDescriptions
 return "", nil
}

func(np *PppSessionPoolL2tp) GetSourceIpAddressModifier ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetSourceIpAddressModifier
 return "", nil
}

func(np *PppSessionPoolL2tp) GetStartingSourceIpAddress ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetStartingSourceIpAddress
 return "", nil
}

func(np *PppSessionPoolL2tp) GetState ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolL2tp GetState
 return "", nil
}

func(np *PppSessionPoolL2tp) GetStateSummary ()(string, error) {
 //parameters: 
 //AgtPppSessionPoolL2tp GetStateSummary
 return "", nil
}

func(np *PppSessionPoolL2tp) GetUniqueAuthenticationParameterFlag ()(string, error) {
 //parameters: SessionPoolHandle Parameter
 //AgtPppSessionPoolL2tp GetUniqueAuthenticationParameterFlag
 return "", nil
}

func(np *PppSessionPoolL2tp) EnableAllSessions () error {
 //parameters: 
 //AgtPppSessionPoolL2tp EnableAllSessions
 return nil
}

func(np *PppSessionPoolL2tp) DisableAllSessions () error {
 //parameters: 
 //AgtPppSessionPoolL2tp DisableAllSessions
 return nil
}

