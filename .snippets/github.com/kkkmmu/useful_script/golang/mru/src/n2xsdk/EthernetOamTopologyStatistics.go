package n2xsdk

type EthernetOamTopologyStatistics struct {
 Handler string
}

func(np *EthernetOamTopologyStatistics) GetSessionAccumulatedValues ()(string, error) {
 //parameters: SessionPoolHandle
 //AgtEthernetOamTopologyStatistics GetSessionAccumulatedValues
 return "", nil
}

func(np *EthernetOamTopologyStatistics) GetStatistic ()(string, error) {
 //parameters: SessionPoolHandle Statistic
 //AgtEthernetOamTopologyStatistics GetStatistic
 return "", nil
}

