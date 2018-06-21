package n2xsdk

type SipIpv4UserAgentPool struct {
 Handler string
}

func(np *SipIpv4UserAgentPool) IsNetworkAddressTranslationEnabled () error {
 //parameters: DeviceHandle
 //AgtSipIpv4UserAgentPool IsNetworkAddressTranslationEnabled
 return nil
}

func(np *SipIpv4UserAgentPool) EnableNetworkAddressTranslation () error {
 //parameters: DeviceHandle
 //AgtSipIpv4UserAgentPool EnableNetworkAddressTranslation
 return nil
}

func(np *SipIpv4UserAgentPool) DisableNetworkAddressTranslation () error {
 //parameters: DeviceHandle
 //AgtSipIpv4UserAgentPool DisableNetworkAddressTranslation
 return nil
}

func(np *SipIpv4UserAgentPool) SetNatedInternaIpv4AddressRange () error {
 //parameters: DeviceHandle Ipv4Address Increment Repeat
 //AgtSipIpv4UserAgentPool SetNatedInternaIpv4AddressRange
 return nil
}

func(np *SipIpv4UserAgentPool) GetNatedInternaIpv4AddressRange ()(string, error) {
 //parameters: DeviceHandle
 //AgtSipIpv4UserAgentPool GetNatedInternaIpv4AddressRange
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetNatedExternalIpv4AddressRange () error {
 //parameters: DeviceHandle Ipv4Address Increment Repeat
 //AgtSipIpv4UserAgentPool SetNatedExternalIpv4AddressRange
 return nil
}

func(np *SipIpv4UserAgentPool) GetNatedExternalIpv4AddressRange ()(string, error) {
 //parameters: DeviceHandle
 //AgtSipIpv4UserAgentPool GetNatedExternalIpv4AddressRange
 return "", nil
}

func(np *SipIpv4UserAgentPool) EstablishSession () error {
 //parameters: DeviceOrSession
 //AgtSipIpv4UserAgentPool EstablishSession
 return nil
}

func(np *SipIpv4UserAgentPool) TerminateSession () error {
 //parameters: DeviceOrSession
 //AgtSipIpv4UserAgentPool TerminateSession
 return nil
}

func(np *SipIpv4UserAgentPool) RegisterSession () error {
 //parameters: DeviceOrSession
 //AgtSipIpv4UserAgentPool RegisterSession
 return nil
}

func(np *SipIpv4UserAgentPool) DeregisterSession () error {
 //parameters: DeviceOrSession
 //AgtSipIpv4UserAgentPool DeregisterSession
 return nil
}

func(np *SipIpv4UserAgentPool) SetLocalUsernamePrefix () error {
 //parameters: EmulationHandle LocalUserNamePrefix
 //AgtSipIpv4UserAgentPool SetLocalUsernamePrefix
 return nil
}

func(np *SipIpv4UserAgentPool) SetLocalUsernameSuffixIncrementingRange () error {
 //parameters: EmulationHandle LocalUserNameSuffix Increment Repeat
 //AgtSipIpv4UserAgentPool SetLocalUsernameSuffixIncrementingRange
 return nil
}

func(np *SipIpv4UserAgentPool) SetLocalDomain () error {
 //parameters: EmulationHandle LocalDomain
 //AgtSipIpv4UserAgentPool SetLocalDomain
 return nil
}

func(np *SipIpv4UserAgentPool) GetLocalDomain ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetLocalDomain
 return "", nil
}

func(np *SipIpv4UserAgentPool) GetLocalUsernamePrefix ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetLocalUsernamePrefix
 return "", nil
}

func(np *SipIpv4UserAgentPool) GetLocalUsernameSuffixIncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetLocalUsernameSuffixIncrementingRange
 return "", nil
}

func(np *SipIpv4UserAgentPool) IsUseRemoteAddressOfRecordEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool IsUseRemoteAddressOfRecordEnabled
 return nil
}

func(np *SipIpv4UserAgentPool) EnableUseRemoteAddressOfRecord () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool EnableUseRemoteAddressOfRecord
 return nil
}

func(np *SipIpv4UserAgentPool) DisableUseRemoteAddressOfRecord () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool DisableUseRemoteAddressOfRecord
 return nil
}

func(np *SipIpv4UserAgentPool) SetRemoteUsernamePrefix () error {
 //parameters: EmulationHandle RemoteUserNamePrefix
 //AgtSipIpv4UserAgentPool SetRemoteUsernamePrefix
 return nil
}

func(np *SipIpv4UserAgentPool) SetRemoteUsernameSuffixIncrementingRange () error {
 //parameters: EmulationHandle RemoteUserNameSuffix Increment Repeat
 //AgtSipIpv4UserAgentPool SetRemoteUsernameSuffixIncrementingRange
 return nil
}

func(np *SipIpv4UserAgentPool) GetRemoteUsernameSuffixIncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetRemoteUsernameSuffixIncrementingRange
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetRemoteDomain () error {
 //parameters: EmulationHandle RemoteDomain
 //AgtSipIpv4UserAgentPool SetRemoteDomain
 return nil
}

func(np *SipIpv4UserAgentPool) GetRemoteUsernamePrefix ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetRemoteUsernamePrefix
 return "", nil
}

func(np *SipIpv4UserAgentPool) GetRemoteDomain ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetRemoteDomain
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetRemoteHostAddressV4IncrementingRange () error {
 //parameters: EmulationHandle Ipv4Address Increment Repeat
 //AgtSipIpv4UserAgentPool SetRemoteHostAddressV4IncrementingRange
 return nil
}

func(np *SipIpv4UserAgentPool) SetRemoteHostL3AddressFamily () error {
 //parameters: EmulationHandle RegistrarAddressType
 //AgtSipIpv4UserAgentPool SetRemoteHostL3AddressFamily
 return nil
}

func(np *SipIpv4UserAgentPool) GetRemoteHostL3AddressFamily ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetRemoteHostL3AddressFamily
 return "", nil
}

func(np *SipIpv4UserAgentPool) GetRemoteHostAddressV4IncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetRemoteHostAddressV4IncrementingRange
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetRemoteHostAddressV6IncrementingRange () error {
 //parameters: EmulationHandle Ipv6Address Increment Repeat
 //AgtSipIpv4UserAgentPool SetRemoteHostAddressV6IncrementingRange
 return nil
}

func(np *SipIpv4UserAgentPool) GetRemoteHostAddressV6IncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetRemoteHostAddressV6IncrementingRange
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetMediaPayloadType () error {
 //parameters: EmulationHandle MediaPayloadType
 //AgtSipIpv4UserAgentPool SetMediaPayloadType
 return nil
}

func(np *SipIpv4UserAgentPool) GetMediaPayloadType ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetMediaPayloadType
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetMediaPort () error {
 //parameters: EmulationHandle Port
 //AgtSipIpv4UserAgentPool SetMediaPort
 return nil
}

func(np *SipIpv4UserAgentPool) GetMediaPort ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetMediaPort
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetSipPort () error {
 //parameters: EmulationHandle Port
 //AgtSipIpv4UserAgentPool SetSipPort
 return nil
}

func(np *SipIpv4UserAgentPool) GetSipPort ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetSipPort
 return "", nil
}

func(np *SipIpv4UserAgentPool) IsRegistrarEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool IsRegistrarEnabled
 return nil
}

func(np *SipIpv4UserAgentPool) EnableRegistrar () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool EnableRegistrar
 return nil
}

func(np *SipIpv4UserAgentPool) DisableRegistrar () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool DisableRegistrar
 return nil
}

func(np *SipIpv4UserAgentPool) SetRegistrarIpv4Address () error {
 //parameters: EmulationHandle Ipv4Address
 //AgtSipIpv4UserAgentPool SetRegistrarIpv4Address
 return nil
}

func(np *SipIpv4UserAgentPool) GetRegistrarIpv4Address ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetRegistrarIpv4Address
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetRegistrarIpv6Address () error {
 //parameters: EmulationHandle Ipv6Address
 //AgtSipIpv4UserAgentPool SetRegistrarIpv6Address
 return nil
}

func(np *SipIpv4UserAgentPool) GetRegistrarIpv6Address ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetRegistrarIpv6Address
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetRegistrarL3AddressFamily () error {
 //parameters: EmulationHandle L3AddressFamily
 //AgtSipIpv4UserAgentPool SetRegistrarL3AddressFamily
 return nil
}

func(np *SipIpv4UserAgentPool) GetRegistrarL3AddressFamily ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetRegistrarL3AddressFamily
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetRegistrarHostname () error {
 //parameters: EmulationHandle RegistrarHostname
 //AgtSipIpv4UserAgentPool SetRegistrarHostname
 return nil
}

func(np *SipIpv4UserAgentPool) GetRegistrarHostname ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetRegistrarHostname
 return "", nil
}

func(np *SipIpv4UserAgentPool) IsAutomaticRegistrationEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool IsAutomaticRegistrationEnabled
 return nil
}

func(np *SipIpv4UserAgentPool) EnableAutomaticRegistration () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool EnableAutomaticRegistration
 return nil
}

func(np *SipIpv4UserAgentPool) DisableAutomaticRegistration () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool DisableAutomaticRegistration
 return nil
}

func(np *SipIpv4UserAgentPool) SetDesiredRegistrationExpiryTime () error {
 //parameters: EmulationHandle DesiredRegistrationExpiryTime
 //AgtSipIpv4UserAgentPool SetDesiredRegistrationExpiryTime
 return nil
}

func(np *SipIpv4UserAgentPool) GetDesiredRegistrationExpiryTime ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetDesiredRegistrationExpiryTime
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetMaxRegistrationAttempts () error {
 //parameters: EmulationHandle MaxRegistrationAttempts
 //AgtSipIpv4UserAgentPool SetMaxRegistrationAttempts
 return nil
}

func(np *SipIpv4UserAgentPool) GetMaxRegistrationAttempts ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetMaxRegistrationAttempts
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetRegistrationAttemptDelayTime () error {
 //parameters: EmulationHandle RegistrationAttemptDelayTime
 //AgtSipIpv4UserAgentPool SetRegistrationAttemptDelayTime
 return nil
}

func(np *SipIpv4UserAgentPool) GetRegistrationAttemptDelayTime ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetRegistrationAttemptDelayTime
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetTransactionTimerT1 () error {
 //parameters: EmulationHandle Time
 //AgtSipIpv4UserAgentPool SetTransactionTimerT1
 return nil
}

func(np *SipIpv4UserAgentPool) SetTransactionTimerT2 () error {
 //parameters: EmulationHandle Time
 //AgtSipIpv4UserAgentPool SetTransactionTimerT2
 return nil
}

func(np *SipIpv4UserAgentPool) SetTransactionTimerT4 () error {
 //parameters: EmulationHandle Time
 //AgtSipIpv4UserAgentPool SetTransactionTimerT4
 return nil
}

func(np *SipIpv4UserAgentPool) GetTransactionTimerT1 ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetTransactionTimerT1
 return "", nil
}

func(np *SipIpv4UserAgentPool) GetTransactionTimerT2 ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetTransactionTimerT2
 return "", nil
}

func(np *SipIpv4UserAgentPool) GetTransactionTimerT4 ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetTransactionTimerT4
 return "", nil
}

func(np *SipIpv4UserAgentPool) IsResourcePriorityEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool IsResourcePriorityEnabled
 return nil
}

func(np *SipIpv4UserAgentPool) EnableResourcePriority () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool EnableResourcePriority
 return nil
}

func(np *SipIpv4UserAgentPool) DisableResourcePriority () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool DisableResourcePriority
 return nil
}

func(np *SipIpv4UserAgentPool) SetResourcePriority () error {
 //parameters: EmulationHandle ResourcePriority
 //AgtSipIpv4UserAgentPool SetResourcePriority
 return nil
}

func(np *SipIpv4UserAgentPool) GetResourcePriority ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetResourcePriority
 return "", nil
}

func(np *SipIpv4UserAgentPool) IsSessionRefreshTimerEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool IsSessionRefreshTimerEnabled
 return nil
}

func(np *SipIpv4UserAgentPool) EnableSessionRefreshTimer () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool EnableSessionRefreshTimer
 return nil
}

func(np *SipIpv4UserAgentPool) DisableSessionRefreshTimer () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool DisableSessionRefreshTimer
 return nil
}

func(np *SipIpv4UserAgentPool) SetDesiredSessionInterval () error {
 //parameters: EmulationHandle DesiredSessionInterval
 //AgtSipIpv4UserAgentPool SetDesiredSessionInterval
 return nil
}

func(np *SipIpv4UserAgentPool) GetDesiredSessionInterval ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetDesiredSessionInterval
 return "", nil
}

func(np *SipIpv4UserAgentPool) GetMinimumSessionInterval ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetMinimumSessionInterval
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetMinimumSessionInterval () error {
 //parameters: EmulationHandle MinimumSessionInterval
 //AgtSipIpv4UserAgentPool SetMinimumSessionInterval
 return nil
}

func(np *SipIpv4UserAgentPool) SetRoundTripDelay () error {
 //parameters: EmulationHandle RoundTripInterval
 //AgtSipIpv4UserAgentPool SetRoundTripDelay
 return nil
}

func(np *SipIpv4UserAgentPool) GetRoundTripDelay ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetRoundTripDelay
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetSessionRefresher () error {
 //parameters: EmulationHandle UserAgentType
 //AgtSipIpv4UserAgentPool SetSessionRefresher
 return nil
}

func(np *SipIpv4UserAgentPool) GetSessionRefresher ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetSessionRefresher
 return "", nil
}

func(np *SipIpv4UserAgentPool) IsCallAcceptDelayEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool IsCallAcceptDelayEnabled
 return nil
}

func(np *SipIpv4UserAgentPool) EnableCallAcceptDelay () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool EnableCallAcceptDelay
 return nil
}

func(np *SipIpv4UserAgentPool) DisableCallAcceptDelay () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool DisableCallAcceptDelay
 return nil
}

func(np *SipIpv4UserAgentPool) SetResponseDelayInterval () error {
 //parameters: EmulationHandle Time
 //AgtSipIpv4UserAgentPool SetResponseDelayInterval
 return nil
}

func(np *SipIpv4UserAgentPool) GetResponseDelayInterval ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetResponseDelayInterval
 return "", nil
}

func(np *SipIpv4UserAgentPool) IsPrivacyHeaderEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool IsPrivacyHeaderEnabled
 return nil
}

func(np *SipIpv4UserAgentPool) EnablePrivacyHeader () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool EnablePrivacyHeader
 return nil
}

func(np *SipIpv4UserAgentPool) DisablePrivacyHeader () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool DisablePrivacyHeader
 return nil
}

func(np *SipIpv4UserAgentPool) SetPrivacyType () error {
 //parameters: EmulationHandle PrivacyType
 //AgtSipIpv4UserAgentPool SetPrivacyType
 return nil
}

func(np *SipIpv4UserAgentPool) GetPrivacyType ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetPrivacyType
 return "", nil
}

func(np *SipIpv4UserAgentPool) IsAuthenticationEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool IsAuthenticationEnabled
 return nil
}

func(np *SipIpv4UserAgentPool) EnableAuthentication () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool EnableAuthentication
 return nil
}

func(np *SipIpv4UserAgentPool) DisableAuthentication () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool DisableAuthentication
 return nil
}

func(np *SipIpv4UserAgentPool) IsAuthenticationCredentialEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool IsAuthenticationCredentialEnabled
 return nil
}

func(np *SipIpv4UserAgentPool) EnableAuthenticationCredential () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool EnableAuthenticationCredential
 return nil
}

func(np *SipIpv4UserAgentPool) DisableAuthenticationCredential () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool DisableAuthenticationCredential
 return nil
}

func(np *SipIpv4UserAgentPool) IsCredentialInDialogEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool IsCredentialInDialogEnabled
 return nil
}

func(np *SipIpv4UserAgentPool) EnableCredentialInDialog () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool EnableCredentialInDialog
 return nil
}

func(np *SipIpv4UserAgentPool) DisableCredentialInDialog () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool DisableCredentialInDialog
 return nil
}

func(np *SipIpv4UserAgentPool) SetAuthCredentialUsernamePrefix () error {
 //parameters: EmulationHandle AuthCredentialUsernamePrefix
 //AgtSipIpv4UserAgentPool SetAuthCredentialUsernamePrefix
 return nil
}

func(np *SipIpv4UserAgentPool) GetAuthCredentialUsernamePrefix ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetAuthCredentialUsernamePrefix
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetAuthCredentialUsernameSuffixIncrementingRange () error {
 //parameters: EmulationHandle AuthCredentialUsernameSuffix Increment Repeat
 //AgtSipIpv4UserAgentPool SetAuthCredentialUsernameSuffixIncrementingRange
 return nil
}

func(np *SipIpv4UserAgentPool) GetAuthCredentialUsernameSuffixIncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetAuthCredentialUsernameSuffixIncrementingRange
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetAuthCredentialPassword () error {
 //parameters: EmulationHandle AuthCredentialPassword
 //AgtSipIpv4UserAgentPool SetAuthCredentialPassword
 return nil
}

func(np *SipIpv4UserAgentPool) GetAuthCredentialPassword ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetAuthCredentialPassword
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetAuthCredentialPasswordModifier () error {
 //parameters: EmulationHandle AuthCredentialPasswordModifier Increment Repeat
 //AgtSipIpv4UserAgentPool SetAuthCredentialPasswordModifier
 return nil
}

func(np *SipIpv4UserAgentPool) GetAuthCredentialPasswordModifier ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetAuthCredentialPasswordModifier
 return "", nil
}

func(np *SipIpv4UserAgentPool) IsAuthenticationChallengeEnabled () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool IsAuthenticationChallengeEnabled
 return nil
}

func(np *SipIpv4UserAgentPool) EnableAuthenticationChallenge () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool EnableAuthenticationChallenge
 return nil
}

func(np *SipIpv4UserAgentPool) DisableAuthenticationChallenge () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool DisableAuthenticationChallenge
 return nil
}

func(np *SipIpv4UserAgentPool) GetChallengeAuthenticationScheme ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetChallengeAuthenticationScheme
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetChallengeAuthenticationRealm () error {
 //parameters: EmulationHandle ChallengeAuthenticationRealm
 //AgtSipIpv4UserAgentPool SetChallengeAuthenticationRealm
 return nil
}

func(np *SipIpv4UserAgentPool) GetChallengeAuthenticationRealm ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetChallengeAuthenticationRealm
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetAuthChallengeUsernamePrefix () error {
 //parameters: EmulationHandle AuthChallengeUsernamePrefix
 //AgtSipIpv4UserAgentPool SetAuthChallengeUsernamePrefix
 return nil
}

func(np *SipIpv4UserAgentPool) GetAuthChallengeUsernamePrefix ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetAuthChallengeUsernamePrefix
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetAuthChallengeUsernameSuffixIncrementingRange () error {
 //parameters: EmulationHandle AuthChallengeUsernameSuffix Increment Repeat
 //AgtSipIpv4UserAgentPool SetAuthChallengeUsernameSuffixIncrementingRange
 return nil
}

func(np *SipIpv4UserAgentPool) GetAuthChallengeUsernameSuffixIncrementingRange ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetAuthChallengeUsernameSuffixIncrementingRange
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetAuthChallengePassword () error {
 //parameters: EmulationHandle AuthChallengePassword
 //AgtSipIpv4UserAgentPool SetAuthChallengePassword
 return nil
}

func(np *SipIpv4UserAgentPool) GetAuthChallengePassword ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetAuthChallengePassword
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetAuthChallengePasswordModifier () error {
 //parameters: EmulationHandle AuthChallengePasswordModifier Increment Repeat
 //AgtSipIpv4UserAgentPool SetAuthChallengePasswordModifier
 return nil
}

func(np *SipIpv4UserAgentPool) GetAuthChallengePasswordModifier ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetAuthChallengePasswordModifier
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetChallengeOpaque () error {
 //parameters: EmulationHandle ChallengeOpaque
 //AgtSipIpv4UserAgentPool SetChallengeOpaque
 return nil
}

func(np *SipIpv4UserAgentPool) GetChallengeOpaque ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetChallengeOpaque
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetChallengeAlgorithm () error {
 //parameters: EmulationHandle ChallengeAlgorithm
 //AgtSipIpv4UserAgentPool SetChallengeAlgorithm
 return nil
}

func(np *SipIpv4UserAgentPool) GetChallengeAlgorithm ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetChallengeAlgorithm
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetChallengeQop () error {
 //parameters: EmulationHandle ChallengeQop
 //AgtSipIpv4UserAgentPool SetChallengeQop
 return nil
}

func(np *SipIpv4UserAgentPool) GetChallengeQop ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetChallengeQop
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetNonceExpireTime () error {
 //parameters: EmulationHandle NonceExpireTime
 //AgtSipIpv4UserAgentPool SetNonceExpireTime
 return nil
}

func(np *SipIpv4UserAgentPool) GetNonceExpireTime ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetNonceExpireTime
 return "", nil
}

func(np *SipIpv4UserAgentPool) SetMaximumChallengeAttempts () error {
 //parameters: EmulationHandle MaximumChallengeAttempts
 //AgtSipIpv4UserAgentPool SetMaximumChallengeAttempts
 return nil
}

func(np *SipIpv4UserAgentPool) GetMaximumChallengeAttempts ()(string, error) {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool GetMaximumChallengeAttempts
 return "", nil
}

func(np *SipIpv4UserAgentPool) SendTrafficParams () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool SendTrafficParams
 return nil
}

func(np *SipIpv4UserAgentPool) ResetTrafficParams () error {
 //parameters: EmulationHandle
 //AgtSipIpv4UserAgentPool ResetTrafficParams
 return nil
}

