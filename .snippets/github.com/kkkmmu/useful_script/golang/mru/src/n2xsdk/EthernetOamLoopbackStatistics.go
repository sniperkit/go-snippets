package n2xsdk

type EthernetOamLoopbackStatistics struct {
 Handler string
}

func(np *EthernetOamLoopbackStatistics) GetSessionAccumulatedValues ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtEthernetOamLoopbackStatistics GetSessionAccumulatedValues
 return "", nil
}

func(np *EthernetOamLoopbackStatistics) GetStatistic ()(string, error) {
 //parameters: SessionPoolHandle Statistic
 //AgtEthernetOamLoopbackStatistics GetStatistic
 return "", nil
}

