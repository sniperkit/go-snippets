package n2xsdk

type HttpIpv4ClientHeaders struct {
 Handler string
}

func(np *HttpIpv4ClientHeaders) Add () error {
 //parameters: SessionPoolHandle HeaderFieldName HeaderFieldValue
 //AgtHttpIpv4ClientHeaders Add
 return nil
}

func(np *HttpIpv4ClientHeaders) List ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtHttpIpv4ClientHeaders List
 return "", nil
}

func(np *HttpIpv4ClientHeaders) Remove () error {
 //parameters: SessionPoolHandle HeaderRowIndex
 //AgtHttpIpv4ClientHeaders Remove
 return nil
}

func(np *HttpIpv4ClientHeaders) RemoveAll () error {
 //parameters: SessionPoolHandle
 //AgtHttpIpv4ClientHeaders RemoveAll
 return nil
}

func(np *HttpIpv4ClientHeaders) SetFieldName () error {
 //parameters: SessionPoolHandle HeaderRowIndex HeaderFieldName
 //AgtHttpIpv4ClientHeaders SetFieldName
 return nil
}

func(np *HttpIpv4ClientHeaders) GetFieldName ()(string, error) {
 //parameters: SessionPoolHandle HeaderRowIndex
 //AgtHttpIpv4ClientHeaders GetFieldName
 return "", nil
}

func(np *HttpIpv4ClientHeaders) SetFieldValue () error {
 //parameters: SessionPoolHandle HeaderRowIndex HeaderFieldValue
 //AgtHttpIpv4ClientHeaders SetFieldValue
 return nil
}

func(np *HttpIpv4ClientHeaders) GetFieldValue ()(string, error) {
 //parameters: SessionPoolHandle HeaderRowIndex
 //AgtHttpIpv4ClientHeaders GetFieldValue
 return "", nil
}

