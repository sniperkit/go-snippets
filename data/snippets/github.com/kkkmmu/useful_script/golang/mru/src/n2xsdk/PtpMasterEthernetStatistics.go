package n2xsdk

type PtpMasterEthernetStatistics struct {
 Handler string
}

func(np *PtpMasterEthernetStatistics) ListAvailableStatisticsTypes ()(string, error) {
 //parameters: 
 //AgtPtpMasterEthernetStatistics ListAvailableStatisticsTypes
 return "", nil
}

func(np *PtpMasterEthernetStatistics) GetStatistics ()(string, error) {
 //parameters: ClockIdentifier StatisticsTypes
 //AgtPtpMasterEthernetStatistics GetStatistics
 return "", nil
}

func(np *PtpMasterEthernetStatistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: ClockIdentifier
 //AgtPtpMasterEthernetStatistics GetStatisticsUpdateCount
 return "", nil
}

