package n2xsdk

type HttpIpv4ClientDnsEntries struct {
 Handler string
}

func(np *HttpIpv4ClientDnsEntries) Add () error {
 //parameters: SessionPoolHandle HostName
 //AgtHttpIpv4ClientDnsEntries Add
 return nil
}

func(np *HttpIpv4ClientDnsEntries) List ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtHttpIpv4ClientDnsEntries List
 return "", nil
}

func(np *HttpIpv4ClientDnsEntries) Remove () error {
 //parameters: SessionPoolHandle DnsEntryRowIndex
 //AgtHttpIpv4ClientDnsEntries Remove
 return nil
}

func(np *HttpIpv4ClientDnsEntries) RemoveAll () error {
 //parameters: SessionPoolHandle
 //AgtHttpIpv4ClientDnsEntries RemoveAll
 return nil
}

func(np *HttpIpv4ClientDnsEntries) SetHostName () error {
 //parameters: SessionPoolHandle DnsEntryRowIndex HostName
 //AgtHttpIpv4ClientDnsEntries SetHostName
 return nil
}

func(np *HttpIpv4ClientDnsEntries) GetHostName ()(string, error) {
 //parameters: SessionPoolHandle DnsEntryRowIndex
 //AgtHttpIpv4ClientDnsEntries GetHostName
 return "", nil
}

func(np *HttpIpv4ClientDnsEntries) SetIpv4Address () error {
 //parameters: SessionPoolHandle DnsEntryRowIndex Ipv4Address
 //AgtHttpIpv4ClientDnsEntries SetIpv4Address
 return nil
}

func(np *HttpIpv4ClientDnsEntries) GetIpv4Address ()(string, error) {
 //parameters: SessionPoolHandle DnsEntryRowIndex
 //AgtHttpIpv4ClientDnsEntries GetIpv4Address
 return "", nil
}

func(np *HttpIpv4ClientDnsEntries) SetIpv6Address () error {
 //parameters: SessionPoolHandle DnsEntryRowIndex Ipv6Address
 //AgtHttpIpv4ClientDnsEntries SetIpv6Address
 return nil
}

func(np *HttpIpv4ClientDnsEntries) GetIpv6Address ()(string, error) {
 //parameters: SessionPoolHandle DnsEntryRowIndex
 //AgtHttpIpv4ClientDnsEntries GetIpv6Address
 return "", nil
}

func(np *HttpIpv4ClientDnsEntries) SetDefaultIpv4Address () error {
 //parameters: SessionPoolHandle DefaultIpv4Address
 //AgtHttpIpv4ClientDnsEntries SetDefaultIpv4Address
 return nil
}

func(np *HttpIpv4ClientDnsEntries) GetDefaultIpv4Address ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtHttpIpv4ClientDnsEntries GetDefaultIpv4Address
 return "", nil
}

func(np *HttpIpv4ClientDnsEntries) SetDefaultIpv6Address () error {
 //parameters: SessionPoolHandle DefaultIpv6Address
 //AgtHttpIpv4ClientDnsEntries SetDefaultIpv6Address
 return nil
}

func(np *HttpIpv4ClientDnsEntries) GetDefaultIpv6Address ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtHttpIpv4ClientDnsEntries GetDefaultIpv6Address
 return "", nil
}

