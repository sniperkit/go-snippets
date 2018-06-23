package n2xsdk

type LnsStatistics struct {
 Handler string
}

func(np *LnsStatistics) GetStatistics ()(string, error) {
 //parameters: SessionIdentifiers StatisticsTypes
 //AgtLnsStatistics GetStatistics
 return "", nil
}

func(np *LnsStatistics) GetInstanceStatistics ()(string, error) {
 //parameters: SessionIdentifiers StatisticsTypes
 //AgtLnsStatistics GetInstanceStatistics
 return "", nil
}

func(np *LnsStatistics) ListAvailableStatisticsTypes ()(string, error) {
 //parameters: 
 //AgtLnsStatistics ListAvailableStatisticsTypes
 return "", nil
}

func(np *LnsStatistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: SessionIdentifiers
 //AgtLnsStatistics GetStatisticsUpdateCount
 return "", nil
}

