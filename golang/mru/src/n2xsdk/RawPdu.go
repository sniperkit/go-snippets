package n2xsdk

type RawPdu struct {
 Handler string
}

func(np *RawPdu) GetAllPduBytes ()(string, error) {
 //parameters: PduHandle
 //AgtRawPdu GetAllPduBytes
 return "", nil
}

func(np *RawPdu) GetPduBytes ()(string, error) {
 //parameters: PduHandle Offset Length
 //AgtRawPdu GetPduBytes
 return "", nil
}

func(np *RawPdu) SetAllPduBytes () error {
 //parameters: PduHandle PduBytes
 //AgtRawPdu SetAllPduBytes
 return nil
}

func(np *RawPdu) SetPduBytes () error {
 //parameters: PduHandle Offset Length PduBytes
 //AgtRawPdu SetPduBytes
 return nil
}

func(np *RawPdu) InsertPduBytes () error {
 //parameters: PduHandle Offset PduBytes
 //AgtRawPdu InsertPduBytes
 return nil
}

func(np *RawPdu) UnsetAllPduBytes () error {
 //parameters: PduHandle
 //AgtRawPdu UnsetAllPduBytes
 return nil
}

