package n2xsdk

type PppoEServerPool struct {
 Handler string
}

func(np *PppoEServerPool) GetServiceName ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetServiceName
 return "", nil
}

func(np *PppoEServerPool) SetServiceName () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetServiceName
 return nil
}

func(np *PppoEServerPool) EnableMaxpayloadTag () error {
 //parameters: DeviceHandle
 //AgtPppoEServerPool EnableMaxpayloadTag
 return nil
}

func(np *PppoEServerPool) GetAcName ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetAcName
 return "", nil
}

func(np *PppoEServerPool) SetAcName () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetAcName
 return nil
}

func(np *PppoEServerPool) DisableMaxpayloadTag () error {
 //parameters: DeviceHandle
 //AgtPppoEServerPool DisableMaxpayloadTag
 return nil
}

func(np *PppoEServerPool) IsMaxpayloadTagEnabled () error {
 //parameters: DeviceHandle
 //AgtPppoEServerPool IsMaxpayloadTagEnabled
 return nil
}

func(np *PppoEServerPool) GetMaxpayloadSize ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetMaxpayloadSize
 return "", nil
}

func(np *PppoEServerPool) SetMaxpayloadSize () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetMaxpayloadSize
 return nil
}

func(np *PppoEServerPool) GetDestinationMacAddress ()(string, error) {
 //parameters: SessionIdentifiers
 //AgtPppoEServerPool GetDestinationMacAddress
 return "", nil
}

func(np *PppoEServerPool) EnableVlan () error {
 //parameters: DeviceHandle TagLevel
 //AgtPppoEServerPool EnableVlan
 return nil
}

func(np *PppoEServerPool) DisableVlan () error {
 //parameters: DeviceHandle TagLevel
 //AgtPppoEServerPool DisableVlan
 return nil
}

func(np *PppoEServerPool) GetVlanFlag ()(string, error) {
 //parameters: DeviceHandle TagLevel
 //AgtPppoEServerPool GetVlanFlag
 return "", nil
}

func(np *PppoEServerPool) SetVlanId () error {
 //parameters: DeviceHandle TagLevel IteratedValue IteratedValueIncrement
 //AgtPppoEServerPool SetVlanId
 return nil
}

func(np *PppoEServerPool) GetVlanId ()(string, error) {
 //parameters: DeviceHandle TagLevel
 //AgtPppoEServerPool GetVlanId
 return "", nil
}

func(np *PppoEServerPool) SetVlanIdRange () error {
 //parameters: DeviceHandle TagLevel IteratedValue IteratedValueCount IteratedValueIncrement IteratedValueRepeat
 //AgtPppoEServerPool SetVlanIdRange
 return nil
}

func(np *PppoEServerPool) GetVlanIdRange ()(string, error) {
 //parameters: DeviceHandle TagLevel
 //AgtPppoEServerPool GetVlanIdRange
 return "", nil
}

func(np *PppoEServerPool) SetVlanIdCountOverride () error {
 //parameters: DeviceHandle TagLevel EnableFlag
 //AgtPppoEServerPool SetVlanIdCountOverride
 return nil
}

func(np *PppoEServerPool) GetVlanIdCountOverride ()(string, error) {
 //parameters: DeviceHandle TagLevel
 //AgtPppoEServerPool GetVlanIdCountOverride
 return "", nil
}

func(np *PppoEServerPool) SetInitialAckTimeout () error {
 //parameters: DeviceHandle Timeout
 //AgtPppoEServerPool SetInitialAckTimeout
 return nil
}

func(np *PppoEServerPool) GetInitialAckTimeout ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetInitialAckTimeout
 return "", nil
}

func(np *PppoEServerPool) SetRetransmitTimeout () error {
 //parameters: DeviceHandle Timeout
 //AgtPppoEServerPool SetRetransmitTimeout
 return nil
}

func(np *PppoEServerPool) GetRetransmitTimeout ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetRetransmitTimeout
 return "", nil
}

func(np *PppoEServerPool) SetMaxTxRetries () error {
 //parameters: DeviceHandle RetransmitAttempts
 //AgtPppoEServerPool SetMaxTxRetries
 return nil
}

func(np *PppoEServerPool) GetMaxTxRetries ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetMaxTxRetries
 return "", nil
}

func(np *PppoEServerPool) GetSessionId ()(string, error) {
 //parameters: SessionIdentifiers
 //AgtPppoEServerPool GetSessionId
 return "", nil
}

func(np *PppoEServerPool) GetSessionIdList ()(string, error) {
 //parameters: SessionIdentifiers
 //AgtPppoEServerPool GetSessionIdList
 return "", nil
}

func(np *PppoEServerPool) GetRemoteMacList ()(string, error) {
 //parameters: SessionIdentifiers
 //AgtPppoEServerPool GetRemoteMacList
 return "", nil
}

func(np *PppoEServerPool) GetIpv6cpEnableFlag ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetIpv6cpEnableFlag
 return "", nil
}

func(np *PppoEServerPool) SetIpv6cpEnableFlag () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetIpv6cpEnableFlag
 return nil
}

func(np *PppoEServerPool) GetServerAddressPoolIncrementingRange ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetServerAddressPoolIncrementingRange
 return "", nil
}

func(np *PppoEServerPool) SetServerAddressPoolIncrementingRange () error {
 //parameters: DeviceHandle IncrementorValue IncrementorValueIncrement
 //AgtPppoEServerPool SetServerAddressPoolIncrementingRange
 return nil
}

func(np *PppoEServerPool) GetNegotiatedIpv6cpIid ()(string, error) {
 //parameters: SessionIdentifiers
 //AgtPppoEServerPool GetNegotiatedIpv6cpIid
 return "", nil
}

func(np *PppoEServerPool) OpenIpcp () error {
 //parameters: DeviceOrSession
 //AgtPppoEServerPool OpenIpcp
 return nil
}

func(np *PppoEServerPool) CloseIpcp () error {
 //parameters: DeviceOrSession
 //AgtPppoEServerPool CloseIpcp
 return nil
}

func(np *PppoEServerPool) OpenIpv6cp () error {
 //parameters: DeviceOrSession
 //AgtPppoEServerPool OpenIpv6cp
 return nil
}

func(np *PppoEServerPool) CloseIpv6cp () error {
 //parameters: DeviceOrSession
 //AgtPppoEServerPool CloseIpv6cp
 return nil
}

func(np *PppoEServerPool) IsMlPppLcpEPDOptionEnabled () error {
 //parameters: DeviceHandle
 //AgtPppoEServerPool IsMlPppLcpEPDOptionEnabled
 return nil
}

func(np *PppoEServerPool) DisableMlPppLcpEPDOption () error {
 //parameters: DeviceHandle
 //AgtPppoEServerPool DisableMlPppLcpEPDOption
 return nil
}

func(np *PppoEServerPool) EnableMlPppLcpEPDOption () error {
 //parameters: DeviceHandle
 //AgtPppoEServerPool EnableMlPppLcpEPDOption
 return nil
}

func(np *PppoEServerPool) GetMlPppLcpEPDClass ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetMlPppLcpEPDClass
 return "", nil
}

func(np *PppoEServerPool) SetMlPppLcpEPDClass () error {
 //parameters: DeviceHandle EpdClassType
 //AgtPppoEServerPool SetMlPppLcpEPDClass
 return nil
}

func(np *PppoEServerPool) GetEPDIpAddressIncrementor ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetEPDIpAddressIncrementor
 return "", nil
}

func(np *PppoEServerPool) SetEPDIpAddressIncrementor () error {
 //parameters: DeviceHandle IncrementorValue IncrementorValueIncrement IncrementorValueRepeat
 //AgtPppoEServerPool SetEPDIpAddressIncrementor
 return nil
}

func(np *PppoEServerPool) GetEPDMacAddressIncrementor ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetEPDMacAddressIncrementor
 return "", nil
}

func(np *PppoEServerPool) SetEPDMacAddressIncrementor () error {
 //parameters: DeviceHandle IncrementorValue IncrementorValueIncrement IncrementorValueRepeat
 //AgtPppoEServerPool SetEPDMacAddressIncrementor
 return nil
}

func(np *PppoEServerPool) GetEnableFlag ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetEnableFlag
 return "", nil
}

func(np *PppoEServerPool) Enable () error {
 //parameters: DeviceHandle
 //AgtPppoEServerPool Enable
 return nil
}

func(np *PppoEServerPool) Disable () error {
 //parameters: DeviceHandle
 //AgtPppoEServerPool Disable
 return nil
}

func(np *PppoEServerPool) Cancel () error {
 //parameters: DeviceOrSession
 //AgtPppoEServerPool Cancel
 return nil
}

func(np *PppoEServerPool) Reset () error {
 //parameters: DeviceOrSession
 //AgtPppoEServerPool Reset
 return nil
}

func(np *PppoEServerPool) Open () error {
 //parameters: DeviceOrSession
 //AgtPppoEServerPool Open
 return nil
}

func(np *PppoEServerPool) Close () error {
 //parameters: DeviceOrSession
 //AgtPppoEServerPool Close
 return nil
}

func(np *PppoEServerPool) CancelAttempts () error {
 //parameters: DeviceOrSession
 //AgtPppoEServerPool CancelAttempts
 return nil
}

func(np *PppoEServerPool) GetNumberOfEstablishedSessions ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetNumberOfEstablishedSessions
 return "", nil
}

func(np *PppoEServerPool) SetConnectionRetryFlags () error {
 //parameters: DeviceHandle ConnectionRetryFlags
 //AgtPppoEServerPool SetConnectionRetryFlags
 return nil
}

func(np *PppoEServerPool) GetConnectionRetryFlags ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetConnectionRetryFlags
 return "", nil
}

func(np *PppoEServerPool) SetReestablishmentFlags () error {
 //parameters: DeviceHandle ReestablishmentFlags
 //AgtPppoEServerPool SetReestablishmentFlags
 return nil
}

func(np *PppoEServerPool) GetReestablishmentFlags ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetReestablishmentFlags
 return "", nil
}

func(np *PppoEServerPool) SetUnlimitedConnectionAttemptsFlag () error {
 //parameters: DeviceHandle EnableFlag
 //AgtPppoEServerPool SetUnlimitedConnectionAttemptsFlag
 return nil
}

func(np *PppoEServerPool) GetUnlimitedConnectionAttemptsFlag ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetUnlimitedConnectionAttemptsFlag
 return "", nil
}

func(np *PppoEServerPool) EnableUnlimitedReestablishment () error {
 //parameters: DeviceHandle
 //AgtPppoEServerPool EnableUnlimitedReestablishment
 return nil
}

func(np *PppoEServerPool) DisableUnlimitedReestablishment () error {
 //parameters: DeviceHandle
 //AgtPppoEServerPool DisableUnlimitedReestablishment
 return nil
}

func(np *PppoEServerPool) IsUnlimitedReestablishmentEnabled () error {
 //parameters: DeviceHandle
 //AgtPppoEServerPool IsUnlimitedReestablishmentEnabled
 return nil
}

func(np *PppoEServerPool) GetStateCount ()(string, error) {
 //parameters: DeviceHandle SessionStateName
 //AgtPppoEServerPool GetStateCount
 return "", nil
}

func(np *PppoEServerPool) GetSessionState ()(string, error) {
 //parameters: SessionIdentifiers
 //AgtPppoEServerPool GetSessionState
 return "", nil
}

func(np *PppoEServerPool) GetSessionSubprotocolState ()(string, error) {
 //parameters: SessionIdentifiers SubprotocolName
 //AgtPppoEServerPool GetSessionSubprotocolState
 return "", nil
}

func(np *PppoEServerPool) GetNumberOfSessions ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetNumberOfSessions
 return "", nil
}

func(np *PppoEServerPool) SetNumberOfSessions () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetNumberOfSessions
 return nil
}

func(np *PppoEServerPool) GetAutoStartFlag ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetAutoStartFlag
 return "", nil
}

func(np *PppoEServerPool) SetAutoStartFlag () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetAutoStartFlag
 return nil
}

func(np *PppoEServerPool) GetSessionLifetime ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetSessionLifetime
 return "", nil
}

func(np *PppoEServerPool) SetSessionLifetime () error {
 //parameters: DeviceHandle Value Value
 //AgtPppoEServerPool SetSessionLifetime
 return nil
}

func(np *PppoEServerPool) GetOpeningWindow ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetOpeningWindow
 return "", nil
}

func(np *PppoEServerPool) SetOpeningWindow () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetOpeningWindow
 return nil
}

func(np *PppoEServerPool) GetMaxConnectionAttempts ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetMaxConnectionAttempts
 return "", nil
}

func(np *PppoEServerPool) SetMaxConnectionAttempts () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetMaxConnectionAttempts
 return nil
}

func(np *PppoEServerPool) GetConnectionAttempts ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetConnectionAttempts
 return "", nil
}

func(np *PppoEServerPool) GetMaxReestablishments ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetMaxReestablishments
 return "", nil
}

func(np *PppoEServerPool) SetMaxReestablishments () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetMaxReestablishments
 return nil
}

func(np *PppoEServerPool) RetryConnections () error {
 //parameters: DeviceOrSession
 //AgtPppoEServerPool RetryConnections
 return nil
}

func(np *PppoEServerPool) GetInitiationRateLimitFlag ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetInitiationRateLimitFlag
 return "", nil
}

func(np *PppoEServerPool) SetInitiationRateLimitFlag () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetInitiationRateLimitFlag
 return nil
}

func(np *PppoEServerPool) GetSessionInitiationRate ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetSessionInitiationRate
 return "", nil
}

func(np *PppoEServerPool) SetSessionInitiationRate () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetSessionInitiationRate
 return nil
}

func(np *PppoEServerPool) GetLimitOpeningRateMode ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetLimitOpeningRateMode
 return "", nil
}

func(np *PppoEServerPool) SetLimitOpeningRateMode () error {
 //parameters: DeviceHandle SetupInterval
 //AgtPppoEServerPool SetLimitOpeningRateMode
 return nil
}

func(np *PppoEServerPool) GetSessionInitiationRateInMilliseconds ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetSessionInitiationRateInMilliseconds
 return "", nil
}

func(np *PppoEServerPool) SetSessionInitiationRateInMilliseconds () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetSessionInitiationRateInMilliseconds
 return nil
}

func(np *PppoEServerPool) GetLcpOption ()(string, error) {
 //parameters: DeviceHandle PppOption
 //AgtPppoEServerPool GetLcpOption
 return "", nil
}

func(np *PppoEServerPool) SetLcpOption () error {
 //parameters: DeviceHandle PppOption Value
 //AgtPppoEServerPool SetLcpOption
 return nil
}

func(np *PppoEServerPool) GetPppOption ()(string, error) {
 //parameters: DeviceHandle PppOption
 //AgtPppoEServerPool GetPppOption
 return "", nil
}

func(np *PppoEServerPool) SetPppOption () error {
 //parameters: DeviceHandle PppOption Value
 //AgtPppoEServerPool SetPppOption
 return nil
}

func(np *PppoEServerPool) IsMlPppLcpMrruOptionEnabled () error {
 //parameters: DeviceHandle
 //AgtPppoEServerPool IsMlPppLcpMrruOptionEnabled
 return nil
}

func(np *PppoEServerPool) DisableMlPppLcpMrruOption () error {
 //parameters: DeviceHandle
 //AgtPppoEServerPool DisableMlPppLcpMrruOption
 return nil
}

func(np *PppoEServerPool) EnableMlPppLcpMrruOption () error {
 //parameters: DeviceHandle
 //AgtPppoEServerPool EnableMlPppLcpMrruOption
 return nil
}

func(np *PppoEServerPool) GetMlPppLcpMrruSize ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetMlPppLcpMrruSize
 return "", nil
}

func(np *PppoEServerPool) SetMlPppLcpMrruSize () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetMlPppLcpMrruSize
 return nil
}

func(np *PppoEServerPool) GetIpcpEnableFlag ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetIpcpEnableFlag
 return "", nil
}

func(np *PppoEServerPool) SetIpcpEnableFlag () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetIpcpEnableFlag
 return nil
}

func(np *PppoEServerPool) GetOfferNetmaskFlag ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetOfferNetmaskFlag
 return "", nil
}

func(np *PppoEServerPool) SetOfferNetmaskFlag () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetOfferNetmaskFlag
 return nil
}

func(np *PppoEServerPool) GetMandatoryNetmaskNegotiationFlag ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetMandatoryNetmaskNegotiationFlag
 return "", nil
}

func(np *PppoEServerPool) SetMandatoryNetmaskNegotiationFlag () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetMandatoryNetmaskNegotiationFlag
 return nil
}

func(np *PppoEServerPool) GetOfferedNetmaskLength ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetOfferedNetmaskLength
 return "", nil
}

func(np *PppoEServerPool) SetOfferedNetmaskLength () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetOfferedNetmaskLength
 return nil
}

func(np *PppoEServerPool) GetOfferedNetmaskLengthModifier ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetOfferedNetmaskLengthModifier
 return "", nil
}

func(np *PppoEServerPool) SetOfferedNetmaskLengthModifier () error {
 //parameters: DeviceHandle Value Repeat Increment
 //AgtPppoEServerPool SetOfferedNetmaskLengthModifier
 return nil
}

func(np *PppoEServerPool) GetOfferNameServerAddressFlag ()(string, error) {
 //parameters: DeviceHandle NameAddressServerType
 //AgtPppoEServerPool GetOfferNameServerAddressFlag
 return "", nil
}

func(np *PppoEServerPool) SetOfferNameServerAddressFlag () error {
 //parameters: DeviceHandle NameAddressServerType Value
 //AgtPppoEServerPool SetOfferNameServerAddressFlag
 return nil
}

func(np *PppoEServerPool) GetMandatoryNameServerAddressNegotiationFlag ()(string, error) {
 //parameters: DeviceHandle NameAddressServerType
 //AgtPppoEServerPool GetMandatoryNameServerAddressNegotiationFlag
 return "", nil
}

func(np *PppoEServerPool) SetMandatoryNameServerAddressNegotiationFlag () error {
 //parameters: DeviceHandle NameAddressServerType Value
 //AgtPppoEServerPool SetMandatoryNameServerAddressNegotiationFlag
 return nil
}

func(np *PppoEServerPool) GetOfferedNameServerAddress ()(string, error) {
 //parameters: DeviceHandle NameAddressServerType
 //AgtPppoEServerPool GetOfferedNameServerAddress
 return "", nil
}

func(np *PppoEServerPool) SetOfferedNameServerAddress () error {
 //parameters: DeviceHandle NameAddressServerType Value
 //AgtPppoEServerPool SetOfferedNameServerAddress
 return nil
}

func(np *PppoEServerPool) GetOfferedNameServerAddressIncrementingValue ()(string, error) {
 //parameters: DeviceHandle NameAddressServerType
 //AgtPppoEServerPool GetOfferedNameServerAddressIncrementingValue
 return "", nil
}

func(np *PppoEServerPool) SetOfferedNameServerAddressIncrementingValue () error {
 //parameters: DeviceHandle NameAddressServerType Value Repeat Increment
 //AgtPppoEServerPool SetOfferedNameServerAddressIncrementingValue
 return nil
}

func(np *PppoEServerPool) GetOfferedNameServerAddressList ()(string, error) {
 //parameters: DeviceHandle NameAddressServerType
 //AgtPppoEServerPool GetOfferedNameServerAddressList
 return "", nil
}

func(np *PppoEServerPool) SetOfferedNameServerAddressList ()(string, error) {
 //parameters: DeviceHandle NameAddressServerType Value
 //AgtPppoEServerPool SetOfferedNameServerAddressList
 return "", nil
}

func(np *PppoEServerPool) GetAuthenticationProtocol ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetAuthenticationProtocol
 return "", nil
}

func(np *PppoEServerPool) SetAuthenticationProtocol () error {
 //parameters: DeviceHandle AuthProtocol
 //AgtPppoEServerPool SetAuthenticationProtocol
 return nil
}

func(np *PppoEServerPool) SetAuthenticationRestartTimer () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetAuthenticationRestartTimer
 return nil
}

func(np *PppoEServerPool) GetAuthenticationRestartTimer ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetAuthenticationRestartTimer
 return "", nil
}

func(np *PppoEServerPool) SetAuthenticationMaximumAttempts () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetAuthenticationMaximumAttempts
 return nil
}

func(np *PppoEServerPool) GetAuthenticationMaximumAttempts ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetAuthenticationMaximumAttempts
 return "", nil
}

func(np *PppoEServerPool) GetTesterCredentials ()(string, error) {
 //parameters: DeviceHandle CredentialParameterName
 //AgtPppoEServerPool GetTesterCredentials
 return "", nil
}

func(np *PppoEServerPool) SetTesterCredentials () error {
 //parameters: DeviceHandle CredentialParameterName Value
 //AgtPppoEServerPool SetTesterCredentials
 return nil
}

func(np *PppoEServerPool) GetTesterCredentialsIncrementor ()(string, error) {
 //parameters: DeviceHandle CredentialParameterName
 //AgtPppoEServerPool GetTesterCredentialsIncrementor
 return "", nil
}

func(np *PppoEServerPool) SetTesterCredentialsIncrementor () error {
 //parameters: DeviceHandle CredentialParameterName IteratedValue IteratedValueIncrement IteratedValueRepeat
 //AgtPppoEServerPool SetTesterCredentialsIncrementor
 return nil
}

func(np *PppoEServerPool) AddTesterCredentialsNumberedIncrementor () error {
 //parameters: DeviceHandle CredentialParameterName
 //AgtPppoEServerPool AddTesterCredentialsNumberedIncrementor
 return nil
}

func(np *PppoEServerPool) RemoveTesterCredentialsNumberedIncrementor () error {
 //parameters: DeviceHandle CredentialParameterName IncrementorNumber
 //AgtPppoEServerPool RemoveTesterCredentialsNumberedIncrementor
 return nil
}

func(np *PppoEServerPool) GetTesterCredentialsNumberedIncrementor ()(string, error) {
 //parameters: DeviceHandle CredentialParameterName IncrementorNumber
 //AgtPppoEServerPool GetTesterCredentialsNumberedIncrementor
 return "", nil
}

func(np *PppoEServerPool) SetTesterCredentialsNumberedIncrementor () error {
 //parameters: DeviceHandle CredentialParameterName IncrementorNumber IteratedValue IteratedValueCount IteratedValueIncrement IteratedValueRepeat
 //AgtPppoEServerPool SetTesterCredentialsNumberedIncrementor
 return nil
}

func(np *PppoEServerPool) GetTesterCredentialsIncrementorCountOverride ()(string, error) {
 //parameters: DeviceHandle CredentialParameterName IncrementorNumber
 //AgtPppoEServerPool GetTesterCredentialsIncrementorCountOverride
 return "", nil
}

func(np *PppoEServerPool) SetTesterCredentialsIncrementorCountOverride () error {
 //parameters: DeviceHandle CredentialParameterName IncrementorNumber IsIncrementCountOverridden
 //AgtPppoEServerPool SetTesterCredentialsIncrementorCountOverride
 return nil
}

func(np *PppoEServerPool) EnableAcceptAnySutCredentialFlag () error {
 //parameters: DeviceHandle
 //AgtPppoEServerPool EnableAcceptAnySutCredentialFlag
 return nil
}

func(np *PppoEServerPool) DisableAcceptAnySutCredentialFlag () error {
 //parameters: DeviceHandle
 //AgtPppoEServerPool DisableAcceptAnySutCredentialFlag
 return nil
}

func(np *PppoEServerPool) IsAcceptAnySutCredentialFlagEnabled () error {
 //parameters: DeviceHandle
 //AgtPppoEServerPool IsAcceptAnySutCredentialFlagEnabled
 return nil
}

func(np *PppoEServerPool) SetAcceptedSutNameOrPeerIdValue () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetAcceptedSutNameOrPeerIdValue
 return nil
}

func(np *PppoEServerPool) GetAcceptedSutNameOrPeerIdValue ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetAcceptedSutNameOrPeerIdValue
 return "", nil
}

func(np *PppoEServerPool) SetAcceptedSutPasswordOrSecretValue () error {
 //parameters: DeviceHandle Value
 //AgtPppoEServerPool SetAcceptedSutPasswordOrSecretValue
 return nil
}

func(np *PppoEServerPool) GetAcceptedSutPasswordOrSecretValue ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetAcceptedSutPasswordOrSecretValue
 return "", nil
}

func(np *PppoEServerPool) SetAcceptedSutNameOrPeerIdModifier () error {
 //parameters: DeviceHandle Value Repeat Increment
 //AgtPppoEServerPool SetAcceptedSutNameOrPeerIdModifier
 return nil
}

func(np *PppoEServerPool) GetAcceptedSutNameOrPeerIdModifier ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetAcceptedSutNameOrPeerIdModifier
 return "", nil
}

func(np *PppoEServerPool) SetAcceptedSutPasswordOrSecretModifier () error {
 //parameters: DeviceHandle Value Repeat Increment
 //AgtPppoEServerPool SetAcceptedSutPasswordOrSecretModifier
 return nil
}

func(np *PppoEServerPool) GetAcceptedSutPasswordOrSecretModifier ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetAcceptedSutPasswordOrSecretModifier
 return "", nil
}

func(np *PppoEServerPool) ListSelectedPools ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool ListSelectedPools
 return "", nil
}

func(np *PppoEServerPool) SelectPools () error {
 //parameters: SessionIdentifiers
 //AgtPppoEServerPool SelectPools
 return nil
}

func(np *PppoEServerPool) DeselectPools () error {
 //parameters: SessionIdentifiers
 //AgtPppoEServerPool DeselectPools
 return nil
}

func(np *PppoEServerPool) GetAccumulatedValues ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoEServerPool GetAccumulatedValues
 return "", nil
}

func(np *PppoEServerPool) GetDeviceAccumulatedValues ()(string, error) {
 //parameters: DeviceHandle DeviceStatisticsList
 //AgtPppoEServerPool GetDeviceAccumulatedValues
 return "", nil
}

func(np *PppoEServerPool) GetSessionAccumulatedSpecifiedValues ()(string, error) {
 //parameters: SessionIdentifiers SessionStatisticsList
 //AgtPppoEServerPool GetSessionAccumulatedSpecifiedValues
 return "", nil
}

func(np *PppoEServerPool) ClearStatistics () error {
 //parameters: DeviceHandle
 //AgtPppoEServerPool ClearStatistics
 return nil
}

func(np *PppoEServerPool) GetLastFailureReason ()(string, error) {
 //parameters: SessionIdentifiers
 //AgtPppoEServerPool GetLastFailureReason
 return "", nil
}

func(np *PppoEServerPool) GetVlanPriority ()(string, error) {
 //parameters: DeviceHandle VlanTagIndex
 //AgtPppoEServerPool GetVlanPriority
 return "", nil
}

func(np *PppoEServerPool) SetVlanPriority () error {
 //parameters: DeviceHandle VlanTagIndex VlanPriority
 //AgtPppoEServerPool SetVlanPriority
 return nil
}

