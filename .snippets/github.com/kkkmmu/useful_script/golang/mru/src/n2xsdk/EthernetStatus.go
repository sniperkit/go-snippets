package n2xsdk

type EthernetStatus struct {
 Handler string
}

func(np *EthernetStatus) GetStatus ()(string, error) {
 //parameters: PortHandle
 //AgtEthernetStatus GetStatus
 return "", nil
}

func(np *EthernetStatus) GetStatusDescription ()(string, error) {
 //parameters: PortHandle
 //AgtEthernetStatus GetStatusDescription
 return "", nil
}

func(np *EthernetStatus) SaveStatus () error {
 //parameters: PortHandle LogFile
 //AgtEthernetStatus SaveStatus
 return nil
}

func(np *EthernetStatus) GetStatusSummary ()(string, error) {
 //parameters: 
 //AgtEthernetStatus GetStatusSummary
 return "", nil
}

func(np *EthernetStatus) GetStatusSummaryDescription ()(string, error) {
 //parameters: 
 //AgtEthernetStatus GetStatusSummaryDescription
 return "", nil
}

func(np *EthernetStatus) SaveStatusSummary () error {
 //parameters: LogFile
 //AgtEthernetStatus SaveStatusSummary
 return nil
}

func(np *EthernetStatus) GetStatusHistory ()(string, error) {
 //parameters: 
 //AgtEthernetStatus GetStatusHistory
 return "", nil
}

func(np *EthernetStatus) GetStatusHistoryDescription ()(string, error) {
 //parameters: 
 //AgtEthernetStatus GetStatusHistoryDescription
 return "", nil
}

func(np *EthernetStatus) SaveStatusHistory () error {
 //parameters: LogFile
 //AgtEthernetStatus SaveStatusHistory
 return nil
}

func(np *EthernetStatus) ClearHistory () error {
 //parameters: 
 //AgtEthernetStatus ClearHistory
 return nil
}

func(np *EthernetStatus) SetHistorySize () error {
 //parameters: HistorySize
 //AgtEthernetStatus SetHistorySize
 return nil
}

func(np *EthernetStatus) GetHistorySize ()(string, error) {
 //parameters: 
 //AgtEthernetStatus GetHistorySize
 return "", nil
}

func(np *EthernetStatus) GetMaximumHistorySize ()(string, error) {
 //parameters: 
 //AgtEthernetStatus GetMaximumHistorySize
 return "", nil
}

