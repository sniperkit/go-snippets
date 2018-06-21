package n2xsdk

type StpStatistics struct {
 Handler string
}

func(np *StpStatistics) GetSessionAccumulatedValues ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtStpStatistics GetSessionAccumulatedValues
 return "", nil
}

func(np *StpStatistics) GetStatistic ()(string, error) {
 //parameters: SessionPoolHandle Statistic
 //AgtStpStatistics GetStatistic
 return "", nil
}

