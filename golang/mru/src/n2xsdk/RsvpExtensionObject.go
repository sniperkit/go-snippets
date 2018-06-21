package n2xsdk

type RsvpExtensionObject struct {
 Handler string
}

func(np *RsvpExtensionObjec) SetPduHandle () error {
 //parameters: ObjectHandle PduHandle
 //AgtRsvpExtensionObject SetPduHandle
 return nil
}

func(np *RsvpExtensionObjec) GetPduHandle ()(string, error) {
 //parameters: ObjectHandle
 //AgtRsvpExtensionObject GetPduHandle
 return "", nil
}

func(np *RsvpExtensionObjec) IsGeneric () error {
 //parameters: ObjectHandle
 //AgtRsvpExtensionObject IsGeneric
 return nil
}

