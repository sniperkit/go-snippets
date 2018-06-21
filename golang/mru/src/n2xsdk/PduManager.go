package n2xsdk

type PduManager struct {
 Handler string
}

func(np *PduManager) GetPduBuilderMessageLog ()(string, error) {
 //parameters: 
 //AgtPduManager GetPduBuilderMessageLog
 return "", nil
}

func(np *PduManager) ClearPduBuilderMessageLog () error {
 //parameters: 
 //AgtPduManager ClearPduBuilderMessageLog
 return nil
}

func(np *PduManager) GetLastError ()(string, error) {
 //parameters: 
 //AgtPduManager GetLastError
 return "", nil
}

