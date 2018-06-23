package n2xsdk

type SipSubscriberPool struct {
 Handler string
}

func(np *SipSubscriberPool) Enable () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool Enable
 return nil
}

func(np *SipSubscriberPool) Disable () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool Disable
 return nil
}

func(np *SipSubscriberPool) SetSubscriberPerGatewayCount () error {
 //parameters: EmulationHandle SipSubscriberPerGateway
 //AgtSipSubscriberPool SetSubscriberPerGatewayCount
 return nil
}

func(np *SipSubscriberPool) GetSubscriberPerGatewayCount ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetSubscriberPerGatewayCount
 return "", nil
}

func(np *SipSubscriberPool) GetSubscriberPoolCount ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetSubscriberPoolCount
 return "", nil
}

func(np *SipSubscriberPool) GetGatewayPoolCount ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetGatewayPoolCount
 return "", nil
}

func(np *SipSubscriberPool) GetCustomStateCount ()(string, error) {
 //parameters: EmulationHandle CustomState
 //AgtSipSubscriberPool GetCustomStateCount
 return "", nil
}

func(np *SipSubscriberPool) EstablishSession () error {
 //parameters: DeviceOrSession
 //AgtSipSubscriberPool EstablishSession
 return nil
}

func(np *SipSubscriberPool) TerminateSession () error {
 //parameters: DeviceOrSession
 //AgtSipSubscriberPool TerminateSession
 return nil
}

func(np *SipSubscriberPool) RegisterSession () error {
 //parameters: DeviceOrSession
 //AgtSipSubscriberPool RegisterSession
 return nil
}

func(np *SipSubscriberPool) DeregisterSession () error {
 //parameters: DeviceOrSession
 //AgtSipSubscriberPool DeregisterSession
 return nil
}

func(np *SipSubscriberPool) SetLocalUsernamePrefix () error {
 //parameters: EmulationHandle LocalUserNamePrefix
 //AgtSipSubscriberPool SetLocalUsernamePrefix
 return nil
}

func(np *SipSubscriberPool) SetLocalUsernameSuffixIncrementingRange () error {
 //parameters: EmulationHandle LocalUserNameSuffix Increment Repeat
 //AgtSipSubscriberPool SetLocalUsernameSuffixIncrementingRange
 return nil
}

func(np *SipSubscriberPool) SetLocalDomain () error {
 //parameters: EmulationHandle LocalDomain
 //AgtSipSubscriberPool SetLocalDomain
 return nil
}

func(np *SipSubscriberPool) GetLocalDomain ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetLocalDomain
 return "", nil
}

func(np *SipSubscriberPool) GetLocalUsernamePrefix ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetLocalUsernamePrefix
 return "", nil
}

func(np *SipSubscriberPool) GetLocalUsernameSuffixIncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetLocalUsernameSuffixIncrementingRange
 return "", nil
}

func(np *SipSubscriberPool) IsUseRemoteAddressOfRecordEnabled () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool IsUseRemoteAddressOfRecordEnabled
 return nil
}

func(np *SipSubscriberPool) EnableUseRemoteAddressOfRecord () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool EnableUseRemoteAddressOfRecord
 return nil
}

func(np *SipSubscriberPool) DisableUseRemoteAddressOfRecord () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool DisableUseRemoteAddressOfRecord
 return nil
}

func(np *SipSubscriberPool) SetRemoteUsernamePrefix () error {
 //parameters: EmulationHandle RemoteUserNamePrefix
 //AgtSipSubscriberPool SetRemoteUsernamePrefix
 return nil
}

func(np *SipSubscriberPool) SetRemoteUsernameSuffixIncrementingRange () error {
 //parameters: EmulationHandle RemoteUserNameSuffix Increment Repeat
 //AgtSipSubscriberPool SetRemoteUsernameSuffixIncrementingRange
 return nil
}

func(np *SipSubscriberPool) GetRemoteUsernameSuffixIncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetRemoteUsernameSuffixIncrementingRange
 return "", nil
}

func(np *SipSubscriberPool) SetRemoteDomain () error {
 //parameters: EmulationHandle RemoteDomain
 //AgtSipSubscriberPool SetRemoteDomain
 return nil
}

func(np *SipSubscriberPool) GetRemoteUsernamePrefix ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetRemoteUsernamePrefix
 return "", nil
}

func(np *SipSubscriberPool) GetRemoteDomain ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetRemoteDomain
 return "", nil
}

func(np *SipSubscriberPool) SetRemoteHostAddressV4IncrementingRange () error {
 //parameters: EmulationHandle Ipv4Address Increment Repeat
 //AgtSipSubscriberPool SetRemoteHostAddressV4IncrementingRange
 return nil
}

func(np *SipSubscriberPool) SetRemoteHostL3AddressFamily () error {
 //parameters: EmulationHandle RegistrarAddressType
 //AgtSipSubscriberPool SetRemoteHostL3AddressFamily
 return nil
}

func(np *SipSubscriberPool) GetRemoteHostL3AddressFamily ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetRemoteHostL3AddressFamily
 return "", nil
}

func(np *SipSubscriberPool) GetRemoteHostAddressV4IncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetRemoteHostAddressV4IncrementingRange
 return "", nil
}

func(np *SipSubscriberPool) SetRemoteHostAddressV6IncrementingRange () error {
 //parameters: EmulationHandle Ipv6Address Increment Repeat
 //AgtSipSubscriberPool SetRemoteHostAddressV6IncrementingRange
 return nil
}

func(np *SipSubscriberPool) GetRemoteHostAddressV6IncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetRemoteHostAddressV6IncrementingRange
 return "", nil
}

func(np *SipSubscriberPool) SetMediaPayloadType () error {
 //parameters: EmulationHandle MediaPayloadType
 //AgtSipSubscriberPool SetMediaPayloadType
 return nil
}

func(np *SipSubscriberPool) GetMediaPayloadType ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetMediaPayloadType
 return "", nil
}

func(np *SipSubscriberPool) SetMediaPort () error {
 //parameters: EmulationHandle Port
 //AgtSipSubscriberPool SetMediaPort
 return nil
}

func(np *SipSubscriberPool) GetMediaPort ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetMediaPort
 return "", nil
}

func(np *SipSubscriberPool) SetSipPort () error {
 //parameters: EmulationHandle Port
 //AgtSipSubscriberPool SetSipPort
 return nil
}

func(np *SipSubscriberPool) GetSipPort ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetSipPort
 return "", nil
}

func(np *SipSubscriberPool) IsRegistrarEnabled () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool IsRegistrarEnabled
 return nil
}

func(np *SipSubscriberPool) EnableRegistrar () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool EnableRegistrar
 return nil
}

func(np *SipSubscriberPool) DisableRegistrar () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool DisableRegistrar
 return nil
}

func(np *SipSubscriberPool) SetRegistrarIpv4Address () error {
 //parameters: EmulationHandle Ipv4Address
 //AgtSipSubscriberPool SetRegistrarIpv4Address
 return nil
}

func(np *SipSubscriberPool) GetRegistrarIpv4Address ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetRegistrarIpv4Address
 return "", nil
}

func(np *SipSubscriberPool) SetRegistrarIpv6Address () error {
 //parameters: EmulationHandle Ipv6Address
 //AgtSipSubscriberPool SetRegistrarIpv6Address
 return nil
}

func(np *SipSubscriberPool) GetRegistrarIpv6Address ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetRegistrarIpv6Address
 return "", nil
}

func(np *SipSubscriberPool) SetRegistrarL3AddressFamily () error {
 //parameters: EmulationHandle L3AddressFamily
 //AgtSipSubscriberPool SetRegistrarL3AddressFamily
 return nil
}

func(np *SipSubscriberPool) GetRegistrarL3AddressFamily ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetRegistrarL3AddressFamily
 return "", nil
}

func(np *SipSubscriberPool) SetRegistrarHostname () error {
 //parameters: EmulationHandle RegistrarHostname
 //AgtSipSubscriberPool SetRegistrarHostname
 return nil
}

func(np *SipSubscriberPool) GetRegistrarHostname ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetRegistrarHostname
 return "", nil
}

func(np *SipSubscriberPool) IsAutomaticRegistrationEnabled () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool IsAutomaticRegistrationEnabled
 return nil
}

func(np *SipSubscriberPool) EnableAutomaticRegistration () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool EnableAutomaticRegistration
 return nil
}

func(np *SipSubscriberPool) DisableAutomaticRegistration () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool DisableAutomaticRegistration
 return nil
}

func(np *SipSubscriberPool) SetDesiredRegistrationExpiryTime () error {
 //parameters: EmulationHandle DesiredRegistrationExpiryTime
 //AgtSipSubscriberPool SetDesiredRegistrationExpiryTime
 return nil
}

func(np *SipSubscriberPool) GetDesiredRegistrationExpiryTime ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetDesiredRegistrationExpiryTime
 return "", nil
}

func(np *SipSubscriberPool) SetMaxRegistrationAttempts () error {
 //parameters: EmulationHandle MaxRegistrationAttempts
 //AgtSipSubscriberPool SetMaxRegistrationAttempts
 return nil
}

func(np *SipSubscriberPool) GetMaxRegistrationAttempts ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetMaxRegistrationAttempts
 return "", nil
}

func(np *SipSubscriberPool) SetRegistrationAttemptDelayTime () error {
 //parameters: EmulationHandle RegistrationAttemptDelayTime
 //AgtSipSubscriberPool SetRegistrationAttemptDelayTime
 return nil
}

func(np *SipSubscriberPool) GetRegistrationAttemptDelayTime ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetRegistrationAttemptDelayTime
 return "", nil
}

func(np *SipSubscriberPool) SetTransactionTimerT1 () error {
 //parameters: EmulationHandle Time
 //AgtSipSubscriberPool SetTransactionTimerT1
 return nil
}

func(np *SipSubscriberPool) SetTransactionTimerT2 () error {
 //parameters: EmulationHandle Time
 //AgtSipSubscriberPool SetTransactionTimerT2
 return nil
}

func(np *SipSubscriberPool) SetTransactionTimerT4 () error {
 //parameters: EmulationHandle Time
 //AgtSipSubscriberPool SetTransactionTimerT4
 return nil
}

func(np *SipSubscriberPool) GetTransactionTimerT1 ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetTransactionTimerT1
 return "", nil
}

func(np *SipSubscriberPool) GetTransactionTimerT2 ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetTransactionTimerT2
 return "", nil
}

func(np *SipSubscriberPool) GetTransactionTimerT4 ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetTransactionTimerT4
 return "", nil
}

func(np *SipSubscriberPool) IsResourcePriorityEnabled () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool IsResourcePriorityEnabled
 return nil
}

func(np *SipSubscriberPool) EnableResourcePriority () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool EnableResourcePriority
 return nil
}

func(np *SipSubscriberPool) DisableResourcePriority () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool DisableResourcePriority
 return nil
}

func(np *SipSubscriberPool) SetResourcePriority () error {
 //parameters: EmulationHandle ResourcePriority
 //AgtSipSubscriberPool SetResourcePriority
 return nil
}

func(np *SipSubscriberPool) GetResourcePriority ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetResourcePriority
 return "", nil
}

func(np *SipSubscriberPool) IsSessionRefreshTimerEnabled () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool IsSessionRefreshTimerEnabled
 return nil
}

func(np *SipSubscriberPool) EnableSessionRefreshTimer () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool EnableSessionRefreshTimer
 return nil
}

func(np *SipSubscriberPool) DisableSessionRefreshTimer () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool DisableSessionRefreshTimer
 return nil
}

func(np *SipSubscriberPool) SetDesiredSessionInterval () error {
 //parameters: EmulationHandle DesiredSessionInterval
 //AgtSipSubscriberPool SetDesiredSessionInterval
 return nil
}

func(np *SipSubscriberPool) GetDesiredSessionInterval ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetDesiredSessionInterval
 return "", nil
}

func(np *SipSubscriberPool) GetMinimumSessionInterval ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetMinimumSessionInterval
 return "", nil
}

func(np *SipSubscriberPool) SetMinimumSessionInterval () error {
 //parameters: EmulationHandle MinimumSessionInterval
 //AgtSipSubscriberPool SetMinimumSessionInterval
 return nil
}

func(np *SipSubscriberPool) SetRoundTripDelay () error {
 //parameters: EmulationHandle RoundTripInterval
 //AgtSipSubscriberPool SetRoundTripDelay
 return nil
}

func(np *SipSubscriberPool) GetRoundTripDelay ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetRoundTripDelay
 return "", nil
}

func(np *SipSubscriberPool) SetSessionRefresher () error {
 //parameters: EmulationHandle UserAgentType
 //AgtSipSubscriberPool SetSessionRefresher
 return nil
}

func(np *SipSubscriberPool) GetSessionRefresher ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetSessionRefresher
 return "", nil
}

func(np *SipSubscriberPool) IsCallAcceptDelayEnabled () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool IsCallAcceptDelayEnabled
 return nil
}

func(np *SipSubscriberPool) EnableCallAcceptDelay () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool EnableCallAcceptDelay
 return nil
}

func(np *SipSubscriberPool) DisableCallAcceptDelay () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool DisableCallAcceptDelay
 return nil
}

func(np *SipSubscriberPool) SetResponseDelayInterval () error {
 //parameters: EmulationHandle Time
 //AgtSipSubscriberPool SetResponseDelayInterval
 return nil
}

func(np *SipSubscriberPool) GetResponseDelayInterval ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetResponseDelayInterval
 return "", nil
}

func(np *SipSubscriberPool) IsPrivacyHeaderEnabled () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool IsPrivacyHeaderEnabled
 return nil
}

func(np *SipSubscriberPool) EnablePrivacyHeader () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool EnablePrivacyHeader
 return nil
}

func(np *SipSubscriberPool) DisablePrivacyHeader () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool DisablePrivacyHeader
 return nil
}

func(np *SipSubscriberPool) SetPrivacyType () error {
 //parameters: EmulationHandle PrivacyType
 //AgtSipSubscriberPool SetPrivacyType
 return nil
}

func(np *SipSubscriberPool) GetPrivacyType ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetPrivacyType
 return "", nil
}

func(np *SipSubscriberPool) IsAuthenticationEnabled () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool IsAuthenticationEnabled
 return nil
}

func(np *SipSubscriberPool) EnableAuthentication () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool EnableAuthentication
 return nil
}

func(np *SipSubscriberPool) DisableAuthentication () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool DisableAuthentication
 return nil
}

func(np *SipSubscriberPool) IsAuthenticationCredentialEnabled () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool IsAuthenticationCredentialEnabled
 return nil
}

func(np *SipSubscriberPool) EnableAuthenticationCredential () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool EnableAuthenticationCredential
 return nil
}

func(np *SipSubscriberPool) DisableAuthenticationCredential () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool DisableAuthenticationCredential
 return nil
}

func(np *SipSubscriberPool) IsCredentialInDialogEnabled () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool IsCredentialInDialogEnabled
 return nil
}

func(np *SipSubscriberPool) EnableCredentialInDialog () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool EnableCredentialInDialog
 return nil
}

func(np *SipSubscriberPool) DisableCredentialInDialog () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool DisableCredentialInDialog
 return nil
}

func(np *SipSubscriberPool) SetAuthCredentialUsernamePrefix () error {
 //parameters: EmulationHandle AuthCredentialUsernamePrefix
 //AgtSipSubscriberPool SetAuthCredentialUsernamePrefix
 return nil
}

func(np *SipSubscriberPool) GetAuthCredentialUsernamePrefix ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetAuthCredentialUsernamePrefix
 return "", nil
}

func(np *SipSubscriberPool) SetAuthCredentialUsernameSuffixIncrementingRange () error {
 //parameters: EmulationHandle AuthCredentialUsernameSuffix Increment Repeat
 //AgtSipSubscriberPool SetAuthCredentialUsernameSuffixIncrementingRange
 return nil
}

func(np *SipSubscriberPool) GetAuthCredentialUsernameSuffixIncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetAuthCredentialUsernameSuffixIncrementingRange
 return "", nil
}

func(np *SipSubscriberPool) SetAuthCredentialPassword () error {
 //parameters: EmulationHandle AuthCredentialPassword
 //AgtSipSubscriberPool SetAuthCredentialPassword
 return nil
}

func(np *SipSubscriberPool) GetAuthCredentialPassword ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetAuthCredentialPassword
 return "", nil
}

func(np *SipSubscriberPool) SetAuthCredentialPasswordModifier () error {
 //parameters: EmulationHandle AuthCredentialPasswordModifier Increment Repeat
 //AgtSipSubscriberPool SetAuthCredentialPasswordModifier
 return nil
}

func(np *SipSubscriberPool) GetAuthCredentialPasswordModifier ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetAuthCredentialPasswordModifier
 return "", nil
}

func(np *SipSubscriberPool) IsAuthenticationChallengeEnabled () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool IsAuthenticationChallengeEnabled
 return nil
}

func(np *SipSubscriberPool) EnableAuthenticationChallenge () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool EnableAuthenticationChallenge
 return nil
}

func(np *SipSubscriberPool) DisableAuthenticationChallenge () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool DisableAuthenticationChallenge
 return nil
}

func(np *SipSubscriberPool) GetChallengeAuthenticationScheme ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetChallengeAuthenticationScheme
 return "", nil
}

func(np *SipSubscriberPool) SetChallengeAuthenticationRealm () error {
 //parameters: EmulationHandle ChallengeAuthenticationRealm
 //AgtSipSubscriberPool SetChallengeAuthenticationRealm
 return nil
}

func(np *SipSubscriberPool) GetChallengeAuthenticationRealm ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetChallengeAuthenticationRealm
 return "", nil
}

func(np *SipSubscriberPool) SetAuthChallengeUsernamePrefix () error {
 //parameters: EmulationHandle AuthChallengeUsernamePrefix
 //AgtSipSubscriberPool SetAuthChallengeUsernamePrefix
 return nil
}

func(np *SipSubscriberPool) GetAuthChallengeUsernamePrefix ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetAuthChallengeUsernamePrefix
 return "", nil
}

func(np *SipSubscriberPool) SetAuthChallengeUsernameSuffixIncrementingRange () error {
 //parameters: EmulationHandle AuthChallengeUsernameSuffix Increment Repeat
 //AgtSipSubscriberPool SetAuthChallengeUsernameSuffixIncrementingRange
 return nil
}

func(np *SipSubscriberPool) GetAuthChallengeUsernameSuffixIncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetAuthChallengeUsernameSuffixIncrementingRange
 return "", nil
}

func(np *SipSubscriberPool) SetAuthChallengePassword () error {
 //parameters: EmulationHandle AuthChallengePassword
 //AgtSipSubscriberPool SetAuthChallengePassword
 return nil
}

func(np *SipSubscriberPool) GetAuthChallengePassword ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetAuthChallengePassword
 return "", nil
}

func(np *SipSubscriberPool) SetAuthChallengePasswordModifier () error {
 //parameters: EmulationHandle AuthChallengePasswordModifier Increment Repeat
 //AgtSipSubscriberPool SetAuthChallengePasswordModifier
 return nil
}

func(np *SipSubscriberPool) GetAuthChallengePasswordModifier ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetAuthChallengePasswordModifier
 return "", nil
}

func(np *SipSubscriberPool) SetChallengeOpaque () error {
 //parameters: EmulationHandle ChallengeOpaque
 //AgtSipSubscriberPool SetChallengeOpaque
 return nil
}

func(np *SipSubscriberPool) GetChallengeOpaque ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetChallengeOpaque
 return "", nil
}

func(np *SipSubscriberPool) SetChallengeAlgorithm () error {
 //parameters: EmulationHandle ChallengeAlgorithm
 //AgtSipSubscriberPool SetChallengeAlgorithm
 return nil
}

func(np *SipSubscriberPool) GetChallengeAlgorithm ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetChallengeAlgorithm
 return "", nil
}

func(np *SipSubscriberPool) SetChallengeQop () error {
 //parameters: EmulationHandle ChallengeQop
 //AgtSipSubscriberPool SetChallengeQop
 return nil
}

func(np *SipSubscriberPool) GetChallengeQop ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetChallengeQop
 return "", nil
}

func(np *SipSubscriberPool) SetNonceExpireTime () error {
 //parameters: EmulationHandle NonceExpireTime
 //AgtSipSubscriberPool SetNonceExpireTime
 return nil
}

func(np *SipSubscriberPool) GetNonceExpireTime ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetNonceExpireTime
 return "", nil
}

func(np *SipSubscriberPool) SetMaximumChallengeAttempts () error {
 //parameters: EmulationHandle MaximumChallengeAttempts
 //AgtSipSubscriberPool SetMaximumChallengeAttempts
 return nil
}

func(np *SipSubscriberPool) GetMaximumChallengeAttempts ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool GetMaximumChallengeAttempts
 return "", nil
}

func(np *SipSubscriberPool) SendTrafficParams () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool SendTrafficParams
 return nil
}

func(np *SipSubscriberPool) ResetTrafficParams () error {
 //parameters: EmulationHandle
 //AgtSipSubscriberPool ResetTrafficParams
 return nil
}

