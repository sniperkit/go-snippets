package n2xsdk

type LinkOamStatistics struct {
 Handler string
}

func(np *LinkOamStatistics) GetStatistics ()(string, error) {
 //parameters: SessionIdentifier StatisticsTypeList
 //AgtLinkOamStatistics GetStatistics
 return "", nil
}

func(np *LinkOamStatistics) GetSessionAccumulatedValues ()(string, error) {
 //parameters: SessionIdentifier
 //AgtLinkOamStatistics GetSessionAccumulatedValues
 return "", nil
}

func(np *LinkOamStatistics) GetInstanceStatistics ()(string, error) {
 //parameters: SessionPoolInstance StatisticsInstanceType
 //AgtLinkOamStatistics GetInstanceStatistics
 return "", nil
}

func(np *LinkOamStatistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: SessionIdentifier
 //AgtLinkOamStatistics GetStatisticsUpdateCount
 return "", nil
}

func(np *LinkOamStatistics) ListAvailableStatisticsTypes ()(string, error) {
 //parameters: 
 //AgtLinkOamStatistics ListAvailableStatisticsTypes
 return "", nil
}

