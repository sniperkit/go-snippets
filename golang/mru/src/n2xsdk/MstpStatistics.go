package n2xsdk

type MstpStatistics struct {
 Handler string
}

func(np *MstpStatistics) GetMstiStatistic ()(string, error) {
 //parameters: SessionPoolHandle MstiIdentifier MstiStatistic
 //AgtMstpStatistics GetMstiStatistic
 return "", nil
}

func(np *MstpStatistics) GetSessionAccumulatedValues ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtMstpStatistics GetSessionAccumulatedValues
 return "", nil
}

func(np *MstpStatistics) GetStatistic ()(string, error) {
 //parameters: SessionPoolHandle Statistic
 //AgtMstpStatistics GetStatistic
 return "", nil
}

