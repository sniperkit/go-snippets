package n2xsdk

type EthernetIpv6AddressPool struct {
 Handler string
}

func(np *EthernetIpv6AddressPool) SetTesterAddresses () error {
 //parameters: hAddressPool FirstIpv6Address PrefixLength NumAddresses AddressModifier VlanNum psaVlanIds
 //AgtEthernetIpv6AddressPool SetTesterAddresses
 return nil
}

func(np *EthernetIpv6AddressPool) GetTesterAddresses ()(string, error) {
 //parameters: hAddressPool
 //AgtEthernetIpv6AddressPool GetTesterAddresses
 return "", nil
}

func(np *EthernetIpv6AddressPool) GetTesterLinkLocalAddresses ()(string, error) {
 //parameters: hAddressPool
 //AgtEthernetIpv6AddressPool GetTesterLinkLocalAddresses
 return "", nil
}

func(np *EthernetIpv6AddressPool) CheckSetTesterAddresses () error {
 //parameters: hAddressPool FirstIpv6Address PrefixLength NumAddresses AddressModifier VlanNum psaVlanIds
 //AgtEthernetIpv6AddressPool CheckSetTesterAddresses
 return nil
}

func(np *EthernetIpv6AddressPool) GetNumAddresses ()(string, error) {
 //parameters: hAddressPool
 //AgtEthernetIpv6AddressPool GetNumAddresses
 return "", nil
}

func(np *EthernetIpv6AddressPool) GetNthTesterAddress ()(string, error) {
 //parameters: hAddressPool Index
 //AgtEthernetIpv6AddressPool GetNthTesterAddress
 return "", nil
}

func(np *EthernetIpv6AddressPool) EnableTesterAddressAutoconfiguration () error {
 //parameters: hAddressPool
 //AgtEthernetIpv6AddressPool EnableTesterAddressAutoconfiguration
 return nil
}

func(np *EthernetIpv6AddressPool) DisableTesterAddressAutoconfiguration () error {
 //parameters: hAddressPool
 //AgtEthernetIpv6AddressPool DisableTesterAddressAutoconfiguration
 return nil
}

func(np *EthernetIpv6AddressPool) IsTesterAddressAutoconfigurationEnabled () error {
 //parameters: hAddressPool
 //AgtEthernetIpv6AddressPool IsTesterAddressAutoconfigurationEnabled
 return nil
}

func(np *EthernetIpv6AddressPool) SetTesterVlanId () error {
 //parameters: hAddressPool VlanId
 //AgtEthernetIpv6AddressPool SetTesterVlanId
 return nil
}

func(np *EthernetIpv6AddressPool) SetTesterVlanIds () error {
 //parameters: hAddressPool VlanNum psaVlanIds
 //AgtEthernetIpv6AddressPool SetTesterVlanIds
 return nil
}

func(np *EthernetIpv6AddressPool) GetTesterVlanId ()(string, error) {
 //parameters: hAddressPool
 //AgtEthernetIpv6AddressPool GetTesterVlanId
 return "", nil
}

func(np *EthernetIpv6AddressPool) GetTesterVlanIds ()(string, error) {
 //parameters: hAddressPool
 //AgtEthernetIpv6AddressPool GetTesterVlanIds
 return "", nil
}

func(np *EthernetIpv6AddressPool) SetTesterMacAddressUniqueFlag () error {
 //parameters: hAddressPool UniqueMacAddressFlag
 //AgtEthernetIpv6AddressPool SetTesterMacAddressUniqueFlag
 return nil
}

func(np *EthernetIpv6AddressPool) GetTesterMacAddressUniqueFlag ()(string, error) {
 //parameters: hAddressPool
 //AgtEthernetIpv6AddressPool GetTesterMacAddressUniqueFlag
 return "", nil
}

func(np *EthernetIpv6AddressPool) SetTesterMacAddresses () error {
 //parameters: hAddressPool FirstMacAddress UniqueMacAddressFlag
 //AgtEthernetIpv6AddressPool SetTesterMacAddresses
 return nil
}

func(np *EthernetIpv6AddressPool) GetTesterMacAddresses ()(string, error) {
 //parameters: hAddressPool
 //AgtEthernetIpv6AddressPool GetTesterMacAddresses
 return "", nil
}

func(np *EthernetIpv6AddressPool) GetNthTesterMacAddress ()(string, error) {
 //parameters: hAddressPool Index
 //AgtEthernetIpv6AddressPool GetNthTesterMacAddress
 return "", nil
}

func(np *EthernetIpv6AddressPool) GetAddressPoolInfo ()(string, error) {
 //parameters: hAddressPool
 //AgtEthernetIpv6AddressPool GetAddressPoolInfo
 return "", nil
}

func(np *EthernetIpv6AddressPool) EnableTrafficDestinations () error {
 //parameters: hAddressPool
 //AgtEthernetIpv6AddressPool EnableTrafficDestinations
 return nil
}

func(np *EthernetIpv6AddressPool) EnableTrafficDestinationsByPort () error {
 //parameters: hPort
 //AgtEthernetIpv6AddressPool EnableTrafficDestinationsByPort
 return nil
}

func(np *EthernetIpv6AddressPool) EnableTrafficDestinationsByList ()(string, error) {
 //parameters: hPort Count hAddressPools
 //AgtEthernetIpv6AddressPool EnableTrafficDestinationsByList
 return "", nil
}

func(np *EthernetIpv6AddressPool) DisableTrafficDestinations () error {
 //parameters: hAddressPool
 //AgtEthernetIpv6AddressPool DisableTrafficDestinations
 return nil
}

func(np *EthernetIpv6AddressPool) DisableTrafficDestinationsByPort () error {
 //parameters: hPort
 //AgtEthernetIpv6AddressPool DisableTrafficDestinationsByPort
 return nil
}

func(np *EthernetIpv6AddressPool) DisableTrafficDestinationsByList ()(string, error) {
 //parameters: hPort Count hAddressPools
 //AgtEthernetIpv6AddressPool DisableTrafficDestinationsByList
 return "", nil
}

func(np *EthernetIpv6AddressPool) IsTrafficDestinationEnabled () error {
 //parameters: hAddressPool
 //AgtEthernetIpv6AddressPool IsTrafficDestinationEnabled
 return nil
}

func(np *EthernetIpv6AddressPool) GetVlanStackSize ()(string, error) {
 //parameters: hAddressPool
 //AgtEthernetIpv6AddressPool GetVlanStackSize
 return "", nil
}

