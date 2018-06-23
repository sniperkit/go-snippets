package n2xsdk

type SipStatistics struct {
 Handler string
}

func(np *SipStatistics) GetStatistics ()(string, error) {
 //parameters: SessionPoolHandle StatisticsTypeList
 //AgtSipStatistics GetStatistics
 return "", nil
}

func(np *SipStatistics) GetInstanceStatistics ()(string, error) {
 //parameters: SessionPoolInstance StatisticsInstanceTypeList
 //AgtSipStatistics GetInstanceStatistics
 return "", nil
}

func(np *SipStatistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtSipStatistics GetStatisticsUpdateCount
 return "", nil
}

