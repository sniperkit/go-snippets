package n2xsdk

type Ospfv2RouterPoolStatistics struct {
 Handler string
}

func(np *Ospfv2RouterPoolStatistics) GetStatistics ()(string, error) {
 //parameters: RouterIdentifiers StatisticsTypes
 //AgtOspfv2RouterPoolStatistics GetStatistics
 return "", nil
}

func(np *Ospfv2RouterPoolStatistics) GetInstanceStatistics ()(string, error) {
 //parameters: RouterIdentifiers StatisticsTypes
 //AgtOspfv2RouterPoolStatistics GetInstanceStatistics
 return "", nil
}

func(np *Ospfv2RouterPoolStatistics) ListAvailableStatisticsTypes ()(string, error) {
 //parameters: 
 //AgtOspfv2RouterPoolStatistics ListAvailableStatisticsTypes
 return "", nil
}

func(np *Ospfv2RouterPoolStatistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: RouterIdentifiers
 //AgtOspfv2RouterPoolStatistics GetStatisticsUpdateCount
 return "", nil
}

