package n2xsdk

type IpStatistics struct {
 Handler string
}

func(np *IpStatistics) GetType ()(string, error) {
 //parameters: Handle
 //AgtIpStatistics GetType
 return "", nil
}

func(np *IpStatistics) GetName ()(string, error) {
 //parameters: Handle
 //AgtIpStatistics GetName
 return "", nil
}

func(np *IpStatistics) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtIpStatistics GetLockCount
 return "", nil
}

func(np *IpStatistics) ListSelectedPorts ()(string, error) {
 //parameters: StatisticsHandle
 //AgtIpStatistics ListSelectedPorts
 return "", nil
}

func(np *IpStatistics) SelectPorts () error {
 //parameters: StatisticsHandle NumPorts PortHandles
 //AgtIpStatistics SelectPorts
 return nil
}

func(np *IpStatistics) GetAccumulatedValues ()(string, error) {
 //parameters: StatisticsHandle
 //AgtIpStatistics GetAccumulatedValues
 return "", nil
}

func(np *IpStatistics) LockItem () error {
 //parameters: StatisticsHandle
 //AgtIpStatistics LockItem
 return nil
}

func(np *IpStatistics) UnlockItem () error {
 //parameters: StatisticsHandle
 //AgtIpStatistics UnlockItem
 return nil
}

func(np *IpStatistics) ListAvailableStatistics ()(string, error) {
 //parameters: 
 //AgtIpStatistics ListAvailableStatistics
 return "", nil
}

func(np *IpStatistics) ListSelectedStatistics ()(string, error) {
 //parameters: Handle
 //AgtIpStatistics ListSelectedStatistics
 return "", nil
}

func(np *IpStatistics) SelectStatistics () error {
 //parameters: Handle NumStatistics Statistics
 //AgtIpStatistics SelectStatistics
 return nil
}

