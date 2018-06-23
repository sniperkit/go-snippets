package n2xsdk

type LacStatistics struct {
 Handler string
}

func(np *LacStatistics) GetStatistics ()(string, error) {
 //parameters: SessionIdentifiers StatisticsTypes
 //AgtLacStatistics GetStatistics
 return "", nil
}

func(np *LacStatistics) GetInstanceStatistics ()(string, error) {
 //parameters: SessionIdentifiers StatisticsTypes
 //AgtLacStatistics GetInstanceStatistics
 return "", nil
}

func(np *LacStatistics) ListAvailableStatisticsTypes ()(string, error) {
 //parameters: 
 //AgtLacStatistics ListAvailableStatisticsTypes
 return "", nil
}

func(np *LacStatistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: SessionIdentifiers
 //AgtLacStatistics GetStatisticsUpdateCount
 return "", nil
}

