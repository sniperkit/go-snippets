package n2xsdk

type IpStatus struct {
 Handler string
}

func(np *IpStatus) GetStatus ()(string, error) {
 //parameters: PortHandle
 //AgtIpStatus GetStatus
 return "", nil
}

func(np *IpStatus) GetStatusDescription ()(string, error) {
 //parameters: PortHandle
 //AgtIpStatus GetStatusDescription
 return "", nil
}

func(np *IpStatus) SaveStatus () error {
 //parameters: PortHandle LogFile
 //AgtIpStatus SaveStatus
 return nil
}

func(np *IpStatus) GetStatusSummary ()(string, error) {
 //parameters: 
 //AgtIpStatus GetStatusSummary
 return "", nil
}

func(np *IpStatus) GetStatusSummaryDescription ()(string, error) {
 //parameters: 
 //AgtIpStatus GetStatusSummaryDescription
 return "", nil
}

func(np *IpStatus) SaveStatusSummary () error {
 //parameters: LogFile
 //AgtIpStatus SaveStatusSummary
 return nil
}

func(np *IpStatus) GetStatusHistory ()(string, error) {
 //parameters: 
 //AgtIpStatus GetStatusHistory
 return "", nil
}

func(np *IpStatus) GetStatusHistoryDescription ()(string, error) {
 //parameters: 
 //AgtIpStatus GetStatusHistoryDescription
 return "", nil
}

func(np *IpStatus) SaveStatusHistory () error {
 //parameters: LogFile
 //AgtIpStatus SaveStatusHistory
 return nil
}

func(np *IpStatus) ClearHistory () error {
 //parameters: 
 //AgtIpStatus ClearHistory
 return nil
}

func(np *IpStatus) SetHistorySize () error {
 //parameters: HistorySize
 //AgtIpStatus SetHistorySize
 return nil
}

func(np *IpStatus) GetHistorySize ()(string, error) {
 //parameters: 
 //AgtIpStatus GetHistorySize
 return "", nil
}

func(np *IpStatus) GetMaximumHistorySize ()(string, error) {
 //parameters: 
 //AgtIpStatus GetMaximumHistorySize
 return "", nil
}

