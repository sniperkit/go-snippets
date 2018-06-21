package n2xsdk

type HttpIpv6ServerStatistics struct {
 Handler string
}

func(np *HttpIpv6ServerStatistics) GetStatistics ()(string, error) {
 //parameters: SessionIdentifiers StatisticsTypes
 //AgtHttpIpv6ServerStatistics GetStatistics
 return "", nil
}

func(np *HttpIpv6ServerStatistics) GetInstanceStatistics ()(string, error) {
 //parameters: SessionIdentifiers StatisticsTypes
 //AgtHttpIpv6ServerStatistics GetInstanceStatistics
 return "", nil
}

func(np *HttpIpv6ServerStatistics) ListAvailableStatisticsTypes ()(string, error) {
 //parameters: 
 //AgtHttpIpv6ServerStatistics ListAvailableStatisticsTypes
 return "", nil
}

func(np *HttpIpv6ServerStatistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: SessionIdentifiers
 //AgtHttpIpv6ServerStatistics GetStatisticsUpdateCount
 return "", nil
}

