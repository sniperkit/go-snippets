package n2xsdk

type EthernetPonStatus struct {
 Handler string
}

func(np *EthernetPonStatus) GetStatus ()(string, error) {
 //parameters: PortHandle
 //AgtEthernetPonStatus GetStatus
 return "", nil
}

func(np *EthernetPonStatus) GetStatusDescription ()(string, error) {
 //parameters: PortHandle
 //AgtEthernetPonStatus GetStatusDescription
 return "", nil
}

func(np *EthernetPonStatus) SaveStatus () error {
 //parameters: PortHandle LogFile
 //AgtEthernetPonStatus SaveStatus
 return nil
}

func(np *EthernetPonStatus) GetStatusSummary ()(string, error) {
 //parameters: 
 //AgtEthernetPonStatus GetStatusSummary
 return "", nil
}

func(np *EthernetPonStatus) GetStatusSummaryDescription ()(string, error) {
 //parameters: 
 //AgtEthernetPonStatus GetStatusSummaryDescription
 return "", nil
}

func(np *EthernetPonStatus) SaveStatusSummary () error {
 //parameters: LogFile
 //AgtEthernetPonStatus SaveStatusSummary
 return nil
}

func(np *EthernetPonStatus) GetStatusHistory ()(string, error) {
 //parameters: 
 //AgtEthernetPonStatus GetStatusHistory
 return "", nil
}

func(np *EthernetPonStatus) GetStatusHistoryDescription ()(string, error) {
 //parameters: 
 //AgtEthernetPonStatus GetStatusHistoryDescription
 return "", nil
}

func(np *EthernetPonStatus) SaveStatusHistory () error {
 //parameters: LogFile
 //AgtEthernetPonStatus SaveStatusHistory
 return nil
}

func(np *EthernetPonStatus) ClearHistory () error {
 //parameters: 
 //AgtEthernetPonStatus ClearHistory
 return nil
}

func(np *EthernetPonStatus) SetHistorySize () error {
 //parameters: HistorySize
 //AgtEthernetPonStatus SetHistorySize
 return nil
}

func(np *EthernetPonStatus) GetHistorySize ()(string, error) {
 //parameters: 
 //AgtEthernetPonStatus GetHistorySize
 return "", nil
}

func(np *EthernetPonStatus) GetMaximumHistorySize ()(string, error) {
 //parameters: 
 //AgtEthernetPonStatus GetMaximumHistorySize
 return "", nil
}

