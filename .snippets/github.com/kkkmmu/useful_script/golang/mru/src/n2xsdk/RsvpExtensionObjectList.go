package n2xsdk

type RsvpExtensionObjectList struct {
 Handler string
}

func(np *RsvpExtensionObjectLis) AddExtensionObjectsByPduHandle () error {
 //parameters: Handle NumPduHandles psaPduHandles
 //AgtRsvpExtensionObjectList AddExtensionObjectsByPduHandle
 return nil
}

func(np *RsvpExtensionObjectLis) ListExtensionObjects ()(string, error) {
 //parameters: Handle
 //AgtRsvpExtensionObjectList ListExtensionObjects
 return "", nil
}

func(np *RsvpExtensionObjectLis) RemoveExtensionObjects () error {
 //parameters: Handle Count psaObjectHandles
 //AgtRsvpExtensionObjectList RemoveExtensionObjects
 return nil
}

func(np *RsvpExtensionObjectLis) ClearExtensionObjects () error {
 //parameters: Handle
 //AgtRsvpExtensionObjectList ClearExtensionObjects
 return nil
}

