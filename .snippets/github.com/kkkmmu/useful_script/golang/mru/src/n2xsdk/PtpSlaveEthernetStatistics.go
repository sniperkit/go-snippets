package n2xsdk

type PtpSlaveEthernetStatistics struct {
 Handler string
}

func(np *PtpSlaveEthernetStatistics) ListAvailableStatisticsTypes ()(string, error) {
 //parameters: 
 //AgtPtpSlaveEthernetStatistics ListAvailableStatisticsTypes
 return "", nil
}

func(np *PtpSlaveEthernetStatistics) GetStatistics ()(string, error) {
 //parameters: ClockIdentifier StatisticsTypes
 //AgtPtpSlaveEthernetStatistics GetStatistics
 return "", nil
}

func(np *PtpSlaveEthernetStatistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: ClockIdentifier
 //AgtPtpSlaveEthernetStatistics GetStatisticsUpdateCount
 return "", nil
}

func(np *PtpSlaveEthernetStatistics) GetInstanceStatistics ()(string, error) {
 //parameters: ClockIdentifier StatisticsTypes
 //AgtPtpSlaveEthernetStatistics GetInstanceStatistics
 return "", nil
}

