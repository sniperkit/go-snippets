package n2xsdk

type RstpStatistics struct {
 Handler string
}

func(np *RstpStatistics) GetSessionAccumulatedValues ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtRstpStatistics GetSessionAccumulatedValues
 return "", nil
}

func(np *RstpStatistics) GetStatistic ()(string, error) {
 //parameters: SessionPoolHandle Statistic
 //AgtRstpStatistics GetStatistic
 return "", nil
}

