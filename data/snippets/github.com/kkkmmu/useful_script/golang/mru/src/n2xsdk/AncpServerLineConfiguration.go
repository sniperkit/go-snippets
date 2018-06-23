package n2xsdk

type AncpServerLineConfiguration struct {
 Handler string
}

func(np *ncpServerLineConfiguration) DisableAutoLineConfigPortManagement () error {
 //parameters: SessionPoolHandle
 //AgtAncpServerLineConfiguration DisableAutoLineConfigPortManagement
 return nil
}

func(np *ncpServerLineConfiguration) EnableAutoLineConfigPortManagement () error {
 //parameters: SessionPoolHandle
 //AgtAncpServerLineConfiguration EnableAutoLineConfigPortManagement
 return nil
}

func(np *ncpServerLineConfiguration) IsAutoLineConfigPortManagementEnabled () error {
 //parameters: SessionPoolHandle
 //AgtAncpServerLineConfiguration IsAutoLineConfigPortManagementEnabled
 return nil
}

func(np *ncpServerLineConfiguration) AddProfileEntry () error {
 //parameters: SessionPoolHandle CircuitId ProfileName
 //AgtAncpServerLineConfiguration AddProfileEntry
 return nil
}

func(np *ncpServerLineConfiguration) ListProfileEntries ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtAncpServerLineConfiguration ListProfileEntries
 return "", nil
}

func(np *ncpServerLineConfiguration) RemoveProfileEntry () error {
 //parameters: SessionPoolHandle CircuitId
 //AgtAncpServerLineConfiguration RemoveProfileEntry
 return nil
}

func(np *ncpServerLineConfiguration) RemoveAllProfileEntries () error {
 //parameters: SessionPoolHandle
 //AgtAncpServerLineConfiguration RemoveAllProfileEntries
 return nil
}

func(np *ncpServerLineConfiguration) GetProfileName ()(string, error) {
 //parameters: SessionPoolHandle CircuitId
 //AgtAncpServerLineConfiguration GetProfileName
 return "", nil
}

func(np *ncpServerLineConfiguration) SetDefaultProfileName () error {
 //parameters: SessionPoolHandle ProfileName
 //AgtAncpServerLineConfiguration SetDefaultProfileName
 return nil
}

func(np *ncpServerLineConfiguration) GetDefaultProfileName ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtAncpServerLineConfiguration GetDefaultProfileName
 return "", nil
}

func(np *ncpServerLineConfiguration) SetResultField () error {
 //parameters: SessionPoolHandle AncpResult
 //AgtAncpServerLineConfiguration SetResultField
 return nil
}

func(np *ncpServerLineConfiguration) GetResultField ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtAncpServerLineConfiguration GetResultField
 return "", nil
}

func(np *ncpServerLineConfiguration) AddTlv () error {
 //parameters: SessionPoolHandle TlvType TlvLength TlvValue
 //AgtAncpServerLineConfiguration AddTlv
 return nil
}

func(np *ncpServerLineConfiguration) ListTlvs ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtAncpServerLineConfiguration ListTlvs
 return "", nil
}

func(np *ncpServerLineConfiguration) RemoveTlv () error {
 //parameters: SessionPoolHandle TlvType
 //AgtAncpServerLineConfiguration RemoveTlv
 return nil
}

func(np *ncpServerLineConfiguration) RemoveAllTlvs () error {
 //parameters: SessionPoolHandle
 //AgtAncpServerLineConfiguration RemoveAllTlvs
 return nil
}

func(np *ncpServerLineConfiguration) SetTlv () error {
 //parameters: SessionPoolHandle TlvType TlvLength TlvValue
 //AgtAncpServerLineConfiguration SetTlv
 return nil
}

func(np *ncpServerLineConfiguration) GetTlv ()(string, error) {
 //parameters: SessionPoolHandle TlvType
 //AgtAncpServerLineConfiguration GetTlv
 return "", nil
}

func(np *ncpServerLineConfiguration) AddSubTlv () error {
 //parameters: SessionPoolHandle TlvType SubTlvType SubTlvLength SubTlvValue
 //AgtAncpServerLineConfiguration AddSubTlv
 return nil
}

func(np *ncpServerLineConfiguration) ListSubTlvs ()(string, error) {
 //parameters: SessionPoolHandle TlvType
 //AgtAncpServerLineConfiguration ListSubTlvs
 return "", nil
}

func(np *ncpServerLineConfiguration) RemoveSubTlv () error {
 //parameters: SessionPoolHandle TlvType SubTlvType
 //AgtAncpServerLineConfiguration RemoveSubTlv
 return nil
}

func(np *ncpServerLineConfiguration) RemoveAllSubTlvs () error {
 //parameters: SessionPoolHandle TlvType
 //AgtAncpServerLineConfiguration RemoveAllSubTlvs
 return nil
}

func(np *ncpServerLineConfiguration) SetSubTlv () error {
 //parameters: SessionPoolHandle TlvType SubTlvType SubTlvLength SubTlvValue
 //AgtAncpServerLineConfiguration SetSubTlv
 return nil
}

func(np *ncpServerLineConfiguration) GetSubTlv ()(string, error) {
 //parameters: SessionPoolHandle TlvType SubTlvType
 //AgtAncpServerLineConfiguration GetSubTlv
 return "", nil
}

