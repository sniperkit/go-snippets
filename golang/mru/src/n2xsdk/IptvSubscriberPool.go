package n2xsdk

type IptvSubscriberPool struct {
 Handler string
}

func(np *IptvSubscriberPool) SetNetworkTopology () error {
 //parameters: DeviceHandle NetworkTopology
 //AgtIptvSubscriberPool SetNetworkTopology
 return nil
}

func(np *IptvSubscriberPool) GetNetworkTopology ()(string, error) {
 //parameters: DeviceHandle
 //AgtIptvSubscriberPool GetNetworkTopology
 return "", nil
}

func(np *IptvSubscriberPool) SetVideoVlanId () error {
 //parameters: DeviceHandle VlanId
 //AgtIptvSubscriberPool SetVideoVlanId
 return nil
}

func(np *IptvSubscriberPool) GetVideoVlanId ()(string, error) {
 //parameters: DeviceHandle
 //AgtIptvSubscriberPool GetVideoVlanId
 return "", nil
}

func(np *IptvSubscriberPool) AddTestScenario () error {
 //parameters: DeviceHandle ScenarioType
 //AgtIptvSubscriberPool AddTestScenario
 return nil
}

func(np *IptvSubscriberPool) RemoveTestScenario () error {
 //parameters: DeviceHandle ScenarioHandle
 //AgtIptvSubscriberPool RemoveTestScenario
 return nil
}

func(np *IptvSubscriberPool) RemoveAllTestScenarios () error {
 //parameters: DeviceHandle
 //AgtIptvSubscriberPool RemoveAllTestScenarios
 return nil
}

func(np *IptvSubscriberPool) ListTestScenarios ()(string, error) {
 //parameters: DeviceHandle
 //AgtIptvSubscriberPool ListTestScenarios
 return "", nil
}

func(np *IptvSubscriberPool) ListAvailableTestScenarioTypes ()(string, error) {
 //parameters: DeviceHandle
 //AgtIptvSubscriberPool ListAvailableTestScenarioTypes
 return "", nil
}

func(np *IptvSubscriberPool) GetTestScenarioType ()(string, error) {
 //parameters: DeviceHandle ScenarioHandle
 //AgtIptvSubscriberPool GetTestScenarioType
 return "", nil
}

func(np *IptvSubscriberPool) SetTestScenarioName () error {
 //parameters: DeviceHandle ScenarioHandle ScenarioName
 //AgtIptvSubscriberPool SetTestScenarioName
 return nil
}

func(np *IptvSubscriberPool) GetTestScenarioName ()(string, error) {
 //parameters: DeviceHandle ScenarioHandle
 //AgtIptvSubscriberPool GetTestScenarioName
 return "", nil
}

