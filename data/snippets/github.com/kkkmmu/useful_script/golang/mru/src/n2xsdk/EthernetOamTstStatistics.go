package n2xsdk

type EthernetOamTstStatistics struct {
 Handler string
}

func(np *EthernetOamTstStatistics) GetSessionAccumulatedValues ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtEthernetOamTstStatistics GetSessionAccumulatedValues
 return "", nil
}

func(np *EthernetOamTstStatistics) GetStatistic ()(string, error) {
 //parameters: SessionPoolHandle Statistic
 //AgtEthernetOamTstStatistics GetStatistic
 return "", nil
}

