package n2xsdk

type IptvTestScenarioStatistics struct {
 Handler string
}

func(np *IptvTestScenarioStatistics) ListAvailableStatisticsSets ()(string, error) {
 //parameters: ScenarioHandle
 //AgtIptvTestScenarioStatistics ListAvailableStatisticsSets
 return "", nil
}

func(np *IptvTestScenarioStatistics) EnableStatisticsSet () error {
 //parameters: ScenarioHandle StatisticsSet
 //AgtIptvTestScenarioStatistics EnableStatisticsSet
 return nil
}

func(np *IptvTestScenarioStatistics) DisableStatisticsSet () error {
 //parameters: ScenarioHandle StatisticsSet
 //AgtIptvTestScenarioStatistics DisableStatisticsSet
 return nil
}

func(np *IptvTestScenarioStatistics) IsStatisticsSetEnabled () error {
 //parameters: ScenarioHandle StatisticsSet
 //AgtIptvTestScenarioStatistics IsStatisticsSetEnabled
 return nil
}

func(np *IptvTestScenarioStatistics) IsStatisticsEnabled () error {
 //parameters: ScenarioHandle
 //AgtIptvTestScenarioStatistics IsStatisticsEnabled
 return nil
}

func(np *IptvTestScenarioStatistics) ListAvailableThresholdTypes ()(string, error) {
 //parameters: ScenarioHandle
 //AgtIptvTestScenarioStatistics ListAvailableThresholdTypes
 return "", nil
}

func(np *IptvTestScenarioStatistics) SetThreshold () error {
 //parameters: ScenarioHandle ThresholdType ThresholdValue
 //AgtIptvTestScenarioStatistics SetThreshold
 return nil
}

func(np *IptvTestScenarioStatistics) GetThreshold ()(string, error) {
 //parameters: ScenarioHandle ThresholdType
 //AgtIptvTestScenarioStatistics GetThreshold
 return "", nil
}

func(np *IptvTestScenarioStatistics) ListAvailableThresholdByChannelTypes ()(string, error) {
 //parameters: ScenarioHandle
 //AgtIptvTestScenarioStatistics ListAvailableThresholdByChannelTypes
 return "", nil
}

func(np *IptvTestScenarioStatistics) SetThresholdByChannel () error {
 //parameters: ScenarioHandle GroupPoolHandle Threshold ThresholdValue
 //AgtIptvTestScenarioStatistics SetThresholdByChannel
 return nil
}

func(np *IptvTestScenarioStatistics) GetThresholdByChannel ()(string, error) {
 //parameters: ScenarioHandle GroupPoolHandle Threshold
 //AgtIptvTestScenarioStatistics GetThresholdByChannel
 return "", nil
}

func(np *IptvTestScenarioStatistics) SetMaxNumberOfReportedViolations () error {
 //parameters: ScenarioHandle ViolationReportSize
 //AgtIptvTestScenarioStatistics SetMaxNumberOfReportedViolations
 return nil
}

func(np *IptvTestScenarioStatistics) GetMaxNumberOfReportedViolations ()(string, error) {
 //parameters: ScenarioHandle
 //AgtIptvTestScenarioStatistics GetMaxNumberOfReportedViolations
 return "", nil
}

func(np *IptvTestScenarioStatistics) ResetStatistics () error {
 //parameters: ScenarioHandle
 //AgtIptvTestScenarioStatistics ResetStatistics
 return nil
}

func(np *IptvTestScenarioStatistics) ListAvailableStatisticsTypes ()(string, error) {
 //parameters: ScenarioHandle
 //AgtIptvTestScenarioStatistics ListAvailableStatisticsTypes
 return "", nil
}

func(np *IptvTestScenarioStatistics) GetStatistics ()(string, error) {
 //parameters: ScenarioHandle StatisticsList
 //AgtIptvTestScenarioStatistics GetStatistics
 return "", nil
}

func(np *IptvTestScenarioStatistics) GetSubscriberPerChannelStatistics ()(string, error) {
 //parameters: ScenarioInstance StatisticsList
 //AgtIptvTestScenarioStatistics GetSubscriberPerChannelStatistics
 return "", nil
}

func(np *IptvTestScenarioStatistics) GetNumberOfThresholdViolations ()(string, error) {
 //parameters: ScenarioInstance Threshold
 //AgtIptvTestScenarioStatistics GetNumberOfThresholdViolations
 return "", nil
}

func(np *IptvTestScenarioStatistics) GetThresholdViolationReports ()(string, error) {
 //parameters: ScenarioInstance Threshold
 //AgtIptvTestScenarioStatistics GetThresholdViolationReports
 return "", nil
}

func(np *IptvTestScenarioStatistics) GetTestStartTime ()(string, error) {
 //parameters: ScenarioHandle
 //AgtIptvTestScenarioStatistics GetTestStartTime
 return "", nil
}

