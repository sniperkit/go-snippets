package n2xsdk

type EthernetIpv6Addresses2 struct {
 Handler string
}

func(np *EthernetIpv6Addresses2) AddSutIpv6Address () error {
 //parameters: PortHandle VlanNum psaVlanIds SutIpv6Address
 //AgtEthernetIpv6Addresses2 AddSutIpv6Address
 return nil
}

func(np *EthernetIpv6Addresses2) RemoveSutIpv6Address () error {
 //parameters: PortHandle VlanNum psaVlanIds SutIpv6Address
 //AgtEthernetIpv6Addresses2 RemoveSutIpv6Address
 return nil
}

func(np *EthernetIpv6Addresses2) RemoveAllIpv6SutAddresses () error {
 //parameters: PortHandle
 //AgtEthernetIpv6Addresses2 RemoveAllIpv6SutAddresses
 return nil
}

func(np *EthernetIpv6Addresses2) AddSutIpv6Addresses () error {
 //parameters: PortHandle SutFirstIpv6Address SutNumAddresses SutAddressModifier VlanNum psaFirstVlanIds IncrementsNum psaVlanIdIncrements
 //AgtEthernetIpv6Addresses2 AddSutIpv6Addresses
 return nil
}

func(np *EthernetIpv6Addresses2) CheckAddSutIpv6Addresses () error {
 //parameters: PortHandle SutFirstIpv6Address SutNumAddresses SutAddressModifier VlanNum psaFirstVlanIds IncrementsNum psaVlanIdIncrements
 //AgtEthernetIpv6Addresses2 CheckAddSutIpv6Addresses
 return nil
}

func(np *EthernetIpv6Addresses2) GetSutIpv6AddressCount ()(string, error) {
 //parameters: PortHandle
 //AgtEthernetIpv6Addresses2 GetSutIpv6AddressCount
 return "", nil
}

func(np *EthernetIpv6Addresses2) SetSutMacAddress () error {
 //parameters: PortHandle VlanNum psaVlanIds SutIpv6Address MacAddress
 //AgtEthernetIpv6Addresses2 SetSutMacAddress
 return nil
}

func(np *EthernetIpv6Addresses2) GetSutMacAddress ()(string, error) {
 //parameters: PortHandle VlanNum psaVlanIds SutIpv6Address
 //AgtEthernetIpv6Addresses2 GetSutMacAddress
 return "", nil
}

func(np *EthernetIpv6Addresses2) AddAddressPool () error {
 //parameters: PortHandle TesterFirstIpv6Address TesterPrefixLength TesterNumAddresses TesterAddressModifier VlanNum psaVlanIds
 //AgtEthernetIpv6Addresses2 AddAddressPool
 return nil
}

func(np *EthernetIpv6Addresses2) CheckAddAddressPool () error {
 //parameters: PortHandle TesterFirstIpv6Address TesterPrefixLength TesterNumAddresses TesterAddressModifier VlanNum psaVlanIds
 //AgtEthernetIpv6Addresses2 CheckAddAddressPool
 return nil
}

func(np *EthernetIpv6Addresses2) AddAddressPools () error {
 //parameters: PortHandle NumOfAddressPools TesterFirstIpv6Address TesterPrefixLength TesterNumAddressesPerPool InPoolAddressModifier AcrossPoolAddressModifier VlanNum psaFirstVlanIds IncrementsNum psaVlanIdIncrements
 //AgtEthernetIpv6Addresses2 AddAddressPools
 return nil
}

func(np *EthernetIpv6Addresses2) CheckAddAddressPools () error {
 //parameters: PortHandle NumOfAddressPools TesterFirstIpv6Address TesterPrefixLength TesterNumAddressesPerPool InPoolAddressModifier AcrossPoolAddressModifier VlanNum psaFirstVlanIds IncrementsNum psaVlanIdIncrements
 //AgtEthernetIpv6Addresses2 CheckAddAddressPools
 return nil
}

func(np *EthernetIpv6Addresses2) RemoveAddressPool () error {
 //parameters: PortHandle AddressPoolHandle
 //AgtEthernetIpv6Addresses2 RemoveAddressPool
 return nil
}

func(np *EthernetIpv6Addresses2) RemoveAllAddressPools () error {
 //parameters: PortHandle
 //AgtEthernetIpv6Addresses2 RemoveAllAddressPools
 return nil
}

func(np *EthernetIpv6Addresses2) GetAddressPoolCount ()(string, error) {
 //parameters: PortHandle
 //AgtEthernetIpv6Addresses2 GetAddressPoolCount
 return "", nil
}

func(np *EthernetIpv6Addresses2) GetAddressPoolLimit ()(string, error) {
 //parameters: PortHandle
 //AgtEthernetIpv6Addresses2 GetAddressPoolLimit
 return "", nil
}

func(np *EthernetIpv6Addresses2) ListAddressPools ()(string, error) {
 //parameters: PortHandle
 //AgtEthernetIpv6Addresses2 ListAddressPools
 return "", nil
}

func(np *EthernetIpv6Addresses2) GetVlanCount ()(string, error) {
 //parameters: PortHandle
 //AgtEthernetIpv6Addresses2 GetVlanCount
 return "", nil
}

func(np *EthernetIpv6Addresses2) ListVlans ()(string, error) {
 //parameters: PortHandle
 //AgtEthernetIpv6Addresses2 ListVlans
 return "", nil
}

func(np *EthernetIpv6Addresses2) GetSutIpv6AddressCountByVlan ()(string, error) {
 //parameters: PortHandle VlanNum psaVlanIds
 //AgtEthernetIpv6Addresses2 GetSutIpv6AddressCountByVlan
 return "", nil
}

func(np *EthernetIpv6Addresses2) ListSutIpv6AddressesByVlan ()(string, error) {
 //parameters: PortHandle VlanNum psaVlanIds
 //AgtEthernetIpv6Addresses2 ListSutIpv6AddressesByVlan
 return "", nil
}

func(np *EthernetIpv6Addresses2) ListSutIpv6LinkLocalAddressesByVlan ()(string, error) {
 //parameters: PortHandle VlanNum psaVlanIds
 //AgtEthernetIpv6Addresses2 ListSutIpv6LinkLocalAddressesByVlan
 return "", nil
}

func(np *EthernetIpv6Addresses2) GetAddressPoolCountByVlan ()(string, error) {
 //parameters: PortHandle VlanNum psaVlanIds
 //AgtEthernetIpv6Addresses2 GetAddressPoolCountByVlan
 return "", nil
}

func(np *EthernetIpv6Addresses2) ListAddressPoolsByVlan ()(string, error) {
 //parameters: PortHandle VlanNum psaVlanIds
 //AgtEthernetIpv6Addresses2 ListAddressPoolsByVlan
 return "", nil
}

func(np *EthernetIpv6Addresses2) RemoveVlan () error {
 //parameters: PortHandle VlanNum psaVlanIds
 //AgtEthernetIpv6Addresses2 RemoveVlan
 return nil
}

func(np *EthernetIpv6Addresses2) RemoveAllVlans () error {
 //parameters: PortHandle
 //AgtEthernetIpv6Addresses2 RemoveAllVlans
 return nil
}

func(np *EthernetIpv6Addresses2) SetVlanEtherType () error {
 //parameters: PortHandle VlanIndex EtherType
 //AgtEthernetIpv6Addresses2 SetVlanEtherType
 return nil
}

func(np *EthernetIpv6Addresses2) GetVlanEtherType ()(string, error) {
 //parameters: PortHandle VlanIndex
 //AgtEthernetIpv6Addresses2 GetVlanEtherType
 return "", nil
}

func(np *EthernetIpv6Addresses2) ListStackedVlans ()(string, error) {
 //parameters: PortHandle
 //AgtEthernetIpv6Addresses2 ListStackedVlans
 return "", nil
}

