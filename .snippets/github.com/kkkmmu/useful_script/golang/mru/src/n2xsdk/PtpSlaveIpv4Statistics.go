package n2xsdk

type PtpSlaveIpv4Statistics struct {
 Handler string
}

func(np *PtpSlaveIpv4Statistics) ListAvailableStatisticsTypes ()(string, error) {
 //parameters: 
 //AgtPtpSlaveIpv4Statistics ListAvailableStatisticsTypes
 return "", nil
}

func(np *PtpSlaveIpv4Statistics) GetStatistics ()(string, error) {
 //parameters: ClockIdentifier StatisticsTypes
 //AgtPtpSlaveIpv4Statistics GetStatistics
 return "", nil
}

func(np *PtpSlaveIpv4Statistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: ClockIdentifier
 //AgtPtpSlaveIpv4Statistics GetStatisticsUpdateCount
 return "", nil
}

func(np *PtpSlaveIpv4Statistics) GetInstanceStatistics ()(string, error) {
 //parameters: ClockIdentifier StatisticsTypes
 //AgtPtpSlaveIpv4Statistics GetInstanceStatistics
 return "", nil
}

