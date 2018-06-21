package n2xsdk

type SipIpv6UserAgentPool struct {
 Handler string
}

func(np *SipIpv6UserAgentPool) EstablishSession () error {
 //parameters: DeviceOrSession
 //AgtSipIpv6UserAgentPool EstablishSession
 return nil
}

func(np *SipIpv6UserAgentPool) TerminateSession () error {
 //parameters: DeviceOrSession
 //AgtSipIpv6UserAgentPool TerminateSession
 return nil
}

func(np *SipIpv6UserAgentPool) RegisterSession () error {
 //parameters: DeviceOrSession
 //AgtSipIpv6UserAgentPool RegisterSession
 return nil
}

func(np *SipIpv6UserAgentPool) DeregisterSession () error {
 //parameters: DeviceOrSession
 //AgtSipIpv6UserAgentPool DeregisterSession
 return nil
}

func(np *SipIpv6UserAgentPool) SetLocalUsernamePrefix () error {
 //parameters: EmulationHandle LocalUserNamePrefix
 //AgtSipIpv6UserAgentPool SetLocalUsernamePrefix
 return nil
}

func(np *SipIpv6UserAgentPool) SetLocalUsernameSuffixIncrementingRange () error {
 //parameters: EmulationHandle LocalUserNameSuffix Increment Repeat
 //AgtSipIpv6UserAgentPool SetLocalUsernameSuffixIncrementingRange
 return nil
}

func(np *SipIpv6UserAgentPool) SetLocalDomain () error {
 //parameters: EmulationHandle LocalDomain
 //AgtSipIpv6UserAgentPool SetLocalDomain
 return nil
}

func(np *SipIpv6UserAgentPool) GetLocalDomain ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetLocalDomain
 return "", nil
}

func(np *SipIpv6UserAgentPool) GetLocalUsernamePrefix ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetLocalUsernamePrefix
 return "", nil
}

func(np *SipIpv6UserAgentPool) GetLocalUsernameSuffixIncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetLocalUsernameSuffixIncrementingRange
 return "", nil
}

func(np *SipIpv6UserAgentPool) IsUseRemoteAddressOfRecordEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool IsUseRemoteAddressOfRecordEnabled
 return nil
}

func(np *SipIpv6UserAgentPool) EnableUseRemoteAddressOfRecord () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool EnableUseRemoteAddressOfRecord
 return nil
}

func(np *SipIpv6UserAgentPool) DisableUseRemoteAddressOfRecord () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool DisableUseRemoteAddressOfRecord
 return nil
}

func(np *SipIpv6UserAgentPool) SetRemoteUsernamePrefix () error {
 //parameters: EmulationHandle RemoteUserNamePrefix
 //AgtSipIpv6UserAgentPool SetRemoteUsernamePrefix
 return nil
}

func(np *SipIpv6UserAgentPool) SetRemoteUsernameSuffixIncrementingRange () error {
 //parameters: EmulationHandle RemoteUserNameSuffix Increment Repeat
 //AgtSipIpv6UserAgentPool SetRemoteUsernameSuffixIncrementingRange
 return nil
}

func(np *SipIpv6UserAgentPool) GetRemoteUsernameSuffixIncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetRemoteUsernameSuffixIncrementingRange
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetRemoteDomain () error {
 //parameters: EmulationHandle RemoteDomain
 //AgtSipIpv6UserAgentPool SetRemoteDomain
 return nil
}

func(np *SipIpv6UserAgentPool) GetRemoteUsernamePrefix ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetRemoteUsernamePrefix
 return "", nil
}

func(np *SipIpv6UserAgentPool) GetRemoteDomain ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetRemoteDomain
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetRemoteHostAddressV4IncrementingRange () error {
 //parameters: EmulationHandle Ipv4Address Increment Repeat
 //AgtSipIpv6UserAgentPool SetRemoteHostAddressV4IncrementingRange
 return nil
}

func(np *SipIpv6UserAgentPool) SetRemoteHostL3AddressFamily () error {
 //parameters: EmulationHandle RegistrarAddressType
 //AgtSipIpv6UserAgentPool SetRemoteHostL3AddressFamily
 return nil
}

func(np *SipIpv6UserAgentPool) GetRemoteHostL3AddressFamily ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetRemoteHostL3AddressFamily
 return "", nil
}

func(np *SipIpv6UserAgentPool) GetRemoteHostAddressV4IncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetRemoteHostAddressV4IncrementingRange
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetRemoteHostAddressV6IncrementingRange () error {
 //parameters: EmulationHandle Ipv6Address Increment Repeat
 //AgtSipIpv6UserAgentPool SetRemoteHostAddressV6IncrementingRange
 return nil
}

func(np *SipIpv6UserAgentPool) GetRemoteHostAddressV6IncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetRemoteHostAddressV6IncrementingRange
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetMediaPayloadType () error {
 //parameters: EmulationHandle MediaPayloadType
 //AgtSipIpv6UserAgentPool SetMediaPayloadType
 return nil
}

func(np *SipIpv6UserAgentPool) GetMediaPayloadType ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetMediaPayloadType
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetMediaPort () error {
 //parameters: EmulationHandle Port
 //AgtSipIpv6UserAgentPool SetMediaPort
 return nil
}

func(np *SipIpv6UserAgentPool) GetMediaPort ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetMediaPort
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetSipPort () error {
 //parameters: EmulationHandle Port
 //AgtSipIpv6UserAgentPool SetSipPort
 return nil
}

func(np *SipIpv6UserAgentPool) GetSipPort ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetSipPort
 return "", nil
}

func(np *SipIpv6UserAgentPool) IsRegistrarEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool IsRegistrarEnabled
 return nil
}

func(np *SipIpv6UserAgentPool) EnableRegistrar () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool EnableRegistrar
 return nil
}

func(np *SipIpv6UserAgentPool) DisableRegistrar () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool DisableRegistrar
 return nil
}

func(np *SipIpv6UserAgentPool) SetRegistrarIpv4Address () error {
 //parameters: EmulationHandle Ipv4Address
 //AgtSipIpv6UserAgentPool SetRegistrarIpv4Address
 return nil
}

func(np *SipIpv6UserAgentPool) GetRegistrarIpv4Address ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetRegistrarIpv4Address
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetRegistrarIpv6Address () error {
 //parameters: EmulationHandle Ipv6Address
 //AgtSipIpv6UserAgentPool SetRegistrarIpv6Address
 return nil
}

func(np *SipIpv6UserAgentPool) GetRegistrarIpv6Address ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetRegistrarIpv6Address
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetRegistrarL3AddressFamily () error {
 //parameters: EmulationHandle L3AddressFamily
 //AgtSipIpv6UserAgentPool SetRegistrarL3AddressFamily
 return nil
}

func(np *SipIpv6UserAgentPool) GetRegistrarL3AddressFamily ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetRegistrarL3AddressFamily
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetRegistrarHostname () error {
 //parameters: EmulationHandle RegistrarHostname
 //AgtSipIpv6UserAgentPool SetRegistrarHostname
 return nil
}

func(np *SipIpv6UserAgentPool) GetRegistrarHostname ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetRegistrarHostname
 return "", nil
}

func(np *SipIpv6UserAgentPool) IsAutomaticRegistrationEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool IsAutomaticRegistrationEnabled
 return nil
}

func(np *SipIpv6UserAgentPool) EnableAutomaticRegistration () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool EnableAutomaticRegistration
 return nil
}

func(np *SipIpv6UserAgentPool) DisableAutomaticRegistration () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool DisableAutomaticRegistration
 return nil
}

func(np *SipIpv6UserAgentPool) SetDesiredRegistrationExpiryTime () error {
 //parameters: EmulationHandle DesiredRegistrationExpiryTime
 //AgtSipIpv6UserAgentPool SetDesiredRegistrationExpiryTime
 return nil
}

func(np *SipIpv6UserAgentPool) GetDesiredRegistrationExpiryTime ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetDesiredRegistrationExpiryTime
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetMaxRegistrationAttempts () error {
 //parameters: EmulationHandle MaxRegistrationAttempts
 //AgtSipIpv6UserAgentPool SetMaxRegistrationAttempts
 return nil
}

func(np *SipIpv6UserAgentPool) GetMaxRegistrationAttempts ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetMaxRegistrationAttempts
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetRegistrationAttemptDelayTime () error {
 //parameters: EmulationHandle RegistrationAttemptDelayTime
 //AgtSipIpv6UserAgentPool SetRegistrationAttemptDelayTime
 return nil
}

func(np *SipIpv6UserAgentPool) GetRegistrationAttemptDelayTime ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetRegistrationAttemptDelayTime
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetTransactionTimerT1 () error {
 //parameters: EmulationHandle Time
 //AgtSipIpv6UserAgentPool SetTransactionTimerT1
 return nil
}

func(np *SipIpv6UserAgentPool) SetTransactionTimerT2 () error {
 //parameters: EmulationHandle Time
 //AgtSipIpv6UserAgentPool SetTransactionTimerT2
 return nil
}

func(np *SipIpv6UserAgentPool) SetTransactionTimerT4 () error {
 //parameters: EmulationHandle Time
 //AgtSipIpv6UserAgentPool SetTransactionTimerT4
 return nil
}

func(np *SipIpv6UserAgentPool) GetTransactionTimerT1 ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetTransactionTimerT1
 return "", nil
}

func(np *SipIpv6UserAgentPool) GetTransactionTimerT2 ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetTransactionTimerT2
 return "", nil
}

func(np *SipIpv6UserAgentPool) GetTransactionTimerT4 ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetTransactionTimerT4
 return "", nil
}

func(np *SipIpv6UserAgentPool) IsResourcePriorityEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool IsResourcePriorityEnabled
 return nil
}

func(np *SipIpv6UserAgentPool) EnableResourcePriority () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool EnableResourcePriority
 return nil
}

func(np *SipIpv6UserAgentPool) DisableResourcePriority () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool DisableResourcePriority
 return nil
}

func(np *SipIpv6UserAgentPool) SetResourcePriority () error {
 //parameters: EmulationHandle ResourcePriority
 //AgtSipIpv6UserAgentPool SetResourcePriority
 return nil
}

func(np *SipIpv6UserAgentPool) GetResourcePriority ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetResourcePriority
 return "", nil
}

func(np *SipIpv6UserAgentPool) IsSessionRefreshTimerEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool IsSessionRefreshTimerEnabled
 return nil
}

func(np *SipIpv6UserAgentPool) EnableSessionRefreshTimer () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool EnableSessionRefreshTimer
 return nil
}

func(np *SipIpv6UserAgentPool) DisableSessionRefreshTimer () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool DisableSessionRefreshTimer
 return nil
}

func(np *SipIpv6UserAgentPool) SetDesiredSessionInterval () error {
 //parameters: EmulationHandle DesiredSessionInterval
 //AgtSipIpv6UserAgentPool SetDesiredSessionInterval
 return nil
}

func(np *SipIpv6UserAgentPool) GetDesiredSessionInterval ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetDesiredSessionInterval
 return "", nil
}

func(np *SipIpv6UserAgentPool) GetMinimumSessionInterval ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetMinimumSessionInterval
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetMinimumSessionInterval () error {
 //parameters: EmulationHandle MinimumSessionInterval
 //AgtSipIpv6UserAgentPool SetMinimumSessionInterval
 return nil
}

func(np *SipIpv6UserAgentPool) SetRoundTripDelay () error {
 //parameters: EmulationHandle RoundTripInterval
 //AgtSipIpv6UserAgentPool SetRoundTripDelay
 return nil
}

func(np *SipIpv6UserAgentPool) GetRoundTripDelay ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetRoundTripDelay
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetSessionRefresher () error {
 //parameters: EmulationHandle UserAgentType
 //AgtSipIpv6UserAgentPool SetSessionRefresher
 return nil
}

func(np *SipIpv6UserAgentPool) GetSessionRefresher ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetSessionRefresher
 return "", nil
}

func(np *SipIpv6UserAgentPool) IsCallAcceptDelayEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool IsCallAcceptDelayEnabled
 return nil
}

func(np *SipIpv6UserAgentPool) EnableCallAcceptDelay () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool EnableCallAcceptDelay
 return nil
}

func(np *SipIpv6UserAgentPool) DisableCallAcceptDelay () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool DisableCallAcceptDelay
 return nil
}

func(np *SipIpv6UserAgentPool) SetResponseDelayInterval () error {
 //parameters: EmulationHandle Time
 //AgtSipIpv6UserAgentPool SetResponseDelayInterval
 return nil
}

func(np *SipIpv6UserAgentPool) GetResponseDelayInterval ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetResponseDelayInterval
 return "", nil
}

func(np *SipIpv6UserAgentPool) IsPrivacyHeaderEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool IsPrivacyHeaderEnabled
 return nil
}

func(np *SipIpv6UserAgentPool) EnablePrivacyHeader () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool EnablePrivacyHeader
 return nil
}

func(np *SipIpv6UserAgentPool) DisablePrivacyHeader () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool DisablePrivacyHeader
 return nil
}

func(np *SipIpv6UserAgentPool) SetPrivacyType () error {
 //parameters: EmulationHandle PrivacyType
 //AgtSipIpv6UserAgentPool SetPrivacyType
 return nil
}

func(np *SipIpv6UserAgentPool) GetPrivacyType ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetPrivacyType
 return "", nil
}

func(np *SipIpv6UserAgentPool) IsAuthenticationEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool IsAuthenticationEnabled
 return nil
}

func(np *SipIpv6UserAgentPool) EnableAuthentication () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool EnableAuthentication
 return nil
}

func(np *SipIpv6UserAgentPool) DisableAuthentication () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool DisableAuthentication
 return nil
}

func(np *SipIpv6UserAgentPool) IsAuthenticationCredentialEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool IsAuthenticationCredentialEnabled
 return nil
}

func(np *SipIpv6UserAgentPool) EnableAuthenticationCredential () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool EnableAuthenticationCredential
 return nil
}

func(np *SipIpv6UserAgentPool) DisableAuthenticationCredential () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool DisableAuthenticationCredential
 return nil
}

func(np *SipIpv6UserAgentPool) IsCredentialInDialogEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool IsCredentialInDialogEnabled
 return nil
}

func(np *SipIpv6UserAgentPool) EnableCredentialInDialog () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool EnableCredentialInDialog
 return nil
}

func(np *SipIpv6UserAgentPool) DisableCredentialInDialog () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool DisableCredentialInDialog
 return nil
}

func(np *SipIpv6UserAgentPool) SetAuthCredentialUsernamePrefix () error {
 //parameters: EmulationHandle AuthCredentialUsernamePrefix
 //AgtSipIpv6UserAgentPool SetAuthCredentialUsernamePrefix
 return nil
}

func(np *SipIpv6UserAgentPool) GetAuthCredentialUsernamePrefix ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetAuthCredentialUsernamePrefix
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetAuthCredentialUsernameSuffixIncrementingRange () error {
 //parameters: EmulationHandle AuthCredentialUsernameSuffix Increment Repeat
 //AgtSipIpv6UserAgentPool SetAuthCredentialUsernameSuffixIncrementingRange
 return nil
}

func(np *SipIpv6UserAgentPool) GetAuthCredentialUsernameSuffixIncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetAuthCredentialUsernameSuffixIncrementingRange
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetAuthCredentialPassword () error {
 //parameters: EmulationHandle AuthCredentialPassword
 //AgtSipIpv6UserAgentPool SetAuthCredentialPassword
 return nil
}

func(np *SipIpv6UserAgentPool) GetAuthCredentialPassword ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetAuthCredentialPassword
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetAuthCredentialPasswordModifier () error {
 //parameters: EmulationHandle AuthCredentialPasswordModifier Increment Repeat
 //AgtSipIpv6UserAgentPool SetAuthCredentialPasswordModifier
 return nil
}

func(np *SipIpv6UserAgentPool) GetAuthCredentialPasswordModifier ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetAuthCredentialPasswordModifier
 return "", nil
}

func(np *SipIpv6UserAgentPool) IsAuthenticationChallengeEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool IsAuthenticationChallengeEnabled
 return nil
}

func(np *SipIpv6UserAgentPool) EnableAuthenticationChallenge () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool EnableAuthenticationChallenge
 return nil
}

func(np *SipIpv6UserAgentPool) DisableAuthenticationChallenge () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool DisableAuthenticationChallenge
 return nil
}

func(np *SipIpv6UserAgentPool) GetChallengeAuthenticationScheme ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetChallengeAuthenticationScheme
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetChallengeAuthenticationRealm () error {
 //parameters: EmulationHandle ChallengeAuthenticationRealm
 //AgtSipIpv6UserAgentPool SetChallengeAuthenticationRealm
 return nil
}

func(np *SipIpv6UserAgentPool) GetChallengeAuthenticationRealm ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetChallengeAuthenticationRealm
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetAuthChallengeUsernamePrefix () error {
 //parameters: EmulationHandle AuthChallengeUsernamePrefix
 //AgtSipIpv6UserAgentPool SetAuthChallengeUsernamePrefix
 return nil
}

func(np *SipIpv6UserAgentPool) GetAuthChallengeUsernamePrefix ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetAuthChallengeUsernamePrefix
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetAuthChallengeUsernameSuffixIncrementingRange () error {
 //parameters: EmulationHandle AuthChallengeUsernameSuffix Increment Repeat
 //AgtSipIpv6UserAgentPool SetAuthChallengeUsernameSuffixIncrementingRange
 return nil
}

func(np *SipIpv6UserAgentPool) GetAuthChallengeUsernameSuffixIncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetAuthChallengeUsernameSuffixIncrementingRange
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetAuthChallengePassword () error {
 //parameters: EmulationHandle AuthChallengePassword
 //AgtSipIpv6UserAgentPool SetAuthChallengePassword
 return nil
}

func(np *SipIpv6UserAgentPool) GetAuthChallengePassword ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetAuthChallengePassword
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetAuthChallengePasswordModifier () error {
 //parameters: EmulationHandle AuthChallengePasswordModifier Increment Repeat
 //AgtSipIpv6UserAgentPool SetAuthChallengePasswordModifier
 return nil
}

func(np *SipIpv6UserAgentPool) GetAuthChallengePasswordModifier ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetAuthChallengePasswordModifier
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetChallengeOpaque () error {
 //parameters: EmulationHandle ChallengeOpaque
 //AgtSipIpv6UserAgentPool SetChallengeOpaque
 return nil
}

func(np *SipIpv6UserAgentPool) GetChallengeOpaque ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetChallengeOpaque
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetChallengeAlgorithm () error {
 //parameters: EmulationHandle ChallengeAlgorithm
 //AgtSipIpv6UserAgentPool SetChallengeAlgorithm
 return nil
}

func(np *SipIpv6UserAgentPool) GetChallengeAlgorithm ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetChallengeAlgorithm
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetChallengeQop () error {
 //parameters: EmulationHandle ChallengeQop
 //AgtSipIpv6UserAgentPool SetChallengeQop
 return nil
}

func(np *SipIpv6UserAgentPool) GetChallengeQop ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetChallengeQop
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetNonceExpireTime () error {
 //parameters: EmulationHandle NonceExpireTime
 //AgtSipIpv6UserAgentPool SetNonceExpireTime
 return nil
}

func(np *SipIpv6UserAgentPool) GetNonceExpireTime ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetNonceExpireTime
 return "", nil
}

func(np *SipIpv6UserAgentPool) SetMaximumChallengeAttempts () error {
 //parameters: EmulationHandle MaximumChallengeAttempts
 //AgtSipIpv6UserAgentPool SetMaximumChallengeAttempts
 return nil
}

func(np *SipIpv6UserAgentPool) GetMaximumChallengeAttempts ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool GetMaximumChallengeAttempts
 return "", nil
}

func(np *SipIpv6UserAgentPool) SendTrafficParams () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool SendTrafficParams
 return nil
}

func(np *SipIpv6UserAgentPool) ResetTrafficParams () error {
 //parameters: EmulationHandle
 //AgtSipIpv6UserAgentPool ResetTrafficParams
 return nil
}

