package n2xsdk

type SipIpv6Statistics struct {
 Handler string
}

func(np *SipIpv6Statistics) GetStatistics ()(string, error) {
 //parameters: SessionPoolHandle StatisticsTypeList
 //AgtSipIpv6Statistics GetStatistics
 return "", nil
}

func(np *SipIpv6Statistics) GetInstanceStatistics ()(string, error) {
 //parameters: SessionPoolInstance StatisticsInstanceTypeList
 //AgtSipIpv6Statistics GetInstanceStatistics
 return "", nil
}

func(np *SipIpv6Statistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtSipIpv6Statistics GetStatisticsUpdateCount
 return "", nil
}

