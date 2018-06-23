package n2xsdk

type EthernetAddressPool struct {
 Handler string
}

func(np *EthernetAddressPool) SetTesterIpAddresses () error {
 //parameters: hAddressPool FirstIpAddress PrefixLength NumAddresses Modifier
 //AgtEthernetAddressPool SetTesterIpAddresses
 return nil
}

func(np *EthernetAddressPool) GetTesterIpAddresses ()(string, error) {
 //parameters: hAddressPool
 //AgtEthernetAddressPool GetTesterIpAddresses
 return "", nil
}

func(np *EthernetAddressPool) GetNthTesterIpAddress ()(string, error) {
 //parameters: hAddressPool Index
 //AgtEthernetAddressPool GetNthTesterIpAddress
 return "", nil
}

func(np *EthernetAddressPool) SetTesterAndSutIpAddresses () error {
 //parameters: hAddressPool FirstIpAddress PrefixLength NumAddresses Modifier OldSutIpAddress NewSutIpAddress
 //AgtEthernetAddressPool SetTesterAndSutIpAddresses
 return nil
}

func(np *EthernetAddressPool) SetVlanId () error {
 //parameters: hAddressPool VlanId
 //AgtEthernetAddressPool SetVlanId
 return nil
}

func(np *EthernetAddressPool) SetVlanIds () error {
 //parameters: hAddressPool VlanNum psaVlanIds
 //AgtEthernetAddressPool SetVlanIds
 return nil
}

func(np *EthernetAddressPool) GetVlanId ()(string, error) {
 //parameters: hAddressPool
 //AgtEthernetAddressPool GetVlanId
 return "", nil
}

func(np *EthernetAddressPool) GetVlanIds ()(string, error) {
 //parameters: hAddressPool
 //AgtEthernetAddressPool GetVlanIds
 return "", nil
}

func(np *EthernetAddressPool) EnableVlan () error {
 //parameters: hAddressPool
 //AgtEthernetAddressPool EnableVlan
 return nil
}

func(np *EthernetAddressPool) DisableVlan () error {
 //parameters: hAddressPool
 //AgtEthernetAddressPool DisableVlan
 return nil
}

func(np *EthernetAddressPool) IsVlanEnabled () error {
 //parameters: hAddressPool
 //AgtEthernetAddressPool IsVlanEnabled
 return nil
}

func(np *EthernetAddressPool) EnableTrafficDestinations () error {
 //parameters: hAddressPool
 //AgtEthernetAddressPool EnableTrafficDestinations
 return nil
}

func(np *EthernetAddressPool) EnableTrafficDestinationsByPort () error {
 //parameters: PortHandle
 //AgtEthernetAddressPool EnableTrafficDestinationsByPort
 return nil
}

func(np *EthernetAddressPool) DisableTrafficDestinations () error {
 //parameters: hAddressPool
 //AgtEthernetAddressPool DisableTrafficDestinations
 return nil
}

func(np *EthernetAddressPool) DisableTrafficDestinationsByPort () error {
 //parameters: PortHandle
 //AgtEthernetAddressPool DisableTrafficDestinationsByPort
 return nil
}

func(np *EthernetAddressPool) IsTrafficDestinationEnabled () error {
 //parameters: hAddressPool
 //AgtEthernetAddressPool IsTrafficDestinationEnabled
 return nil
}

func(np *EthernetAddressPool) SetTesterMacAddressUniqueFlag () error {
 //parameters: hAddressPool UniqueMacAddressFlag
 //AgtEthernetAddressPool SetTesterMacAddressUniqueFlag
 return nil
}

func(np *EthernetAddressPool) GetTesterMacAddressUniqueFlag ()(string, error) {
 //parameters: hAddressPool
 //AgtEthernetAddressPool GetTesterMacAddressUniqueFlag
 return "", nil
}

func(np *EthernetAddressPool) SetTesterMacAddresses () error {
 //parameters: hAddressPool FirstMacAddress UniqueMacAddressFlag
 //AgtEthernetAddressPool SetTesterMacAddresses
 return nil
}

func(np *EthernetAddressPool) GetTesterMacAddresses ()(string, error) {
 //parameters: hAddressPool
 //AgtEthernetAddressPool GetTesterMacAddresses
 return "", nil
}

func(np *EthernetAddressPool) GetNthTesterMacAddress ()(string, error) {
 //parameters: hAddressPool Index
 //AgtEthernetAddressPool GetNthTesterMacAddress
 return "", nil
}

func(np *EthernetAddressPool) GetAddressPoolInfo ()(string, error) {
 //parameters: hAddressPool
 //AgtEthernetAddressPool GetAddressPoolInfo
 return "", nil
}

func(np *EthernetAddressPool) GetVlanStackSize ()(string, error) {
 //parameters: hAddressPool
 //AgtEthernetAddressPool GetVlanStackSize
 return "", nil
}

