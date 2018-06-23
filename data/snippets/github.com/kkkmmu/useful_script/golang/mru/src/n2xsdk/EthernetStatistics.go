package n2xsdk

type EthernetStatistics struct {
 Handler string
}

func(np *EthernetStatistics) GetType ()(string, error) {
 //parameters: Handle
 //AgtEthernetStatistics GetType
 return "", nil
}

func(np *EthernetStatistics) GetName ()(string, error) {
 //parameters: Handle
 //AgtEthernetStatistics GetName
 return "", nil
}

func(np *EthernetStatistics) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtEthernetStatistics GetLockCount
 return "", nil
}

func(np *EthernetStatistics) ListSelectedPorts ()(string, error) {
 //parameters: StatisticsHandle
 //AgtEthernetStatistics ListSelectedPorts
 return "", nil
}

func(np *EthernetStatistics) SelectPorts () error {
 //parameters: StatisticsHandle NumPorts PortHandles
 //AgtEthernetStatistics SelectPorts
 return nil
}

func(np *EthernetStatistics) GetAccumulatedValues ()(string, error) {
 //parameters: StatisticsHandle
 //AgtEthernetStatistics GetAccumulatedValues
 return "", nil
}

func(np *EthernetStatistics) LockItem () error {
 //parameters: StatisticsHandle
 //AgtEthernetStatistics LockItem
 return nil
}

func(np *EthernetStatistics) UnlockItem () error {
 //parameters: StatisticsHandle
 //AgtEthernetStatistics UnlockItem
 return nil
}

func(np *EthernetStatistics) ListAvailableStatistics ()(string, error) {
 //parameters: 
 //AgtEthernetStatistics ListAvailableStatistics
 return "", nil
}

func(np *EthernetStatistics) ListSelectedStatistics ()(string, error) {
 //parameters: Handle
 //AgtEthernetStatistics ListSelectedStatistics
 return "", nil
}

func(np *EthernetStatistics) SelectStatistics () error {
 //parameters: Handle NumStatistics Statistics
 //AgtEthernetStatistics SelectStatistics
 return nil
}

