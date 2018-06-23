package n2xsdk

type HttpIpv6ClientDnsEntries struct {
 Handler string
}

func(np *HttpIpv6ClientDnsEntries) Add () error {
 //parameters: SessionPoolHandle HostName
 //AgtHttpIpv6ClientDnsEntries Add
 return nil
}

func(np *HttpIpv6ClientDnsEntries) List ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientDnsEntries List
 return "", nil
}

func(np *HttpIpv6ClientDnsEntries) Remove () error {
 //parameters: SessionPoolHandle DnsEntryRowIndex
 //AgtHttpIpv6ClientDnsEntries Remove
 return nil
}

func(np *HttpIpv6ClientDnsEntries) RemoveAll () error {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientDnsEntries RemoveAll
 return nil
}

func(np *HttpIpv6ClientDnsEntries) SetHostName () error {
 //parameters: SessionPoolHandle DnsEntryRowIndex HostName
 //AgtHttpIpv6ClientDnsEntries SetHostName
 return nil
}

func(np *HttpIpv6ClientDnsEntries) GetHostName ()(string, error) {
 //parameters: SessionPoolHandle DnsEntryRowIndex
 //AgtHttpIpv6ClientDnsEntries GetHostName
 return "", nil
}

func(np *HttpIpv6ClientDnsEntries) SetIpv4Address () error {
 //parameters: SessionPoolHandle DnsEntryRowIndex Ipv4Address
 //AgtHttpIpv6ClientDnsEntries SetIpv4Address
 return nil
}

func(np *HttpIpv6ClientDnsEntries) GetIpv4Address ()(string, error) {
 //parameters: SessionPoolHandle DnsEntryRowIndex
 //AgtHttpIpv6ClientDnsEntries GetIpv4Address
 return "", nil
}

func(np *HttpIpv6ClientDnsEntries) SetIpv6Address () error {
 //parameters: SessionPoolHandle DnsEntryRowIndex Ipv6Address
 //AgtHttpIpv6ClientDnsEntries SetIpv6Address
 return nil
}

func(np *HttpIpv6ClientDnsEntries) GetIpv6Address ()(string, error) {
 //parameters: SessionPoolHandle DnsEntryRowIndex
 //AgtHttpIpv6ClientDnsEntries GetIpv6Address
 return "", nil
}

func(np *HttpIpv6ClientDnsEntries) SetDefaultIpv4Address () error {
 //parameters: SessionPoolHandle DefaultIpv4Address
 //AgtHttpIpv6ClientDnsEntries SetDefaultIpv4Address
 return nil
}

func(np *HttpIpv6ClientDnsEntries) GetDefaultIpv4Address ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientDnsEntries GetDefaultIpv4Address
 return "", nil
}

func(np *HttpIpv6ClientDnsEntries) SetDefaultIpv6Address () error {
 //parameters: SessionPoolHandle DefaultIpv6Address
 //AgtHttpIpv6ClientDnsEntries SetDefaultIpv6Address
 return nil
}

func(np *HttpIpv6ClientDnsEntries) GetDefaultIpv6Address ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtHttpIpv6ClientDnsEntries GetDefaultIpv6Address
 return "", nil
}

