package n2xsdk

type SonetStatus struct {
 Handler string
}

func(np *SonetStatus) GetStatus ()(string, error) {
 //parameters: PortHandle
 //AgtSonetStatus GetStatus
 return "", nil
}

func(np *SonetStatus) GetStatusDescription ()(string, error) {
 //parameters: PortHandle
 //AgtSonetStatus GetStatusDescription
 return "", nil
}

func(np *SonetStatus) SaveStatus () error {
 //parameters: PortHandle LogFile
 //AgtSonetStatus SaveStatus
 return nil
}

func(np *SonetStatus) GetStatusSummary ()(string, error) {
 //parameters: 
 //AgtSonetStatus GetStatusSummary
 return "", nil
}

func(np *SonetStatus) GetStatusSummaryDescription ()(string, error) {
 //parameters: 
 //AgtSonetStatus GetStatusSummaryDescription
 return "", nil
}

func(np *SonetStatus) SaveStatusSummary () error {
 //parameters: LogFile
 //AgtSonetStatus SaveStatusSummary
 return nil
}

func(np *SonetStatus) GetStatusHistory ()(string, error) {
 //parameters: 
 //AgtSonetStatus GetStatusHistory
 return "", nil
}

func(np *SonetStatus) GetStatusHistoryDescription ()(string, error) {
 //parameters: 
 //AgtSonetStatus GetStatusHistoryDescription
 return "", nil
}

func(np *SonetStatus) SaveStatusHistory () error {
 //parameters: LogFile
 //AgtSonetStatus SaveStatusHistory
 return nil
}

func(np *SonetStatus) ClearHistory () error {
 //parameters: 
 //AgtSonetStatus ClearHistory
 return nil
}

func(np *SonetStatus) SetHistorySize () error {
 //parameters: HistorySize
 //AgtSonetStatus SetHistorySize
 return nil
}

func(np *SonetStatus) GetHistorySize ()(string, error) {
 //parameters: 
 //AgtSonetStatus GetHistorySize
 return "", nil
}

func(np *SonetStatus) GetMaximumHistorySize ()(string, error) {
 //parameters: 
 //AgtSonetStatus GetMaximumHistorySize
 return "", nil
}

