package n2xsdk

type SipUserAgentPool struct {
 Handler string
}

func(np *SipUserAgentPool) SetOutboundProxyIpv4Address () error {
 //parameters: DeviceHandle Ipv4Address Increment Repeat
 //AgtSipUserAgentPool SetOutboundProxyIpv4Address
 return nil
}

func(np *SipUserAgentPool) GetOutboundProxyIpv4Address ()(string, error) {
 //parameters: DeviceHandle
 //AgtSipUserAgentPool GetOutboundProxyIpv4Address
 return "", nil
}

func(np *SipUserAgentPool) IsNetworkAddressTranslationEnabled () error {
 //parameters: DeviceHandle
 //AgtSipUserAgentPool IsNetworkAddressTranslationEnabled
 return nil
}

func(np *SipUserAgentPool) EnableNetworkAddressTranslation () error {
 //parameters: DeviceHandle
 //AgtSipUserAgentPool EnableNetworkAddressTranslation
 return nil
}

func(np *SipUserAgentPool) DisableNetworkAddressTranslation () error {
 //parameters: DeviceHandle
 //AgtSipUserAgentPool DisableNetworkAddressTranslation
 return nil
}

func(np *SipUserAgentPool) SetNatedInternaIpv4AddressRange () error {
 //parameters: DeviceHandle Ipv4Address Increment Repeat
 //AgtSipUserAgentPool SetNatedInternaIpv4AddressRange
 return nil
}

func(np *SipUserAgentPool) GetNatedInternaIpv4AddressRange ()(string, error) {
 //parameters: DeviceHandle
 //AgtSipUserAgentPool GetNatedInternaIpv4AddressRange
 return "", nil
}

func(np *SipUserAgentPool) SetNatedExternalIpv4AddressRange () error {
 //parameters: DeviceHandle Ipv4Address Increment Repeat
 //AgtSipUserAgentPool SetNatedExternalIpv4AddressRange
 return nil
}

func(np *SipUserAgentPool) GetNatedExternalIpv4AddressRange ()(string, error) {
 //parameters: DeviceHandle
 //AgtSipUserAgentPool GetNatedExternalIpv4AddressRange
 return "", nil
}

func(np *SipUserAgentPool) EstablishSession () error {
 //parameters: DeviceOrSession
 //AgtSipUserAgentPool EstablishSession
 return nil
}

func(np *SipUserAgentPool) TerminateSession () error {
 //parameters: DeviceOrSession
 //AgtSipUserAgentPool TerminateSession
 return nil
}

func(np *SipUserAgentPool) RegisterSession () error {
 //parameters: DeviceOrSession
 //AgtSipUserAgentPool RegisterSession
 return nil
}

func(np *SipUserAgentPool) DeregisterSession () error {
 //parameters: DeviceOrSession
 //AgtSipUserAgentPool DeregisterSession
 return nil
}

func(np *SipUserAgentPool) SetLocalUsernamePrefix () error {
 //parameters: EmulationHandle LocalUserNamePrefix
 //AgtSipUserAgentPool SetLocalUsernamePrefix
 return nil
}

func(np *SipUserAgentPool) SetLocalUsernameSuffixIncrementingRange () error {
 //parameters: EmulationHandle LocalUserNameSuffix Increment Repeat
 //AgtSipUserAgentPool SetLocalUsernameSuffixIncrementingRange
 return nil
}

func(np *SipUserAgentPool) SetLocalDomain () error {
 //parameters: EmulationHandle LocalDomain
 //AgtSipUserAgentPool SetLocalDomain
 return nil
}

func(np *SipUserAgentPool) GetLocalDomain ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetLocalDomain
 return "", nil
}

func(np *SipUserAgentPool) GetLocalUsernamePrefix ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetLocalUsernamePrefix
 return "", nil
}

func(np *SipUserAgentPool) GetLocalUsernameSuffixIncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetLocalUsernameSuffixIncrementingRange
 return "", nil
}

func(np *SipUserAgentPool) IsUseRemoteAddressOfRecordEnabled () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool IsUseRemoteAddressOfRecordEnabled
 return nil
}

func(np *SipUserAgentPool) EnableUseRemoteAddressOfRecord () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool EnableUseRemoteAddressOfRecord
 return nil
}

func(np *SipUserAgentPool) DisableUseRemoteAddressOfRecord () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool DisableUseRemoteAddressOfRecord
 return nil
}

func(np *SipUserAgentPool) SetRemoteUsernamePrefix () error {
 //parameters: EmulationHandle RemoteUserNamePrefix
 //AgtSipUserAgentPool SetRemoteUsernamePrefix
 return nil
}

func(np *SipUserAgentPool) SetRemoteUsernameSuffixIncrementingRange () error {
 //parameters: EmulationHandle RemoteUserNameSuffix Increment Repeat
 //AgtSipUserAgentPool SetRemoteUsernameSuffixIncrementingRange
 return nil
}

func(np *SipUserAgentPool) GetRemoteUsernameSuffixIncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetRemoteUsernameSuffixIncrementingRange
 return "", nil
}

func(np *SipUserAgentPool) SetRemoteDomain () error {
 //parameters: EmulationHandle RemoteDomain
 //AgtSipUserAgentPool SetRemoteDomain
 return nil
}

func(np *SipUserAgentPool) GetRemoteUsernamePrefix ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetRemoteUsernamePrefix
 return "", nil
}

func(np *SipUserAgentPool) GetRemoteDomain ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetRemoteDomain
 return "", nil
}

func(np *SipUserAgentPool) SetRemoteHostAddressV4IncrementingRange () error {
 //parameters: EmulationHandle Ipv4Address Increment Repeat
 //AgtSipUserAgentPool SetRemoteHostAddressV4IncrementingRange
 return nil
}

func(np *SipUserAgentPool) SetRemoteHostL3AddressFamily () error {
 //parameters: EmulationHandle RegistrarAddressType
 //AgtSipUserAgentPool SetRemoteHostL3AddressFamily
 return nil
}

func(np *SipUserAgentPool) GetRemoteHostL3AddressFamily ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetRemoteHostL3AddressFamily
 return "", nil
}

func(np *SipUserAgentPool) GetRemoteHostAddressV4IncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetRemoteHostAddressV4IncrementingRange
 return "", nil
}

func(np *SipUserAgentPool) SetRemoteHostAddressV6IncrementingRange () error {
 //parameters: EmulationHandle Ipv6Address Increment Repeat
 //AgtSipUserAgentPool SetRemoteHostAddressV6IncrementingRange
 return nil
}

func(np *SipUserAgentPool) GetRemoteHostAddressV6IncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetRemoteHostAddressV6IncrementingRange
 return "", nil
}

func(np *SipUserAgentPool) SetMediaPayloadType () error {
 //parameters: EmulationHandle MediaPayloadType
 //AgtSipUserAgentPool SetMediaPayloadType
 return nil
}

func(np *SipUserAgentPool) GetMediaPayloadType ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetMediaPayloadType
 return "", nil
}

func(np *SipUserAgentPool) SetMediaPort () error {
 //parameters: EmulationHandle Port
 //AgtSipUserAgentPool SetMediaPort
 return nil
}

func(np *SipUserAgentPool) GetMediaPort ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetMediaPort
 return "", nil
}

func(np *SipUserAgentPool) SetSipPort () error {
 //parameters: EmulationHandle Port
 //AgtSipUserAgentPool SetSipPort
 return nil
}

func(np *SipUserAgentPool) GetSipPort ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetSipPort
 return "", nil
}

func(np *SipUserAgentPool) IsRegistrarEnabled () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool IsRegistrarEnabled
 return nil
}

func(np *SipUserAgentPool) EnableRegistrar () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool EnableRegistrar
 return nil
}

func(np *SipUserAgentPool) DisableRegistrar () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool DisableRegistrar
 return nil
}

func(np *SipUserAgentPool) SetRegistrarIpv4Address () error {
 //parameters: EmulationHandle Ipv4Address
 //AgtSipUserAgentPool SetRegistrarIpv4Address
 return nil
}

func(np *SipUserAgentPool) GetRegistrarIpv4Address ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetRegistrarIpv4Address
 return "", nil
}

func(np *SipUserAgentPool) SetRegistrarIpv6Address () error {
 //parameters: EmulationHandle Ipv6Address
 //AgtSipUserAgentPool SetRegistrarIpv6Address
 return nil
}

func(np *SipUserAgentPool) GetRegistrarIpv6Address ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetRegistrarIpv6Address
 return "", nil
}

func(np *SipUserAgentPool) SetRegistrarL3AddressFamily () error {
 //parameters: EmulationHandle L3AddressFamily
 //AgtSipUserAgentPool SetRegistrarL3AddressFamily
 return nil
}

func(np *SipUserAgentPool) GetRegistrarL3AddressFamily ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetRegistrarL3AddressFamily
 return "", nil
}

func(np *SipUserAgentPool) SetRegistrarHostname () error {
 //parameters: EmulationHandle RegistrarHostname
 //AgtSipUserAgentPool SetRegistrarHostname
 return nil
}

func(np *SipUserAgentPool) GetRegistrarHostname ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetRegistrarHostname
 return "", nil
}

func(np *SipUserAgentPool) IsAutomaticRegistrationEnabled () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool IsAutomaticRegistrationEnabled
 return nil
}

func(np *SipUserAgentPool) EnableAutomaticRegistration () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool EnableAutomaticRegistration
 return nil
}

func(np *SipUserAgentPool) DisableAutomaticRegistration () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool DisableAutomaticRegistration
 return nil
}

func(np *SipUserAgentPool) SetDesiredRegistrationExpiryTime () error {
 //parameters: EmulationHandle DesiredRegistrationExpiryTime
 //AgtSipUserAgentPool SetDesiredRegistrationExpiryTime
 return nil
}

func(np *SipUserAgentPool) GetDesiredRegistrationExpiryTime ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetDesiredRegistrationExpiryTime
 return "", nil
}

func(np *SipUserAgentPool) SetMaxRegistrationAttempts () error {
 //parameters: EmulationHandle MaxRegistrationAttempts
 //AgtSipUserAgentPool SetMaxRegistrationAttempts
 return nil
}

func(np *SipUserAgentPool) GetMaxRegistrationAttempts ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetMaxRegistrationAttempts
 return "", nil
}

func(np *SipUserAgentPool) SetRegistrationAttemptDelayTime () error {
 //parameters: EmulationHandle RegistrationAttemptDelayTime
 //AgtSipUserAgentPool SetRegistrationAttemptDelayTime
 return nil
}

func(np *SipUserAgentPool) GetRegistrationAttemptDelayTime ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetRegistrationAttemptDelayTime
 return "", nil
}

func(np *SipUserAgentPool) SetTransactionTimerT1 () error {
 //parameters: EmulationHandle Time
 //AgtSipUserAgentPool SetTransactionTimerT1
 return nil
}

func(np *SipUserAgentPool) SetTransactionTimerT2 () error {
 //parameters: EmulationHandle Time
 //AgtSipUserAgentPool SetTransactionTimerT2
 return nil
}

func(np *SipUserAgentPool) SetTransactionTimerT4 () error {
 //parameters: EmulationHandle Time
 //AgtSipUserAgentPool SetTransactionTimerT4
 return nil
}

func(np *SipUserAgentPool) GetTransactionTimerT1 ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetTransactionTimerT1
 return "", nil
}

func(np *SipUserAgentPool) GetTransactionTimerT2 ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetTransactionTimerT2
 return "", nil
}

func(np *SipUserAgentPool) GetTransactionTimerT4 ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetTransactionTimerT4
 return "", nil
}

func(np *SipUserAgentPool) IsResourcePriorityEnabled () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool IsResourcePriorityEnabled
 return nil
}

func(np *SipUserAgentPool) EnableResourcePriority () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool EnableResourcePriority
 return nil
}

func(np *SipUserAgentPool) DisableResourcePriority () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool DisableResourcePriority
 return nil
}

func(np *SipUserAgentPool) SetResourcePriority () error {
 //parameters: EmulationHandle ResourcePriority
 //AgtSipUserAgentPool SetResourcePriority
 return nil
}

func(np *SipUserAgentPool) GetResourcePriority ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetResourcePriority
 return "", nil
}

func(np *SipUserAgentPool) IsSessionRefreshTimerEnabled () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool IsSessionRefreshTimerEnabled
 return nil
}

func(np *SipUserAgentPool) EnableSessionRefreshTimer () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool EnableSessionRefreshTimer
 return nil
}

func(np *SipUserAgentPool) DisableSessionRefreshTimer () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool DisableSessionRefreshTimer
 return nil
}

func(np *SipUserAgentPool) SetDesiredSessionInterval () error {
 //parameters: EmulationHandle DesiredSessionInterval
 //AgtSipUserAgentPool SetDesiredSessionInterval
 return nil
}

func(np *SipUserAgentPool) GetDesiredSessionInterval ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetDesiredSessionInterval
 return "", nil
}

func(np *SipUserAgentPool) GetMinimumSessionInterval ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetMinimumSessionInterval
 return "", nil
}

func(np *SipUserAgentPool) SetMinimumSessionInterval () error {
 //parameters: EmulationHandle MinimumSessionInterval
 //AgtSipUserAgentPool SetMinimumSessionInterval
 return nil
}

func(np *SipUserAgentPool) SetRoundTripDelay () error {
 //parameters: EmulationHandle RoundTripInterval
 //AgtSipUserAgentPool SetRoundTripDelay
 return nil
}

func(np *SipUserAgentPool) GetRoundTripDelay ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetRoundTripDelay
 return "", nil
}

func(np *SipUserAgentPool) SetSessionRefresher () error {
 //parameters: EmulationHandle UserAgentType
 //AgtSipUserAgentPool SetSessionRefresher
 return nil
}

func(np *SipUserAgentPool) GetSessionRefresher ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetSessionRefresher
 return "", nil
}

func(np *SipUserAgentPool) IsCallAcceptDelayEnabled () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool IsCallAcceptDelayEnabled
 return nil
}

func(np *SipUserAgentPool) EnableCallAcceptDelay () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool EnableCallAcceptDelay
 return nil
}

func(np *SipUserAgentPool) DisableCallAcceptDelay () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool DisableCallAcceptDelay
 return nil
}

func(np *SipUserAgentPool) SetResponseDelayInterval () error {
 //parameters: EmulationHandle Time
 //AgtSipUserAgentPool SetResponseDelayInterval
 return nil
}

func(np *SipUserAgentPool) GetResponseDelayInterval ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetResponseDelayInterval
 return "", nil
}

func(np *SipUserAgentPool) IsPrivacyHeaderEnabled () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool IsPrivacyHeaderEnabled
 return nil
}

func(np *SipUserAgentPool) EnablePrivacyHeader () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool EnablePrivacyHeader
 return nil
}

func(np *SipUserAgentPool) DisablePrivacyHeader () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool DisablePrivacyHeader
 return nil
}

func(np *SipUserAgentPool) SetPrivacyType () error {
 //parameters: EmulationHandle PrivacyType
 //AgtSipUserAgentPool SetPrivacyType
 return nil
}

func(np *SipUserAgentPool) GetPrivacyType ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetPrivacyType
 return "", nil
}

func(np *SipUserAgentPool) IsAuthenticationEnabled () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool IsAuthenticationEnabled
 return nil
}

func(np *SipUserAgentPool) EnableAuthentication () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool EnableAuthentication
 return nil
}

func(np *SipUserAgentPool) DisableAuthentication () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool DisableAuthentication
 return nil
}

func(np *SipUserAgentPool) IsAuthenticationCredentialEnabled () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool IsAuthenticationCredentialEnabled
 return nil
}

func(np *SipUserAgentPool) EnableAuthenticationCredential () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool EnableAuthenticationCredential
 return nil
}

func(np *SipUserAgentPool) DisableAuthenticationCredential () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool DisableAuthenticationCredential
 return nil
}

func(np *SipUserAgentPool) IsCredentialInDialogEnabled () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool IsCredentialInDialogEnabled
 return nil
}

func(np *SipUserAgentPool) EnableCredentialInDialog () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool EnableCredentialInDialog
 return nil
}

func(np *SipUserAgentPool) DisableCredentialInDialog () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool DisableCredentialInDialog
 return nil
}

func(np *SipUserAgentPool) SetAuthCredentialUsernamePrefix () error {
 //parameters: EmulationHandle AuthCredentialUsernamePrefix
 //AgtSipUserAgentPool SetAuthCredentialUsernamePrefix
 return nil
}

func(np *SipUserAgentPool) GetAuthCredentialUsernamePrefix ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetAuthCredentialUsernamePrefix
 return "", nil
}

func(np *SipUserAgentPool) SetAuthCredentialUsernameSuffixIncrementingRange () error {
 //parameters: EmulationHandle AuthCredentialUsernameSuffix Increment Repeat
 //AgtSipUserAgentPool SetAuthCredentialUsernameSuffixIncrementingRange
 return nil
}

func(np *SipUserAgentPool) GetAuthCredentialUsernameSuffixIncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetAuthCredentialUsernameSuffixIncrementingRange
 return "", nil
}

func(np *SipUserAgentPool) SetAuthCredentialPassword () error {
 //parameters: EmulationHandle AuthCredentialPassword
 //AgtSipUserAgentPool SetAuthCredentialPassword
 return nil
}

func(np *SipUserAgentPool) GetAuthCredentialPassword ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetAuthCredentialPassword
 return "", nil
}

func(np *SipUserAgentPool) SetAuthCredentialPasswordModifier () error {
 //parameters: EmulationHandle AuthCredentialPasswordModifier Increment Repeat
 //AgtSipUserAgentPool SetAuthCredentialPasswordModifier
 return nil
}

func(np *SipUserAgentPool) GetAuthCredentialPasswordModifier ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetAuthCredentialPasswordModifier
 return "", nil
}

func(np *SipUserAgentPool) IsAuthenticationChallengeEnabled () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool IsAuthenticationChallengeEnabled
 return nil
}

func(np *SipUserAgentPool) EnableAuthenticationChallenge () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool EnableAuthenticationChallenge
 return nil
}

func(np *SipUserAgentPool) DisableAuthenticationChallenge () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool DisableAuthenticationChallenge
 return nil
}

func(np *SipUserAgentPool) GetChallengeAuthenticationScheme ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetChallengeAuthenticationScheme
 return "", nil
}

func(np *SipUserAgentPool) SetChallengeAuthenticationRealm () error {
 //parameters: EmulationHandle ChallengeAuthenticationRealm
 //AgtSipUserAgentPool SetChallengeAuthenticationRealm
 return nil
}

func(np *SipUserAgentPool) GetChallengeAuthenticationRealm ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetChallengeAuthenticationRealm
 return "", nil
}

func(np *SipUserAgentPool) SetAuthChallengeUsernamePrefix () error {
 //parameters: EmulationHandle AuthChallengeUsernamePrefix
 //AgtSipUserAgentPool SetAuthChallengeUsernamePrefix
 return nil
}

func(np *SipUserAgentPool) GetAuthChallengeUsernamePrefix ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetAuthChallengeUsernamePrefix
 return "", nil
}

func(np *SipUserAgentPool) SetAuthChallengeUsernameSuffixIncrementingRange () error {
 //parameters: EmulationHandle AuthChallengeUsernameSuffix Increment Repeat
 //AgtSipUserAgentPool SetAuthChallengeUsernameSuffixIncrementingRange
 return nil
}

func(np *SipUserAgentPool) GetAuthChallengeUsernameSuffixIncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetAuthChallengeUsernameSuffixIncrementingRange
 return "", nil
}

func(np *SipUserAgentPool) SetAuthChallengePassword () error {
 //parameters: EmulationHandle AuthChallengePassword
 //AgtSipUserAgentPool SetAuthChallengePassword
 return nil
}

func(np *SipUserAgentPool) GetAuthChallengePassword ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetAuthChallengePassword
 return "", nil
}

func(np *SipUserAgentPool) SetAuthChallengePasswordModifier () error {
 //parameters: EmulationHandle AuthChallengePasswordModifier Increment Repeat
 //AgtSipUserAgentPool SetAuthChallengePasswordModifier
 return nil
}

func(np *SipUserAgentPool) GetAuthChallengePasswordModifier ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetAuthChallengePasswordModifier
 return "", nil
}

func(np *SipUserAgentPool) SetChallengeOpaque () error {
 //parameters: EmulationHandle ChallengeOpaque
 //AgtSipUserAgentPool SetChallengeOpaque
 return nil
}

func(np *SipUserAgentPool) GetChallengeOpaque ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetChallengeOpaque
 return "", nil
}

func(np *SipUserAgentPool) SetChallengeAlgorithm () error {
 //parameters: EmulationHandle ChallengeAlgorithm
 //AgtSipUserAgentPool SetChallengeAlgorithm
 return nil
}

func(np *SipUserAgentPool) GetChallengeAlgorithm ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetChallengeAlgorithm
 return "", nil
}

func(np *SipUserAgentPool) SetChallengeQop () error {
 //parameters: EmulationHandle ChallengeQop
 //AgtSipUserAgentPool SetChallengeQop
 return nil
}

func(np *SipUserAgentPool) GetChallengeQop ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetChallengeQop
 return "", nil
}

func(np *SipUserAgentPool) SetNonceExpireTime () error {
 //parameters: EmulationHandle NonceExpireTime
 //AgtSipUserAgentPool SetNonceExpireTime
 return nil
}

func(np *SipUserAgentPool) GetNonceExpireTime ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetNonceExpireTime
 return "", nil
}

func(np *SipUserAgentPool) SetMaximumChallengeAttempts () error {
 //parameters: EmulationHandle MaximumChallengeAttempts
 //AgtSipUserAgentPool SetMaximumChallengeAttempts
 return nil
}

func(np *SipUserAgentPool) GetMaximumChallengeAttempts ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool GetMaximumChallengeAttempts
 return "", nil
}

func(np *SipUserAgentPool) SendTrafficParams () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool SendTrafficParams
 return nil
}

func(np *SipUserAgentPool) ResetTrafficParams () error {
 //parameters: EmulationHandle
 //AgtSipUserAgentPool ResetTrafficParams
 return nil
}

