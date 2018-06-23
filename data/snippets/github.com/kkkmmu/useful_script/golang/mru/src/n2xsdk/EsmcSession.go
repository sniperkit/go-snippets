package n2xsdk

type EsmcSession struct {
 Handler string
}

func(np *EsmcSession) SetClockQuality () error {
 //parameters: SessionHandle SsmQuality
 //AgtEsmcSession SetClockQuality
 return nil
}

func(np *EsmcSession) GetClockQuality ()(string, error) {
 //parameters: SessionHandle
 //AgtEsmcSession GetClockQuality
 return "", nil
}

func(np *EsmcSession) SetMessageRate () error {
 //parameters: SessionHandle MessageRate
 //AgtEsmcSession SetMessageRate
 return nil
}

func(np *EsmcSession) GetMessageRate ()(string, error) {
 //parameters: SessionHandle
 //AgtEsmcSession GetMessageRate
 return "", nil
}

func(np *EsmcSession) SetEventMode () error {
 //parameters: SessionHandle EventMode
 //AgtEsmcSession SetEventMode
 return nil
}

func(np *EsmcSession) GetEventMode ()(string, error) {
 //parameters: SessionHandle
 //AgtEsmcSession GetEventMode
 return "", nil
}

