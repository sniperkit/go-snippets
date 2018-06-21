package n2xsdk

type VsrStatus struct {
 Handler string
}

func(np *VsrStatus) GetStatus ()(string, error) {
 //parameters: PortHandle
 //AgtVsrStatus GetStatus
 return "", nil
}

func(np *VsrStatus) GetStatusDescription ()(string, error) {
 //parameters: PortHandle
 //AgtVsrStatus GetStatusDescription
 return "", nil
}

func(np *VsrStatus) SaveStatus () error {
 //parameters: PortHandle LogFile
 //AgtVsrStatus SaveStatus
 return nil
}

func(np *VsrStatus) GetStatusSummary ()(string, error) {
 //parameters: 
 //AgtVsrStatus GetStatusSummary
 return "", nil
}

func(np *VsrStatus) GetStatusSummaryDescription ()(string, error) {
 //parameters: 
 //AgtVsrStatus GetStatusSummaryDescription
 return "", nil
}

func(np *VsrStatus) SaveStatusSummary () error {
 //parameters: LogFile
 //AgtVsrStatus SaveStatusSummary
 return nil
}

func(np *VsrStatus) GetStatusHistory ()(string, error) {
 //parameters: 
 //AgtVsrStatus GetStatusHistory
 return "", nil
}

func(np *VsrStatus) GetStatusHistoryDescription ()(string, error) {
 //parameters: 
 //AgtVsrStatus GetStatusHistoryDescription
 return "", nil
}

func(np *VsrStatus) SaveStatusHistory () error {
 //parameters: LogFile
 //AgtVsrStatus SaveStatusHistory
 return nil
}

func(np *VsrStatus) ClearHistory () error {
 //parameters: 
 //AgtVsrStatus ClearHistory
 return nil
}

func(np *VsrStatus) SetHistorySize () error {
 //parameters: HistorySize
 //AgtVsrStatus SetHistorySize
 return nil
}

func(np *VsrStatus) GetHistorySize ()(string, error) {
 //parameters: 
 //AgtVsrStatus GetHistorySize
 return "", nil
}

func(np *VsrStatus) GetMaximumHistorySize ()(string, error) {
 //parameters: 
 //AgtVsrStatus GetMaximumHistorySize
 return "", nil
}

