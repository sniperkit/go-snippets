package n2xsdk

type HttpIpv6ClientL4RedirectSettings struct {
 Handler string
}

func(np *HttpIpv6ClientL4RedirectSettings) EnableL4RedirectionSettings () error {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientL4RedirectSettings EnableL4RedirectionSettings
 return nil
}

func(np *HttpIpv6ClientL4RedirectSettings) DisableL4RedirectionSettings () error {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientL4RedirectSettings DisableL4RedirectionSettings
 return nil
}

func(np *HttpIpv6ClientL4RedirectSettings) IsL4RedirectionSettingsEnabled () error {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientL4RedirectSettings IsL4RedirectionSettingsEnabled
 return nil
}

func(np *HttpIpv6ClientL4RedirectSettings) SetAuthenticationLevel () error {
 //parameters: SessionPoolHandle AuthenticationLevel
 //AgtHttpIpv6ClientL4RedirectSettings SetAuthenticationLevel
 return nil
}

func(np *HttpIpv6ClientL4RedirectSettings) GetAuthenticationLevel ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientL4RedirectSettings GetAuthenticationLevel
 return "", nil
}

func(np *HttpIpv6ClientL4RedirectSettings) SetFirstCredentialNameValue () error {
 //parameters: SessionPoolHandle CredentialFieldName CredentialFieldValue
 //AgtHttpIpv6ClientL4RedirectSettings SetFirstCredentialNameValue
 return nil
}

func(np *HttpIpv6ClientL4RedirectSettings) GetFirstCredentialNameValue ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientL4RedirectSettings GetFirstCredentialNameValue
 return "", nil
}

func(np *HttpIpv6ClientL4RedirectSettings) GetFirstCredentialFieldValue ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientL4RedirectSettings GetFirstCredentialFieldValue
 return "", nil
}

func(np *HttpIpv6ClientL4RedirectSettings) SetFirstCredentialValueModifierIncrementor () error {
 //parameters: SessionPoolHandle FieldValueModifier IncrementingRange Repeat
 //AgtHttpIpv6ClientL4RedirectSettings SetFirstCredentialValueModifierIncrementor
 return nil
}

func(np *HttpIpv6ClientL4RedirectSettings) GetFirstCredentialValueModifierIncrementor ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientL4RedirectSettings GetFirstCredentialValueModifierIncrementor
 return "", nil
}

func(np *HttpIpv6ClientL4RedirectSettings) SetSecondCredentialNameValue () error {
 //parameters: SessionPoolHandle CredentialFieldName CredentialFieldValue
 //AgtHttpIpv6ClientL4RedirectSettings SetSecondCredentialNameValue
 return nil
}

func(np *HttpIpv6ClientL4RedirectSettings) GetSecondCredentialNameValue ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientL4RedirectSettings GetSecondCredentialNameValue
 return "", nil
}

func(np *HttpIpv6ClientL4RedirectSettings) GetSecondCredentialFieldValue ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientL4RedirectSettings GetSecondCredentialFieldValue
 return "", nil
}

func(np *HttpIpv6ClientL4RedirectSettings) SetSecondCredentialValueModifierIncrementor () error {
 //parameters: SessionPoolHandle FieldValueModifier IncrementingRange Repeat
 //AgtHttpIpv6ClientL4RedirectSettings SetSecondCredentialValueModifierIncrementor
 return nil
}

func(np *HttpIpv6ClientL4RedirectSettings) GetSecondCredentialValueModifierIncrementor ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientL4RedirectSettings GetSecondCredentialValueModifierIncrementor
 return "", nil
}

func(np *HttpIpv6ClientL4RedirectSettings) AddAuthenticationParameter () error {
 //parameters: SessionPoolHandle RowIndex
 //AgtHttpIpv6ClientL4RedirectSettings AddAuthenticationParameter
 return nil
}

func(np *HttpIpv6ClientL4RedirectSettings) GetNumAuthenticationParameters ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientL4RedirectSettings GetNumAuthenticationParameters
 return "", nil
}

func(np *HttpIpv6ClientL4RedirectSettings) ListAuthenticationParameters ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientL4RedirectSettings ListAuthenticationParameters
 return "", nil
}

func(np *HttpIpv6ClientL4RedirectSettings) RemoveAuthenticationParameter () error {
 //parameters: SessionPoolHandle RowIndex
 //AgtHttpIpv6ClientL4RedirectSettings RemoveAuthenticationParameter
 return nil
}

func(np *HttpIpv6ClientL4RedirectSettings) RemoveAllAuthenticationParameters () error {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientL4RedirectSettings RemoveAllAuthenticationParameters
 return nil
}

func(np *HttpIpv6ClientL4RedirectSettings) SetAuthenticationParameterValue () error {
 //parameters: SessionPoolHandle RowIndex AuthenticationParameterValue
 //AgtHttpIpv6ClientL4RedirectSettings SetAuthenticationParameterValue
 return nil
}

func(np *HttpIpv6ClientL4RedirectSettings) GetAuthenticationParameterValue ()(string, error) {
 //parameters: SessionPoolHandle RowIndex
 //AgtHttpIpv6ClientL4RedirectSettings GetAuthenticationParameterValue
 return "", nil
}

func(np *HttpIpv6ClientL4RedirectSettings) SetSuccessIndicatorType () error {
 //parameters: SessionPoolHandle SuccessIndicatorType
 //AgtHttpIpv6ClientL4RedirectSettings SetSuccessIndicatorType
 return nil
}

func(np *HttpIpv6ClientL4RedirectSettings) GetSuccessIndicatorType ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientL4RedirectSettings GetSuccessIndicatorType
 return "", nil
}

func(np *HttpIpv6ClientL4RedirectSettings) SetSuccessIndicatorValue () error {
 //parameters: SessionPoolHandle SuccessIndicatorValue
 //AgtHttpIpv6ClientL4RedirectSettings SetSuccessIndicatorValue
 return nil
}

func(np *HttpIpv6ClientL4RedirectSettings) GetSuccessIndicatorValue ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientL4RedirectSettings GetSuccessIndicatorValue
 return "", nil
}

