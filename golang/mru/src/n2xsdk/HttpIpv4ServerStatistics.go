package n2xsdk

type HttpIpv4ServerStatistics struct {
 Handler string
}

func(np *HttpIpv4ServerStatistics) GetStatistics ()(string, error) {
 //parameters: SessionIdentifiers StatisticsTypes
 //AgtHttpIpv4ServerStatistics GetStatistics
 return "", nil
}

func(np *HttpIpv4ServerStatistics) GetInstanceStatistics ()(string, error) {
 //parameters: SessionIdentifiers StatisticsTypes
 //AgtHttpIpv4ServerStatistics GetInstanceStatistics
 return "", nil
}

func(np *HttpIpv4ServerStatistics) ListAvailableStatisticsTypes ()(string, error) {
 //parameters: 
 //AgtHttpIpv4ServerStatistics ListAvailableStatisticsTypes
 return "", nil
}

func(np *HttpIpv4ServerStatistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: SessionIdentifiers
 //AgtHttpIpv4ServerStatistics GetStatisticsUpdateCount
 return "", nil
}

