package n2xsdk

type PimVpnStatistics struct {
 Handler string
}

func(np *PimVpnStatistics) ListAvailableStatisticsTypes ()(string, error) {
 //parameters: 
 //AgtPimVpnStatistics ListAvailableStatisticsTypes
 return "", nil
}

func(np *PimVpnStatistics) GetStatistics ()(string, error) {
 //parameters: RouterIdentifiers StatisticsTypes
 //AgtPimVpnStatistics GetStatistics
 return "", nil
}

func(np *PimVpnStatistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: RouterIdentifiers
 //AgtPimVpnStatistics GetStatisticsUpdateCount
 return "", nil
}

