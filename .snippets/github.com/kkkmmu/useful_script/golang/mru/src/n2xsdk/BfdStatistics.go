package n2xsdk

type BfdStatistics struct {
 Handler string
}

func(np *BfdStatistics) GetSessionAccumulatedValues ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtBfdStatistics GetSessionAccumulatedValues
 return "", nil
}

func(np *BfdStatistics) GetStatistic ()(string, error) {
 //parameters: SessionPoolHandle Statistic
 //AgtBfdStatistics GetStatistic
 return "", nil
}

func(np *BfdStatistics) ClearSessionStatistics () error {
 //parameters: SessionPoolHandle
 //AgtBfdStatistics ClearSessionStatistics
 return nil
}

func(np *BfdStatistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtBfdStatistics GetStatisticsUpdateCount
 return "", nil
}

