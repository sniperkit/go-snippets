package n2xsdk

type PppSessionPoolAtm struct {
 Handler string
}

func(np *PppSessionPoolAtm) SetAtmEncapsulation () error {
 //parameters: SessionPoolHandle Encapsulation
 //AgtPppSessionPoolAtm SetAtmEncapsulation
 return nil
}

func(np *PppSessionPoolAtm) SetStartingVpiVci () error {
 //parameters: SessionPoolHandle VpiStart VciStart
 //AgtPppSessionPoolAtm SetStartingVpiVci
 return nil
}

func(np *PppSessionPoolAtm) SetVpiVciIncrement () error {
 //parameters: SessionPoolHandle VpiIncrement VciIncrement
 //AgtPppSessionPoolAtm SetVpiVciIncrement
 return nil
}

func(np *PppSessionPoolAtm) GetAtmEncapsulation ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetAtmEncapsulation
 return "", nil
}

func(np *PppSessionPoolAtm) GetStartingVpiVci ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetStartingVpiVci
 return "", nil
}

func(np *PppSessionPoolAtm) GetVpiVciIncrement ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetVpiVciIncrement
 return "", nil
}

func(np *PppSessionPoolAtm) Enable () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm Enable
 return nil
}

func(np *PppSessionPoolAtm) Disable () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm Disable
 return nil
}

func(np *PppSessionPoolAtm) SetNumberOfSessions () error {
 //parameters: SessionPoolHandle NumSessions
 //AgtPppSessionPoolAtm SetNumberOfSessions
 return nil
}

func(np *PppSessionPoolAtm) SetSessionLifetime () error {
 //parameters: SessionPoolHandle MinLifetime MaxLifetime
 //AgtPppSessionPoolAtm SetSessionLifetime
 return nil
}

func(np *PppSessionPoolAtm) SetSessionInitiationRate () error {
 //parameters: SessionPoolHandle NumSessions Period
 //AgtPppSessionPoolAtm SetSessionInitiationRate
 return nil
}

func(np *PppSessionPoolAtm) SetLcpOption () error {
 //parameters: SessionPoolHandle LcpOption Value
 //AgtPppSessionPoolAtm SetLcpOption
 return nil
}

func(np *PppSessionPoolAtm) EnableAccm () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm EnableAccm
 return nil
}

func(np *PppSessionPoolAtm) DisableAccm () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm DisableAccm
 return nil
}

func(np *PppSessionPoolAtm) EnableAddressControlFieldCompression () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm EnableAddressControlFieldCompression
 return nil
}

func(np *PppSessionPoolAtm) DisableAddressControlFieldCompression () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm DisableAddressControlFieldCompression
 return nil
}

func(np *PppSessionPoolAtm) SetMode () error {
 //parameters: SessionPoolHandle Mode
 //AgtPppSessionPoolAtm SetMode
 return nil
}

func(np *PppSessionPoolAtm) SetAuthenticationProtocol () error {
 //parameters: SessionPoolHandle Protocol
 //AgtPppSessionPoolAtm SetAuthenticationProtocol
 return nil
}

func(np *PppSessionPoolAtm) SetAuthenticationParameter () error {
 //parameters: SessionPoolHandle Parameter Value
 //AgtPppSessionPoolAtm SetAuthenticationParameter
 return nil
}

func(np *PppSessionPoolAtm) EnableUniqueAuthenticationParameter () error {
 //parameters: SessionPoolHandle Parameter
 //AgtPppSessionPoolAtm EnableUniqueAuthenticationParameter
 return nil
}

func(np *PppSessionPoolAtm) DisableUniqueAuthenticationParameter () error {
 //parameters: SessionPoolHandle Parameter
 //AgtPppSessionPoolAtm DisableUniqueAuthenticationParameter
 return nil
}

func(np *PppSessionPoolAtm) SetAuthenticationParameterStartValue () error {
 //parameters: SessionPoolHandle Parameter StartValue
 //AgtPppSessionPoolAtm SetAuthenticationParameterStartValue
 return nil
}

func(np *PppSessionPoolAtm) GetAuthenticationParameter ()(string, error) {
 //parameters: SessionPoolHandle Parameter
 //AgtPppSessionPoolAtm GetAuthenticationParameter
 return "", nil
}

func(np *PppSessionPoolAtm) SetPeerAuthenticationParameter () error {
 //parameters: SessionPoolHandle PeerParameter Value
 //AgtPppSessionPoolAtm SetPeerAuthenticationParameter
 return nil
}

func(np *PppSessionPoolAtm) GetPeerAuthenticationParameter ()(string, error) {
 //parameters: SessionPoolHandle PeerParameter
 //AgtPppSessionPoolAtm GetPeerAuthenticationParameter
 return "", nil
}

func(np *PppSessionPoolAtm) EnableAcceptAllAuthenticationValues () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm EnableAcceptAllAuthenticationValues
 return nil
}

func(np *PppSessionPoolAtm) DisableAcceptAllAuthenticationValues () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm DisableAcceptAllAuthenticationValues
 return nil
}

func(np *PppSessionPoolAtm) IsAcceptAllAuthenticationValuesEnabled () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm IsAcceptAllAuthenticationValuesEnabled
 return nil
}

func(np *PppSessionPoolAtm) SetUserToDomainRatio () error {
 //parameters: SessionPoolHandle UserToDomainRatio
 //AgtPppSessionPoolAtm SetUserToDomainRatio
 return nil
}

func(np *PppSessionPoolAtm) GetUserToDomainRatio ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetUserToDomainRatio
 return "", nil
}

func(np *PppSessionPoolAtm) SetPasswordIncrementMethod () error {
 //parameters: SessionPoolHandle PasswordIncrementMethod
 //AgtPppSessionPoolAtm SetPasswordIncrementMethod
 return nil
}

func(np *PppSessionPoolAtm) GetPasswordIncrementMethod ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetPasswordIncrementMethod
 return "", nil
}

func(np *PppSessionPoolAtm) SetDomainCountLimit () error {
 //parameters: SessionPoolHandle DomainCountLimit
 //AgtPppSessionPoolAtm SetDomainCountLimit
 return nil
}

func(np *PppSessionPoolAtm) GetDomainCountLimit ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetDomainCountLimit
 return "", nil
}

func(np *PppSessionPoolAtm) EnableDomainNameIncrementAfterExtension () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm EnableDomainNameIncrementAfterExtension
 return nil
}

func(np *PppSessionPoolAtm) DisableDomainNameIncrementAfterExtension () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm DisableDomainNameIncrementAfterExtension
 return nil
}

func(np *PppSessionPoolAtm) IsDomainNameIncrementAfterExtensionEnabled () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm IsDomainNameIncrementAfterExtensionEnabled
 return nil
}

func(np *PppSessionPoolAtm) SetAcceptedAuthenticationParameter () error {
 //parameters: SessionPoolHandle AcceptedParameter Value
 //AgtPppSessionPoolAtm SetAcceptedAuthenticationParameter
 return nil
}

func(np *PppSessionPoolAtm) GetAcceptedAuthenticationParameter ()(string, error) {
 //parameters: SessionPoolHandle AcceptedParameter
 //AgtPppSessionPoolAtm GetAcceptedAuthenticationParameter
 return "", nil
}

func(np *PppSessionPoolAtm) SetStartingSourceIpAddress () error {
 //parameters: SessionPoolHandle IpAddress
 //AgtPppSessionPoolAtm SetStartingSourceIpAddress
 return nil
}

func(np *PppSessionPoolAtm) SetSourceIpAddressModifier () error {
 //parameters: SessionPoolHandle Modifier
 //AgtPppSessionPoolAtm SetSourceIpAddressModifier
 return nil
}

func(np *PppSessionPoolAtm) SetMaxConnectionAttempts () error {
 //parameters: SessionPoolHandle MaxConnectionAttempts
 //AgtPppSessionPoolAtm SetMaxConnectionAttempts
 return nil
}

func(np *PppSessionPoolAtm) GetMaxConnectionAttempts ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetMaxConnectionAttempts
 return "", nil
}

func(np *PppSessionPoolAtm) RetryConnections () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm RetryConnections
 return nil
}

func(np *PppSessionPoolAtm) GetConnectionAttempts ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetConnectionAttempts
 return "", nil
}

func(np *PppSessionPoolAtm) Cancel () error {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm Cancel
 return nil
}

func(np *PppSessionPoolAtm) GetAccmFlag ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetAccmFlag
 return "", nil
}

func(np *PppSessionPoolAtm) GetAddressControlFieldCompressionFlag ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetAddressControlFieldCompressionFlag
 return "", nil
}

func(np *PppSessionPoolAtm) GetAuthenticationParameterStartValue ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetAuthenticationParameterStartValue
 return "", nil
}

func(np *PppSessionPoolAtm) GetAuthenticationProtocol ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetAuthenticationProtocol
 return "", nil
}

func(np *PppSessionPoolAtm) GetEnableFlag ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetEnableFlag
 return "", nil
}

func(np *PppSessionPoolAtm) GetLcpOption ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetLcpOption
 return "", nil
}

func(np *PppSessionPoolAtm) GetMode ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetMode
 return "", nil
}

func(np *PppSessionPoolAtm) GetName ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetName
 return "", nil
}

func(np *PppSessionPoolAtm) SetName () error {
 //parameters: SessionPoolHandle Name
 //AgtPppSessionPoolAtm SetName
 return nil
}

func(np *PppSessionPoolAtm) GetNegotiatedIpAddresses ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetNegotiatedIpAddresses
 return "", nil
}

func(np *PppSessionPoolAtm) GetNumberOfEstablishedSessions ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetNumberOfEstablishedSessions
 return "", nil
}

func(np *PppSessionPoolAtm) GetNumberOfSessions ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetNumberOfSessions
 return "", nil
}

func(np *PppSessionPoolAtm) GetRemoteIpAddresses ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetRemoteIpAddresses
 return "", nil
}

func(np *PppSessionPoolAtm) GetSessionInitiationRate ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetSessionInitiationRate
 return "", nil
}

func(np *PppSessionPoolAtm) GetSessionLifetime ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetSessionLifetime
 return "", nil
}

func(np *PppSessionPoolAtm) GetSessionStateDescriptions ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetSessionStateDescriptions
 return "", nil
}

func(np *PppSessionPoolAtm) GetSourceIpAddressModifier ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetSourceIpAddressModifier
 return "", nil
}

func(np *PppSessionPoolAtm) GetStartingSourceIpAddress ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetStartingSourceIpAddress
 return "", nil
}

func(np *PppSessionPoolAtm) GetState ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtPppSessionPoolAtm GetState
 return "", nil
}

func(np *PppSessionPoolAtm) GetStateSummary ()(string, error) {
 //parameters: 
 //AgtPppSessionPoolAtm GetStateSummary
 return "", nil
}

func(np *PppSessionPoolAtm) GetUniqueAuthenticationParameterFlag ()(string, error) {
 //parameters: SessionPoolHandle Parameter
 //AgtPppSessionPoolAtm GetUniqueAuthenticationParameterFlag
 return "", nil
}

func(np *PppSessionPoolAtm) EnableAllSessions () error {
 //parameters: 
 //AgtPppSessionPoolAtm EnableAllSessions
 return nil
}

func(np *PppSessionPoolAtm) DisableAllSessions () error {
 //parameters: 
 //AgtPppSessionPoolAtm DisableAllSessions
 return nil
}

