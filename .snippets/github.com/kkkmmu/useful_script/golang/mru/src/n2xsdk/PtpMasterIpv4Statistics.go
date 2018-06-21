package n2xsdk

type PtpMasterIpv4Statistics struct {
 Handler string
}

func(np *PtpMasterIpv4Statistics) ListAvailableStatisticsTypes ()(string, error) {
 //parameters: 
 //AgtPtpMasterIpv4Statistics ListAvailableStatisticsTypes
 return "", nil
}

func(np *PtpMasterIpv4Statistics) GetStatistics ()(string, error) {
 //parameters: ClockIdentifier StatisticsTypes
 //AgtPtpMasterIpv4Statistics GetStatistics
 return "", nil
}

func(np *PtpMasterIpv4Statistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: ClockIdentifier
 //AgtPtpMasterIpv4Statistics GetStatisticsUpdateCount
 return "", nil
}

