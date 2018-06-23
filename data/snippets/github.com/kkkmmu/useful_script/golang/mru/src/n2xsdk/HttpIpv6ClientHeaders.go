package n2xsdk

type HttpIpv6ClientHeaders struct {
 Handler string
}

func(np *HttpIpv6ClientHeaders) Add () error {
 //parameters: SessionPoolHandle HeaderFieldName HeaderFieldValue
 //AgtHttpIpv6ClientHeaders Add
 return nil
}

func(np *HttpIpv6ClientHeaders) List ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientHeaders List
 return "", nil
}

func(np *HttpIpv6ClientHeaders) Remove () error {
 //parameters: SessionPoolHandle HeaderRowIndex
 //AgtHttpIpv6ClientHeaders Remove
 return nil
}

func(np *HttpIpv6ClientHeaders) RemoveAll () error {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientHeaders RemoveAll
 return nil
}

func(np *HttpIpv6ClientHeaders) SetFieldName () error {
 //parameters: SessionPoolHandle HeaderRowIndex HeaderFieldName
 //AgtHttpIpv6ClientHeaders SetFieldName
 return nil
}

func(np *HttpIpv6ClientHeaders) GetFieldName ()(string, error) {
 //parameters: SessionPoolHandle HeaderRowIndex
 //AgtHttpIpv6ClientHeaders GetFieldName
 return "", nil
}

func(np *HttpIpv6ClientHeaders) SetFieldValue () error {
 //parameters: SessionPoolHandle HeaderRowIndex HeaderFieldValue
 //AgtHttpIpv6ClientHeaders SetFieldValue
 return nil
}

func(np *HttpIpv6ClientHeaders) GetFieldValue ()(string, error) {
 //parameters: SessionPoolHandle HeaderRowIndex
 //AgtHttpIpv6ClientHeaders GetFieldValue
 return "", nil
}

