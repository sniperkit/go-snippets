package n2xsdk

type PduPayload struct {
 Handler string
}

func(np *PduPayload) GetPayloadBytes ()(string, error) {
 //parameters: PduHandle
 //AgtPduPayload GetPayloadBytes
 return "", nil
}

func(np *PduPayload) SetPayloadBytes () error {
 //parameters: PduHandle PayloadBytes
 //AgtPduPayload SetPayloadBytes
 return nil
}

func(np *PduPayload) GetPayloadFill ()(string, error) {
 //parameters: PduHandle
 //AgtPduPayload GetPayloadFill
 return "", nil
}

func(np *PduPayload) SetPayloadFill () error {
 //parameters: PduHandle FillType PayloadData
 //AgtPduPayload SetPayloadFill
 return nil
}

func(np *PduPayload) GetSupportedPayloadFillTypes ()(string, error) {
 //parameters: PduHandle
 //AgtPduPayload GetSupportedPayloadFillTypes
 return "", nil
}

