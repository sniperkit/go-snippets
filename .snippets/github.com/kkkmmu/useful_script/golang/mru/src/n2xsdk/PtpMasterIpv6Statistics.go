package n2xsdk

type PtpMasterIpv6Statistics struct {
 Handler string
}

func(np *PtpMasterIpv6Statistics) ListAvailableStatisticsTypes ()(string, error) {
 //parameters: 
 //AgtPtpMasterIpv6Statistics ListAvailableStatisticsTypes
 return "", nil
}

func(np *PtpMasterIpv6Statistics) GetStatistics ()(string, error) {
 //parameters: ClockIdentifier StatisticsTypes
 //AgtPtpMasterIpv6Statistics GetStatistics
 return "", nil
}

func(np *PtpMasterIpv6Statistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: ClockIdentifier
 //AgtPtpMasterIpv6Statistics GetStatisticsUpdateCount
 return "", nil
}

