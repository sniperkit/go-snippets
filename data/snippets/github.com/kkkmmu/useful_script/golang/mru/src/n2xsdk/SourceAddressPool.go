package n2xsdk

type SourceAddressPool struct {
 Handler string
}

func(np *SourceAddressPool) GetType ()(string, error) {
 //parameters: Handle
 //AgtSourceAddressPool GetType
 return "", nil
}

func(np *SourceAddressPool) GetName ()(string, error) {
 //parameters: Handle
 //AgtSourceAddressPool GetName
 return "", nil
}

func(np *SourceAddressPool) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtSourceAddressPool GetLockCount
 return "", nil
}

func(np *SourceAddressPool) SetSourceAddressRange () error {
 //parameters: PoolHandle FirstSourceAddress PrefixLength NumAddresses Modifier
 //AgtSourceAddressPool SetSourceAddressRange
 return nil
}

func(np *SourceAddressPool) GetSourceAddressRange ()(string, error) {
 //parameters: PoolHandle
 //AgtSourceAddressPool GetSourceAddressRange
 return "", nil
}

func(np *SourceAddressPool) SetSourceAddressPoolVirtualInterface () error {
 //parameters: PoolHandle Count pInterfaceHandles
 //AgtSourceAddressPool SetSourceAddressPoolVirtualInterface
 return nil
}

func(np *SourceAddressPool) GetSourceAddressPoolVirtualInterface ()(string, error) {
 //parameters: PoolHandle
 //AgtSourceAddressPool GetSourceAddressPoolVirtualInterface
 return "", nil
}

