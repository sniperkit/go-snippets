package n2xsdk

type MsdpPdu struct {
 Handler string
}

func(np *MsdpPdu) BuildPdu () error {
 //parameters: SessionHandle PduType
 //AgtMsdpPdu BuildPdu
 return nil
}

func(np *MsdpPdu) BuildPduUsingId () error {
 //parameters: PduType ReferencePduId
 //AgtMsdpPdu BuildPduUsingId
 return nil
}

func(np *MsdpPdu) BuildPduUsingContents () error {
 //parameters: SessionHandle PduContents
 //AgtMsdpPdu BuildPduUsingContents
 return nil
}

func(np *MsdpPdu) DeletePdu () error {
 //parameters: PduId
 //AgtMsdpPdu DeletePdu
 return nil
}

func(np *MsdpPdu) DeleteAllPdusInSession () error {
 //parameters: SessionHandle
 //AgtMsdpPdu DeleteAllPdusInSession
 return nil
}

func(np *MsdpPdu) DeleteAllPdus () error {
 //parameters: 
 //AgtMsdpPdu DeleteAllPdus
 return nil
}

func(np *MsdpPdu) DecodePdu () error {
 //parameters: FormatMode PduId
 //AgtMsdpPdu DecodePdu
 return nil
}

func(np *MsdpPdu) GetFieldCount ()(string, error) {
 //parameters: PduId FieldType
 //AgtMsdpPdu GetFieldCount
 return "", nil
}

func(np *MsdpPdu) GetSubfield ()(string, error) {
 //parameters: PduId FieldType Index SubfieldType
 //AgtMsdpPdu GetSubfield
 return "", nil
}

func(np *MsdpPdu) SetSubfield () error {
 //parameters: PduId FieldType Index SubfieldType SubfieldValue
 //AgtMsdpPdu SetSubfield
 return nil
}

func(np *MsdpPdu) DeleteField () error {
 //parameters: PduId FieldType Index
 //AgtMsdpPdu DeleteField
 return nil
}

func(np *MsdpPdu) SendPdu () error {
 //parameters: DestinationIpAddress PduId
 //AgtMsdpPdu SendPdu
 return nil
}

func(np *MsdpPdu) GetNextPduInfo ()(string, error) {
 //parameters: TimeOutSec Count psaSessionHandle
 //AgtMsdpPdu GetNextPduInfo
 return "", nil
}

func(np *MsdpPdu) GetPduDetails ()(string, error) {
 //parameters: PduId PduInfoType
 //AgtMsdpPdu GetPduDetails
 return "", nil
}

func(np *MsdpPdu) Prohibit () error {
 //parameters: SessionHandle PduType
 //AgtMsdpPdu Prohibit
 return nil
}

func(np *MsdpPdu) Permit () error {
 //parameters: SessionHandle PduType
 //AgtMsdpPdu Permit
 return nil
}

func(np *MsdpPdu) IsProhibited () error {
 //parameters: SessionHandle PduType
 //AgtMsdpPdu IsProhibited
 return nil
}

func(np *MsdpPdu) AddPduFilter () error {
 //parameters: SessionHandle PduType ModifierType SourceIpAddress GroupIpAddress
 //AgtMsdpPdu AddPduFilter
 return nil
}

func(np *MsdpPdu) RemovePduFilter () error {
 //parameters: SessionHandle FilterId
 //AgtMsdpPdu RemovePduFilter
 return nil
}

func(np *MsdpPdu) RemoveAllPduFilters () error {
 //parameters: SessionHandle
 //AgtMsdpPdu RemoveAllPduFilters
 return nil
}

func(np *MsdpPdu) ListPduFilters ()(string, error) {
 //parameters: SessionHandle
 //AgtMsdpPdu ListPduFilters
 return "", nil
}

func(np *MsdpPdu) GetPduFilterDetails ()(string, error) {
 //parameters: SessionHandle FilterId
 //AgtMsdpPdu GetPduFilterDetails
 return "", nil
}

func(np *MsdpPdu) StartPduSave () error {
 //parameters: SessionHandle
 //AgtMsdpPdu StartPduSave
 return nil
}

func(np *MsdpPdu) StopPduSave () error {
 //parameters: SessionHandle
 //AgtMsdpPdu StopPduSave
 return nil
}

func(np *MsdpPdu) IsPduSaveEnabled () error {
 //parameters: SessionHandle
 //AgtMsdpPdu IsPduSaveEnabled
 return nil
}

func(np *MsdpPdu) EnablePduNotification () error {
 //parameters: SessionHandle
 //AgtMsdpPdu EnablePduNotification
 return nil
}

func(np *MsdpPdu) DisablePduNotification () error {
 //parameters: SessionHandle
 //AgtMsdpPdu DisablePduNotification
 return nil
}

func(np *MsdpPdu) IsPduNotificationEnabled () error {
 //parameters: SessionHandle
 //AgtMsdpPdu IsPduNotificationEnabled
 return nil
}

func(np *MsdpPdu) ListSavedPdus ()(string, error) {
 //parameters: SessionHandle
 //AgtMsdpPdu ListSavedPdus
 return "", nil
}

