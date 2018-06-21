package n2xsdk

type PtpSlaveIpv6Statistics struct {
 Handler string
}

func(np *PtpSlaveIpv6Statistics) ListAvailableStatisticsTypes ()(string, error) {
 //parameters: 
 //AgtPtpSlaveIpv6Statistics ListAvailableStatisticsTypes
 return "", nil
}

func(np *PtpSlaveIpv6Statistics) GetStatistics ()(string, error) {
 //parameters: ClockIdentifier StatisticsTypes
 //AgtPtpSlaveIpv6Statistics GetStatistics
 return "", nil
}

func(np *PtpSlaveIpv6Statistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: ClockIdentifier
 //AgtPtpSlaveIpv6Statistics GetStatisticsUpdateCount
 return "", nil
}

func(np *PtpSlaveIpv6Statistics) GetInstanceStatistics ()(string, error) {
 //parameters: ClockIdentifier StatisticsTypes
 //AgtPtpSlaveIpv6Statistics GetInstanceStatistics
 return "", nil
}

