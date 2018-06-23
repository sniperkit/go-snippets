package n2xsdk

type RoutingEngine struct {
 Handler string
}

func(np *RoutingEngine) Start () error {
 //parameters: 
 //AgtRoutingEngine Start
 return nil
}

func(np *RoutingEngine) Stop () error {
 //parameters: 
 //AgtRoutingEngine Stop
 return nil
}

func(np *RoutingEngine) GetState ()(string, error) {
 //parameters: 
 //AgtRoutingEngine GetState
 return "", nil
}

func(np *RoutingEngine) SetTraceLogLevel () error {
 //parameters: LogLevel
 //AgtRoutingEngine SetTraceLogLevel
 return nil
}

func(np *RoutingEngine) GetTraceLogLevel ()(string, error) {
 //parameters: 
 //AgtRoutingEngine GetTraceLogLevel
 return "", nil
}

func(np *RoutingEngine) SetTraceLogFile () error {
 //parameters: FileName
 //AgtRoutingEngine SetTraceLogFile
 return nil
}

func(np *RoutingEngine) GetTraceLogFile ()(string, error) {
 //parameters: 
 //AgtRoutingEngine GetTraceLogFile
 return "", nil
}

