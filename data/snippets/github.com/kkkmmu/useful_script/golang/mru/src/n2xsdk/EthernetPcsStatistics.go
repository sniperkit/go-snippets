package n2xsdk

type EthernetPcsStatistics struct {
 Handler string
}

func(np *EthernetPcsStatistics) GetType ()(string, error) {
 //parameters: Handle
 //AgtEthernetPcsStatistics GetType
 return "", nil
}

func(np *EthernetPcsStatistics) GetName ()(string, error) {
 //parameters: Handle
 //AgtEthernetPcsStatistics GetName
 return "", nil
}

func(np *EthernetPcsStatistics) GetLockCount ()(string, error) {
 //parameters: Handle
 //AgtEthernetPcsStatistics GetLockCount
 return "", nil
}

func(np *EthernetPcsStatistics) ListSelectedPorts ()(string, error) {
 //parameters: StatisticsHandle
 //AgtEthernetPcsStatistics ListSelectedPorts
 return "", nil
}

func(np *EthernetPcsStatistics) SelectPorts () error {
 //parameters: StatisticsHandle NumPorts PortHandles
 //AgtEthernetPcsStatistics SelectPorts
 return nil
}

func(np *EthernetPcsStatistics) GetAccumulatedValues ()(string, error) {
 //parameters: StatisticsHandle
 //AgtEthernetPcsStatistics GetAccumulatedValues
 return "", nil
}

func(np *EthernetPcsStatistics) LockItem () error {
 //parameters: StatisticsHandle
 //AgtEthernetPcsStatistics LockItem
 return nil
}

func(np *EthernetPcsStatistics) UnlockItem () error {
 //parameters: StatisticsHandle
 //AgtEthernetPcsStatistics UnlockItem
 return nil
}

func(np *EthernetPcsStatistics) ListAvailableStatistics ()(string, error) {
 //parameters: 
 //AgtEthernetPcsStatistics ListAvailableStatistics
 return "", nil
}

func(np *EthernetPcsStatistics) ListSelectedStatistics ()(string, error) {
 //parameters: Handle
 //AgtEthernetPcsStatistics ListSelectedStatistics
 return "", nil
}

func(np *EthernetPcsStatistics) SelectStatistics () error {
 //parameters: Handle NumStatistics Statistics
 //AgtEthernetPcsStatistics SelectStatistics
 return nil
}

