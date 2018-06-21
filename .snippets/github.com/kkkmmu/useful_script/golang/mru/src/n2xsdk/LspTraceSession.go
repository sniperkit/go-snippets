package n2xsdk

type LspTraceSession struct {
 Handler string
}

func(np *LspTraceSession) Start () error {
 //parameters: SessionHandle
 //AgtLspTraceSession Start
 return nil
}

func(np *LspTraceSession) Reset () error {
 //parameters: SessionHandle
 //AgtLspTraceSession Reset
 return nil
}

func(np *LspTraceSession) Renew () error {
 //parameters: SessionHandle
 //AgtLspTraceSession Renew
 return nil
}

