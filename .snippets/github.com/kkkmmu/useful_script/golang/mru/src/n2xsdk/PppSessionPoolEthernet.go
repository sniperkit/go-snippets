package n2xsdk

type PppSessionPoolEthernet struct {
 Handler string
}

func(np *PppSessionPoolEtherne) SetNumberOfSessionsPerSrcMacAddr () error {
 //parameters: SessionPoolHandle SessionsPerMacAddr
 //AgtPppSessionPoolEthernet SetNumberOfSessionsPerSrcMacAddr
 return nil
}

func(np *PppSessionPoolEtherne) SetInitialAckTimeout () error {
 //parameters: SessionPoolHandle Timeout
 //AgtPppSessionPoolEthernet SetInitialAckTimeout
 return nil
}

func(np *PppSessionPoolEtherne) SetMaxTxRetries () error {
 //parameters: SessionPoolHandle Retries
 //AgtPppSessionPoolEthernet SetMaxTxRetries
 return nil
}

func(np *PppSessionPoolEtherne) EnableVlan () error {
 //parameters: SessionPoolHandle VlanTagNumber
 //AgtPppSessionPoolEthernet EnableVlan
 return nil
}

func(np *PppSessionPoolEtherne) DisableVlan () error {
 //parameters: SessionPoolHandle VlanTagNumber
 //AgtPppSessionPoolEthernet DisableVlan
 return nil
}

func(np *PppSessionPoolEtherne) SetEthertype () error {
 //parameters: SessionPoolHandle VlanTagNumber Ethertype
 //AgtPppSessionPoolEthernet SetEthertype
 return nil
}

func(np *PppSessionPoolEtherne) SetVlanId () error {
 //parameters: SessionPoolHandle VlanTagNumber StartingVlanId VlanIdIncrement
 //AgtPppSessionPoolEthernet SetVlanId
 return nil
}

func(np *PppSessionPoolEtherne) SetVlanIdCountOverride () error {
 //parameters: SessionPoolHandle VlanTagNumber IsCountOverride
 //AgtPppSessionPoolEthernet SetVlanIdCountOverride
 return nil
}

func(np *PppSessionPoolEtherne) GetVlanIdCountOverride ()(string, error) {
 //parameters: SessionPoolHandle VlanTagNumber
 //AgtPppSessionPoolEthernet GetVlanIdCountOverride
 return "", nil
}

func(np *PppSessionPoolEtherne) SetVlanIdRange () error {
 //parameters: SessionPoolHandle VlanTagNumber StartingVlanId VlanIdCount VlanIdIncrement VlanIdRepeat
 //AgtPppSessionPoolEthernet SetVlanIdRange
 return nil
}

func(np *PppSessionPoolEtherne) GetVlanIdRange ()(string, error) {
 //parameters: SessionPoolHandle VlanTagNumber
 //AgtPppSessionPoolEthernet GetVlanIdRange
 return "", nil
}

func(np *PppSessionPoolEtherne) SetServiceName () error {
 //parameters: SessionPoolHandle ServiceName
 //AgtPppSessionPoolEthernet SetServiceName
 return nil
}

func(np *PppSessionPoolEtherne) SetAccessConcentratorName () error {
 //parameters: SessionPoolHandle AccessConcentratorName
 //AgtPppSessionPoolEthernet SetAccessConcentratorName
 return nil
}

func(np *PppSessionPoolEtherne) GetAccessConcentratorName ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetAccessConcentratorName
 return "", nil
}

func(np *PppSessionPoolEtherne) EnableEchoVendorSpecificTag () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet EnableEchoVendorSpecificTag
 return nil
}

func(np *PppSessionPoolEtherne) DisableEchoVendorSpecificTag () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet DisableEchoVendorSpecificTag
 return nil
}

func(np *PppSessionPoolEtherne) IsEchoVendorSpecificTagEnabled () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet IsEchoVendorSpecificTagEnabled
 return nil
}

func(np *PppSessionPoolEtherne) SetStartingSourceMacAddress () error {
 //parameters: SessionPoolHandle MacAddress
 //AgtPppSessionPoolEthernet SetStartingSourceMacAddress
 return nil
}

func(np *PppSessionPoolEtherne) SetSourceMacAddressIncrement () error {
 //parameters: SessionPoolHandle Increment
 //AgtPppSessionPoolEthernet SetSourceMacAddressIncrement
 return nil
}

func(np *PppSessionPoolEtherne) GetServiceName ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetServiceName
 return "", nil
}

func(np *PppSessionPoolEtherne) GetNumberOfSessionsPerSrcMacAddr ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetNumberOfSessionsPerSrcMacAddr
 return "", nil
}

func(np *PppSessionPoolEtherne) GetStartingSourceMacAddress ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetStartingSourceMacAddress
 return "", nil
}

func(np *PppSessionPoolEtherne) GetSourceMacAddressIncrement ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetSourceMacAddressIncrement
 return "", nil
}

func(np *PppSessionPoolEtherne) GetInitialAckTimeout ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetInitialAckTimeout
 return "", nil
}

func(np *PppSessionPoolEtherne) GetMaxTxRetries ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetMaxTxRetries
 return "", nil
}

func(np *PppSessionPoolEtherne) GetVlanFlag ()(string, error) {
 //parameters: SessionPoolHandle VlanTagNumber
 //AgtPppSessionPoolEthernet GetVlanFlag
 return "", nil
}

func(np *PppSessionPoolEtherne) GetEthertype ()(string, error) {
 //parameters: SessionPoolHandle VlanTagNumber
 //AgtPppSessionPoolEthernet GetEthertype
 return "", nil
}

func(np *PppSessionPoolEtherne) GetVlanId ()(string, error) {
 //parameters: SessionPoolHandle VlanTagNumber
 //AgtPppSessionPoolEthernet GetVlanId
 return "", nil
}

func(np *PppSessionPoolEtherne) GetDestinationMacAddresses ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetDestinationMacAddresses
 return "", nil
}

func(np *PppSessionPoolEtherne) GetSessionIds ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetSessionIds
 return "", nil
}

func(np *PppSessionPoolEtherne) Enable () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet Enable
 return nil
}

func(np *PppSessionPoolEtherne) Disable () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet Disable
 return nil
}

func(np *PppSessionPoolEtherne) SetNumberOfSessions () error {
 //parameters: SessionPoolHandle NumSessions
 //AgtPppSessionPoolEthernet SetNumberOfSessions
 return nil
}

func(np *PppSessionPoolEtherne) SetSessionLifetime () error {
 //parameters: SessionPoolHandle MinLifetime MaxLifetime
 //AgtPppSessionPoolEthernet SetSessionLifetime
 return nil
}

func(np *PppSessionPoolEtherne) SetSessionInitiationRate () error {
 //parameters: SessionPoolHandle NumSessions Period
 //AgtPppSessionPoolEthernet SetSessionInitiationRate
 return nil
}

func(np *PppSessionPoolEtherne) SetLcpOption () error {
 //parameters: SessionPoolHandle LcpOption Value
 //AgtPppSessionPoolEthernet SetLcpOption
 return nil
}

func(np *PppSessionPoolEtherne) EnableAccm () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet EnableAccm
 return nil
}

func(np *PppSessionPoolEtherne) DisableAccm () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet DisableAccm
 return nil
}

func(np *PppSessionPoolEtherne) EnableAddressControlFieldCompression () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet EnableAddressControlFieldCompression
 return nil
}

func(np *PppSessionPoolEtherne) DisableAddressControlFieldCompression () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet DisableAddressControlFieldCompression
 return nil
}

func(np *PppSessionPoolEtherne) SetMode () error {
 //parameters: SessionPoolHandle Mode
 //AgtPppSessionPoolEthernet SetMode
 return nil
}

func(np *PppSessionPoolEtherne) SetAuthenticationProtocol () error {
 //parameters: SessionPoolHandle Protocol
 //AgtPppSessionPoolEthernet SetAuthenticationProtocol
 return nil
}

func(np *PppSessionPoolEtherne) SetAuthenticationParameter () error {
 //parameters: SessionPoolHandle Parameter Value
 //AgtPppSessionPoolEthernet SetAuthenticationParameter
 return nil
}

func(np *PppSessionPoolEtherne) EnableUniqueAuthenticationParameter () error {
 //parameters: SessionPoolHandle Parameter
 //AgtPppSessionPoolEthernet EnableUniqueAuthenticationParameter
 return nil
}

func(np *PppSessionPoolEtherne) DisableUniqueAuthenticationParameter () error {
 //parameters: SessionPoolHandle Parameter
 //AgtPppSessionPoolEthernet DisableUniqueAuthenticationParameter
 return nil
}

func(np *PppSessionPoolEtherne) SetAuthenticationParameterStartValue () error {
 //parameters: SessionPoolHandle Parameter StartValue
 //AgtPppSessionPoolEthernet SetAuthenticationParameterStartValue
 return nil
}

func(np *PppSessionPoolEtherne) GetAuthenticationParameter ()(string, error) {
 //parameters: SessionPoolHandle Parameter
 //AgtPppSessionPoolEthernet GetAuthenticationParameter
 return "", nil
}

func(np *PppSessionPoolEtherne) SetPeerAuthenticationParameter () error {
 //parameters: SessionPoolHandle PeerParameter Value
 //AgtPppSessionPoolEthernet SetPeerAuthenticationParameter
 return nil
}

func(np *PppSessionPoolEtherne) GetPeerAuthenticationParameter ()(string, error) {
 //parameters: SessionPoolHandle PeerParameter
 //AgtPppSessionPoolEthernet GetPeerAuthenticationParameter
 return "", nil
}

func(np *PppSessionPoolEtherne) EnableAcceptAllAuthenticationValues () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet EnableAcceptAllAuthenticationValues
 return nil
}

func(np *PppSessionPoolEtherne) DisableAcceptAllAuthenticationValues () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet DisableAcceptAllAuthenticationValues
 return nil
}

func(np *PppSessionPoolEtherne) IsAcceptAllAuthenticationValuesEnabled () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet IsAcceptAllAuthenticationValuesEnabled
 return nil
}

func(np *PppSessionPoolEtherne) SetUserToDomainRatio () error {
 //parameters: SessionPoolHandle UserToDomainRatio
 //AgtPppSessionPoolEthernet SetUserToDomainRatio
 return nil
}

func(np *PppSessionPoolEtherne) GetUserToDomainRatio ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetUserToDomainRatio
 return "", nil
}

func(np *PppSessionPoolEtherne) SetPasswordIncrementMethod () error {
 //parameters: SessionPoolHandle PasswordIncrementMethod
 //AgtPppSessionPoolEthernet SetPasswordIncrementMethod
 return nil
}

func(np *PppSessionPoolEtherne) GetPasswordIncrementMethod ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetPasswordIncrementMethod
 return "", nil
}

func(np *PppSessionPoolEtherne) SetDomainCountLimit () error {
 //parameters: SessionPoolHandle DomainCountLimit
 //AgtPppSessionPoolEthernet SetDomainCountLimit
 return nil
}

func(np *PppSessionPoolEtherne) GetDomainCountLimit ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetDomainCountLimit
 return "", nil
}

func(np *PppSessionPoolEtherne) EnableDomainNameIncrementAfterExtension () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet EnableDomainNameIncrementAfterExtension
 return nil
}

func(np *PppSessionPoolEtherne) DisableDomainNameIncrementAfterExtension () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet DisableDomainNameIncrementAfterExtension
 return nil
}

func(np *PppSessionPoolEtherne) IsDomainNameIncrementAfterExtensionEnabled () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet IsDomainNameIncrementAfterExtensionEnabled
 return nil
}

func(np *PppSessionPoolEtherne) SetAcceptedAuthenticationParameter () error {
 //parameters: SessionPoolHandle AcceptedParameter Value
 //AgtPppSessionPoolEthernet SetAcceptedAuthenticationParameter
 return nil
}

func(np *PppSessionPoolEtherne) GetAcceptedAuthenticationParameter ()(string, error) {
 //parameters: SessionPoolHandle AcceptedParameter
 //AgtPppSessionPoolEthernet GetAcceptedAuthenticationParameter
 return "", nil
}

func(np *PppSessionPoolEtherne) SetStartingSourceIpAddress () error {
 //parameters: SessionPoolHandle IpAddress
 //AgtPppSessionPoolEthernet SetStartingSourceIpAddress
 return nil
}

func(np *PppSessionPoolEtherne) SetSourceIpAddressModifier () error {
 //parameters: SessionPoolHandle Modifier
 //AgtPppSessionPoolEthernet SetSourceIpAddressModifier
 return nil
}

func(np *PppSessionPoolEtherne) SetMaxConnectionAttempts () error {
 //parameters: SessionPoolHandle MaxConnectionAttempts
 //AgtPppSessionPoolEthernet SetMaxConnectionAttempts
 return nil
}

func(np *PppSessionPoolEtherne) GetMaxConnectionAttempts ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetMaxConnectionAttempts
 return "", nil
}

func(np *PppSessionPoolEtherne) RetryConnections () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet RetryConnections
 return nil
}

func(np *PppSessionPoolEtherne) GetConnectionAttempts ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetConnectionAttempts
 return "", nil
}

func(np *PppSessionPoolEtherne) Cancel () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet Cancel
 return nil
}

func(np *PppSessionPoolEtherne) GetAccmFlag ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetAccmFlag
 return "", nil
}

func(np *PppSessionPoolEtherne) GetAddressControlFieldCompressionFlag ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetAddressControlFieldCompressionFlag
 return "", nil
}

func(np *PppSessionPoolEtherne) GetAuthenticationParameterStartValue ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetAuthenticationParameterStartValue
 return "", nil
}

func(np *PppSessionPoolEtherne) GetAuthenticationProtocol ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetAuthenticationProtocol
 return "", nil
}

func(np *PppSessionPoolEtherne) GetEnableFlag ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetEnableFlag
 return "", nil
}

func(np *PppSessionPoolEtherne) GetLcpOption ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetLcpOption
 return "", nil
}

func(np *PppSessionPoolEtherne) GetMode ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetMode
 return "", nil
}

func(np *PppSessionPoolEtherne) GetName ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetName
 return "", nil
}

func(np *PppSessionPoolEtherne) SetName () error {
 //parameters: SessionPoolHandle Name
 //AgtPppSessionPoolEthernet SetName
 return nil
}

func(np *PppSessionPoolEtherne) GetNegotiatedIpAddresses ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetNegotiatedIpAddresses
 return "", nil
}

func(np *PppSessionPoolEtherne) GetNumberOfEstablishedSessions ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetNumberOfEstablishedSessions
 return "", nil
}

func(np *PppSessionPoolEtherne) GetNumberOfSessions ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetNumberOfSessions
 return "", nil
}

func(np *PppSessionPoolEtherne) GetRemoteIpAddresses ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetRemoteIpAddresses
 return "", nil
}

func(np *PppSessionPoolEtherne) GetSessionInitiationRate ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetSessionInitiationRate
 return "", nil
}

func(np *PppSessionPoolEtherne) GetSessionLifetime ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetSessionLifetime
 return "", nil
}

func(np *PppSessionPoolEtherne) GetSessionStateDescriptions ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetSessionStateDescriptions
 return "", nil
}

func(np *PppSessionPoolEtherne) GetSourceIpAddressModifier ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetSourceIpAddressModifier
 return "", nil
}

func(np *PppSessionPoolEtherne) GetStartingSourceIpAddress ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetStartingSourceIpAddress
 return "", nil
}

func(np *PppSessionPoolEtherne) GetState ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolEthernet GetState
 return "", nil
}

func(np *PppSessionPoolEtherne) GetStateSummary ()(string, error) {
 //parameters: 
 //AgtPppSessionPoolEthernet GetStateSummary
 return "", nil
}

func(np *PppSessionPoolEtherne) GetUniqueAuthenticationParameterFlag ()(string, error) {
 //parameters: SessionPoolHandle Parameter
 //AgtPppSessionPoolEthernet GetUniqueAuthenticationParameterFlag
 return "", nil
}

func(np *PppSessionPoolEtherne) EnableAllSessions () error {
 //parameters: 
 //AgtPppSessionPoolEthernet EnableAllSessions
 return nil
}

func(np *PppSessionPoolEtherne) DisableAllSessions () error {
 //parameters: 
 //AgtPppSessionPoolEthernet DisableAllSessions
 return nil
}

