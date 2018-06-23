package n2xsdk

type EsmcStatistics struct {
 Handler string
}

func(np *EsmcStatistics) ListAvailableStatisticsTypes ()(string, error) {
 //parameters: 
 //AgtEsmcStatistics ListAvailableStatisticsTypes
 return "", nil
}

func(np *EsmcStatistics) GetStatistics ()(string, error) {
 //parameters: SessionHandle StatisticsTypes
 //AgtEsmcStatistics GetStatistics
 return "", nil
}

func(np *EsmcStatistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: SessionHandle
 //AgtEsmcStatistics GetStatisticsUpdateCount
 return "", nil
}

