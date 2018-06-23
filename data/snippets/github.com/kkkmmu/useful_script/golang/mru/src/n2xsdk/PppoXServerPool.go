package n2xsdk

type PppoXServerPool struct {
 Handler string
}

func(np *PppoXServerPool) GetEnableFlag ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetEnableFlag
 return "", nil
}

func(np *PppoXServerPool) Enable () error {
 //parameters: DeviceHandle
 //AgtPppoXServerPool Enable
 return nil
}

func(np *PppoXServerPool) Disable () error {
 //parameters: DeviceHandle
 //AgtPppoXServerPool Disable
 return nil
}

func(np *PppoXServerPool) Cancel () error {
 //parameters: DeviceOrSession
 //AgtPppoXServerPool Cancel
 return nil
}

func(np *PppoXServerPool) Reset () error {
 //parameters: DeviceOrSession
 //AgtPppoXServerPool Reset
 return nil
}

func(np *PppoXServerPool) Open () error {
 //parameters: DeviceOrSession
 //AgtPppoXServerPool Open
 return nil
}

func(np *PppoXServerPool) Close () error {
 //parameters: DeviceOrSession
 //AgtPppoXServerPool Close
 return nil
}

func(np *PppoXServerPool) CancelAttempts () error {
 //parameters: DeviceOrSession
 //AgtPppoXServerPool CancelAttempts
 return nil
}

func(np *PppoXServerPool) GetNumberOfEstablishedSessions ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetNumberOfEstablishedSessions
 return "", nil
}

func(np *PppoXServerPool) SetConnectionRetryFlags () error {
 //parameters: DeviceHandle ConnectionRetryFlags
 //AgtPppoXServerPool SetConnectionRetryFlags
 return nil
}

func(np *PppoXServerPool) GetConnectionRetryFlags ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetConnectionRetryFlags
 return "", nil
}

func(np *PppoXServerPool) SetReestablishmentFlags () error {
 //parameters: DeviceHandle ReestablishmentFlags
 //AgtPppoXServerPool SetReestablishmentFlags
 return nil
}

func(np *PppoXServerPool) GetReestablishmentFlags ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetReestablishmentFlags
 return "", nil
}

func(np *PppoXServerPool) SetUnlimitedConnectionAttemptsFlag () error {
 //parameters: DeviceHandle EnableFlag
 //AgtPppoXServerPool SetUnlimitedConnectionAttemptsFlag
 return nil
}

func(np *PppoXServerPool) GetUnlimitedConnectionAttemptsFlag ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetUnlimitedConnectionAttemptsFlag
 return "", nil
}

func(np *PppoXServerPool) EnableUnlimitedReestablishment () error {
 //parameters: DeviceHandle
 //AgtPppoXServerPool EnableUnlimitedReestablishment
 return nil
}

func(np *PppoXServerPool) DisableUnlimitedReestablishment () error {
 //parameters: DeviceHandle
 //AgtPppoXServerPool DisableUnlimitedReestablishment
 return nil
}

func(np *PppoXServerPool) IsUnlimitedReestablishmentEnabled () error {
 //parameters: DeviceHandle
 //AgtPppoXServerPool IsUnlimitedReestablishmentEnabled
 return nil
}

func(np *PppoXServerPool) GetStateCount ()(string, error) {
 //parameters: DeviceHandle SessionStateName
 //AgtPppoXServerPool GetStateCount
 return "", nil
}

func(np *PppoXServerPool) GetSessionState ()(string, error) {
 //parameters: SessionIdentifiers
 //AgtPppoXServerPool GetSessionState
 return "", nil
}

func(np *PppoXServerPool) GetSessionSubprotocolState ()(string, error) {
 //parameters: SessionIdentifiers SubprotocolName
 //AgtPppoXServerPool GetSessionSubprotocolState
 return "", nil
}

func(np *PppoXServerPool) GetNumberOfSessions ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetNumberOfSessions
 return "", nil
}

func(np *PppoXServerPool) SetNumberOfSessions () error {
 //parameters: DeviceHandle Value
 //AgtPppoXServerPool SetNumberOfSessions
 return nil
}

func(np *PppoXServerPool) GetAutoStartFlag ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetAutoStartFlag
 return "", nil
}

func(np *PppoXServerPool) SetAutoStartFlag () error {
 //parameters: DeviceHandle Value
 //AgtPppoXServerPool SetAutoStartFlag
 return nil
}

func(np *PppoXServerPool) GetSessionLifetime ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetSessionLifetime
 return "", nil
}

func(np *PppoXServerPool) SetSessionLifetime () error {
 //parameters: DeviceHandle Value Value
 //AgtPppoXServerPool SetSessionLifetime
 return nil
}

func(np *PppoXServerPool) GetOpeningWindow ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetOpeningWindow
 return "", nil
}

func(np *PppoXServerPool) SetOpeningWindow () error {
 //parameters: DeviceHandle Value
 //AgtPppoXServerPool SetOpeningWindow
 return nil
}

func(np *PppoXServerPool) GetMaxConnectionAttempts ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetMaxConnectionAttempts
 return "", nil
}

func(np *PppoXServerPool) SetMaxConnectionAttempts () error {
 //parameters: DeviceHandle Value
 //AgtPppoXServerPool SetMaxConnectionAttempts
 return nil
}

func(np *PppoXServerPool) GetConnectionAttempts ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetConnectionAttempts
 return "", nil
}

func(np *PppoXServerPool) GetMaxReestablishments ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetMaxReestablishments
 return "", nil
}

func(np *PppoXServerPool) SetMaxReestablishments () error {
 //parameters: DeviceHandle Value
 //AgtPppoXServerPool SetMaxReestablishments
 return nil
}

func(np *PppoXServerPool) RetryConnections () error {
 //parameters: DeviceOrSession
 //AgtPppoXServerPool RetryConnections
 return nil
}

func(np *PppoXServerPool) GetInitiationRateLimitFlag ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetInitiationRateLimitFlag
 return "", nil
}

func(np *PppoXServerPool) SetInitiationRateLimitFlag () error {
 //parameters: DeviceHandle Value
 //AgtPppoXServerPool SetInitiationRateLimitFlag
 return nil
}

func(np *PppoXServerPool) GetSessionInitiationRate ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetSessionInitiationRate
 return "", nil
}

func(np *PppoXServerPool) SetSessionInitiationRate () error {
 //parameters: DeviceHandle Value
 //AgtPppoXServerPool SetSessionInitiationRate
 return nil
}

func(np *PppoXServerPool) GetLimitOpeningRateMode ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetLimitOpeningRateMode
 return "", nil
}

func(np *PppoXServerPool) SetLimitOpeningRateMode () error {
 //parameters: DeviceHandle SetupInterval
 //AgtPppoXServerPool SetLimitOpeningRateMode
 return nil
}

func(np *PppoXServerPool) GetSessionInitiationRateInMilliseconds ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetSessionInitiationRateInMilliseconds
 return "", nil
}

func(np *PppoXServerPool) SetSessionInitiationRateInMilliseconds () error {
 //parameters: DeviceHandle Value
 //AgtPppoXServerPool SetSessionInitiationRateInMilliseconds
 return nil
}

func(np *PppoXServerPool) GetLcpOption ()(string, error) {
 //parameters: DeviceHandle PppOption
 //AgtPppoXServerPool GetLcpOption
 return "", nil
}

func(np *PppoXServerPool) SetLcpOption () error {
 //parameters: DeviceHandle PppOption Value
 //AgtPppoXServerPool SetLcpOption
 return nil
}

func(np *PppoXServerPool) GetPppOption ()(string, error) {
 //parameters: DeviceHandle PppOption
 //AgtPppoXServerPool GetPppOption
 return "", nil
}

func(np *PppoXServerPool) SetPppOption () error {
 //parameters: DeviceHandle PppOption Value
 //AgtPppoXServerPool SetPppOption
 return nil
}

func(np *PppoXServerPool) IsMlPppLcpMrruOptionEnabled () error {
 //parameters: DeviceHandle
 //AgtPppoXServerPool IsMlPppLcpMrruOptionEnabled
 return nil
}

func(np *PppoXServerPool) DisableMlPppLcpMrruOption () error {
 //parameters: DeviceHandle
 //AgtPppoXServerPool DisableMlPppLcpMrruOption
 return nil
}

func(np *PppoXServerPool) EnableMlPppLcpMrruOption () error {
 //parameters: DeviceHandle
 //AgtPppoXServerPool EnableMlPppLcpMrruOption
 return nil
}

func(np *PppoXServerPool) GetMlPppLcpMrruSize ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetMlPppLcpMrruSize
 return "", nil
}

func(np *PppoXServerPool) SetMlPppLcpMrruSize () error {
 //parameters: DeviceHandle Value
 //AgtPppoXServerPool SetMlPppLcpMrruSize
 return nil
}

func(np *PppoXServerPool) GetIpcpEnableFlag ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetIpcpEnableFlag
 return "", nil
}

func(np *PppoXServerPool) SetIpcpEnableFlag () error {
 //parameters: DeviceHandle Value
 //AgtPppoXServerPool SetIpcpEnableFlag
 return nil
}

func(np *PppoXServerPool) GetOfferNetmaskFlag ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetOfferNetmaskFlag
 return "", nil
}

func(np *PppoXServerPool) SetOfferNetmaskFlag () error {
 //parameters: DeviceHandle Value
 //AgtPppoXServerPool SetOfferNetmaskFlag
 return nil
}

func(np *PppoXServerPool) GetMandatoryNetmaskNegotiationFlag ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetMandatoryNetmaskNegotiationFlag
 return "", nil
}

func(np *PppoXServerPool) SetMandatoryNetmaskNegotiationFlag () error {
 //parameters: DeviceHandle Value
 //AgtPppoXServerPool SetMandatoryNetmaskNegotiationFlag
 return nil
}

func(np *PppoXServerPool) GetOfferedNetmaskLength ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetOfferedNetmaskLength
 return "", nil
}

func(np *PppoXServerPool) SetOfferedNetmaskLength () error {
 //parameters: DeviceHandle Value
 //AgtPppoXServerPool SetOfferedNetmaskLength
 return nil
}

func(np *PppoXServerPool) GetOfferedNetmaskLengthModifier ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetOfferedNetmaskLengthModifier
 return "", nil
}

func(np *PppoXServerPool) SetOfferedNetmaskLengthModifier () error {
 //parameters: DeviceHandle Value Repeat Increment
 //AgtPppoXServerPool SetOfferedNetmaskLengthModifier
 return nil
}

func(np *PppoXServerPool) GetOfferNameServerAddressFlag ()(string, error) {
 //parameters: DeviceHandle NameAddressServerType
 //AgtPppoXServerPool GetOfferNameServerAddressFlag
 return "", nil
}

func(np *PppoXServerPool) SetOfferNameServerAddressFlag () error {
 //parameters: DeviceHandle NameAddressServerType Value
 //AgtPppoXServerPool SetOfferNameServerAddressFlag
 return nil
}

func(np *PppoXServerPool) GetMandatoryNameServerAddressNegotiationFlag ()(string, error) {
 //parameters: DeviceHandle NameAddressServerType
 //AgtPppoXServerPool GetMandatoryNameServerAddressNegotiationFlag
 return "", nil
}

func(np *PppoXServerPool) SetMandatoryNameServerAddressNegotiationFlag () error {
 //parameters: DeviceHandle NameAddressServerType Value
 //AgtPppoXServerPool SetMandatoryNameServerAddressNegotiationFlag
 return nil
}

func(np *PppoXServerPool) GetOfferedNameServerAddress ()(string, error) {
 //parameters: DeviceHandle NameAddressServerType
 //AgtPppoXServerPool GetOfferedNameServerAddress
 return "", nil
}

func(np *PppoXServerPool) SetOfferedNameServerAddress () error {
 //parameters: DeviceHandle NameAddressServerType Value
 //AgtPppoXServerPool SetOfferedNameServerAddress
 return nil
}

func(np *PppoXServerPool) GetOfferedNameServerAddressIncrementingValue ()(string, error) {
 //parameters: DeviceHandle NameAddressServerType
 //AgtPppoXServerPool GetOfferedNameServerAddressIncrementingValue
 return "", nil
}

func(np *PppoXServerPool) SetOfferedNameServerAddressIncrementingValue () error {
 //parameters: DeviceHandle NameAddressServerType Value Repeat Increment
 //AgtPppoXServerPool SetOfferedNameServerAddressIncrementingValue
 return nil
}

func(np *PppoXServerPool) GetOfferedNameServerAddressList ()(string, error) {
 //parameters: DeviceHandle NameAddressServerType
 //AgtPppoXServerPool GetOfferedNameServerAddressList
 return "", nil
}

func(np *PppoXServerPool) SetOfferedNameServerAddressList ()(string, error) {
 //parameters: DeviceHandle NameAddressServerType Value
 //AgtPppoXServerPool SetOfferedNameServerAddressList
 return "", nil
}

func(np *PppoXServerPool) GetAuthenticationProtocol ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetAuthenticationProtocol
 return "", nil
}

func(np *PppoXServerPool) SetAuthenticationProtocol () error {
 //parameters: DeviceHandle AuthProtocol
 //AgtPppoXServerPool SetAuthenticationProtocol
 return nil
}

func(np *PppoXServerPool) SetAuthenticationRestartTimer () error {
 //parameters: DeviceHandle Value
 //AgtPppoXServerPool SetAuthenticationRestartTimer
 return nil
}

func(np *PppoXServerPool) GetAuthenticationRestartTimer ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetAuthenticationRestartTimer
 return "", nil
}

func(np *PppoXServerPool) SetAuthenticationMaximumAttempts () error {
 //parameters: DeviceHandle Value
 //AgtPppoXServerPool SetAuthenticationMaximumAttempts
 return nil
}

func(np *PppoXServerPool) GetAuthenticationMaximumAttempts ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetAuthenticationMaximumAttempts
 return "", nil
}

func(np *PppoXServerPool) GetTesterCredentials ()(string, error) {
 //parameters: DeviceHandle CredentialParameterName
 //AgtPppoXServerPool GetTesterCredentials
 return "", nil
}

func(np *PppoXServerPool) SetTesterCredentials () error {
 //parameters: DeviceHandle CredentialParameterName Value
 //AgtPppoXServerPool SetTesterCredentials
 return nil
}

func(np *PppoXServerPool) GetTesterCredentialsIncrementor ()(string, error) {
 //parameters: DeviceHandle CredentialParameterName
 //AgtPppoXServerPool GetTesterCredentialsIncrementor
 return "", nil
}

func(np *PppoXServerPool) SetTesterCredentialsIncrementor () error {
 //parameters: DeviceHandle CredentialParameterName IteratedValue IteratedValueIncrement IteratedValueRepeat
 //AgtPppoXServerPool SetTesterCredentialsIncrementor
 return nil
}

func(np *PppoXServerPool) AddTesterCredentialsNumberedIncrementor () error {
 //parameters: DeviceHandle CredentialParameterName
 //AgtPppoXServerPool AddTesterCredentialsNumberedIncrementor
 return nil
}

func(np *PppoXServerPool) RemoveTesterCredentialsNumberedIncrementor () error {
 //parameters: DeviceHandle CredentialParameterName IncrementorNumber
 //AgtPppoXServerPool RemoveTesterCredentialsNumberedIncrementor
 return nil
}

func(np *PppoXServerPool) GetTesterCredentialsNumberedIncrementor ()(string, error) {
 //parameters: DeviceHandle CredentialParameterName IncrementorNumber
 //AgtPppoXServerPool GetTesterCredentialsNumberedIncrementor
 return "", nil
}

func(np *PppoXServerPool) SetTesterCredentialsNumberedIncrementor () error {
 //parameters: DeviceHandle CredentialParameterName IncrementorNumber IteratedValue IteratedValueCount IteratedValueIncrement IteratedValueRepeat
 //AgtPppoXServerPool SetTesterCredentialsNumberedIncrementor
 return nil
}

func(np *PppoXServerPool) GetTesterCredentialsIncrementorCountOverride ()(string, error) {
 //parameters: DeviceHandle CredentialParameterName IncrementorNumber
 //AgtPppoXServerPool GetTesterCredentialsIncrementorCountOverride
 return "", nil
}

func(np *PppoXServerPool) SetTesterCredentialsIncrementorCountOverride () error {
 //parameters: DeviceHandle CredentialParameterName IncrementorNumber IsIncrementCountOverridden
 //AgtPppoXServerPool SetTesterCredentialsIncrementorCountOverride
 return nil
}

func(np *PppoXServerPool) EnableAcceptAnySutCredentialFlag () error {
 //parameters: DeviceHandle
 //AgtPppoXServerPool EnableAcceptAnySutCredentialFlag
 return nil
}

func(np *PppoXServerPool) DisableAcceptAnySutCredentialFlag () error {
 //parameters: DeviceHandle
 //AgtPppoXServerPool DisableAcceptAnySutCredentialFlag
 return nil
}

func(np *PppoXServerPool) IsAcceptAnySutCredentialFlagEnabled () error {
 //parameters: DeviceHandle
 //AgtPppoXServerPool IsAcceptAnySutCredentialFlagEnabled
 return nil
}

func(np *PppoXServerPool) SetAcceptedSutNameOrPeerIdValue () error {
 //parameters: DeviceHandle Value
 //AgtPppoXServerPool SetAcceptedSutNameOrPeerIdValue
 return nil
}

func(np *PppoXServerPool) GetAcceptedSutNameOrPeerIdValue ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetAcceptedSutNameOrPeerIdValue
 return "", nil
}

func(np *PppoXServerPool) SetAcceptedSutPasswordOrSecretValue () error {
 //parameters: DeviceHandle Value
 //AgtPppoXServerPool SetAcceptedSutPasswordOrSecretValue
 return nil
}

func(np *PppoXServerPool) GetAcceptedSutPasswordOrSecretValue ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetAcceptedSutPasswordOrSecretValue
 return "", nil
}

func(np *PppoXServerPool) SetAcceptedSutNameOrPeerIdModifier () error {
 //parameters: DeviceHandle Value Repeat Increment
 //AgtPppoXServerPool SetAcceptedSutNameOrPeerIdModifier
 return nil
}

func(np *PppoXServerPool) GetAcceptedSutNameOrPeerIdModifier ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetAcceptedSutNameOrPeerIdModifier
 return "", nil
}

func(np *PppoXServerPool) SetAcceptedSutPasswordOrSecretModifier () error {
 //parameters: DeviceHandle Value Repeat Increment
 //AgtPppoXServerPool SetAcceptedSutPasswordOrSecretModifier
 return nil
}

func(np *PppoXServerPool) GetAcceptedSutPasswordOrSecretModifier ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetAcceptedSutPasswordOrSecretModifier
 return "", nil
}

func(np *PppoXServerPool) ListSelectedPools ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool ListSelectedPools
 return "", nil
}

func(np *PppoXServerPool) SelectPools () error {
 //parameters: SessionIdentifiers
 //AgtPppoXServerPool SelectPools
 return nil
}

func(np *PppoXServerPool) DeselectPools () error {
 //parameters: SessionIdentifiers
 //AgtPppoXServerPool DeselectPools
 return nil
}

func(np *PppoXServerPool) GetAccumulatedValues ()(string, error) {
 //parameters: DeviceHandle
 //AgtPppoXServerPool GetAccumulatedValues
 return "", nil
}

func(np *PppoXServerPool) GetDeviceAccumulatedValues ()(string, error) {
 //parameters: DeviceHandle DeviceStatisticsList
 //AgtPppoXServerPool GetDeviceAccumulatedValues
 return "", nil
}

func(np *PppoXServerPool) GetSessionAccumulatedSpecifiedValues ()(string, error) {
 //parameters: SessionIdentifiers SessionStatisticsList
 //AgtPppoXServerPool GetSessionAccumulatedSpecifiedValues
 return "", nil
}

func(np *PppoXServerPool) ClearStatistics () error {
 //parameters: DeviceHandle
 //AgtPppoXServerPool ClearStatistics
 return nil
}

func(np *PppoXServerPool) GetLastFailureReason ()(string, error) {
 //parameters: SessionIdentifiers
 //AgtPppoXServerPool GetLastFailureReason
 return "", nil
}

