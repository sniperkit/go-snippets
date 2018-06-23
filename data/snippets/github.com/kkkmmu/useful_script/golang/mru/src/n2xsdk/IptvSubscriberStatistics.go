package n2xsdk

type IptvSubscriberStatistics struct {
 Handler string
}

func(np *IptvSubscriberStatistics) GetStatistics ()(string, error) {
 //parameters: SubscriberPoolHandle StatisticsList
 //AgtIptvSubscriberStatistics GetStatistics
 return "", nil
}

func(np *IptvSubscriberStatistics) ListAvailableStatisticsTypes ()(string, error) {
 //parameters: 
 //AgtIptvSubscriberStatistics ListAvailableStatisticsTypes
 return "", nil
}

func(np *IptvSubscriberStatistics) GetStatisticsUpdateCount ()(string, error) {
 //parameters: SubscriberPoolHandle
 //AgtIptvSubscriberStatistics GetStatisticsUpdateCount
 return "", nil
}

