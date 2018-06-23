package n2xsdk

type PppAddresses struct {
 Handler string
}

func(np *PppAddresses) SetSutIpAddress () error {
 //parameters: PortHandle SutIpAddress PrefixLength
 //AgtPppAddresses SetSutIpAddress
 return nil
}

func(np *PppAddresses) SetSutAndTesterIpAddresses () error {
 //parameters: PortHandle SutIpAddress TesterIpAddress PrefixLength
 //AgtPppAddresses SetSutAndTesterIpAddresses
 return nil
}

func(np *PppAddresses) GetSutIpAddress ()(string, error) {
 //parameters: PortHandle
 //AgtPppAddresses GetSutIpAddress
 return "", nil
}

func(np *PppAddresses) GetTesterIpAddress ()(string, error) {
 //parameters: PortHandle
 //AgtPppAddresses GetTesterIpAddress
 return "", nil
}

func(np *PppAddresses) IsTesterIpAddressCalculated () error {
 //parameters: PortHandle
 //AgtPppAddresses IsTesterIpAddressCalculated
 return nil
}

func(np *PppAddresses) SetSutIpv6Address () error {
 //parameters: PortHandle SutIpv6Address PrefixLength
 //AgtPppAddresses SetSutIpv6Address
 return nil
}

func(np *PppAddresses) SetSutAndTesterIpv6Addresses () error {
 //parameters: PortHandle SutIpv6Address TesterIpv6Address PrefixLength
 //AgtPppAddresses SetSutAndTesterIpv6Addresses
 return nil
}

func(np *PppAddresses) GetSutIpv6Address ()(string, error) {
 //parameters: PortHandle
 //AgtPppAddresses GetSutIpv6Address
 return "", nil
}

func(np *PppAddresses) GetTesterIpv6Address ()(string, error) {
 //parameters: PortHandle
 //AgtPppAddresses GetTesterIpv6Address
 return "", nil
}

func(np *PppAddresses) IsTesterIpv6AddressCalculated () error {
 //parameters: PortHandle
 //AgtPppAddresses IsTesterIpv6AddressCalculated
 return nil
}

