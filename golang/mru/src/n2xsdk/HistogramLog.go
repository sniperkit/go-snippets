package n2xsdk

type HistogramLog struct {
 Handler string
}

func(np *HistogramLo) EnableLogging () error {
 //parameters: 
 //AgtHistogramLog EnableLogging
 return nil
}

func(np *HistogramLo) DisableLogging () error {
 //parameters: 
 //AgtHistogramLog DisableLogging
 return nil
}

func(np *HistogramLo) IsLoggingEnabled () error {
 //parameters: 
 //AgtHistogramLog IsLoggingEnabled
 return nil
}

func(np *HistogramLo) SetLogFile () error {
 //parameters: LogFile
 //AgtHistogramLog SetLogFile
 return nil
}

func(np *HistogramLo) GetLogFile ()(string, error) {
 //parameters: 
 //AgtHistogramLog GetLogFile
 return "", nil
}

func(np *HistogramLo) SetLoggingInterval () error {
 //parameters: Multiple
 //AgtHistogramLog SetLoggingInterval
 return nil
}

func(np *HistogramLo) GetLoggingInterval ()(string, error) {
 //parameters: 
 //AgtHistogramLog GetLoggingInterval
 return "", nil
}

