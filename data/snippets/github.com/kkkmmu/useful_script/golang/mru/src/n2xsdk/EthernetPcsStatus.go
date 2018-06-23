package n2xsdk

type EthernetPcsStatus struct {
 Handler string
}

func(np *EthernetPcsStatus) GetStatus ()(string, error) {
 //parameters: PortHandle
 //AgtEthernetPcsStatus GetStatus
 return "", nil
}

func(np *EthernetPcsStatus) GetStatusDescription ()(string, error) {
 //parameters: PortHandle
 //AgtEthernetPcsStatus GetStatusDescription
 return "", nil
}

func(np *EthernetPcsStatus) SaveStatus () error {
 //parameters: PortHandle LogFile
 //AgtEthernetPcsStatus SaveStatus
 return nil
}

func(np *EthernetPcsStatus) GetStatusSummary ()(string, error) {
 //parameters: 
 //AgtEthernetPcsStatus GetStatusSummary
 return "", nil
}

func(np *EthernetPcsStatus) GetStatusSummaryDescription ()(string, error) {
 //parameters: 
 //AgtEthernetPcsStatus GetStatusSummaryDescription
 return "", nil
}

func(np *EthernetPcsStatus) SaveStatusSummary () error {
 //parameters: LogFile
 //AgtEthernetPcsStatus SaveStatusSummary
 return nil
}

func(np *EthernetPcsStatus) GetStatusHistory ()(string, error) {
 //parameters: 
 //AgtEthernetPcsStatus GetStatusHistory
 return "", nil
}

func(np *EthernetPcsStatus) GetStatusHistoryDescription ()(string, error) {
 //parameters: 
 //AgtEthernetPcsStatus GetStatusHistoryDescription
 return "", nil
}

func(np *EthernetPcsStatus) SaveStatusHistory () error {
 //parameters: LogFile
 //AgtEthernetPcsStatus SaveStatusHistory
 return nil
}

func(np *EthernetPcsStatus) ClearHistory () error {
 //parameters: 
 //AgtEthernetPcsStatus ClearHistory
 return nil
}

func(np *EthernetPcsStatus) SetHistorySize () error {
 //parameters: HistorySize
 //AgtEthernetPcsStatus SetHistorySize
 return nil
}

func(np *EthernetPcsStatus) GetHistorySize ()(string, error) {
 //parameters: 
 //AgtEthernetPcsStatus GetHistorySize
 return "", nil
}

func(np *EthernetPcsStatus) GetMaximumHistorySize ()(string, error) {
 //parameters: 
 //AgtEthernetPcsStatus GetMaximumHistorySize
 return "", nil
}

