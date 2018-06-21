package n2xsdk

type ConnectionGroupPdu struct {
 Handler string
}

func(np *ConnectionGroupPdu) ListAvailableFields ()(string, error) {
 //parameters: 
 //AgtConnectionGroupPdu ListAvailableFields
 return "", nil
}

func(np *ConnectionGroupPdu) ListOptionalFields ()(string, error) {
 //parameters: PduHandle
 //AgtConnectionGroupPdu ListOptionalFields
 return "", nil
}

func(np *ConnectionGroupPdu) EnableOptionalField () error {
 //parameters: PduHandle OptionalPduField
 //AgtConnectionGroupPdu EnableOptionalField
 return nil
}

func(np *ConnectionGroupPdu) DisableOptionalField () error {
 //parameters: PduHandle OptionalPduField
 //AgtConnectionGroupPdu DisableOptionalField
 return nil
}

func(np *ConnectionGroupPdu) IsOptionalFieldEnabled () error {
 //parameters: PduHandle OptionalPduField
 //AgtConnectionGroupPdu IsOptionalFieldEnabled
 return nil
}

func(np *ConnectionGroupPdu) ListEditableFields ()(string, error) {
 //parameters: PduHandle
 //AgtConnectionGroupPdu ListEditableFields
 return "", nil
}

func(np *ConnectionGroupPdu) GetFieldValueType ()(string, error) {
 //parameters: PduHandle PduField
 //AgtConnectionGroupPdu GetFieldValueType
 return "", nil
}

func(np *ConnectionGroupPdu) SetValue () error {
 //parameters: PduHandle PduField Value
 //AgtConnectionGroupPdu SetValue
 return nil
}

func(np *ConnectionGroupPdu) GetValue ()(string, error) {
 //parameters: PduHandle PduField
 //AgtConnectionGroupPdu GetValue
 return "", nil
}

func(np *ConnectionGroupPdu) SetValueRange () error {
 //parameters: PduHandle PduField StartValue Count Step
 //AgtConnectionGroupPdu SetValueRange
 return nil
}

func(np *ConnectionGroupPdu) GetValueRange ()(string, error) {
 //parameters: PduHandle PduField
 //AgtConnectionGroupPdu GetValueRange
 return "", nil
}

func(np *ConnectionGroupPdu) SetValueList ()(string, error) {
 //parameters: PduHandle PduField NumberOfValues psaValueList
 //AgtConnectionGroupPdu SetValueList
 return "", nil
}

func(np *ConnectionGroupPdu) GetValueList ()(string, error) {
 //parameters: PduHandle PduField
 //AgtConnectionGroupPdu GetValueList
 return "", nil
}

func(np *ConnectionGroupPdu) SetAddress () error {
 //parameters: PduHandle PduField Address
 //AgtConnectionGroupPdu SetAddress
 return nil
}

func(np *ConnectionGroupPdu) GetAddress ()(string, error) {
 //parameters: PduHandle PduField
 //AgtConnectionGroupPdu GetAddress
 return "", nil
}

func(np *ConnectionGroupPdu) SetAddressRange () error {
 //parameters: PduHandle PduField StartAddress Count IpPrefixOrMacOffset Step
 //AgtConnectionGroupPdu SetAddressRange
 return nil
}

func(np *ConnectionGroupPdu) GetAddressRange ()(string, error) {
 //parameters: PduHandle PduField
 //AgtConnectionGroupPdu GetAddressRange
 return "", nil
}

func(np *ConnectionGroupPdu) SetAddressList ()(string, error) {
 //parameters: PduHandle PduField NumberOfAddresses psaAddressList
 //AgtConnectionGroupPdu SetAddressList
 return "", nil
}

func(np *ConnectionGroupPdu) GetAddressList ()(string, error) {
 //parameters: PduHandle PduField
 //AgtConnectionGroupPdu GetAddressList
 return "", nil
}

