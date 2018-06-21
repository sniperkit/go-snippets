package n2xsdk

type BufferOverflow struct {
 Handler string
}

func(np *BufferOverflow) GetOverflowHistory ()(string, error) {
 //parameters: 
 //AgtBufferOverflow GetOverflowHistory
 return "", nil
}

func(np *BufferOverflow) ClearHistory () error {
 //parameters: 
 //AgtBufferOverflow ClearHistory
 return nil
}

func(np *BufferOverflow) SetHistorySize () error {
 //parameters: HistorySize
 //AgtBufferOverflow SetHistorySize
 return nil
}

func(np *BufferOverflow) GetHistorySize ()(string, error) {
 //parameters: 
 //AgtBufferOverflow GetHistorySize
 return "", nil
}

func(np *BufferOverflow) GetMaximumHistorySize ()(string, error) {
 //parameters: 
 //AgtBufferOverflow GetMaximumHistorySize
 return "", nil
}

