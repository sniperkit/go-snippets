package n2xsdk

type SipIpv4Statistics struct {
 Handler string
}

func(np *SipIpv4Statistics) GetStatistics ()(string, error) {
 //parameters: SessionPoolHandle StatisticsTypeList
 //AgtSipIpv4Statistics GetStatistics
 return "", nil
}

func(np *SipIpv4Statistics) GetInstanceStatistics ()(string, error) {
 //parameters: SessionPoolInstance StatisticsInstanceTypeList
 //AgtSipIpv4Statistics GetInstanceStatistics
 return "", nil
}

func(np *SipIpv4Statistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtSipIpv4Statistics GetStatisticsUpdateCount
 return "", nil
}

